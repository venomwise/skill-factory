package postgres

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
	r, err := db.QueryRows(ctx, conn, `SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT LIKE 'pg_%' AND schema_name <> 'information_schema' ORDER BY schema_name`)
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
	return introspect.RelationSchema{Table: introspect.Relation{Schema: schema, Name: table, Type: typeName}, Columns: cols, Indexes: idx, ForeignKeys: fks}, nil
}

func (Adapter) Data(ctx context.Context, conn *sql.DB, name string, limit int) (introspect.QueryResult, error) {
	schema, table, err := splitName(name)
	if err != nil {
		return introspect.QueryResult{}, err
	}
	if limit <= 0 {
		limit = 10
	}
	return db.QueryRows(ctx, conn, fmt.Sprintf("SELECT * FROM %s.%s LIMIT %d", quote(schema), quote(table), limit))
}
func (Adapter) Query(ctx context.Context, conn *sql.DB, sql string) (introspect.QueryResult, error) {
	return db.QueryRows(ctx, conn, sql)
}

func relations(ctx context.Context, conn *sql.DB, tableType string) ([]introspect.Relation, error) {
	r, err := db.QueryRows(ctx, conn, `
SELECT t.table_schema, t.table_name, t.table_type, c.reltuples::bigint
FROM information_schema.tables t
LEFT JOIN pg_class c ON c.relname = t.table_name
LEFT JOIN pg_namespace n ON n.oid = c.relnamespace AND n.nspname = t.table_schema
WHERE t.table_schema NOT LIKE 'pg_%' AND t.table_schema <> 'information_schema' AND t.table_type = $1
ORDER BY t.table_schema, t.table_name`, tableType)
	if err != nil {
		return nil, err
	}
	out := make([]introspect.Relation, 0, len(r.Rows))
	for _, row := range r.Rows {
		out = append(out, introspect.Relation{Schema: fmt.Sprint(row[0]), Name: fmt.Sprint(row[1]), Type: strings.ToLower(fmt.Sprint(row[2])), RowEstimate: toInt64Ptr(row[3]), RowEstimateKind: "estimate"})
	}
	return out, nil
}

func relationType(ctx context.Context, conn *sql.DB, schema, table string) (string, error) {
	row, ok, err := db.QuerySingle(ctx, conn, `SELECT table_type FROM information_schema.tables WHERE table_schema=$1 AND table_name=$2`, schema, table)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", &db.Error{Code: "TABLE_NOT_FOUND", Detail: schema + "." + table}
	}
	return strings.ToLower(fmt.Sprint(row[0])), nil
}

func columns(ctx context.Context, conn *sql.DB, schema, table string) ([]introspect.Column, error) {
	r, err := db.QueryRows(ctx, conn, `
SELECT c.column_name, c.data_type, c.is_nullable, c.column_default,
       EXISTS (SELECT 1 FROM information_schema.table_constraints tc JOIN information_schema.key_column_usage kcu ON tc.constraint_name=kcu.constraint_name AND tc.table_schema=kcu.table_schema WHERE tc.constraint_type='PRIMARY KEY' AND tc.table_schema=c.table_schema AND tc.table_name=c.table_name AND kcu.column_name=c.column_name)
FROM information_schema.columns c WHERE c.table_schema=$1 AND c.table_name=$2 ORDER BY c.ordinal_position`, schema, table)
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
	r, err := db.QueryRows(ctx, conn, `
SELECT i.relname, ix.indisunique, array_to_string(array_agg(a.attname ORDER BY x.ord), ',')
FROM pg_class t JOIN pg_namespace ns ON ns.oid=t.relnamespace JOIN pg_index ix ON t.oid=ix.indrelid JOIN pg_class i ON i.oid=ix.indexrelid JOIN unnest(ix.indkey) WITH ORDINALITY AS x(attnum, ord) ON true JOIN pg_attribute a ON a.attrelid=t.oid AND a.attnum=x.attnum
WHERE ns.nspname=$1 AND t.relname=$2 GROUP BY i.relname, ix.indisunique ORDER BY i.relname`, schema, table)
	if err != nil {
		return nil, err
	}
	out := make([]introspect.Index, 0, len(r.Rows))
	for _, row := range r.Rows {
		out = append(out, introspect.Index{Name: fmt.Sprint(row[0]), Unique: truthy(row[1]), Columns: splitCSV(fmt.Sprint(row[2]))})
	}
	return out, nil
}

func foreignKeys(ctx context.Context, conn *sql.DB, schema, table string) ([]introspect.ForeignKey, error) {
	r, err := db.QueryRows(ctx, conn, `
SELECT tc.constraint_name, kcu.column_name, ccu.table_schema, ccu.table_name, ccu.column_name
FROM information_schema.table_constraints tc JOIN information_schema.key_column_usage kcu ON tc.constraint_name=kcu.constraint_name AND tc.table_schema=kcu.table_schema JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name=tc.constraint_name AND ccu.table_schema=tc.table_schema
WHERE tc.constraint_type='FOREIGN KEY' AND tc.table_schema=$1 AND tc.table_name=$2 ORDER BY tc.constraint_name, kcu.ordinal_position`, schema, table)
	if err != nil {
		return nil, err
	}
	out := make([]introspect.ForeignKey, 0, len(r.Rows))
	for _, row := range r.Rows {
		out = append(out, introspect.ForeignKey{Name: fmt.Sprint(row[0]), Columns: []string{fmt.Sprint(row[1])}, ReferencedSchema: fmt.Sprint(row[2]), ReferencedTable: fmt.Sprint(row[3]), ReferencedColumns: []string{fmt.Sprint(row[4])}})
	}
	return out, nil
}

func splitName(name string) (string, string, error) {
	parts := strings.Split(name, ".")
	if len(parts) == 1 {
		if !ident.MatchString(parts[0]) {
			return "", "", &db.Error{Code: "QUERY_FAILED", Detail: "invalid identifier: " + name}
		}
		return "public", parts[0], nil
	}
	if len(parts) == 2 && ident.MatchString(parts[0]) && ident.MatchString(parts[1]) {
		return parts[0], parts[1], nil
	}
	return "", "", &db.Error{Code: "QUERY_FAILED", Detail: "invalid identifier: " + name}
}
func quote(s string) string { return `"` + s + `"` }
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
