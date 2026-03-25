package todopersistence

import (
	"fmt"

	"starter/internal/domain/todo"

	"github.com/google/uuid"
)

// Translator converts between the domain aggregate and the GORM persistence model.
// This is the Anti-Corruption Layer between domain and infrastructure.
type Translator struct{}

// ToModel converts the TodoList domain aggregate (and all its lines) to GORM models.
// The caller uses the returned models inside a transaction.
func (t Translator) ToModel(tl *todo.TodoList) (*TodoListModel, []TodoListLineModel) {
	listModel := &TodoListModel{
		ID:     tl.ID().String(),
		Title:  tl.Title(),
		Status: tl.Status().String(),
	}

	lineModels := make([]TodoListLineModel, 0, len(tl.Lines()))
	for _, l := range tl.Lines() {
		lineModels = append(lineModels, TodoListLineModel{
			ID:          l.ID().String(),
			TodoListID:  tl.ID().String(),
			Description: l.Description(),
			Done:        l.Done(),
		})
	}
	return listModel, lineModels
}

// ToDomain converts GORM persistence models back to the domain aggregate.
// Lines should already be preloaded on the listModel.
func (t Translator) ToDomain(m *TodoListModel) (*todo.TodoList, error) {
	listID, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("translator: invalid todo list id %q: %w", m.ID, err)
	}

	status, err := todo.NewTodoStatus(m.Status)
	if err != nil {
		return nil, fmt.Errorf("translator: %w", err)
	}

	lines := make([]*todo.TodoListLine, 0, len(m.Lines))
	for _, lm := range m.Lines {
		lineID, err := uuid.Parse(lm.ID)
		if err != nil {
			return nil, fmt.Errorf("translator: invalid line id %q: %w", lm.ID, err)
		}
		lines = append(lines, todo.ReconstituteLine(lineID, listID, lm.Description, lm.Done))
	}

	return todo.Reconstitute(listID, m.Title, status, lines), nil
}
