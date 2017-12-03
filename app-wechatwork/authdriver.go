package wechatwork

import (
	"net/http"
	"net/url"

	cache "github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache-session"
	auth "github.com/jarlyyn/herb-go-experimental/app-externalauth"
)

const FieldName = "externalauthdriver-wechatwork"
const StateLength = 128
const oauthURL = "https://open.weixin.qq.com/connect/oauth2/authorize"
const qrauthURL = "https://open.work.weixin.qq.com/wwopen/sso/qrConnect"

type Session struct {
	State string
}

func authRequestWithAgent(agent *Agent, service *auth.Service, r *http.Request) (*auth.Result, error) {
	var authsession = &Session{}
	var code = r.URL.Query().Get("code")
	if code == "" {
		return nil, nil
	}
	var state = r.URL.Query().Get("state")
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
	info, err := agent.GetUserInfo(code)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	result := &auth.Result{}
	result.Keyword = service.Keyword
	result.Account = info.UserID
	result.Data.SetValue(auth.DataIndexAvatar, info.Avatar)
	result.Data.SetValue(auth.DataIndexEmail, info.Email)
	switch info.Gender {
	case ApiResultGenderMale:
		result.Data.SetValue(auth.DataIndexGender, auth.GenderMale)
	case ApiResultGenderFemale:
		result.Data.SetValue(auth.DataIndexGender, auth.GenderFemale)
	}
	result.Data.SetValue(auth.DataIndexName, info.Name)
	result.Data.SetValue(auth.DataIndexNickname, info.Name)
	return result, nil
}

type OauthAuthDriver struct {
	agent *Agent
	scope string
}

func NewOauthDriver(agent *Agent, scope string) *OauthAuthDriver {
	return &OauthAuthDriver{
		agent: agent,
		scope: scope,
	}
}

func (d *OauthAuthDriver) ExternalLogin(service *auth.Service, w http.ResponseWriter, r *http.Request) {
	bytes, err := cache.RandMaskedBytes(cache.TokenMask, StateLength)
	if err != nil {
		panic(err)
	}
	state := string(bytes)
	authsession := Session{
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
	u.Query().Set("appid", d.agent.CorpID)
	u.Query().Set("agentid", d.agent.AgentID)
	u.Query().Set("scope", d.scope)
	u.Query().Set("state", state)
	u.Query().Set("redirect_uri", service.GetAuthUrl())
	u.Fragment = "wechat_redirect"
	http.Redirect(w, r, u.String(), http.StatusFound)
}
func (d *OauthAuthDriver) AuthRequest(service *auth.Service, r *http.Request) (*auth.Result, error) {
	return authRequestWithAgent(d.agent, service, r)
}

type QRAuthDriver struct {
	agent *Agent
}

func NewQRAuthDriver(agent *Agent) *QRAuthDriver {
	return &QRAuthDriver{
		agent: agent,
	}
}

func (d *QRAuthDriver) ExternalLogin(service *auth.Service, w http.ResponseWriter, r *http.Request) {
	bytes, err := cache.RandMaskedBytes(cache.TokenMask, StateLength)
	if err != nil {
		panic(err)
	}
	state := string(bytes)
	authsession := Session{
		State: state,
	}
	err = service.Auth.SessionStore.MustGetRequestSession(r).Set(FieldName, authsession)
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(qrauthURL)
	if err != nil {
		panic(err)
	}
	u.Query().Set("appid", d.agent.CorpID)
	u.Query().Set("agentid", d.agent.AgentID)
	u.Query().Set("state", state)
	u.Query().Set("redirect_uri", service.GetAuthUrl())
	http.Redirect(w, r, u.String(), http.StatusFound)
}
func (d *QRAuthDriver) AuthRequest(service *auth.Service, r *http.Request) (*auth.Result, error) {
	return authRequestWithAgent(d.agent, service, r)
}
