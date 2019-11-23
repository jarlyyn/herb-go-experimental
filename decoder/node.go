package decoder

type Child struct {
	Field Step
	Node  Node
}
type Node interface {
	Children() ([]Child, error)
	GetData(path Path) (interface{}, error)
}

type MapNode struct {
	Value interface{}
}

func NewMapNode(v interface{}) *MapNode {
	return &MapNode{
		Value: v,
	}
}
func (d *MapNode) GetData(path Path) (interface{}, error) {
	if path == nil {
		return d.Value, nil
	}
	var ok bool
	var m map[string]interface{}
	var v interface{}
	var step Step
	v = d.Value
	path = path.Clone()
	for {
		step, path = path.Unshift()
		m, ok = v.(map[string]interface{})
		if ok == false {

		}
		v, ok = m[step.String()]
		if ok == false {

		}
		if path == nil {
			return v, nil
		}
	}
}
