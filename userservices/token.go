package usersystem

import "github.com/herb-go/herbsecurity/authority/service/token"

type TokenService interface {
	token.Manager
	//Reload reload user data
	Reload(string) error
}
