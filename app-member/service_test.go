package member

import "github.com/herb-go/herb/user"
import "strconv"
import "time"

type testAccountService struct {
	AccountsMap map[string]user.UserAccounts
}

func newTestAccount(uid string) *user.UserAccount {
	return &user.UserAccount{
		Keyword: "test",
		Account: uid,
	}
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

type testRevokeService struct {
	Tokens map[string]string
}

func (s *testRevokeService) RevokeTokens(uid ...string) (RevokeTokens, error) {
	var r = RevokeTokens{}
	for _, v := range uid {
		r[v] = s.Tokens[v]
	}
	return r, nil

}
func (s *testRevokeService) Revoke(uid string) (string, error) {
	var ts = strconv.FormatInt(time.Now().UnixNano(), 10)
	s.Tokens[uid] = ts
	return ts, nil
}
func newTestRevokeService() *testRevokeService {
	return &testRevokeService{
		Tokens: map[string]string{},
	}
}

type testBannedService struct {
	BannedMap BannedMap
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

func (p *userProfiles) NewMapElement(key string) error {
	(*p)[key] = user.Profile{}
	return nil
}
func (p *userProfiles) LoadMapElements(keys ...string) error {
	for _, v := range keys {
		(*p)[v] = rawUserProfiles[v]
	}
	return nil
}

func newTestUesrProfiles() *userProfiles {
	return &userProfiles{}
}

type testRoleService Roles

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
