package auth

import (
	"github.com/herb-go/herb/user"
)

type Result struct {
	user.UserAccount
	Data Profile
}

func NewResult() *Result {
	return &Result{
		Data: map[ProfileIndex][]string{},
	}
}
