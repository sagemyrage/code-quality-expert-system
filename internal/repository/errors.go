package repository

import "errors"

var (
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
)
