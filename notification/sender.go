package notification

type Sender interface {
	SendNotification(NotificationInstance) error
}
