package logic

import (
	"blog-system/services/blog/types/blog"
	"context"
	"encoding/json"
	"fmt"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePostLogic) CreatePost(req *types.CreatePostReq) (resp *types.CreatePostResp, err error) {
	// go-zero 会自动解析 JWT 并将 claims 放入 context
	// 获取用户 ID - 这是 go-zero 的标准方式
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	if userId == 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	// 调用 Blog RPC 服务
	res, err := l.svcCtx.BlogRpc.CreatePost(l.ctx, &blog.CreatePostRequest{
		Title:   req.Title,
		Content: req.Content,
		UserId:  userId,
	})
	if err != nil {
		logx.Errorf("create post error: %v", err)
		return nil, err
	}

	return &types.CreatePostResp{
		Id:      res.Id,
		Message: res.Message,
	}, nil
}
