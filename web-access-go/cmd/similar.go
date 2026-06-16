package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/venomwise/skill-factory/web-access/internal/config"
	"github.com/venomwise/skill-factory/web-access/internal/output"
	"github.com/venomwise/skill-factory/web-access/internal/providers/exa"
)

var (
	similarURL string
	similarNum int
)

// similarCmd represents the similar command (Exa provider)
var similarCmd = &cobra.Command{
	Use:   "similar",
	Short: "Find similar pages (Exa)",
	Long:  `Find pages similar to a given URL using Exa API`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if similarURL == "" {
			return fmt.Errorf("--url is required")
		}

		startTime := time.Now()

		cfg, err := config.Load(config.Options{
			ConfigPath: cfgFile,
			ProfileID:  profileID,
			Timeout:    timeout,
			ExaAPIKey:  exaAPIKey,
		})
		if err != nil {
			return output.RenderError(getOutputFormat(), "exa", "similar", "", similarURL, 0, nil, err)
		}

		if len(cfg.Exa.Profiles) == 0 {
			return output.RenderError(getOutputFormat(), "exa", "similar", "", similarURL, 0, nil,
				fmt.Errorf("missing_api_key: exa"))
		}

		req := exa.SimilarRequest{
			URL:        similarURL,
			NumResults: similarNum,
		}

		ctx := context.Background()
		resp, attempts, err := exa.ExecuteSimilar(ctx, cfg.Exa, req)
		elapsed := time.Since(startTime).Milliseconds()

		if err != nil {
			return output.RenderError(getOutputFormat(), "exa", "similar", "", similarURL, elapsed, attempts, err)
		}

		var profileID, profileSource string
		for _, a := range attempts {
			if a.OK {
				profileID = a.ProfileID
				profileSource = a.ProfileSource
				break
			}
		}

		return output.RenderExaSimilarSuccess(getOutputFormat(), similarURL, profileID, profileSource, resp, attempts, elapsed)
	},
}

func init() {
	rootCmd.AddCommand(similarCmd)

	similarCmd.Flags().StringVar(&similarURL, "url", "", "reference URL (required)")
	similarCmd.Flags().IntVar(&similarNum, "num", 5, "number of results")
}
