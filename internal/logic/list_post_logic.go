package logic

import (
	"context"

	"github.com/xh-polaris/meowchat-post-rpc/internal/convertor"
	"github.com/xh-polaris/meowchat-post-rpc/internal/model/post"
	"github.com/xh-polaris/meowchat-post-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-post-rpc/pb"
	"github.com/xh-polaris/paginator-go"

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

func parseFilter(fopts *pb.FilterOptions) *post.FilterOptions {
	if fopts != nil {
		return &post.FilterOptions{
			OnlyUserId:   fopts.OnlyUserId,
			OnlyOfficial: fopts.OnlyOfficial,
		}
	}
	return &post.FilterOptions{}
}

func parsePaginator(popts *pb.PaginationOptions) *paginator.PaginationOptions {
	if popts != nil {
		return &paginator.PaginationOptions{
			Limit:     popts.Limit,
			Offset:    popts.Offset,
			Backward:  popts.Backward,
			LastToken: popts.LastToken,
		}
	}
	return &paginator.PaginationOptions{}
}

func (l *ListPostLogic) ListPost(in *pb.ListPostReq) (*pb.ListPostResp, error) {
	resp := new(pb.ListPostResp)
	var posts []*post.Post
	var total int64
	var err error

	filter := parseFilter(in.FilterOptions)

	p := parsePaginator(in.PaginationOptions)

	if in.SearchOptions == nil {
		posts, total, err = l.svcCtx.PostModel.FindManyAndCount(l.ctx, filter, p, post.IdSorter)
		if err != nil {
			return nil, err
		}
	} else {
		switch o := in.SearchOptions.Query.(type) {
		case *pb.SearchOptions_AllFieldsKey:
			posts, total, err = l.svcCtx.PostModel.Search(l.ctx, post.ConvertAllFieldsSearchQuery(o), filter, p, post.ScoreSorter)
		case *pb.SearchOptions_MultiFieldsKey:
			posts, total, err = l.svcCtx.PostModel.Search(l.ctx, post.ConvertMultiFieldsSearchQuery(o), filter, p, post.ScoreSorter)
		}
		if err != nil {
			return nil, err
		}
	}

	resp.Total = total
	if p.LastToken != nil {
		resp.Token = *p.LastToken
	}
	resp.Posts = make([]*pb.Post, 0, len(posts))
	for _, post_ := range posts {
		resp.Posts = append(resp.Posts, convertor.ConvertPost(post_))
	}
	return resp, nil
}
