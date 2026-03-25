package userpersistence

import (
	"time"

	"gorm.io/gorm"
)

// UserModel is the GORM persistence model for the User aggregate.
// It is intentionally decoupled from the domain entity.
type UserModel struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (UserModel) TableName() string { return "users" }
