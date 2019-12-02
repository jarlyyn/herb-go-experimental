package unmarshaler

//Part part interface
//Part is object used to assembler struct data.
//Unmarshal driver should create part to unmarshal data with assembler
type Part interface {
	//Iter return part iter.
	//Nil should be returned If part is not iterable.
	Iter() (*PartIter, error)
	//Value return part value as empty interface.
	//Shold only used when part is not iterable
	Value() (interface{}, error)
}

//PartIter part iter struct
type PartIter struct {
	//Step current iter step
	Step Step
	//Part current iter part
	Part Part
	//Next return next part iter and any error if rasied.
	//nil shold be returned if iter finished.
	Next func() (*PartIter, error)
}
