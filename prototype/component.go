package prototype

type Component struct {
	Type        Type
	Name        *Name
	Label       *Label
	Icon        *Icon
	Description *Description
	Tags        *Tags
	URL         *URL
	Children    *Components
	Data        *Data
	Options     *Options
}

func NewComponent() *Component {
	return &Component{}
}

type Components []*Component

func NewComponents() *Components {
	return &Components{}
}

type ComponentBuilder interface {
	Build(*Component) error
}

type ComponentBuilderFunc func(*Component) error

func (f ComponentBuilderFunc) Build(c *Component) error {
	return f(c)
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
	Value interface{}
}

type Options []*Option
