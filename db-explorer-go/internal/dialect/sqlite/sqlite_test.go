package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/venomwise/skill-factory/db-explorer/internal/introspect"
	_ "modernc.org/sqlite"
)

func TestAdapterMetadataAndQueries(t *testing.T) {
	ctx := context.Background()
	conn := newFixtureDB(t)
	defer conn.Close()

	adapter := New()
	if err := adapter.Test(ctx, conn); err != nil {
		t.Fatal(err)
	}

	schemas, err := adapter.Schemas(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	if len(schemas) == 0 || schemas[0].Name != "main" {
		t.Fatalf("unexpected schemas: %#v", schemas)
	}

	tables, err := adapter.Tables(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	if !hasRelation(tables, "users") || !hasRelation(tables, "orders") {
		t.Fatalf("unexpected tables: %#v", tables)
	}
	for _, table := range tables {
		if table.RowEstimate != nil || table.RowEstimateKind != "unknown" {
			t.Fatalf("unexpected row estimate metadata: %#v", table)
		}
	}

	views, err := adapter.Views(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	if !hasRelation(views, "active_users") {
		t.Fatalf("unexpected views: %#v", views)
	}

	schema, err := adapter.Schema(ctx, conn, "orders")
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 3 {
		t.Fatalf("unexpected columns: %#v", schema.Columns)
	}
	if !hasIndex(schema.Indexes, "idx_orders_user") {
		t.Fatalf("missing index: %#v", schema.Indexes)
	}
	if len(schema.ForeignKeys) != 1 || schema.ForeignKeys[0].ReferencedTable != "users" {
		t.Fatalf("unexpected foreign keys: %#v", schema.ForeignKeys)
	}

	data, err := adapter.Data(ctx, conn, "users", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Rows) != 1 || len(data.Columns) != 3 {
		t.Fatalf("unexpected data result: %#v", data)
	}

	query, err := adapter.Query(ctx, conn, "SELECT name FROM users WHERE id = 1")
	if err != nil {
		t.Fatal(err)
	}
	if len(query.Rows) != 1 || query.Rows[0][0] != "Alice" {
		t.Fatalf("unexpected query result: %#v", query)
	}

	empty, err := adapter.Query(ctx, conn, "SELECT name FROM users WHERE id = -1")
	if err != nil {
		t.Fatal(err)
	}
	if len(empty.Columns) != 1 || len(empty.Rows) != 0 {
		t.Fatalf("empty result did not preserve columns: %#v", empty)
	}
}

func TestAdapterRejectsInvalidIdentifier(t *testing.T) {
	conn := newFixtureDB(t)
	defer conn.Close()

	_, err := New().Schema(context.Background(), conn, "users;drop")
	if err == nil {
		t.Fatal("expected invalid identifier error")
	}
}

func newFixtureDB(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec(`
		PRAGMA foreign_keys = ON;
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE
		);
		CREATE TABLE orders (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			amount REAL NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
		CREATE INDEX idx_orders_user ON orders(user_id);
		CREATE VIEW active_users AS SELECT id, name FROM users;
		INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com'), ('Bob', 'bob@example.com');
		INSERT INTO orders (user_id, amount) VALUES (1, 10.5);
	`)
	if err != nil {
		conn.Close()
		t.Fatal(err)
	}
	return conn
}

func hasRelation(relations []introspect.Relation, name string) bool {
	for _, relation := range relations {
		if relation.Name == name {
			return true
		}
	}
	return false
}

func hasIndex(indexes []introspect.Index, name string) bool {
	for _, index := range indexes {
		if index.Name == name {
			return true
		}
	}
	return false
}
