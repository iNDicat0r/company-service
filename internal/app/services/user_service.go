package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iNDicat0r/company/internal/app/models"
	"github.com/iNDicat0r/company/internal/app/repositories"
	"github.com/iNDicat0r/company/internal/app/utils"
	"golang.org/x/crypto/bcrypt"
)

// User defines the behaviours of the user functionalities in this service.
type User interface {
	Save(ctx context.Context, name, username, password string) (uuid.UUID, error)
	Authenticate(ctx context.Context, username, password string) (string, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.User, error)
}

// UserService represents a user service.
type UserService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

// NewUserService creates a new user service.
func NewUserService(userRepo repositories.UserRepository, jwtSecret string) (*UserService, error) {
	if userRepo == nil {
		return nil, errors.New("user repository is nil")
	}

	if jwtSecret == "" {
		return nil, errors.New("jwtSecret is empty")
	}

	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}, nil
}

// Save a user into db.
func (us *UserService) Save(ctx context.Context, name, username, password string) (uuid.UUID, error) {
	if name == "" {
		return uuid.UUID{}, errors.New("name is empty")
	}

	if username == "" {
		return uuid.UUID{}, errors.New("username is empty")
	}

	if password == "" {
		return uuid.UUID{}, errors.New("password is empty")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to save user due to password hashing: %w", err)
	}

	user := models.User{
		Name:     name,
		Password: hashedPassword,
		Username: username,
	}

	userID, err := us.userRepo.Save(ctx, user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to save user: %w", err)
	}

	return userID, nil
}

// Authenticate a user and returns a JWT token.
func (us *UserService) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := us.userRepo.FindByUserName(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong username/password combination")
	}

	token, err := utils.GenerateJWT(us.jwtSecret, user.ID.String())
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt token: %w", err)
	}

	return token, nil
}

// Introspect a user with details.
func (us *UserService) GetUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	user, err := us.userRepo.FindByID(ctx, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}
