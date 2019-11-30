package unmarshaler

import (
	"reflect"
)

type Assembler struct {
	config *Config
	part   Part
	path   Path
	step   Step
}

func (a *Assembler) Value() (interface{}, error) {
	return a.part.Value()
}
func (a *Assembler) Assemble(v interface{}) (ok bool, err error) {
	return a.config.Unifiers.Unify(a, v)
}
func (a *Assembler) CheckType(rt reflect.Type) (tp Type, err error) {
	return a.Config().CheckType(a, rt)
}
func (a *Assembler) WithConfig(c *Config) *Assembler {
	return &Assembler{
		config: c,
		part:   a.part,
		path:   a.path,
		step:   a.step,
	}
}

func (a *Assembler) WithPart(p Part) *Assembler {
	return a.WithChild(p, nil)
}

func (a *Assembler) WithChild(p Part, step Step) *Assembler {
	return &Assembler{
		config: a.config,
		part:   p,
		path:   a.path.Join(step),
		step:   step,
	}
}
func (a *Assembler) Config() *Config {
	return a.config
}

func (a *Assembler) Part() Part {
	return a.part
}

func (a *Assembler) Path() Path {
	return a.path
}

func (a *Assembler) Step() Step {
	return a.step
}

var EmptyAssembler = &Assembler{
	path: NewSteps(),
}
