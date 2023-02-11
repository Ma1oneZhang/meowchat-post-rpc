package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-post-rpc/internal/common"

	"github.com/xh-polaris/meowchat-post-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-post-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostByUserIdLogic {
	return &ListPostByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPostByUserIdLogic) ListPostByUserId(in *pb.ListPostByUserIdReq) (*pb.ListPostByUserIdResp, error) {
	data, total, err := l.svcCtx.PostModel.FindManyByUserId(l.ctx, in.UserId, in.Skip, in.Count)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Post, 0, len(data))
	for _, val := range data {
		res = append(res, common.PostTransform(val))
	}
	return &pb.ListPostByUserIdResp{Posts: res, Total: total}, nil
}
