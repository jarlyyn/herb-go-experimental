package guarder

import "net/http"

type IDTokenHeaders struct {
	IDHeader    string
	TokenHeader string
}

func NewIDTokenHeaders() *IDTokenHeaders {
	return &IDTokenHeaders{}
}
func (h *IDTokenHeaders) ReadParamsFromRequest(r *http.Request) (*RequestParams, error) {
	p := NewRequestParams()
	if h.IDHeader != "" {
		p.SetID(r.Header.Get(h.IDHeader))
	}
	p.SetToken(r.Header.Get(h.TokenHeader))
	return p, nil
}
func (h *IDTokenHeaders) WriteParamsToRequest(r *http.Request, p *RequestParams) error {
	if h.IDHeader != "" {
		r.Header.Set(h.IDHeader, p.ID())
	}
	r.Header.Set(h.TokenHeader, p.Token())
	return nil
}

func createIDTokenHeadersWithConfig(conf Config, prefix string) (*IDTokenHeaders, error) {
	var err error
	v := NewIDTokenHeaders()
	if err != nil {
		return nil, err
	}
	err = conf.Get("IDHeader", &v.IDHeader)
	if err != nil {
		return nil, err
	}
	err = conf.Get("TokenHeader", &v.TokenHeader)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func idTokenHeadersMapperFactory(conf Config, prefix string) (RequestParamsMapper, error) {
	return createIDTokenHeadersWithConfig(conf, prefix)
}

func registerIDTokenHeadersFactories() {
	RegisterMapper("headers", idTokenHeadersMapperFactory)
}

func init() {
	registerIDTokenHeadersFactories()
}
