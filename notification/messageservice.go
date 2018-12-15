package notification

const (
	MessageInstanceStatusNew         = 0
	MessageInstanceStatusSuccess     = 1
	MessageInstanceStatusFail        = 2
	MessageInstanceStatusPending     = 3
	MessageInstanceStatusUnsupported = 10
)

type MessageFactory func() (Message, error)

type MessageInstance struct {
	Message    Message
	Sender     string
	InstanceID string
	Status     int
	Logs       []string
}

func NewMessageInstance(message Message) *MessageInstance {
	return &MessageInstance{
		Message: message,
		Status:  MessageInstanceStatusNew,
		Logs:    []string{},
	}
}

type MessageService interface {
	RegiterMessageType(string, MessageFactory)
	NewMessageInstance(instancetype string) (*MessageInstance, error)
}
