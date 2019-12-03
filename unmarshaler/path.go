package unmarshaler

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//InterfaceStep interface step struct
type InterfaceStep struct {
	value interface{}
}

//Type return step type
func (s *InterfaceStep) Type() interface{} {
	return TypeEmptyInterface
}

//String return step value as string.
func (s *InterfaceStep) String() string {
	return fmt.Sprint(s.value)
}

//Int return step value as int and any error if rasied.
func (s *InterfaceStep) Int() (int, bool) {
	return 0, false
}

//Interface return step value as interface.
func (s *InterfaceStep) Interface() interface{} {
	return s.value
}

//NewInterfaceStep create new interface step
func NewInterfaceStep(i interface{}) *InterfaceStep {
	s := InterfaceStep{
		value: i,
	}
	return &s
}

//StringStep type
type StringStep string

//Type return step type
func (s *StringStep) Type() interface{} {
	return TypeString
}

//String return step value as string.
func (s *StringStep) String() string {
	return string(*s)
}

//Int return step value as int and any error if rasied.
func (s *StringStep) Int() (int, bool) {
	return 0, false
}

//Interface return step value as interface.
func (s *StringStep) Interface() interface{} {
	return string(*s)
}

//NewStringStep create new string step.
func NewStringStep(str string) *StringStep {
	s := StringStep(str)
	return &s
}

//FieldStep field step struct
type FieldStep struct {
	*reflect.StructField
}

//Type return step type
func (s *FieldStep) Type() interface{} {
	return TypeStructField
}

//String return step value as string.
func (s *FieldStep) String() string {
	return s.Name
}

//Int return step value as int and any error if rasied.
func (s *FieldStep) Int() (int, bool) {
	return 0, true
}

//Interface return step value as interface.
func (s *FieldStep) Interface() interface{} {
	return *((*s).StructField)
}

//NewFieldStep create new field step
func NewFieldStep(f *reflect.StructField) *FieldStep {
	return &FieldStep{
		f,
	}
}

//ArrayStep arry step type
type ArrayStep int

//Type return step type
func (s *ArrayStep) Type() interface{} {
	return TypeArray
}

//String return step value as string.
func (s *ArrayStep) String() string {
	return strconv.Itoa(int(*s))
}

//Int return step value as int and any error if rasied.
func (s *ArrayStep) Int() (int, bool) {
	return int(*s), true
}

//Interface return step value as interface.
func (s *ArrayStep) Interface() interface{} {
	return int(*s)
}

//NewArrayStep create new array step
func NewArrayStep(i int) *ArrayStep {
	s := ArrayStep(i)
	return &s
}

//Step part iter step interface
type Step interface {
	Type() interface{}
	String() string
	Int() (int, bool)
	Interface() interface{}
}

//Steps steps list struct
type Steps struct {
	step   Step
	parent *Steps
}

//Join join step to path and return new path
func (s *Steps) Join(step Step) Path {
	steps := NewSteps()
	steps.parent = s
	steps.step = step
	return steps
}

//Pop pop last step and parent path
func (s *Steps) Pop() (Step, Path) {
	return s.step, s.parent
}

//NewSteps create new steps.
func NewSteps() *Steps {
	s := Steps{}
	return &s
}

//Path assembler path interface
type Path interface {
	//Join join step to path and return new path
	Join(Step) Path
	//Pop pop last step and parent path
	Pop() (Step, Path)
}

//ConvertPathToString convert assembler path to dot spilted string
func ConvertPathToString(p Path) string {
	var results = []string{}
	var step Step
	for p != nil {
		step, p = p.Pop()
		results = append(results, step.String())
	}
	return strings.Join(results, ".")
}
