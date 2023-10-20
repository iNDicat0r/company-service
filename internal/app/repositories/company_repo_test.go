package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/iNDicat0r/company/common"
	"github.com/iNDicat0r/company/internal/app/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewSQLCompanyRepository(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		db     *gorm.DB
		expErr string
	}{
		"no database": {
			db:     nil,
			expErr: "db is nil",
		},
		"success": {
			db: &gorm.DB{},
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			repo, err := NewSQLCompanyRepository(tt.db)
			if tt.expErr != "" {
				assert.Nil(t, repo)
				assert.EqualError(t, err, tt.expErr)
			} else {
				assert.NotNil(t, repo)
			}
		})
	}
}

func TestSQLCompanyRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)

	repo, err := NewSQLCompanyRepository(db)
	if err != nil {
		t.Fatal("Failed to create repository: ", err)
	}

	company := models.Company{
		Name:            "Test Company",
		EmployeesAmount: 22,
		Registered:      true,
		Type:            common.Corporations,
		UserID:          uuid.MustParse("b6000e46-809f-4684-abd9-dc8f445b5ca9"),
	}
	db.Create(&company)

	id := company.ID
	foundCompany, err := repo.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if foundCompany.ID != id {
		t.Fatalf("Expected company with ID %s, but got %s", id, foundCompany.ID)
	}
}

func TestSQLCompanyRepository_Save(t *testing.T) {
	db := setupTestDB(t)

	repo, err := NewSQLCompanyRepository(db)
	if err != nil {
		t.Fatal("Failed to create repository: ", err)
	}

	company := models.Company{
		Name:            "Test Company",
		EmployeesAmount: 22,
		Registered:      true,
		Type:            common.Corporations,
		UserID:          uuid.MustParse("b6000e46-809f-4684-abd9-dc8f445b5ca9"),
	}

	_, err = repo.Save(context.Background(), company)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
}

func TestSQLCompanyRepository_Delete(t *testing.T) {
	db := setupTestDB(t)

	repo, err := NewSQLCompanyRepository(db)
	if err != nil {
		t.Fatal("Failed to create repository: ", err)
	}

	company := models.Company{
		Name:            "Test Company",
		EmployeesAmount: 22,
		Registered:      true,
		Type:            common.Corporations,
		UserID:          uuid.MustParse("b6000e46-809f-4684-abd9-dc8f445b5ca9"),
	}
	db.Create(&company)

	id := company.ID
	userID := uuid.New()
	err = repo.Delete(context.Background(), userID, id)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestSQLCompanyRepository_Update(t *testing.T) {
	db := setupTestDB(t)

	repo, err := NewSQLCompanyRepository(db)
	if err != nil {
		t.Fatal("Failed to create repository: ", err)
	}

	company := models.Company{
		Name:            "Test Company",
		EmployeesAmount: 22,
		Registered:      true,
		Type:            common.Corporations,
		UserID:          uuid.MustParse("b6000e46-809f-4684-abd9-dc8f445b5ca9"),
	}
	db.Create(&company)

	company.Name = "Updated Company"
	updatedCompany, err := repo.Update(context.Background(), company)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedCompany.Name != company.Name {
		t.Fatalf("Expected updated company name %s, but got %s", company.Name, updatedCompany.Name)
	}
}
