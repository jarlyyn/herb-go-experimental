package usersystem

import "github.com/herb-go/herb/user"

type PasswordService interface {
	user.PasswordVerifier
	//PasswordChangeable return password changeable
	PasswordChangeable() bool
	//UpdatePassword update user password
	//Return any error if raised
	UpdatePassword(uid string, password string) error
	//Start start service
	Start() error
	//Stop stop service
	Stop() error
	//Reload reload user data
	Reload(string) error
}
