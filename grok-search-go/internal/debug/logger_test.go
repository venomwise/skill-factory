package debug

import "testing"

func TestRedact(t *testing.T) {
	if got := Redact("secret-key-1234"); got != "secr...1234" {
		t.Fatalf("Redact() = %q", got)
	}
	if got := Redact("short"); got != "********" {
		t.Fatalf("Redact(short) = %q", got)
	}
	if got := Redact("   "); got != "" {
		t.Fatalf("Redact(empty) = %q", got)
	}
}
