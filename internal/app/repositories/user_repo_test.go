package repositories

import (
	"context"
	"testing"

	"github.com/iNDicat0r/company/internal/app/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestNewSQLLiteUserRepository(t *testing.T) {
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
			userRepo, err := NewSQLUserRepository(tt.db)
			if tt.expErr != "" {
				assert.Nil(t, userRepo)
				assert.EqualError(t, err, tt.expErr)
			} else {
				assert.NotNil(t, userRepo)
			}
		})
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	_ = db.AutoMigrate(&models.User{}, &models.Company{})
	return db
}

func TestSQLLiteUserRepository_Save(t *testing.T) {
	db := setupTestDB(t)

	repo, err := NewSQLUserRepository(db)
	assert.NoError(t, err)

	user := models.User{
		Name:     "John Doe",
		Username: "johndoe",
		Password: "password",
	}

	ctx := context.Background()
	userID, err := repo.Save(ctx, user)

	assert.NoError(t, err)
	assert.NotZero(t, userID)

	retrievedUser, err := repo.FindByID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Password, retrievedUser.Password)
}

func TestSQLLiteUserRepository_FindByUserName(t *testing.T) {
	db := setupTestDB(t)

	repo, err := NewSQLUserRepository(db)
	assert.NoError(t, err)

	user := models.User{
		Name:     "Alice Smith",
		Username: "alicesmith",
		Password: "password123",
	}

	ctx := context.Background()
	_, err = repo.Save(ctx, user)
	assert.NoError(t, err)

	retrievedUser, err := repo.FindByUserName(ctx, user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Password, retrievedUser.Password)
}

func TestSQLLiteUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)

	repo, err := NewSQLUserRepository(db)
	assert.NoError(t, err)

	user := models.User{
		Name:     "Bob Johnson",
		Username: "bobjohnson",
		Password: "pass123",
	}

	ctx := context.Background()
	userID, err := repo.Save(ctx, user)
	assert.NoError(t, err)

	retrievedUser, err := repo.FindByID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Password, retrievedUser.Password)
}
