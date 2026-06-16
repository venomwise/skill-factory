package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/web-access/internal/config"
	"github.com/venomwise/skill-factory/web-access/internal/output"
	"github.com/venomwise/skill-factory/web-access/internal/prompts"
	"github.com/venomwise/skill-factory/web-access/internal/providers/grok"
)

var docsCompareQuery string

// docsCompareCmd represents the docs-compare command (Grok provider)
var docsCompareCmd = &cobra.Command{
	Use:   "docs-compare",
	Short: "Compare official docs with community interpretations (Grok)",
	Long:  `Compare official documentation with community interpretations and discussions using Grok live synthesis`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if docsCompareQuery == "" {
			return fmt.Errorf("--query is required")
		}

		startTime := time.Now()

		cfg, err := config.Load(config.Options{
			ConfigPath:       cfgFile,
			ProfileID:        profileID,
			Timeout:          timeout,
			GrokAPIKey:       grokAPIKey,
			GrokModel:        grokModel,
			ExtraBodyJSON:    extraBodyJSON,
			ExtraHeadersJSON: extraHeadersJSON,
		})
		if err != nil {
			return output.RenderError(getOutputFormat(), "grok", "docs-compare", docsCompareQuery, "", 0, nil, err)
		}

		if len(cfg.Grok.Profiles) == 0 {
			return output.RenderError(getOutputFormat(), "grok", "docs-compare", docsCompareQuery, "", 0, nil,
				fmt.Errorf("missing_api_key: grok"))
		}

		prompt, err := prompts.ForMode("docs-compare")
		if err != nil {
			return output.RenderError(getOutputFormat(), "grok", "docs-compare", docsCompareQuery, "", 0, nil, err)
		}

		ctx := context.Background()
		resp, attempts, err := grok.ExecuteResearch(ctx, cfg.Grok, prompt, docsCompareQuery)
		elapsed := time.Since(startTime).Milliseconds()

		if err != nil {
			return output.RenderError(getOutputFormat(), "grok", "docs-compare", docsCompareQuery, "", elapsed, nil, err)
		}

		var profileID, profileSource string
		for _, a := range attempts {
			if a.OK {
				profileID = a.ProfileID
				profileSource = a.ProfileSource
				break
			}
		}

		return output.RenderGrokSuccess(getOutputFormat(), "docs-compare", docsCompareQuery, profileID, profileSource, resp, attempts, elapsed)
	},
}

func init() {
	rootCmd.AddCommand(docsCompareCmd)
	docsCompareCmd.Flags().StringVar(&docsCompareQuery, "query", "", "comparison query (required)")
}
