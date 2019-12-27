package apiservice

import (
	"github.com/herb-go/herb/service/httpservice/channel"
)

type APIChannel struct {
	Channel           *channel.Channel
	MiddlewareBuilder MiddlewareBuilder
	ValuesLoader      ValuesLoader
}

type APIChanelConfig struct {
	Channel      *channel.Channel
	Guarder      string
	ValuesLoader string
	Config       func(interface{}) error
}
