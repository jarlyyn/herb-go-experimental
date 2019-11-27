package unmarshaler

type Part interface {
	Child(step Step) (Part, error)
	Iter() (*PartIter, error)
	Value() (interface{}, error)
}

type PartIter struct {
	Step Step
	Part Part
	Next func() (*PartIter, error)
}
