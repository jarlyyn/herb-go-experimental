package notification

type NotificationFactory func() (Notification, error)

type NotificationService interface {
	RegiterNotificationType(string, NotificationFactory)
	NewNotificationInstance(instancetype string) (*NotificationInstance, error)
}
