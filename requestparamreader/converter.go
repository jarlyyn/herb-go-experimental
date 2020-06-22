package requestparamreader

type Converter interface {
	ConvertParam([]byte) ([]byte, error)
}

type ConverterFactory interface {
	CreateConverter(loader func(interface{}) error) (Converter, error)
}
