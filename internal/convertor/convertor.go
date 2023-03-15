package convertor

import (
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post"
	"github.com/xh-polaris/meowchat-post-rpc/pb"
)

func ConvertPost(in *post.Post) *pb.Post {
	return &pb.Post{
		Id:         in.ID.Hex(),
		CreateAt:   in.CreateAt.Unix(),
		UpdateAt:   in.UpdateAt.Unix(),
		Title:      in.Title,
		Text:       in.Text,
		CoverUrl:   in.CoverUrl,
		Tags:       in.Tags,
		UserId:     in.UserId,
		IsOfficial: in.Flags.GetFlag(post.OfficialFlag),
	}
}
