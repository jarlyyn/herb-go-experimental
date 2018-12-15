package notification

type NotificationFactory func() (Notification, error)
type Notifier interface {
	RegiterNotificationType(name string, f NotificationFactory) error
	NewNotificationInstance(n Notification) (*NotificationInstance, error)
	RegisterSender(notificationtype string, s Sender)
	SendersByType(notificationtype string) ([]Sender, error)
	Recover() func()
	SetRecover(func())
	Start() error
	Stop() error
	NotificationChan() chan Notification
}

var DefaultRecover = func() {
}

type CommonNotifier struct {
	registeredFactories map[string]NotificationFactory
	registeredSender    map[string][]Sender
	idGenerator         IDGenerator
	c                   chan Notification
	closeChan           chan bool
	recover             func()
}

func NewCommonNotifier() *CommonNotifier {
	return &CommonNotifier{
		registeredFactories: map[string]NotificationFactory{},
		registeredSender:    map[string][]Sender{},
		idGenerator:         DefaultIDGenerator,
		closeChan:           make(chan bool),
		recover:             DefaultRecover,
	}
}
func (m *CommonNotifier) Recover() func() {
	return m.recover
}
func (m *CommonNotifier) SetRecover(r func()) {
	m.recover = r
}
func (m *CommonNotifier) Start() error {
	c := m.NotificationChan()
	go func() {
		for {
			select {
			case n := <-c:
				go func() {
					defer m.Recover()()
					err := SendNotificationByNotifier(m, n)
					if err != nil {
						panic(err)
					}
				}()
			}
		}
	}()
	return nil
}

func (m *CommonNotifier) Stop() error {
	return nil
}
func (m *CommonNotifier) NotificationChan() chan Notification {
	return m.c
}
func (m *CommonNotifier) RegiterNotificationType(name string, f NotificationFactory) error {
	m.registeredFactories[name] = f
	return nil
}
func (m *CommonNotifier) NewNotificationInstance(n Notification) (*NotificationInstance, error) {
	id, err := m.idGenerator()
	if err != nil {
		return nil, err
	}
	i := NewNotificationInstance(n)
	i.InstanceID = id
	return i, nil
}
func (m *CommonNotifier) RegisterSender(notificationtype string, s Sender) {
	_, ok := m.registeredSender[notificationtype]
	if ok == false {
		m.registeredSender[notificationtype] = []Sender{s}
	} else {
		m.registeredSender[notificationtype] = append(m.registeredSender[notificationtype], s)
	}
}

func (m *CommonNotifier) SendersByType(notificationtype string) ([]Sender, error) {
	return m.registeredSender[notificationtype], nil
}

var DefaultNotifier = NewCommonNotifier()

var SendNotificationByNotifier = func(m Notifier, n Notification) error {
	nt, err := n.NotificationType()
	if err != nil {
		return err
	}
	Senders, err := m.SendersByType(nt)
	if err != nil {
		return err
	}
	for k := range Senders {
		i, err := m.NewNotificationInstance(n)
		if err != nil {
			return err
		}
		go func() {
			Senders[k].MustSendNotification(i)
		}()
	}
	return nil
}

var SendNotification = func(n Notification) error {
	return SendNotificationByNotifier(DefaultNotifier, n)
}
