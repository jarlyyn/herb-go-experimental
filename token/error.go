package token

import "errors"

var ErrEmptyID = errors.New("empty id")
var ErrTokenNotFound = errors.New("token not found")
var ErrTokenExists = errors.New("token exists")
