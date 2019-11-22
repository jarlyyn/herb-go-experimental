package decoder

type DataSource interface {
	GetData(path Path) (interface{}, error)
}

type MapDataSource map[string]interface{}

func (d *MapDataSource) GetData(path Path) (interface{}, error) {
	var ok bool
	var m map[string]interface{} = *d
	var v interface{}
	if path == nil {
		return m, nil
	}
	path.Clone()
	var step Step
	for {
		step, path = path.Unshift()
		v, ok = m[step.String()]
		if ok == false {

		}
		m, ok = v.(map[string]interface{})
		if ok == false {

		}
		if path == nil {
			return v, nil
		}
	}
}
