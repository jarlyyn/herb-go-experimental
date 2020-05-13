package prototype

type Prototype struct {
	Type        Type
	Name        *Name
	Label       *Label
	Icon        *Icon
	Description *Description
	Tags        *Tags
	URL         *URL
	Value       *Value
	Children    *Children
	Options     *Options
}

func New(t Type) *Prototype {
	return &Prototype{
		Type: t,
	}
}

type Children []*Prototype

func NewChildren() *Children {
	return &Children{}
}

type Factory interface {
	Create() (*Prototype, error)
}

type BuilderFunc func() (*Prototype, error)

func (f BuilderFunc) Create() (*Prototype, error) {
	return f()
}

type URL string

type Description string

type Icon string

type Label string

type Name string

type Tag string

type Tags []Tag

type Type string

type Value interface{}

type OptionName string

type Option struct {
	Name  OptionName
	Value Value
}

type Options []*Option
