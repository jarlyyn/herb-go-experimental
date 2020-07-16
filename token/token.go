package token

type Owner string

func (o Owner) NewToken() *Token {
	t := New()
	t.Owner = o
	return t
}

type ID string

type Secret []byte

type Token struct {
	Owner  Owner
	ID     ID
	Secret Secret
}

func New() *Token {
	return &Token{}
}
