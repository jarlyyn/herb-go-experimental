package github

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/jarlyyn/herb-go-experimental/httpclient"

	"github.com/herb-go/herb/cache"
	session "github.com/herb-go/herb/cache-session"
	auth "github.com/jarlyyn/herb-go-experimental/app-externalauth"
)

const StateLength = 128

var TokenMask = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")

const oauthURL = "https://github.com/login/oauth/authorize"

const FieldName = "externalauthdriver-github"

type StateSession struct {
	State string
}
type OauthAuthDriver struct {
	client *Client
	scope  string
}

func NewOauthDriver(client *Client, scope string) *OauthAuthDriver {
	return &OauthAuthDriver{
		client: client,
		scope:  scope,
	}
}
func (d *OauthAuthDriver) ExternalLogin(service *auth.Service, w http.ResponseWriter, r *http.Request) {
	bytes, err := cache.RandMaskedBytes(TokenMask, StateLength)
	if err != nil {
		panic(err)
	}
	state := string(bytes)
	authsession := StateSession{
		State: state,
	}
	err = service.Auth.SessionStore.MustGetRequestSession(r).Set(FieldName, authsession)
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(oauthURL)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set("client_id", d.client.ClientID)
	q.Set("scope", d.scope)
	q.Set("state", state)
	q.Set("redirect_uri", service.GetAuthUrl())
	u.RawQuery = q.Encode()
	u.Fragment = "wechat_redirect"
	http.Redirect(w, r, u.String(), 302)
}

func (d *OauthAuthDriver) AuthRequest(service *auth.Service, r *http.Request) (*auth.Result, error) {
	var authsession = &StateSession{}
	q := r.URL.Query()
	var code = q.Get("code")
	if code == "" {
		return nil, nil
	}
	var state = q.Get("state")
	if state == "" {
		return nil, auth.ErrAuthParamsError
	}
	err := service.Auth.SessionStore.MustGetRequestSession(r).Get(FieldName, authsession)
	if err == session.ErrDataNotFound {
		return nil, nil
	}
	if authsession.State == "" || authsession.State != state {
		return nil, auth.ErrAuthParamsError
	}
	err = service.Auth.SessionStore.MustGetRequestSession(r).Del(FieldName)
	if err != nil {
		return nil, err
	}
	result, err := d.client.GetAccessToken(code)
	if err != nil {
		statuscode := httpclient.GetErrorStatusCode(err)
		if statuscode > 400 && statuscode < 500 {
			return nil, auth.ErrAuthParamsError
		}
		return nil, err
	}
	if result.AccessToken == "" {
		return nil, auth.ErrAuthParamsError
	}
	user, err := d.client.GetUser(result.AccessToken)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	authresult := auth.NewResult()
	authresult.Account = user.Login
	authresult.Keyword = service.Keyword
	authresult.Data.SetValue(auth.DataIndexAccessToken, result.AccessToken)
	authresult.Data.SetValue(auth.DataIndexAvatar, user.AvatarURL)
	authresult.Data.SetValue(auth.DataIndexEmail, user.Email)
	authresult.Data.SetValue(auth.DataIndexName, user.Name)
	authresult.Data.SetValue(auth.DataIndexNickname, user.Login)
	authresult.Data.SetValue(auth.DataIndexProfileURL, user.HTMLURL)
	authresult.Data.SetValue(auth.DataIndexID, strconv.Itoa(user.ID))
	authresult.Data.SetValue(auth.DataIndexCompany, user.Company)
	authresult.Data.SetValue(auth.DataIndexLocation, user.Location)
	authresult.Data.SetValue(auth.DataIndexWebsite, user.Blog)
	return authresult, nil
}
