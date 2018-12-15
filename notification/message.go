package notification

type Message interface {
	MessageID() (string, error)
	SetMessageID(string) error
	MessageType() (string, error)
	SetMessageType(string) error
	MessageRecipient() (string, error)
	SetMessageRecipient(string) error
}

type CommonMessage struct {
	ID        string
	Type      string
	Recipient string
}

func (m *CommonMessage) MessageID() (string, error) {
	return m.ID, nil
}

func (m *CommonMessage) SetMessageID(id string) error {
	m.ID = id
	return nil
}

func (m *CommonMessage) MessageType() (string, error) {
	return m.Type, nil
}

func (m *CommonMessage) SetMessageType(mesasgetype string) error {
	m.Type = mesasgetype
	return nil
}

func (m *CommonMessage) MessageRecipient() (string, error) {
	return m.Recipient, nil
}

func (m *CommonMessage) SetMessageRecipient(recipient string) error {
	m.Recipient = recipient
	return nil
}
