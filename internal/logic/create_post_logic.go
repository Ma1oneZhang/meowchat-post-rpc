package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model"
	"github.com/xh-polaris/meowchat-post-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-post-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(in *pb.CreatePostReq) (*pb.CreatePostResp, error) {
	post := &model.Post{
		Title:    in.Title,
		Text:     in.Text,
		CoverUrl: in.CoverUrl,
		Tags:     in.Tags,
		UserId:   in.UserId,
		Status:   1,
	}
	err := l.svcCtx.PostModel.Insert(l.ctx, post)
	if err != nil {
		return nil, err
	}
	return &pb.CreatePostResp{PostId: post.ID.Hex()}, nil
}
