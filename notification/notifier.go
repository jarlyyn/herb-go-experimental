package notification

type NotificationFactory func() (Notification, error)

type Service interface {
	RegiterNotificationType(name string, f NotificationFactory) error
	NewNotificationInstance(n Notification) (*NotificationInstance, error)
	RegisterSender(notificationtype string, s Sender)
	SendersByType(notificationtype string) ([]Sender, error)
	Recover() func()
	SetRecover(func())
	Start() error
	Stop() error
	Notifier
}

var DefaultRecover = func() {
}

type CommonService struct {
	registeredFactories map[string]NotificationFactory
	registeredSender    map[string][]Sender
	idGenerator         IDGenerator
	c                   chan Notification
	closeChan           chan bool
	recover             func()
}

func NewCommonService() *CommonService {
	return &CommonService{
		registeredFactories: map[string]NotificationFactory{},
		registeredSender:    map[string][]Sender{},
		idGenerator:         DefaultIDGenerator,
		closeChan:           make(chan bool),
		recover:             DefaultRecover,
	}
}
func (m *CommonService) Recover() func() {
	return m.recover
}
func (m *CommonService) SetRecover(r func()) {
	m.recover = r
}
func (m *CommonService) Start() error {
	c := m.NotificationChan()
	go func() {
		for {
			select {
			case n := <-c:
				go func() {
					defer m.Recover()()
					err := SendNotificationByService(m, n)
					if err != nil {
						panic(err)
					}
				}()
			}
		}
	}()
	return nil
}

func (m *CommonService) Stop() error {
	return nil
}
func (m *CommonService) NotificationChan() chan Notification {
	return m.c
}
func (m *CommonService) RegiterNotificationType(name string, f NotificationFactory) error {
	m.registeredFactories[name] = f
	return nil
}
func (m *CommonService) NewNotificationInstance(n Notification) (*NotificationInstance, error) {
	id, err := m.idGenerator()
	if err != nil {
		return nil, err
	}
	i := NewNotificationInstance(n)
	i.InstanceID = id
	return i, nil
}
func (m *CommonService) RegisterSender(notificationtype string, s Sender) {
	_, ok := m.registeredSender[notificationtype]
	if ok == false {
		m.registeredSender[notificationtype] = []Sender{s}
	} else {
		m.registeredSender[notificationtype] = append(m.registeredSender[notificationtype], s)
	}
}

func (m *CommonService) SendersByType(notificationtype string) ([]Sender, error) {
	return m.registeredSender[notificationtype], nil
}

var SendNotificationByService = func(m Service, n Notification) error {
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
	return SendNotificationByService(DefaultService, n)
}

var DefaultService = NewCommonService()

type Notifier interface {
	NotificationChan() chan Notification
}

var DefaultNotifier = DefaultService

var NotifyTo = func(notifier Notifier, n Notification) {
	go func() {
		notifier.NotificationChan() <- n
	}()
}

var Notify = func(n Notification) {
	NotifyTo(DefaultNotifier, n)
}
