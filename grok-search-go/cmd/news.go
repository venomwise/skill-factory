package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var newsQuery string

var newsCmd = &cobra.Command{
	Use:   "news",
	Short: "Research fresh updates and breaking news",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(newsQuery)
		if query == "" {
			return fmt.Errorf("--query is required")
		}
		return runResearchMode("news", query)
	},
}

func init() {
	newsCmd.Flags().StringVar(&newsQuery, "query", "", "search query / research task")
	rootCmd.AddCommand(newsCmd)
}
