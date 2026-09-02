package main

import (
	"flag"

	"gitee.com/meinongyihe/travel-api/internal/config"
	"gitee.com/meinongyihe/travel-api/internal/handler"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var configFile = flag.String("f", "etc/travel-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	conn, err := grpc.NewClient(c.TravelRpc.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { panic(err) }
	defer conn.Close()
	ctx := svc.NewServiceContext(c, conn)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	handler.RegisterHandlers(server, ctx)
	server.Start()
}
