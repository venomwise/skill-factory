package output

import (
	"fmt"

	"github.com/venomwise/skill-factory/web-access/internal/providers/grok"
)

// GrokResult represents a successful Grok provider result
type GrokResult struct {
	OK            bool                     `json:"ok"`
	Mode          string                   `json:"mode"`
	Provider      string                   `json:"provider"`
	Query         string                   `json:"query"`
	ProfileID     string                   `json:"profileId"`
	ProfileSource string                   `json:"profileSource"`
	Attempts      []map[string]interface{} `json:"attempts"`
	ElapsedMS     int64                    `json:"elapsedMS"`
	Content       string                   `json:"content"`
	Sources       []grok.Source            `json:"sources,omitempty"`
	Usage         *grok.Usage              `json:"usage,omitempty"`
	Raw           map[string]interface{}   `json:"raw,omitempty"`
}

// RenderGrokSuccess renders a successful Grok response
func RenderGrokSuccess(format, mode, query, profileID, profileSource string, resp *grok.Response, attempts []grok.Attempt, elapsedMS int64) error {
	result := GrokResult{
		OK:            true,
		Mode:          mode,
		Provider:      "grok",
		Query:         query,
		ProfileID:     profileID,
		ProfileSource: profileSource,
		Attempts:      grok.FormatAttempts(attempts),
		ElapsedMS:     elapsedMS,
		Content:       resp.Content,
		Sources:       resp.Sources,
		Usage:         resp.Usage,
		Raw:           resp.Raw,
	}

	switch format {
	case "plain":
		return renderGrokPlain(result)
	case "urls":
		return renderGrokURLs(result)
	default:
		return renderJSON(result)
	}
}

// renderGrokPlain renders Grok results in plain text
func renderGrokPlain(result GrokResult) error {
	fmt.Printf("Grok %s\n", result.Mode)
	fmt.Printf("Provider: %s (profile: %s)\n", result.Provider, result.ProfileID)
	fmt.Printf("Elapsed: %dms\n\n", result.ElapsedMS)

	fmt.Println(result.Content)
	fmt.Println()

	if len(result.Sources) > 0 {
		fmt.Printf("Sources (%d):\n", len(result.Sources))
		for i, s := range result.Sources {
			fmt.Printf("%d. %s\n", i+1, s.URL)
		}
	}

	if result.Usage != nil {
		fmt.Printf("\nTokens: %d total (%d prompt + %d completion)\n",
			result.Usage.TotalTokens, result.Usage.PromptTokens, result.Usage.CompletionTokens)
	}

	return nil
}

// renderGrokURLs renders only URLs from Grok sources
func renderGrokURLs(result GrokResult) error {
	for _, s := range result.Sources {
		fmt.Println(s.URL)
	}
	return nil
}
