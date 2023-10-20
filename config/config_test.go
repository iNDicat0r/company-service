package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	configPath := "test_config.yaml"
	defer os.Remove(configPath)

	configData := `
global:
  jwt_signer_key: "test_signer_key"
server:
  port: 8080
  host: "localhost"
database:
  name: "testdb"
`

	// nolint:gofumpt
	err := os.WriteFile(configPath, []byte(configData), 0600)
	if err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	cfg, err := NewConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	expectedJWTSignerKey := "test_signer_key"
	assert.Equal(t, expectedJWTSignerKey, cfg.Global.JWTSignerKey)

	expectedServerPort := 8080
	assert.Equal(t, expectedServerPort, cfg.Server.Port)

	expectedServerHost := "localhost"
	assert.Equal(t, expectedServerHost, cfg.Server.Host)

	expectedDatabaseName := "testdb"
	assert.Equal(t, expectedDatabaseName, cfg.Database.Name)
}
