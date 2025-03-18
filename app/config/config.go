package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseConfig DatabaseConfig
	ServerConfig   ServerConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"db"`
	Port     string `env:"DB_PORT" envDefault:"27017"`
	Username string `env:"MONGO_INITDB_USERNAME" envDefault:""`
	Password string `env:"MONGO_INITDB_PASSWORD" envDefault:""`
	Database string `env:"MONGO_INITDB_DATABASE" envDefault:"library"`
}

type ServerConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return cfg, nil
}
