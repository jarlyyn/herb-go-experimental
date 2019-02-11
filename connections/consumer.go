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

func Consume(i ConnectionsInput, c ConnectionsConsumer) {
	for {
		select {
		case m := <-i.Messages():
			c.OnMessage(m)
		case e := <-i.Errors():
			c.OnError(e)
		case conn := <-i.OnCloseEvents():
			c.OnClose(conn)
		case conn := <-i.OnOpenEvents():
			c.OnOpen(conn)
		}
	}
}
