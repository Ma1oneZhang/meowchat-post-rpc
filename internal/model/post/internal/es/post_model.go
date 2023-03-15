package es

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xh-polaris/meowchat-post-rpc/internal/config"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/pagination"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/pagination/esp"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PostCollectionName = "post"
const prefixPostPaginatorKey = "cache:paginator:post:"

type (
	PostEsModel interface {
		Search(ctx context.Context, query []types.Query, fopts *internal.FilterOptions, popts *pagination.PaginationOptions, sorter int64) ([]*internal.Post, int64, error)
	}

	defaultPostModel struct {
		es             *elasticsearch.TypedClient
		indexName      string
		paginatorCache cache.Cache
	}
)

// NewPostModel returns a model for the elasticsearch.
func NewPostModel(db string, es config.ElasticsearchConf, c cache.CacheConf) PostEsModel {
	esClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: es.Addresses,
		Username:  es.Username,
		Password:  es.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return &defaultPostModel{
		es:             esClient,
		indexName:      fmt.Sprintf("%s.%s-alias", db, PostCollectionName),
		paginatorCache: cache.New(c, syncx.NewSingleFlight(), cache.NewStat("paginator-es"), model.ErrPaginatorTokenExpired),
	}
}

func (m *defaultPostModel) Search(ctx context.Context, query []types.Query, fopts *internal.FilterOptions, popts *pagination.PaginationOptions, sorter int64) ([]*internal.Post, int64, error) {
	p := esp.NewEsPaginator(m.paginatorCache, Sorters[sorter], popts)
	filter := newPostFilter(fopts)
	s, sa, err := p.MakeSortOptions(ctx, prefixPostPaginatorKey)
	if err != nil {
		return nil, 0, err
	}
	res, err := m.es.Search().From(int(*popts.Offset)).Size(int(*popts.Limit)).Index(m.indexName).Request(&search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must:   query,
				Filter: filter,
			},
		},
		SearchAfter: sa,
		Sort:        s,
	}).Do(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, 0, err
		} else {
			return nil, 0, errors.Errorf("[%s] %s: %s",
				res.Status,
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]any
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	hits := r["hits"].(map[string]any)["hits"].([]any)
	total := int64(r["hits"].(map[string]any)["total"].(map[string]any)["value"].(float64))
	posts := make([]*internal.Post, 0, len(hits))
	for i := range hits {
		hit := hits[i].(map[string]any)
		post := &internal.Post{}
		source := hit["_source"].(map[string]any)
		if source["createAt"], err = time.Parse("2006-01-02T15:04:05Z07:00", source["createAt"].(string)); err != nil {
			return nil, 0, err
		}
		if source["updateAt"], err = time.Parse("2006-01-02T15:04:05Z07:00", source["updateAt"].(string)); err != nil {
			return nil, 0, err
		}
		hit["_source"] = source
		err := mapstructure.Decode(hit["_source"], post)
		if err != nil {
			return nil, 0, err
		}
		oid := hit["_id"].(string)
		id, err := primitive.ObjectIDFromHex(oid)
		if err != nil {
			return nil, 0, err
		}
		post.ID = id
		post.Score_ = hit["_score"].(float64)
		posts = append(posts, post)
	}
	// 如果是反向查询，反转数据
	if *popts.Backward {
		for i := 0; i < len(posts)/2; i++ {
			posts[i], posts[len(posts)-i-1] = posts[len(posts)-i-1], posts[i]
		}
	}
	if len(posts) > 0 {
		err = p.StoreSorter(ctx, prefixPostPaginatorKey, posts[0], posts[len(posts)-1])
		if err != nil {
			return nil, 0, err
		}
	}
	return posts, total, nil
}
