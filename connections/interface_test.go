package connections

import "testing"

func TestInterface(t *testing.T) {
	var c ConnectionOutput
	c = New()
	id := c.ID()
	if id != "" {
		t.Error(id)
	}
}
