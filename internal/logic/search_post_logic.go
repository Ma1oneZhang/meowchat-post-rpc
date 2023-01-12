package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-post-rpc/internal/common"

	"github.com/xh-polaris/meowchat-post-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-post-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchPostLogic {
	return &SearchPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchPostLogic) SearchPost(in *pb.SearchPostReq) (*pb.SearchPostResp, error) {
	data, total, err := l.svcCtx.PostModel.Search(l.ctx, in.Keyword, in.Count, in.Skip)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Post, 0, len(data))
	for _, val := range data {
		res = append(res, common.PostTransform(val))
	}
	return &pb.SearchPostResp{Posts: res, Total: total}, nil
}
