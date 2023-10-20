package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password, err := HashPassword("123")
	assert.NoError(t, err)
	assert.NotEmpty(t, password)
}

func TestGenerateJWT(t *testing.T) {
	jwtKey := "123"
	userID := "12"
	tok, err := GenerateJWT(jwtKey, userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	derivedUserID, err := ParseJWT(jwtKey, tok)
	assert.NoError(t, err)
	assert.Equal(t, userID, derivedUserID)
}

func TestParseJWT(t *testing.T) {
	jwtKey := "123"
	userID := "12"
	tok, err := GenerateJWT(jwtKey, userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	invalidTok := "invalid" + tok

	derivedUserID, err := ParseJWT(jwtKey, invalidTok)
	assert.EqualError(t, err, "invalid token signature")
	assert.Equal(t, "", derivedUserID)
}
