package components

import "reflect"
import "errors"
type ComponentFactory  func([]byte) (interface{},error)

var ErrMustRegisterPtrOrInterface=errors.New("only pointer or interface can be registered as component")

type Component struct{
	Module string
	Name string
	target reflect.Value
	origin reflect.Value
}

func NewComponent (value interface{})(*Component,error){
	tp:=reflect.TypeOf(value)
	kind:= tp.Kind()
	if kind!=reflect.Ptr||kind!=reflect.Interface{
		return nil,ErrMustRegisterPtrOrInterface
	}
	target:=reflect.ValueOf(value)
	origin:=target.Elem()
	return &Component{
		target:target,
		origin:origin,
	},nil
}	
func (c *Component) Overwrite(v interface{}){
	
}
func (c *Component)Rollback(){

}
type Components []*Component

func(c *Components) MustRegister(module string,name string,value interface{})*Component{
	component,err:=NewComponent(value)
	if err!=nil{
		panic(err)
	}
	component.Module=module
	component.Name=name
	(*c)=append(*c,component)
	return component
}