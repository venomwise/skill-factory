package cmd

import (
	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/exa-search/internal/debug"
)

var (
	// Global flags
	cfgFile    string
	apiKey     string
	baseURL    string
	timeout    int
	profileID  string
	debugMode  bool
	
	// Output flags
	plainOutput bool
	urlsOutput  bool
	jsonOutput  bool
	
	// Version info (set by main)
	version   string
	commit    string
	date      string
	goVersion string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "exa-search",
	Short: "Neural web search powered by Exa API",
	Long: `exa-search is a command-line tool for neural web search using the Exa API.
It supports multiple search modes, failover across API keys, and flexible output formats.`,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default: ~/.config/ai-skills/exa-search.toml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Exa API key (overrides config)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "API base URL (default: https://api.exa.ai)")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 0, "request timeout in seconds (default: 30)")
	rootCmd.PersistentFlags().StringVar(&profileID, "profile", "", "use specific profile from config")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "enable debug logging")
	
	// Output flags
	rootCmd.PersistentFlags().BoolVar(&plainOutput, "plain", false, "output plain text format")
	rootCmd.PersistentFlags().BoolVar(&urlsOutput, "urls", false, "output URLs only")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output JSON format (default)")
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
