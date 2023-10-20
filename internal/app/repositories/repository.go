package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/iNDicat0r/company/internal/app/models"
)

// UserRepository defines the functionality of user repository.
type UserRepository interface {
	Save(ctx context.Context, user models.User) (uuid.UUID, error)
	FindByUserName(ctx context.Context, username string) (models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (models.User, error)
}

// CompanyRepository defines the functionality of company repository.
type CompanyRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (models.Company, error)
	Delete(ctx context.Context, userID, companyID uuid.UUID) error
	Save(ctx context.Context, company models.Company) (uuid.UUID, error)
	Update(ctx context.Context, company models.Company) (models.Company, error)
}
