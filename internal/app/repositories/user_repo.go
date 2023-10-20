package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iNDicat0r/company/internal/app/models"
	"gorm.io/gorm"
)

// SQLUserRepository implements the user storage, querying and db related logic.
type SQLUserRepository struct {
	db *gorm.DB
}

// NewSQLUserRepository creates a new mysql user repository.
func NewSQLUserRepository(db *gorm.DB) (*SQLUserRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &SQLUserRepository{
		db: db,
	}, nil
}

// Save a user into db.
func (u *SQLUserRepository) Save(ctx context.Context, user models.User) (uuid.UUID, error) {
	result := u.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return uuid.UUID{}, fmt.Errorf("failed to save user: %w", result.Error)
	}

	return user.ID, nil
}

// FindByUserName finds a user by username.
func (u *SQLUserRepository) FindByUserName(ctx context.Context, username string) (models.User, error) {
	var user models.User
	result := u.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return user, nil
}

// FindByID finds a user by id.
func (u *SQLUserRepository) FindByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	var user models.User
	result := u.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return user, nil
}
