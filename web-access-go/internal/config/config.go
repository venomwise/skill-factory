package config

import (
	"errors"
)

// Config represents the unified web-access configuration
type Config struct {
	Exa  ProviderConfig `toml:"exa"`
	Grok ProviderConfig `toml:"grok"`
}

// ProviderConfig represents provider-specific configuration
type ProviderConfig struct {
	BaseURL  string    `toml:"base_url"`
	Timeout  int       `toml:"timeout"` // timeout in seconds
	Model    string    `toml:"model"`   // for Grok only
	Profiles []Profile `toml:"profiles"`

	// Grok-specific fields
	ExtraBody    map[string]interface{} `toml:"extra_body"`
	ExtraHeaders map[string]string      `toml:"extra_headers"`
}

// Profile represents an API key profile with optional overrides
type Profile struct {
	ID      string `toml:"id"`
	APIKey  string `toml:"api_key"`
	BaseURL string `toml:"base_url"` // optional override
	Model   string `toml:"model"`    // optional override for Grok
}

// ResolvedConfig represents the final resolved configuration after applying precedence
type ResolvedConfig struct {
	Exa  ResolvedProviderConfig
	Grok ResolvedProviderConfig
}

// ResolvedProviderConfig represents resolved provider-specific configuration
type ResolvedProviderConfig struct {
	BaseURL      string
	Timeout      int
	Model        string // for Grok only
	Profiles     []ResolvedProfile
	ExtraBody    map[string]interface{} // for Grok only
	ExtraHeaders map[string]string      // for Grok only
}

// ResolvedProfile represents a resolved profile with all settings applied
type ResolvedProfile struct {
	ID           string
	APIKey       string
	BaseURL      string
	Model        string // for Grok only
	ProfileSource string // "cli", "env", "toml"
}

// Common errors
var (
	ErrNoAPIKey       = errors.New("no API key configured")
	ErrConfigParse    = errors.New("config parse error")
	ErrInvalidJSON    = errors.New("invalid JSON")
)
