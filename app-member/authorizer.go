package member

import (
	"net/http"

	role "github.com/herb-go/herb/user-role"
)

type Authorizer struct {
	Service     *Service
	RuleService role.RuleService
}

func (a *Authorizer) Authorize(r *http.Request) (bool, error) {
	uid, err := a.Service.IdentifyRequest(r)
	if err != nil {
		return false, err
	}
	if uid == "" {
		return false, nil
	}
	var members = a.Service.GetMembersFromRequest(r)
	if a.Service.BannedProvider != nil {
		_, err = members.LoadBanned(uid)
		if err != nil {
			return false, err
		}
		if members.BannedMap[uid] == true {
			return false, nil
		}
	}

	if a.Service.RoleProvider == nil {
		return true, nil
	}
	if a.RuleService == nil {
		return true, nil
	}
	rolesmap, err := members.LoadRoles(uid)
	if err != nil {
		return false, err
	}
	if rolesmap == nil {
		return false, err
	}
	roles := rolesmap[uid]
	if roles == nil {
		return false, nil
	}
	rm, err := a.RuleService.Rule(r)
	if err != nil {
		return false, err
	}
	return rm.Execute(roles...)
}
