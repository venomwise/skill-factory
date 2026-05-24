package output

import (
	"fmt"
	"strings"
)

// FormatErrorMessage generates user-friendly error messages with actionable guidance
func FormatErrorMessage(errorCode, detail string, attempts []Attempt) string {
	switch errorCode {
	case "missing_api_key", "no_api_key":
		return `Error: No API key configured

Provide an API key using one of these methods:
  1. Command line:  exa-search --api-key YOUR_KEY search ...
  2. Environment:   export EXA_API_KEY=YOUR_KEY
  3. Config file:   ~/.config/ai-skills/exa-search.toml

Example config file:
  [[profiles]]
  id = "main"
  api_key = "your-key-here"

Get your API key at: https://exa.ai/`

	case "rate_limit_exceeded", "http_429":
		if len(attempts) > 1 {
			// Multiple profiles failed
			var failedProfiles []string
			for _, a := range attempts {
				if !a.OK {
					reason := "unknown"
					if a.Detail != "" {
						reason = truncate(a.Detail, 50)
					}
					failedProfiles = append(failedProfiles, fmt.Sprintf("  - %s: %s", a.ProfileID, reason))
				}
			}
			return fmt.Sprintf(`Error: All API keys exhausted

Tried %d profiles, all failed:
%s

Check your usage at: https://exa.ai/dashboard`, len(attempts), strings.Join(failedProfiles, "\n"))
		}
		// Single profile
		return `Error: Rate limit exceeded

Your API key has reached its rate limit or quota.
Check your usage at: https://exa.ai/dashboard

Consider:
  - Waiting before retrying
  - Upgrading your plan
  - Adding a backup API key in config`

	case "network_timeout", "request_failed":
		return `Error: Request failed

Could not connect to Exa API (timeout or network error).
Check your network connection and try again.

If the problem persists, check Exa status: https://status.exa.ai/`

	case "invalid_config", "config_parse_error":
		return fmt.Sprintf(`Error: Invalid configuration

%s

Fix the syntax error or delete the file to use defaults.`, detail)

	case "http_401", "http_403":
		return `Error: Authentication failed

Your API key is invalid or has been revoked.
Check your API key at: https://exa.ai/dashboard

Update your configuration:
  - Config file: ~/.config/ai-skills/exa-search.toml
  - Environment: export EXA_API_KEY=YOUR_KEY
  - Command line: --api-key YOUR_KEY`

	default:
		// Generic error with detail
		if detail != "" {
			return fmt.Sprintf("Error: %s\n\n%s", errorCode, detail)
		}
		return fmt.Sprintf("Error: %s", errorCode)
	}
}
