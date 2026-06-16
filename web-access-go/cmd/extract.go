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
	extractQuery          string
	extractNum            int
	extractType           string
	extractText           bool
	extractHighlights     bool
	extractStartDate      string
	extractIncludeDomains []string
	extractExcludeDomains []string
	extractCategory       string
	extractNoAutoprompt   bool
)

// extractCmd represents the extract command (Exa provider)
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract text content from search results (Exa)",
	Long:  `Extract text content from search results using Exa API. Defaults to text extraction when neither --text nor --highlights is set.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if extractQuery == "" {
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
			return output.RenderError(getOutputFormat(), "exa", "extract", extractQuery, "", 0, nil, err)
		}

		if len(cfg.Exa.Profiles) == 0 {
			return output.RenderError(getOutputFormat(), "exa", "extract", extractQuery, "", 0, nil,
				fmt.Errorf("missing_api_key: exa"))
		}

		// Default to text extraction when neither --text nor --highlights is explicitly set
		includeText := extractText
		includeHighlights := extractHighlights
		if !cmd.Flags().Changed("text") && !cmd.Flags().Changed("highlights") {
			includeText = true
		}

		req := exa.SearchRequest{
			Query:              extractQuery,
			NumResults:         extractNum,
			Type:               extractType,
			UseAutoprompt:      !extractNoAutoprompt,
			IncludeText:        includeText,
			IncludeHighlights:  includeHighlights,
			StartPublishedDate: extractStartDate,
			IncludeDomains:     extractIncludeDomains,
			ExcludeDomains:     extractExcludeDomains,
			Category:           extractCategory,
		}

		ctx := context.Background()
		resp, attempts, err := exa.ExecuteSearch(ctx, cfg.Exa, req)
		elapsed := time.Since(startTime).Milliseconds()

		if err != nil {
			return output.RenderError(getOutputFormat(), "exa", "extract", extractQuery, "", elapsed, attempts, err)
		}

		var profileID, profileSource string
		for _, a := range attempts {
			if a.OK {
				profileID = a.ProfileID
				profileSource = a.ProfileSource
				break
			}
		}

		return output.RenderExaSuccess(getOutputFormat(), "extract", extractQuery, profileID, profileSource, resp, attempts, elapsed)
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringVar(&extractQuery, "query", "", "search query (required)")
	extractCmd.Flags().IntVar(&extractNum, "num", 5, "number of results")
	extractCmd.Flags().StringVar(&extractType, "type", "neural", "search type (neural, keyword, auto)")
	extractCmd.Flags().BoolVar(&extractText, "text", false, "include text content")
	extractCmd.Flags().BoolVar(&extractHighlights, "highlights", false, "include highlights")
	extractCmd.Flags().StringVar(&extractStartDate, "start-date", "", "start date filter (YYYY-MM-DD)")
	extractCmd.Flags().StringSliceVar(&extractIncludeDomains, "include-domains", []string{}, "domains to include")
	extractCmd.Flags().StringSliceVar(&extractExcludeDomains, "exclude-domains", []string{}, "domains to exclude")
	extractCmd.Flags().StringVar(&extractCategory, "category", "", "category filter")
	extractCmd.Flags().BoolVar(&extractNoAutoprompt, "no-autoprompt", false, "disable autoprompt")
}
