package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iNDicat0r/company/internal/app/models"
	"gorm.io/gorm"
)

// SQLCompanyRepository implements the company storage, querying and db related logic.
type SQLCompanyRepository struct {
	db *gorm.DB
}

// NewSQLCompanyRepository creates a new mysql company repository.
func NewSQLCompanyRepository(db *gorm.DB) (*SQLCompanyRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &SQLCompanyRepository{
		db: db,
	}, nil
}

// FindByID returns a company by id.
func (br *SQLCompanyRepository) FindByID(ctx context.Context, id uuid.UUID) (models.Company, error) {
	var comp models.Company
	result := br.db.WithContext(ctx).Where("id = ?", id).First(&comp)
	if result.Error != nil {
		return models.Company{}, fmt.Errorf("failed to find company: %w", result.Error)
	}

	return comp, nil
}

// Save a company into db.
func (br *SQLCompanyRepository) Save(ctx context.Context, company models.Company) (uuid.UUID, error) {
	result := br.db.WithContext(ctx).Create(&company)
	if result.Error != nil {
		return uuid.UUID{}, fmt.Errorf("failed to save company: %w", result.Error)
	}

	return company.ID, nil
}

// Delete a company from db.
func (br *SQLCompanyRepository) Delete(ctx context.Context, userID, companyID uuid.UUID) error {
	var comp models.Company
	result := br.db.WithContext(ctx).Where("user_id = ?", userID).Where("id = ?", companyID).First(&comp)
	if result.Error != nil {
		return fmt.Errorf("failed to find company: %w", result.Error)
	}

	if err := br.db.Delete(&comp).Error; err != nil {
		return fmt.Errorf("failed to delete a company: %w", err)
	}

	return nil
}

// Delete a company from db.
func (br *SQLCompanyRepository) Update(ctx context.Context, company models.Company) (models.Company, error) {
	result := br.db.WithContext(ctx).Save(&company)
	if result.Error != nil {
		return models.Company{}, fmt.Errorf("failed to update company: %w", result.Error)
	}
	return company, nil
}
