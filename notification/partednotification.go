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

func NewPartedNotificationWithID() (*PartedNotification, error) {
	id, err := DefaultIDGenerator()
	if err != nil {
		return nil, err
	}
	n := NewPartedNotification()
	n.SetNotificationID(id)
	return n, nil
}
