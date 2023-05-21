package es

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post/internal"
	"github.com/xh-polaris/meowchat-post-rpc/pb"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func ConvertAllFieldsSearchQuery(in *pb.SearchOptions_AllFieldsKey) []types.Query {
	return []types.Query{{
		MultiMatch: &types.MultiMatchQuery{
			Query:  in.AllFieldsKey,
			Fields: []string{internal.Title + "^3", internal.Text, internal.Tags},
		}},
	}
}

func ConvertMultiFieldsSearchQuery(in *pb.SearchOptions_MultiFieldsKey) []types.Query {
	var q []types.Query
	if in.MultiFieldsKey.Title != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				internal.Title: {
					Query: *in.MultiFieldsKey.Title + "^3",
				},
			},
		})
	}
	if in.MultiFieldsKey.Text != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				internal.Text: {
					Query: *in.MultiFieldsKey.Text,
				},
			},
		})
	}
	if in.MultiFieldsKey.Tag != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				internal.Tags: {
					Query: *in.MultiFieldsKey.Tag,
				},
			},
		})
	}
	return q
}
