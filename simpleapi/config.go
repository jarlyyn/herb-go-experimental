package simpleapi

import (
	"github.com/herb-go/util/httpserver"
)

var defaultConfig = &httpserver.Config{
	Net:  "tcp",
	Addr: ":6789",
}
