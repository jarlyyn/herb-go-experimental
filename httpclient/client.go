package httpclient

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var ErrMsgLengthLimit = 512
var DefaultTimeout = 120

type Service struct {
	TimeoutInSecond int
}

var DefaultService = &Service{}

func (s *Service) Client() *Client {
	var cs *Service
	if s != nil {
		cs = s
	} else {
		cs = DefaultService
	}

	timeout := cs.TimeoutInSecond
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	c := http.Client{
		Timeout: time.Duration(cs.TimeoutInSecond) * time.Second,
	}
	return &Client{Client: &c}
}

func (s *Service) Fetch(req *http.Request) (*Result, error) {
	return s.Client().Fetch(req)
}

type Client struct {
	*http.Client
}

func (c *Client) Fetch(req *http.Request) (*Result, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := Result{
		Response:    resp,
		BodyContent: bodyContent,
	}
	return &result, nil
}

type Result struct {
	*http.Response
	BodyContent []byte
}

func (r *Result) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(r.BodyContent, &v)
}
func (r *Result) UnmarshalXML(v interface{}) error {
	return xml.Unmarshal(r.BodyContent, &v)
}
func (r Result) Error() string {
	msg := fmt.Sprintf("http error [ %s ] %s : %s", r.Response.Request.URL.String(), r.Status, string(r.BodyContent))
	if len(msg) > ErrMsgLengthLimit {
		msg = msg[:ErrMsgLengthLimit]
	}
	return msg
}

func (r *Result) NewAPICodeErr(code interface{}) *APICodeErr {
	return NewAPICodeErr(r.Response.Request.URL.String(), code, r.BodyContent)

}
func GetErrorStatusCode(err error) int {
	r, ok := err.(Result)
	if ok {
		return r.StatusCode
	}
	return 0
}

func NewAPICodeErr(url string, code interface{}, content []byte) *APICodeErr {
	return &APICodeErr{
		URI:     url,
		Code:    fmt.Sprint(code),
		Content: content,
	}
}

type APICodeErr struct {
	URI     string
	Code    string
	Content []byte
}

func (r APICodeErr) Error() string {
	msg := fmt.Sprintf("api error [ %s] code %d : %s", r.URI, r.Code, string(r.Content))
	if len(msg) > ErrMsgLengthLimit {
		msg = msg[:ErrMsgLengthLimit]
	}
	return msg
}

func GetAPIErrCode(err error) string {
	r, ok := err.(APICodeErr)
	if ok {
		return r.Code
	}
	return ""

}

func CompareApiErrCode(err error, code interface{}) bool {
	return GetAPIErrCode(err) == fmt.Sprint(code)
}
