package prototype

type FieldType interface {
	ApplyTo(*Field)
}
type FieldTypeFunc func(*Field)

func (f FieldTypeFunc) ApplyTo(field *Field) {
	f(field)
}

type Text struct {
}

type TextArea struct {
}
type Password struct {
}
type Integer struct {
}

type Float struct {
}

type Boolean struct {
}

type DropDown struct {
}

type Radio struct {
}

type Hidden struct {
}
