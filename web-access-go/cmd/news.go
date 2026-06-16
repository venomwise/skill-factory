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

var newsQuery string

// newsCmd represents the news command (Grok provider)
var newsCmd = &cobra.Command{
	Use:   "news",
	Short: "Get fresh news and recent developments (Grok)",
	Long:  `Get fresh news and recent developments using Grok live synthesis`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if newsQuery == "" {
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
			return output.RenderError(getOutputFormat(), "grok", "news", newsQuery, "", 0, nil, err)
		}

		if len(cfg.Grok.Profiles) == 0 {
			return output.RenderError(getOutputFormat(), "grok", "news", newsQuery, "", 0, nil,
				fmt.Errorf("missing_api_key: grok"))
		}

		prompt, err := prompts.ForMode("news")
		if err != nil {
			return output.RenderError(getOutputFormat(), "grok", "news", newsQuery, "", 0, nil, err)
		}

		ctx := context.Background()
		resp, attempts, err := grok.ExecuteResearch(ctx, cfg.Grok, prompt, newsQuery)
		elapsed := time.Since(startTime).Milliseconds()

		if err != nil {
			return output.RenderError(getOutputFormat(), "grok", "news", newsQuery, "", elapsed, nil, err)
		}

		var profileID, profileSource string
		for _, a := range attempts {
			if a.OK {
				profileID = a.ProfileID
				profileSource = a.ProfileSource
				break
			}
		}

		return output.RenderGrokSuccess(getOutputFormat(), "news", newsQuery, profileID, profileSource, resp, attempts, elapsed)
	},
}

func init() {
	rootCmd.AddCommand(newsCmd)
	newsCmd.Flags().StringVar(&newsQuery, "query", "", "news query (required)")
}
