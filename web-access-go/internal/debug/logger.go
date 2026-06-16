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

// Log prints debug messages to stderr if debug mode is enabled
func Log(format string, args ...interface{}) {
	if enabled {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}
