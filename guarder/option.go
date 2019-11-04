package guarder

import (
	"encoding/json"

	"github.com/herb-go/fetch"
)

type Config interface {
	Get(key string, v interface{}) error
}

//ConfigMap config in map format.
type ConfigMap map[string]interface{}

//Get get value from config map.
//Return any error if raised.
func (c *ConfigMap) Get(key string, v interface{}) error {
	i, ok := (*c)[key]
	if !ok {
		return nil
	}
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

//Set set value to config map.
//Return any error if raised.
func (c *ConfigMap) Set(key string, v interface{}) error {
	(*c)[key] = v
	return nil
}

type CredentialOption struct {
	Clients fetch.Clients
	Driver  string
	Config  Config
}

type CredentialOptionConfigMap struct {
	Clients fetch.Clients
	Driver  string
	Config  ConfigMap
}

func (o *CredentialOptionConfigMap) Option() *CredentialOption {
	return &CredentialOption{
		Clients: o.Clients,
		Driver:  o.Driver,
		Config:  &o.Config,
	}
}

type GuarderOption struct {
	Driver string
	Config Config
}

type GuarderOptionConfigMap struct {
	Driver string
	Config ConfigMap
}

func (o *GuarderOptionConfigMap) Option() *GuarderOption {
	return &GuarderOption{
		Driver: o.Driver,
		Config: &o.Config,
	}
}
