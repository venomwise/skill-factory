package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/venomwise/skill-factory/db-explorer/internal/db"
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

func (e *Error) Unwrap() error { return e.Err }

// LoadConfig loads a TOML config file from path.
func LoadConfig(path string) (Config, error) {
	var cfg Config
	if strings.TrimSpace(path) == "" {
		return cfg, &Error{Code: "INVALID_CONFIG", Detail: "config path is empty"}
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, &Error{Code: "INVALID_CONFIG", Detail: err.Error(), Err: err}
	}
	return cfg, nil
}

// FindProjectConfig returns the project config path if it exists.
func FindProjectConfig(startDir string) (string, bool, error) {
	if strings.TrimSpace(startDir) == "" {
		var err error
		startDir, err = os.Getwd()
		if err != nil {
			return "", false, err
		}
	}
	path := filepath.Join(startDir, DefaultProjectConfig)
	_, err := os.Stat(path)
	if err == nil {
		return path, true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	return "", false, err
}

// LoadGlobalConfig loads the default global config if present.
func LoadGlobalConfig() (Config, string, bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, "", false, err
	}
	path := filepath.Join(home, DefaultConfigDir, DefaultConfigFile)
	_, err = os.Stat(path)
	if err == nil {
		cfg, err := LoadConfig(path)
		return cfg, path, err == nil, err
	}
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, path, false, nil
	}
	return Config{}, path, false, err
}

// ResolveConnection resolves the effective connection using CLI, project config,
// global config, and environment fallback precedence.
func ResolveConnection(opts Options) (ResolvedConnection, error) {
	if conn, ok, err := resolveCLI(opts); ok || err != nil {
		return conn, err
	}

	projectCfg, projectSource, projectOK, err := loadProjectConfig(opts.ConfigPath)
	if err != nil {
		return ResolvedConnection{}, err
	}
	globalCfg, globalSource, globalOK, err := LoadGlobalConfig()
	if err != nil {
		return ResolvedConnection{}, &Error{Code: "INVALID_CONFIG", Detail: err.Error(), Err: err}
	}

	if strings.TrimSpace(opts.ProfileID) != "" {
		if projectOK {
			if conn, ok, err := resolveProfile(projectCfg, opts.ProfileID, projectSource); ok || err != nil {
				return conn, err
			}
		}
		if globalOK {
			if conn, ok, err := resolveProfile(globalCfg, opts.ProfileID, globalSource); ok || err != nil {
				return conn, err
			}
		}
		return ResolvedConnection{}, &Error{Code: "PROFILE_NOT_FOUND", Detail: opts.ProfileID}
	}

	if projectOK && strings.TrimSpace(projectCfg.DefaultProfile) != "" {
		if conn, ok, err := resolveProfile(projectCfg, projectCfg.DefaultProfile, projectSource); ok || err != nil {
			return conn, err
		}
	}
	if globalOK && strings.TrimSpace(globalCfg.DefaultProfile) != "" {
		if conn, ok, err := resolveProfile(globalCfg, globalCfg.DefaultProfile, globalSource); ok || err != nil {
			return conn, err
		}
	}

	if conn, ok := resolveEnvFallback(); ok {
		return conn, nil
	}
	return ResolvedConnection{}, &Error{Code: "MISSING_CONNECTION", Detail: "provide --db/--url, --profile, config, or environment fallback"}
}

func loadProjectConfig(explicitPath string) (Config, string, bool, error) {
	if strings.TrimSpace(explicitPath) != "" {
		cfg, err := LoadConfig(explicitPath)
		if err != nil {
			return Config{}, explicitPath, false, err
		}
		return cfg, explicitPath, true, nil
	}
	path, ok, err := FindProjectConfig("")
	if err != nil || !ok {
		return Config{}, path, false, err
	}
	cfg, err := LoadConfig(path)
	return cfg, path, err == nil, err
}

func resolveCLI(opts Options) (ResolvedConnection, bool, error) {
	dbType := strings.TrimSpace(opts.DB)
	url := strings.TrimSpace(opts.URL)
	urlEnv := strings.TrimSpace(opts.URLEnv)
	if dbType == "" && url == "" && urlEnv == "" {
		return ResolvedConnection{}, false, nil
	}
	parsed, err := parseDBType(dbType)
	if err != nil {
		return ResolvedConnection{}, true, err
	}
	resolvedURL, err := resolveURLValue(url, urlEnv)
	if err != nil {
		return ResolvedConnection{}, true, err
	}
	return ResolvedConnection{DB: parsed, URL: resolvedURL, URLEnv: urlEnv, Source: "cli"}, true, nil
}

func resolveProfile(cfg Config, id, source string) (ResolvedConnection, bool, error) {
	id = strings.TrimSpace(id)
	for _, profile := range cfg.Profiles {
		if strings.TrimSpace(profile.ID) != id {
			continue
		}
		dbType, err := parseDBType(profile.DB)
		if err != nil {
			return ResolvedConnection{}, true, err
		}
		resolvedURL, err := resolveURLValue(profile.URL, profile.URLEnv)
		if err != nil {
			if cfgErr, ok := err.(*Error); ok && cfgErr.Code == "MISSING_CONNECTION" {
				cfgErr.Detail = fmt.Sprintf("profile %q has no url or url_env", id)
			}
			return ResolvedConnection{}, true, err
		}
		return ResolvedConnection{DB: dbType, URL: resolvedURL, URLEnv: strings.TrimSpace(profile.URLEnv), Profile: id, Source: source}, true, nil
	}
	return ResolvedConnection{}, false, nil
}

func resolveURLValue(rawURL, rawURLEnv string) (string, error) {
	url := strings.TrimSpace(rawURL)
	urlEnv := strings.TrimSpace(rawURLEnv)
	if url != "" {
		return url, nil
	}
	if urlEnv != "" {
		return ResolveURLEnv(urlEnv)
	}
	return "", &Error{Code: "MISSING_CONNECTION", Detail: "url or url_env is required"}
}

func resolveEnvFallback() (ResolvedConnection, bool) {
	for _, item := range []struct {
		Name string
		DB   db.Type
	}{
		{Name: "POSTGRES_URL", DB: db.TypePostgres},
		{Name: "MYSQL_URL", DB: db.TypeMySQL},
		{Name: "DATABASE_URL"},
		{Name: "DB_URL"},
	} {
		value := strings.TrimSpace(os.Getenv(item.Name))
		if value == "" {
			continue
		}
		dbType := item.DB
		if dbType == "" {
			dbType = inferDBType(value)
		}
		if dbType == "" {
			continue
		}
		return ResolvedConnection{DB: dbType, URL: value, Source: item.Name}, true
	}
	return ResolvedConnection{}, false
}

func parseDBType(value string) (db.Type, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "sqlite":
		return db.TypeSQLite, nil
	case "postgres", "postgresql":
		return db.TypePostgres, nil
	case "mysql":
		return db.TypeMySQL, nil
	default:
		return "", &Error{Code: "UNSUPPORTED_DB", Detail: value}
	}
}

func inferDBType(value string) db.Type {
	lower := strings.ToLower(strings.TrimSpace(value))
	switch {
	case strings.HasPrefix(lower, "postgres://"), strings.HasPrefix(lower, "postgresql://"):
		return db.TypePostgres
	case strings.HasPrefix(lower, "mysql://"):
		return db.TypeMySQL
	case strings.HasPrefix(lower, "sqlite://"), strings.HasSuffix(lower, ".db"), strings.HasSuffix(lower, ".sqlite"), strings.HasSuffix(lower, ".sqlite3"):
		return db.TypeSQLite
	default:
		return ""
	}
}
