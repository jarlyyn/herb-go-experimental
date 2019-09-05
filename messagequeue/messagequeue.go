package messagequeue

type ConsumerStatus int

const ConsumerStatusSuccess = ConsumerStatus(0)
const ConsumerStatusFail = ConsumerStatus(-1)

type Broker struct {
	Driver Driver
}

type Driver interface {
	Start() error
	Close() error
	ProduceMessages(...[]byte) (sent []bool, err error)
	SetConsumer(func([]byte) ConsumerStatus)
}
