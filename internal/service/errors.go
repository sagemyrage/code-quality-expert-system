package service

import "errors"

var ErrUnauthenticated = errors.New("unauthenticated")

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
