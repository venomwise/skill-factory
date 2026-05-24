package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/exa-search/internal/client"
	"github.com/venomwise/skill-factory/exa-search/internal/config"
)

var researchCmd = &cobra.Command{
	Use:   "research",
	Short: "Deep research with text extraction",
	Long:  `Execute a search with text extraction enabled by default for deep research.`,
	RunE:  runResearch,
}

var (
	researchQuery            string
	researchNum              int
	researchType             string
	researchText             bool
	researchHighlights       bool
	researchStartDate        string
	researchIncludeDomains   string
	researchExcludeDomains   string
	researchCategory         string
	researchNoAutoprompt     bool
)

func init() {
	rootCmd.AddCommand(researchCmd)

	researchCmd.Flags().StringVar(&researchQuery, "query", "", "search query (required)")
	researchCmd.Flags().IntVar(&researchNum, "num", 5, "number of results")
	researchCmd.Flags().StringVar(&researchType, "type", "neural", "search type (neural, keyword, magic)")
	researchCmd.Flags().BoolVar(&researchText, "text", false, "include full text (default: true if neither text nor highlights specified)")
	researchCmd.Flags().BoolVar(&researchHighlights, "highlights", false, "include highlights")
	researchCmd.Flags().StringVar(&researchStartDate, "start-date", "", "filter by published date (ISO format)")
	researchCmd.Flags().StringVar(&researchIncludeDomains, "include-domains", "", "comma-separated domains to include")
	researchCmd.Flags().StringVar(&researchExcludeDomains, "exclude-domains", "", "comma-separated domains to exclude")
	researchCmd.Flags().StringVar(&researchCategory, "category", "", "filter by category")
	researchCmd.Flags().BoolVar(&researchNoAutoprompt, "no-autoprompt", false, "disable Exa autoprompt")

	researchCmd.MarkFlagRequired("query")
}

func runResearch(cmd *cobra.Command, args []string) error {
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

	// Determine if text should be enabled by default
	includeText := researchText
	includeHighlights := researchHighlights
	
	// If neither text nor highlights was explicitly set, default text to true
	if !cmd.Flags().Changed("text") && !cmd.Flags().Changed("highlights") {
		includeText = true
		debugLog("Defaulting to text extraction for research mode")
	}

	// Build search request
	req := client.SearchRequest{
		Query:              researchQuery,
		NumResults:         researchNum,
		Type:               researchType,
		UseAutoprompt:      !researchNoAutoprompt,
		IncludeText:        includeText,
		IncludeHighlights:  includeHighlights,
		StartPublishedDate: researchStartDate,
		Category:           researchCategory,
	}

	if researchIncludeDomains != "" {
		req.IncludeDomains = splitDomains(researchIncludeDomains)
	}
	if researchExcludeDomains != "" {
		req.ExcludeDomains = splitDomains(researchExcludeDomains)
	}

	debugLog("Executing research search: query=%s, num=%d, type=%s, text=%v", researchQuery, researchNum, researchType, includeText)

	// Execute search
	resp, attempts, err := c.Search(context.Background(), req)
	if err != nil {
		return outputError("request_failed", err.Error(), attempts, startTime)
	}

	// Output results
	return outputSuccess("research", researchQuery, "", cfg.Profiles[0].ID, "config", resp, attempts, startTime)
}
