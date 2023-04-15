package mongo

import (
	"context"
	"sync"

	"github.com/xh-polaris/meowchat-post-rpc/internal/model"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"
	"github.com/xh-polaris/paginator-go"
	"github.com/xh-polaris/paginator-go/mongop"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"github.com/zeromicro/go-zero/core/syncx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PostCollectionName = "post"
const prefixPostPaginatorKey = "cache:paginator:post:"

var _ PostMongoModel = (*customPostModel)(nil)

type (
	// PostMongoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPostModel.
	PostMongoModel interface {
		postModel
		FindMany(ctx context.Context, fopts *internal.FilterOptions, popts *paginator.PaginationOptions, sorter int64) ([]*internal.Post, error)
		Count(ctx context.Context, fopts *internal.FilterOptions) (int64, error)
		FindManyAndCount(ctx context.Context, fopts *internal.FilterOptions, popts *paginator.PaginationOptions, sorter int64) ([]*internal.Post, int64, error)
		UpdateFlags(ctx context.Context, id string, flags map[internal.PostFlag]bool) error
	}

	customPostModel struct {
		*defaultPostModel
		paginatorCache cache.Cache
	}
)

// NewPostModel returns a model for the mongo.
func NewPostModel(url, db string, c cache.CacheConf) PostMongoModel {
	conn := monc.MustNewModel(url, db, PostCollectionName, c)

	return &customPostModel{
		defaultPostModel: newDefaultPostModel(conn),
		paginatorCache:   cache.New(c, syncx.NewSingleFlight(), cache.NewStat("paginator-mongo"), model.ErrPaginatorTokenExpired),
	}
}

func (m *customPostModel) UpdateFlags(ctx context.Context, id string, flags map[internal.PostFlag]bool) error {
	var or, and internal.PostFlag
	for flag, v := range flags {
		if v {
			or += flag
		} else {
			and += flag
		}
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.ErrInvalidObjectId
	}
	_, err = m.conn.UpdateOne(ctx, prefixPostCacheKey+id, bson.M{internal.ID: oid}, bson.M{
		"$bit": bson.M{
			internal.Flags: bson.M{
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

func (m *customPostModel) FindMany(ctx context.Context, fopts *internal.FilterOptions, popts *paginator.PaginationOptions, sorter int64) ([]*internal.Post, error) {
	p := mongop.NewMongoPaginator(paginator.NewCacheStore(m.paginatorCache, Sorters[sorter], prefixPostPaginatorKey), popts)

	filter := MakeBsonFilter(fopts)
	sort, err := p.MakeSortOptions(ctx, filter)
	if err != nil {
		return nil, err
	}

	var data []*internal.Post
	if err := m.conn.Find(ctx, &data, filter, &options.FindOptions{
		Sort:  sort,
		Limit: popts.Limit,
		Skip:  popts.Offset,
	}); err != nil {
		return nil, err
	}

	// 如果是反向查询，反转数据
	if *popts.Backward {
		for i := 0; i < len(data)/2; i++ {
			data[i], data[len(data)-i-1] = data[len(data)-i-1], data[i]
		}
	}
	if len(data) > 0 {
		err = p.StoreSorter(ctx, data[0], data[len(data)-1])
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (m *customPostModel) Count(ctx context.Context, filter *internal.FilterOptions) (int64, error) {
	f := MakeBsonFilter(filter)
	return m.conn.CountDocuments(ctx, f)
}

func (m *customPostModel) FindManyAndCount(ctx context.Context, fopts *internal.FilterOptions, popts *paginator.PaginationOptions, sorter int64) ([]*internal.Post, int64, error) {
	var posts []*internal.Post
	var total int64
	wg := sync.WaitGroup{}
	wg.Add(2)
	c := make(chan error)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		defer wg.Done()
		var err error
		posts, err = m.FindMany(ctx, fopts, popts, sorter)
		if err != nil {
			c <- err
			return
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		total, err = m.Count(ctx, fopts)
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
