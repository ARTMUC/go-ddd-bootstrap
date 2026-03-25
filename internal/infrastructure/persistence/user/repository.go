package userpersistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	domainuser "starter/internal/domain/user"
	"starter/internal/repository"
)

// UserRepository is the GORM-backed adapter for domainuser.Repository.
// It embeds BaseRepo[UserModel] to inherit standard CRUD helpers.
type UserRepository struct {
	repository.BaseRepo[UserModel]
	translator Translator
}

// NewUserRepository creates a new repository adapter.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepo:   repository.NewBaseRepo[UserModel](db),
		translator: Translator{},
	}
}

// -----------------------------------------------------------------
// domainuser.Repository implementation
// -----------------------------------------------------------------

// Save persists the User inside a transaction.
//   - ID == 0 → INSERT; the DB-assigned ID is written back to the domain entity.
//   - ID  > 0 → UPDATE (save all fields).
func (r *UserRepository) Save(ctx context.Context, u *domainuser.User) error {
	m := r.translator.ToModel(u)

	return r.Transaction(ctx, func(tx *gorm.DB) error {
		if u.ID() == 0 {
			if err := r.CreateTx(m, tx); err != nil {
				return fmt.Errorf("create user: %w", err)
			}
			u.SetID(m.ID)
		} else {
			if err := r.SaveTx(m, tx); err != nil {
				return fmt.Errorf("update user: %w", err)
			}
		}
		return nil
	})
}

// FindByID loads a User by its auto-increment primary key.
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*domainuser.User, error) {
	scope := r.DB().WithContext(ctx).Where("id = ?", id)

	m, err := r.FirstScoped(scope)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainuser.ErrNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return r.translator.ToDomain(m), nil
}

// FindAll returns every User.
func (r *UserRepository) FindAll(ctx context.Context) ([]*domainuser.User, error) {
	models, err := r.FindAllScoped(r.DB().WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("find all users: %w", err)
	}

	result := make([]*domainuser.User, 0, len(models))
	for _, m := range models {
		result = append(result, r.translator.ToDomain(m))
	}
	return result, nil
}
