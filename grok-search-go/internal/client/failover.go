package client

import "strings"

var failoverStatuses = map[int]struct{}{
	401: {},
	403: {},
	429: {},
}

var failoverTextPatterns = []string{
	"rate limit",
	"rate_limit",
	"too many requests",
	"quota",
	"credits",
	"billing",
	"exhaust",
	"usage limit",
	"unauthorized",
	"forbidden",
	"invalid api key",
	"api key invalid",
	"no available tokens",
	"token unavailable",
}

// ShouldFailover reports whether a request failure should try the next profile.
func ShouldFailover(status int, detail string) bool {
	if _, ok := failoverStatuses[status]; ok {
		return true
	}
	lower := strings.ToLower(detail)
	for _, pattern := range failoverTextPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}
