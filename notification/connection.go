package notification

type NotificationConnection interface {
	SetInput(NotificationInput)
	SetOutput(NotificationOutput)
}

type CommonConnection struct {
	input  NotificationInput
	output NotificationOutput
}

func (c *CommonConnection) SetInput(i NotificationInput) {
	c.input = i
}

func (c *CommonConnection) SetOnput(i NotificationOutput) {
	c.output = i
}

type NotificationInput interface {
	NotificationChanIn() chan *Notification
}

type NotificationOutput interface {
	NotificationChanOut() chan *Notification
}
