package unmarshaler

import (
	"fmt"
)

//Unmarshaler unmarshaler interface
type Unmarshaler interface {
	//Unmarshal data to giver interface.
	//Return any error if raised.
	Unmarshal(data []byte, v interface{}) error
}

//unmarshalers all registered unmarshaler
var unmarshalers = map[string]Unmarshaler{}

//RegisterUnmarshaler register unmarshaler with given name.
func RegisterUnmarshaler(name string, u Unmarshaler) {
	unmarshalers[name] = u
}

//UnregisterAllUnmarshaler unreister all unmarshaler.
func UnregisterAllUnmarshaler() {
	unmarshalers = map[string]Unmarshaler{}
}

//Unmarshal unmarshal byte slice to data by given munarshaler.
//Return any error if raised
func Unmarshal(name string, data []byte, v interface{}) error {
	u := unmarshalers[name]
	if u == nil {
		return fmt.Errorf("unmarshaler : %w (%s)", ErrUnmarshalerNotRegistered, name)
	}
	return u.Unmarshal(data, v)
}
