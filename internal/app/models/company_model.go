package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/iNDicat0r/company/common"
	"gorm.io/gorm"
)

// Company represents the company stored in DB.
type Company struct {
	ID              uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Name            string         `gorm:"size:15;unique"`
	Description     string         `gorm:"size:3000"`
	EmployeesAmount int
	Registered      bool
	Type            common.Type `gorm:"type:enum('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship')"`
	UserID          uuid.UUID   `gorm:"type:uuid"`
}
