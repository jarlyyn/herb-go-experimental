package identifier

import (
	"sync"
	"time"

	"github.com/jarlyyn/herb-go-experimental/connections"
)

var GenerateDefaultMapOnLogout = func(m *Map) func(id string, conn *connections.Conn) error {
	return func(id string, conn *connections.Conn) error {
		if conn.RawConnection != nil {
			conn.Close()
			return nil
		}
		m.Connections.Close(conn.Info.ID)
		return nil
	}
}

type MapIdentity struct {
	Conn      *connections.Conn
	Timestamp int64
}

func NewMapIdentity() *MapIdentity {
	return &MapIdentity{
		Conn:      nil,
		Timestamp: time.Now().Unix(),
	}
}

type Map struct {
	Connections connections.Connections
	Identities  sync.Map
	lock        sync.Mutex
	onLogout    func(id string, conn *connections.Conn) error
}

func (m *Map) conn(id string) (*connections.Conn, bool) {
	data, ok := m.Identities.Load(id)
	if ok == false {
		return nil, false
	}
	conn, ok := data.(*connections.Conn)
	return conn, ok
}
func (m *Map) Login(id string, conn *connections.Conn) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	conn, ok := m.conn(id)
	if ok {
		err := m.onLogout(id, conn)
		if err != nil {
			return err
		}
	}
	m.Identities.Delete(id)
	m.Identities.Store(id, conn)
	return nil
}
func (m *Map) Logout(id string, c *connections.Conn) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	conn, ok := m.conn(id)
	if ok {
		if c != nil && conn.Info.ID != c.Info.ID {
			return nil
		}
		err := m.onLogout(id, conn)
		if err != nil {
			return err
		}
	}
	m.Identities.Delete(id)
	return nil
}
func (m *Map) Verify(id string, conn *connections.Conn) (bool, error) {
	conn, ok := m.conn(id)
	if ok == false {
		return false, nil
	}
	return conn.Info.ID == id, nil
}
func (m *Map) SendByID(id string, msg []byte) error {
	conn, ok := m.conn(id)
	if ok == false {
		return nil
	}
	return conn.Send(msg)
}
func (m *Map) OnLogout() func(id string, conn *connections.Conn) error {
	return m.onLogout
}
func (m *Map) SetOnLogout(f func(id string, conn *connections.Conn) error) {
	m.onLogout = f
}

func NewMap() *Map {

	m := &Map{
		Connections: nil,
		Identities:  sync.Map{},
	}
	m.onLogout = GenerateDefaultMapOnLogout(m)
	return m
}
