package userapplication

import (
	"context"
	"fmt"

	domainuser "starter/internal/domain/user"
	"starter/internal/eventbus"
)

// Service is the application service for the User bounded context.
//
// Responsibilities:
//   - Coordinate use-case execution (create domain entity → save → publish events).
//   - Never contain business logic – that lives in the domain entity.
//   - Translate between DTOs and domain objects.
type Service struct {
	repo domainuser.Repository
	bus  *eventbus.Bus
}

// NewService creates a new application service.
func NewService(repo domainuser.Repository, bus *eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// -----------------------------------------------------------------
// Use cases
// -----------------------------------------------------------------

// Create registers a new user and persists it inside a transaction.
// The UserCreated event is dispatched only after a successful commit.
func (s *Service) Create(ctx context.Context, req CreateUserRequest) (UserDTO, error) {
	u, err := domainuser.NewUser(req.Name)
	if err != nil {
		return UserDTO{}, fmt.Errorf("create user: %w", err)
	}

	// Save runs in a transaction; the auto-increment ID is set on u after return.
	if err := s.repo.Save(ctx, u); err != nil {
		return UserDTO{}, fmt.Errorf("persist user: %w", err)
	}

	// Publish events only after a successful commit.
	eventbus.PublishAll(s.bus, u.PullEvents())

	return ToUserDTO(u), nil
}

// GetAll returns all registered users as DTOs.
func (s *Service) GetAll(ctx context.Context) ([]UserDTO, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}

	dtos := make([]UserDTO, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, ToUserDTO(u))
	}
	return dtos, nil
}

// GetByID loads a single user by its ID.
func (s *Service) GetByID(ctx context.Context, id uint) (UserDTO, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return UserDTO{}, err
	}
	return ToUserDTO(u), nil
}
