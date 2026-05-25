package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/venomwise/skill-factory/grok-search/internal/config"
)

func sampleResult() Result {
	return Result{
		OK:            true,
		Mode:          "news",
		Query:         "query",
		ProfileID:     "main",
		ProfileSource: "config.profiles",
		Attempts:      []Attempt{{ProfileID: "main", OK: true}},
		ConfigPath:    "/tmp/grok.toml",
		ConfigPaths:   []string{"/tmp/grok.toml"},
		BaseURL:       "https://api.x.ai",
		Model:         "grok-4.1-fast",
		Content:       "answer",
		Sources:       []Source{{URL: "https://example.com", Title: "Example", Snippet: "Snippet"}},
		Raw:           "",
		Usage:         map[string]any{"total_tokens": 1},
		ElapsedMS:     12,
	}
}

func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderJSON(&buf, sampleResult()); err != nil {
		t.Fatalf("RenderJSON() error = %v", err)
	}
	text := buf.String()
	if !strings.Contains(text, "\n  \"ok\": true") {
		t.Fatalf("expected pretty JSON, got %q", text)
	}
	var decoded map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, field := range []string{"ok", "mode", "query", "profileId", "profileSource", "attempts", "config_path", "config_paths", "base_url", "model", "content", "sources", "raw", "usage", "elapsed_ms"} {
		if _, ok := decoded[field]; !ok {
			t.Fatalf("missing JSON field %q in %+v", field, decoded)
		}
	}
}

func TestRenderPlain(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderPlain(&buf, sampleResult()); err != nil {
		t.Fatalf("RenderPlain() error = %v", err)
	}
	text := buf.String()
	for _, want := range []string{"Profile: main [config.profiles]", "Attempts:", "answer", "Sources:", "https://example.com"} {
		if !strings.Contains(text, want) {
			t.Fatalf("plain output missing %q in %q", want, text)
		}
	}
}

func TestRenderPlainError(t *testing.T) {
	var buf bytes.Buffer
	resp := ErrorResponse{OK: false, Error: "missing_api_key", Detail: "configure key", Attempts: []Attempt{{ProfileID: "main", OK: false, Detail: "missing"}}}
	if err := RenderPlainError(&buf, resp); err != nil {
		t.Fatalf("RenderPlainError() error = %v", err)
	}
	text := buf.String()
	for _, want := range []string{"ERROR: missing_api_key", "configure key", "Attempts:", "main"} {
		if !strings.Contains(text, want) {
			t.Fatalf("plain error output missing %q in %q", want, text)
		}
	}
}

func TestRenderURLs(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderURLs(&buf, []Source{{URL: "https://a.example"}, {}, {URL: "https://b.example"}}); err != nil {
		t.Fatalf("RenderURLs() error = %v", err)
	}
	if got := buf.String(); got != "https://a.example\nhttps://b.example\n" {
		t.Fatalf("URLs output = %q", got)
	}
}

func TestFromError(t *testing.T) {
	cfgErr := &config.Error{Code: "invalid_config", Detail: "bad toml"}
	resp := FromError(cfgErr)
	if resp.OK || resp.Error != "invalid_config" || resp.Detail != "bad toml" {
		t.Fatalf("config error response = %+v", resp)
	}

	cmdErr := NewCommandError("all_profiles_failed", "failed", errors.New("wrapped"))
	resp = FromError(cmdErr)
	if resp.OK || resp.Error != "all_profiles_failed" || resp.Detail != "failed" {
		t.Fatalf("command error response = %+v", resp)
	}
}
