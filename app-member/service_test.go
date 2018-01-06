package member

import "github.com/herb-go/herb/user"
import "strconv"
import "time"

type testAccountService struct {
	AccountsMap map[string]user.UserAccounts
}

func newTestAccount(uid string) *user.UserAccount {
	account, err := service.NewAccount("test", uid)
	if err != nil {
		panic(err)
	}
	return account
}
func (s *testAccountService) InstallToService(service *Service) {
	service.AccountsService = s
}
func (s *testAccountService) Accounts(uid ...string) (Accounts, error) {
	var a = map[string]user.UserAccounts{}
	for _, v := range uid {
		a[v] = s.AccountsMap[v]
	}
	return a, nil
}
func (s *testAccountService) AccountToUID(account user.UserAccount) (uid string, err error) {
	for uid, v := range s.AccountsMap {
		if v.Exists(&account) {
			return uid, nil
		}
	}
	return "", nil
}
func (s *testAccountService) Register(account user.UserAccount) (uid string, err error) {
	for _, v := range s.AccountsMap {
		if v.Exists(&account) {
			return "", ErrAccountRegisterExists
		}
	}
	uid = strconv.Itoa(len(s.AccountsMap))
	s.AccountsMap[uid] = []user.UserAccount{account}
	return uid, nil
}
func (s *testAccountService) AccountToUIDOrRegister(account user.UserAccount) (uid string, err error) {
	for uid, v := range s.AccountsMap {
		if v.Exists(&account) {
			return uid, nil
		}
	}
	s.AccountsMap[account.Account] = []user.UserAccount{account}
	return account.Account, nil
}
func (s *testAccountService) BindAccounts(uid string, account user.UserAccount) error {
	for _, v := range s.AccountsMap {
		if v.Exists(&account) {
			return user.ErrAccountBindExists
		}
	}
	if s.AccountsMap[uid] == nil {
		s.AccountsMap[uid] = user.UserAccounts{}
	}
	accounts := s.AccountsMap[uid]
	err := accounts.Bind(&account)
	if err != nil {
		return err
	}
	s.AccountsMap[uid] = accounts
	return nil
}
func (s *testAccountService) UnbindAccounts(uid string, account user.UserAccount) error {
	if s.AccountsMap[uid] == nil {
		return user.ErrAccountUnbindNotExists
	}
	accounts := s.AccountsMap[uid]
	err := accounts.Unbind(&account)
	if err != nil {
		return err
	}
	s.AccountsMap[uid] = accounts
	return nil
}

func newTestAccountService() *testAccountService {
	return &testAccountService{
		AccountsMap: map[string]user.UserAccounts{},
	}
}

type testTokenService struct {
	MemberTokens map[string]string
}

func (s *testTokenService) InstallToService(service *Service) {
	service.TokenService = s
}
func (s *testTokenService) Tokens(uid ...string) (Tokens, error) {
	var r = Tokens{}
	for _, v := range uid {
		r[v] = s.MemberTokens[v]
	}
	return r, nil

}
func (s *testTokenService) Revoke(uid string) (string, error) {
	var ts = strconv.FormatInt(time.Now().UnixNano(), 10)
	s.MemberTokens[uid] = ts
	return ts, nil
}
func newTestTokenService() *testTokenService {
	return &testTokenService{
		MemberTokens: map[string]string{},
	}
}

type testBannedService struct {
	BannedMap BannedMap
}

func (s *testBannedService) InstallToService(service *Service) {
	service.BannedService = s
}

func (s *testBannedService) Banned(uid ...string) (BannedMap, error) {
	var r = BannedMap{}
	for _, v := range uid {
		r[v] = s.BannedMap[v]
	}
	return r, nil

}
func (s *testBannedService) Ban(uid string, banned bool) error {
	s.BannedMap[uid] = banned
	return nil
}

func newTestBannedService() *testBannedService {
	return &testBannedService{
		BannedMap: BannedMap{},
	}
}

type testPasswordService struct {
	Passwords map[string]string
}

func (s *testPasswordService) InstallToService(service *Service) {
	service.PasswordService = s
}
func (s *testPasswordService) VerifyPassword(uid string, password string) (bool, error) {

	pass := s.Passwords[uid]
	if pass == "" {
		return false, nil
	}
	return pass == password, nil

}
func (s *testPasswordService) UpdatePassword(uid string, password string) error {
	s.Passwords[uid] = password
	return nil
}

func newTestPasswordService() *testPasswordService {
	return &testPasswordService{
		Passwords: map[string]string{},
	}
}

type userProfiles map[string]user.Profile

var rawUserProfiles = map[string]user.Profile{}

func (p userProfiles) NewMapElement(key string) error {
	p[key] = user.Profile{}
	return nil
}
func (p userProfiles) LoadMapElements(keys ...string) error {
	for _, v := range keys {
		p[v] = rawUserProfiles[v]
	}
	return nil
}

func newTestUesrProfiles() *userProfiles {
	return &userProfiles{}
}

type testRoleService Roles

func (s *testRoleService) InstallToService(service *Service) {
	service.RoleService = s
}
func (s *testRoleService) Roles(uid ...string) (Roles, error) {
	r := Roles{}
	for _, v := range uid {
		r[v] = (*s)[v]
	}
	return r, nil

}

func newTestRoleService() *testRoleService {
	return &testRoleService{}
}
