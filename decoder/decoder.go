package decoder

import "reflect"

// type Decoder struct {
// 	Checkers []TypeChecker
// 	Unifiers Unifiers
// }

// func (d *Decoder) DecodeDataSource(v interface{}, n Node) error {
// 	ctx := NewContext()
// 	ctx.Decoder = d
// 	ctx.Node = n
// 	data, err := n.GetData(nil)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = d.Unifiers.Unify(ctx, reflect.ValueOf(v), data)
// 	return err
// }

// func (d *Decoder) CheckType(ctx *Context, rt reflect.Type) (tp interface{}, err error) {
// 	for k := range d.Checkers {
// 		result, err := d.Checkers[k].CheckType(ctx, rt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if result {
// 			return d.Checkers[k].Type, nil
// 		}
// 	}
// 	return nil, nil
// }

// func NewDecoder() *Decoder {
// 	return &Decoder{}
// }

// func NewCommonDecoder() *Decoder {
// 	d := NewDecoder()
// 	return d
// }

type Decoder struct {
	config *Config
	node   Node
	path   Path
	parent Node
	step   Step
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
