package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/venomwise/skill-factory/web-access/internal/providers/exa"
)

// ExaResult represents a successful Exa provider result
type ExaResult struct {
	OK            bool                     `json:"ok"`
	Mode          string                   `json:"mode"`
	Provider      string                   `json:"provider"`
	Query         string                   `json:"query,omitempty"`
	URL           string                   `json:"url,omitempty"`
	ProfileID     string                   `json:"profileId"`
	ProfileSource string                   `json:"profileSource"`
	Attempts      []map[string]interface{} `json:"attempts"`
	ElapsedMS     int64                    `json:"elapsedMS"`
	Results       []exa.Result             `json:"results"`
	Metadata      map[string]interface{}   `json:"metadata,omitempty"`
}

// RenderExaSuccess renders a successful Exa response
func RenderExaSuccess(format, mode, query, profileID, profileSource string, resp *exa.Response, attempts []exa.Attempt, elapsedMS int64) error {
	result := ExaResult{
		OK:            true,
		Mode:          mode,
		Provider:      "exa",
		Query:         query,
		ProfileID:     profileID,
		ProfileSource: profileSource,
		Attempts:      exa.FormatAttempts(attempts),
		ElapsedMS:     elapsedMS,
		Results:       resp.Results,
	}

	// Add metadata if present
	if resp.ResolvedSearchType != "" || resp.RequestID != "" {
		result.Metadata = map[string]interface{}{}
		if resp.ResolvedSearchType != "" {
			result.Metadata["resolvedSearchType"] = resp.ResolvedSearchType
		}
		if resp.RequestID != "" {
			result.Metadata["requestId"] = resp.RequestID
		}
		if resp.SearchTime > 0 {
			result.Metadata["searchTime"] = resp.SearchTime
		}
		if resp.CostDollars > 0 {
			result.Metadata["costDollars"] = resp.CostDollars
		}
	}

	switch format {
	case "plain":
		return renderExaPlain(result)
	case "urls":
		return renderExaURLs(result)
	default:
		return renderJSON(result)
	}
}

// RenderExaSimilarSuccess renders a successful Exa similar response
func RenderExaSimilarSuccess(format, url, profileID, profileSource string, resp *exa.Response, attempts []exa.Attempt, elapsedMS int64) error {
	result := ExaResult{
		OK:            true,
		Mode:          "similar",
		Provider:      "exa",
		URL:           url,
		ProfileID:     profileID,
		ProfileSource: profileSource,
		Attempts:      exa.FormatAttempts(attempts),
		ElapsedMS:     elapsedMS,
		Results:       resp.Results,
	}

	// Add metadata if present
	if resp.ResolvedSearchType != "" || resp.RequestID != "" {
		result.Metadata = map[string]interface{}{}
		if resp.ResolvedSearchType != "" {
			result.Metadata["resolvedSearchType"] = resp.ResolvedSearchType
		}
		if resp.RequestID != "" {
			result.Metadata["requestId"] = resp.RequestID
		}
		if resp.SearchTime > 0 {
			result.Metadata["searchTime"] = resp.SearchTime
		}
		if resp.CostDollars > 0 {
			result.Metadata["costDollars"] = resp.CostDollars
		}
	}

	switch format {
	case "plain":
		return renderExaPlain(result)
	case "urls":
		return renderExaURLs(result)
	default:
		return renderJSON(result)
	}
}

// renderExaPlain renders Exa results in plain text
func renderExaPlain(result ExaResult) error {
	fmt.Printf("Exa %s: %d results\n", result.Mode, len(result.Results))
	fmt.Printf("Provider: %s (profile: %s)\n", result.Provider, result.ProfileID)
	fmt.Printf("Elapsed: %dms\n\n", result.ElapsedMS)

	for i, r := range result.Results {
		fmt.Printf("%d. %s\n", i+1, r.Title)
		fmt.Printf("   %s\n", r.URL)
		if r.PublishedDate != "" {
			fmt.Printf("   Published: %s\n", r.PublishedDate)
		}
		if r.Text != "" {
			fmt.Printf("   %s\n", truncateText(r.Text, 200))
		}
		fmt.Println()
	}

	return nil
}

// renderExaURLs renders only URLs from Exa results
func renderExaURLs(result ExaResult) error {
	for _, r := range result.Results {
		fmt.Println(r.URL)
	}
	return nil
}

// renderJSON renders result as JSON
func renderJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

// truncateText truncates text to maxLen characters
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}
