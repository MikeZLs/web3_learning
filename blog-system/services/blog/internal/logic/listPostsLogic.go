package logic

import (
	"context"

	"blog-system/services/blog/blog"
	"blog-system/services/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPostsLogic) ListPosts(in *blog.ListPostsRequest) (*blog.ListPostsResponse, error) {
	// todo: add your logic here and delete this line

	return &blog.ListPostsResponse{}, nil
}
