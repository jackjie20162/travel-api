package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
    rest.RestConf
    TravelRpc struct { Target string }
    Payment struct {
        PayPal struct {
            Enabled   bool
            BaseURL   string
            ClientID  string
            Secret    string
            ReturnURL string
            CancelURL string
        }
    }
}
