package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveConnectionCLIOverridesConfig(t *testing.T) {
	t.Chdir(writeProjectConfig(t, `default_profile = "local"

[[profiles]]
id = "local"
db = "sqlite"
url = "./local.db"
`))

	conn, err := ResolveConnection(Options{DB: "postgres", URL: "postgres://user:pass@example/db"})
	if err != nil {
		t.Fatal(err)
	}
	if conn.DB != "postgres" || conn.URL != "postgres://user:pass@example/db" || conn.Source != "cli" {
		t.Fatalf("unexpected connection: %+v", conn)
	}
}

func TestResolveConnectionProjectPrecedesGlobal(t *testing.T) {
	projectDir := writeProjectConfig(t, `default_profile = "local"

[[profiles]]
id = "local"
db = "sqlite"
url = "./project.db"
`)
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeFile(t, filepath.Join(home, DefaultConfigDir, DefaultConfigFile), `default_profile = "local"

[[profiles]]
id = "local"
db = "sqlite"
url = "./global.db"
`)
	t.Chdir(projectDir)

	conn, err := ResolveConnection(Options{})
	if err != nil {
		t.Fatal(err)
	}
	if conn.URL != "./project.db" {
		t.Fatalf("url = %q, want project config", conn.URL)
	}
}

func TestResolveConnectionDefaultProfileAndURLEnv(t *testing.T) {
	t.Setenv("TEST_DB_URL", "sqlite:///tmp/test.db")
	t.Chdir(writeProjectConfig(t, `default_profile = "local"

[[profiles]]
id = "local"
db = "sqlite"
url_env = "TEST_DB_URL"
`))

	conn, err := ResolveConnection(Options{})
	if err != nil {
		t.Fatal(err)
	}
	if conn.URL != "sqlite:///tmp/test.db" || conn.URLEnv != "TEST_DB_URL" || conn.Profile != "local" {
		t.Fatalf("unexpected connection: %+v", conn)
	}
}

func TestResolveConnectionErrors(t *testing.T) {
	t.Chdir(writeProjectConfig(t, `[[profiles]]
id = "local"
db = "sqlite"
url_env = "MISSING_DB_URL"
`))

	_, err := ResolveConnection(Options{ProfileID: "missing"})
	assertConfigCode(t, err, "PROFILE_NOT_FOUND")

	_, err = ResolveConnection(Options{ProfileID: "local"})
	assertConfigCode(t, err, "ENV_NOT_SET")

	_, err = ResolveConnection(Options{DB: "oracle", URL: "x"})
	assertConfigCode(t, err, "UNSUPPORTED_DB")
}

func TestResolveConnectionEnvFallback(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@example/db")
	t.Chdir(t.TempDir())
	t.Setenv("HOME", t.TempDir())

	conn, err := ResolveConnection(Options{})
	if err != nil {
		t.Fatal(err)
	}
	if conn.DB != "postgres" || conn.Source != "DATABASE_URL" {
		t.Fatalf("unexpected connection: %+v", conn)
	}
}

func TestMaskSecrets(t *testing.T) {
	masked := MaskSecrets("postgres://user:secret@example/db?token=abc&password=def")
	if masked == "postgres://user:secret@example/db?token=abc&password=def" {
		t.Fatal("secret was not masked")
	}
	for _, secret := range []string{"secret", "token=abc", "password=def"} {
		if strings.Contains(masked, secret) {
			t.Fatalf("masked value still contains %q: %s", secret, masked)
		}
	}
}

func assertConfigCode(t *testing.T, err error, want string) {
	t.Helper()
	var cfgErr *Error
	if !errors.As(err, &cfgErr) {
		t.Fatalf("expected config error %q, got %v", want, err)
	}
	if cfgErr.Code != want {
		t.Fatalf("code = %q, want %q", cfgErr.Code, want)
	}
}

func writeProjectConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, DefaultProjectConfig), content)
	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
