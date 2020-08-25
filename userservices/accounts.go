package usersystem

import "github.com/herb-go/user"

type AccountsService interface {
	LoadAccounts(...string) (map[string]*user.Accounts, error)
	//AccountToUID query uid by user account.
	//Return user id and any error if raised.
	//Return empty string as userid if account not found.
	AccountToUID(account *user.Account) (uid string, err error)
	//Register create new user with given account.
	//Return created user id and any error if raised.
	//Privoder should return ErrAccountRegisterExists if account is used.
	Register(account *user.Account) (uid string, err error)
	//BindAccount bind account to user.
	//Return any error if raised.
	//If account exists,user.ErrAccountBindingExists should be rasied.
	BindAccount(uid string, account *user.Account) error
	//UnbindAccount unbind account from user.
	//Return any error if raised.
	//If account not exists,user.ErrAccountUnbindingNotExists should be rasied.
	UnbindAccount(uid string, account *user.Account) error
	//Start start service
	Start() error
	//Stop stop service
	Stop() error
	//Reload reload user data
	Reload(string) error
}
