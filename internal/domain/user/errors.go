package user

import "errors"

var (
	ErrEmptyName = errors.New("user name cannot be empty")
	ErrNotFound  = errors.New("user not found")
)
