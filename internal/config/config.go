package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	TravelRpc struct {
		Target string
	}
}
