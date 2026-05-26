package config

import (
	"net/url"
	"os"
	"regexp"
	"strings"
)

var credentialPattern = regexp.MustCompile(`(?i)(password|passwd|pwd|token|api[_-]?key|secret)=([^\s&]+)`)

// ResolveURLEnv reads a database URL from an environment variable name.
func ResolveURLEnv(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", &Error{Code: "ENV_NOT_SET", Detail: "environment variable name is empty"}
	}
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return "", &Error{Code: "ENV_NOT_SET", Detail: name}
	}
	return value, nil
}

// MaskSecrets redacts credentials from strings before they are shown to users.
func MaskSecrets(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return value
	}
	masked := maskURLPassword(trimmed)
	return credentialPattern.ReplaceAllString(masked, `$1=***`)
}

func maskURLPassword(value string) string {
	parsed, err := url.Parse(value)
	if err != nil || parsed.User == nil {
		return value
	}
	username := parsed.User.Username()
	if _, ok := parsed.User.Password(); ok {
		parsed.User = url.UserPassword(username, "***")
	} else if username != "" {
		parsed.User = url.User(username)
	}
	return parsed.String()
}
