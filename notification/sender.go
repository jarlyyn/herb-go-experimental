package notification

type Sender interface {
	MustSendNotification(*NotificationInstance)
}
