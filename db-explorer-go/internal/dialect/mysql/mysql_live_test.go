package mysql

import (
	"context"
	"os"
	"testing"

	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/introspect"
	"github.com/venomwise/skill-factory/db-explorer/internal/safety"
)

func TestLiveMySQLAdapter(t *testing.T) {
	url := os.Getenv("DBX_MYSQL_URL")
	if url == "" {
		t.Skip("DBX_MYSQL_URL not set")
	}
	ctx := context.Background()
	conn, err := db.Open(ctx, db.TypeMySQL, url)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	setupLiveMySQL(t, ctx, conn)
	defer cleanupLiveMySQL(ctx, conn)

	adapter := New()
	if err := adapter.Test(ctx, conn.DB); err != nil {
		t.Fatal(err)
	}
	if schemas, err := adapter.Schemas(ctx, conn.DB); err != nil || len(schemas) == 0 {
		t.Fatalf("schemas err=%v schemas=%#v", err, schemas)
	}
	if tables, err := adapter.Tables(ctx, conn.DB); err != nil || !hasRelation(tables, "dbx_live_users") {
		t.Fatalf("tables err=%v tables=%#v", err, tables)
	}
	if views, err := adapter.Views(ctx, conn.DB); err != nil || !hasRelation(views, "dbx_live_active_users") {
		t.Fatalf("views err=%v views=%#v", err, views)
	}
	schema, err := adapter.Schema(ctx, conn.DB, "dbx_live_orders")
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) == 0 || len(schema.Indexes) == 0 || len(schema.ForeignKeys) == 0 {
		t.Fatalf("incomplete schema metadata: %#v", schema)
	}
	query, err := adapter.Query(ctx, conn.DB, `SELECT name FROM dbx_live_users WHERE id = 1`)
	if err != nil || len(query.Rows) != 1 {
		t.Fatalf("query err=%v result=%#v", err, query)
	}
	if err := safety.ValidateReadOnly("DROP TABLE dbx_live_users"); err == nil {
		t.Fatal("expected safety rejection")
	}
}

func setupLiveMySQL(t *testing.T, ctx context.Context, conn *db.Connection) {
	t.Helper()
	cleanupLiveMySQL(ctx, conn)
	_, err := conn.DB.ExecContext(ctx, `
CREATE TABLE dbx_live_users (
  id INT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE
);
CREATE TABLE dbx_live_orders (
  id INT PRIMARY KEY,
  user_id INT NOT NULL,
  amount DECIMAL(10,2) NOT NULL,
  CONSTRAINT fk_dbx_live_orders_user FOREIGN KEY (user_id) REFERENCES dbx_live_users(id)
);
CREATE INDEX idx_dbx_live_orders_user ON dbx_live_orders(user_id);
CREATE VIEW dbx_live_active_users AS SELECT id, name FROM dbx_live_users;
INSERT INTO dbx_live_users (id, name, email) VALUES (1, 'Alice', 'alice@example.com');
INSERT INTO dbx_live_orders (id, user_id, amount) VALUES (1, 1, 10.50);
`)
	if err != nil {
		t.Fatal(err)
	}
}

func cleanupLiveMySQL(ctx context.Context, conn *db.Connection) {
	_, _ = conn.DB.ExecContext(ctx, `DROP VIEW IF EXISTS dbx_live_active_users`)
	_, _ = conn.DB.ExecContext(ctx, `DROP TABLE IF EXISTS dbx_live_orders`)
	_, _ = conn.DB.ExecContext(ctx, `DROP TABLE IF EXISTS dbx_live_users`)
}

func hasRelation(relations []introspect.Relation, name string) bool {
	for _, relation := range relations {
		if relation.Name == name {
			return true
		}
	}
	return false
}
