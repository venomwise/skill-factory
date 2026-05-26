package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "db-explorer version %s\n", version)
		fmt.Fprintf(out, "commit: %s\n", commit)
		fmt.Fprintf(out, "built: %s\n", date)
		fmt.Fprintf(out, "go: %s\n", goVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
