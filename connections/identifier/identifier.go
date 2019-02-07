package identifier

import "github.com/jarlyyn/herb-go-experimental/connections"

type Identifier interface {
	Login(id string, conn *connections.Conn) error
	Logout(conn *connections.Conn) error
	Verify(id string, connid *connections.Conn) (bool, error)
	SendByID(id string, msg []byte) error
	OnLogout() func(id string, connid *connections.Conn)
	SetOnLogout(func(id string, connid *connections.Conn))
}
