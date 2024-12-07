package config

import (
	"time"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Server configuration
	ServerPort string `envconfig:"SERVER_PORT" default:"8080"`

	// Database configuration
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     string `envconfig:"DB_PORT" default:"3306"`
	DBUser     string `envconfig:"DB_USER" default:"root"`
	DBPassword string `envconfig:"DB_PASSWORD" default:""`
	DBName     string `envconfig:"DB_NAME" default:"ewallet_api"`

	// JWT configuration
	JWTSecret                string        `envconfig:"JWT_SECRET" default:"secretKeysJwt"`
	JWTExpirationHours      time.Duration `envconfig:"JWT_EXPIRATION_HOURS" default:"24h"`
	RefreshTokenSecret      string        `envconfig:"REFRESH_TOKEN_SECRET" default:"refreshSecretKeysJwt"`
	RefreshTokenExpirationDays time.Duration `envconfig:"REFRESH_TOKEN_EXPIRATION_DAYS" default:"168h"`
}

var cfg Config

// Load reads configuration from environment variables
func Load() (*Config, error) {
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Get returns the current configuration
func Get() *Config {
	return &cfg
}
