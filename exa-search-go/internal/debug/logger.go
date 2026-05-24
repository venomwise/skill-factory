package debug

import (
	"fmt"
	"os"
)

var enabled bool

// Enable turns on debug logging
func Enable() {
	enabled = true
}

// IsEnabled returns whether debug mode is active
func IsEnabled() bool {
	return enabled
}

// Log prints a debug message to stderr if debug mode is enabled
func Log(format string, args ...interface{}) {
	if !enabled {
		return
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
}

// RedactAPIKey shows only the first 8 characters of an API key
func RedactAPIKey(key string) string {
	if len(key) <= 8 {
		return key
	}
	return key[:8] + "..."
}

// LogHTTPRequest logs an HTTP request with redacted API key
func LogHTTPRequest(method, url, apiKey string) {
	if !enabled {
		return
	}
	redacted := RedactAPIKey(apiKey)
	Log("HTTP %s %s (API key: %s)", method, url, redacted)
}

// LogHTTPResponse logs an HTTP response
func LogHTTPResponse(statusCode int, bodyPreview string) {
	if !enabled {
		return
	}
	preview := bodyPreview
	if len(preview) > 200 {
		preview = preview[:200] + "..."
	}
	Log("HTTP Response: %d, body: %s", statusCode, preview)
}

// LogFailover logs a failover decision
func LogFailover(fromProfile, toProfile, reason string) {
	if !enabled {
		return
	}
	Log("Failover: %s failed (%s), trying %s", fromProfile, reason, toProfile)
}

// LogConfigResolution logs configuration loading steps
func LogConfigResolution(step string, details ...interface{}) {
	if !enabled {
		return
	}
	msg := step
	if len(details) > 0 {
		parts := make([]string, len(details))
		for i, d := range details {
			parts[i] = fmt.Sprintf("%v", d)
		}
		msg = fmt.Sprintf(step, details...)
	}
	Log("Config: %s", msg)
}
