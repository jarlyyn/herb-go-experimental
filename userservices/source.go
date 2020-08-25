package usersystem

import (
	"sync"

	"github.com/herb-go/user"
	"github.com/herb-go/user/profile"
	"github.com/herb-go/user/status"
)

type Source interface {
	LoadProfile(system *System, idlist ...string) (map[string]*profile.Profile, error)
	SetProfile(map[string]*profile.Profile) error
	LoadAccounts(system *System, idlist ...string) (map[string]*user.Accounts, error)
	SetAccounts(map[string]*user.Accounts) error
	LoadStatus(system *System, idlist ...string) (map[string]status.Status, error)
	SetStatus(map[string]status.Status) error
	//Reload reload user data
	Reload(string) error
}

type MapSource struct {
	profilemap  sync.Map
	accountsmap sync.Map
	statusmap   sync.Map
}

func (s *MapSource) LoadProfile(system *System, idlist ...string) (map[string]*profile.Profile, error) {
	result := map[string]*profile.Profile{}
	for k := range idlist {
		v, ok := s.profilemap.Load(idlist[k])
		if ok {
			result[idlist[k]] = v.(*profile.Profile)
		}
	}
	return result, nil
}
func (s *MapSource) SetProfile(m map[string]*profile.Profile) error {
	for k := range m {
		s.profilemap.Store(k, m[k])
	}
	return nil
}
func (s *MapSource) LoadAccounts(system *System, idlist ...string) (map[string]*user.Accounts, error) {
	result := map[string]*user.Accounts{}
	for k := range idlist {
		v, ok := s.accountsmap.Load(idlist[k])
		if ok {
			result[idlist[k]] = v.(*user.Accounts)
		}
	}
	return result, nil
}
func (s *MapSource) SetAccounts(m map[string]*user.Accounts) error {
	for k := range m {
		s.accountsmap.Store(k, m[k])
	}
	return nil
}

func (s *MapSource) LoadStatus(system *System, idlist ...string) (map[string]status.Status, error) {
	result := map[string]status.Status{}
	for k := range idlist {
		v, ok := s.statusmap.Load(idlist[k])
		if ok {
			result[idlist[k]] = v.(status.Status)
		}
	}
	return result, nil

}
func (s *MapSource) SetStatus(m map[string]status.Status) error {
	for k := range m {
		s.statusmap.Store(k, m[k])
	}
	return nil
}

//Reload reload user data
func (s *MapSource) Reload(id string) error {
	s.statusmap.Delete(id)
	s.profilemap.Delete(id)
	s.accountsmap.Delete(id)
	return nil
}

func NewMapSource() *MapSource {
	return &MapSource{}
}

type SourceService interface {
	CreateSource() (Source, error)
	//Start start service
	Start() error
	//Stop stop service
	Stop() error
	//Reload reload user data
	Reload(string) error
}

type SourceServiceFunc func() (Source, error)

func (f SourceServiceFunc) CreateSource() (Source, error) {
	return f()
}

func (f SourceServiceFunc) Start() error {
	return nil
}
func (f SourceServiceFunc) Stop() error {
	return nil
}

var MapSourceService = SourceServiceFunc(func() (Source, error) {
	return NewMapSource(), nil
})
