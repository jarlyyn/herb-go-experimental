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

type DriverConfig interface {
	MapperDriverName() string
	DriverName() string
	DriverConfig() Config
}

func ApplyToGuarder(g *Guarder, c DriverConfig) error {
	config := c.DriverConfig()
	if g.Mapper == nil {
		d := c.MapperDriverName()
		driver, err := NewMapperDriver(d, config, "")
		if err != nil {
			return err
		}
		g.Mapper = driver
	}
	if g.Identifier == nil {
		d := c.DriverName()
		driver, err := NewIdentifierDriver(d, config, "")
		if err != nil {
			return err
		}
		g.Identifier = driver

	}
	return nil

}

func ApplyToVisitor(v *Visitor, c DriverConfig) error {
	config := c.DriverConfig()
	if v.Mapper == nil {
		d := c.MapperDriverName()
		driver, err := NewMapperDriver(d, config, "")
		if err != nil {
			return err
		}
		v.Mapper = driver
	}
	if v.Credential == nil {
		d := c.DriverName()
		driver, err := NewCredentialDriver(d, config, "")
		if err != nil {
			return err
		}
		v.Credential = driver

	}
	return nil

}

type DriverField struct {
	Driver       string
	staticDriver string
}

func (f *DriverField) SetStaticDriver(d string) {
	f.staticDriver = d
}
func (f *DriverField) DriverName() string {
	if f.staticDriver == "" {
		return f.Driver
	}
	return f.staticDriver
}

type MapperDriverField struct {
	MapperDriver       string
	staticMapperDriver string
}

func (f *MapperDriverField) SetStaticMapperDriver(d string) {
	f.staticMapperDriver = d
}
func (f *MapperDriverField) MapperDriverName() string {
	if f.staticMapperDriver == "" {
		return f.MapperDriver
	}
	return f.staticMapperDriver
}

func NewDriverConfigMap() *DirverConfigMap {
	return &DirverConfigMap{}
}

type DirverConfigMap struct {
	DriverField
	MapperDriverField
	Config ConfigMap
}

func (c *DirverConfigMap) DriverConfig() Config {
	return &c.Config
}

func (c *DirverConfigMap) ApplyToGuarder(g *Guarder) error {
	return ApplyToGuarder(g, c)
}

func (c *DirverConfigMap) ApplyToVisitor(v *Visitor) error {
	return ApplyToVisitor(v, c)
}
