package main

import (
	"flag"
	"fmt"

	"gitee.com/meinongyihe/travel-api/internal/config"
	"gitee.com/meinongyihe/travel-api/internal/handler"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/travel-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	fmt.Printf("Starting travel-api at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
