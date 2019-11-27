package unmarshaler

type AssemblerLazyLoader struct {
	Assembler *Assembler
}

func (l *AssemblerLazyLoader) LazyLoad(v interface{}) error {
	_, err := l.Assembler.Assemble(v)
	return err
}

func NewLazyLoader() *AssemblerLazyLoader {
	return &AssemblerLazyLoader{}
}

type LazyLoaderFunc func(v interface{}) error

type LazyLoader interface {
	LazyLoad(v interface{}) error
}
