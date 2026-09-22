package config

import (
	"github.com/suyuan32/simple-admin-common/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	TravelRpc    zrpc.RpcClientConf
	DatabaseConf config.DatabaseConf
	AllowInit    bool
	Email        EmailConf
	Redis        redis.RedisConf
	Captcha      CaptchaConf
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
		Stripe struct {
			Enabled        bool
			SecretKey      string
			PublishableKey string
			WebhookSecret  string
		}
	}
	Shopify struct {
		Enabled       bool
		StoreDomain   string
		AccessToken   string
		APIVersion    string
		WebhookSecret string
		SyncEnabled   bool
	}
}

// EmailConf 邮件发送配置
type EmailConf struct {
	To       string
	From     string
	Host     string
	Secret   string
	Nickname string
	Port     int
	IsSSL    bool
}

// CaptchaConf 验证码配置
type CaptchaConf struct {
	KeyPrefix string `mapstructure:"key-prefix"`
	Expire    int    `mapstructure:"expire"`     // 验证码过期时间（秒）
	EmailCode int    `mapstructure:"email-code"` // 邮件验证码过期时间（秒）
}
