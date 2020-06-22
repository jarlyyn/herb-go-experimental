package requestparamreader

type FieldConfig struct {
	Type   string
	Config func(v interface{}) error
}
type ConvertConfig struct {
	Type   string
	Config func(v interface{}) error
}
type Config struct {
	Fields   []*FieldConfig
	Converts []*ConvertConfig
}
