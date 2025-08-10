package db

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
	"web3_learning/golang/gocode/blog-api/internal/config"
)

func NewMysql(mysqlConfig config.Database) sqlx.SqlConn {
	mysql := sqlx.NewMysql(mysqlConfig.DataSource)
	db, err := mysql.RawDB()
	if err != nil {
		panic(err)
	}
	cxt, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(mysqlConfig.ConnectTimeOut))
	defer cancel()
	err = db.PingContext(cxt)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	return mysql
}
