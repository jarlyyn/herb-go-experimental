package notificationmanager

import (
	"github.com/herb-go/notification"
)

type DataNotification interface {
	NotificationData() (interface{}, error)
	notification.Notification
}
type DataInterfaceNotification struct {
	Data interface{}
	notification.CommonNotification
}

func (n *DataInterfaceNotification) NotificationData() (interface{}, error) {
	return n.Data, nil
}

func NewDataInterfaceNotification(data interface{}) *DataInterfaceNotification {
	return &DataInterfaceNotification{
		Data:               data,
		CommonNotification: notification.CommonNotification{},
	}
}
