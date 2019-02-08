package connections

import (
	"net"
)

type RawConnection interface {
	Close() error
	Send([]byte) error
	Messages() chan []byte
	Errors() chan error
	RemoteAddr() net.Addr
	C() chan int
}

type OutputConn interface {
	Close() error
	Send([]byte) error
	ID() string
}
type Info struct {
	ID        string
	Timestamp int64
}

type Conn struct {
	RawConnection
	Info *Info
}

func New() *Conn {
	return &Conn{}
}
func (c *Conn) ID() string {
	if c.Info != nil {
		return c.Info.ID
	}
	return ""
}

type Message struct {
	Message []byte
	Conn    OutputConn
}

type Error struct {
	Error error
	Conn  OutputConn
}
