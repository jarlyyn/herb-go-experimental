package websocketmanager

import (
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

var DefaultIDGenerator = func() (string, error) {
	unid, err := uuid.NewV1()
	if err != nil {
		return "", err
	}
	return unid.String(), nil
}

type ConnMessage struct {
	Message []byte
	Info    *ConnInfo
}

type ConnError struct {
	Error error
	Info  *ConnInfo
}

func New() *Manager {
	return &Manager{
		IDGenerator:  DefaultIDGenerator,
		ConnMessages: make(chan *ConnMessage),
		ConnErrors:   make(chan *ConnError),
	}
}

type Manager struct {
	ID           string
	IDGenerator  func() (string, error)
	Registered   sync.Map
	ConnMessages chan *ConnMessage
	ConnErrors   chan *ConnError
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
				m.ConnMessages <- &ConnMessage{
					Message: message,
					Info:    r.Info,
				}
			case err := <-conn.Errors():
				m.ConnErrors <- &ConnError{
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
func (m *Manager) SendByID(id string, msg []byte) error {
	c := m.ConnByID(id)
	if c == nil {
		return nil
	}
	return c.Send(msg)
}
