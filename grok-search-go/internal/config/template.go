package config

import (
	"errors"
	"os"
	"path/filepath"
)

const templateConfig = `# Grok Search Configuration
# Get your xAI API key from https://console.x.ai/

base_url = "https://api.x.ai"
model = "grok-4.1-fast"
timeout = 120

[[profiles]]
id = "main"
api_key = "YOUR_GROK_API_KEY"

# Optional backup profile
# [[profiles]]
# id = "backup"
# api_key = "YOUR_BACKUP_GROK_API_KEY"

# Optional OpenAI-compatible proxy profile
# [[profiles]]
# id = "proxy"
# api_key = "YOUR_PROXY_KEY"
# base_url = "https://your-compatible-endpoint.example"
# model = "grok-custom-model"

[extra_body]
# Add endpoint-specific JSON body fields here.

[extra_headers]
# X-Custom-Header = "value"

[cooldown]
enabled = true
state_file = "runtime/cooldowns.json"
default_minutes = 15
rate_limit_minutes = 20
quota_minutes = 60
auth_minutes = 360
`

// EnsureTemplate creates a template config at path when it does not already exist.
func EnsureTemplate(path string) error {
	if path == "" {
		return nil
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(templateConfig), 0o644)
}
