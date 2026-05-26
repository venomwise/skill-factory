package dialect

import (
	"context"
	"database/sql"

	"github.com/venomwise/skill-factory/db-explorer/internal/introspect"
)

// Adapter defines database-specific metadata and query behavior.
type Adapter interface {
	Test(ctx context.Context, conn *sql.DB) error
	Schemas(ctx context.Context, conn *sql.DB) ([]introspect.Schema, error)
	Tables(ctx context.Context, conn *sql.DB) ([]introspect.Relation, error)
	Views(ctx context.Context, conn *sql.DB) ([]introspect.Relation, error)
	Schema(ctx context.Context, conn *sql.DB, name string) (introspect.RelationSchema, error)
	Data(ctx context.Context, conn *sql.DB, name string, limit int) (introspect.QueryResult, error)
	Query(ctx context.Context, conn *sql.DB, sql string) (introspect.QueryResult, error)
}
