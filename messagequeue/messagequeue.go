package messagequeue

var EmptyConsumed = func(bool) error {
	return nil
}

type Message struct {
	ID   string
	Type string
	Data []byte
}

type ConsumerMessage struct {
	Message
	Consumed func(bool) error
}

type Broker struct {
	Driver
}
type Driver interface {
	Start() error
	Close() error
	ProduceMessages(...*Message) error
	MessageCuonsumer() chan *ConsumerMessage
}
