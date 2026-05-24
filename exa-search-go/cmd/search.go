package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/exa-search/internal/client"
	"github.com/venomwise/skill-factory/exa-search/internal/config"
	"github.com/venomwise/skill-factory/exa-search/internal/output"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "General neural search",
	Long:  `Execute a general neural search with various filters and options.`,
	RunE:  runSearch,
}

var (
	searchQuery            string
	searchNum              int
	searchType             string
	searchText             bool
	searchHighlights       bool
	searchStartDate        string
	searchIncludeDomains   string
	searchExcludeDomains   string
	searchCategory         string
	searchNoAutoprompt     bool
)

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchQuery, "query", "", "search query (required)")
	searchCmd.Flags().IntVar(&searchNum, "num", 5, "number of results")
	searchCmd.Flags().StringVar(&searchType, "type", "neural", "search type (neural, keyword, magic)")
	searchCmd.Flags().BoolVar(&searchText, "text", false, "include full text")
	searchCmd.Flags().BoolVar(&searchHighlights, "highlights", false, "include highlights")
	searchCmd.Flags().StringVar(&searchStartDate, "start-date", "", "filter by published date (ISO format)")
	searchCmd.Flags().StringVar(&searchIncludeDomains, "include-domains", "", "comma-separated domains to include")
	searchCmd.Flags().StringVar(&searchExcludeDomains, "exclude-domains", "", "comma-separated domains to exclude")
	searchCmd.Flags().StringVar(&searchCategory, "category", "", "filter by category (e.g., company, research paper, news)")
	searchCmd.Flags().BoolVar(&searchNoAutoprompt, "no-autoprompt", false, "disable Exa autoprompt")

	searchCmd.MarkFlagRequired("query")
}

func runSearch(cmd *cobra.Command, args []string) error {
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

	// Build search request
	req := client.SearchRequest{
		Query:             searchQuery,
		NumResults:        searchNum,
		Type:              searchType,
		UseAutoprompt:     !searchNoAutoprompt,
		IncludeText:       searchText,
		IncludeHighlights: searchHighlights,
		StartPublishedDate: searchStartDate,
		Category:          searchCategory,
	}

	if searchIncludeDomains != "" {
		req.IncludeDomains = splitDomains(searchIncludeDomains)
	}
	if searchExcludeDomains != "" {
		req.ExcludeDomains = splitDomains(searchExcludeDomains)
	}

	debugLog("Executing search: query=%s, num=%d, type=%s", searchQuery, searchNum, searchType)

	// Execute search
	resp, attempts, err := c.Search(context.Background(), req)
	if err != nil {
		return outputError("request_failed", err.Error(), attempts, startTime)
	}

	// Output results
	return outputSuccess("search", searchQuery, "", cfg.Profiles[0].ID, "config", resp, attempts, startTime)
}

func splitDomains(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func outputSuccess(mode, query, url, profileID, profileSource string, resp *client.SearchResponse, attempts []client.Attempt, startTime time.Time) error {
	data := output.OutputData{
		OK:                 true,
		Mode:               mode,
		Query:              query,
		URL:                url,
		ProfileID:          profileID,
		ProfileSource:      profileSource,
		Attempts:           output.ConvertAttempts(attempts),
		BaseURL:            baseURL,
		Results:            resp.Results,
		ResolvedSearchType: resp.ResolvedSearchType,
		RequestID:          resp.RequestID,
		SearchTime:         resp.SearchTime,
		CostDollars:        resp.CostDollars,
		ElapsedMS:          time.Since(startTime).Milliseconds(),
	}

	format := getOutputFormat()
	switch format {
	case "plain":
		return output.RenderPlain(data)
	case "urls":
		return output.RenderURLs(data)
	default:
		return output.RenderJSON(data)
	}
}

func outputError(errorCode, detail string, attempts []client.Attempt, startTime time.Time) error {
	data := output.OutputData{
		OK:        false,
		Error:     errorCode,
		Detail:    output.FormatErrorMessage(errorCode, detail, output.ConvertAttempts(attempts)),
		Attempts:  output.ConvertAttempts(attempts),
		ElapsedMS: time.Since(startTime).Milliseconds(),
	}

	format := getOutputFormat()
	if format == "plain" {
		fmt.Fprintln(os.Stderr, data.Detail)
		return fmt.Errorf("%s", errorCode)
	}

	// JSON output for errors
	return output.RenderJSON(data)
}
