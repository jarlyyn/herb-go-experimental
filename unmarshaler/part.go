package unmarshaler

type Part interface {
	Iter() (*PartIter, error)
	Value() (interface{}, error)
}

type PartIter struct {
	Step Step
	Part Part
	Next func() (*PartIter, error)
}
