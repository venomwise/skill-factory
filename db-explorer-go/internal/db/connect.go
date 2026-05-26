package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// Error is a structured database error.
type Error struct {
	Code   string
	Detail string
	Err    error
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return e.Code + ": " + e.Detail
	}
	return e.Code
}

func (e *Error) Unwrap() error { return e.Err }

// Open creates and verifies a database connection.
func Open(ctx context.Context, dbType Type, rawURL string) (*Connection, error) {
	driver, dsn, err := driverAndDSN(dbType, rawURL)
	if err != nil {
		return nil, err
	}
	handle, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, &Error{Code: "CONNECTION_FAILED", Detail: err.Error(), Err: err}
	}
	if err := handle.PingContext(ctx); err != nil {
		_ = handle.Close()
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, &Error{Code: "QUERY_TIMEOUT", Detail: err.Error(), Err: err}
		}
		return nil, &Error{Code: "CONNECTION_FAILED", Detail: err.Error(), Err: err}
	}
	return &Connection{Type: dbType, DB: handle}, nil
}

// Close releases the underlying database handle.
func (c *Connection) Close() error {
	if c == nil || c.DB == nil {
		return nil
	}
	return c.DB.Close()
}

func driverAndDSN(dbType Type, rawURL string) (string, string, error) {
	switch dbType {
	case TypeSQLite:
		return "sqlite", sqliteDSN(rawURL), nil
	case TypePostgres:
		return "pgx", rawURL, nil
	case TypeMySQL:
		dsn, err := mysqlDSN(rawURL)
		if err != nil {
			return "", "", &Error{Code: "CONNECTION_FAILED", Detail: err.Error(), Err: err}
		}
		return "mysql", dsn, nil
	default:
		return "", "", &Error{Code: "UNSUPPORTED_DB", Detail: string(dbType)}
	}
}

func sqliteDSN(rawURL string) string {
	value := strings.TrimSpace(rawURL)
	if strings.HasPrefix(value, "sqlite:///") {
		return strings.TrimPrefix(value, "sqlite://")
	}
	if strings.HasPrefix(value, "sqlite://") {
		return strings.TrimPrefix(value, "sqlite://")
	}
	return value
}

func mysqlDSN(rawURL string) (string, error) {
	value := strings.TrimSpace(rawURL)
	if value == "" || !strings.Contains(value, "://") {
		return value, nil
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "mysql" {
		return "", fmt.Errorf("unsupported MySQL URL scheme: %s", parsed.Scheme)
	}
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = parsed.Host
	cfg.User = parsed.User.Username()
	if password, ok := parsed.User.Password(); ok {
		cfg.Passwd = password
	}
	cfg.DBName = strings.TrimPrefix(parsed.Path, "/")
	cfg.ParseTime = true
	params := parsed.Query()
	for key, values := range params {
		if len(values) > 0 {
			cfg.Params[key] = values[len(values)-1]
		}
	}
	return cfg.FormatDSN(), nil
}
