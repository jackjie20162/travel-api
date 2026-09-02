package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
    rest.RestConf
    TravelRpc struct { Target string }
    Payment struct {
        FrontendReturnURL string
        PayPal struct {
            Enabled   bool
            BaseURL   string
            ClientID  string
            Secret    string
            WebhookID string
            ReturnURL string
            CancelURL string
        }
    }
}
