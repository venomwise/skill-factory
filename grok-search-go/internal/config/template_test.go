package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureTemplateCreatesConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "grok-search.toml")
	if err := EnsureTemplate(path); err != nil {
		t.Fatalf("EnsureTemplate() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(data)
	for _, want := range []string{
		`base_url = "https://api.x.ai"`,
		`model = "grok-4.1-fast"`,
		"[[profiles]]",
		"[extra_body]",
		"[extra_headers]",
		"[cooldown]",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("template missing %q", want)
		}
	}
}

func TestEnsureTemplateDoesNotOverwriteExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "grok-search.toml")
	if err := os.WriteFile(path, []byte("custom = true\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := EnsureTemplate(path); err != nil {
		t.Fatalf("EnsureTemplate() error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if got := string(data); got != "custom = true\n" {
		t.Fatalf("existing file overwritten: %q", got)
	}
}
