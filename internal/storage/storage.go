package storage

import (
	"errors"
)

var (
	ErrUsernameNotFound = errors.New("username is not found")
	ErrEmailNotFound    = errors.New("email is not found")
	ErrUserExists       = errors.New("user already exists")
)
