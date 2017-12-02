package client

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var ErrMsgLengthLimit = 1024
var DefaultTimeout = 120

type Service struct {
	TimeoutInSecond int
}

func (s *Service) Client() *Client {
	timeout := s.TimeoutInSecond
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	c := http.Client{
		Timeout: time.Duration(s.TimeoutInSecond) * time.Second,
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
	msg := fmt.Sprintf("http error %s : %s", r.StatusCode, string(r.BodyContent))
	return msg[:ErrMsgLengthLimit]
}

func GetErrorStatusCode(err error) int {
	r, ok := err.(Result)
	if ok {
		return r.StatusCode
	}
	return 0
}
