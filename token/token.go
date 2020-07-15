package token

type Token struct {
	Owner  Owner
	ID     ID
	Secret Secret
}

func New() *Token {
	return &Token{}
}
