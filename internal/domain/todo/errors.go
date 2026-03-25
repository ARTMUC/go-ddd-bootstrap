package todo

import "errors"

var (
	ErrEmptyTitle       = errors.New("todo list title cannot be empty")
	ErrEmptyDescription = errors.New("todo list line description cannot be empty")
	ErrLineNotFound     = errors.New("todo list line not found")
	ErrAlreadyCompleted = errors.New("todo list is already completed")
	ErrNotFound         = errors.New("todo list not found")
)
