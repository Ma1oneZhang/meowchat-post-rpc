package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-post-rpc/internal/convertor"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/pagination"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post"
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
	var posts []*post.Post
	var total int64
	var err error

	filter := &post.FilterOptions{
		OnlyUserId:   in.OnlyUserId,
		OnlyOfficial: in.OnlyOfficial,
	}
	p := &pagination.PaginationOptions{
		Limit:     in.Limit,
		Offset:    in.Offset,
		Backward:  in.Backward,
		LastToken: in.LastToken,
	}

	if in.SearchOptions == nil {
		posts, total, err = l.svcCtx.PostModel.FindManyAndCount(l.ctx, filter, p, post.IdSorter)
		if err != nil {
			return nil, err
		}
	} else {
		switch o := in.SearchOptions.(type) {
		case *pb.ListPostReq_AllFieldsKey:
			posts, total, err = l.svcCtx.PostModel.Search(l.ctx, post.ConvertAllFieldsSearchQuery(o), filter, p, post.ScoreSorter)
		case *pb.ListPostReq_MultiFieldsKey:
			posts, total, err = l.svcCtx.PostModel.Search(l.ctx, post.ConvertMultiFieldsSearchQuery(o), filter, p, post.ScoreSorter)
		}
		if err != nil {
			return nil, err
		}
	}

	resp.Total = total
	resp.Token = *p.LastToken
	resp.Posts = make([]*pb.Post, 0, len(posts))
	for _, post_ := range posts {
		resp.Posts = append(resp.Posts, convertor.ConvertPost(post_))
	}
	return resp, nil
}
