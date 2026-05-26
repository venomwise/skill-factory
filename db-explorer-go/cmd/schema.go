package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/db-explorer/internal/config"
	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/dialect"
)

func newSchemaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schema <table>",
		Short: "Show relation schema metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithConnection(cmd, "schema", func(ctx context.Context, conn *db.Connection, resolved config.ResolvedConnection) (commandResult, error) {
				adapter, err := dialect.For(resolved.DB)
				if err != nil {
					return commandResult{}, err
				}
				result, err := adapter.Schema(ctx, conn.DB, args[0])
				if err != nil {
					return commandResult{}, err
				}
				return commandResult{Data: result}, nil
			})
		},
	}
}
