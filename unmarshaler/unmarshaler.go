package unmarshaler

type Unmarshaler interface {
	Unmarshal(data []byte, v interface{}) error
}

type DataConverter struct {
	marshaler   func(v interface{}) ([]byte, error)
	unmarshaler func(data []byte, v interface{}) error
}

func (c *DataConverter) Marshal(v interface{}) ([]byte, error) {
	return c.marshaler(v)
}

func (c *DataConverter) Unmarshal(data []byte, v interface{}) error {
	return c.unmarshaler(data, v)
}
