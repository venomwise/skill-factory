package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var socialQuery string

var socialCmd = &cobra.Command{
	Use:   "social",
	Short: "Research current social and community discourse",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(socialQuery)
		if query == "" {
			return fmt.Errorf("--query is required")
		}
		return runResearchMode("social", query)
	},
}

func init() {
	socialCmd.Flags().StringVar(&socialQuery, "query", "", "search query / research task")
	rootCmd.AddCommand(socialCmd)
}
