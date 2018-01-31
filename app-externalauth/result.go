package auth

type Result struct {
	Account string
	Data    Profile
}

func NewResult() *Result {
	return &Result{
		Data: map[ProfileIndex][]string{},
	}
}
