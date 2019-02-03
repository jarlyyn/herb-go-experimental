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
	messages    chan *ConnMessage
	errors      chan *ConnError
}

type ConnMessage struct {
	*Message
	Info *ConnInfo
}

type ConnError struct {
	Error error
	Info  *ConnInfo
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
	go func() {
		defer func() {
			m.Registered.Delete(r.Info.ID)
		}()
		for {
			select {
			case message := <-conn.Messages():
				m.messages <- &ConnMessage{
					Message: message,
					Info:    r.Info,
				}
			case err := <-conn.Errors():
				m.errors <- &ConnError{
					Error: err,
					Info:  r.Info,
				}
			case <-conn.C():
				break
			}
		}
	}()
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
