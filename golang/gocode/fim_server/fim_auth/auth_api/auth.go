package main

import (
	"fim_server/common/etcd"
	"fim_server/fim_auth/auth_api/internal/config"
	"fim_server/fim_auth/auth_api/internal/handler"
	"fim_server/fim_auth/auth_api/internal/svc"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "fim_auth/auth_api/etc/auth.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
