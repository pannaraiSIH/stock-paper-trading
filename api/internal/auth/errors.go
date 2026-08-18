package auth

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
