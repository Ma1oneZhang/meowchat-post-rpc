package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post"
	"github.com/xh-polaris/meowchat-post-rpc/internal/scheduled"
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
	p := &post.Post{
		Title:    in.Title,
		Text:     in.Text,
		CoverUrl: in.CoverUrl,
		Tags:     in.Tags,
		UserId:   in.UserId,
	}
	err := l.svcCtx.PostModel.Insert(l.ctx, p)
	if err != nil {
		return nil, err
	}
	go scheduled.SendUrlUsedMessageToSts(&l.svcCtx.Config, &[]string{in.CoverUrl})
	return &pb.CreatePostResp{PostId: p.ID.Hex()}, nil
}
