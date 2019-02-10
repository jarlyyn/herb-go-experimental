package connections

type ConnectionsConsumer interface {
	OnMessage(*Message)
	OnError(*Error)
	OnClose(ConnectionOutput)
	OnOpen(ConnectionOutput)
}

type EmptyConsumer struct {
}

func (e EmptyConsumer) OnMessage(*Message) {

}
func (e EmptyConsumer) OnError(*Error) {

}
func (e EmptyConsumer) OnClose(ConnectionOutput) {

}
func (e EmptyConsumer) OnOpen(ConnectionOutput) {

}
