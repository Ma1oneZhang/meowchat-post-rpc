package logic

import (
	"context"

	"github.com/xh-polaris/meowchat-post-rpc/errorx"
	"github.com/xh-polaris/meowchat-post-rpc/internal/convertor"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model"
	"github.com/xh-polaris/meowchat-post-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-post-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RetrievePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRetrievePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetrievePostLogic {
	return &RetrievePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RetrievePostLogic) RetrievePost(in *pb.RetrievePostReq) (*pb.RetrievePostResp, error) {
	data, err := l.svcCtx.PostModel.FindOne(l.ctx, in.PostId)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errorx.ErrNoSuchPost
	default:
		return nil, err
	}
	return &pb.RetrievePostResp{Post: convertor.ConvertPost(data)}, nil
}
