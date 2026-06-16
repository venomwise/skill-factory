package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Options represents configuration loading options
type Options struct {
	ConfigPath string
	ProfileID  string
	Timeout    int

	// Exa CLI flags
	ExaAPIKey string

	// Grok CLI flags
	GrokAPIKey       string
	GrokModel        string
	ExtraBodyJSON    string
	ExtraHeadersJSON string
}

// Load loads and resolves configuration with precedence:
// CLI flags > WEB_ACCESS_* env > provider-specific env > TOML config > defaults
func Load(opts Options) (*ResolvedConfig, error) {
	// Resolve config file path
	configPath := opts.ConfigPath
	if configPath == "" {
		// Check WEB_ACCESS_CONFIG env var
		if envPath := os.Getenv("WEB_ACCESS_CONFIG"); envPath != "" {
			configPath = envPath
		} else {
			// Use default path
			defaultPath, err := GetDefaultConfigPath()
			if err != nil {
				return nil, err
			}
			configPath = defaultPath
		}
	}

	// Ensure template exists
	if err := EnsureTemplate(configPath); err != nil {
		return nil, fmt.Errorf("failed to create config template: %w", err)
	}

	// Load TOML config
	var cfg Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config doesn't exist, use empty config
			cfg = Config{}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("%w: %s: %v", ErrConfigParse, configPath, err)
		}
	}

	// Resolve Exa provider config
	exaResolved, err := resolveExaProvider(cfg.Exa, opts)
	if err != nil {
		return nil, err
	}

	// Resolve Grok provider config
	grokResolved, err := resolveGrokProvider(cfg.Grok, opts)
	if err != nil {
		return nil, err
	}

	return &ResolvedConfig{
		Exa:  exaResolved,
		Grok: grokResolved,
	}, nil
}

// resolveExaProvider resolves Exa provider configuration with precedence
func resolveExaProvider(cfg ProviderConfig, opts Options) (ResolvedProviderConfig, error) {
	resolved := ResolvedProviderConfig{
		BaseURL: defaultExaBaseURL,
		Timeout: defaultTimeout,
	}

	// Apply TOML config
	if cfg.BaseURL != "" {
		resolved.BaseURL = cfg.BaseURL
	}
	if cfg.Timeout > 0 {
		resolved.Timeout = cfg.Timeout
	}

	// Apply provider-specific env vars
	if baseURL := os.Getenv("EXA_BASE_URL"); baseURL != "" {
		resolved.BaseURL = baseURL
	}
	if timeoutStr := os.Getenv("EXA_TIMEOUT"); timeoutStr != "" {
		var timeout int
		if _, err := fmt.Sscanf(timeoutStr, "%d", &timeout); err == nil && timeout > 0 {
			resolved.Timeout = timeout
		}
	}

	// Apply CLI flags
	if opts.Timeout > 0 {
		resolved.Timeout = opts.Timeout
	}

	// Resolve profiles with precedence
	var profiles []ResolvedProfile

	// 1. TOML profiles
	for _, p := range cfg.Profiles {
		if !isPlaceholder(p.APIKey) {
			profiles = append(profiles, ResolvedProfile{
				ID:            p.ID,
				APIKey:        p.APIKey,
				BaseURL:       getOrDefault(p.BaseURL, resolved.BaseURL),
				ProfileSource: "toml",
			})
		}
	}

	// 2. Provider-specific env vars (EXA_API_KEYS, EXA_API_KEY)
	if keys := os.Getenv("EXA_API_KEYS"); keys != "" {
		parts := strings.Split(keys, ",")
		for i, key := range parts {
			key = strings.TrimSpace(key)
			if key != "" && !isPlaceholder(key) {
				profiles = append(profiles, ResolvedProfile{
					ID:            fmt.Sprintf("env-exa-%d", i+1),
					APIKey:        key,
					BaseURL:       resolved.BaseURL,
					ProfileSource: "env",
				})
			}
		}
	} else if key := os.Getenv("EXA_API_KEY"); key != "" && !isPlaceholder(key) {
		profiles = append(profiles, ResolvedProfile{
			ID:            "env-exa",
			APIKey:        key,
			BaseURL:       resolved.BaseURL,
			ProfileSource: "env",
		})
	}

	// 3. WEB_ACCESS_EXA_API_KEY env var (higher priority)
	if key := os.Getenv("WEB_ACCESS_EXA_API_KEY"); key != "" && !isPlaceholder(key) {
		profiles = append(profiles, ResolvedProfile{
			ID:            "web-access-exa",
			APIKey:        key,
			BaseURL:       resolved.BaseURL,
			ProfileSource: "env",
		})
	}

	// 4. CLI flag (highest priority)
	if opts.ExaAPIKey != "" && !isPlaceholder(opts.ExaAPIKey) {
		profiles = []ResolvedProfile{{
			ID:            "cli",
			APIKey:        opts.ExaAPIKey,
			BaseURL:       resolved.BaseURL,
			ProfileSource: "cli",
		}}
	}

	// Filter by --profile if specified
	if opts.ProfileID != "" && len(profiles) > 0 {
		filtered := []ResolvedProfile{}
		for _, p := range profiles {
			if p.ID == opts.ProfileID {
				filtered = append(filtered, p)
			}
		}
		profiles = filtered
	}

	resolved.Profiles = profiles
	return resolved, nil
}

// resolveGrokProvider resolves Grok provider configuration with precedence
func resolveGrokProvider(cfg ProviderConfig, opts Options) (ResolvedProviderConfig, error) {
	resolved := ResolvedProviderConfig{
		BaseURL:      defaultGrokBaseURL,
		Timeout:      defaultGrokTimeout,
		Model:        defaultGrokModel,
		ExtraBody:    make(map[string]interface{}),
		ExtraHeaders: make(map[string]string),
	}

	// Apply TOML config
	if cfg.BaseURL != "" {
		resolved.BaseURL = cfg.BaseURL
	}
	if cfg.Timeout > 0 {
		resolved.Timeout = cfg.Timeout
	}
	if cfg.Model != "" {
		resolved.Model = cfg.Model
	}
	if cfg.ExtraBody != nil {
		resolved.ExtraBody = cfg.ExtraBody
	}
	if cfg.ExtraHeaders != nil {
		resolved.ExtraHeaders = cfg.ExtraHeaders
	}

	// Apply provider-specific env vars
	if baseURL := os.Getenv("GROK_BASE_URL"); baseURL != "" {
		resolved.BaseURL = baseURL
	}
	if model := os.Getenv("GROK_MODEL"); model != "" {
		resolved.Model = model
	}
	if timeoutStr := os.Getenv("GROK_TIMEOUT"); timeoutStr != "" {
		var timeout int
		if _, err := fmt.Sscanf(timeoutStr, "%d", &timeout); err == nil && timeout > 0 {
			resolved.Timeout = timeout
		}
	}

	// Merge GROK_EXTRA_BODY_JSON
	if extraBodyJSON := os.Getenv("GROK_EXTRA_BODY_JSON"); extraBodyJSON != "" {
		var extraBody map[string]interface{}
		if err := json.Unmarshal([]byte(extraBodyJSON), &extraBody); err != nil {
			return resolved, fmt.Errorf("%w: GROK_EXTRA_BODY_JSON: %v", ErrInvalidJSON, err)
		}
		for k, v := range extraBody {
			resolved.ExtraBody[k] = v
		}
	}

	// Merge GROK_EXTRA_HEADERS_JSON
	if extraHeadersJSON := os.Getenv("GROK_EXTRA_HEADERS_JSON"); extraHeadersJSON != "" {
		var extraHeaders map[string]string
		if err := json.Unmarshal([]byte(extraHeadersJSON), &extraHeaders); err != nil {
			return resolved, fmt.Errorf("%w: GROK_EXTRA_HEADERS_JSON: %v", ErrInvalidJSON, err)
		}
		for k, v := range extraHeaders {
			resolved.ExtraHeaders[k] = v
		}
	}

	// Apply CLI flags
	if opts.Timeout > 0 {
		resolved.Timeout = opts.Timeout
	}
	if opts.GrokModel != "" {
		resolved.Model = opts.GrokModel
	}

	// Merge CLI extra body JSON
	if opts.ExtraBodyJSON != "" {
		var extraBody map[string]interface{}
		if err := json.Unmarshal([]byte(opts.ExtraBodyJSON), &extraBody); err != nil {
			return resolved, fmt.Errorf("%w: --extra-body-json: %v", ErrInvalidJSON, err)
		}
		for k, v := range extraBody {
			resolved.ExtraBody[k] = v
		}
	}

	// Merge CLI extra headers JSON
	if opts.ExtraHeadersJSON != "" {
		var extraHeaders map[string]string
		if err := json.Unmarshal([]byte(opts.ExtraHeadersJSON), &extraHeaders); err != nil {
			return resolved, fmt.Errorf("%w: --extra-headers-json: %v", ErrInvalidJSON, err)
		}
		for k, v := range extraHeaders {
			resolved.ExtraHeaders[k] = v
		}
	}

	// Resolve profiles with precedence
	var profiles []ResolvedProfile

	// 1. TOML profiles
	for _, p := range cfg.Profiles {
		if !isPlaceholder(p.APIKey) {
			profiles = append(profiles, ResolvedProfile{
				ID:            p.ID,
				APIKey:        p.APIKey,
				BaseURL:       getOrDefault(p.BaseURL, resolved.BaseURL),
				Model:         getOrDefault(p.Model, resolved.Model),
				ProfileSource: "toml",
			})
		}
	}

	// 2. Provider-specific env vars (GROK_API_KEYS, GROK_API_KEY)
	if keys := os.Getenv("GROK_API_KEYS"); keys != "" {
		parts := strings.Split(keys, ",")
		for i, key := range parts {
			key = strings.TrimSpace(key)
			if key != "" && !isPlaceholder(key) {
				profiles = append(profiles, ResolvedProfile{
					ID:            fmt.Sprintf("env-grok-%d", i+1),
					APIKey:        key,
					BaseURL:       resolved.BaseURL,
					Model:         resolved.Model,
					ProfileSource: "env",
				})
			}
		}
	} else if key := os.Getenv("GROK_API_KEY"); key != "" && !isPlaceholder(key) {
		profiles = append(profiles, ResolvedProfile{
			ID:            "env-grok",
			APIKey:        key,
			BaseURL:       resolved.BaseURL,
			Model:         resolved.Model,
			ProfileSource: "env",
		})
	}

	// 3. WEB_ACCESS_GROK_API_KEY env var (higher priority)
	if key := os.Getenv("WEB_ACCESS_GROK_API_KEY"); key != "" && !isPlaceholder(key) {
		profiles = append(profiles, ResolvedProfile{
			ID:            "web-access-grok",
			APIKey:        key,
			BaseURL:       resolved.BaseURL,
			Model:         resolved.Model,
			ProfileSource: "env",
		})
	}

	// 4. CLI flag (highest priority)
	if opts.GrokAPIKey != "" && !isPlaceholder(opts.GrokAPIKey) {
		profiles = []ResolvedProfile{{
			ID:            "cli",
			APIKey:        opts.GrokAPIKey,
			BaseURL:       resolved.BaseURL,
			Model:         resolved.Model,
			ProfileSource: "cli",
		}}
	}

	// Filter by --profile if specified
	if opts.ProfileID != "" && len(profiles) > 0 {
		filtered := []ResolvedProfile{}
		for _, p := range profiles {
			if p.ID == opts.ProfileID {
				filtered = append(filtered, p)
			}
		}
		profiles = filtered
	}

	resolved.Profiles = profiles
	return resolved, nil
}

// isPlaceholder checks if an API key is a placeholder or empty
func isPlaceholder(key string) bool {
	if key == "" {
		return true
	}

	placeholders := []string{
		"YOUR_EXA_API_KEY",
		"YOUR_GROK_API_KEY",
		"YOUR_API_KEY",
		"YOUR_BACKUP_EXA_KEY",
		"YOUR_BACKUP_GROK_KEY",
		"YOUR_BACKUP_KEY",
		"API_KEY",
		"CHANGE_ME",
		"REPLACE_ME",
	}

	for _, placeholder := range placeholders {
		if key == placeholder {
			return true
		}
	}

	return false
}

// getOrDefault returns value if not empty, otherwise returns defaultValue
func getOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}
