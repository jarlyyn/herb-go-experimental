package connections

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

func NewGateway() *Gateway {
	return &Gateway{
		IDGenerator: DefaultIDGenerator,
		Messages:    make(chan *Message),
		Errors:      make(chan *Error),
	}
}

type Gateway struct {
	ID          string
	IDGenerator func() (string, error)
	Connections sync.Map
	Messages    chan *Message
	Errors      chan *Error
}

func (m *Gateway) Register(conn RawConnection) (*Conn, error) {
	id, err := m.IDGenerator()
	if err != nil {
		return nil, err
	}
	r := &Conn{
		RawConnection: conn,
		Info: &Info{
			ID:        id,
			GatewayID: m.ID,
			Timestamp: time.Now().Unix(),
		},
	}
	go func() {
		defer func() {
			m.Connections.Delete(r.Info.ID)
		}()
		for {
			select {
			case message := <-conn.Messages():
				m.Messages <- &Message{
					Message: message,
					Info:    r.Info,
				}
			case err := <-conn.Errors():
				m.Errors <- &Error{
					Error: err,
					Info:  r.Info,
				}
			case <-conn.C():
				break
			}
		}
	}()
	m.Connections.Store(id, r)
	return r, nil
}
func (m *Gateway) Conn(id string) *Conn {
	val, ok := m.Connections.Load(id)
	if ok == false {
		return nil
	}
	r := val.(*Conn)
	return r
}
func (m *Gateway) Send(id string, msg []byte) error {
	c := m.Conn(id)
	if c == nil {
		return nil
	}
	return c.Send(msg)
}
