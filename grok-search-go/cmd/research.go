package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var researchQuery string

var researchCmd = &cobra.Command{
	Use:   "research",
	Short: "Run broad multi-source live research synthesis",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(researchQuery)
		if query == "" {
			return fmt.Errorf("--query is required")
		}
		return runResearchMode("research", query)
	},
}

func init() {
	researchCmd.Flags().StringVar(&researchQuery, "query", "", "search query / research task")
	rootCmd.AddCommand(researchCmd)
}
