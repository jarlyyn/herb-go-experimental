package requestparam

type CommonFieldConfig struct {
	Field string
}
type FieldConfig struct {
	Name   string
	Type   string
	Config func(v interface{}) error
}
type ConvertConfig struct {
	Name   string
	Type   string
	Config func(v interface{}) error
}
type Config struct {
	Fields   []*FieldConfig
	Converts []*ConvertConfig
}
