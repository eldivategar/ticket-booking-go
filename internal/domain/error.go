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

	// Specific errors
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidID             = errors.New("invalid id")

	// Auth errors
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Event errors
	ErrInvalidStock   = errors.New("invalid stock")
	ErrInvalidPrice   = errors.New("invalid price")
	ErrInvalidDate    = errors.New("invalid date")
	ErrEventNotFound  = errors.New("event not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)
