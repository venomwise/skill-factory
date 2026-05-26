package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/db-explorer/internal/config"
	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/dialect"
	"github.com/venomwise/skill-factory/db-explorer/internal/safety"
)

func newQueryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "query <sql>",
		Short: "Run a read-only SQL query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sql := args[0]
			if err := safety.ValidateReadOnly(sql); err != nil {
				return writeValidationError(cmd, "query", err)
			}
			return runWithConnection(cmd, "query", func(ctx context.Context, conn *db.Connection, resolved config.ResolvedConnection) (commandResult, error) {
				adapter, err := dialect.For(resolved.DB)
				if err != nil {
					return commandResult{}, err
				}
				result, err := adapter.Query(ctx, conn.DB, sql)
				if err != nil {
					return commandResult{}, err
				}
				return commandResult{Data: result}, nil
			})
		},
	}
}
