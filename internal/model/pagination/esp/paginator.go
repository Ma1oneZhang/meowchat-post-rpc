package esp

import (
	"context"

	"github.com/xh-polaris/meowchat-post-rpc/internal/model/pagination"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type (
	EsPaginator struct {
		sorter *pagination.CachePaginator
		opts   *pagination.PaginationOptions
	}
)

func NewEsPaginator(c cache.Cache, sorter EsSorter, opts *pagination.PaginationOptions) *EsPaginator {
	opts.EnsureSafe()
	return &EsPaginator{
		sorter: pagination.NewCachePaginator(c, sorter),
		opts:   opts,
	}
}

// MakeSortOptions 生成ID分页查询选项
func (p *EsPaginator) MakeSortOptions(ctx context.Context, prefix string) ([]types.SortCombinations, []types.FieldValue, error) {
	if p.opts.LastToken != nil {
		err := p.sorter.LoadSorter(ctx, prefix+*p.opts.LastToken, *p.opts.Backward)
		if err != nil {
			return nil, nil, err
		}
	}

	sorter := p.sorter.GetSorter()
	sort, sa, err := sorter.(EsSorter).MakeSortOptions(*p.opts.Backward)
	if err != nil {
		return nil, nil, err
	}
	return sort, sa, nil
}

func (p *EsPaginator) StoreSorter(ctx context.Context, prefix string, first, last any) error {
	token, err := p.sorter.StoreSorter(ctx, prefix, p.opts.LastToken, first, last)
	p.opts.LastToken = token
	return err
}
