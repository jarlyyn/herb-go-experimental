package identifier

import "github.com/jarlyyn/herb-go-experimental/connections"

type Identifier interface {
	Login(id string, conn connections.OutputConn) error
	Logout(id string, conn connections.OutputConn) error
	Verify(id string, conn connections.OutputConn) (bool, error)
	SendByID(id string, msg []byte) error
	OnLogout() func(id string, conn connections.OutputConn) error
	SetOnLogout(func(id string, conn connections.OutputConn) error)
}
