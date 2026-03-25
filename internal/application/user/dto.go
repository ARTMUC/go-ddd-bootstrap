package userapplication

import domainuser "starter/internal/domain/user"

// -------------------------------------------------------
// Request DTOs  (HTTP → Application)
// -------------------------------------------------------

// CreateUserRequest carries data needed to register a new user.
type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// -------------------------------------------------------
// Response DTOs  (Domain → HTTP)
// -------------------------------------------------------

// UserDTO is the read representation of a User entity.
type UserDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// -------------------------------------------------------
// Translator  (Domain → DTO)
// -------------------------------------------------------

// ToUserDTO converts the domain User entity to its DTO representation.
func ToUserDTO(u *domainuser.User) UserDTO {
	return UserDTO{
		ID:   u.ID(),
		Name: u.Name(),
	}
}
