package unmarshaler

import (
	"encoding/json"
	"errors"
	"fmt"
)

var a = json.Unmarshal

var ErrNotSetable = errors.New("value cannot set")

var ErrNotAssignable = errors.New("value is not assignable")

type AssemblerError struct {
	a   *Assembler
	err error
}

func (e *AssemblerError) Unwrap() error {
	return e.err
}

func NewAssemblerError(a *Assembler, err error) *AssemblerError {
	return &AssemblerError{
		a:   a,
		err: err,
	}
}

func (e *AssemblerError) Error() string {
	return fmt.Sprintf("unmarshaler: error: %s (%s)", e.err.Error(), ConvertPathToString(e.a.Path()))
}
