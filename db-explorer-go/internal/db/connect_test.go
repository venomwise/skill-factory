package db

import (
	"context"
	"errors"
	"testing"
)

func TestOpenRejectsUnsupportedDB(t *testing.T) {
	_, err := Open(context.Background(), Type("oracle"), "example")
	assertDBCode(t, err, "UNSUPPORTED_DB")
}

func TestOpenSQLiteMemoryAndQueryRows(t *testing.T) {
	conn, err := Open(context.Background(), TypeSQLite, ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.DB.ExecContext(context.Background(), `CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)`); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.DB.ExecContext(context.Background(), `INSERT INTO users (name) VALUES ('Alice')`); err != nil {
		t.Fatal(err)
	}

	result, err := QueryRows(context.Background(), conn.DB, `SELECT id, name FROM users`)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Columns) != 2 || result.Columns[0] != "id" || result.Columns[1] != "name" {
		t.Fatalf("unexpected columns: %#v", result.Columns)
	}
	if len(result.Rows) != 1 || result.Rows[0][1] != "Alice" {
		t.Fatalf("unexpected rows: %#v", result.Rows)
	}
}

func TestQueryRowsPreservesEmptyColumns(t *testing.T) {
	conn, err := Open(context.Background(), TypeSQLite, ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	result, err := QueryRows(context.Background(), conn.DB, `SELECT 1 AS one WHERE 0`)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Columns) != 1 || result.Columns[0] != "one" {
		t.Fatalf("unexpected columns: %#v", result.Columns)
	}
	if len(result.Rows) != 0 {
		t.Fatalf("unexpected rows: %#v", result.Rows)
	}
}

func assertDBCode(t *testing.T, err error, want string) {
	t.Helper()
	var dbErr *Error
	if !errors.As(err, &dbErr) {
		t.Fatalf("expected db error %q, got %v", want, err)
	}
	if dbErr.Code != want {
		t.Fatalf("code = %q, want %q", dbErr.Code, want)
	}
}
