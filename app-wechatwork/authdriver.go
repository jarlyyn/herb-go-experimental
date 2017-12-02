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
const oauthUrl = "https://open.weixin.qq.com/connect/oauth2/authorize"

type Session struct {
	State string
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
	u, err := url.Parse(oauthUrl)
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
	var authsession = &Session{}
	var code = r.URL.Query().Get("code")
	if code == "" {
		return nil, auth.ErrAuthParamsError
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
	return nil, nil
}
