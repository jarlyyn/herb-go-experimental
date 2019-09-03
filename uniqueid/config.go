package uniqueid

import (
	"encoding/json"
)

//Option generator option interface.
type Option interface {
	ApplyTo(*Generator) error
}

// OptionConfigJSON option config in json format.
type OptionConfigJSON struct {
	Driver string
	Config ConfigJSON
}

//NewOptionConfigJSON create new option config.
func NewOptionConfigJSON() *OptionConfigJSON {
	return &OptionConfigJSON{
		Config: map[string]string{},
	}
}

//ApplyTo apply option to generator.
func (o *OptionConfigJSON) ApplyTo(g *Generator) error {
	driver, err := NewDriver(o.Driver, &o.Config, "")
	if err != nil {
		return err
	}
	g.Driver = driver
	return nil
}

// OptionConfigMap option config in map format.
type OptionConfigMap struct {
	Driver string
	Config ConfigMap
}

//ApplyTo apply option to file store.
func (o *OptionConfigMap) ApplyTo(g *Generator) error {
	driver, err := NewDriver(o.Driver, &o.Config, "")
	if err != nil {
		return err
	}
	g.Driver = driver
	return nil
}

//NewOptionConfigMap create new option config.
func NewOptionConfigMap() *OptionConfigMap {
	return &OptionConfigMap{
		Config: map[string]interface{}{},
	}
}

// Config confit interface
type Config interface {
	//Get get value form given key.
	//Return any error if raised.
	Get(key string, v interface{}) error
}

//ConfigJSON JSON format config
type ConfigJSON map[string]string

//Get get value form given key.
//Return any error if raised.
func (c *ConfigJSON) Get(key string, v interface{}) error {
	s, ok := (*c)[key]
	if !ok {
		return nil
	}
	return json.Unmarshal([]byte(s), v)
}

//Set set value to given key.
//Return any error if raised.
func (c *ConfigJSON) Set(key string, v interface{}) error {
	s, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	(*c)[key] = string(s)
	return nil
}

//ConfigMap Map  format config
type ConfigMap map[string]interface{}

//Get get value form given key.
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

//Set set value to given key.
//Return any error if raised.
func (c *ConfigMap) Set(key string, v interface{}) error {
	(*c)[key] = v
	return nil
}
