package notification

type Notification interface {
	NotificationID() (string, error)
	SetNotificationID(string) error
	NotificationType() (string, error)
	SetNotificationType(string) error
	NotificationRecipient() (string, error)
	SetNotificationRecipient(string) error
}

type CommonNotification struct {
	ID        string
	Type      string
	Recipient string
}

func (m *CommonNotification) NotificationID() (string, error) {
	return m.ID, nil
}

func (m *CommonNotification) SetNotificationID(id string) error {
	m.ID = id
	return nil
}

func (m *CommonNotification) NotificationType() (string, error) {
	return m.Type, nil
}

func (m *CommonNotification) SetNotificationType(mesasgetype string) error {
	m.Type = mesasgetype
	return nil
}

func (m *CommonNotification) NotificationRecipient() (string, error) {
	return m.Recipient, nil
}

func (m *CommonNotification) SetNotificationRecipient(recipient string) error {
	m.Recipient = recipient
	return nil
}
