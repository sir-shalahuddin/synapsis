package error_helper

import "errors"

var (
	ErrorDuplicateEmail = errors.New("email already registered")
	ErrorAuthentication = errors.New("email or password invalid")
	ErrorNotFound       = errors.New("entity not found")
	ErrorDuplicate      = errors.New("entity duplicate")
	ErrorUnAuthorized   = errors.New("forbidden")
)
