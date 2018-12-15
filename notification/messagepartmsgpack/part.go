package messagepartmsgpack

import "github.com/jarlyyn/herb-go-experimental/notification"
import "github.com/vmihailenco/msgpack"

var Marshal = func(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}
var Unmarshal = func(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

type StringPart struct {
	Name string
}

func NewStringPart(name string) *StringPart {
	return &StringPart{
		Name: name,
	}
}

func (p *StringPart) Set(message *notification.PartedMessage, data string) error {
	message.Parts[p.Name] = []byte(data)
	return nil
}
func (p *StringPart) Get(message *notification.PartedMessage) (string, error) {
	bytes := message.Parts[p.Name]
	return string(bytes), nil
}

type BinaryPart struct {
	Name string
}

func NewBinaryPart(name string) *StringPart {
	return &StringPart{
		Name: name,
	}
}

func (p *BinaryPart) Set(message *notification.PartedMessage, data []byte) error {
	message.Parts[p.Name] = data
	return nil
}
func (p *BinaryPart) Get(message *notification.PartedMessage) ([]byte, error) {
	data := message.Parts[p.Name]
	return data, nil
}

type StringListPart struct {
	Name string
}

func NewStringListPart(name string) *StringListPart {
	return &StringListPart{
		Name: name,
	}
}

func (p *StringListPart) Set(message *notification.PartedMessage, data []string) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	message.Parts[p.Name] = bs
	return nil
}
func (p *StringListPart) Get(message *notification.PartedMessage) ([]string, error) {
	bytes := message.Parts[p.Name]
	data := []string{}
	err := Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type StringMapPart struct {
	Name string
}

func NewStringMapPart(name string) *StringMapPart {
	return &StringMapPart{
		Name: name,
	}
}

func (p *StringMapPart) Set(message *notification.PartedMessage, data map[string]string) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	message.Parts[p.Name] = bs
	return nil
}
func (p *StringMapPart) Get(message *notification.PartedMessage) (map[string]string, error) {
	bytes := message.Parts[p.Name]
	data := map[string]string{}
	err := Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type BinaryListPart struct {
	Name string
}

func NewBinaryListPart(name string) *StringListPart {
	return &StringListPart{
		Name: name,
	}
}

func (p *BinaryListPart) Set(message *notification.PartedMessage, data [][]byte) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	message.Parts[p.Name] = bs
	return nil
}
func (p *BinaryListPart) Get(message *notification.PartedMessage) ([][]byte, error) {
	bytes := message.Parts[p.Name]
	data := [][]byte{}
	err := Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type BinaryMapPart struct {
	Name string
}

func NewBinaryMapPart(name string) *BinaryMapPart {
	return &BinaryMapPart{
		Name: name,
	}
}

func (p *BinaryMapPart) Set(message *notification.PartedMessage, data map[string][]byte) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	message.Parts[p.Name] = bs
	return nil
}
func (p *BinaryMapPart) Get(message *notification.PartedMessage) (map[string][]byte, error) {
	bytes := message.Parts[p.Name]
	data := map[string][]byte{}
	err := Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type InterfacePart struct {
	Name string
}

func NewInterfacePart(name string) *InterfacePart {
	return &InterfacePart{
		Name: name,
	}
}

func (p *InterfacePart) Marshal(message *notification.PartedMessage, data interface{}) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	message.Parts[p.Name] = bs
	return nil
}
func (p *InterfacePart) Unmarshal(message *notification.PartedMessage, data interface{}) error {
	bytes := message.Parts[p.Name]
	return Unmarshal(bytes, &data)
}
