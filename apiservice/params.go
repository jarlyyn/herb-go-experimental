package apiservice

import (
	"net/http"

	"github.com/herb-go/fetcher"
)

const ValueFieldID = "id"
const ValueFieldCredential = "credential"

type Values map[string][]byte

func (v *Values) SetValue(field string, value []byte) {
	(*v)[field] = value
}
func (v *Values) Value(field string) []byte {
	return (*v)[field]
}
func (v *Values) SetID(id string) {
	v.SetValue(ValueFieldID, []byte(id))
}

func (v *Values) ID() string {
	return string(v.Value(ValueFieldID))
}

func (v *Values) SetCredential(credential string) {
	v.SetValue(ValueFieldCredential, []byte(credential))
}

func (v *Values) Credential() string {
	return string(v.Value(ValueFieldCredential))
}

func NewValues() *Values {
	return &Values{}
}

type CommandBuilder interface {
	BuildCommand(*Values) (fetcher.Command, error)
}

type ValuesProvider interface {
	Values() (*Values, error)
}

type ValuesLoader interface {
	LoadValues(*http.Request) (*Values, error)
}

type ValuesMiddleware interface {
	ServeMiddlewareWithValues(v *Values, w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type MiddlewareBuilder interface {
	BuildMiddleware() ValuesMiddleware
}
