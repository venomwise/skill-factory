package cmd

import "github.com/spf13/cobra"

var (
	cfgFile string
	profile string
	dbType  string
	url     string
	urlEnv  string
	format  string
	timeout int
	debug   bool

	version   string
	commit    string
	date      string
	goVersion string
)

var rootCmd = &cobra.Command{
	Use:           "db-explorer",
	Short:         "Agent-first read-only database explorer",
	SilenceUsage:  true,
	SilenceErrors: true,
	Long: `db-explorer is an Agent-first command-line tool for read-only database exploration.
It supports SQLite, PostgreSQL, and MySQL with JSON-first output.`,
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
	rootCmd.AddCommand(
		newTestCmd(),
		newSchemasCmd(),
		newTablesCmd(),
		newViewsCmd(),
		newSchemaCmd(),
		newDataCmd(),
		newQueryCmd(),
	)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "use specific profile from config")
	rootCmd.PersistentFlags().StringVar(&dbType, "db", "", "database type: sqlite, postgres, or mysql")
	rootCmd.PersistentFlags().StringVar(&url, "url", "", "database connection URL or SQLite path")
	rootCmd.PersistentFlags().StringVar(&urlEnv, "url-env", "", "environment variable containing the database connection URL")
	rootCmd.PersistentFlags().StringVar(&format, "format", "json", "output format: json, table, markdown, or csv")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 30, "query timeout in seconds")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug diagnostics")
}
