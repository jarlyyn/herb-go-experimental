package requestparamreader

type CommonFieldConfig struct {
	Field string
}
type FieldConfig struct {
	Name     string
	Type     string
	Config   func(v interface{}) error
	Converts []*ConvertConfig
}
type ConvertConfig struct {
	Name   string
	Type   string
	Config func(v interface{}) error
}
type Config struct {
	Fields []*FieldConfig
}
