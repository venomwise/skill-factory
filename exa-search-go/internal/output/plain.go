package output

import (
	"fmt"
	"strings"
)

// RenderPlain outputs data as human-readable text to stdout
func RenderPlain(data OutputData) error {
	// Print profile info if available
	if data.ProfileID != "" {
		source := ""
		if data.ProfileSource != "" {
			source = fmt.Sprintf(" [%s]", data.ProfileSource)
		}
		fmt.Printf("Profile: %s%s\n", data.ProfileID, source)
	}

	// Print attempts summary if available
	if len(data.Attempts) > 0 {
		fmt.Println("Attempts:")
		for _, attempt := range data.Attempts {
			parts := []string{fmt.Sprintf("- %s", attempt.ProfileID)}
			if attempt.OK {
				parts = append(parts, "ok")
			} else {
				parts = append(parts, "fail")
			}
			if attempt.Status != 0 {
				parts = append(parts, fmt.Sprintf("status=%d", attempt.Status))
			}
			if attempt.Failover {
				parts = append(parts, "failover")
			}
			if attempt.Detail != "" {
				parts = append(parts, truncate(attempt.Detail, 160))
			}
			fmt.Println(strings.Join(parts, " | "))
		}
		fmt.Println()
	}

	// Print error if present
	if data.Error != "" {
		fmt.Printf("ERROR: %s\n", data.Error)
		if data.Detail != "" {
			fmt.Println(data.Detail)
		}
		return nil
	}

	// Print results
	for i, result := range data.Results {
		fmt.Printf("[%d] %s\n", i+1, result.Title)
		fmt.Printf("URL: %s\n", result.URL)
		if result.Score != 0 {
			fmt.Printf("Score: %.4f\n", result.Score)
		}
		if result.PublishedDate != "" {
			fmt.Printf("Published: %s\n", result.PublishedDate)
		}
		if result.Author != "" {
			fmt.Printf("Author: %s\n", result.Author)
		}
		if result.Text != "" {
			preview := truncate(result.Text, 1200)
			fmt.Println("Text:")
			fmt.Println(preview)
			if len(result.Text) > 1200 {
				fmt.Println("...[truncated]")
			}
		}
		if len(result.Highlights) > 0 {
			fmt.Println("Highlights:")
			for j, h := range result.Highlights {
				if j >= 5 {
					break
				}
				fmt.Printf("- %s\n", h)
			}
		}
		fmt.Println()
	}

	return nil
}

// truncate shortens a string to maxLen characters
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
