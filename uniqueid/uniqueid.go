package uniqueid

import (
	"fmt"
	"sort"
	"sync"
)

//Generator unique id generator
type Generator struct {
	Driver Driver
}

//GenerateID generate unique id.
//Return  generated id and any error if rasied.
func (g *Generator) GenerateID() (string, error) {
	return g.Driver.GenerateID()
}

//MustGenerateID generate unique id.
//Return  generated id.
//Panic if any error raised
func (g *Generator) MustGenerateID() string {
	id, err := g.Driver.GenerateID()
	if err != nil {
		panic(err)
	}
	return id
}

//NewGenerator create new Generator
func NewGenerator() *Generator {
	return &Generator{}
}

//Driver unique id generator driver interface
type Driver interface {
	//GenerateID generate unique id.
	//Return  generated id and any error if rasied.
	GenerateID() (string, error)
}

// Factory unique id generator driver create factory.
type Factory func(conf Config, prefix string) (Driver, error)

var (
	factorysMu sync.RWMutex
	factories  = make(map[string]Factory)
)

// Register makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, f Factory) {
	factorysMu.Lock()
	defer factorysMu.Unlock()
	if f == nil {
		panic("unique: Register cache factory is nil")
	}
	if _, dup := factories[name]; dup {
		panic("unique: Register called twice for factory " + name)
	}
	factories[name] = f
}

//UnregisterAll unregister all driver
func UnregisterAll() {
	factorysMu.Lock()
	defer factorysMu.Unlock()
	// For tests.
	factories = make(map[string]Factory)
}

//Factories returns a sorted list of the names of the registered factories.
func Factories() []string {
	factorysMu.RLock()
	defer factorysMu.RUnlock()
	var list []string
	for name := range factories {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//NewDriver create new driver with given name,config and prefix.
//Reutrn driver created and any error if raised.
func NewDriver(name string, conf Config, prefix string) (Driver, error) {
	factorysMu.RLock()
	factoryi, ok := factories[name]
	factorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("uniqueid: unknown driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}
