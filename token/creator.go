package token

type Creator interface {
	Create(Owner) (*Token, error)
}

type CreatorFunc func(owner Owner) (*Token, error)

func (f CreatorFunc) Create(owner Owner) (*Token, error) {
	return f(owner)
}
