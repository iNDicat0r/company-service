package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password using bcrypt with the default cost.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %v", err)
	}

	return string(hashedPassword), nil
}

// GenerateJWT generates a jwt token given the username.
func GenerateJWT(jwtKey string, userID string) (string, error) {
	expirationTime := time.Now().Add(6 * time.Hour)
	claims := &claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ParseJWT parses and validates a jwt and returns the username.
func ParseJWT(jwtKey, token string) (string, error) {
	claims := &claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return "", errors.New("invalid token signature")
	}

	if !tkn.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}
