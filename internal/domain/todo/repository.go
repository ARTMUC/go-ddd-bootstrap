package todo

import (
	"context"

	"github.com/google/uuid"
)

// Repository is the port that the domain defines.
// The infrastructure layer provides the concrete adapter.
type Repository interface {
	// Save persists the full aggregate (TodoList + all TodoListLines) in one transaction.
	// If the aggregate already exists it is updated; otherwise it is inserted.
	Save(ctx context.Context, todoList *TodoList) error

	// FindByID loads the full aggregate – including all its lines.
	FindByID(ctx context.Context, id uuid.UUID) (*TodoList, error)

	// FindAll returns all TodoList aggregates with their lines.
	FindAll(ctx context.Context) ([]*TodoList, error)
}
