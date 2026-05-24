package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/exa-search/internal/client"
	"github.com/venomwise/skill-factory/exa-search/internal/config"
)

var similarCmd = &cobra.Command{
	Use:   "similar",
	Short: "Find similar pages",
	Long:  `Find pages similar to a given URL.`,
	RunE:  runSimilar,
}

var (
	similarURL string
	similarNum int
)

func init() {
	rootCmd.AddCommand(similarCmd)

	similarCmd.Flags().StringVar(&similarURL, "url", "", "canonical URL to find similar pages (required)")
	similarCmd.Flags().IntVar(&similarNum, "num", 5, "number of results")

	similarCmd.MarkFlagRequired("url")
}

func runSimilar(cmd *cobra.Command, args []string) error {
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

	debugLog("Finding similar pages for: %s, num=%d", similarURL, similarNum)

	// Execute findSimilar
	resp, attempts, err := c.FindSimilar(context.Background(), similarURL, similarNum)
	if err != nil {
		return outputError("request_failed", err.Error(), attempts, startTime)
	}

	// Output results
	return outputSuccess("similar", "", similarURL, cfg.Profiles[0].ID, "config", resp, attempts, startTime)
}
