package apperror

import "errors"

var (
	ErrAlreadyExists = errors.New("resource already exists")
	ErrUserNotFound  = errors.New("user not found")

	ErrUsernameInvalidLength = errors.New("username must be between 3 and 20 characters")
	ErrUsernameContainsSpace = errors.New("username cannot contain spaces")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrDatabase              = errors.New("database error")

	ErrInvalidRole    = errors.New("invalid role spcified")
	ErrSelfRoleChange = errors.New("cannot change your own role")
	ErrForbidden      = errors.New("permission denied")

	ErrInvalidTitle = errors.New("invalid title")
	ErrPostNotFound = errors.New("post not found")

	ErrInvalidInput = errors.New("invalid input parameters")
	ErrUnauthorized = errors.New("unauthorized")

	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token has been expired")
	ErrRefreshTokenRevoked = errors.New("refresh token has been revoked")
)
