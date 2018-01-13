package sqluser

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"crypto/rand"
	"crypto/sha256"

	"github.com/herb-go/herb/user"
	"github.com/jarlyyn/herb-go-experimental/app-member"
	"github.com/jarlyyn/herb-go-experimental/model-sql-query"
	"github.com/jarlyyn/herb-go-experimental/model-sqlxdatamapper"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

const (
	FlagEmpty        = 0
	FlagWithAccount  = 1
	FlagWithPassword = 2
	FlagWithToken    = 4
	FlagWithUser     = 8
)

const (
	UserStatusNormal = 0
	UserStatusBanned = 1
)

var RandomBytesLength = 32
var ErrHashMethodNotFound = errors.New("password hash method not found")

type HashFunc func(key string, salt string, password string) ([]byte, error)

var DefaultAccountTableName = "account"
var DefaultPasswordTableName = "password"
var DefaultTokenTableName = "token"
var DefaultUserTableName = "user"
var DefaultHashMethod = "sha256"
var HashFuncMap = map[string]HashFunc{
	"sha256": func(key string, salt string, password string) ([]byte, error) {
		var val = []byte(key + salt + password)
		var s256 = sha256.New()
		s256.Write(val)
		val = s256.Sum(nil)
		s256.Write(val)
		return []byte(hex.EncodeToString(s256.Sum(nil))), nil
	},
}

func New(db *sqlx.DB, prefix string, flag int) *User {
	return &User{
		DB: xdatamapper.NewDB(db, prefix),
		Tables: Tables{
			AccountTableName:  DefaultAccountTableName,
			PasswordTableName: DefaultPasswordTableName,
			TokenTableName:    DefaultTokenTableName,
			UserTableName:     DefaultUserTableName,
		},
		HashMethod:     DefaultHashMethod,
		UIDGenerater:   UUID,
		TokenGenerater: Timestamp,
		SaltGenerater:  RandomBytes,
		Flag:           flag,
	}
}

type Tables struct {
	AccountTableName  string
	PasswordTableName string
	TokenTableName    string
	UserTableName     string
}
type User struct {
	DB             xdatamapper.DB
	Tables         Tables
	Flag           int
	UIDGenerater   func() (string, error)
	TokenGenerater func() (string, error)
	SaltGenerater  func() (string, error)
	HashMethod     string
	PasswordKey    string
}

func RandomBytes() (string, error) {
	var bytes = make([]byte, RandomBytesLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func UUID() (string, error) {
	return uuid.NewV1().String(), nil
}
func Timestamp() (string, error) {
	return strconv.FormatInt(time.Now().UnixNano(), 10), nil
}
func (u *User) HasFlag(flag int) bool {
	return u.Flag&flag != 0
}
func (u *User) AccountTableName() string {
	return u.DB.TableName(u.Tables.AccountTableName)
}
func (u *User) PasswordTableName() string {
	return u.DB.TableName(u.Tables.PasswordTableName)
}
func (u *User) TokenTableName() string {
	return u.DB.TableName(u.Tables.TokenTableName)
}
func (u *User) UserTableName() string {
	return u.DB.TableName(u.Tables.UserTableName)
}
func (u *User) Account() *AccountDataMapper {
	return &AccountDataMapper{
		DataMapper: xdatamapper.New(u.DB, u.Tables.AccountTableName),
		User:       u,
	}
}

func (u *User) Password() *PasswordDataMapper {
	return &PasswordDataMapper{
		DataMapper: xdatamapper.New(u.DB, u.Tables.PasswordTableName),
		User:       u,
	}
}
func (u *User) Token() *TokenDataMapper {
	return &TokenDataMapper{
		DataMapper: xdatamapper.New(u.DB, u.Tables.TokenTableName),
		User:       u,
	}
}
func (u *User) User() *UserDataMapper {
	return &UserDataMapper{
		DataMapper: xdatamapper.New(u.DB, u.Tables.UserTableName),
		User:       u,
	}
}

type AccountDataMapper struct {
	xdatamapper.DataMapper
	User    *User
	Service *member.Service
}

func (a *AccountDataMapper) InstallToService(service *member.Service) {
	service.AccountsService = a
	a.Service = service
}

func (a *AccountDataMapper) Unbind(uid string, account user.UserAccount) error {
	stmt, err := a.DB().Beginx()
	if err != nil {
		return err
	}
	defer stmt.Rollback()
	Delete := query.NewDelete(a.DBTableName())
	Delete.Where.Condition = query.And(
		query.New("account.uid = ?", uid),
		query.New("account.keyword = ?", account.Keyword),
		query.New("account.account = ?", account.Account),
	)
	Query := Delete.Query()
	_, err = stmt.Exec(Query.Command, Query.QueryArgs()...)
	// _, err = stmt.Exec("DELETE From "+a.DBTableName()+" WHERE uid=? and keyword=? and account=?", uid, account.Keyword, account.Account)
	if err != nil {
		return err
	}
	return stmt.Commit()

}

func (a *AccountDataMapper) Bind(uid string, account user.UserAccount) error {
	stmt, err := a.DB().Beginx()
	if err != nil {
		return err
	}
	defer stmt.Rollback()
	var u = ""
	row := stmt.QueryRow("SELECT uid from "+a.DBTableName()+" where keyword= ? and account = ?", account.Keyword, account.Account)
	err = row.Scan(&u)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return user.ErrAccountBindExists

	}

	var CreatedTime = time.Now().Unix()
	_, err = stmt.Exec("Insert Into "+a.DBTableName()+" (uid,keyword,account,created_time) VALUES(?,?,?,?)", uid, account.Keyword, account.Account, CreatedTime)
	if err != nil {
		return err
	}
	return stmt.Commit()
}
func (a *AccountDataMapper) FindOrInsert(UIDGenerater func() (string, error), account user.UserAccount) (string, error) {
	var result = AccountModel{}
	stmt, err := a.DB().Beginx()
	if err != nil {
		return "", err
	}
	defer stmt.Rollback()
	row := a.DB().QueryRow("SELECT uid,keyword,account,created_time from "+a.DBTableName()+" where keyword= ? and account = ?", account.Keyword, account.Account)
	err = row.Scan(&result.UID, &result.Keyword, &result.Account, &result.CreatedTime)
	if err == nil {
		return result.UID, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}
	uid, err := UIDGenerater()
	var CreatedTime = time.Now().Unix()
	_, err = stmt.Exec("Insert Into "+a.DBTableName()+" (uid,keyword,account,created_time) VALUES(?,?,?,?)", uid, account.Keyword, account.Account, CreatedTime)
	if err != nil {
		return "", err
	}
	if a.User.HasFlag(FlagWithUser) {
		_, err = stmt.Exec("Insert Into "+a.User.UserTableName()+" (uid,status,created_time,updated_time) VALUES(?,?,?,?)", uid, UserStatusNormal, CreatedTime, CreatedTime)
		if err != nil {
			return "", err
		}
	}
	return uid, stmt.Commit()
}
func (a *AccountDataMapper) Insert(uid string, keyword string, account string) error {
	stmt, err := a.DB().Beginx()
	if err != nil {
		return err
	}
	defer stmt.Rollback()
	var u = ""
	row := stmt.QueryRow("SELECT uid from "+a.DBTableName()+" where keyword= ? and account = ?", keyword, account)
	err = row.Scan(&u)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return member.ErrAccountRegisterExists
	}
	var CreatedTime = time.Now().Unix()
	_, err = stmt.Exec("Insert Into "+a.DBTableName()+" (uid,keyword,account,created_time) VALUES(?,?,?,?)", uid, keyword, account, CreatedTime)
	if err != nil {
		return err
	}
	if a.User.HasFlag(FlagWithUser) {
		_, err = stmt.Exec("Insert Into "+a.User.UserTableName()+" (uid,status,created_time,updated_time) VALUES(?,?,?,?)", uid, UserStatusNormal, CreatedTime, CreatedTime)
		if err != nil {
			return err
		}
	}
	return stmt.Commit()
}
func (a *AccountDataMapper) Find(keyword string, account string) (AccountModel, error) {
	var result = AccountModel{}
	if keyword == "" || account == "" {
		return result, sql.ErrNoRows
	}
	row := a.DB().QueryRow("SELECT uid,keyword,account,created_time from "+a.DBTableName()+" where keyword= ? and account = ?", keyword, account)
	err := row.Scan(&result.UID, &result.Keyword, &result.Account, &result.CreatedTime)
	return result, err
}
func (a *AccountDataMapper) FindAllByUID(uids ...string) ([]AccountModel, error) {
	var result = []AccountModel{}
	if len(uids) == 0 {
		return result, nil
	}
	query, args, err := sqlx.In("SELECT uid,keyword,account from "+a.DBTableName()+" where uid IN (?) ", uids)
	if err != nil {
		return nil, err
	}
	query = a.DB().Rebind(query)
	rows, err := a.DB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := AccountModel{}
		err = rows.Scan(&v.UID, &v.Keyword, &v.Account)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}
func (a *AccountDataMapper) Accounts(uid ...string) (member.Accounts, error) {
	models, err := a.FindAllByUID(uid...)
	if err != nil {
		return nil, err
	}
	result := member.Accounts{}
	for _, v := range models {
		if result[v.UID] == nil {
			result[v.UID] = user.UserAccounts{}
		}
		account := user.UserAccount{Keyword: v.Keyword, Account: v.Account}
		result[v.UID] = append(result[v.UID], account)
	}
	return result, nil
}
func (a *AccountDataMapper) AccountToUID(account user.UserAccount) (uid string, err error) {
	model, err := a.Find(account.Keyword, account.Account)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return model.UID, err
}
func (a *AccountDataMapper) Register(account user.UserAccount) (uid string, err error) {
	uid, err = a.User.UIDGenerater()
	if err != nil {
		return
	}
	err = a.Insert(uid, account.Keyword, account.Account)
	return
}
func (a *AccountDataMapper) AccountToUIDOrRegister(account user.UserAccount) (uid string, err error) {
	return a.FindOrInsert(a.User.UIDGenerater, account)
}
func (a *AccountDataMapper) BindAccounts(uid string, account user.UserAccount) error {
	return a.Bind(uid, account)
}
func (a *AccountDataMapper) UnbindAccounts(uid string, account user.UserAccount) error {
	return a.Unbind(uid, account)
}

type AccountModel struct {
	UID         string
	Keyword     string
	Account     string
	CreatedTime int64
}
type PasswordDataMapper struct {
	xdatamapper.DataMapper
	User    *User
	Service *member.Service
}

func (p *PasswordDataMapper) InstallToService(service *member.Service) {
	service.PasswordService = p
	p.Service = service
}
func (p *PasswordDataMapper) Find(uid string) (PasswordModel, error) {
	var result = PasswordModel{}
	if uid == "" {
		return result, sql.ErrNoRows
	}
	row := p.DB().QueryRow("SELECT hash_method,salt,password,updated_time from "+p.DBTableName()+" where uid= ? ", uid)
	result.UID = uid
	err := row.Scan(&result.HashMethod, &result.Salt, &result.Password, &result.UpdatedTime)
	return result, err
}
func (p *PasswordDataMapper) InsertOrUpdate(model *PasswordModel) error {
	stmt, err := p.DB().Beginx()
	if err != nil {
		return err
	}
	defer stmt.Rollback()
	r, err := stmt.Exec("UPDATE "+p.DBTableName()+" SET hash_method = ? ,salt = ? ,password=? ,updated_time=? where uid = ?", model.HashMethod, model.Salt, model.Password, model.UpdatedTime, model.UID)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 0 {
		return stmt.Commit()
	}
	_, err = stmt.Exec("Insert Into "+p.DBTableName()+" (uid,hash_method,salt,password,updated_time) VALUES(?,?,?,?,?)", model.UID, model.HashMethod, model.Salt, model.Password, model.UpdatedTime)
	if err != nil {
		return err
	}
	return stmt.Commit()
}
func (p *PasswordDataMapper) VerifyPassword(uid string, password string) (bool, error) {
	model, err := p.Find(uid)
	if err == sql.ErrNoRows {
		return false, member.ErrUserNotFound
	}
	if err != nil {
		return false, err
	}
	hash := HashFuncMap[model.HashMethod]
	if hash == nil {
		return false, ErrHashMethodNotFound
	}
	hashed, err := hash(p.User.PasswordKey, model.Salt, password)
	if err != nil {
		return false, err
	}
	return bytes.Compare(hashed, model.Password) == 0, nil
}
func (p *PasswordDataMapper) UpdatePassword(uid string, password string) error {
	salt, err := p.User.SaltGenerater()
	if err != nil {
		return err
	}
	hash := HashFuncMap[p.User.HashMethod]
	if hash == nil {
		return ErrHashMethodNotFound
	}
	hashed, err := hash(p.User.PasswordKey, salt, password)
	if err != nil {
		return err
	}
	model := &PasswordModel{
		UID:         uid,
		HashMethod:  p.User.HashMethod,
		Salt:        salt,
		Password:    hashed,
		UpdatedTime: time.Now().Unix(),
	}
	return p.InsertOrUpdate(model)
}

type PasswordModel struct {
	UID         string
	HashMethod  string
	Salt        string
	Password    []byte
	UpdatedTime int64
}

type TokenDataMapper struct {
	xdatamapper.DataMapper
	User    *User
	Service *member.Service
}

func (t *TokenDataMapper) InstallToService(service *member.Service) {
	service.TokenService = t
	t.Service = service
}

func (t *TokenDataMapper) InsertOrUpdate(uid string, token string) error {
	stmt, err := t.DB().Beginx()
	if err != nil {
		return err
	}
	defer stmt.Rollback()
	var CreatedTime = time.Now().Unix()
	r, err := stmt.Exec("UPDATE "+t.DBTableName()+" SET token = ? ,updated_time = ? where uid = ?", token, CreatedTime, uid)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 0 {
		return stmt.Commit()
	}
	_, err = stmt.Exec("Insert Into "+t.DBTableName()+" (uid,token,updated_time) VALUES(?,?,?)", uid, token, CreatedTime)
	if err != nil {
		return err
	}
	return stmt.Commit()
}
func (t *TokenDataMapper) FindAllByUID(uids ...string) ([]TokenModel, error) {
	var result = []TokenModel{}
	if len(uids) == 0 {
		return result, nil
	}
	query, args, err := sqlx.In("SELECT uid,token from "+t.DBTableName()+" where uid IN (?) ", uids)
	if err != nil {
		return nil, err
	}
	query = t.DB().Rebind(query)
	rows, err := t.DB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := TokenModel{}
		err = rows.Scan(&v.UID, &v.Token)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}
func (t *TokenDataMapper) Tokens(uid ...string) (member.Tokens, error) {
	models, err := t.FindAllByUID(uid...)
	if err != nil {
		return nil, err
	}
	result := member.Tokens{}
	for _, v := range models {
		result[v.UID] = v.Token
	}
	return result, nil

}
func (t *TokenDataMapper) Revoke(uid string) (string, error) {
	token, err := t.User.TokenGenerater()
	if err != nil {
		return "", err
	}
	return token, t.InsertOrUpdate(uid, token)
}

type TokenModel struct {
	UID         string
	Token       string
	UpdatedTime string
}
type UserDataMapper struct {
	xdatamapper.DataMapper
	User    *User
	Service *member.Service
}

func (u *UserDataMapper) InstallToService(service *member.Service) {
	service.BannedService = u
	u.Service = service
}

func (u *UserDataMapper) FindAllByUID(uids ...string) ([]UserModel, error) {
	var result = []UserModel{}
	if len(uids) == 0 {
		return result, nil
	}
	query, args, err := sqlx.In("SELECT uid,status from "+u.DBTableName()+" where uid IN (?) ", uids)
	if err != nil {
		return nil, err
	}
	query = u.DB().Rebind(query)
	rows, err := u.DB().Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := UserModel{}
		err = rows.Scan(&v.UID, &v.Status)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func (u *UserDataMapper) InsertOrUpdate(uid string, status int) error {
	stmt, err := u.DB().Beginx()
	if err != nil {
		return err
	}
	defer stmt.Rollback()
	var CreatedTime = time.Now().Unix()
	r, err := stmt.Exec("UPDATE "+u.DBTableName()+" SET status = ? ,updated_time = ? where uid = ?", status, CreatedTime, uid)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 0 {
		return stmt.Commit()
	}
	_, err = stmt.Exec("Insert Into "+u.DBTableName()+" (uid,status,updated_time,created_time) VALUES(?,?,?,?)", uid, status, CreatedTime, CreatedTime)
	if err != nil {
		return err
	}
	return stmt.Commit()
}

func (u *UserDataMapper) Banned(uid ...string) (member.BannedMap, error) {
	models, err := u.FindAllByUID(uid...)
	if err != nil {
		return nil, err
	}
	result := member.BannedMap{}
	for _, v := range models {
		result[v.UID] = (v.Status == UserStatusBanned)
	}
	return result, nil
}
func (u *UserDataMapper) Ban(uid string, banned bool) error {
	var status int
	if banned {
		status = UserStatusBanned
	} else {
		status = UserStatusNormal
	}
	return u.InsertOrUpdate(uid, status)

}

type UserModel struct {
	UID         string
	CreatedTime int64
	UpdateTIme  int64
	Status      int
}
