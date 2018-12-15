package notification

type PartedNotification struct {
	CommonNotification
	Parts map[string][]byte
}

func NewPartedNotification() *PartedNotification {
	return &PartedNotification{
		Parts: map[string][]byte{},
	}
}
