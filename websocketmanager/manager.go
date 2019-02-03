package websocketmanager

import (
	"sync"
	"time"
)

func New() *Manager {
	return &Manager{}
}

type Manager struct {
	ID          string
	IDGenerator func() (string, error)
	Registered  sync.Map
}

func (m *Manager) Register(conn Conn) (*RegisteredConn, error) {
	id, err := m.IDGenerator()
	if err != nil {
		return nil, err
	}
	r := &RegisteredConn{
		Conn: conn,
		Info: &ConnInfo{
			ID:        id,
			ManagerID: m.ID,
			Timestamp: time.Now().Unix(),
		},
	}
	m.Registered.Store(id, r)
	return r, nil
}

func (m *Manager) ConnByID(id string) *RegisteredConn {
	val, ok := m.Registered.Load(id)
	if ok == false {
		return nil
	}
	r := val.(*RegisteredConn)
	return r
}
func (m *Manager) SendByID(id string, msg *Message) error {
	c := m.ConnByID(id)
	if c == nil {
		return nil
	}
	return c.Send(msg)
}

func (m *Manager) OnClose(r *RegisteredConn) {
	m.Registered.Delete(r.Info.ID)
}
