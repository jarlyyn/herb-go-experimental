package auth

import (
	"net/http"
)

type Session interface {
	Set(r *http.Request, fieldname string, v interface{}) error
	Get(r *http.Request, fieldname string, v interface{}) error
	Del(r *http.Request, fieldname string) error
	IsNotFound(err error) bool
}
