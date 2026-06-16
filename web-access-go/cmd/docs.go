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
	docsQuery          string
	docsNum            int
	docsType           string
	docsText           bool
	docsHighlights     bool
	docsStartDate      string
	docsIncludeDomains []string
	docsExcludeDomains []string
	docsCategory       string
	docsNoAutoprompt   bool
)

// docsCmd represents the docs command (Exa provider)
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Search official documentation (Exa)",
	Long:  `Search official documentation using Exa API with default domain filtering to docs.openclaw.ai`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if docsQuery == "" {
			return fmt.Errorf("--query is required")
		}

		startTime := time.Now()

		// Load configuration
		cfg, err := config.Load(config.Options{
			ConfigPath: cfgFile,
			ProfileID:  profileID,
			Timeout:    timeout,
			ExaAPIKey:  exaAPIKey,
		})
		if err != nil {
			return output.RenderError(getOutputFormat(), "exa", "docs", docsQuery, "", 0, nil, err)
		}

		// Check for Exa API key
		if len(cfg.Exa.Profiles) == 0 {
			return output.RenderError(getOutputFormat(), "exa", "docs", docsQuery, "", 0, nil,
				fmt.Errorf("missing_api_key: exa"))
		}

		// Build Exa request
		req := exa.SearchRequest{
			Query:              docsQuery,
			NumResults:         docsNum,
			Type:               docsType,
			UseAutoprompt:      !docsNoAutoprompt,
			IncludeText:        docsText,
			IncludeHighlights:  docsHighlights,
			StartPublishedDate: docsStartDate,
			IncludeDomains:     docsIncludeDomains,
			ExcludeDomains:     docsExcludeDomains,
			Category:           docsCategory,
		}

		// Execute search
		ctx := context.Background()
		resp, attempts, err := exa.ExecuteSearch(ctx, cfg.Exa, req)

		elapsed := time.Since(startTime).Milliseconds()

		if err != nil {
			return output.RenderError(getOutputFormat(), "exa", "docs", docsQuery, "", elapsed, attempts, err)
		}

		// Get profile info from first successful attempt
		var profileID, profileSource string
		for _, a := range attempts {
			if a.OK {
				profileID = a.ProfileID
				profileSource = a.ProfileSource
				break
			}
		}

		// Render output
		return output.RenderExaSuccess(getOutputFormat(), "docs", docsQuery, profileID, profileSource, resp, attempts, elapsed)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)

	docsCmd.Flags().StringVar(&docsQuery, "query", "", "search query (required)")
	docsCmd.Flags().IntVar(&docsNum, "num", 5, "number of results")
	docsCmd.Flags().StringVar(&docsType, "type", "neural", "search type (neural, keyword, auto)")
	docsCmd.Flags().BoolVar(&docsText, "text", false, "include text content")
	docsCmd.Flags().BoolVar(&docsHighlights, "highlights", false, "include highlights")
	docsCmd.Flags().StringVar(&docsStartDate, "start-date", "", "start date filter (YYYY-MM-DD)")
	docsCmd.Flags().StringSliceVar(&docsIncludeDomains, "include-domains", []string{"docs.openclaw.ai"}, "domains to include")
	docsCmd.Flags().StringSliceVar(&docsExcludeDomains, "exclude-domains", []string{}, "domains to exclude")
	docsCmd.Flags().StringVar(&docsCategory, "category", "", "category filter")
	docsCmd.Flags().BoolVar(&docsNoAutoprompt, "no-autoprompt", false, "disable autoprompt")
}
