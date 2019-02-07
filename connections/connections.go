package connections

type Connections interface {
	Send(id string, msg []byte)
	Close(id string) error
}
