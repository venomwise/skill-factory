package config

import "github.com/venomwise/skill-factory/db-explorer/internal/db"

const (
	DefaultProjectConfig = ".db-explorer.toml"
	DefaultConfigDir     = ".config/ai-skills"
	DefaultConfigFile    = "db-explorer.toml"
)

// Config represents a db-explorer TOML configuration file.
type Config struct {
	DefaultProfile string    `toml:"default_profile"`
	Profiles       []Profile `toml:"profiles"`
}

// Profile represents one named database connection.
type Profile struct {
	ID     string `toml:"id"`
	DB     string `toml:"db"`
	URL    string `toml:"url"`
	URLEnv string `toml:"url_env"`
}

// Options contains CLI flag values and process-level overrides.
type Options struct {
	ConfigPath string
	ProfileID  string
	DB         string
	URL        string
	URLEnv     string
}

// ResolvedConnection is the effective connection after precedence is applied.
type ResolvedConnection struct {
	DB      db.Type
	URL     string
	URLEnv  string
	Profile string
	Source  string
}
