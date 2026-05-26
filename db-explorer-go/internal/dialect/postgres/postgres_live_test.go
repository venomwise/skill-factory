package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/introspect"
	"github.com/venomwise/skill-factory/db-explorer/internal/safety"
)

func TestLivePostgresAdapter(t *testing.T) {
	url := os.Getenv("DBX_POSTGRES_URL")
	if url == "" {
		t.Skip("DBX_POSTGRES_URL not set")
	}
	ctx := context.Background()
	conn, err := db.Open(ctx, db.TypePostgres, url)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	setupLivePostgres(t, ctx, conn)
	defer conn.DB.ExecContext(ctx, `DROP SCHEMA IF EXISTS dbx_live CASCADE`)

	adapter := New()
	if err := adapter.Test(ctx, conn.DB); err != nil {
		t.Fatal(err)
	}
	if schemas, err := adapter.Schemas(ctx, conn.DB); err != nil || len(schemas) == 0 {
		t.Fatalf("schemas err=%v schemas=%#v", err, schemas)
	}
	if tables, err := adapter.Tables(ctx, conn.DB); err != nil || !hasRelation(tables, "dbx_users") {
		t.Fatalf("tables err=%v tables=%#v", err, tables)
	}
	if views, err := adapter.Views(ctx, conn.DB); err != nil || !hasRelation(views, "dbx_active_users") {
		t.Fatalf("views err=%v views=%#v", err, views)
	}
	schema, err := adapter.Schema(ctx, conn.DB, "dbx_live.dbx_orders")
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) == 0 || len(schema.Indexes) == 0 || len(schema.ForeignKeys) == 0 {
		t.Fatalf("incomplete schema metadata: %#v", schema)
	}
	query, err := adapter.Query(ctx, conn.DB, `SELECT name FROM dbx_live.dbx_users WHERE id = 1`)
	if err != nil || len(query.Rows) != 1 {
		t.Fatalf("query err=%v result=%#v", err, query)
	}
	if err := safety.ValidateReadOnly("DROP TABLE dbx_live.dbx_users"); err == nil {
		t.Fatal("expected safety rejection")
	}
}

func setupLivePostgres(t *testing.T, ctx context.Context, conn *db.Connection) {
	t.Helper()
	_, err := conn.DB.ExecContext(ctx, `
DROP SCHEMA IF EXISTS dbx_live CASCADE;
CREATE SCHEMA dbx_live;
CREATE TABLE dbx_live.dbx_users (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE
);
CREATE TABLE dbx_live.dbx_orders (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES dbx_live.dbx_users(id),
  amount NUMERIC(10,2) NOT NULL
);
CREATE INDEX idx_dbx_live_orders_user ON dbx_live.dbx_orders(user_id);
CREATE VIEW dbx_live.dbx_active_users AS SELECT id, name FROM dbx_live.dbx_users;
INSERT INTO dbx_live.dbx_users (id, name, email) VALUES (1, 'Alice', 'alice@example.com');
INSERT INTO dbx_live.dbx_orders (id, user_id, amount) VALUES (1, 1, 10.50);
`)
	if err != nil {
		t.Fatal(err)
	}
}

func hasRelation(relations []introspect.Relation, name string) bool {
	for _, relation := range relations {
		if relation.Name == name {
			return true
		}
	}
	return false
}
