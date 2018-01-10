package sqluser

import (
	"testing"

	"github.com/jarlyyn/herb-go-experimental/app-member"

	"github.com/herb-go/herb/user"
	"github.com/jmoiron/sqlx"
)

const accountype = "test"

func InitDB() *sqlx.DB {
	var db = sqlx.MustOpen(config.Driver, config.Conn)
	db.MustExec("TRUNCATE account")
	db.MustExec("TRUNCATE password")
	db.MustExec("TRUNCATE token")
	db.MustExec("TRUNCATE user")
	return db
}
func TestInterface(t *testing.T) {
	var U = New(nil, config.Prefix, FlagWithAccount|FlagWithPassword|FlagWithToken|FlagWithUser)
	var service = member.New(nil)
	service.Install(U.Account())
	service.Install(U.Password())
	service.Install(U.Token())
	service.Install(U.User())
}

func TestSqluser(t *testing.T) {
	var unusedUID = "-test"
	var password = "password"
	var newpassword = "newpassword"
	var wrongpassword = "rongpassword"
	account1, err := user.CaseSensitiveAcountType.NewAccount(accountype, "account1")
	if err != nil {
		panic(err)
	}
	account1plus, err := user.CaseSensitiveAcountType.NewAccount(accountype, "account1plus")
	if err != nil {
		panic(err)
	}
	account2, err := user.CaseSensitiveAcountType.NewAccount(accountype, "account2")
	if err != nil {
		panic(err)
	}
	var U = New(InitDB(), config.Prefix, FlagWithAccount|FlagWithPassword|FlagWithToken|FlagWithUser)
	account := U.Account()
	if account.DBTableName() != U.AccountTableName() {
		t.Error(account.DBTableName())
	}
	a, err := account.Accounts(account1plus.Account)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 0 {
		t.Error(a)
	}
	uid, err := account.AccountToUID(*account1)
	if err != nil {
		t.Error(err)
	}
	if uid != "" {
		t.Error(uid)
	}
	uid1, err := account.Register(*account1)
	if err != nil {
		t.Error(err)
	}
	uid, err = account.AccountToUID(*account1)
	if err != nil {
		t.Error(err)
	}
	if uid != uid1 {
		t.Error(uid)
	}
	uid2, err := account.AccountToUIDOrRegister(*account2)
	if err != nil {
		t.Error(err)
	}
	uid, err = account.AccountToUID(*account2)
	if err != nil {
		t.Error(err)
	}
	if uid != uid2 {
		t.Error(uid)
	}
	uid, err = account.AccountToUIDOrRegister(*account2)
	if err != nil {
		t.Error(err)
	}
	if uid != uid2 {
		t.Error(uid)
	}
	a, err = account.Accounts(uid1, uid2, account1plus.Account)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 2 {
		t.Error(a)
	}
	if len(a[uid1]) != 1 || a[uid1][0].Account != account1.Account || a[uid1][0].Keyword != accountype {
		t.Error(a[uid1])
	}
	if len(a[uid2]) != 1 || a[uid2][0].Account != account2.Account || a[uid2][0].Keyword != accountype {
		t.Error(a[uid2])
	}
	uid, err = account.Register(*account1)
	if err != member.ErrAccountRegisterExists {
		t.Error(err)
	}
	uid, err = account.AccountToUID(*account1plus)
	if err != nil {
		t.Error(err)
	}
	if uid != "" {
		t.Error(uid)
	}
	err = account.BindAccounts(uid1, *account1plus)
	if err != nil {
		t.Fatal(err)
	}
	uid, err = account.AccountToUID(*account1plus)
	if err != nil {
		t.Error(err)
	}
	if uid != uid1 {
		t.Error(uid)
	}
	err = account.BindAccounts(uid1, *account1plus)
	if err != user.ErrAccountBindExists {
		t.Error(err)
	}
	err = account.BindAccounts(uid2, *account1plus)
	if err != user.ErrAccountBindExists {
		t.Error(err)
	}
	err = account.UnbindAccounts(uid1, *account1plus)
	if err != nil {
		t.Error(err)
	}
	uid, err = account.AccountToUID(*account1plus)
	if err != nil {
		t.Error(err)
	}
	if uid != "" {
		t.Error(uid)
	}
	userdm := U.User()
	if userdm.DBTableName() != U.UserTableName() {
		t.Error(userdm.DBTableName())
	}
	u, err := userdm.Banned(uid1, uid2, account1plus.Account)
	if err != nil {
		t.Fatal(err)
	}
	if len(u) != 2 {
		t.Error(a)
	}
	if u[uid1] != false {
		t.Error(u[uid1])
	}
	if u[uid2] != false {
		t.Error(u[uid2])
	}
	err = userdm.Ban(uid1, true)
	u, err = userdm.Banned(uid1, uid2, account1plus.Account)
	if err != nil {
		t.Fatal(err)
	}
	if u[uid1] != true {
		t.Error(u[uid1])
	}
	err = userdm.Ban(uid1, false)
	u, err = userdm.Banned(uid1, uid2, account1plus.Account)
	if err != nil {
		t.Fatal(err)
	}
	if u[uid1] != false {
		t.Error(u[uid1])
	}
	err = userdm.Ban(unusedUID, true)
	u, err = userdm.Banned(unusedUID)
	if err != nil {
		t.Fatal(err)
	}
	if u[unusedUID] != true {
		t.Error(u[unusedUID])
	}
	var token = U.Token()
	if token.DBTableName() != U.TokenTableName() {
		t.Error(token.DBTableName())
	}

	tk, err := token.Tokens(uid1)
	if err != nil {
		t.Fatal(err)
	}
	tokenresult := tk[uid1]
	tokenresult2, err := token.Revoke(uid1)
	if tokenresult == tokenresult2 {
		t.Error(tokenresult, tokenresult2)
	}
	tk, err = token.Tokens(uid1)
	if err != nil {
		t.Fatal(err)
	}
	tokenresult = tk[uid1]
	if tokenresult != tokenresult2 {
		t.Error(tokenresult, tokenresult2)
	}
	tokenresult2, err = token.Revoke(uid1)
	if tokenresult == tokenresult2 {
		t.Error(tokenresult, tokenresult2)
	}
	tk, err = token.Tokens(uid1)
	if err != nil {
		t.Fatal(err)
	}
	tokenresult = tk[uid1]
	if tokenresult != tokenresult2 {
		t.Error(tokenresult, tokenresult2)
	}
	p := U.Password()

	if p.DBTableName() != U.PasswordTableName() {
		t.Error(p.DBTableName())
	}

	_, err = p.VerifyPassword(uid1, password)
	if err != member.ErrUserNotFound {
		t.Fatal(err)
	}
	err = p.UpdatePassword(uid1, password)
	if err != nil {
		panic(err)
	}
	bresult, err := p.VerifyPassword(uid1, password)
	if err != nil {
		t.Fatal(err)
	}
	if bresult != true {
		t.Error(bresult)
	}
	bresult, err = p.VerifyPassword(uid1, wrongpassword)
	if err != nil {
		t.Fatal(err)
	}
	if bresult != false {
		t.Error(bresult)
	}
	err = p.UpdatePassword(uid1, newpassword)
	if err != nil {
		panic(err)
	}
	bresult, err = p.VerifyPassword(uid1, newpassword)
	if err != nil {
		t.Fatal(err)
	}
	if bresult != true {
		t.Error(bresult)
	}
	bresult, err = p.VerifyPassword(uid1, password)
	if err != nil {
		t.Fatal(err)
	}
	if bresult != false {
		t.Error(bresult)
	}
	bresult, err = p.VerifyPassword(uid1, wrongpassword)
	if err != nil {
		t.Fatal(err)
	}
	if bresult != false {
		t.Error(bresult)
	}

}

type dbConfig struct {
	Driver string
	Conn   string
	Prefix string
}
