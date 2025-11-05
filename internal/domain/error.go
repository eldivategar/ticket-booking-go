package domain

import "errors"

var (
	// Common auth errors
	ErrUnauthorized = errors.New("unauthorized")

	// General errors
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource already exists")
	ErrInvalidInput    = errors.New("invalid input")
	ErrInternal        = errors.New("internal server error")
	ErrBadRequest      = errors.New("bad request")
	ErrTooManyRequests = errors.New("too many requests")
)
