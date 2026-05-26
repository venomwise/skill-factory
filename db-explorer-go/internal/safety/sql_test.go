package safety

import "testing"

func TestValidateReadOnlyAcceptsSafeSQL(t *testing.T) {
	for _, sql := range []string{
		"SELECT 1",
		"WITH items AS (SELECT 1) SELECT * FROM items",
		"EXPLAIN SELECT * FROM users",
		"PRAGMA table_info(users)",
		"PRAGMA index_list(users)",
	} {
		if err := ValidateReadOnly(sql); err != nil {
			t.Fatalf("ValidateReadOnly(%q) returned error: %v", sql, err)
		}
	}
}

func TestValidateReadOnlyRejectsUnsafeSQL(t *testing.T) {
	cases := map[string]string{
		"SELECT 1; DROP TABLE users":  "SQL_MULTIPLE_STATEMENTS",
		"DROP TABLE users":            "SQL_NOT_READONLY",
		"UPDATE users SET id = 1":     "SQL_NOT_READONLY",
		"VACUUM":                      "SQL_NOT_READONLY",
		"INSERT INTO users VALUES(1)": "SQL_NOT_READONLY",
		"BEGIN":                       "SQL_NOT_READONLY",
		"PRAGMA journal_mode=WAL":     "UNSAFE_PRAGMA",
	}
	for sql, want := range cases {
		err := ValidateReadOnly(sql)
		if err == nil {
			t.Fatalf("ValidateReadOnly(%q) succeeded, want %s", sql, want)
		}
		got, ok := err.(*Error)
		if !ok {
			t.Fatalf("ValidateReadOnly(%q) error type = %T", sql, err)
		}
		if got.Code != want {
			t.Fatalf("ValidateReadOnly(%q) code = %s, want %s", sql, got.Code, want)
		}
	}
}

func TestHasMultipleStatementsIgnoresQuotedSemicolon(t *testing.T) {
	if HasMultipleStatements("SELECT ';' AS semi") {
		t.Fatal("quoted semicolon detected as multiple statements")
	}
	if !HasMultipleStatements("SELECT 1; SELECT 2") {
		t.Fatal("multiple statements not detected")
	}
}

func TestStripComments(t *testing.T) {
	got := StripComments("/* hidden */ SELECT 1 -- trailing")
	if got != "SELECT 1" {
		t.Fatalf("StripComments() = %q", got)
	}
}
