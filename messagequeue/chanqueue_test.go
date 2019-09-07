package messagequeue

import "testing"

func newTestBroker() *Broker {
	b := NewBroker()
	c := NewOptionConfigMap()
	c.Driver = "chan"
	c.ApplyTo(b)
	return b
}

func TestBroker(t *testing.T) {
	b := newTestBroker()
	err := b.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := b.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

}
