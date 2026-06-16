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

var socialQuery string

// socialCmd represents the social command (Grok provider)
var socialCmd = &cobra.Command{
	Use:   "social",
	Short: "Analyze social discourse and community discussions (Grok)",
	Long:  `Analyze social discourse and community discussions using Grok live synthesis`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if socialQuery == "" {
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
			return output.RenderError(getOutputFormat(), "grok", "social", socialQuery, "", 0, nil, err)
		}

		if len(cfg.Grok.Profiles) == 0 {
			return output.RenderError(getOutputFormat(), "grok", "social", socialQuery, "", 0, nil,
				fmt.Errorf("missing_api_key: grok"))
		}

		prompt, err := prompts.ForMode("social")
		if err != nil {
			return output.RenderError(getOutputFormat(), "grok", "social", socialQuery, "", 0, nil, err)
		}

		ctx := context.Background()
		resp, attempts, err := grok.ExecuteResearch(ctx, cfg.Grok, prompt, socialQuery)
		elapsed := time.Since(startTime).Milliseconds()

		if err != nil {
			return output.RenderError(getOutputFormat(), "grok", "social", socialQuery, "", elapsed, nil, err)
		}

		var profileID, profileSource string
		for _, a := range attempts {
			if a.OK {
				profileID = a.ProfileID
				profileSource = a.ProfileSource
				break
			}
		}

		return output.RenderGrokSuccess(getOutputFormat(), "social", socialQuery, profileID, profileSource, resp, attempts, elapsed)
	},
}

func init() {
	rootCmd.AddCommand(socialCmd)
	socialCmd.Flags().StringVar(&socialQuery, "query", "", "social query (required)")
}
