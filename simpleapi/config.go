package simpleapi

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/herb-go/util/httpserver"
)

var defaultConfig = &httpserver.Config{
	Net:  "tcp",
	Addr: ":6789",
}

var configLock = sync.Mutex{}

var config *httpserver.Config

var errConfigSetted = func() error {
	configcontent, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return fmt.Errorf("simpleapi : config has been setted as \"%s\"", string(configcontent))
}

var CleanConfig = func() {
	configLock.Lock()
	defer configLock.Unlock()
	config = nil
}
var SetConfig = func(c *httpserver.Config) error {
	configLock.Lock()
	defer configLock.Unlock()
	if config != nil {
		return errConfigSetted()
	}
	config = c
	return nil
}

var Config = func() *httpserver.Config {
	configLock.Lock()
	defer configLock.Unlock()
	if config == nil {
		config = defaultConfig
	}
	return config
}
