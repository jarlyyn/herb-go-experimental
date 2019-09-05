package messagequeue

type Message struct {
	ID   string
	Type string
	Data []byte
}

type Consumer func(*Message) bool

type Broker struct {
	Driver Driver
}
type Driver interface {
	Start() error
	Close() error
	ProduceMessages(...*Message) (unsend []*Message, err error)
	SetConsumer(Consumer)
}
