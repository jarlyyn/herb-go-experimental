package websocketmanager

type Conn interface {
	Close() error
	Send(*Message) error
	Messages() chan *Message
	C() chan int
}
type ConnInfo struct {
	ID        string
	ManagerID string
	Timestamp int64
}

type RegisteredConn struct {
	Conn
	Info *ConnInfo
}
