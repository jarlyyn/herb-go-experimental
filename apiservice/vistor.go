package apiservice

import (
	"github.com/herb-go/fetcher"
)

type APICaller struct {
	Preset         *fetcher.Preset
	CommandBuilder CommandBuilder
	ValuesProvider ValuesProvider
}

type APICallerConfig struct {
	EndPoint  *fetcher.Server
	Vistor    string
	ValueType string
	Config    func(v interface{}) error
}
