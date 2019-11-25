package assembler

import (
	"fmt"
	"reflect"
	"strconv"
)

type InterfaceStep struct {
	value interface{}
}

func (s InterfaceStep) Type() interface{} {
	return TypeInterface
}
func (s *InterfaceStep) String() string {
	return fmt.Sprint(s.value)
}
func (s *InterfaceStep) Int() (int, bool) {
	return 0, false
}
func (s *InterfaceStep) Interface() interface{} {
	return s.value
}

func NewInterfaceStep(i interface{}) *InterfaceStep {
	s := InterfaceStep{
		value: i,
	}
	return &s
}

type StringStep string

func (s StringStep) Type() interface{} {
	return TypeString
}
func (s *StringStep) String() string {
	return string(*s)
}
func (s *StringStep) Int() (int, bool) {
	return 0, false
}
func (s *StringStep) Interface() interface{} {
	return string(*s)
}

func NewStringStep(str string) *StringStep {
	s := StringStep(str)
	return &s
}

type FieldStep struct {
	*reflect.StructField
}

func (s *FieldStep) Type() interface{} {
	return TypeStructField
}
func (s *FieldStep) String() string {
	return s.Name
}
func (s *FieldStep) Int() (int, bool) {
	return 0, true
}
func (s *FieldStep) Interface() interface{} {
	return *s
}

func NewFieldStep(f *reflect.StructField) *FieldStep {
	return &FieldStep{
		f,
	}
}

type ArrayStep int

func (s *ArrayStep) Type() interface{} {
	return TypeArray
}
func (s *ArrayStep) String() string {
	return strconv.Itoa(int(*s))
}
func (s *ArrayStep) Int() (int, bool) {
	return int(*s), true
}
func (s *ArrayStep) Interface() interface{} {
	return int(*s)
}

func NewArrayStep(i int) *ArrayStep {
	s := ArrayStep(i)
	return &s
}

type Step interface {
	Type() interface{}
	String() string
	Int() (int, bool)
	Interface() interface{}
}
type Steps []Step

func (s *Steps) Join(steps ...Step) Path {
	p := s.Clone().(*Steps)
	*p = append(*p, steps...)
	return p
}
func (s *Steps) Clone() Path {
	newpath := make([]Step, len(*s))
	copy(newpath, *s)
	p := Steps(newpath)
	return &p
}
func (s *Steps) Unshift() (Step, Path) {
	if len(*s) == 0 {
		return nil, nil
	}
	steps := make([]Step, len(*s)-1)
	copy(steps, (*s)[1:])
	newpath := Steps(steps)
	return (*s)[0], &newpath
}
func NewSteps() *Steps {
	s := Steps([]Step{})
	return &s
}

type Path interface {
	Join(...Step) Path
	Unshift() (Step, Path)
	Clone() Path
}
