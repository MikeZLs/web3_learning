package svc

import (
	"fim_server/core"
	"fim_server/lim_user/user_rpc/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
	}
}
