package logic

import (
	"context"

	"github.com/xh-polaris/meowchat-post-rpc/internal/convertor"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/paginator"
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
	resp := new(pb.ListPostResp)
	var posts []*model.Post
	var total int64
	var err error

	filter := model.Filter{
		UserId:       in.UserId,
		OnlyOfficial: in.OnlyOfficial,
	}
	p := paginator.BasePaginator{
		Limit:    &in.Count,
		Offset:   &in.Skip,
		Backward: &in.Backward,
	}

	if in.SearchOptions == nil {
		posts, total, err = l.svcCtx.PostModel.FindManyAndCount(l.ctx, model.MongoFilter{
			Filter: filter,
		}, paginator.IdPaginator{
			LastId:        in.LastId,
			BasePaginator: p,
		})
		if err != nil {
			return nil, err
		}
	} else {
		switch o := in.SearchOptions.(type) {
		case *pb.ListPostReq_AllFieldsKey:
			posts, total, err = l.svcCtx.PostModel.Search(l.ctx, convertor.ConvertAllFieldsSearchQuery(o), model.EsFilter{Filter: filter}, p)
		case *pb.ListPostReq_MultiFieldsKey:
			posts, total, err = l.svcCtx.PostModel.Search(l.ctx, convertor.ConvertMultiFieldsSearchQuery(o), model.EsFilter{Filter: filter}, p)
		}
		if err != nil {
			return nil, err
		}
	}

	resp.Total = total
	resp.Posts = make([]*pb.Post, 0, len(posts))
	for _, post := range posts {
		resp.Posts = append(resp.Posts, convertor.ConvertPost(post))
	}
	return resp, nil
}
