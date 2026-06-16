package grok

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/venomwise/skill-factory/web-access/internal/config"
)

var (
	// ErrNoAPIKey is returned when no valid Grok API key is configured
	ErrNoAPIKey = errors.New("no valid Grok API key configured")
)

// ExecuteResearch executes a research request with failover across profiles
func ExecuteResearch(ctx context.Context, cfg config.ResolvedProviderConfig, prompt, query string) (*Response, []Attempt, error) {
	if len(cfg.Profiles) == 0 {
		return nil, nil, ErrNoAPIKey
	}

	timeout := time.Duration(cfg.Timeout) * time.Second
	client := New(cfg.Profiles, cfg.BaseURL, cfg.Model, timeout, cfg.ExtraBody, cfg.ExtraHeaders)

	return client.DoResearch(ctx, prompt, query)
}

// ClassifyFailure determines if a failure is retryable
func ClassifyFailure(attempts []Attempt) string {
	if len(attempts) == 0 {
		return "unknown"
	}

	lastAttempt := attempts[len(attempts)-1]

	// Check if any attempt indicated failover was appropriate
	for _, a := range attempts {
		if a.Failover {
			return "retryable"
		}
	}

	// Check specific failure patterns
	switch lastAttempt.Status {
	case 401, 403:
		return "auth_failure"
	case 429:
		return "rate_limit"
	case 0:
		return "network_failure"
	default:
		if lastAttempt.Status >= 500 {
			return "server_error"
		}
		return "client_error"
	}
}

// FormatAttempts formats attempts for output
func FormatAttempts(attempts []Attempt) []map[string]interface{} {
	result := make([]map[string]interface{}, len(attempts))
	for i, a := range attempts {
		result[i] = map[string]interface{}{
			"profileId":     a.ProfileID,
			"profileSource": a.ProfileSource,
			"ok":            a.OK,
			"status":        a.Status,
			"failover":      a.Failover,
			"detail":        a.Detail,
		}
	}
	return result
}

// BuildRetryGuidance provides guidance for retry based on failure classification
func BuildRetryGuidance(classification string, attemptCount int) string {
	switch classification {
	case "auth_failure":
		return "Authentication failed. Check your Grok API key configuration."
	case "rate_limit":
		return fmt.Sprintf("Rate limit exceeded on %d profile(s). Wait before retrying or add more API keys for failover.", attemptCount)
	case "network_failure":
		return "Network request failed. Check your internet connection and API endpoint configuration."
	case "server_error":
		return "Grok API server error. The service may be temporarily unavailable."
	case "retryable":
		return fmt.Sprintf("All %d configured profile(s) failed. Add more API keys for better failover coverage.", attemptCount)
	default:
		return "Request failed. Check the error details for more information."
	}
}
