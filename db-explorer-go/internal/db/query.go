package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/venomwise/skill-factory/db-explorer/internal/introspect"
)

// QueryRows executes a query and returns column names plus row values.
func QueryRows(ctx context.Context, conn *sql.DB, query string, args ...any) (introspect.QueryResult, error) {
	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return introspect.QueryResult{}, mapQueryError(ctx, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return introspect.QueryResult{}, &Error{Code: "QUERY_FAILED", Detail: err.Error(), Err: err}
	}

	result := introspect.QueryResult{Columns: columns, Rows: [][]interface{}{}}
	for rows.Next() {
		values := make([]any, len(columns))
		scans := make([]any, len(columns))
		for i := range values {
			scans[i] = &values[i]
		}
		if err := rows.Scan(scans...); err != nil {
			return introspect.QueryResult{}, &Error{Code: "QUERY_FAILED", Detail: err.Error(), Err: err}
		}
		for i, value := range values {
			values[i] = normalizeValue(value)
		}
		result.Rows = append(result.Rows, values)
	}
	if err := rows.Err(); err != nil {
		return introspect.QueryResult{}, mapQueryError(ctx, err)
	}
	return result, nil
}

// QuerySingle executes a query and returns the first row if present.
func QuerySingle(ctx context.Context, conn *sql.DB, query string, args ...any) ([]interface{}, bool, error) {
	result, err := QueryRows(ctx, conn, query, args...)
	if err != nil {
		return nil, false, err
	}
	if len(result.Rows) == 0 {
		return nil, false, nil
	}
	return result.Rows[0], true, nil
}

func mapQueryError(ctx context.Context, err error) error {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return &Error{Code: "QUERY_TIMEOUT", Detail: err.Error(), Err: err}
	}
	return &Error{Code: "QUERY_FAILED", Detail: err.Error(), Err: err}
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case []byte:
		return string(typed)
	default:
		return value
	}
}
