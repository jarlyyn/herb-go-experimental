package github

import (
	"net/http"
	"net/url"

	"github.com/herb-go/herb/fetch"
)

type Client struct {
	ClientID     string
	ClientSecret string
	Fetcher      fetch.Fetcher
}

func (c *Client) GetAccessToken(code string) (*ResultAPIAccessToken, error) {
	params := url.Values{}
	params.Set("client_id", c.ClientID)
	params.Set("client_secret", c.ClientSecret)
	params.Set("code", code)

	req, err := apiAccessToken.NewRequest(nil, []byte(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	rep, err := c.Fetcher.Fetch(req)
	if err != nil {
		return nil, err
	}
	if rep.StatusCode != http.StatusOK {
		return nil, rep
	}
	result := &ResultAPIAccessToken{}
	err = rep.UnmarshalJSON(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Client) GetUser(accessToken string) (*ResultAPIUser, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)

	req, err := apiUser.NewRequest(params, nil)
	if err != nil {
		return nil, err
	}
	rep, err := c.Fetcher.Fetch(req)
	if err != nil {
		return nil, err
	}
	if rep.StatusCode != http.StatusOK {
		return nil, rep
	}
	result := &ResultAPIUser{}
	err = rep.UnmarshalJSON(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
