package websocket

import "testing"
import "github.com/jarlyyn/herb-go-experimental/websocketmanager"

func TestInterface(t *testing.T) {
	var c websocketmanager.Conn = New()
	c.Close()
}
