package messagequeue

import (
	"sync"
	"testing"
)

func newTestBroker() *Broker {
	b := NewBroker()
	c := NewOptionConfigMap()
	c.Driver = "chan"
	c.ApplyTo(b)
	return b
}

func TestBrokerOrderedMessages(t *testing.T) {
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
	received := [][]byte{}
	wg := &sync.WaitGroup{}
	l := &sync.Mutex{}
	b.SetConsumer(func(m []byte) ConsumerStatus {
		l.Lock()
		defer l.Unlock()
		received = append(received, m)
		wg.Done()
		return ConsumerStatusSuccess
	})
	messages := [][]byte{}
	for i := byte(0); i < 5; i++ {
		wg.Add(1)
		messages = append(messages, []byte{i})
	}
	go func() {
		sent, err := b.ProduceMessages(messages...)
		if err != nil {
			t.Fatal(err)
		}
		for k := range sent {
			if sent[k] == false {
				t.Fatal(k)
			}
		}
	}()
	wg.Wait()
	if len(received) != 5 {
		t.Fatal(len(received))
	}
	for i := byte(0); i < 5; i++ {
		b := received[i]
		if len(b) != 1 || b[0] != i {
			t.Fatal(i, b)
		}
	}
}
