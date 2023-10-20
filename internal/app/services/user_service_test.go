package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/iNDicat0r/company/internal/app/models"
	"github.com/iNDicat0r/company/internal/app/repositories"
	"github.com/iNDicat0r/company/internal/app/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewUserService(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		userRepo  repositories.UserRepository
		jwtSecret string
		expErr    string
	}{
		"no user repository": {
			expErr: "user repository is nil",
		},
		"empty jwt key signer": {
			userRepo: &mockUserRepository{},
			expErr:   "jwtSecret is empty",
		},
		"success": {
			userRepo:  &mockUserRepository{},
			jwtSecret: "123",
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s, err := NewUserService(tt.userRepo, tt.jwtSecret)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
				assert.Nil(t, s)
			} else {
				assert.NotNil(t, s)
			}
		})
	}
}

func TestUserService_Save(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		userRepo  repositories.UserRepository
		jwtSecret string
		name      string
		username  string
		password  string
		createdID uint64
		expErr    string
	}{
		"empty name": {
			userRepo:  &mockUserRepository{},
			jwtSecret: "123",
			expErr:    "name is empty",
		},
		"empty username": {
			userRepo:  &mockUserRepository{},
			jwtSecret: "123",
			name:      "user1",
			expErr:    "username is empty",
		},
		"empty password": {
			userRepo:  &mockUserRepository{},
			jwtSecret: "123",
			name:      "user1",
			username:  "123433",
			expErr:    "password is empty",
		},
		"empty authorPseudonym": {
			userRepo:  &mockUserRepository{},
			jwtSecret: "123",
			name:      "user1",
			username:  "123433",
			password:  "3cidjhjh",
			expErr:    "authorPseudonym is empty",
		},
		"userRepo error": {
			userRepo: &mockUserRepository{
				err: errors.New("issue with db"),
			},
			jwtSecret: "123",
			name:      "user1",
			username:  "123433",
			password:  "3cidjhjh",
			expErr:    "failed to save user: issue with db",
		},
		"success": {
			userRepo: &mockUserRepository{
				id: uuid.UUID{},
			},
			jwtSecret: "123",
			name:      "user1",
			username:  "123433",
			password:  "3cidjhjh",
			createdID: 33,
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s, err := NewUserService(tt.userRepo, tt.jwtSecret)
			assert.NoError(t, err)
			id, err := s.Save(context.TODO(), tt.name, tt.username, tt.password)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
				assert.Equal(t, uint64(0), id)
			} else {
				assert.Equal(t, tt.createdID, id)
			}
		})
	}
}

func TestUserService_Authenticate(t *testing.T) {
	t.Parallel()
	hashedPass, err := utils.HashPassword("12345")
	assert.NoError(t, err)
	cases := map[string]struct {
		userRepo   repositories.UserRepository
		jwtSecret  string
		username   string
		password   string
		createdJWT string
		expErr     string
	}{
		"userRepo error": {
			userRepo: &mockUserRepository{
				err: errors.New("issue with db"),
			},
			jwtSecret: "123",
			expErr:    "failed to authenticate: issue with db",
		},
		"compare password error": {
			userRepo:  &mockUserRepository{},
			jwtSecret: "123",
			expErr:    "wrong username/password combination",
		},
		"success": {
			userRepo: &mockUserRepository{
				user: models.User{
					Password: hashedPass,
				},
			},
			jwtSecret: "secretkeyhere",
			username:  "john",
			password:  "12345",
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s, err := NewUserService(tt.userRepo, tt.jwtSecret)
			assert.NoError(t, err)
			jwt, err := s.Authenticate(context.TODO(), tt.username, tt.password)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
				assert.Equal(t, "", jwt)
			} else {
				assert.NotEqual(t, "", jwt)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		userRepo  repositories.UserRepository
		jwtSecret string
		userID    uuid.UUID
		expErr    string
	}{
		"userRepo error": {
			userRepo: &mockUserRepository{
				err: errors.New("issue with db"),
			},
			jwtSecret: "123",
			expErr:    "user not found: issue with db",
		},
		"success": {
			userRepo: &mockUserRepository{
				user: models.User{ID: uuid.MustParse("862dedcb-68c5-49f7-a94a-b7190499f16b")},
			},
			userID:    uuid.MustParse("862dedcb-68c5-49f7-a94a-b7190499f16b"),
			jwtSecret: "123",
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s, err := NewUserService(tt.userRepo, tt.jwtSecret)
			assert.NoError(t, err)
			user, err := s.GetUser(context.TODO(), tt.userID)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
			} else {
				assert.Equal(t, tt.userID, user.ID)
			}
		})
	}
}

type mockUserRepository struct {
	id   uuid.UUID
	err  error
	user models.User
}

func (m *mockUserRepository) Save(_ context.Context, _ models.User) (uuid.UUID, error) {
	return m.id, m.err
}

func (m *mockUserRepository) FindByUserName(_ context.Context, _ string) (models.User, error) {
	return m.user, m.err
}

func (m *mockUserRepository) FindByID(_ context.Context, _ uuid.UUID) (models.User, error) {
	return m.user, m.err
}
