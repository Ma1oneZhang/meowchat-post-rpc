package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post"
	"github.com/xh-polaris/meowchat-post-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-post-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountPostLogic {
	return &CountPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CountPostLogic) CountPost(in *pb.CountPostReq) (*pb.CountPostResp, error) {
	var total int64
	var err error

	filter := &post.FilterOptions{
		OnlyUserId:   in.FilterOptions.OnlyUserId,
		OnlyOfficial: in.FilterOptions.OnlyOfficial,
	}

	if in.SearchOptions == nil {
		total, err = l.svcCtx.PostModel.Count(l.ctx, filter)
		if err != nil {
			return nil, err
		}
	} else {
		switch o := in.SearchOptions.Query.(type) {
		case *pb.SearchOptions_AllFieldsKey:
			total, err = l.svcCtx.PostModel.CountWithQuery(l.ctx, post.ConvertAllFieldsSearchQuery(o), filter)
		case *pb.SearchOptions_MultiFieldsKey:
			total, err = l.svcCtx.PostModel.CountWithQuery(l.ctx, post.ConvertMultiFieldsSearchQuery(o), filter)
		}
		if err != nil {
			return nil, err
		}
	}

	return &pb.CountPostResp{Total: total}, nil
}
