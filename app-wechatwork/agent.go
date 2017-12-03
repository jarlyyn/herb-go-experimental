package wechatwork

import (
	"net/http"
	"net/url"
	"sync"

	client "github.com/jarlyyn/herb-go-experimental/httpclient"
)

type Agent struct {
	CorpID        string
	AgentID       string
	Secret        string
	ClientService client.Service
	accessToken   string
	lock          sync.Mutex
}

func (a *Agent) AccessToken() string {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.accessToken
}

func (a *Agent) GrantAccessToken() error {
	params := url.Values{}
	params.Set("corpid", a.CorpID)
	params.Set("corpsecret", a.Secret)
	req, err := apiGetToken.NewRequest(params, nil)
	if err != nil {
		return err
	}
	rep, err := a.ClientService.Fetch(req)
	if err != nil {
		return err
	}
	if rep.StatusCode != http.StatusOK {
		return rep
	}
	result := &resultAccessToken{}
	err = rep.UnmarshalJSON(&result)
	if err != nil {
		return err
	}
	if result.Errcode != 0 || result.Errmsg == "" || result.Access_token == "" {
		return rep
	}
	a.accessToken = result.Access_token
	return nil
}

func (a *Agent) CallApiWithAccessToken(api *client.Api, params url.Values, body interface{}, v interface{}) error {
	var apierr resultApiError
	var err error
	if a.AccessToken() == "" {
		err := a.GrantAccessToken()
		if err != nil {
			return err
		}
	}
	p := url.Values{}
	if params != nil {
		for k, vs := range params {
			for _, v := range vs {
				p.Add(k, v)
			}
		}
	}
	p.Set("access_token", a.AccessToken())
	req, err := api.NewJSONRequest(p, body)
	if err != nil {
		return err
	}
	resp, err := a.ClientService.Fetch(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return resp
	}
	apierr = resultApiError{}
	err = resp.UnmarshalJSON(&apierr)
	if err != nil {
		return err
	}
	if apierr.Errcode == ApiErrAccessTokenOutOfDate || apierr.Errcode == ApiErrAccessTokenWrong {
		err := a.GrantAccessToken()
		if err != nil {
			return err
		}
		p.Set("access_token", a.AccessToken())
		req, err := api.NewJSONRequest(p, body)
		if err != nil {
			return err
		}
		resp, err := a.ClientService.Fetch(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return resp
		}
		apierr = resultApiError{}
		err = resp.UnmarshalJSON(&apierr)
		if err != nil {
			return err
		}
	}
	if apierr.Errcode != 0 {
		return newApiError(apierr.Errcode, string(resp.BodyContent))
	}
	return resp.UnmarshalJSON(&v)
}

type Userinfo struct {
	UserID string
	Name   string
	Mobile string
	Email  string
	Gender string
	Avatar string
}

func (a *Agent) GetUserInfo(code string) (*Userinfo, error) {
	var info = &Userinfo{}
	if code == "" {
		return nil, nil
	}
	var result = &resultUserInfo{}
	params := url.Values{}
	params.Set("code", code)
	err := a.CallApiWithAccessToken(apiGetUserInfo, params, nil, result)
	if err != nil {
		return nil, err
	}
	if result.UserID == "" {
		return nil, nil
	}
	var getuser = &resultUserGet{}
	userGetParam := url.Values{}
	userGetParam.Add("userid", result.UserID)
	err = a.CallApiWithAccessToken(apiUserGet, userGetParam, nil, getuser)
	if err != nil {
		if getApiErrCode(err) == ApiErrUserUnaccessible {
			return nil, nil
		}
		return nil, err
	}
	info.UserID = result.UserID
	info.Avatar = getuser.Avatar
	info.Email = getuser.Email
	info.Gender = getuser.Gender
	info.Mobile = getuser.Mobile
	info.Name = getuser.Name
	return info, nil
}
