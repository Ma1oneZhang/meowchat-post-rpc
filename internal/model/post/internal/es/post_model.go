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
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"
	"github.com/xh-polaris/paginator-go"
	"github.com/xh-polaris/paginator-go/esp"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PostCollectionName = "post"
const prefixPostPaginatorKey = "cache:paginator:post:"

type (
	PostEsModel interface {
		Search(ctx context.Context, query []types.Query, fopts *internal.FilterOptions, popts *paginator.PaginationOptions, sorter int64) ([]*internal.Post, int64, error)
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

func (m *defaultPostModel) Search(ctx context.Context, query []types.Query, fopts *internal.FilterOptions, popts *paginator.PaginationOptions, sorter int64) ([]*internal.Post, int64, error) {
	p := esp.NewEsPaginator(paginator.NewCacheStore(m.paginatorCache, Sorters[sorter], prefixPostPaginatorKey), popts)
	filter := newPostFilter(fopts)
	s, sa, err := p.MakeSortOptions(ctx)
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

	hits := res.Hits.Hits
	total := res.Hits.Total.Value
	posts := make([]*internal.Post, 0, len(hits))
	for i := range hits {
		hit := hits[i]
		post := &internal.Post{}
		source := make(map[string]any)
		err = json.Unmarshal(hit.Source_, &source)
		if err != nil {
			return nil, 0, err
		}
		if source["createAt"], err = time.Parse("2006-01-02T15:04:05Z07:00", source["createAt"].(string)); err != nil {
			return nil, 0, err
		}
		if source["updateAt"], err = time.Parse("2006-01-02T15:04:05Z07:00", source["updateAt"].(string)); err != nil {
			return nil, 0, err
		}
		err = mapstructure.Decode(source, post)
		if err != nil {
			return nil, 0, err
		}

		oid := hit.Id_
		post.ID, err = primitive.ObjectIDFromHex(oid)
		if err != nil {
			return nil, 0, err
		}
		post.Score_ = float64(hit.Score_)
		posts = append(posts, post)
	}
	// 如果是反向查询，反转数据
	if *popts.Backward {
		for i := 0; i < len(posts)/2; i++ {
			posts[i], posts[len(posts)-i-1] = posts[len(posts)-i-1], posts[i]
		}
	}
	if len(posts) > 0 {
		err = p.StoreSorter(ctx, posts[0], posts[len(posts)-1])
		if err != nil {
			return nil, 0, err
		}
	}
	return posts, total, nil
}
