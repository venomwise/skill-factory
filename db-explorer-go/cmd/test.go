package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/db-explorer/internal/config"
	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/dialect"
)

func newTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test the resolved database connection",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithConnection(cmd, "test", func(ctx context.Context, conn *db.Connection, resolved config.ResolvedConnection) (commandResult, error) {
				adapter, err := dialect.For(resolved.DB)
				if err != nil {
					return commandResult{}, err
				}
				if err := adapter.Test(ctx, conn.DB); err != nil {
					return commandResult{}, err
				}
				return commandResult{Data: map[string]any{"connected": true}}, nil
			})
		},
	}
}
