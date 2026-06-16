package main

import (
	"os"

	"github.com/venomwise/skill-factory/web-access/cmd"
)

// Version information injected at build time via -ldflags
var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	goVersion = "unknown"
)

func main() {
	// Pass version info to cmd package
	cmd.SetVersionInfo(version, commit, date, goVersion)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
