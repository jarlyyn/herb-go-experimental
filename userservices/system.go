package usersystem

type System struct {
	Keyword         Keyword
	SourceService   SourceService
	AccountsService AccountsService
	ProfileService  ProfileService
	RolesService    RolesService
	StatusService   StatusService
	PasswordService PasswordService
	TokenService    TokenService
}

//Reload reload user data
func (s *System) Reload(id string) error {
	return Reload(id, s.AccountsService, s.ProfileService, s.RolesService, s.StatusService, s.PasswordService, s.TokenService)
}
