package client

import "testing"

func TestShouldFailover(t *testing.T) {
	for _, status := range []int{401, 403, 429} {
		if !ShouldFailover(status, "") {
			t.Fatalf("expected status %d to failover", status)
		}
	}
	for _, detail := range []string{
		"rate limit exceeded",
		"quota exhausted",
		"billing credits unavailable",
		"invalid api key",
		"token unavailable",
		"too many requests",
	} {
		if !ShouldFailover(500, detail) {
			t.Fatalf("expected detail %q to failover", detail)
		}
	}
	if ShouldFailover(500, "internal server error") {
		t.Fatalf("unexpected failover for generic 500")
	}
}
