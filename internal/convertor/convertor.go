package convertor

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model"
	"github.com/xh-polaris/meowchat-post-rpc/pb"
)

func ConvertPost(in *model.Post) *pb.Post {
	return &pb.Post{
		Id:         in.ID.Hex(),
		CreateAt:   in.CreateAt.Unix(),
		UpdateAt:   in.UpdateAt.Unix(),
		Title:      in.Title,
		Text:       in.Text,
		CoverUrl:   in.CoverUrl,
		Tags:       in.Tags,
		UserId:     in.UserId,
		IsOfficial: in.Flags.GetOfficial(),
	}
}

func ConvertAllFieldsSearchQuery(in *pb.ListPostReq_AllFieldsKey) []types.Query {
	return []types.Query{{
		MultiMatch: &types.MultiMatchQuery{
			Query:  in.AllFieldsKey,
			Fields: []string{model.Title + "^3", model.Text, model.Tags},
		}},
	}
}

func ConvertMultiFieldsSearchQuery(in *pb.ListPostReq_MultiFieldsKey) []types.Query {
	var q []types.Query
	if in.MultiFieldsKey.Title != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				model.Title: {
					Query: *in.MultiFieldsKey.Title,
				},
			},
		})
	}
	if in.MultiFieldsKey.Text != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				model.Text: {
					Query: *in.MultiFieldsKey.Text,
				},
			},
		})
	}
	if in.MultiFieldsKey.Tag != nil {
		q = append(q, types.Query{
			Match: map[string]types.MatchQuery{
				model.Tags: {
					Query: *in.MultiFieldsKey.Tag,
				},
			},
		})
	}
	return q
}
