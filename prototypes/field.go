package prototype

type Field struct {
	*Brief
	Name    string
	Value   interface{}
	Type    string
	Options map[string]interface{}
}

func NewField() *Field {
	return &Field{
		Brief:   NewBreif(),
		Options: map[string]interface{}{},
	}
}
