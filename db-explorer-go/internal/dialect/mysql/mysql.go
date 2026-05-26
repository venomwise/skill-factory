package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/introspect"
)

var ident = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type Adapter struct{}

func New() Adapter { return Adapter{} }

func (Adapter) Test(ctx context.Context, conn *sql.DB) error {
	_, err := db.QueryRows(ctx, conn, "SELECT 1")
	return err
}

func (Adapter) Schemas(ctx context.Context, conn *sql.DB) ([]introspect.Schema, error) {
	r, err := db.QueryRows(ctx, conn, `SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('information_schema','mysql','performance_schema','sys') ORDER BY schema_name`)
	if err != nil {
		return nil, err
	}
	out := make([]introspect.Schema, 0, len(r.Rows))
	for _, row := range r.Rows {
		out = append(out, introspect.Schema{Name: fmt.Sprint(row[0])})
	}
	return out, nil
}

func (Adapter) Tables(ctx context.Context, conn *sql.DB) ([]introspect.Relation, error) {
	return relations(ctx, conn, "BASE TABLE")
}

func (Adapter) Views(ctx context.Context, conn *sql.DB) ([]introspect.Relation, error) {
	return relations(ctx, conn, "VIEW")
}

func (Adapter) Schema(ctx context.Context, conn *sql.DB, name string) (introspect.RelationSchema, error) {
	schema, table, err := splitName(name)
	if err != nil {
		return introspect.RelationSchema{}, err
	}
	typeName, err := relationType(ctx, conn, schema, table)
	if err != nil {
		return introspect.RelationSchema{}, err
	}
	cols, err := columns(ctx, conn, schema, table)
	if err != nil {
		return introspect.RelationSchema{}, err
	}
	if len(cols) == 0 {
		return introspect.RelationSchema{}, &db.Error{Code: "TABLE_NOT_FOUND", Detail: name}
	}
	idx, err := indexes(ctx, conn, schema, table)
	if err != nil {
		return introspect.RelationSchema{}, err
	}
	fks, err := foreignKeys(ctx, conn, schema, table)
	if err != nil {
		return introspect.RelationSchema{}, err
	}
	return introspect.RelationSchema{
		Table:       introspect.Relation{Schema: schema, Name: table, Type: typeName},
		Columns:     cols,
		Indexes:     idx,
		ForeignKeys: fks,
	}, nil
}

func (Adapter) Data(ctx context.Context, conn *sql.DB, name string, limit int) (introspect.QueryResult, error) {
	schema, table, err := splitName(name)
	if err != nil {
		return introspect.QueryResult{}, err
	}
	if limit <= 0 {
		limit = 10
	}
	qualified := quote(table)
	if schema != "" {
		qualified = quote(schema) + "." + quote(table)
	}
	return db.QueryRows(ctx, conn, fmt.Sprintf("SELECT * FROM %s LIMIT %d", qualified, limit))
}

func (Adapter) Query(ctx context.Context, conn *sql.DB, sql string) (introspect.QueryResult, error) {
	return db.QueryRows(ctx, conn, sql)
}

func relations(ctx context.Context, conn *sql.DB, tableType string) ([]introspect.Relation, error) {
	r, err := db.QueryRows(ctx, conn, `
SELECT table_schema, table_name, table_type, table_rows
FROM information_schema.tables
WHERE table_schema = DATABASE() AND table_type = ?
ORDER BY table_name`, tableType)
	if err != nil {
		return nil, err
	}
	out := make([]introspect.Relation, 0, len(r.Rows))
	for _, row := range r.Rows {
		out = append(out, introspect.Relation{
			Schema:          fmt.Sprint(row[0]),
			Name:            fmt.Sprint(row[1]),
			Type:            strings.ToLower(fmt.Sprint(row[2])),
			RowEstimate:     toInt64Ptr(row[3]),
			RowEstimateKind: "estimate",
		})
	}
	return out, nil
}

func relationType(ctx context.Context, conn *sql.DB, schema, table string) (string, error) {
	query, args := schemaFilter(`SELECT table_type FROM information_schema.tables WHERE table_schema = %s AND table_name = ?`, schema, table)
	row, ok, err := db.QuerySingle(ctx, conn, query, args...)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", &db.Error{Code: "TABLE_NOT_FOUND", Detail: table}
	}
	return strings.ToLower(fmt.Sprint(row[0])), nil
}

func columns(ctx context.Context, conn *sql.DB, schema, table string) ([]introspect.Column, error) {
	query, args := schemaFilter(`
SELECT c.column_name, c.column_type, c.is_nullable, c.column_default,
       CASE WHEN kcu.constraint_name = 'PRIMARY' THEN 1 ELSE 0 END AS primary_key
FROM information_schema.columns c
LEFT JOIN information_schema.key_column_usage kcu
  ON kcu.table_schema = c.table_schema AND kcu.table_name = c.table_name AND kcu.column_name = c.column_name AND kcu.constraint_name = 'PRIMARY'
WHERE c.table_schema = %s AND c.table_name = ?
ORDER BY c.ordinal_position`, schema, table)
	r, err := db.QueryRows(ctx, conn, query, args...)
	if err != nil {
		return nil, err
	}
	out := make([]introspect.Column, 0, len(r.Rows))
	for _, row := range r.Rows {
		out = append(out, introspect.Column{Name: fmt.Sprint(row[0]), Type: fmt.Sprint(row[1]), Nullable: fmt.Sprint(row[2]) == "YES", Default: strPtr(row[3]), PrimaryKey: truthy(row[4])})
	}
	return out, nil
}

func indexes(ctx context.Context, conn *sql.DB, schema, table string) ([]introspect.Index, error) {
	query, args := schemaFilter(`
SELECT index_name, non_unique, GROUP_CONCAT(column_name ORDER BY seq_in_index SEPARATOR ',')
FROM information_schema.statistics
WHERE table_schema = %s AND table_name = ?
GROUP BY index_name, non_unique
ORDER BY index_name`, schema, table)
	r, err := db.QueryRows(ctx, conn, query, args...)
	if err != nil {
		return nil, err
	}
	out := make([]introspect.Index, 0, len(r.Rows))
	for _, row := range r.Rows {
		out = append(out, introspect.Index{Name: fmt.Sprint(row[0]), Unique: !truthy(row[1]), Columns: splitCSV(fmt.Sprint(row[2]))})
	}
	return out, nil
}

func foreignKeys(ctx context.Context, conn *sql.DB, schema, table string) ([]introspect.ForeignKey, error) {
	query, args := schemaFilter(`
SELECT constraint_name, column_name, referenced_table_schema, referenced_table_name, referenced_column_name
FROM information_schema.key_column_usage
WHERE table_schema = %s AND table_name = ? AND referenced_table_name IS NOT NULL
ORDER BY constraint_name, ordinal_position`, schema, table)
	r, err := db.QueryRows(ctx, conn, query, args...)
	if err != nil {
		return nil, err
	}
	out := make([]introspect.ForeignKey, 0, len(r.Rows))
	for _, row := range r.Rows {
		out = append(out, introspect.ForeignKey{Name: fmt.Sprint(row[0]), Columns: []string{fmt.Sprint(row[1])}, ReferencedSchema: fmt.Sprint(row[2]), ReferencedTable: fmt.Sprint(row[3]), ReferencedColumns: []string{fmt.Sprint(row[4])}})
	}
	return out, nil
}

func schemaFilter(template, schema, table string) (string, []any) {
	if schema == "" {
		return fmt.Sprintf(template, "DATABASE()"), []any{table}
	}
	return fmt.Sprintf(template, "?"), []any{schema, table}
}

func splitName(name string) (string, string, error) {
	parts := strings.Split(name, ".")
	if len(parts) == 1 {
		if !ident.MatchString(parts[0]) {
			return "", "", &db.Error{Code: "QUERY_FAILED", Detail: "invalid identifier: " + name}
		}
		return "", parts[0], nil
	}
	if len(parts) == 2 && ident.MatchString(parts[0]) && ident.MatchString(parts[1]) {
		return parts[0], parts[1], nil
	}
	return "", "", &db.Error{Code: "QUERY_FAILED", Detail: "invalid identifier: " + name}
}

func quote(s string) string { return "`" + s + "`" }
func strPtr(v any) *string {
	if v == nil {
		return nil
	}
	s := fmt.Sprint(v)
	return &s
}
func truthy(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case int64:
		return x != 0
	case int:
		return x != 0
	case []byte:
		return string(x) != "0" && string(x) != ""
	case string:
		return x != "" && x != "0" && x != "false"
	}
	return fmt.Sprint(v) != "" && fmt.Sprint(v) != "0"
}
func toInt64Ptr(v any) *int64 {
	if v == nil {
		return nil
	}
	switch x := v.(type) {
	case int64:
		return &x
	case int:
		y := int64(x)
		return &y
	case float64:
		y := int64(x)
		return &y
	}
	return nil
}
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
