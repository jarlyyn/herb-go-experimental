package notificationmanager

import (
	"github.com/herb-go/herb/notification"
)

type DataNotification interface {
	NotificationData() (interface{}, error)
	notification.Notification
}
type DataInterfaceNotification struct {
	Data interface{}
	notification.Notification
}

func (n *DataInterfaceNotification) NotificationData() (interface{}, error) {
	return n.Data, nil
}
