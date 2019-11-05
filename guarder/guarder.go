package guarder

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/herb-go/herb/user/httpuser"
)

type Guarder interface {
	httpuser.Authorizer
	httpuser.Identifier
}

type GuarderDriver interface {
	Guarder() (Guarder, error)
}
type GuarderProvider struct {
	Driver GuarderDriver
}

func (g *GuarderProvider) Guarder() (Guarder, error) {
	return g.Driver.Guarder()
}

type Factory func(conf Config, prefix string) (GuarderDriver, error)

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
		panic(errors.New("guarder: Register guarder factory is nil"))
	}
	if _, dup := factories[name]; dup {
		panic(errors.New("guarder: Register called twice for factory " + name))
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
func NewDriver(name string, conf Config, prefix string) (GuarderDriver, error) {
	factorysMu.RLock()
	factoryi, ok := factories[name]
	factorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("guarder: unknown driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}
