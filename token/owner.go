package token

type Owner string

func (o Owner) NewToken(id ID) *Token {
	return &Token{
		Owner: o,
		ID:    id,
	}
}
