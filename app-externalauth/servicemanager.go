package auth

type ServiceManager interface {
	GetService(auth *Auth, keyword string) (*Service, error)
	RegisterService(auth *Auth, keyword string, driver Driver) (*Service, error)
}

type MapServiceManager struct {
	Services map[string]*Service
}

func (m *MapServiceManager) GetService(a *Auth, keyword string) (*Service, error) {
	s, ok := m.Services[keyword]
	if ok {
		return s, nil
	}
	return nil, nil

}

func (m *MapServiceManager) RegisterService(a *Auth, keyword string, driver Driver) (*Service, error) {
	s := &Service{
		Driver:  driver,
		Auth:    a,
		Keyword: keyword,
	}
	m.Services[keyword] = s
	return s, nil

}

var DefaultServiceManager = &MapServiceManager{}
