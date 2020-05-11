package prototype

type UserInput struct {
	Action    API
	Operation *Brief
	*Item
	Fields []*Field
}

func (i *UserInput) WithField(f *Field) {
	i.Fields = append(i.Fields, f)
}

func NewUserInput() *UserInput {
	return &UserInput{
		Operation: NewBreif(),
		Item:      NewItem(),
	}
}

func GetFieldTypes(inputs ...*UserInput) map[string]bool {
	result := map[string]bool{}
	for k := range inputs {
		for fk := range inputs[k].Fields {
			result[inputs[k].Fields[fk].Name] = true
		}
	}
	return result
}
