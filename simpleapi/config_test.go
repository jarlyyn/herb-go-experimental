package simpleapi

import (
	"testing"

	"github.com/herb-go/util/httpserver"
)

func TestConfig(t *testing.T) {
	var err error
	defer func() {
		CleanConfig()
	}()
	CleanConfig()
	if config != nil {
		t.Fatal(defaultConfig)
	}
	c := Config()
	if c != defaultConfig {
		t.Fatal(c)
	}
	CleanConfig()
	config := &httpserver.Config{}
	err = SetConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	c = Config()
	if c != config {
		t.Fatal(c)
	}
	err = SetConfig(config)
	if err == nil {
		t.Fatal(err)
	}
}
