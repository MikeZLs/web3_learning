package logic

import (
	"context"

	"blog-system/services/blog/blog"
	"blog-system/services/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(in *blog.CreatePostRequest) (*blog.CreatePostResponse, error) {
	// todo: add your logic here and delete this line

	return &blog.CreatePostResponse{}, nil
}
