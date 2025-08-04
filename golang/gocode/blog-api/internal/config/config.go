package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth     Auth
	Database Database
}

type Database struct {
	DataSource     string
	ConnectTimeOut int
}

type Auth struct {
	AccessSecret string
	AccessExpire int64
}
