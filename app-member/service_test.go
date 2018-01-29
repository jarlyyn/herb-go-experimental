package member

import "github.com/herb-go/herb/user"
import "strconv"
import "time"

type testAccountProvider struct {
	AccountsMap map[string]user.UserAccounts
}

func newTestAccount(uid string) *user.UserAccount {
	account, err := service.NewAccount("test", uid)
	if err != nil {
		panic(err)
	}
	return account
}
func (s *testAccountProvider) InstallToMember(service *Service) {
	service.AccountsProvider = s
}
func (s *testAccountProvider) Accounts(uid ...string) (Accounts, error) {
	var a = map[string]user.UserAccounts{}
	for _, v := range uid {
		a[v] = s.AccountsMap[v]
	}
	return a, nil
}
func (s *testAccountProvider) AccountToUID(account user.UserAccount) (uid string, err error) {
	for uid, v := range s.AccountsMap {
		if v.Exists(&account) {
			return uid, nil
		}
	}
	return "", nil
}
func (s *testAccountProvider) Register(account user.UserAccount) (uid string, err error) {
	for _, v := range s.AccountsMap {
		if v.Exists(&account) {
			return "", ErrAccountRegisterExists
		}
	}
	uid = strconv.Itoa(len(s.AccountsMap))
	s.AccountsMap[uid] = []user.UserAccount{account}
	return uid, nil
}
func (s *testAccountProvider) AccountToUIDOrRegister(account user.UserAccount) (uid string, err error) {
	for uid, v := range s.AccountsMap {
		if v.Exists(&account) {
			return uid, nil
		}
	}
	s.AccountsMap[account.Account] = []user.UserAccount{account}
	return account.Account, nil
}
func (s *testAccountProvider) BindAccounts(uid string, account user.UserAccount) error {
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
func (s *testAccountProvider) UnbindAccounts(uid string, account user.UserAccount) error {
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

func newTestAccountProvider() *testAccountProvider {
	return &testAccountProvider{
		AccountsMap: map[string]user.UserAccounts{},
	}
}

type testTokenProvider struct {
	MemberTokens map[string]string
}

func (s *testTokenProvider) InstallToMember(service *Service) {
	service.TokenProvider = s
}
func (s *testTokenProvider) Tokens(uid ...string) (Tokens, error) {
	var r = Tokens{}
	for _, v := range uid {
		r[v] = s.MemberTokens[v]
	}
	return r, nil

}
func (s *testTokenProvider) Revoke(uid string) (string, error) {
	var ts = strconv.FormatInt(time.Now().UnixNano(), 10)
	s.MemberTokens[uid] = ts
	return ts, nil
}
func newTestTokenProvider() *testTokenProvider {
	return &testTokenProvider{
		MemberTokens: map[string]string{},
	}
}

type testBannedService struct {
	BannedMap BannedMap
}

func (s *testBannedService) InstallToMember(service *Service) {
	service.BannedProvider = s
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

func newTestBannedProvider() *testBannedService {
	return &testBannedService{
		BannedMap: BannedMap{},
	}
}

type testPasswordProvider struct {
	Passwords map[string]string
}

func (s *testPasswordProvider) InstallToMember(service *Service) {
	service.PasswordProvider = s
}
func (s *testPasswordProvider) VerifyPassword(uid string, password string) (bool, error) {

	pass := s.Passwords[uid]
	if pass == "" {
		return false, nil
	}
	return pass == password, nil

}
func (s *testPasswordProvider) UpdatePassword(uid string, password string) error {
	s.Passwords[uid] = password
	return nil
}

func newTestPasswordProvider() *testPasswordProvider {
	return &testPasswordProvider{
		Passwords: map[string]string{},
	}
}

type userProfiles map[string]map[string][]string

var rawUserProfiles = map[string]map[string][]string{}

func (p userProfiles) NewMapElement(key string) error {
	p[key] = map[string][]string{}
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

type testRoleProvider Roles

func (s *testRoleProvider) InstallToMember(service *Service) {
	service.RoleProvider = s
}
func (s *testRoleProvider) Roles(uid ...string) (Roles, error) {
	r := Roles{}
	for _, v := range uid {
		r[v] = (*s)[v]
	}
	return r, nil

}

func newTestRoleProvider() *testRoleProvider {
	return &testRoleProvider{}
}
