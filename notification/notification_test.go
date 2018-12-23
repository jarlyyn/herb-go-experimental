package notification

import (
	"testing"
	"time"
)

type ChanSender struct {
	C chan *NotificationInstance
}

func (c *ChanSender) MustSendNotification(n *NotificationInstance) {
	c.C <- n
}

func newChanSender() *ChanSender {
	return &ChanSender{
		C: make(chan *NotificationInstance, 10),
	}
}

func TestSend(t *testing.T) {
	n, err := NewPartedNotificationWithID()
	if err != nil {
		t.Fatal(err)
	}
	err = n.SetNotificationAuthor("AuthorTest")
	if err != nil {
		t.Fatal(err)
	}
	err = n.SetNotificationType("testtype")
	if err != nil {
		t.Fatal(err)
	}
	err = n.SetNotificationRecipient("RecipientTest")
	if err != nil {
		t.Fatal(err)
	}
	unusedSender := newChanSender()
	sender := newChanSender()
	DefaultService.RegisterSender("testtype", sender)
	err = DefaultService.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer DefaultService.Stop()
	if len(sender.C) != 0 {
		t.Error(sender.C)
	}
	if len(unusedSender.C) != 0 {
		t.Error(sender.C)
	}
	Notify(n)
	time.Sleep(time.Microsecond)
	if len(sender.C) != 1 {
		t.Error(sender.C)
	}
	if len(unusedSender.C) != 0 {
		t.Error(sender.C)
	}
	ni := <-sender.C
	niNotificationID, err := ni.Notification.NotificationID()
	if err != nil {
		t.Fatal(err)
	}
	niNotificationAuthor, err := ni.Notification.NotificationAuthor()
	if err != nil {
		t.Fatal(err)
	}
	niNotificationRecipient, err := ni.Notification.NotificationRecipient()
	if err != nil {
		t.Fatal(err)
	}
	niNotificationType, err := ni.Notification.NotificationType()
	if err != nil {
		t.Fatal(err)
	}
	nNotificationID, err := n.NotificationID()
	if err != nil {
		t.Fatal(err)
	}
	nNotificationAuthor, err := n.NotificationAuthor()
	if err != nil {
		t.Fatal(err)
	}
	nNotificationRecipient, err := n.NotificationRecipient()
	if err != nil {
		t.Fatal(err)
	}
	nNotificationType, err := n.NotificationType()
	if err != nil {
		t.Fatal(err)
	}
	if niNotificationID != nNotificationID ||
		niNotificationAuthor != nNotificationAuthor ||
		niNotificationRecipient != nNotificationRecipient ||
		niNotificationType != nNotificationType {
		t.Error(*ni)
	}
}
