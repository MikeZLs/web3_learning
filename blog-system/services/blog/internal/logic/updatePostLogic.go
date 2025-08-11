package logic

import (
	"context"

	"blog-system/services/blog/blog"
	"blog-system/services/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePostLogic) UpdatePost(in *blog.UpdatePostRequest) (*blog.UpdatePostResponse, error) {
	// todo: add your logic here and delete this line

	return &blog.UpdatePostResponse{}, nil
}
