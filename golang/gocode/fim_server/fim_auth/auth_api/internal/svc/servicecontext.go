package svc

import (
	"fim_server/core"
	"fim_server/fim_auth/auth_api/internal/config"
	"fim_server/fim_user/user_rpc/types/user_rpc"
	"fim_server/fim_user/user_rpc/users"
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Redis   *redis.Client
	UserRpc user_rpc.UsersClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	redisClient := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)

	return &ServiceContext{
		Config:  c,
		DB:      mysqlDb,
		Redis:   redisClient,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
	}
}
