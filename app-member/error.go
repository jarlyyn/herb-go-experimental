package member

import "errors"

var ErrRegisteredDataNotMap = errors.New("registered data is not a map")

var ErrAccountRegisterExists = errors.New("account registered exists")

var ErrUserBanned = errors.New("user banned")

var ErrUserNotFound = errors.New("user not found")

var ErrFeatureNotSupported = errors.New("feature not supported")

var ErrAccountKeywordNotRegistered = errors.New("account keyword not regietered")
