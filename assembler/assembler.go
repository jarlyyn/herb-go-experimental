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
	for k := range a.config.Checkers {
		result, err := a.config.Checkers[k].CheckType(a)
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
	return &Assembler{
		config: a.config,
		part:   p,
		path:   a.path,
		parent: nil,
		step:   a.step,
	}
}

func (a *Assembler) WithChild(parent reflect.Type, p Part, steps ...Step) *Assembler {
	if len(steps) == 0 {
		return a
	}
	return &Assembler{
		config: a.config,
		part:   p,
		path:   a.path.Join(steps...),
		parent: parent,
		step:   steps[len(steps)-1],
	}
}

var BaseAssembler = &Assembler{
	path: NewSteps(),
}
