package prototype

type Builder interface {
	Build(*Prototype) error
}
