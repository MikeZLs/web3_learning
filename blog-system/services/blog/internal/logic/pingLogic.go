package logic

import (
	"context"

	"blog-system/services/blog/blog"
	"blog-system/services/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *blog.Request) (*blog.Response, error) {
	// todo: add your logic here and delete this line

	return &blog.Response{}, nil
}
