package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDefaultConfigPath verifies the default config path resolution
func TestDefaultConfigPath(t *testing.T) {
	path, err := GetDefaultConfigPath()
	if err != nil {
		t.Fatalf("failed to get default config path: %v", err)
	}

	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got: %s", path)
	}

	if !contains(path, "web-access.toml") {
		t.Errorf("expected path to contain 'web-access.toml', got: %s", path)
	}
}

// TestTemplateCreation verifies template config file creation
func TestTemplateCreation(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.toml")

	err := EnsureTemplate(configPath)
	if err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("template file was not created: %s", configPath)
	}

	// Verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	content := string(data)
	if !contains(content, "[exa]") {
		t.Error("template should contain [exa] section")
	}
	if !contains(content, "[grok]") {
		t.Error("template should contain [grok] section")
	}
	if !contains(content, "YOUR_EXA_API_KEY") {
		t.Error("template should contain placeholder for Exa API key")
	}
	if !contains(content, "YOUR_GROK_API_KEY") {
		t.Error("template should contain placeholder for Grok API key")
	}
}

// TestConfigPrecedence verifies the precedence order
func TestConfigPrecedence(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	// Create a test config with TOML profiles
	configContent := `
[exa]
base_url = "https://exa-from-toml.example"
timeout = 10

[[exa.profiles]]
id = "toml-profile"
api_key = "exa-toml-key"

[grok]
base_url = "https://grok-from-toml.example"
model = "grok-from-toml"
timeout = 20

[[grok.profiles]]
id = "grok-toml-profile"
api_key = "grok-toml-key"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	tests := []struct {
		name      string
		opts      Options
		envVars   map[string]string
		checkExa  func(*testing.T, ResolvedProviderConfig)
		checkGrok func(*testing.T, ResolvedProviderConfig)
	}{
		{
			name: "TOML config only",
			opts: Options{ConfigPath: configPath},
			checkExa: func(t *testing.T, cfg ResolvedProviderConfig) {
				if cfg.BaseURL != "https://exa-from-toml.example" {
					t.Errorf("expected base URL from TOML, got: %s", cfg.BaseURL)
				}
				if len(cfg.Profiles) != 1 || cfg.Profiles[0].APIKey != "exa-toml-key" {
					t.Errorf("expected TOML profile")
				}
			},
			checkGrok: func(t *testing.T, cfg ResolvedProviderConfig) {
				if cfg.Model != "grok-from-toml" {
					t.Errorf("expected model from TOML, got: %s", cfg.Model)
				}
			},
		},
		{
			name: "Provider-specific env overrides TOML",
			opts: Options{ConfigPath: configPath},
			envVars: map[string]string{
				"EXA_API_KEY":  "exa-env-key",
				"GROK_API_KEY": "grok-env-key",
			},
			checkExa: func(t *testing.T, cfg ResolvedProviderConfig) {
				found := false
				for _, p := range cfg.Profiles {
					if p.APIKey == "exa-env-key" {
						found = true
						break
					}
				}
				if !found {
					t.Error("expected env profile to be added")
				}
			},
			checkGrok: func(t *testing.T, cfg ResolvedProviderConfig) {
				found := false
				for _, p := range cfg.Profiles {
					if p.APIKey == "grok-env-key" {
						found = true
						break
					}
				}
				if !found {
					t.Error("expected env profile to be added")
				}
			},
		},
		{
			name: "WEB_ACCESS env overrides provider-specific env",
			opts: Options{ConfigPath: configPath},
			envVars: map[string]string{
				"EXA_API_KEY":             "exa-env-key",
				"WEB_ACCESS_EXA_API_KEY":  "exa-web-access-key",
				"GROK_API_KEY":            "grok-env-key",
				"WEB_ACCESS_GROK_API_KEY": "grok-web-access-key",
			},
			checkExa: func(t *testing.T, cfg ResolvedProviderConfig) {
				found := false
				for _, p := range cfg.Profiles {
					if p.APIKey == "exa-web-access-key" {
						found = true
						break
					}
				}
				if !found {
					t.Error("expected WEB_ACCESS env profile")
				}
			},
			checkGrok: func(t *testing.T, cfg ResolvedProviderConfig) {
				found := false
				for _, p := range cfg.Profiles {
					if p.APIKey == "grok-web-access-key" {
						found = true
						break
					}
				}
				if !found {
					t.Error("expected WEB_ACCESS env profile")
				}
			},
		},
		{
			name: "CLI flags have highest priority",
			opts: Options{
				ConfigPath: configPath,
				ExaAPIKey:  "exa-cli-key",
				GrokAPIKey: "grok-cli-key",
				Timeout:    99,
			},
			checkExa: func(t *testing.T, cfg ResolvedProviderConfig) {
				if len(cfg.Profiles) != 1 || cfg.Profiles[0].APIKey != "exa-cli-key" {
					t.Error("CLI flag should replace all other profiles")
				}
				if cfg.Timeout != 99 {
					t.Errorf("expected timeout 99, got: %d", cfg.Timeout)
				}
			},
			checkGrok: func(t *testing.T, cfg ResolvedProviderConfig) {
				if len(cfg.Profiles) != 1 || cfg.Profiles[0].APIKey != "grok-cli-key" {
					t.Error("CLI flag should replace all other profiles")
				}
				if cfg.Timeout != 99 {
					t.Errorf("expected timeout 99, got: %d", cfg.Timeout)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			cfg, err := Load(tt.opts)
			if err != nil {
				t.Fatalf("Load failed: %v", err)
			}

			if tt.checkExa != nil {
				tt.checkExa(t, cfg.Exa)
			}
			if tt.checkGrok != nil {
				tt.checkGrok(t, cfg.Grok)
			}
		})
	}
}

// TestPlaceholderFiltering verifies that placeholder API keys are filtered out
func TestPlaceholderFiltering(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	configContent := `
[[exa.profiles]]
id = "placeholder"
api_key = "YOUR_EXA_API_KEY"

[[exa.profiles]]
id = "valid"
api_key = "real-exa-key"

[[grok.profiles]]
id = "placeholder"
api_key = "YOUR_GROK_API_KEY"

[[grok.profiles]]
id = "valid"
api_key = "real-grok-key"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(Options{ConfigPath: configPath})
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check Exa profiles
	if len(cfg.Exa.Profiles) != 1 {
		t.Errorf("expected 1 Exa profile after filtering, got: %d", len(cfg.Exa.Profiles))
	}
	if len(cfg.Exa.Profiles) > 0 && cfg.Exa.Profiles[0].APIKey != "real-exa-key" {
		t.Errorf("expected valid Exa key, got: %s", cfg.Exa.Profiles[0].APIKey)
	}

	// Check Grok profiles
	if len(cfg.Grok.Profiles) != 1 {
		t.Errorf("expected 1 Grok profile after filtering, got: %d", len(cfg.Grok.Profiles))
	}
	if len(cfg.Grok.Profiles) > 0 && cfg.Grok.Profiles[0].APIKey != "real-grok-key" {
		t.Errorf("expected valid Grok key, got: %s", cfg.Grok.Profiles[0].APIKey)
	}
}

// TestProfileFiltering verifies --profile flag filters profiles correctly
func TestProfileFiltering(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	configContent := `
[[exa.profiles]]
id = "primary"
api_key = "exa-primary"

[[exa.profiles]]
id = "backup"
api_key = "exa-backup"

[[grok.profiles]]
id = "primary"
api_key = "grok-primary"

[[grok.profiles]]
id = "backup"
api_key = "grok-backup"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(Options{
		ConfigPath: configPath,
		ProfileID:  "backup",
	})
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Exa.Profiles) != 1 || cfg.Exa.Profiles[0].ID != "backup" {
		t.Error("expected only backup Exa profile")
	}
	if len(cfg.Grok.Profiles) != 1 || cfg.Grok.Profiles[0].ID != "backup" {
		t.Error("expected only backup Grok profile")
	}
}

// TestInvalidTOML verifies config parse error is returned for invalid TOML
func TestInvalidTOML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	invalidContent := `
[exa
this is not valid TOML
`
	if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := Load(Options{ConfigPath: configPath})
	if err == nil {
		t.Error("expected config parse error for invalid TOML")
	}
	if !contains(err.Error(), "config parse error") && !contains(err.Error(), "toml") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

// TestInvalidJSON verifies JSON parse error is returned for invalid extra body/headers
func TestInvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "invalid extra body JSON",
			opts: Options{ExtraBodyJSON: "{not valid json}"},
		},
		{
			name: "invalid extra headers JSON",
			opts: Options{ExtraHeadersJSON: "{not valid json}"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.toml")
			EnsureTemplate(configPath)
			tt.opts.ConfigPath = configPath

			_, err := Load(tt.opts)
			if err == nil {
				t.Error("expected invalid JSON error")
			}
			if !contains(err.Error(), "invalid JSON") && !contains(err.Error(), "json") {
				t.Errorf("expected JSON error, got: %v", err)
			}
		})
	}
}

// TestGrokExtraBodyAndHeadersMerge verifies extra body and headers are merged correctly
func TestGrokExtraBodyAndHeadersMerge(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	configContent := `
[grok]
[grok.extra_body]
temperature = 0.5

[grok.extra_headers]
X-From-TOML = "toml-value"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(Options{
		ConfigPath:       configPath,
		ExtraBodyJSON:    `{"temperature": 0.7, "max_tokens": 1000}`,
		ExtraHeadersJSON: `{"X-From-CLI": "cli-value"}`,
	})
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// CLI should override TOML for temperature
	if temp, ok := cfg.Grok.ExtraBody["temperature"].(float64); !ok || temp != 0.7 {
		t.Errorf("expected temperature 0.7, got: %v", cfg.Grok.ExtraBody["temperature"])
	}

	// max_tokens should be added
	if maxTokens, ok := cfg.Grok.ExtraBody["max_tokens"].(float64); !ok || maxTokens != 1000 {
		t.Errorf("expected max_tokens 1000, got: %v", cfg.Grok.ExtraBody["max_tokens"])
	}

	// Headers should be merged
	if cfg.Grok.ExtraHeaders["X-From-TOML"] != "toml-value" {
		t.Error("expected TOML header to be preserved")
	}
	if cfg.Grok.ExtraHeaders["X-From-CLI"] != "cli-value" {
		t.Error("expected CLI header to be added")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr)*2 && s[len(s)/2-len(substr)/2:len(s)/2+len(substr)/2+len(substr)%2] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
