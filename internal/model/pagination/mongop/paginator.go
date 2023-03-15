package mongop

import (
	"context"

	"github.com/xh-polaris/meowchat-post-rpc/internal/model/pagination"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	MongoPaginator struct {
		opts   *pagination.PaginationOptions
		sorter *pagination.CachePaginator
	}
)

func NewMongoPaginator(c cache.Cache, sorter MongoSorter, opts *pagination.PaginationOptions) *MongoPaginator {
	opts.EnsureSafe()
	return &MongoPaginator{
		sorter: pagination.NewCachePaginator(c, sorter),
		opts:   opts,
	}
}

// MakeSortOptions 生成ID分页查询选项，并将filter在原地更新
func (p *MongoPaginator) MakeSortOptions(ctx context.Context, prefix string, filter bson.M) (bson.M, error) {
	if p.opts.LastToken != nil {
		err := p.sorter.LoadSorter(ctx, prefix+*p.opts.LastToken, *p.opts.Backward)
		if err != nil {
			return nil, err
		}
	}

	sorter := p.sorter.GetSorter()
	sort, err := sorter.(MongoSorter).MakeSortOptions(filter, *p.opts.Backward)
	if err != nil {
		return nil, err
	}
	return sort, nil
}

func (p *MongoPaginator) StoreSorter(ctx context.Context, prefix string, first, last any) error {
	token, err := p.sorter.StoreSorter(ctx, prefix, p.opts.LastToken, first, last)
	p.opts.LastToken = token
	return err
}
