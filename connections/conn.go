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
type Info struct {
	ID        string
	GatewayID string
	UID       string
	Timestamp int64
}

type Conn struct {
	RawConnection
	Info *Info
}

type Message struct {
	Message []byte
	Info    *Info
}

type Error struct {
	Error error
	Info  *Info
}
