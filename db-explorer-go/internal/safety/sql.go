package safety

import (
	"regexp"
	"strings"
)

var dangerousKeywords = regexp.MustCompile(`(?i)\b(INSERT|UPDATE|DELETE|DROP|ALTER|TRUNCATE|CREATE|GRANT|REVOKE|EXEC|MERGE|CALL|VACUUM|REINDEX|ATTACH|DETACH)\b`)
var firstTokenPattern = regexp.MustCompile(`^[A-Za-z]+`)
var pragmaNamePattern = regexp.MustCompile(`(?i)^PRAGMA\s+([A-Za-z_][A-Za-z0-9_]*)`)
var blockCommentPattern = regexp.MustCompile(`(?s)/\*.*?\*/`)
var lineCommentPattern = regexp.MustCompile(`(?m)--.*$`)

var allowedPrefixes = map[string]struct{}{
	"SELECT":   {},
	"WITH":     {},
	"SHOW":     {},
	"DESCRIBE": {},
	"DESC":     {},
	"EXPLAIN":  {},
	"PRAGMA":   {},
}

var safePragmas = map[string]struct{}{
	"table_info":        {},
	"index_list":        {},
	"index_info":        {},
	"foreign_key_list":  {},
	"table_xinfo":       {},
	"database_list":     {},
	"table_list":        {},
	"index_xinfo":       {},
	"collation_list":    {},
	"function_list":     {},
	"module_list":       {},
	"pragma_list":       {},
	"compile_options":   {},
	"foreign_key_check": {},
}

// Error is a structured SQL safety validation error.
type Error struct {
	Code   string
	Detail string
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return e.Code + ": " + e.Detail
	}
	return e.Code
}

// StripComments removes SQL comments before simple safety checks.
func StripComments(sql string) string {
	withoutBlocks := blockCommentPattern.ReplaceAllString(sql, " ")
	return strings.TrimSpace(lineCommentPattern.ReplaceAllString(withoutBlocks, " "))
}

// HasMultipleStatements reports whether sql contains more than one statement.
func HasMultipleStatements(sql string) bool {
	inSingle := false
	inDouble := false
	inBacktick := false
	runes := []rune(sql)
	for i, ch := range runes {
		var prev rune
		if i > 0 {
			prev = runes[i-1]
		}
		switch ch {
		case '\'':
			if !inDouble && !inBacktick && prev != '\\' {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle && !inBacktick && prev != '\\' {
				inDouble = !inDouble
			}
		case '`':
			if !inSingle && !inDouble && prev != '\\' {
				inBacktick = !inBacktick
			}
		case ';':
			if !inSingle && !inDouble && !inBacktick && strings.TrimSpace(string(runes[i+1:])) != "" {
				return true
			}
		}
	}
	return false
}

// ValidateReadOnly rejects SQL that is not allowed by the read-only contract.
func ValidateReadOnly(sql string) error {
	stripped := strings.TrimSpace(strings.TrimSuffix(StripComments(sql), ";"))
	if stripped == "" {
		return &Error{Code: "SQL_NOT_READONLY", Detail: "SQL statement is empty"}
	}
	if HasMultipleStatements(stripped) {
		return &Error{Code: "SQL_MULTIPLE_STATEMENTS", Detail: "only one read-only SQL statement is allowed"}
	}
	if dangerousKeywords.MatchString(stripped) {
		return &Error{Code: "SQL_NOT_READONLY", Detail: "dangerous write or DDL keyword detected"}
	}
	first := firstTokenPattern.FindString(stripped)
	upper := strings.ToUpper(first)
	if first == "" {
		return &Error{Code: "SQL_NOT_READONLY", Detail: "SQL statement must start with a supported read-only command"}
	}
	if _, ok := allowedPrefixes[upper]; !ok {
		return &Error{Code: "SQL_NOT_READONLY", Detail: "only SELECT, WITH, SHOW, DESCRIBE, EXPLAIN, or safe PRAGMA statements are allowed"}
	}
	if upper == "PRAGMA" {
		match := pragmaNamePattern.FindStringSubmatch(stripped)
		if len(match) < 2 {
			return &Error{Code: "UNSAFE_PRAGMA", Detail: "PRAGMA name is required"}
		}
		if _, ok := safePragmas[strings.ToLower(match[1])]; !ok {
			return &Error{Code: "UNSAFE_PRAGMA", Detail: "only read-only metadata PRAGMA statements are allowed"}
		}
	}
	return nil
}
