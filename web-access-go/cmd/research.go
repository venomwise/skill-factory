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

var researchQuery string

// researchCmd represents the research command (Grok provider)
var researchCmd = &cobra.Command{
	Use:   "research",
	Short: "Conduct broad live research (Grok)",
	Long:  `Conduct broad live research and synthesis using Grok`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if researchQuery == "" {
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
			return output.RenderError(getOutputFormat(), "grok", "research", researchQuery, "", 0, nil, err)
		}

		if len(cfg.Grok.Profiles) == 0 {
			return output.RenderError(getOutputFormat(), "grok", "research", researchQuery, "", 0, nil,
				fmt.Errorf("missing_api_key: grok"))
		}

		prompt, err := prompts.ForMode("research")
		if err != nil {
			return output.RenderError(getOutputFormat(), "grok", "research", researchQuery, "", 0, nil, err)
		}

		ctx := context.Background()
		resp, attempts, err := grok.ExecuteResearch(ctx, cfg.Grok, prompt, researchQuery)
		elapsed := time.Since(startTime).Milliseconds()

		if err != nil {
			return output.RenderError(getOutputFormat(), "grok", "research", researchQuery, "", elapsed, nil, err)
		}

		var profileID, profileSource string
		for _, a := range attempts {
			if a.OK {
				profileID = a.ProfileID
				profileSource = a.ProfileSource
				break
			}
		}

		return output.RenderGrokSuccess(getOutputFormat(), "research", researchQuery, profileID, profileSource, resp, attempts, elapsed)
	},
}

func init() {
	rootCmd.AddCommand(researchCmd)
	researchCmd.Flags().StringVar(&researchQuery, "query", "", "research query (required)")
}
