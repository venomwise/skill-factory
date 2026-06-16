package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultConfigDir  = ".config/ai-skills"
	defaultConfigFile = "web-access.toml"

	defaultExaBaseURL    = "https://api.exa.ai"
	defaultGrokBaseURL   = "https://api.x.ai/v1"
	defaultGrokModel     = "grok-beta"
	defaultTimeout       = 30
	defaultGrokTimeout   = 60
)

// EnsureTemplate creates a template config file if it doesn't exist
func EnsureTemplate(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return nil // File exists, no need to create
	}

	template := `# Web Access Configuration
# Get your API keys at:
# - Exa: https://exa.ai/
# - Grok: https://x.ai/api

[exa]
base_url = "https://api.exa.ai"
timeout = 30

[[exa.profiles]]
id = "primary"
api_key = "YOUR_EXA_API_KEY"

# Uncomment to add backup profiles
# [[exa.profiles]]
# id = "backup"
# api_key = "YOUR_BACKUP_EXA_KEY"

[grok]
base_url = "https://api.x.ai/v1"
model = "grok-beta"
timeout = 60

[[grok.profiles]]
id = "primary"
api_key = "YOUR_GROK_API_KEY"

# Uncomment to add backup profiles
# [[grok.profiles]]
# id = "backup"
# api_key = "YOUR_BACKUP_GROK_KEY"

# Optional: Extra body and headers for Grok requests
# [grok.extra_body]
# temperature = 0.7

# [grok.extra_headers]
# X-Custom-Header = "value"
`

	return os.WriteFile(path, []byte(template), 0644)
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, defaultConfigDir, defaultConfigFile), nil
}
