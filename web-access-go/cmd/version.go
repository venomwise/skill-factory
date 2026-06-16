package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print version information for web-access CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("web-access %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  date: %s\n", date)
		fmt.Printf("  go: %s\n", goVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
