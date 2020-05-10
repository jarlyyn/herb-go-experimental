package formdefinition

type FieldType interface {
	Apply(*Field)
}

type FieldTypeFunc func(*Field)

func (f FieldTypeFunc) Apply(field *Field) {
	f(field)
}

type Field struct {
	Label       string
	Description string
	Name        string
	Requeired   bool
	Type        string
	Default     interface{}
	Option      interface{}
}

func NewField() *Field {
	return &Field{}
}

type Form struct {
	Title       string
	Description string
	Operation   string
	Fields      []*Field
	URL         string
}

type Template func() ([]byte, error)

type Compiler interface {
	Compile(*Form) (Template, error)
}
