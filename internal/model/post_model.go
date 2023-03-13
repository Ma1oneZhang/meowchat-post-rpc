package model

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/mitchellh/mapstructure"
	"github.com/xh-polaris/meowchat-post-rpc/internal/config"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/paginator"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PostCollectionName = "post"

var _ PostModel = (*customPostModel)(nil)

type (
	// PostModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPostModel.
	PostModel interface {
		postModel
		FindMany(ctx context.Context, filter MongoFilter, paginator paginator.IdPaginator) ([]*Post, error)
		Count(ctx context.Context, filter MongoFilter) (int64, error)
		FindManyAndCount(ctx context.Context, filter MongoFilter, paginator paginator.IdPaginator) ([]*Post, int64, error)
		Search(ctx context.Context, query []types.Query, filter EsFilter, paginator paginator.BasePaginator) ([]*Post, int64, error)
		UpdateFlags(ctx context.Context, id string, flags map[PostFlag]bool) error
	}

	customPostModel struct {
		*defaultPostModel
		es        *elasticsearch.TypedClient
		indexName string
	}
)

// NewPostModel returns a model for the mongo.
func NewPostModel(url, db string, c cache.CacheConf, es config.ElasticsearchConf) PostModel {
	conn := monc.MustNewModel(url, db, PostCollectionName, c)
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
	return &customPostModel{
		defaultPostModel: newDefaultPostModel(conn),
		es:               esClient,
		indexName:        fmt.Sprintf("%s.%s-alias", db, PostCollectionName),
	}
}

func (m *customPostModel) UpdateFlags(ctx context.Context, id string, flags map[PostFlag]bool) error {
	var or, and PostFlag
	for flag, v := range flags {
		if v {
			or += flag
		} else {
			and += flag
		}
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidObjectId
	}
	_, err = m.conn.UpdateOne(ctx, prefixPostCacheKey, bson.M{ID: oid}, bson.M{
		"$bit": bson.M{
			Flags: bson.M{
				"and": ^and,
				"or":  or,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *customPostModel) FindMany(ctx context.Context, filter MongoFilter, paginator paginator.IdPaginator) ([]*Post, error) {
	f := filter.toBson()
	opts, err := paginator.GenQuery(f)
	if err != nil {
		return nil, err
	}

	var data []*Post
	if err := m.conn.Find(ctx, &data, f, opts); err != nil {
		return nil, err
	}
	return data, nil
}

func (m *customPostModel) Count(ctx context.Context, filter MongoFilter) (int64, error) {
	return m.conn.CountDocuments(ctx, filter.toBson())
}

func (m *customPostModel) FindManyAndCount(ctx context.Context, filter MongoFilter, paginator paginator.IdPaginator) ([]*Post, int64, error) {
	var posts []*Post
	var total int64
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := make(chan error)
	go func() {
		defer wg.Done()
		var err error
		posts, err = m.FindMany(ctx, filter, paginator)
		if err != nil {
			c <- err
			return
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		total, err = m.Count(ctx, filter)
		if err != nil {
			c <- err
			return
		}
	}()
	go func() {
		wg.Wait()
		defer close(c)
	}()
	if err := <-c; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (m *customPostModel) Search(ctx context.Context, query []types.Query, filter EsFilter, p paginator.BasePaginator) ([]*Post, int64, error) {
	s := m.es.Search().From(int(*p.Offset)).Size(int(*p.Limit)).Index(m.indexName)
	res, err := s.Request(&search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must:   query,
				Filter: filter.toEsQuery(),
			},
		},
		Sort: []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"_score": {Order: &sortorder.Desc},
					CreateAt: {Order: &sortorder.Desc},
				},
			},
		},
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
	posts := make([]*Post, 0, len(hits))
	for i := range hits {
		hit := hits[i].(map[string]any)
		post := &Post{}
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
		posts = append(posts, post)
	}
	return posts, total, nil
}
