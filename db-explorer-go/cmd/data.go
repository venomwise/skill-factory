package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/db-explorer/internal/config"
	"github.com/venomwise/skill-factory/db-explorer/internal/db"
	"github.com/venomwise/skill-factory/db-explorer/internal/dialect"
)

func newDataCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "data <table>",
		Short: "Sample rows from a relation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				return writeValidationError(cmd, "data", &db.Error{Code: "QUERY_FAILED", Detail: "--limit must be greater than 0"})
			}
			return runWithConnection(cmd, "data", func(ctx context.Context, conn *db.Connection, resolved config.ResolvedConnection) (commandResult, error) {
				adapter, err := dialect.For(resolved.DB)
				if err != nil {
					return commandResult{}, err
				}
				result, err := adapter.Data(ctx, conn.DB, args[0], limit)
				if err != nil {
					return commandResult{}, err
				}
				return commandResult{Data: result, Truncated: len(result.Rows) == limit}, nil
			})
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 10, "maximum rows to return")
	return cmd
}
