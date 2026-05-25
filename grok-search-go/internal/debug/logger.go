package debug

import (
	"fmt"
	"os"
	"strings"
)

var enabled bool

// Enable turns on debug logging.
func Enable() {
	enabled = true
}

// Enabled reports whether debug logging is active.
func Enabled() bool {
	return enabled
}

// Log writes a formatted debug message to stderr when debug logging is enabled.
func Log(format string, args ...interface{}) {
	if !enabled {
		return
	}
	fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
}

// Redact returns a debug-safe representation of a secret value.
func Redact(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "********"
	}
	return value[:4] + "..." + value[len(value)-4:]
}
