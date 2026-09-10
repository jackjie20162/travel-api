package main

import (
	"flag"
	"fmt"
	"net/http"

	"gitee.com/meinongyihe/travel-api/internal/config"
	"gitee.com/meinongyihe/travel-api/internal/handler"
	"gitee.com/meinongyihe/travel-api/internal/i18n"
	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/travel-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	// Load i18n bundles and register the localized error handler so every error
	// is returned as {code, msg} translated by the request locale.
	i18n.Load()
	httpx.SetErrorHandlerCtx(i18n.ErrorHandler)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// Global middleware: inject tenant/merchant from HTTP headers into context
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := middleware.ContextFromRequest(r)
			r = r.WithContext(ctx)
			next(w, r)
		}
	})

	ctx := svc.NewServiceContext(c)
	defer ctx.DB.Close()
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting travel-api at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
