package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var docsCompareQuery string

var docsCompareCmd = &cobra.Command{
	Use:   "docs-compare",
	Short: "Compare official facts with community interpretation",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(docsCompareQuery)
		if query == "" {
			return fmt.Errorf("--query is required")
		}
		return runResearchMode("docs-compare", query)
	},
}

func init() {
	docsCompareCmd.Flags().StringVar(&docsCompareQuery, "query", "", "search query / research task")
	rootCmd.AddCommand(docsCompareCmd)
}
