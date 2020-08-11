package usersystem

import "github.com/herb-go/herbsecurity/authorize/role"

type RolesService interface {
	LoadRoles(...string) (map[string]*role.Roles, error)
	//Start start service
	Start() error
	//Stop stop service
	Stop() error
	//Reload reload user data
	Reload(string) error
}
