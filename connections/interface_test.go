package connections

import "testing"

func TestInterface(t *testing.T) {
	var c OutputConn
	c = New()
	id := c.ID()
	if id != "" {
		t.Error(id)
	}
}
