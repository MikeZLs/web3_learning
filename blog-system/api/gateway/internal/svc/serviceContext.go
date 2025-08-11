package svc

import (
	"blog-system/services/blog/blogclient"
	"blog-system/services/comment/commentclient"
	"blog-system/services/user/userclient"
	"gateway/internal/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	UserRpc    userclient.User       // 使用 goctl 生成的客户端接口
	BlogRpc    blogclient.Blog       // 使用 goctl 生成的客户端接口
	CommentRpc commentclient.Comment // 使用 goctl 生成的客户端接口
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		BlogRpc:    blogclient.NewBlog(zrpc.MustNewClient(c.BlogRpc)),
		CommentRpc: commentclient.NewComment(zrpc.MustNewClient(c.CommentRpc)),
	}
}
