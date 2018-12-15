package notification

type Sender interface {
	SendMessage(Message) error
}
