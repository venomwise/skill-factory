package cooldown

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/venomwise/skill-factory/grok-search/internal/config"
)

func testConfig() config.CooldownConfig {
	return config.CooldownConfig{
		Enabled:          true,
		DefaultMinutes:   15,
		RateLimitMinutes: 20,
		QuotaMinutes:     60,
		AuthMinutes:      360,
	}
}

func TestSecondsForFailure(t *testing.T) {
	cfg := testConfig()
	cases := []struct {
		name   string
		status int
		detail string
		want   int
	}{
		{name: "auth status", status: 401, want: 360 * 60},
		{name: "quota text", status: 500, detail: "quota exhausted", want: 60 * 60},
		{name: "rate status", status: 429, want: 20 * 60},
		{name: "default", status: 500, detail: "temporary", want: 15 * 60},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SecondsForFailure(tc.status, tc.detail, cfg); got != tc.want {
				t.Fatalf("SecondsForFailure() = %d, want %d", got, tc.want)
			}
		})
	}
	cfg.Enabled = false
	if got := SecondsForFailure(429, "rate limit", cfg); got != 0 {
		t.Fatalf("disabled cooldown seconds = %d", got)
	}
}

func TestStateOperations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "runtime", "cooldowns.json")
	state, err := LoadState(path)
	if err != nil {
		t.Fatalf("LoadState(missing) error = %v", err)
	}
	if len(state.Profiles) != 0 {
		t.Fatalf("initial state = %+v", state)
	}

	now := time.Unix(1000, 0)
	entry := Set(&state, "main", 60, "rate limit", 429, now)
	if entry.Seconds != 60 || entry.Status != 429 || entry.UntilText == "" || entry.SetAtText == "" {
		t.Fatalf("entry = %+v", entry)
	}
	if _, ok := Active(state, "main", now); !ok {
		t.Fatalf("expected active cooldown")
	}
	if err := SaveState(path, state); err != nil {
		t.Fatalf("SaveState() error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected saved file: %v", err)
	}

	loaded, err := LoadState(path)
	if err != nil {
		t.Fatalf("LoadState() error = %v", err)
	}
	if _, ok := Active(loaded, "main", now); !ok {
		t.Fatalf("expected loaded active cooldown")
	}
	if changed := PruneExpired(&loaded, now.Add(61*time.Second)); !changed {
		t.Fatalf("expected prune change")
	}
	if _, ok := Active(loaded, "main", now.Add(61*time.Second)); ok {
		t.Fatalf("expected cooldown to expire")
	}

	Set(&loaded, "main", 60, "rate limit", 429, now)
	if !Clear(&loaded, "main") {
		t.Fatalf("expected clear change")
	}
	if Clear(&loaded, "main") {
		t.Fatalf("expected no second clear change")
	}
}
