package auth

type ProviderManager interface {
	GetProvider(auth *Auth, keyword string) (*Provider, error)
	RegisterProvider(auth *Auth, keyword string, driver Driver) (*Provider, error)
}

type MapProviderManager struct {
	Providers map[string]*Provider
}

func NewMapProviderManager() *MapProviderManager {
	return &MapProviderManager{
		Providers: map[string]*Provider{},
	}
}
func (m *MapProviderManager) GetProvider(a *Auth, keyword string) (*Provider, error) {
	s, ok := m.Providers[keyword]
	if ok {
		return s, nil
	}
	return nil, nil

}

func (m *MapProviderManager) RegisterProvider(a *Auth, keyword string, driver Driver) (*Provider, error) {
	s := &Provider{
		Driver:  driver,
		Auth:    a,
		Keyword: keyword,
	}
	m.Providers[keyword] = s
	return s, nil

}

var DefaultProviderManager = NewMapProviderManager()
