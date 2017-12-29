package role

import (
	"net/http"

	"github.com/herb-go/herb/user"
)

type RuleService interface {
	Rule(*http.Request) (Rule, error)
}
type RoleService interface {
	Roles(uid string) (Roles, error)
}

type Authority struct {
	RoleService RoleService
	Identifier  user.Identifier
}

func (a *Authority) Authorize(rs RuleService) func(r *http.Request) (bool, error) {
	return func(r *http.Request) (bool, error) {
		uid, err := a.Identifier.IdentifyRequest(r)
		if err != nil {
			return false, err
		}
		roles, err := a.RoleService.Roles(uid)
		if err != nil {
			return false, err
		}
		rm, err := rs.Rule(r)
		if err != nil {
			return false, err
		}
		return rm.Execute(roles...)
	}
}
