package usersystem

import "github.com/herb-go/herb/user/status"

type StatusService interface {
	LoadStatus(...string) (map[string]status.Status, error)
	UpdateStatus(string, status.Status) error
	status.Service
	//Start start service
	Start() error
	//Stop stop service
	Stop() error
	//Reload reload user data
	Reload(string) error
}
