package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInternal           = errors.New("internal error")
	ErrUserNotLoggedIn    = errors.New("user not logged in")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrPrivateAccount     = errors.New("private account")
)
