package output

import (
	"fmt"
	"io"
)

// RenderPlain writes human-readable successful output.
func RenderPlain(w io.Writer, result Result) error {
	if result.ProfileID != "" {
		suffix := ""
		if result.ProfileSource != "" {
			suffix = " [" + result.ProfileSource + "]"
		}
		if _, err := fmt.Fprintf(w, "Profile: %s%s\n", result.ProfileID, suffix); err != nil {
			return err
		}
	}
	if len(result.Attempts) > 0 {
		if _, err := fmt.Fprintln(w, "Attempts:"); err != nil {
			return err
		}
		for _, attempt := range result.Attempts {
			if err := renderPlainAttempt(w, attempt); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	if result.Content != "" {
		if _, err := fmt.Fprintln(w, result.Content); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	} else if result.Raw != "" {
		if _, err := fmt.Fprintln(w, result.Raw); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	if len(result.Sources) > 0 {
		if _, err := fmt.Fprintln(w, "Sources:"); err != nil {
			return err
		}
		for i, source := range result.Sources {
			if _, err := fmt.Fprintf(w, "[%d] %s\n", i+1, source.Title); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, "URL: %s\n", source.URL); err != nil {
				return err
			}
			if source.Snippet != "" {
				if _, err := fmt.Fprintf(w, "Snippet: %s\n", shorten(source.Snippet, 220)); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
	}
	return nil
}

// RenderPlainError writes a human-readable error response.
func RenderPlainError(w io.Writer, response ErrorResponse) error {
	if _, err := fmt.Fprintf(w, "ERROR: %s\n", response.Error); err != nil {
		return err
	}
	if response.Detail != "" {
		if _, err := fmt.Fprintln(w, response.Detail); err != nil {
			return err
		}
	}
	if len(response.Attempts) > 0 {
		if _, err := fmt.Fprintln(w, "Attempts:"); err != nil {
			return err
		}
		for _, attempt := range response.Attempts {
			if err := renderPlainAttempt(w, attempt); err != nil {
				return err
			}
		}
	}
	return nil
}

func renderPlainAttempt(w io.Writer, attempt Attempt) error {
	status := "fail"
	if attempt.OK {
		status = "ok"
	}
	if attempt.Cooldown {
		status = "cooldown"
	}
	if _, err := fmt.Fprintf(w, "- %s | %s", attempt.ProfileID, status); err != nil {
		return err
	}
	if attempt.Status != 0 {
		if _, err := fmt.Fprintf(w, " | status=%d", attempt.Status); err != nil {
			return err
		}
	}
	if attempt.Failover {
		if _, err := fmt.Fprint(w, " | failover"); err != nil {
			return err
		}
	}
	if attempt.RemainingSeconds > 0 {
		if _, err := fmt.Fprintf(w, " | remaining=%ds", attempt.RemainingSeconds); err != nil {
			return err
		}
	}
	if attempt.UntilText != "" {
		if _, err := fmt.Fprintf(w, " | until=%s", attempt.UntilText); err != nil {
			return err
		}
	}
	if attempt.Detail != "" {
		if _, err := fmt.Fprintf(w, " | %s", shorten(attempt.Detail, 160)); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w)
	return err
}

func shorten(text string, limit int) string {
	if len(text) <= limit {
		return text
	}
	if limit <= 3 {
		return text[:limit]
	}
	return text[:limit-3] + "..."
}
