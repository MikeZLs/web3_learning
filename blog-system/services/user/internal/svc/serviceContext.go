package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"user/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Mysql  sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config: c,
		Mysql:  sqlConn,
	}
}
