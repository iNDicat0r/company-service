package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/iNDicat0r/company/internal/app/models"
	"github.com/iNDicat0r/company/internal/app/repositories"
	"github.com/stretchr/testify/assert"
)

func TestNewCompanyService(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		companyRepo repositories.CompanyRepository
		expErr      string
	}{
		"company repo is nil": {
			expErr: "company repository is nil",
		},
		"success": {
			companyRepo: &mockCompanyRepository{},
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s, err := NewCompanyService(tt.companyRepo)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
				assert.Nil(t, s)
			} else {
				assert.NotNil(t, s)
			}
		})
	}
}

func TestCompanyService_GetCompany(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		companyRepo repositories.CompanyRepository
		companyID   uuid.UUID
		expErr      string
	}{
		"company repo error": {
			companyRepo: &mockCompanyRepository{
				err: errors.New("company repo error"),
			},

			expErr: "failed to find company 0: company repo error",
		},
		"success": {
			companyRepo: &mockCompanyRepository{
				singleCompany: models.Company{
					Name: "hello world",
				},
			},
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s, err := NewCompanyService(tt.companyRepo)
			assert.NoError(t, err)
			comp, err := s.Get(context.TODO(), tt.companyID)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
			} else {
				assert.Equal(t, comp.ID, tt.companyID)
			}
		})
	}
}

// mockCompanyRepository for testing
type mockCompanyRepository struct {
	err           error
	singleCompany models.Company
	id            uuid.UUID
}

func (m *mockCompanyRepository) FindByID(_ context.Context, _ uuid.UUID) (models.Company, error) {
	return m.singleCompany, m.err
}

func (m *mockCompanyRepository) Delete(_ context.Context, _, _ uuid.UUID) error {
	return m.err
}

func (m *mockCompanyRepository) Save(_ context.Context, _ models.Company) (uuid.UUID, error) {
	return m.id, m.err
}

func (m *mockCompanyRepository) Update(_ context.Context, _ models.Company) (models.Company, error) {
	return m.singleCompany, m.err
}
