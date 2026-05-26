package dialect

import (
	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	mysqldialect "github.com/venomwise/skill-factory/db-explorer/internal/dialect/mysql"
	postgresdialect "github.com/venomwise/skill-factory/db-explorer/internal/dialect/postgres"
	sqlitedialect "github.com/venomwise/skill-factory/db-explorer/internal/dialect/sqlite"
)

// For returns the adapter for a supported database type.
func For(dbType db.Type) (Adapter, error) {
	switch dbType {
	case db.TypeSQLite:
		return sqlitedialect.New(), nil
	case db.TypePostgres:
		return postgresdialect.New(), nil
	case db.TypeMySQL:
		return mysqldialect.New(), nil
	default:
		return nil, &db.Error{Code: "UNSUPPORTED_DB", Detail: string(dbType)}
	}
}
