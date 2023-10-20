package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the user stored in DB.
type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:UUID()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Username  string    `gorm:"index; unique"` // Unique and Index Username which will be used for auth.
	Password  string    `json:"-"`             // Hide password when json encoded.
	Companies []Company // Define a one-to-many relationship
}

func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
