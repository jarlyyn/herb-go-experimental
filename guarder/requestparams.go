package guarder

type RequestParamsKey string

func (k RequestParamsKey) LoadFrom(p *RequestParams) string {
	return (*p)[string(k)]
}

func (k RequestParamsKey) SetTo(p *RequestParams, v string) {
	(*p)[string(k)] = v
}

const RequestParamsKeyID = RequestParamsKey("id")
const RequestParamsKeyToken = RequestParamsKey("token")

type RequestParams map[string]string

func (p *RequestParams) ID() string {
	return RequestParamsKeyID.LoadFrom(p)
}

func (p *RequestParams) SetID(v string) {
	RequestParamsKeyID.SetTo(p, v)
}

func (p *RequestParams) Token() string {
	return RequestParamsKeyToken.LoadFrom(p)
}

func (p *RequestParams) SetToken(v string) {
	RequestParamsKeyToken.SetTo(p, v)
}

func NewRequestParams() *RequestParams {
	return &RequestParams{}
}
