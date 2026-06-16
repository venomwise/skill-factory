package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/web-access/internal/config"
	"github.com/venomwise/skill-factory/web-access/internal/output"
	"github.com/venomwise/skill-factory/web-access/internal/providers/exa"
)

var (
	searchQuery          string
	searchNum            int
	searchType           string
	searchText           bool
	searchHighlights     bool
	searchStartDate      string
	searchIncludeDomains []string
	searchExcludeDomains []string
	searchCategory       string
	searchNoAutoprompt   bool
)

// searchCmd represents the search command (Exa provider)
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Neural web search (Exa)",
	Long:  `Perform neural web search using Exa API without default domain filtering`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if searchQuery == "" {
			return fmt.Errorf("--query is required")
		}

		startTime := time.Now()

		cfg, err := config.Load(config.Options{
			ConfigPath: cfgFile,
			ProfileID:  profileID,
			Timeout:    timeout,
			ExaAPIKey:  exaAPIKey,
		})
		if err != nil {
			return output.RenderError(getOutputFormat(), "exa", "search", searchQuery, "", 0, nil, err)
		}

		if len(cfg.Exa.Profiles) == 0 {
			return output.RenderError(getOutputFormat(), "exa", "search", searchQuery, "", 0, nil,
				fmt.Errorf("missing_api_key: exa"))
		}

		req := exa.SearchRequest{
			Query:              searchQuery,
			NumResults:         searchNum,
			Type:               searchType,
			UseAutoprompt:      !searchNoAutoprompt,
			IncludeText:        searchText,
			IncludeHighlights:  searchHighlights,
			StartPublishedDate: searchStartDate,
			IncludeDomains:     searchIncludeDomains,
			ExcludeDomains:     searchExcludeDomains,
			Category:           searchCategory,
		}

		ctx := context.Background()
		resp, attempts, err := exa.ExecuteSearch(ctx, cfg.Exa, req)
		elapsed := time.Since(startTime).Milliseconds()

		if err != nil {
			return output.RenderError(getOutputFormat(), "exa", "search", searchQuery, "", elapsed, attempts, err)
		}

		var profileID, profileSource string
		for _, a := range attempts {
			if a.OK {
				profileID = a.ProfileID
				profileSource = a.ProfileSource
				break
			}
		}

		return output.RenderExaSuccess(getOutputFormat(), "search", searchQuery, profileID, profileSource, resp, attempts, elapsed)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchQuery, "query", "", "search query (required)")
	searchCmd.Flags().IntVar(&searchNum, "num", 5, "number of results")
	searchCmd.Flags().StringVar(&searchType, "type", "neural", "search type (neural, keyword, auto)")
	searchCmd.Flags().BoolVar(&searchText, "text", false, "include text content")
	searchCmd.Flags().BoolVar(&searchHighlights, "highlights", false, "include highlights")
	searchCmd.Flags().StringVar(&searchStartDate, "start-date", "", "start date filter (YYYY-MM-DD)")
	searchCmd.Flags().StringSliceVar(&searchIncludeDomains, "include-domains", []string{}, "domains to include")
	searchCmd.Flags().StringSliceVar(&searchExcludeDomains, "exclude-domains", []string{}, "domains to exclude")
	searchCmd.Flags().StringVar(&searchCategory, "category", "", "category filter")
	searchCmd.Flags().BoolVar(&searchNoAutoprompt, "no-autoprompt", false, "disable autoprompt")
}
