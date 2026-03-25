package user

import "context"

// Repository is the port (interface) for User persistence.
// The infrastructure layer provides the concrete adapter.
type Repository interface {
	// Save persists a User. For a new User (ID == 0) it inserts and sets the ID.
	// For an existing User it updates the record.
	Save(ctx context.Context, u *User) error

	// FindByID loads a User by its auto-increment ID.
	FindByID(ctx context.Context, id uint) (*User, error)

	// FindAll returns every User.
	FindAll(ctx context.Context) ([]*User, error)
}
