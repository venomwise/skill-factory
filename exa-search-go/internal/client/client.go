package client

import (
	"context"
	"net/http"
	"time"

	"github.com/venomwise/skill-factory/exa-search/internal/config"
	"github.com/venomwise/skill-factory/exa-search/internal/debug"
)

// Client handles HTTP requests to the Exa API with failover support
type Client struct {
	profiles   []config.Profile
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
}

// SearchRequest represents a search API request
type SearchRequest struct {
	Query              string
	NumResults         int
	Type               string // neural, keyword, magic
	UseAutoprompt      bool
	IncludeText        bool
	IncludeHighlights  bool
	StartPublishedDate string
	IncludeDomains     []string
	ExcludeDomains     []string
	Category           string
}

// SearchResponse represents a search API response
type SearchResponse struct {
	Results            []Result
	ResolvedSearchType string
	RequestID          string
	SearchTime         float64
	CostDollars        float64
}

// Result represents a single search result
type Result struct {
	ID            string
	Title         string
	URL           string
	PublishedDate string
	Author        string
	Score         float64
	Text          string
	Highlights    []string
	Image         string
	Favicon       string
}

// Attempt tracks a single API request attempt for debugging
type Attempt struct {
	ProfileID string
	OK        bool
	Status    int
	Failover  bool
	Detail    string
}

// New creates a new Exa API client
func New(profiles []config.Profile, baseURL string, timeout time.Duration) *Client {
	return &Client{
		profiles: profiles,
		baseURL:  baseURL,
		timeout:  timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Search executes a search request with failover support
func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, []Attempt, error) {
	var attempts []Attempt
	var lastErr error

	for idx, profile := range c.profiles {
		// Build request payload
		payload := map[string]interface{}{
			"query":         req.Query,
			"numResults":    req.NumResults,
			"type":          req.Type,
			"useAutoprompt": req.UseAutoprompt,
		}

		if req.IncludeText || req.IncludeHighlights {
			payload["contents"] = map[string]bool{
				"text":       req.IncludeText,
				"highlights": req.IncludeHighlights,
			}
		}

		if req.StartPublishedDate != "" {
			payload["startPublishedDate"] = req.StartPublishedDate
		}
		if len(req.IncludeDomains) > 0 {
			payload["includeDomains"] = req.IncludeDomains
		}
		if len(req.ExcludeDomains) > 0 {
			payload["excludeDomains"] = req.ExcludeDomains
		}
		if req.Category != "" {
			payload["category"] = req.Category
		}

		// Determine base URL (profile override or client default)
		baseURL := c.baseURL
		if profile.BaseURL != "" {
			baseURL = profile.BaseURL
		}

		// Execute request
		resp, attempt, err := c.executeRequest(ctx, baseURL+"/search", profile.APIKey, payload)
		attempt.ProfileID = profile.ID
		attempts = append(attempts, attempt)

		if err == nil {
			// Success - parse response
			searchResp, parseErr := parseSearchResponse(resp)
			if parseErr != nil {
				return nil, attempts, parseErr
			}
			return searchResp, attempts, nil
		}

		lastErr = err

		// Check if we should failover
		if !attempt.Failover {
			// Non-failover error, stop immediately
			break
		}

		// Log failover decision
		if idx < len(c.profiles)-1 {
			debug.LogFailover(profile.ID, c.profiles[idx+1].ID, attempt.Detail)
		}
	}

	// All profiles failed
	return nil, attempts, lastErr
}

// FindSimilar finds pages similar to the given URL
func (c *Client) FindSimilar(ctx context.Context, url string, numResults int) (*SearchResponse, []Attempt, error) {
	var attempts []Attempt
	var lastErr error

	for idx, profile := range c.profiles {
		// Build request payload
		payload := map[string]interface{}{
			"url":        url,
			"numResults": numResults,
		}

		// Determine base URL (profile override or client default)
		baseURL := c.baseURL
		if profile.BaseURL != "" {
			baseURL = profile.BaseURL
		}

		// Execute request
		resp, attempt, err := c.executeRequest(ctx, baseURL+"/findSimilar", profile.APIKey, payload)
		attempt.ProfileID = profile.ID
		attempts = append(attempts, attempt)

		if err == nil {
			// Success - parse response
			searchResp, parseErr := parseSearchResponse(resp)
			if parseErr != nil {
				return nil, attempts, parseErr
			}
			return searchResp, attempts, nil
		}

		lastErr = err

		// Check if we should failover
		if !attempt.Failover {
			// Non-failover error, stop immediately
			break
		}

		// Log failover decision
		if idx < len(c.profiles)-1 {
			debug.LogFailover(profile.ID, c.profiles[idx+1].ID, attempt.Detail)
		}
	}

	// All profiles failed
	return nil, attempts, lastErr
}
