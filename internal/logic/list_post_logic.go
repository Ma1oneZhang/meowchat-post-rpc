package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-post-rpc/internal/common"

	"github.com/xh-polaris/meowchat-post-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-post-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostLogic {
	return &ListPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPostLogic) ListPost(in *pb.ListPostReq) (*pb.ListPostResp, error) {
	data, total, err := l.svcCtx.PostModel.FindMany(l.ctx, in.Skip, in.Count)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Post, 0, len(data))
	for _, val := range data {
		res = append(res, common.PostTransform(val))
	}
	return &pb.ListPostResp{Posts: res, Total: total}, nil
}
