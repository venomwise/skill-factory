package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/introspect"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Adapter implements SQLite metadata and query behavior.
type Adapter struct{}

// New creates a SQLite adapter.
func New() Adapter { return Adapter{} }

func (Adapter) Test(ctx context.Context, conn *sql.DB) error {
	_, err := db.QueryRows(ctx, conn, "SELECT 1")
	return err
}

func (Adapter) Schemas(ctx context.Context, conn *sql.DB) ([]introspect.Schema, error) {
	result, err := db.QueryRows(ctx, conn, "PRAGMA database_list")
	if err != nil {
		return nil, err
	}
	schemas := make([]introspect.Schema, 0, len(result.Rows))
	for _, row := range result.Rows {
		if len(row) > 1 {
			schemas = append(schemas, introspect.Schema{Name: fmt.Sprint(row[1])})
		}
	}
	return schemas, nil
}

func (Adapter) Tables(ctx context.Context, conn *sql.DB) ([]introspect.Relation, error) {
	return relations(ctx, conn, "table")
}

func (Adapter) Views(ctx context.Context, conn *sql.DB) ([]introspect.Relation, error) {
	return relations(ctx, conn, "view")
}

func (Adapter) Schema(ctx context.Context, conn *sql.DB, name string) (introspect.RelationSchema, error) {
	if err := validateIdentifier(name); err != nil {
		return introspect.RelationSchema{}, err
	}
	relationType, err := relationType(ctx, conn, name)
	if err != nil {
		return introspect.RelationSchema{}, err
	}
	quoted := quoteIdentifier(name)
	columns, err := columns(ctx, conn, quoted)
	if err != nil {
		return introspect.RelationSchema{}, err
	}
	if len(columns) == 0 {
		return introspect.RelationSchema{}, &db.Error{Code: "TABLE_NOT_FOUND", Detail: name}
	}
	indexes, err := indexes(ctx, conn, quoted)
	if err != nil {
		return introspect.RelationSchema{}, err
	}
	foreignKeys, err := foreignKeys(ctx, conn, quoted)
	if err != nil {
		return introspect.RelationSchema{}, err
	}
	return introspect.RelationSchema{
		Table:       introspect.Relation{Schema: "main", Name: name, Type: relationType, RowEstimateKind: "unknown"},
		Columns:     columns,
		Indexes:     indexes,
		ForeignKeys: foreignKeys,
	}, nil
}

func (Adapter) Data(ctx context.Context, conn *sql.DB, name string, limit int) (introspect.QueryResult, error) {
	if err := validateIdentifier(name); err != nil {
		return introspect.QueryResult{}, err
	}
	if limit <= 0 {
		limit = 10
	}
	return db.QueryRows(ctx, conn, fmt.Sprintf("SELECT * FROM %s LIMIT %d", quoteIdentifier(name), limit))
}

func (Adapter) Query(ctx context.Context, conn *sql.DB, sql string) (introspect.QueryResult, error) {
	return db.QueryRows(ctx, conn, sql)
}

func relations(ctx context.Context, conn *sql.DB, relationType string) ([]introspect.Relation, error) {
	result, err := db.QueryRows(ctx, conn, `
		SELECT name, type
		FROM sqlite_master
		WHERE type = ? AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`, relationType)
	if err != nil {
		return nil, err
	}
	items := make([]introspect.Relation, 0, len(result.Rows))
	for _, row := range result.Rows {
		items = append(items, introspect.Relation{
			Schema:          "main",
			Name:            fmt.Sprint(row[0]),
			Type:            fmt.Sprint(row[1]),
			RowEstimateKind: "unknown",
		})
	}
	return items, nil
}

func relationType(ctx context.Context, conn *sql.DB, name string) (string, error) {
	row, ok, err := db.QuerySingle(ctx, conn, `
		SELECT type
		FROM sqlite_master
		WHERE name = ? AND type IN ('table', 'view')
	`, name)
	if err != nil {
		return "", err
	}
	if !ok || len(row) == 0 {
		return "", &db.Error{Code: "TABLE_NOT_FOUND", Detail: name}
	}
	return fmt.Sprint(row[0]), nil
}

func columns(ctx context.Context, conn *sql.DB, quoted string) ([]introspect.Column, error) {
	result, err := db.QueryRows(ctx, conn, "PRAGMA table_info("+quoted+")")
	if err != nil {
		return nil, err
	}
	items := make([]introspect.Column, 0, len(result.Rows))
	for _, row := range result.Rows {
		if len(row) < 6 {
			continue
		}
		items = append(items, introspect.Column{
			Name:       fmt.Sprint(row[1]),
			Type:       fmt.Sprint(row[2]),
			Nullable:   !truthy(row[3]),
			Default:    nullableString(row[4]),
			PrimaryKey: truthy(row[5]),
		})
	}
	return items, nil
}

func indexes(ctx context.Context, conn *sql.DB, quoted string) ([]introspect.Index, error) {
	result, err := db.QueryRows(ctx, conn, "PRAGMA index_list("+quoted+")")
	if err != nil {
		return nil, err
	}
	items := make([]introspect.Index, 0, len(result.Rows))
	for _, row := range result.Rows {
		if len(row) < 3 {
			continue
		}
		name := fmt.Sprint(row[1])
		cols, err := indexColumns(ctx, conn, name)
		if err != nil {
			return nil, err
		}
		items = append(items, introspect.Index{Name: name, Columns: cols, Unique: truthy(row[2])})
	}
	return items, nil
}

func indexColumns(ctx context.Context, conn *sql.DB, name string) ([]string, error) {
	result, err := db.QueryRows(ctx, conn, "PRAGMA index_info("+quoteIdentifier(name)+")")
	if err != nil {
		return nil, err
	}
	cols := make([]string, 0, len(result.Rows))
	for _, row := range result.Rows {
		if len(row) > 2 {
			cols = append(cols, fmt.Sprint(row[2]))
		}
	}
	return cols, nil
}

func foreignKeys(ctx context.Context, conn *sql.DB, quoted string) ([]introspect.ForeignKey, error) {
	result, err := db.QueryRows(ctx, conn, "PRAGMA foreign_key_list("+quoted+")")
	if err != nil {
		return nil, err
	}
	items := make([]introspect.ForeignKey, 0, len(result.Rows))
	for _, row := range result.Rows {
		if len(row) < 5 {
			continue
		}
		items = append(items, introspect.ForeignKey{
			Name:              fmt.Sprintf("fk_%v_%v", row[0], row[1]),
			Columns:           []string{fmt.Sprint(row[3])},
			ReferencedSchema:  "main",
			ReferencedTable:   fmt.Sprint(row[2]),
			ReferencedColumns: []string{fmt.Sprint(row[4])},
		})
	}
	return items, nil
}

func validateIdentifier(name string) error {
	if !identifierPattern.MatchString(name) {
		return &db.Error{Code: "QUERY_FAILED", Detail: "invalid identifier: " + name}
	}
	return nil
}

func quoteIdentifier(name string) string {
	return `"` + name + `"`
}

func nullableString(value any) *string {
	if value == nil {
		return nil
	}
	str := fmt.Sprint(value)
	return &str
}

func truthy(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case int64:
		return v != 0
	case int32:
		return v != 0
	case float64:
		return v != 0
	case string:
		return v != "" && v != "0"
	default:
		return fmt.Sprint(value) != "" && fmt.Sprint(value) != "0"
	}
}
