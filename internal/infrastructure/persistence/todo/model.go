package todopersistence

import (
	"time"

	"gorm.io/gorm"
)

// TodoListModel is the GORM persistence model for the TodoList aggregate root.
// It is intentionally decoupled from the domain aggregate.
type TodoListModel struct {
	ID        string `gorm:"primaryKey"`
	Title     string
	Status    string
	Lines     []TodoListLineModel `gorm:"foreignKey:TodoListID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (TodoListModel) TableName() string { return "todo_lists" }

// TodoListLineModel is the GORM persistence model for a single TodoListLine entity.
type TodoListLineModel struct {
	ID          string `gorm:"primaryKey"`
	TodoListID  string
	Description string
	Done        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (TodoListLineModel) TableName() string { return "todo_list_lines" }
