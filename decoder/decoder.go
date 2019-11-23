package decoder

import "reflect"

type Decoder struct {
	config *Config
	node   Node
	path   Path
	parent Node
	step   Step
}

func (d *Decoder) Value() (interface{}, error) {
	return d.node.GetData(nil)
}
func (d *Decoder) Decode(v interface{}) (err error) {
	dv, err := d.Value()
	if err != nil {
		return err
	}
	return d.config.Unifiers.Unify(d, reflect.ValueOf(v), dv)
}
func (d *Decoder) CheckType(rt reflect.Type) (tp interface{}, err error) {
	for k := range d.config.Checkers {
		result, err := d.config.Checkers[k].CheckType(d, rt)
		if err != nil {
			return nil, err
		}
		if result {
			return d.config.Checkers[k].Type, nil
		}
	}
	return nil, nil
}
func (d *Decoder) WithConfig(c *Config) *Decoder {
	return &Decoder{
		config: c,
		node:   d.node,
		path:   d.path,
		parent: d.parent,
		step:   d.step,
	}
}

func (d *Decoder) WithNode(n Node) *Decoder {
	return &Decoder{
		config: d.config,
		node:   n,
		path:   d.path,
		parent: nil,
		step:   d.step,
	}
}

func (d *Decoder) WithChild(n Node, steps ...Step) *Decoder {
	if len(steps) == 0 {
		return d
	}
	return &Decoder{
		config: d.config,
		node:   n,
		path:   d.path.Join(steps...),
		parent: d.node,
		step:   steps[len(steps)-1],
	}
}

var BaseDecoder = &Decoder{
	path: NewSteps(),
}
