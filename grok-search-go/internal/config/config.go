package config

const (
	DefaultBaseURL = "https://api.x.ai"
	DefaultModel   = "grok-4.1-fast"
	DefaultTimeout = 120
)

// Config represents the TOML configuration file.
type Config struct {
	BaseURL      string         `toml:"base_url"`
	Model        string         `toml:"model"`
	Timeout      int            `toml:"timeout"`
	Profiles     []Profile      `toml:"profiles"`
	ExtraBody    map[string]any `toml:"extra_body"`
	ExtraHeaders map[string]any `toml:"extra_headers"`
	Cooldown     CooldownConfig `toml:"cooldown"`
}

// Profile represents an API key profile with optional endpoint overrides.
type Profile struct {
	ID      string `toml:"id"`
	APIKey  string `toml:"api_key"`
	BaseURL string `toml:"base_url"`
	Model   string `toml:"model"`
	Enabled *bool  `toml:"enabled"`
}

// CooldownConfig controls temporary suppression of failing profiles.
type CooldownConfig struct {
	Enabled          bool   `toml:"enabled"`
	StateFile        string `toml:"state_file"`
	DefaultMinutes   int    `toml:"default_minutes"`
	RateLimitMinutes int    `toml:"rate_limit_minutes"`
	QuotaMinutes     int    `toml:"quota_minutes"`
	AuthMinutes      int    `toml:"auth_minutes"`
}

// Options contains CLI flag values and process-level overrides for loading config.
type Options struct {
	ConfigPath       string
	APIKey           string
	BaseURL          string
	Model            string
	Timeout          int
	ProfileID        string
	ExtraBodyJSON    string
	ExtraHeadersJSON string
}

// ResolvedConfig is the effective configuration after precedence is applied.
type ResolvedConfig struct {
	ConfigPath   string
	ConfigPaths  []string
	BaseURL      string
	Model        string
	Timeout      int
	Profiles     []ResolvedProfile
	ExtraBody    map[string]any
	ExtraHeaders map[string]any
	Cooldown     CooldownConfig
}

// ResolvedProfile is the effective runtime profile after overrides are applied.
type ResolvedProfile struct {
	ID      string
	APIKey  string
	Source  string
	BaseURL string
	Model   string
}
