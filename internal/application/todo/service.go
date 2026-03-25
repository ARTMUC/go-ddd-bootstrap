package todoapplication

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"starter/internal/domain/todo"
	"starter/internal/eventbus"
)

// Service is the application service for the TodoList bounded context.
//
// Responsibilities:
//   - Coordinate use-case execution (load aggregate → execute command → save → publish events).
//   - Never contain business logic – that lives in the domain aggregate.
//   - Translate between DTOs and domain objects.
type Service struct {
	repo todo.Repository
	bus  *eventbus.Bus
}

// NewService creates a new application service.
func NewService(repo todo.Repository, bus *eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// -----------------------------------------------------------------
// Use cases
// -----------------------------------------------------------------

// Create creates a new TodoList and persists it.
// Domain events recorded during creation are dispatched after the
// successful transaction commit.
func (s *Service) Create(ctx context.Context, req CreateTodoListRequest) (TodoListDTO, error) {
	tl, err := todo.NewTodoList(req.Title)
	if err != nil {
		return TodoListDTO{}, fmt.Errorf("create todo list: %w", err)
	}

	if err := s.repo.Save(ctx, tl); err != nil {
		return TodoListDTO{}, fmt.Errorf("persist todo list: %w", err)
	}

	// Dispatch events only after the transaction has committed successfully.
	eventbus.PublishAll(s.bus, tl.PullEvents())

	return ToTodoListDTO(tl), nil
}

// GetByID loads and returns the full aggregate as a DTO.
func (s *Service) GetByID(ctx context.Context, id string) (TodoListDTO, error) {
	uid, err := parseID(id)
	if err != nil {
		return TodoListDTO{}, err
	}

	tl, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return TodoListDTO{}, err
	}
	return ToTodoListDTO(tl), nil
}

// GetAll returns all TodoLists as DTOs.
func (s *Service) GetAll(ctx context.Context) ([]TodoListDTO, error) {
	lists, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]TodoListDTO, 0, len(lists))
	for _, tl := range lists {
		dtos = append(dtos, ToTodoListDTO(tl))
	}
	return dtos, nil
}

// AddLine loads the aggregate, adds a line, saves, and dispatches events.
// The save (list header update + new line insert) runs in one transaction.
func (s *Service) AddLine(ctx context.Context, todoListID string, req AddLineRequest) (TodoListDTO, error) {
	uid, err := parseID(todoListID)
	if err != nil {
		return TodoListDTO{}, err
	}

	tl, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return TodoListDTO{}, err
	}

	if _, err := tl.AddLine(req.Description); err != nil {
		return TodoListDTO{}, fmt.Errorf("add line: %w", err)
	}

	if err := s.repo.Save(ctx, tl); err != nil {
		return TodoListDTO{}, fmt.Errorf("persist todo list after add line: %w", err)
	}

	eventbus.PublishAll(s.bus, tl.PullEvents())

	return ToTodoListDTO(tl), nil
}

// CompleteLine marks a specific line as done.
// If all lines are done the aggregate transitions to Completed automatically.
// The save (updated line + possibly updated list status) runs in one transaction.
func (s *Service) CompleteLine(ctx context.Context, todoListID, lineID string) (TodoListDTO, error) {
	tlUID, err := parseID(todoListID)
	if err != nil {
		return TodoListDTO{}, err
	}
	lUID, err := parseID(lineID)
	if err != nil {
		return TodoListDTO{}, err
	}

	tl, err := s.repo.FindByID(ctx, tlUID)
	if err != nil {
		return TodoListDTO{}, err
	}

	if err := tl.CompleteLine(lUID); err != nil {
		return TodoListDTO{}, fmt.Errorf("complete line: %w", err)
	}

	if err := s.repo.Save(ctx, tl); err != nil {
		return TodoListDTO{}, fmt.Errorf("persist todo list after complete line: %w", err)
	}

	eventbus.PublishAll(s.bus, tl.PullEvents())

	return ToTodoListDTO(tl), nil
}

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

func parseID(raw string) (uuid.UUID, error) {
	uid, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid id %q: %w", raw, err)
	}
	return uid, nil
}
