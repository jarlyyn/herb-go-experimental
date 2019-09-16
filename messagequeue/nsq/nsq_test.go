package nsq

import (
	"bytes"
	"container/list"
	"encoding/json"
	"testing"
	"time"

	"github.com/jarlyyn/herb-go-experimental/messagequeue"
)

func newTestBroker() *messagequeue.Broker {
	b := messagequeue.NewBroker()
	c := messagequeue.NewOptionConfigMap()
	err := json.Unmarshal([]byte(testConfig), c)
	if err != nil {
		panic(err)
	}
	err = c.ApplyTo(b)
	if err != nil {
		panic(err)
	}
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
	testchan := make(chan *messagequeue.Message, 100)
	b.SetConsumer(messagequeue.NewChanConsumer(testchan))
	messages := [][]byte{}
	unreceived := list.New()
	for i := byte(0); i < 5; i++ {
		messages = append(messages, []byte{i})
		unreceived.PushBack([]byte{i})
	}
	sent, err := b.ProduceMessages(messages...)
	if err != nil {
		t.Fatal(err)
	}
	for k := range sent {
		if sent[k] == false {
			t.Fatal(k)
		}
	}
	time.Sleep(1 * time.Second)
	if len(testchan) != 5 {
		t.Fatal(len(testchan))
	}
	if unreceived.Len() != 5 {
		t.Fatal(unreceived.Len())
	}
	for i := byte(0); i < 5; i++ {
		m := <-testchan
		e := unreceived.Front()
		for e != nil {
			if bytes.Compare(e.Value.([]byte), m.Data) == 0 {
				unreceived.Remove(e)
				break
			}
			e = e.Next()
		}
	}
	if unreceived.Len() != 0 {
		t.Fatal(unreceived)
	}
}

func newLookupTestBroker() *messagequeue.Broker {
	b := messagequeue.NewBroker()
	c := messagequeue.NewOptionConfigMap()
	err := json.Unmarshal([]byte(testLookupConfig), c)
	if err != nil {
		panic(err)
	}
	err = c.ApplyTo(b)
	if err != nil {
		panic(err)
	}
	return b
}

func TestLookupBroker(t *testing.T) {
	b := newLookupTestBroker()
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
	testchan := make(chan *messagequeue.Message, 100)
	b.SetConsumer(messagequeue.NewChanConsumer(testchan))
	messages := [][]byte{}
	unreceived := list.New()
	for i := byte(0); i < 5; i++ {
		messages = append(messages, []byte{i})
		unreceived.PushBack([]byte{i})
	}
	sent, err := b.ProduceMessages(messages...)
	if err != nil {
		t.Fatal(err)
	}
	for k := range sent {
		if sent[k] == false {
			t.Fatal(k)
		}
	}
	time.Sleep(1 * time.Second)
	if len(testchan) != 5 {
		t.Fatal(len(testchan))
	}
	if unreceived.Len() != 5 {
		t.Fatal(unreceived.Len())
	}
	for i := byte(0); i < 5; i++ {
		m := <-testchan
		e := unreceived.Front()
		for e != nil {
			if bytes.Compare(e.Value.([]byte), m.Data) == 0 {
				unreceived.Remove(e)
				break
			}
			e = e.Next()
		}
	}
	if unreceived.Len() != 0 {
		t.Fatal(unreceived)
	}
}
