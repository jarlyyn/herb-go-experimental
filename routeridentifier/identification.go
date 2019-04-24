package routeridentifier

import (
	"context"
	"net/http"
	"strings"

	"github.com/herb-go/herb/middleware/router"
)

const ContextNameRouterIdentification = router.ContextName("routerIdentification")

type Identification struct {
	ID   string
	Tags []string
}

func (i *Identification) String() string {
	tags := strings.Join(i.Tags, ",")
	if tags != "" {
		return i.ID + "|" + tags
	}
	return i.ID
}
func (i *Identification) AddTag(tag string) {
	if tag == "" {
		return
	}
	for k := range i.Tags {
		if i.Tags[k] == tag {
			return
		}
	}
	i.Tags = append(i.Tags, tag)
}

func (i *Identification) HasTag(tag string) bool {
	for k := range i.Tags {
		if i.Tags[k] == tag {
			return true
		}
	}
	return false
}
func NewIdentification() *Identification {
	return &Identification{
		Tags: []string{},
	}
}

func SetIdentificationToRequest(r *http.Request, i *Identification) {
	ctx := context.WithValue(r.Context(), ContextNameRouterIdentification, i)
	*r = *r.WithContext(ctx)

}

func GetIdentificationFromRequest(r *http.Request) *Identification {
	v := r.Context().Value(ContextNameRouterIdentification)
	if v == nil {
		return nil
	}
	return v.(*Identification)
}
