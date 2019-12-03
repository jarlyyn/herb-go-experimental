package unmarshaler

import (
	"reflect"
)

//Assembler data assembler struct
type Assembler struct {
	config *Config
	part   Part
	path   Path
	step   Step
}

//Assemble assemble data to given value.
//Return assemble result and any error if raised.
func (a *Assembler) Assemble(v interface{}) (ok bool, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
			if err == nil {
				panic(r)
			}
			ok = false
			err = NewAssemblerError(a, err)
		}
	}()
	ok, err = a.config.Unifiers.Unify(a, v)
	return ok, NewAssemblerError(a, err)
}

//CheckType check given reflect type type.
//Return type and any error if raised.
func (a *Assembler) CheckType(rt reflect.Type) (tp Type, err error) {
	return a.Config().CheckType(a, rt)
}

//WithConfig create new assembler with given config
func (a *Assembler) WithConfig(c *Config) *Assembler {
	return &Assembler{
		config: c,
		part:   a.part,
		path:   a.path,
		step:   a.step,
	}
}

//WithPart create new assembler with given part
func (a *Assembler) WithPart(p Part) *Assembler {
	return a.WithChild(p, nil)
}

//WithChild create assembler with given child part and step
func (a *Assembler) WithChild(p Part, step Step) *Assembler {
	return &Assembler{
		config: a.config,
		part:   p,
		path:   a.path.Join(step),
		step:   step,
	}
}

//Config return assembler config
func (a *Assembler) Config() *Config {
	return a.config
}

//Part return assembler part
func (a *Assembler) Part() Part {
	return a.part
}

//Path return assembler path
func (a *Assembler) Path() Path {
	return a.path
}

//Step return current assembler step
func (a *Assembler) Step() Step {
	return a.step
}

//EmptyAssembler empty assembler
var EmptyAssembler = &Assembler{
	path: NewSteps(),
}
