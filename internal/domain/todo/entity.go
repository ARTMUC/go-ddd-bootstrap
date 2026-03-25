package todo

import "github.com/google/uuid"

// TodoListLine is an entity inside the TodoList aggregate.
// It must only be created/mutated through the TodoList aggregate root.
type TodoListLine struct {
	id          uuid.UUID
	todoListID  uuid.UUID
	description string
	done        bool
}

func newTodoListLine(todoListID uuid.UUID, description string) (*TodoListLine, error) {
	if description == "" {
		return nil, ErrEmptyDescription
	}
	return &TodoListLine{
		id:          uuid.New(),
		todoListID:  todoListID,
		description: description,
		done:        false,
	}, nil
}

// ReconstituteLine rebuilds a TodoListLine from persistence data (no invariant checks).
func ReconstituteLine(id, todoListID uuid.UUID, description string, done bool) *TodoListLine {
	return &TodoListLine{
		id:          id,
		todoListID:  todoListID,
		description: description,
		done:        done,
	}
}

func (l *TodoListLine) ID() uuid.UUID         { return l.id }
func (l *TodoListLine) TodoListID() uuid.UUID { return l.todoListID }
func (l *TodoListLine) Description() string   { return l.description }
func (l *TodoListLine) Done() bool            { return l.done }

func (l *TodoListLine) complete() {
	l.done = true
}
