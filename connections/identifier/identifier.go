package identifier

import "github.com/jarlyyn/herb-go-experimental/connections"

type Identifier interface {
	Login(id string, conn *connections.Conn) error
	Logout(id string, conn *connections.Conn) error
	Verify(id string, conn *connections.Conn) (bool, error)
	SendByID(id string, msg []byte) error
	OnLogout() func(conn *connections.Conn)
	SetOnLogout(func(conn *connections.Conn))
}
