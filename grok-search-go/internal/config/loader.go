package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	defaultConfigDir  = ".config/ai-skills"
	defaultConfigFile = "grok-search.toml"
)

// Error is a structured configuration error.
type Error struct {
	Code   string
	Detail string
	Err    error
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return e.Code + ": " + e.Detail
	}
	return e.Code
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Load resolves Grok Search configuration from CLI options, environment, TOML, and defaults.
func Load(opts Options) (*ResolvedConfig, error) {
	path, isDefault, err := resolveConfigPath(opts.ConfigPath)
	if err != nil {
		return nil, &Error{Code: "invalid_config", Detail: err.Error(), Err: err}
	}

	if isDefault {
		if err := EnsureTemplate(path); err != nil {
			return nil, &Error{Code: "invalid_config", Detail: err.Error(), Err: err}
		}
	}

	cfg := Config{Cooldown: CooldownConfig{Enabled: true}}
	if _, err := os.Stat(path); err == nil {
		if _, err := toml.DecodeFile(path, &cfg); err != nil {
			return nil, &Error{Code: "invalid_config", Detail: err.Error(), Err: err}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, &Error{Code: "invalid_config", Detail: err.Error(), Err: err}
	} else if !isDefault {
		return nil, &Error{Code: "invalid_config", Detail: fmt.Sprintf("config file not found: %s", path), Err: err}
	}

	applyDefaults(&cfg)
	applyEnv(&cfg)
	applyOptions(&cfg, opts)
	if err := applyJSONOverrides(&cfg, opts); err != nil {
		return nil, err
	}

	profiles := ResolveProfiles(cfg, opts)
	if len(profiles) == 0 {
		return nil, &Error{Code: "missing_api_key", Detail: "Pass --api-key, set GROK_API_KEY/GROK_API_KEYS, or configure ~/.config/ai-skills/grok-search.toml"}
	}

	return &ResolvedConfig{
		ConfigPath:   path,
		ConfigPaths:  []string{path},
		BaseURL:      cfg.BaseURL,
		Model:        cfg.Model,
		Timeout:      cfg.Timeout,
		Profiles:     profiles,
		ExtraBody:    nonNilMap(cfg.ExtraBody),
		ExtraHeaders: nonNilMap(cfg.ExtraHeaders),
		Cooldown:     cfg.Cooldown,
	}, nil
}

func resolveConfigPath(explicit string) (string, bool, error) {
	if strings.TrimSpace(explicit) != "" {
		return explicit, false, nil
	}
	if env := strings.TrimSpace(os.Getenv("GROK_CONFIG")); env != "" {
		return env, false, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", true, err
	}
	return filepath.Join(home, defaultConfigDir, defaultConfigFile), true, nil
}

func applyDefaults(cfg *Config) {
	cfg.BaseURL = firstNonEmpty(cfg.BaseURL, DefaultBaseURL)
	cfg.Model = firstNonEmpty(cfg.Model, DefaultModel)
	if cfg.Timeout <= 0 {
		cfg.Timeout = DefaultTimeout
	}
	if cfg.ExtraBody == nil {
		cfg.ExtraBody = map[string]any{}
	}
	if cfg.ExtraHeaders == nil {
		cfg.ExtraHeaders = map[string]any{}
	}
	cfg.Cooldown = normalizeCooldown(cfg.Cooldown)
}

func applyEnv(cfg *Config) {
	if value := strings.TrimSpace(os.Getenv("GROK_BASE_URL")); value != "" {
		cfg.BaseURL = value
	}
	if value := strings.TrimSpace(os.Getenv("GROK_MODEL")); value != "" {
		cfg.Model = value
	}
	if value := strings.TrimSpace(os.Getenv("GROK_TIMEOUT")); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			cfg.Timeout = parsed
		}
	}
}

func applyOptions(cfg *Config, opts Options) {
	if value := strings.TrimSpace(opts.BaseURL); value != "" {
		cfg.BaseURL = value
	}
	if value := strings.TrimSpace(opts.Model); value != "" {
		cfg.Model = value
	}
	if opts.Timeout > 0 {
		cfg.Timeout = opts.Timeout
	}
}

// ResolveProfiles resolves runtime profiles using the approved credential priority.
func ResolveProfiles(cfg Config, opts Options) []ResolvedProfile {
	globalBaseURL := firstNonEmpty(cfg.BaseURL, DefaultBaseURL)
	globalModel := firstNonEmpty(cfg.Model, DefaultModel)

	if key := normalizeAPIKey(opts.APIKey); key != "" {
		return filterProfile([]ResolvedProfile{{
			ID:      "cli",
			APIKey:  key,
			Source:  "--api-key",
			BaseURL: globalBaseURL,
			Model:   globalModel,
		}}, opts.ProfileID)
	}

	if envKeys := profilesFromEnvKeys(globalBaseURL, globalModel); len(envKeys) > 0 {
		return filterProfile(envKeys, opts.ProfileID)
	}

	if key := normalizeAPIKey(os.Getenv("GROK_API_KEY")); key != "" {
		return filterProfile([]ResolvedProfile{{
			ID:      "env",
			APIKey:  key,
			Source:  "GROK_API_KEY",
			BaseURL: globalBaseURL,
			Model:   globalModel,
		}}, opts.ProfileID)
	}

	profiles := make([]ResolvedProfile, 0, len(cfg.Profiles))
	for i, profile := range cfg.Profiles {
		if profile.Enabled != nil && !*profile.Enabled {
			continue
		}
		key := normalizeAPIKey(profile.APIKey)
		if key == "" {
			continue
		}
		id := strings.TrimSpace(profile.ID)
		if id == "" {
			id = fmt.Sprintf("profile-%d", i+1)
		}
		profiles = append(profiles, ResolvedProfile{
			ID:      id,
			APIKey:  key,
			Source:  "config.profiles",
			BaseURL: firstNonEmpty(profile.BaseURL, globalBaseURL),
			Model:   firstNonEmpty(profile.Model, globalModel),
		})
	}
	return filterProfile(profiles, opts.ProfileID)
}

func profilesFromEnvKeys(globalBaseURL, globalModel string) []ResolvedProfile {
	raw := strings.TrimSpace(os.Getenv("GROK_API_KEYS"))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	profiles := make([]ResolvedProfile, 0, len(parts))
	for i, part := range parts {
		key := normalizeAPIKey(part)
		if key == "" {
			continue
		}
		profiles = append(profiles, ResolvedProfile{
			ID:      fmt.Sprintf("env-%d", i+1),
			APIKey:  key,
			Source:  "GROK_API_KEYS",
			BaseURL: globalBaseURL,
			Model:   globalModel,
		})
	}
	return profiles
}

func filterProfile(profiles []ResolvedProfile, profileID string) []ResolvedProfile {
	profileID = strings.TrimSpace(profileID)
	if profileID == "" {
		return profiles
	}
	for _, profile := range profiles {
		if profile.ID == profileID {
			return []ResolvedProfile{profile}
		}
	}
	return nil
}

func normalizeAPIKey(apiKey string) string {
	key := strings.TrimSpace(apiKey)
	if key == "" {
		return ""
	}
	if isPlaceholder(key) {
		return ""
	}
	return key
}

func isPlaceholder(value string) bool {
	upper := strings.ToUpper(strings.TrimSpace(value))
	placeholders := map[string]struct{}{
		"YOUR_API_KEY":             {},
		"YOUR_GROK_API_KEY":        {},
		"YOUR_BACKUP_GROK_API_KEY": {},
		"API_KEY":                  {},
		"CHANGE_ME":                {},
		"REPLACE_ME":               {},
		"<YOUR_GROK_API_KEY>":      {},
	}
	_, ok := placeholders[upper]
	return ok
}

func applyJSONOverrides(cfg *Config, opts Options) error {
	for _, item := range []struct {
		label string
		raw   string
		dest  *map[string]any
	}{
		{label: "GROK_EXTRA_BODY_JSON", raw: os.Getenv("GROK_EXTRA_BODY_JSON"), dest: &cfg.ExtraBody},
		{label: "GROK_EXTRA_HEADERS_JSON", raw: os.Getenv("GROK_EXTRA_HEADERS_JSON"), dest: &cfg.ExtraHeaders},
		{label: "--extra-body-json", raw: opts.ExtraBodyJSON, dest: &cfg.ExtraBody},
		{label: "--extra-headers-json", raw: opts.ExtraHeadersJSON, dest: &cfg.ExtraHeaders},
	} {
		parsed, err := parseJSONObject(item.raw, item.label)
		if err != nil {
			return err
		}
		mergeMap(*item.dest, parsed)
	}
	return nil
}

func parseJSONObject(raw, label string) (map[string]any, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, &Error{Code: "invalid_json", Detail: fmt.Sprintf("%s: %s", label, err.Error()), Err: err}
	}
	if parsed == nil {
		return nil, &Error{Code: "invalid_json", Detail: fmt.Sprintf("%s must be a JSON object", label)}
	}
	return parsed, nil
}

func mergeMap(dst map[string]any, src map[string]any) {
	for key, value := range src {
		dst[key] = value
	}
}

func normalizeCooldown(c CooldownConfig) CooldownConfig {
	if c.StateFile == "" {
		c.StateFile = "runtime/cooldowns.json"
	}
	if c.DefaultMinutes <= 0 {
		c.DefaultMinutes = 15
	}
	if c.RateLimitMinutes <= 0 {
		c.RateLimitMinutes = 20
	}
	if c.QuotaMinutes <= 0 {
		c.QuotaMinutes = 60
	}
	if c.AuthMinutes <= 0 {
		c.AuthMinutes = 360
	}
	return c
}

func firstNonEmpty(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func nonNilMap(value map[string]any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	return value
}
