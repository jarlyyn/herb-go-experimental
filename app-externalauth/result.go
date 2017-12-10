package auth

import "github.com/jarlyyn/herb-go-experimental/user"

type Result struct {
	user.UserAccount
	Data user.Profile
}

func NewResult() *Result {
	return &Result{
		Data: map[user.ProfileIndex][]string{},
	}
}
