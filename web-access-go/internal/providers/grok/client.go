package grok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/venomwise/skill-factory/web-access/internal/config"
	"github.com/venomwise/skill-factory/web-access/internal/debug"
)

// Client handles HTTP requests to the Grok API
type Client struct {
	profiles     []config.ResolvedProfile
	baseURL      string
	model        string
	timeout      time.Duration
	extraBody    map[string]interface{}
	extraHeaders map[string]string
	httpClient   *http.Client
}

// Attempt tracks a single API request attempt
type Attempt struct {
	ProfileID     string `json:"profileId"`
	ProfileSource string `json:"profileSource"`
	OK            bool   `json:"ok"`
	Status        int    `json:"status,omitempty"`
	Failover      bool   `json:"failover"`
	Detail        string `json:"detail,omitempty"`
}

// New creates a new Grok API client
func New(profiles []config.ResolvedProfile, baseURL, model string, timeout time.Duration, extraBody map[string]interface{}, extraHeaders map[string]string) *Client {
	return &Client{
		profiles:     profiles,
		baseURL:      baseURL,
		model:        model,
		timeout:      timeout,
		extraBody:    extraBody,
		extraHeaders: extraHeaders,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// DoResearch executes a research request with the given prompt
func (c *Client) DoResearch(ctx context.Context, prompt, query string) (*Response, []Attempt, error) {
	var attempts []Attempt
	var lastErr error

	for idx, profile := range c.profiles {
		// Build request
		req := ChatRequest{
			Model: c.model,
			Messages: []Message{
				{
					Role:    "user",
					Content: fmt.Sprintf("%s\n\nQuery: %s", prompt, query),
				},
			},
			ExtraBody: c.extraBody,
		}

		// Use profile-specific model if set
		if profile.Model != "" {
			req.Model = profile.Model
		}

		// Determine base URL (profile override or client default)
		baseURL := c.baseURL
		if profile.BaseURL != "" {
			baseURL = profile.BaseURL
		}

		// Execute request
		resp, attempt, err := c.executeRequest(ctx, baseURL+"/chat/completions", profile.APIKey, req)
		attempt.ProfileID = profile.ID
		attempt.ProfileSource = profile.ProfileSource
		attempts = append(attempts, attempt)

		if err == nil {
			// Success - parse response
			grokResp, parseErr := parseResponse(resp)
			if parseErr != nil {
				return nil, attempts, parseErr
			}
			return grokResp, attempts, nil
		}

		lastErr = err

		// Check if we should failover
		if !attempt.Failover {
			// Non-failover error, stop immediately
			break
		}

		// Log failover decision
		if idx < len(c.profiles)-1 {
			debug.Log("Grok failover: %s -> %s (%s)", profile.ID, c.profiles[idx+1].ID, attempt.Detail)
		}
	}

	// All profiles failed
	return nil, attempts, lastErr
}

// executeRequest performs an HTTP POST request
func (c *Client) executeRequest(ctx context.Context, url, apiKey string, req ChatRequest) (map[string]interface{}, Attempt, error) {
	attempt := Attempt{
		OK:       false,
		Failover: false,
	}

	// Build payload
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
	}

	// Merge extra body
	for k, v := range c.extraBody {
		payload[k] = v
	}
	for k, v := range req.ExtraBody {
		payload[k] = v
	}

	// Marshal payload
	body, err := json.Marshal(payload)
	if err != nil {
		attempt.Detail = fmt.Sprintf("failed to marshal payload: %v", err)
		return nil, attempt, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		attempt.Detail = fmt.Sprintf("failed to create request: %v", err)
		return nil, attempt, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	// Apply extra headers
	for k, v := range c.extraHeaders {
		httpReq.Header.Set(k, v)
	}

	// Log request
	debug.Log("Grok POST %s (model: %s)", url, req.Model)

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		attempt.Detail = fmt.Sprintf("request failed: %v", err)
		return nil, attempt, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	attempt.Status = resp.StatusCode

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		attempt.Detail = fmt.Sprintf("failed to read response: %v", err)
		return nil, attempt, fmt.Errorf("failed to read response: %w", err)
	}

	// Log response
	debug.Log("Grok response: HTTP %d", resp.StatusCode)

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		bodyStr := string(respBody)
		attempt.Detail = truncate(bodyStr, 200)

		// Check if this is a failover-eligible error
		if shouldFailover(resp.StatusCode, bodyStr) {
			attempt.Failover = true
			return nil, attempt, fmt.Errorf("HTTP %d: %s", resp.StatusCode, attempt.Detail)
		}

		return nil, attempt, fmt.Errorf("HTTP %d: %s", resp.StatusCode, bodyStr)
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		attempt.Detail = fmt.Sprintf("failed to parse JSON: %v", err)
		return nil, attempt, fmt.Errorf("failed to parse JSON: %w", err)
	}

	attempt.OK = true
	return result, attempt, nil
}

// parseResponse parses the Grok API response
func parseResponse(data map[string]interface{}) (*Response, error) {
	resp := &Response{
		Raw: data,
	}

	// Extract usage
	if usageData, ok := data["usage"].(map[string]interface{}); ok {
		resp.Usage = &Usage{}
		if v, ok := usageData["prompt_tokens"].(float64); ok {
			resp.Usage.PromptTokens = int(v)
		}
		if v, ok := usageData["completion_tokens"].(float64); ok {
			resp.Usage.CompletionTokens = int(v)
		}
		if v, ok := usageData["total_tokens"].(float64); ok {
			resp.Usage.TotalTokens = int(v)
		}
	}

	// Extract content from choices
	if choices, ok := data["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					resp.Content = content

					// Extract sources from content (simple pattern matching)
					resp.Sources = extractSources(content)
				}
			}
		}
	}

	return resp, nil
}

// extractSources extracts URLs from content as sources
func extractSources(content string) []Source {
	var sources []Source
	seen := make(map[string]bool)

	// Match URLs in markdown links [text](url) and bare URLs
	urlPattern := regexp.MustCompile(`https?://[^\s\)]+`)
	matches := urlPattern.FindAllString(content, -1)

	for _, url := range matches {
		// Clean URL
		url = strings.TrimRight(url, ".,;:")
		if !seen[url] {
			sources = append(sources, Source{URL: url})
			seen[url] = true
		}
	}

	return sources
}

// shouldFailover determines if an error should trigger failover
func shouldFailover(statusCode int, body string) bool {
	// Check status codes
	if statusCode == 429 || statusCode == 401 || statusCode == 403 {
		return true
	}

	// Check response body for rate limit/quota keywords
	lowerBody := strings.ToLower(body)
	keywords := []string{
		"rate limit",
		"quota",
		"credits",
		"insufficient",
		"exhaust",
	}

	for _, keyword := range keywords {
		if strings.Contains(lowerBody, keyword) {
			return true
		}
	}

	return false
}

// truncate shortens a string to the specified length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
