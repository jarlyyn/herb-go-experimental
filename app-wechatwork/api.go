package wechatwork

import (
	"github.com/jarlyyn/herb-go-experimental/httpclient"
)

var Server = client.Server{
	Host: "https://qyapi.weixin.qq.com",
}

var apiGetUserInfo = Server.CreateApi("GET", "/cgi-bin/user/getuserinfo")
var apiGetToken = Server.CreateApi("GET", "/cgi-bin/gettoken")
var apiGetUserDetail = Server.CreateApi("POST", "/cgi-bin/user/getuserdetail")
var apiUserGet = Server.CreateApi("GET", "/cgi-bin/user/get")

const ApiErrAccessTokenWrong = 40014
const ApiErrAccessTokenOutOfDate = 42001
const ApiErrSuccess = 0
const ApiErrUserUnaccessible = 50002

const ApiResultGenderMale = "1"
const ApiResultGenderFemale = "2"

type resultApiError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type resultAccessToken struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

type resultUserInfo struct {
	UserID     string `json:"UserId"`
	UserTicket string `json:"user_ticket"`
}
type paramsUserDetail struct {
	UserTicket string `json:"user_ticket"`
}
type resultUserDetail struct {
	UserID   string `json:"userid"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Mobile   string `json:"mobile"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type resultUserGet struct {
	UserID   string `json:"userid"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Mobile   string `json:"mobile"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
type apiError struct {
	Code int
	Body string
}

func (e apiError) Error() string {
	return e.Body
}

func newApiError(code int, body string) *apiError {
	return &apiError{
		Code: code,
		Body: body,
	}
}

func getApiErrCode(err error) int {
	apierr, ok := err.(apiError)
	if ok {
		return apierr.Code
	}
	return 0
}
