package auth

type Data map[DataIndex][]string

func (d *Data) Value(index DataIndex) string {
	data, ok := (*d)[index]
	if ok == false || len(data) == 0 {
		return ""
	}
	return data[0]
}

func (d *Data) Values(index DataIndex) []string {
	data, ok := (*d)[index]
	if ok == false {
		return nil
	}
	return data
}

func (d *Data) SetValue(index DataIndex, value string) {
	(*d)[index] = []string{value}
}

func (d *Data) SetValues(index DataIndex, values []string) {
	(*d)[index] = values
}

func (d *Data) AddValue(index DataIndex, value string) {
	data, ok := (*d)[index]
	if ok == false {
		data = []string{}
	}
	data = append(data, value)
	(*d)[index] = data
}

type Result struct {
	Keyword string
	Account string
	Data    Data
}

func NewResult() *Result {
	return &Result{
		Data: map[DataIndex][]string{},
	}
}
