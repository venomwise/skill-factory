package config

import (
	"errors"
	"fmt"
)

// Config represents the application configuration
type Config struct {
	Profiles []Profile `mapstructure:"profiles"`
	BaseURL  string    `mapstructure:"base_url"`
	Timeout  int       `mapstructure:"timeout"` // timeout in seconds
}

// Profile represents an API key profile with optional base URL override
type Profile struct {
	ID      string `mapstructure:"id"`
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"` // optional override
}

// ErrNoAPIKey is returned when no valid API key is configured
var ErrNoAPIKey = errors.New("no API key configured")

// GetProfile returns a specific profile by ID
func (c *Config) GetProfile(id string) (*Profile, error) {
	for _, p := range c.Profiles {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("profile not found: %s", id)
}

// Validate checks if the configuration has at least one valid profile
func (c *Config) Validate() error {
	if len(c.Profiles) == 0 {
		return ErrNoAPIKey
	}
	return nil
}
