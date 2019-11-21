package preset

type Config interface {
	Get(key string, v interface{}) error
}
type Option struct {
	Type   string
	Name   string
	Config map[string]interface{}
}
