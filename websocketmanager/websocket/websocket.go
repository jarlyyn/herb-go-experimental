package websocket

import (
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jarlyyn/herb-go-experimental/websocketmanager"
)

var ErrMsgTypeNotMatch = errors.New("websocket message type not match")

type Conn struct {
	*websocket.Conn
	closed      bool
	messageType int
	closelocker sync.Mutex
	messages    chan *websocketmanager.Message
	errors      chan error
	c           chan int
}

func (c *Conn) C() chan int {
	return c.c
}
func (c *Conn) Messages() chan *websocketmanager.Message {
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

func (c *Conn) Send(m *websocketmanager.Message) error {
	c.closelocker.Lock()
	closed := c.closed
	c.closelocker.Unlock()
	if closed {
		return nil
	}
	return c.Conn.WriteMessage(c.messageType, []byte(*m))
}
func New() *Conn {
	return &Conn{
		closed: true,
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
			m := websocketmanager.Message(msg)
			c.messages <- &m
		}
	}()
	return c, nil
}
