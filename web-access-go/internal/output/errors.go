package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/venomwise/skill-factory/web-access/internal/providers/exa"
)

// ErrorResult represents an error response
type ErrorResult struct {
	OK            bool                     `json:"ok"`
	Error         string                   `json:"error"`
	ErrorType     string                   `json:"errorType"`
	Provider      string                   `json:"provider,omitempty"`
	Mode          string                   `json:"mode,omitempty"`
	Query         string                   `json:"query,omitempty"`
	URL           string                   `json:"url,omitempty"`
	ElapsedMS     int64                    `json:"elapsedMS,omitempty"`
	Attempts      []map[string]interface{} `json:"attempts,omitempty"`
	RetryGuidance string                   `json:"retryGuidance,omitempty"`
}

// RenderError renders an error response
func RenderError(format, provider, mode, query, url string, elapsedMS int64, attempts []exa.Attempt, err error) error {
	errorType := classifyError(err)

	result := ErrorResult{
		OK:        false,
		Error:     err.Error(),
		ErrorType: errorType,
		Provider:  provider,
		Mode:      mode,
		Query:     query,
		URL:       url,
		ElapsedMS: elapsedMS,
	}

	// Add attempts if present
	if len(attempts) > 0 {
		result.Attempts = exa.FormatAttempts(attempts)
		classification := exa.ClassifyFailure(attempts)
		result.RetryGuidance = exa.BuildRetryGuidance(classification, len(attempts))
	}

	// Render based on format
	switch format {
	case "plain":
		return renderErrorPlain(result)
	case "urls":
		// URLs format doesn't make sense for errors, fallback to plain
		return renderErrorPlain(result)
	default:
		return renderErrorJSON(result)
	}
}

// classifyError determines the error type from the error message
func classifyError(err error) string {
	errMsg := err.Error()
	lowerMsg := strings.ToLower(errMsg)

	if strings.Contains(lowerMsg, "missing_api_key") {
		return "missing_api_key"
	}
	if strings.Contains(lowerMsg, "config parse error") || strings.Contains(lowerMsg, "config_parse_error") {
		return "config_parse_error"
	}
	if strings.Contains(lowerMsg, "request failed") || strings.Contains(lowerMsg, "network") {
		return "request_failed"
	}
	if strings.Contains(lowerMsg, "parse") || strings.Contains(lowerMsg, "json") {
		return "response_parse_error"
	}
	if strings.Contains(lowerMsg, "all") && strings.Contains(lowerMsg, "failed") {
		return "all_profiles_failed"
	}

	return "unknown_error"
}

// renderErrorPlain renders error in plain text
func renderErrorPlain(result ErrorResult) error {
	fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)
	fmt.Fprintf(os.Stderr, "Type: %s\n", result.ErrorType)

	if result.Provider != "" {
		fmt.Fprintf(os.Stderr, "Provider: %s\n", result.Provider)
	}
	if result.Mode != "" {
		fmt.Fprintf(os.Stderr, "Mode: %s\n", result.Mode)
	}
	if len(result.Attempts) > 0 {
		fmt.Fprintf(os.Stderr, "Attempts: %d\n", len(result.Attempts))
	}
	if result.RetryGuidance != "" {
		fmt.Fprintf(os.Stderr, "\nGuidance: %s\n", result.RetryGuidance)
	}

	return fmt.Errorf("%s", result.Error)
}

// renderErrorJSON renders error as JSON
func renderErrorJSON(result ErrorResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return err
	}
	return fmt.Errorf("%s", result.Error)
}
