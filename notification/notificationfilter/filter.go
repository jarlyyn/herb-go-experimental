package notificationfilter

import "github.com/jarlyyn/herb-go-experimental/notification"

type Filter func(instance *notification.NotificationInstance, next func())

var FilterRecover = func(recovery func()) Filter {
	return func(instance *notification.NotificationInstance, next func()) {
		defer recovery()
		next()
	}
}

type RecipientConvertor func(recipient string) (string, error)

var FilterRecipient = func(convertor RecipientConvertor) Filter {
	return func(instance *notification.NotificationInstance, next func()) {
		recipient, err := instance.Notification.NotificationRecipient()
		if err != nil {
			panic(err)
		}
		id, err := convertor(recipient)
		if err != nil {
			panic(err)
		}
		err = instance.Notification.SetNotificationRecipient(id)
		if err != nil {
			panic(err)
		}
		next()
	}
}
