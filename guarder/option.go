package guarder

import (
	"encoding/json"
)

type GuarderOption interface {
	ApplyToGuarder(g *Guarder) error
}

type VisitorOption interface {
	ApplyToVisitor(v *Visitor) error
}

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

type RequestParamsConfig interface {
	RequestParamsMapperDriver() string
	RequestParamsDriver() string
	DriverConfig() Config
}

func ApplyToGuarder(g *Guarder, c RequestParamsConfig) error {
	config := c.DriverConfig()
	if g.Mapper == nil {
		d := c.RequestParamsMapperDriver()
		driver, err := NewMapperDriver(d, config, "")
		if err != nil {
			return err
		}
		g.Mapper = driver
	}
	if g.Identifier == nil {
		d := c.RequestParamsDriver()
		driver, err := NewIdentifierDriver(d, config, "")
		if err != nil {
			return err
		}
		g.Identifier = driver

	}
	return nil

}

func ApplyToVisitor(v *Visitor, c RequestParamsConfig) error {
	config := c.DriverConfig()
	if v.Mapper == nil {
		d := c.RequestParamsMapperDriver()
		driver, err := NewMapperDriver(d, config, "")
		if err != nil {
			return err
		}
		v.Mapper = driver
	}
	if v.Credential == nil {
		d := c.RequestParamsDriver()
		driver, err := NewCredentialDriver(d, config, "")
		if err != nil {
			return err
		}
		v.Credential = driver

	}
	return nil

}

type RequestParamsConfigMap struct {
	RequestParamsDriverField
	RequestParamsMapperDriverField
	Config ConfigMap
}

func (c *RequestParamsConfigMap) DriverConfig() Config {
	return &c.Config
}

func (c *RequestParamsConfigMap) ApplyToGuarder(g *Guarder) error {
	return ApplyToGuarder(g, c)
}

func (c *RequestParamsConfigMap) ApplyToVisitor(v *Visitor) error {
	return ApplyToVisitor(v, c)
}
