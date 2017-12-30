package role

import (
	"net/http"

	"github.com/herb-go/herb/user"
)

type RuleService interface {
	Rule(*http.Request) (Rule, error)
}
type RoleService interface {
	Roles(uid string) (*Roles, error)
}

type Authority struct {
	RoleService RoleService
	Identifier  user.Identifier
}

func (a *Authority) Authority(rs RuleService) func(r *http.Request) (bool, error) {
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
		return rm.Execute(*roles...)
	}
}

func (a *Authority) AuthorizeMiddleware(rs RuleService, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return user.AuthorizeMiddleware(a.Authority(rs), unauthorizedAction)
}

func (a *Authority) RolesAuthorizeOrForbiddenMiddleware(ruleNames ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rs = NewRoles(ruleNames...)
	return a.AuthorizeMiddleware(rs, nil)
}
