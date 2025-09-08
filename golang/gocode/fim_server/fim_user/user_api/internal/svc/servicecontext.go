package svc

import (
	"fim_server/core"
	"fim_server/fim_user/user_api/internal/config"
	"fim_server/fim_user/user_rpc/types/user_rpc"
	"fim_server/fim_user/user_rpc/users"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user_rpc.UsersClient
	DB      *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config:  c,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		DB:      mysqlDb,
	}
}
