package config

import (
	"github.com/suyuan32/simple-admin-common/config"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	TravelRpc    zrpc.RpcClientConf
	DatabaseConf config.DatabaseConf
	AllowInit    bool
	Email        EmailConf
	Payment      struct {
		FrontendReturnURL string
		PayPal            struct {
			Enabled    bool
			BaseURL    string
			ClientID   string
			Secret     string
			WebhookID  string
			ReturnURL  string
			CancelURL  string
			SkipVerify bool
		}
	}
}

// EmailConf 邮件发送配置
type EmailConf struct {
	To      string
	From    string
	Host    string
	Secret  string
	Nickname string
	Port    int
	IsSSL   bool `mapstructure:"is-ssl"`
}
