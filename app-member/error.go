package member

import "errors"

var ErrRegisteredDataNotMap = errors.New("registered data is not a map")

var ErrAccountRegisterExists = errors.New("account registered exists")

var ErrUserBanned = errors.New("user banned")
