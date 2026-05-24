package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/exa-search/internal/client"
	"github.com/venomwise/skill-factory/exa-search/internal/config"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Search official documentation",
	Long:  `Search official documentation with default domain filter (docs.openclaw.ai).`,
	RunE:  runDocs,
}

var (
	docsQuery            string
	docsNum              int
	docsType             string
	docsText             bool
	docsHighlights       bool
	docsStartDate        string
	docsIncludeDomains   string
	docsExcludeDomains   string
	docsCategory         string
	docsNoAutoprompt     bool
)

func init() {
	rootCmd.AddCommand(docsCmd)

	docsCmd.Flags().StringVar(&docsQuery, "query", "", "search query (required)")
	docsCmd.Flags().IntVar(&docsNum, "num", 5, "number of results")
	docsCmd.Flags().StringVar(&docsType, "type", "neural", "search type (neural, keyword, magic)")
	docsCmd.Flags().BoolVar(&docsText, "text", false, "include full text")
	docsCmd.Flags().BoolVar(&docsHighlights, "highlights", false, "include highlights")
	docsCmd.Flags().StringVar(&docsStartDate, "start-date", "", "filter by published date (ISO format)")
	docsCmd.Flags().StringVar(&docsIncludeDomains, "include-domains", "", "comma-separated domains to include (default: docs.openclaw.ai)")
	docsCmd.Flags().StringVar(&docsExcludeDomains, "exclude-domains", "", "comma-separated domains to exclude")
	docsCmd.Flags().StringVar(&docsCategory, "category", "", "filter by category")
	docsCmd.Flags().BoolVar(&docsNoAutoprompt, "no-autoprompt", false, "disable Exa autoprompt")

	docsCmd.MarkFlagRequired("query")
}

func runDocs(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	// Load configuration
	debugLog("Loading configuration from: %s", cfgFile)
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return outputError("config_parse_error", err.Error(), nil, startTime)
	}

	// Apply CLI flag overrides
	config.ApplyFlags(cfg, apiKey, profileID)
	debugLog("Loaded %d profiles", len(cfg.Profiles))

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return outputError("missing_api_key", "", nil, startTime)
	}

	// Apply base URL and timeout overrides
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if timeout > 0 {
		cfg.Timeout = timeout
	}

	debugLog("Using base URL: %s", cfg.BaseURL)
	debugLog("Timeout: %d seconds", cfg.Timeout)

	// Create client
	c := client.New(cfg.Profiles, cfg.BaseURL, time.Duration(cfg.Timeout)*time.Second)

	// Build search request with docs defaults
	req := client.SearchRequest{
		Query:              docsQuery,
		NumResults:         docsNum,
		Type:               docsType,
		UseAutoprompt:      !docsNoAutoprompt,
		IncludeText:        docsText,
		IncludeHighlights:  docsHighlights,
		StartPublishedDate: docsStartDate,
		Category:           docsCategory,
	}

	// Default to docs.openclaw.ai if no domains specified
	if docsIncludeDomains != "" {
		req.IncludeDomains = splitDomains(docsIncludeDomains)
	} else {
		req.IncludeDomains = []string{"docs.openclaw.ai"}
		debugLog("Using default include domains: docs.openclaw.ai")
	}

	if docsExcludeDomains != "" {
		req.ExcludeDomains = splitDomains(docsExcludeDomains)
	}

	debugLog("Executing docs search: query=%s, num=%d, type=%s", docsQuery, docsNum, docsType)

	// Execute search
	resp, attempts, err := c.Search(context.Background(), req)
	if err != nil {
		return outputError("request_failed", err.Error(), attempts, startTime)
	}

	// Output results
	return outputSuccess("docs", docsQuery, "", cfg.Profiles[0].ID, "config", resp, attempts, startTime)
}
