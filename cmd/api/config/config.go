package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Load returns Configuration struct
func Load() (*Configuration, error) {
	var cfg = new(Configuration)

	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file, %s", err)
	}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return cfg, nil
}

// Configuration holds data necessery for configuring application
type Configuration struct {
	Server *Server
	DB     *Database
	JWT    *JWT
}

// Database holds data necessery for database configuration
type Database struct {
	PSN          string `envconfig:"DSN" required:"true"`
	Log          bool   `envconfig:"DB_LOG" default:"false"`
	CreateSchema bool   `envconfig:"DB_CREATE_SCHEMA" default:"false"`
	Timeout      int    `envconfig:"DB_TIMEOUT" default:"5"`
}

// Server holds data necessery for server configuration
type Server struct {
	Port         string `envconfig:"PORT" default:":3000"`
	Debug        bool   `envconfig:"DEBUG" default:"false"`
	ReadTimeout  int    `envconfig:"READ_TIMEOUT" default:"5"`
	WriteTimeout int    `envconfig:"WRITE_TIMEOUT" default:"5"`
}

// JWT holds data necessery for JWT configuration
type JWT struct {
	Realm            string `envconfig:"JWT_REALM" required:"true"`
	Secret           string `envconfig:"JWT_SECRET" required:"true"`
	Duration         int    `envconfig:"JWT_DURATION" required:"true"`
	RefreshDuration  int
	MaxRefresh       int
	SigningAlgorithm string `envconfig:"JWT_ALGORITHM" default:"HS256"`
}
