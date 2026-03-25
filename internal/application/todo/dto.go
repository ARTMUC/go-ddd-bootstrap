package todoapplication

import "starter/internal/domain/todo"

// -------------------------------------------------------
// Request DTOs  (HTTP → Application)
// -------------------------------------------------------

// CreateTodoListRequest carries data needed to create a new list.
type CreateTodoListRequest struct {
	Title string `json:"title" binding:"required"`
}

// AddLineRequest carries data for adding a new item to a list.
type AddLineRequest struct {
	Description string `json:"description" binding:"required"`
}

// -------------------------------------------------------
// Response DTOs  (Domain → HTTP)
// -------------------------------------------------------

// TodoListLineDTO is the read representation of a single TodoListLine.
type TodoListLineDTO struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// TodoListDTO is the read representation of a full TodoList aggregate.
type TodoListDTO struct {
	ID     string            `json:"id"`
	Title  string            `json:"title"`
	Status string            `json:"status"`
	Lines  []TodoListLineDTO `json:"lines"`
}

// -------------------------------------------------------
// Translator  (Domain → DTO)
// -------------------------------------------------------

// ToTodoListDTO converts the domain aggregate to its DTO representation.
func ToTodoListDTO(tl *todo.TodoList) TodoListDTO {
	lines := make([]TodoListLineDTO, 0, len(tl.Lines()))
	for _, l := range tl.Lines() {
		lines = append(lines, TodoListLineDTO{
			ID:          l.ID().String(),
			Description: l.Description(),
			Done:        l.Done(),
		})
	}
	return TodoListDTO{
		ID:     tl.ID().String(),
		Title:  tl.Title(),
		Status: tl.Status().String(),
		Lines:  lines,
	}
}
