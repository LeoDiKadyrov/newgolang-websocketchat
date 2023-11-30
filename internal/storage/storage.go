package storage

import (
	"errors"
)

var (
	ErrUsernameNotFound = errors.New("username is not found")
	ErrEmailNotFound    = errors.New("email is not found")
	ErrUserExists       = errors.New("user already exists")
)

/* Clean code thoughts & questions to myself
TODO:
[ ] - Shouldn't here be a method that init a storage by your choice (postgres, mysql, etc) and returns needed instance? Right now in main it's just: postgres.New(...)
*/
