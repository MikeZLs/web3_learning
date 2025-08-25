package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	OpenLoginList []struct {
		Name string
		Icon string
		Href string
	}
	Etcd string
}
