package notificationmanager

import "sync"
import "github.com/herb-go/herb/notification"

type Gateway struct {
	notification.CommonInstancesBuilder
	Fields   map[string]*Field
	complied map[string]complied
	engine   Engine
	lock     sync.Mutex
}

func (m *Gateway) RenderTemplate(name string, data interface{}) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	var err error
	complied := m.complied[name]
	if complied == nil {
		f := m.Fields[name]
		if f == nil || f.Template == "" {
			return nil, nil
		}
		complied, err = m.engine.Compile(f.Template)
		if err != nil {
			return nil, err
		}
		m.complied[name] = complied
	}
	return complied(data)
}

func (m *Gateway) UpdateTemplate(name string, template string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	f := m.Fields[name]
	if f == nil {
		return ErrTemplateFieldNotRegistered
	}

	f.Template = template
	delete(m.complied, name)
	return nil
}
