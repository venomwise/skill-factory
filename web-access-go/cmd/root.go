package cmd

import (
	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/web-access/internal/debug"
)

var (
	// Global flags
	cfgFile    string
	profileID  string
	timeout    int
	debugMode  bool

	// Output flags
	plainOutput bool
	urlsOutput  bool
	jsonOutput  bool

	// Exa provider flags
	exaAPIKey string

	// Grok provider flags
	grokAPIKey       string
	grokModel        string
	extraBodyJSON    string
	extraHeadersJSON string

	// Version info (set by main)
	version   string
	commit    string
	date      string
	goVersion string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "web-access",
	Short: "Unified web access CLI for AI coding agents",
	Long: `web-access is a unified command-line tool for web access using Exa (source-first)
and Grok (live synthesis) providers. It supports official docs lookup, neural search,
text extraction, similar pages, fresh news, social discourse, and live research.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugMode {
			debug.Enable()
			debug.Log("Debug mode enabled")
		}
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets version information from main
func SetVersionInfo(v, c, d, g string) {
	version = v
	commit = c
	date = d
	goVersion = g
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default: ~/.config/ai-skills/web-access.toml)")
	rootCmd.PersistentFlags().StringVar(&profileID, "profile", "", "use specific profile from config")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 0, "request timeout in seconds (default: 30)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "enable debug logging")

	// Output flags
	rootCmd.PersistentFlags().BoolVar(&plainOutput, "plain", false, "output plain text format")
	rootCmd.PersistentFlags().BoolVar(&urlsOutput, "urls", false, "output URLs only")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output JSON format (default)")

	// Exa provider flags
	rootCmd.PersistentFlags().StringVar(&exaAPIKey, "exa-api-key", "", "Exa API key (overrides config)")

	// Grok provider flags
	rootCmd.PersistentFlags().StringVar(&grokAPIKey, "grok-api-key", "", "Grok API key (overrides config)")
	rootCmd.PersistentFlags().StringVar(&grokModel, "grok-model", "", "Grok model (default: grok-beta)")
	rootCmd.PersistentFlags().StringVar(&extraBodyJSON, "extra-body-json", "", "Extra body JSON for Grok requests")
	rootCmd.PersistentFlags().StringVar(&extraHeadersJSON, "extra-headers-json", "", "Extra headers JSON for Grok requests")
}

// getOutputFormat determines the output format from flags
func getOutputFormat() string {
	if plainOutput {
		return "plain"
	}
	if urlsOutput {
		return "urls"
	}
	return "json"
}

// debugLog prints debug messages to stderr if debug mode is enabled
func debugLog(format string, args ...interface{}) {
	debug.Log(format, args...)
}
