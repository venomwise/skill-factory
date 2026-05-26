package db

import "database/sql"

// Type is a supported database engine identifier.
type Type string

const (
	TypeSQLite   Type = "sqlite"
	TypePostgres Type = "postgres"
	TypeMySQL    Type = "mysql"
)

// Connection wraps a database handle and its resolved engine type.
type Connection struct {
	Type Type
	DB   *sql.DB
}
