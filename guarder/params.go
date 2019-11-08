package guarder

type ParamsKey string

func (k ParamsKey) LoadFrom(p *Params) string {
	return (*p)[string(k)]
}

func (k ParamsKey) SetTo(p *Params, v string) {
	(*p)[string(k)] = v
}

const ParamsKeyID = ParamsKey("id")
const ParamsKeyToken = ParamsKey("token")

type Params map[string]string

func (p *Params) ID() string {
	return ParamsKeyID.LoadFrom(p)
}

func (p *Params) SetID(v string) {
	ParamsKeyID.SetTo(p, v)
}

func (p *Params) Token() string {
	return ParamsKeyToken.LoadFrom(p)
}

func (p *Params) SetToken(v string) {
	ParamsKeyToken.SetTo(p, v)
}

func NewParams() *Params {
	return &Params{}
}
