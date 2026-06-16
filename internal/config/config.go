package config

import (
	"errors"
	"fmt"
	"net"
	"os"
)

// Recognized deployment environments.
const (
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvProduction  = "production"
)

// Config holds runtime configuration loaded from the environment at boot.
type Config struct {
	// Env is the deployment environment: development, test, or production.
	Env string
	// HTTPAddr is the address the HTTP server listens on, for example ":8080".
	HTTPAddr string
	// LogLevel is one of debug, info, warn, error.
	LogLevel string
	// LogFormat is json or text.
	LogFormat string
	// DatabaseURL is the PostgreSQL connection string. It is required in
	// production and optional otherwise.
	DatabaseURL string
}

var (
	validEnvs       = map[string]bool{EnvDevelopment: true, EnvTest: true, EnvProduction: true}
	validLogLevels  = map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	validLogFormats = map[string]bool{"json": true, "text": true}
)

// Load reads and validates configuration from the OS environment. It returns
// an error if any value is invalid or a required value is missing, so callers
// can fail fast at boot.
func Load() (*Config, error) {
	return load(os.Getenv)
}

// load reads configuration using the provided getenv function. It is separated
// from Load so tests can supply a fake environment.
func load(getenv func(string) string) (*Config, error) {
	c := &Config{
		Env:         orDefault(getenv("FARAID_ENV"), EnvDevelopment),
		HTTPAddr:    orDefault(getenv("FARAID_HTTP_ADDR"), ":8080"),
		LogLevel:    orDefault(getenv("FARAID_LOG_LEVEL"), "info"),
		LogFormat:   orDefault(getenv("FARAID_LOG_FORMAT"), "json"),
		DatabaseURL: getenv("FARAID_DATABASE_URL"),
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// orDefault returns value when it is non-empty, otherwise def.
func orDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}

// validate checks every field and reports the first problem found.
func (c *Config) validate() error {
	if !validEnvs[c.Env] {
		return fmt.Errorf("config: invalid FARAID_ENV %q (want one of development, test, production)", c.Env)
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("config: invalid FARAID_LOG_LEVEL %q (want one of debug, info, warn, error)", c.LogLevel)
	}
	if !validLogFormats[c.LogFormat] {
		return fmt.Errorf("config: invalid FARAID_LOG_FORMAT %q (want one of json, text)", c.LogFormat)
	}
	if _, _, err := net.SplitHostPort(c.HTTPAddr); err != nil {
		return fmt.Errorf("config: invalid FARAID_HTTP_ADDR %q: %w", c.HTTPAddr, err)
	}
	if c.Env == EnvProduction && c.DatabaseURL == "" {
		return errors.New("config: FARAID_DATABASE_URL is required when FARAID_ENV=production")
	}
	return nil
}
