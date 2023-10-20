package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iNDicat0r/company/common"
	"github.com/iNDicat0r/company/internal/app/models"
	"github.com/iNDicat0r/company/internal/app/repositories"
)

// CompanyGetCreateUpdateDeleter defines the functionality related to company service.
type CompanyGetCreateUpdateDeleter interface {
	Get(ctx context.Context, companyID uuid.UUID) (models.Company, error)
	Create(ctx context.Context, userID uuid.UUID, payload CreateUpdateCompanyPayload) (models.Company, error)
	Update(ctx context.Context, companyID uuid.UUID, payload CreateUpdateCompanyPayload) (models.Company, error)
	Delete(ctx context.Context, userID, companyID uuid.UUID) error
}

// CreateUpdateCompanyPayload represents the payload for creating or updating a company.
type CreateUpdateCompanyPayload struct {
	Name            string
	Description     string
	EmployeesAmount int
	Registered      bool
	Type            common.Type
}

// CompanyService represents the company service.
type CompanyService struct {
	companyRepo repositories.CompanyRepository
}

// NewCompanyService creates a new company service.
func NewCompanyService(companyRepo repositories.CompanyRepository) (*CompanyService, error) {
	if companyRepo == nil {
		return nil, errors.New("company repository is nil")
	}

	return &CompanyService{
		companyRepo: companyRepo,
	}, nil
}

// Get a company.
func (s *CompanyService) Get(ctx context.Context, companyID uuid.UUID) (models.Company, error) {
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return models.Company{}, fmt.Errorf("failed to get company: %w", err)
	}
	return comp, nil
}

// Create a company.
func (s *CompanyService) Create(ctx context.Context, userID uuid.UUID, payload CreateUpdateCompanyPayload) (models.Company, error) {
	if payload.Name == "" {
		return models.Company{}, errors.New("company name is empty")
	}

	if payload.EmployeesAmount == 0 {
		return models.Company{}, errors.New("company employees amount is empty")
	}

	if payload.Type == "" {
		return models.Company{}, errors.New("company type is empty")
	}

	b := models.Company{
		Name:            payload.Name,
		Description:     payload.Description,
		EmployeesAmount: payload.EmployeesAmount,
		Registered:      payload.Registered,
		Type:            payload.Type,
		UserID:          userID,
	}

	id, err := s.companyRepo.Save(ctx, b)
	if err != nil {
		return models.Company{}, fmt.Errorf("failed to save company: %w", err)
	}

	retrievedCompany, _ := s.companyRepo.FindByID(ctx, id)

	return retrievedCompany, nil
}

// Update a company.
func (s *CompanyService) Update(ctx context.Context, companyID uuid.UUID, payload CreateUpdateCompanyPayload) (models.Company, error) {
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return models.Company{}, fmt.Errorf("failed to find company: %w", err)
	}

	if payload.Name != "" {
		company.Name = payload.Name
	}

	if payload.EmployeesAmount != 0 {
		company.EmployeesAmount = payload.EmployeesAmount
	}

	if payload.Type == "" {
		company.EmployeesAmount = payload.EmployeesAmount
	}

	updated, err := s.companyRepo.Update(ctx, company)
	if err != nil {
		return models.Company{}, fmt.Errorf("failed to update company: %w", err)
	}

	return updated, nil
}

// Delete a company.
func (s *CompanyService) Delete(ctx context.Context, userID, companyID uuid.UUID) error {
	err := s.companyRepo.Delete(ctx, userID, companyID)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	return nil
}
