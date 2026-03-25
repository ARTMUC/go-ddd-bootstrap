package todopersistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"starter/internal/domain/todo"
	"starter/internal/repository"
)

// TodoListRepository is the GORM-backed adapter for todo.Repository.
// It embeds BaseRepo[TodoListModel] to inherit standard CRUD helpers and
// adds aggregate-specific methods on top.
type TodoListRepository struct {
	repository.BaseRepo[TodoListModel]
	translator Translator
}

// NewTodoListRepository creates a new repository adapter.
func NewTodoListRepository(db *gorm.DB) *TodoListRepository {
	return &TodoListRepository{
		BaseRepo:   repository.NewBaseRepo[TodoListModel](db),
		translator: Translator{},
	}
}

// -----------------------------------------------------------------
// todo.Repository implementation
// -----------------------------------------------------------------

// Save upserts the TodoList aggregate (header + lines) in a single transaction.
//
// Strategy inside the transaction:
//  1. Upsert todo_lists row via BaseRepo.UpsertTx.
//  2. Delete any lines that are no longer on the aggregate.
//  3. Upsert all current lines via a per-line BaseRepo[TodoListLineModel].UpsertTx.
func (r *TodoListRepository) Save(ctx context.Context, tl *todo.TodoList) error {
	listModel, lineModels := r.translator.ToModel(tl)

	return r.Transaction(ctx, func(tx *gorm.DB) error {
		// 1. Upsert the aggregate root row.
		if err := r.UpsertTx(
			listModel,
			[]string{"id"},
			[]string{"title", "status", "updated_at"},
			tx,
		); err != nil {
			return fmt.Errorf("save todo list header: %w", err)
		}

		// 2. Delete lines that were removed from the aggregate.
		currentIDs := make([]string, 0, len(lineModels))
		for _, l := range lineModels {
			currentIDs = append(currentIDs, l.ID)
		}
		// Sentinel prevents "IN ()" which is a SQL error on some drivers.
		if err := tx.
			Where("todo_list_id = ? AND id NOT IN ?", listModel.ID, append(currentIDs, "")).
			Delete(&TodoListLineModel{}).Error; err != nil {
			return fmt.Errorf("delete removed lines: %w", err)
		}

		// 3. Upsert all current lines.
		//    TodoListLineModel is a different type from T so we create a scoped
		//    BaseRepo for it and delegate to UpsertTx.
		if len(lineModels) > 0 {
			lineRepo := repository.NewBaseRepo[TodoListLineModel](tx)
			for i := range lineModels {
				if err := lineRepo.UpsertTx(
					&lineModels[i],
					[]string{"id"},
					[]string{"description", "done", "updated_at"},
					tx,
				); err != nil {
					return fmt.Errorf("upsert line %s: %w", lineModels[i].ID, err)
				}
			}
		}

		return nil
	})
}

// FindByID loads the full TodoList aggregate including all its lines.
func (r *TodoListRepository) FindByID(ctx context.Context, id uuid.UUID) (*todo.TodoList, error) {
	scope := r.DB().WithContext(ctx).Preload("Lines").Where("id = ?", id.String())

	m, err := r.FirstScoped(scope)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, todo.ErrNotFound
		}
		return nil, fmt.Errorf("find todo list by id: %w", err)
	}

	return r.translator.ToDomain(m)
}

// FindAll loads every TodoList with its lines.
func (r *TodoListRepository) FindAll(ctx context.Context) ([]*todo.TodoList, error) {
	models, err := r.FindAllScoped(r.DB().WithContext(ctx).Preload("Lines"))
	if err != nil {
		return nil, fmt.Errorf("find all todo lists: %w", err)
	}

	result := make([]*todo.TodoList, 0, len(models))
	for i := range models {
		tl, err := r.translator.ToDomain(models[i])
		if err != nil {
			return nil, err
		}
		result = append(result, tl)
	}
	return result, nil
}
