package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/venomwise/skill-factory/exa-search/internal/debug"
)

const (
	defaultConfigDir  = ".config/ai-skills"
	defaultConfigFile = "exa-search.toml"
	defaultBaseURL    = "https://api.exa.ai"
	defaultTimeout    = 30
)

// Load reads configuration from the specified path or default location
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("toml")

	debug.LogConfigResolution("Loading configuration...")

	// Determine config file path
	if configPath != "" {
		// Use explicit path
		debug.LogConfigResolution("Using explicit config path: %s", configPath)
		v.SetConfigFile(configPath)
	} else {
		// Use default path: ~/.config/ai-skills/exa-search.toml
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir := filepath.Join(home, defaultConfigDir)
		configFile := filepath.Join(configDir, defaultConfigFile)
		debug.LogConfigResolution("Using default config path: %s", configFile)
		
		// Check if config file exists, if not create template
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			debug.LogConfigResolution("Config file not found, creating template")
			if err := createConfigTemplate(configFile); err != nil {
				return nil, fmt.Errorf("failed to create config template: %w", err)
			}
		}
		
		v.SetConfigFile(configFile)
	}

	// Set defaults
	v.SetDefault("base_url", defaultBaseURL)
	v.SetDefault("timeout", defaultTimeout)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, return empty config with defaults
			return &Config{
				BaseURL: defaultBaseURL,
				Timeout: defaultTimeout,
			}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply defaults if not set
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}

	// Merge environment variables
	if err := mergeEnvVars(&cfg); err != nil {
		return nil, err
	}

	// Filter out profiles with placeholder keys
	validProfiles := make([]Profile, 0, len(cfg.Profiles))
	for _, p := range cfg.Profiles {
		if !isPlaceholder(p.APIKey) {
			validProfiles = append(validProfiles, p)
		}
	}
	cfg.Profiles = validProfiles

	return &cfg, nil
}

// mergeEnvVars merges environment variables into config
func mergeEnvVars(cfg *Config) error {
	// EXA_API_KEYS: comma-separated list of keys
	if keys := os.Getenv("EXA_API_KEYS"); keys != "" {
		debug.LogConfigResolution("Found EXA_API_KEYS environment variable")
		parts := strings.Split(keys, ",")
		for i, key := range parts {
			key = strings.TrimSpace(key)
			if key != "" && !isPlaceholder(key) {
				cfg.Profiles = append(cfg.Profiles, Profile{
					ID:     fmt.Sprintf("env-%d", i+1),
					APIKey: key,
				})
				debug.LogConfigResolution("Added profile from EXA_API_KEYS: env-%d (key: %s)", i+1, debug.RedactAPIKey(key))
			}
		}
	}

	// EXA_API_KEY: single key
	if key := os.Getenv("EXA_API_KEY"); key != "" && !isPlaceholder(key) {
		debug.LogConfigResolution("Found EXA_API_KEY environment variable")
		// Only add if no profiles from EXA_API_KEYS
		hasEnvProfiles := false
		for _, p := range cfg.Profiles {
			if strings.HasPrefix(p.ID, "env-") {
				hasEnvProfiles = true
				break
			}
		}
		if !hasEnvProfiles {
			cfg.Profiles = append(cfg.Profiles, Profile{
				ID:     "env",
				APIKey: key,
			})
		}
	}

	// EXA_BASE_URL: override base URL
	if baseURL := os.Getenv("EXA_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	// EXA_TIMEOUT: override timeout
	if timeoutStr := os.Getenv("EXA_TIMEOUT"); timeoutStr != "" {
		var timeout int
		if _, err := fmt.Sscanf(timeoutStr, "%d", &timeout); err == nil && timeout > 0 {
			cfg.Timeout = timeout
		}
	}

	return nil
}

// ApplyFlags applies CLI flag overrides to the configuration
func ApplyFlags(cfg *Config, apiKey, profileID string) {
	// If --api-key is provided, create a single CLI profile and clear others
	if apiKey != "" && !isPlaceholder(apiKey) {
		cfg.Profiles = []Profile{
			{
				ID:     "cli",
				APIKey: apiKey,
			},
		}
		return
	}

	// If --profile is provided, filter to that profile only
	if profileID != "" {
		for _, p := range cfg.Profiles {
			if p.ID == profileID {
				cfg.Profiles = []Profile{p}
				return
			}
		}
		// Profile not found, clear profiles (will trigger missing API key error)
		cfg.Profiles = nil
	}
}

// createConfigTemplate creates a template config file with comments
func createConfigTemplate(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	template := `# Exa Search Configuration
# Get your API key at: https://exa.ai/

# API key profiles for failover support
[[profiles]]
id = "main"
api_key = "YOUR_EXA_API_KEY"  # Replace with your actual API key
# base_url = "https://api.exa.ai"  # Optional: override base URL

# Add additional profiles for failover
# [[profiles]]
# id = "backup"
# api_key = "YOUR_BACKUP_KEY"

# Global settings
# base_url = "https://api.exa.ai"  # Default API endpoint
# timeout = 30  # Request timeout in seconds
`

	return os.WriteFile(path, []byte(template), 0644)
}

// isPlaceholder checks if an API key is a placeholder or empty
func isPlaceholder(key string) bool {
	if key == "" {
		return true
	}
	
	placeholders := []string{
		"YOUR_EXA_API_KEY",
		"YOUR_API_KEY",
		"YOUR_BACKUP_KEY",
		"API_KEY",
		"CHANGE_ME",
		"REPLACE_ME",
		"<YOUR_EXA_API_KEY>",
	}
	
	for _, placeholder := range placeholders {
		if key == placeholder {
			return true
		}
	}
	
	return false
}
