package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/db-explorer/internal/config"
	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/dialect"
)

func newViewsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "views",
		Short: "List database views",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithConnection(cmd, "views", func(ctx context.Context, conn *db.Connection, resolved config.ResolvedConnection) (commandResult, error) {
				adapter, err := dialect.For(resolved.DB)
				if err != nil {
					return commandResult{}, err
				}
				items, err := adapter.Views(ctx, conn.DB)
				if err != nil {
					return commandResult{}, err
				}
				return commandResult{Data: map[string]any{"views": items}}, nil
			})
		},
	}
}
