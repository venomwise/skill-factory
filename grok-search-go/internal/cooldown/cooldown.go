package cooldown

import (
	"strings"
	"time"

	"github.com/venomwise/skill-factory/grok-search/internal/config"
)

// SecondsForFailure returns the cooldown duration in seconds for a failure.
func SecondsForFailure(status int, detail string, cfg config.CooldownConfig) int {
	if !cfg.Enabled {
		return 0
	}
	return int(DurationForFailure(status, detail, cfg).Seconds())
}

// DurationForFailure maps failure type to configured cooldown duration.
func DurationForFailure(status int, detail string, cfg config.CooldownConfig) time.Duration {
	lower := strings.ToLower(detail)
	minutes := cfg.DefaultMinutes
	if minutes <= 0 {
		minutes = 15
	}

	if status == 401 || status == 403 || containsAny(lower, "invalid api key", "authentication_error", "unauthorized", "forbidden") {
		minutes = cfg.AuthMinutes
		if minutes <= 0 {
			minutes = 360
		}
		return time.Duration(minutes) * time.Minute
	}
	if containsAny(lower, "quota", "credits", "billing", "usage limit", "exhaust", "insufficient") {
		minutes = cfg.QuotaMinutes
		if minutes <= 0 {
			minutes = 60
		}
		return time.Duration(minutes) * time.Minute
	}
	if status == 429 || containsAny(lower, "rate limit", "rate_limit", "too many requests", "no available tokens", "token unavailable") {
		minutes = cfg.RateLimitMinutes
		if minutes <= 0 {
			minutes = 20
		}
		return time.Duration(minutes) * time.Minute
	}
	return time.Duration(minutes) * time.Minute
}

func containsAny(value string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}
