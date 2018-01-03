package member

import (
	"reflect"

	"github.com/jarlyyn/herb-go-experimental/cache-cachedmap"
)

type Members struct {
	Service      *Service
	Accounts     Accounts
	BannedMap    BannedMap
	RevokeTokens RevokeTokens
	Roles        Roles
	Dataset      map[string]cachedmap.CachedMap
}

func (m *Members) LoadBanned(keys ...string) (BannedMap, error) {
	return m.BannedMap, m.Service.Banned().Load(&m.BannedMap, keys...)
}
func (m *Members) LoadRevokeTokens(keys ...string) (RevokeTokens, error) {
	return m.RevokeTokens, m.Service.Revoke().Load(&m.RevokeTokens, keys...)
}
func (m *Members) LoadAccount(keys ...string) (Accounts, error) {
	return m.Accounts, m.Service.Accounts().Load(&m.Accounts, keys...)
}
func (m *Members) LoadRoles(keys ...string) (Roles, error) {
	return m.Roles, m.Service.Roles().Load(&m.Roles, keys...)
}
func (m *Members) LoadData(field string, keys ...string) (cachedmap.CachedMap, error) {
	return m.Dataset[field], m.Service.Data().Load(field, m.Dataset[field], keys...)
}

func (m *Members) Data(field string) cachedmap.CachedMap {
	return m.Dataset[field]
}
func NewMembers(s *Service) *Members {
	var member = &Members{
		Service:      s,
		Accounts:     Accounts{},
		BannedMap:    BannedMap{},
		Roles:        Roles{},
		RevokeTokens: RevokeTokens{},
	}
	member.Dataset = make(map[string]cachedmap.CachedMap, len(s.DataServices))
	var mapvalue = reflect.ValueOf(member.Data)
	for k := range s.DataServices {
		d := reflect.MakeMap(s.DataServices[k])
		mapvalue.SetMapIndex(reflect.ValueOf(k), d)
	}
	return member
}
