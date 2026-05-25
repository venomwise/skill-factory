package cmd

import (
	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/grok-search/internal/debug"
)

var (
	// Global flags.
	cfgFile          string
	apiKey           string
	baseURL          string
	model            string
	timeout          int
	profileID        string
	ignoreCooldown   bool
	extraBodyJSON    string
	extraHeadersJSON string
	debugMode        bool

	// Output flags.
	plainOutput bool
	urlsOutput  bool
	jsonOutput  bool

	// Version info set by main.
	version   string
	commit    string
	date      string
	goVersion string
)

var rootCmd = &cobra.Command{
	Use:   "grok-search",
	Short: "Real-time web research powered by Grok-compatible chat completions",
	Long: `grok-search is a command-line tool for real-time research using an OpenAI-compatible Grok endpoint.
It supports mode-specific research commands, failover across API keys, cooldowns, and multiple output formats.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugMode {
			debug.Enable()
			debug.Log("Debug mode enabled")
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo stores build metadata injected by main.
func SetVersionInfo(v, c, d, g string) {
	version = v
	commit = c
	date = d
	goVersion = g
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default: ~/.config/ai-skills/grok-search.toml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Grok API key (overrides config)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "API base URL (default: https://api.x.ai)")
	rootCmd.PersistentFlags().StringVar(&model, "model", "", "model name (default: grok-4.1-fast)")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 0, "request timeout in seconds (default: 120)")
	rootCmd.PersistentFlags().StringVar(&profileID, "profile", "", "use specific profile from config")
	rootCmd.PersistentFlags().BoolVar(&ignoreCooldown, "ignore-cooldown", false, "attempt requests even if a profile is cooling down")
	rootCmd.PersistentFlags().StringVar(&extraBodyJSON, "extra-body-json", "", "extra JSON object merged into request body")
	rootCmd.PersistentFlags().StringVar(&extraHeadersJSON, "extra-headers-json", "", "extra JSON object merged into request headers")
	rootCmd.PersistentFlags().BoolVar(&plainOutput, "plain", false, "output plain text format")
	rootCmd.PersistentFlags().BoolVar(&urlsOutput, "urls", false, "output URLs only")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output JSON format (default)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "enable debug logging")
}

func getOutputFormat() string {
	if plainOutput {
		return "plain"
	}
	if urlsOutput {
		return "urls"
	}
	return "json"
}
