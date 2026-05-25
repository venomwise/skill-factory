package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/venomwise/skill-factory/grok-search/internal/cooldown"
)

func resetRunGlobals(t *testing.T) *bytes.Buffer {
	t.Helper()
	cfgFile = ""
	apiKey = ""
	baseURL = ""
	model = ""
	timeout = 0
	profileID = ""
	ignoreCooldown = false
	extraBodyJSON = ""
	extraHeadersJSON = ""
	plainOutput = false
	urlsOutput = false
	jsonOutput = false
	debugMode = false

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
	})
	return &out
}

func writeRunConfig(t *testing.T, serverURL, cooldownPath string, keys ...string) string {
	t.Helper()
	profiles := strings.Builder{}
	for i, key := range keys {
		profiles.WriteString("\n[[profiles]]\n")
		profiles.WriteString("id = \"")
		if i == 0 {
			profiles.WriteString("main")
		} else {
			profiles.WriteString("backup")
		}
		profiles.WriteString("\"\n")
		profiles.WriteString("api_key = \"")
		profiles.WriteString(key)
		profiles.WriteString("\"\n")
	}
	content := "base_url = \"" + serverURL + "\"\nmodel = \"grok-test\"\ntimeout = 5\n" + profiles.String() + "\n[cooldown]\nenabled = true\nstate_file = \"" + cooldownPath + "\"\ndefault_minutes = 15\nrate_limit_minutes = 20\nquota_minutes = 60\nauth_minutes = 360\n"
	path := filepath.Join(t.TempDir(), "grok-search.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func TestRunResearchModeFailoverThenSuccess(t *testing.T) {
	out := resetRunGlobals(t)
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			http.Error(w, "rate limit", http.StatusTooManyRequests)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model":   "grok-success",
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": `{"content":"ok","sources":[]}`}}},
		})
	}))
	defer server.Close()

	cooldownPath := filepath.Join(t.TempDir(), "cooldowns.json")
	cfgFile = writeRunConfig(t, server.URL, cooldownPath, "key-1", "key-2")

	if err := runResearchMode("news", "query"); err != nil {
		t.Fatalf("runResearchMode() error = %v", err)
	}
	if requestCount != 2 {
		t.Fatalf("requestCount = %d", requestCount)
	}
	var result map[string]any
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("Unmarshal output error = %v; output=%s", err, out.String())
	}
	if result["profileId"] != "backup" || result["content"] != "ok" {
		t.Fatalf("unexpected result: %+v", result)
	}
	attempts := result["attempts"].([]any)
	if len(attempts) != 2 {
		t.Fatalf("attempts = %+v", attempts)
	}
	state, err := cooldown.LoadState(cooldownPath)
	if err != nil {
		t.Fatalf("LoadState() error = %v", err)
	}
	if _, ok := cooldown.Active(state, "main", time.Now()); !ok {
		t.Fatalf("expected main profile in cooldown")
	}
}

func TestRunResearchModeAllProfilesInCooldown(t *testing.T) {
	resetRunGlobals(t)
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	cooldownPath := filepath.Join(t.TempDir(), "cooldowns.json")
	state := cooldown.State{Profiles: map[string]cooldown.Entry{}}
	cooldown.Set(&state, "main", 60, "rate limit", 429, time.Now())
	if err := cooldown.SaveState(cooldownPath, state); err != nil {
		t.Fatalf("SaveState() error = %v", err)
	}
	cfgFile = writeRunConfig(t, server.URL, cooldownPath, "key-1")

	err := runResearchMode("news", "query")
	if err == nil || !strings.Contains(err.Error(), "all_profiles_in_cooldown") {
		t.Fatalf("expected all_profiles_in_cooldown, got %v", err)
	}
}

func TestRunResearchModeAllProfilesFailed(t *testing.T) {
	resetRunGlobals(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "rate limit", http.StatusTooManyRequests)
	}))
	defer server.Close()

	cooldownPath := filepath.Join(t.TempDir(), "cooldowns.json")
	cfgFile = writeRunConfig(t, server.URL, cooldownPath, "key-1", "key-2")

	err := runResearchMode("news", "query")
	if err == nil || !strings.Contains(err.Error(), "all_profiles_failed") {
		t.Fatalf("expected all_profiles_failed, got %v", err)
	}
}
