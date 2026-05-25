package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func clearGrokEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"GROK_CONFIG",
		"GROK_BASE_URL",
		"GROK_MODEL",
		"GROK_TIMEOUT",
		"GROK_API_KEY",
		"GROK_API_KEYS",
		"GROK_EXTRA_BODY_JSON",
		"GROK_EXTRA_HEADERS_JSON",
	} {
		t.Setenv(key, "")
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "grok-search.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func requireConfigErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	var cfgErr *Error
	if !errors.As(err, &cfgErr) {
		t.Fatalf("expected config Error, got %T: %v", err, err)
	}
	if cfgErr.Code != code {
		t.Fatalf("expected error code %q, got %q", code, cfgErr.Code)
	}
}

func TestLoadDefaultPathCreatesTemplateBeforeMissingAPIKey(t *testing.T) {
	clearGrokEnv(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	_, err := Load(Options{})
	if err == nil {
		t.Fatalf("expected missing api key error")
	}
	requireConfigErrorCode(t, err, "missing_api_key")

	path := filepath.Join(home, ".config", "ai-skills", "grok-search.toml")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected template config at %s: %v", path, err)
	}
}

func TestLoadPrecedence(t *testing.T) {
	clearGrokEnv(t)
	path := writeConfig(t, `
base_url = "https://config.example"
model = "config-model"
timeout = 10

[[profiles]]
id = "main"
api_key = "config-key"
`)
	t.Setenv("GROK_BASE_URL", "https://env.example")
	t.Setenv("GROK_MODEL", "env-model")
	t.Setenv("GROK_TIMEOUT", "20")

	cfg, err := Load(Options{
		ConfigPath: path,
		BaseURL:    "https://cli.example",
		Model:      "cli-model",
		Timeout:    30,
	})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.BaseURL != "https://cli.example" || cfg.Model != "cli-model" || cfg.Timeout != 30 {
		t.Fatalf("unexpected resolved config: %+v", cfg)
	}
}

func TestLoadBuiltInDefaults(t *testing.T) {
	clearGrokEnv(t)
	path := writeConfig(t, `
[[profiles]]
id = "main"
api_key = "config-key"
`)

	cfg, err := Load(Options{ConfigPath: path})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.BaseURL != DefaultBaseURL || cfg.Model != DefaultModel || cfg.Timeout != DefaultTimeout {
		t.Fatalf("defaults not applied: %+v", cfg)
	}
}

func TestResolveProfilesPriorityAndFiltering(t *testing.T) {
	clearGrokEnv(t)
	cfg := Config{
		BaseURL: DefaultBaseURL,
		Model:   DefaultModel,
		Profiles: []Profile{
			{ID: "main", APIKey: "config-main"},
			{ID: "backup", APIKey: "config-backup"},
		},
	}

	profiles := ResolveProfiles(cfg, Options{APIKey: "cli-key"})
	if len(profiles) != 1 || profiles[0].ID != "cli" || profiles[0].APIKey != "cli-key" || profiles[0].Source != "--api-key" {
		t.Fatalf("unexpected cli profiles: %+v", profiles)
	}

	t.Setenv("GROK_API_KEYS", "env-key-1,YOUR_GROK_API_KEY,env-key-3")
	profiles = ResolveProfiles(cfg, Options{})
	if len(profiles) != 2 || profiles[0].ID != "env-1" || profiles[1].ID != "env-3" {
		t.Fatalf("unexpected GROK_API_KEYS profiles: %+v", profiles)
	}

	profiles = ResolveProfiles(cfg, Options{ProfileID: "env-3"})
	if len(profiles) != 1 || profiles[0].APIKey != "env-key-3" {
		t.Fatalf("unexpected filtered env profiles: %+v", profiles)
	}

	t.Setenv("GROK_API_KEYS", "")
	t.Setenv("GROK_API_KEY", "single-env-key")
	profiles = ResolveProfiles(cfg, Options{})
	if len(profiles) != 1 || profiles[0].ID != "env" || profiles[0].Source != "GROK_API_KEY" {
		t.Fatalf("unexpected GROK_API_KEY profile: %+v", profiles)
	}

	t.Setenv("GROK_API_KEY", "")
	profiles = ResolveProfiles(cfg, Options{ProfileID: "backup"})
	if len(profiles) != 1 || profiles[0].ID != "backup" || profiles[0].APIKey != "config-backup" {
		t.Fatalf("unexpected config profile filter: %+v", profiles)
	}
}

func TestResolveProfilesPlaceholderMissingAndOverrides(t *testing.T) {
	clearGrokEnv(t)
	disabled := false
	cfg := Config{
		BaseURL: "https://global.example",
		Model:   "global-model",
		Profiles: []Profile{
			{ID: "placeholder", APIKey: "YOUR_GROK_API_KEY"},
			{ID: "disabled", APIKey: "disabled-key", Enabled: &disabled},
			{ID: "proxy", APIKey: "proxy-key", BaseURL: "https://proxy.example", Model: "proxy-model"},
		},
	}

	profiles := ResolveProfiles(cfg, Options{})
	if len(profiles) != 1 {
		t.Fatalf("expected one usable profile, got %+v", profiles)
	}
	profile := profiles[0]
	if profile.ID != "proxy" || profile.BaseURL != "https://proxy.example" || profile.Model != "proxy-model" {
		t.Fatalf("profile overrides not applied: %+v", profile)
	}

	profiles = ResolveProfiles(Config{Profiles: []Profile{{ID: "placeholder", APIKey: "YOUR_GROK_API_KEY"}}}, Options{})
	if len(profiles) != 0 {
		t.Fatalf("expected no profiles for placeholders, got %+v", profiles)
	}

	path := writeConfig(t, `
[[profiles]]
id = "main"
api_key = "YOUR_GROK_API_KEY"
`)
	_, err := Load(Options{ConfigPath: path})
	if err == nil {
		t.Fatalf("expected missing api key error")
	}
	requireConfigErrorCode(t, err, "missing_api_key")
}

func TestLoadInvalidTOMLAndJSONOverrides(t *testing.T) {
	clearGrokEnv(t)
	badTOML := writeConfig(t, `base_url = [`)
	_, err := Load(Options{ConfigPath: badTOML})
	if err == nil {
		t.Fatalf("expected invalid TOML error")
	}
	requireConfigErrorCode(t, err, "invalid_config")

	valid := writeConfig(t, `
[[profiles]]
id = "main"
api_key = "config-key"
`)
	_, err = Load(Options{ConfigPath: valid, ExtraBodyJSON: `[`})
	if err == nil {
		t.Fatalf("expected invalid JSON error")
	}
	requireConfigErrorCode(t, err, "invalid_json")

	t.Setenv("GROK_EXTRA_HEADERS_JSON", `{"X-Test":"env"}`)
	cfg, err := Load(Options{ConfigPath: valid, ExtraHeadersJSON: `{"X-Test":"cli","X-CLI":"yes"}`})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.ExtraHeaders["X-Test"] != "cli" || cfg.ExtraHeaders["X-CLI"] != "yes" {
		t.Fatalf("extra header overrides not merged with CLI precedence: %+v", cfg.ExtraHeaders)
	}
}
