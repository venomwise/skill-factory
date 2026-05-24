package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/venomwise/skill-factory/exa-search/internal/debug"
)

// executeRequest performs an HTTP POST request with error detection
func (c *Client) executeRequest(ctx context.Context, url, apiKey string, payload map[string]interface{}) (map[string]interface{}, Attempt, error) {
	attempt := Attempt{
		OK:       false,
		Failover: false,
	}

	// Marshal payload
	body, err := json.Marshal(payload)
	if err != nil {
		attempt.Detail = fmt.Sprintf("failed to marshal payload: %v", err)
		return nil, attempt, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		attempt.Detail = fmt.Sprintf("failed to create request: %v", err)
		return nil, attempt, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Log request
	debug.LogHTTPRequest("POST", url, apiKey)

	// Execute request
	resp, err := c.httpClient.Do(req)
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
	debug.LogHTTPResponse(resp.StatusCode, string(respBody))

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
		"credit balance",
		"insufficient",
		"exhaust",
		"usage limit",
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

// parseSearchResponse parses the API response into SearchResponse
func parseSearchResponse(data map[string]interface{}) (*SearchResponse, error) {
	resp := &SearchResponse{}

	// Extract top-level fields
	if v, ok := data["resolvedSearchType"].(string); ok {
		resp.ResolvedSearchType = v
	}
	if v, ok := data["requestId"].(string); ok {
		resp.RequestID = v
	}
	if v, ok := data["searchTime"].(float64); ok {
		resp.SearchTime = v
	}
	if v, ok := data["costDollars"].(float64); ok {
		resp.CostDollars = v
	}

	// Extract results
	if results, ok := data["results"].([]interface{}); ok {
		for _, item := range results {
			if itemMap, ok := item.(map[string]interface{}); ok {
				result := normalizeResult(itemMap)
				resp.Results = append(resp.Results, result)
			}
		}
	}

	return resp, nil
}

// normalizeResult extracts and normalizes a single result
func normalizeResult(data map[string]interface{}) Result {
	result := Result{}

	if v, ok := data["id"].(string); ok {
		result.ID = v
	}
	if v, ok := data["title"].(string); ok {
		result.Title = v
	}
	if v, ok := data["url"].(string); ok {
		result.URL = v
	}
	if v, ok := data["publishedDate"].(string); ok {
		result.PublishedDate = v
	}
	if v, ok := data["author"].(string); ok {
		result.Author = v
	}
	if v, ok := data["score"].(float64); ok {
		result.Score = v
	}
	if v, ok := data["text"].(string); ok {
		result.Text = v
	}
	if v, ok := data["image"].(string); ok {
		result.Image = v
	}
	if v, ok := data["favicon"].(string); ok {
		result.Favicon = v
	}

	// Handle highlights array
	if highlights, ok := data["highlights"].([]interface{}); ok {
		for _, h := range highlights {
			if str, ok := h.(string); ok {
				result.Highlights = append(result.Highlights, str)
			}
		}
	}

	return result
}
