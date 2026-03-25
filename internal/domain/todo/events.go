package todo

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent is the base interface for all domain events.
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

// -------------------------------------------------------
// TodoListCreated
// -------------------------------------------------------

// TodoListCreated is raised when a new TodoList aggregate is created.
type TodoListCreated struct {
	TodoListID uuid.UUID
	Title      string
	OccurredOn time.Time
}

func (e TodoListCreated) EventName() string     { return "todo_list.created" }
func (e TodoListCreated) OccurredAt() time.Time { return e.OccurredOn }

// -------------------------------------------------------
// TodoListLineAdded
// -------------------------------------------------------

// TodoListLineAdded is raised when a new line is added to the list.
type TodoListLineAdded struct {
	TodoListID uuid.UUID
	LineID     uuid.UUID
	Desc       string
	OccurredOn time.Time
}

func (e TodoListLineAdded) EventName() string     { return "todo_list.line_added" }
func (e TodoListLineAdded) OccurredAt() time.Time { return e.OccurredOn }

// -------------------------------------------------------
// TodoListLineCompleted
// -------------------------------------------------------

// TodoListLineCompleted is raised when a specific line is marked as done.
type TodoListLineCompleted struct {
	TodoListID uuid.UUID
	LineID     uuid.UUID
	OccurredOn time.Time
}

func (e TodoListLineCompleted) EventName() string     { return "todo_list.line_completed" }
func (e TodoListLineCompleted) OccurredAt() time.Time { return e.OccurredOn }

// -------------------------------------------------------
// TodoListCompleted
// -------------------------------------------------------

// TodoListCompleted is raised when every line in the list has been done.
type TodoListCompleted struct {
	TodoListID uuid.UUID
	OccurredOn time.Time
}

func (e TodoListCompleted) EventName() string     { return "todo_list.completed" }
func (e TodoListCompleted) OccurredAt() time.Time { return e.OccurredOn }
