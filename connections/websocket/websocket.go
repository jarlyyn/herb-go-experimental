package websocket

import (
	"errors"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const MsgTypeText = websocket.TextMessage
const MsgTypeBinary = websocket.BinaryMessage

var ErrMsgTypeNotMatch = errors.New("websocket message type not match")

type Conn struct {
	*websocket.Conn
	closed      bool
	messageType int
	closelocker sync.Mutex
	messages    chan []byte
	errors      chan error
	c           chan int
}

func (c *Conn) C() chan int {
	return c.c
}
func (c *Conn) Messages() chan []byte {
	return c.messages
}
func (c *Conn) Errors() chan error {
	return c.errors
}
func (c *Conn) Close() error {
	defer c.closelocker.Unlock()
	c.closelocker.Lock()
	if c.closed {
		return nil
	}
	close(c.c)
	c.closed = true
	return c.Conn.Close()
}

func (c *Conn) Send(m []byte) error {
	c.closelocker.Lock()
	closed := c.closed
	c.closelocker.Unlock()
	if closed {
		return nil
	}
	return c.Conn.WriteMessage(c.messageType, m)
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func New() *Conn {
	return &Conn{
		closed:   true,
		c:        make(chan int),
		messages: make(chan []byte),
		errors:   make(chan error),
	}
}

var upgrader = websocket.Upgrader{} // use default options

func Upgrade(w http.ResponseWriter, r *http.Request, msgtype int) (*Conn, error) {
	wc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	c := New()
	c.closed = false
	c.Conn = wc
	c.messageType = msgtype
	go func() {
		defer func() {

		}()
		defer func() {
			recover()
		}()
		for {
			mt, msg, err := c.ReadMessage()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.closelocker.Lock()
				closed := c.closed
				c.closelocker.Unlock()
				if closed {
					break
				}
				if websocket.IsUnexpectedCloseError(err) {
					c.Close()
					break
				}
				c.errors <- err
				continue
			}
			if mt != c.messageType {
				c.errors <- ErrMsgTypeNotMatch
				continue
			}
			c.messages <- msg
		}
	}()
	return c, nil
}
