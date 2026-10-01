package apperror

import "errors"

var (
	ErrAlreadyExists = errors.New("resource already exists")
	ErrUserNotFound  = errors.New("user not found")

	ErrUsernameInvalidLength = errors.New("username must be between 3 and 20 characters")
	ErrUsernameContainsSpace = errors.New("username cannot contain spaces")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrDatabase              = errors.New("database error")
)
