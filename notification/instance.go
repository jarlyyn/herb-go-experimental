package notification

const (
	NotificationInstanceStatusCanceled    = -1
	NotificationInstanceStatusNew         = 0
	NotificationInstanceStatusSuccess     = 1
	NotificationInstanceStatusFail        = 2
	NotificationInstanceStatusPending     = 3
	NotificationInstanceStatusUnsupported = 10
)

type NotificationInstance struct {
	Notification Notification
	Sender       string
	InstanceID   string
	Status       int
	Output       string
	Logs         []string
}

func NewNotificationInstance(Notification Notification) *NotificationInstance {
	return &NotificationInstance{
		Notification: Notification,
		Status:       NotificationInstanceStatusNew,
		Logs:         []string{},
	}
}
