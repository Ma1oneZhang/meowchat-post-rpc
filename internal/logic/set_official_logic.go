package logic

import (
	"context"

	"github.com/xh-polaris/meowchat-post-rpc/internal/model"
	"github.com/xh-polaris/meowchat-post-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-post-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetOfficialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetOfficialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetOfficialLogic {
	return &SetOfficialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetOfficialLogic) SetOfficial(in *pb.SetOfficialReq) (*pb.SetOfficialResp, error) {
	err := l.svcCtx.PostModel.UpdateFlags(l.ctx, in.PostId, map[model.PostFlag]bool{
		model.OfficialFlag: !in.IsRemove,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SetOfficialResp{}, nil
}
