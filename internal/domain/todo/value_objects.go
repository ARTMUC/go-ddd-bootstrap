package todo

import "errors"

// TodoStatus is a value object representing the state of a TodoList.
type TodoStatus string

const (
	StatusDraft     TodoStatus = "draft"
	StatusActive    TodoStatus = "active"
	StatusCompleted TodoStatus = "completed"
)

var validStatuses = map[TodoStatus]struct{}{
	StatusDraft:     {},
	StatusActive:    {},
	StatusCompleted: {},
}

// NewTodoStatus validates and returns a TodoStatus value object.
func NewTodoStatus(s string) (TodoStatus, error) {
	ts := TodoStatus(s)
	if _, ok := validStatuses[ts]; !ok {
		return "", errors.New("invalid todo status: " + s)
	}
	return ts, nil
}

func (s TodoStatus) String() string {
	return string(s)
}

func (s TodoStatus) IsCompleted() bool {
	return s == StatusCompleted
}
