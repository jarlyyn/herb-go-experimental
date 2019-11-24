package assembler

import (
	"reflect"
)

type Assembler struct {
	config *Config
	part   Part
	path   Path
	parent reflect.Type
	step   Step
}

func (a *Assembler) Value() (interface{}, error) {
	return a.part.GetData(nil)
}
func (a *Assembler) Assemble(v interface{}) (err error) {
	dv, err := a.Value()
	if err != nil {
		return err
	}
	return a.config.Unifiers.Unify(a, dv)
}
func (a *Assembler) CheckType() (tp interface{}, err error) {
	v, err := a.Value()
	if err != nil {
		return nil, err
	}
	rt := getReflectType(v)
	for k := range a.config.Checkers {
		result, err := a.config.Checkers[k].CheckType(a, rt)
		if err != nil {
			return nil, err
		}
		if result {
			return a.config.Checkers[k].Type, nil
		}
	}
	return nil, nil
}
func (a *Assembler) WithConfig(c *Config) *Assembler {
	return &Assembler{
		config: c,
		part:   a.part,
		path:   a.path,
		parent: a.parent,
		step:   a.step,
	}
}

func (a *Assembler) WithPart(p Part) *Assembler {
	return a.WithChild(p, nil, nil)
}

func (a *Assembler) WithChild(p Part, parent reflect.Type, step Step) *Assembler {
	return &Assembler{
		config: a.config,
		part:   p,
		path:   a.path.Join(step),
		parent: parent,
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

func (a *Assembler) Parent() reflect.Type {
	return a.parent
}

var BaseAssembler = &Assembler{
	path: NewSteps(),
}
