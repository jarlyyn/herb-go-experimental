package notification

type PartedMessage struct {
	CommonMessage
	Parts map[string][]byte
}

func NewPartedMessage() *PartedMessage {
	return &PartedMessage{
		Parts: map[string][]byte{},
	}
}
