package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config represents the configuration.
type Config struct {
	Global struct {
		JWTSignerKey string `yaml:"jwt_signer_key" envconfig:"GLOBAL_JWTSIGNERKEY"`
	} `yaml:"global"`
	Server struct {
		Port int    `yaml:"port" envconfig:"SERVER_PORT"`
		Host string `yaml:"host" envconfig:"SERVER_HOST"`
	} `yaml:"server"`
	Database struct {
		Name     string `yaml:"name" envconfig:"DATABASE_NAME"`
		Host     string `yaml:"host" envconfig:"DATABASE_HOST"`
		Port     int    `yaml:"port" envconfig:"DATABASE_PORT"`
		User     string `yaml:"user" envconfig:"DATABASE_USER"`
		Password string `yaml:"password" envconfig:"DATABASE_PASSWORD"`
	} `yaml:"database"`
}

// NewConfig returns a new configuration by parsing yml and env vars.
func NewConfig(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load env variables: %w", err)
	}

	return &cfg, nil
}
