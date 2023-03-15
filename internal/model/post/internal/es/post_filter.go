package es

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"
)

// 检验是否完成Filter接口
var _ internal.Filter = (*postFilter)(nil)

// postFilter is filter used for Elasticsearch
type postFilter struct {
	q []types.Query
	*internal.BaseFilter
}

func newPostFilter(options *internal.FilterOptions) []types.Query {
	return (&postFilter{
		q: make([]types.Query, 0),
		BaseFilter: &internal.BaseFilter{
			FilterOptions: options,
		},
	}).toEsQuery()
}

func (f *postFilter) toEsQuery() []types.Query {
	f.CheckOnlyUserId()
	f.CheckOnlyOfficial()
	f.CheckFlags()
	return f.q
}

func (f *postFilter) CheckFlags() {
	if f.MustFlags != nil && *f.MustFlags != 0 {
		f.q = append(f.q, types.Query{
			//TODO 也许会造成潜在的性能风险
			Script: &types.ScriptQuery{
				Script: types.InlineScript{
					Source: fmt.Sprintf("doc['%s'].size() != 0 && "+
						"(doc['%s'].value & params.%s) == params.%s", internal.Flags, internal.Flags, internal.Flags, internal.Flags),
					Params: map[string]any{
						internal.Flags: *f.MustFlags,
					},
				},
			},
		})
	}
	if f.MustNotFlags != nil {
		f.q = append(f.q, types.Query{
			//TODO 也许会造成潜在的性能风险
			Script: &types.ScriptQuery{
				Script: types.InlineScript{
					Source: fmt.Sprintf("doc['%s'].size() == 0 || "+
						"(doc['%s'].value & params.%s) == 0", internal.Flags, internal.Flags, internal.Flags),
					Params: map[string]any{
						internal.Flags: *f.MustFlags,
					},
				},
			},
		})
	}
}

func (f *postFilter) CheckOnlyUserId() {
	if f.OnlyUserId != nil {
		f.q = append(f.q, types.Query{
			Term: map[string]types.TermQuery{
				internal.UserId: {Value: *f.OnlyUserId},
			},
		})
	}
}
