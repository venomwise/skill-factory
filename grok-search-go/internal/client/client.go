package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultUserAgent = "grok-search/1.0"

// Client sends OpenAI-compatible chat completion requests.
type Client struct {
	httpClient *http.Client
}

// New creates a client with the provided timeout in seconds.
func New(timeoutSeconds int) *Client {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 120
	}
	return &Client{httpClient: &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}}
}

// DoResearch sends one non-streaming Grok research request.
func (c *Client) DoResearch(ctx context.Context, req ResearchRequest) (*ResearchResponse, error) {
	body := BuildChatRequest(req.Model, req.SystemPrompt, req.Query, req.ExtraBody)
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, chatCompletionsURL(req.BaseURL), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", defaultUserAgent)
	for key, value := range req.ExtraHeaders {
		httpReq.Header.Set(key, fmt.Sprint(value))
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, &RequestError{Detail: err.Error(), Err: err}
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RequestError{StatusCode: resp.StatusCode, Detail: err.Error(), Err: err}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &RequestError{StatusCode: resp.StatusCode, Detail: string(raw)}
	}

	var parsed *ResearchResponse
	if strings.HasPrefix(strings.TrimSpace(string(raw)), "data:") {
		parsed, err = ParseSSELikeResponse(raw)
	} else {
		parsed, err = ParseChatResponse(raw)
	}
	if err != nil {
		return nil, &RequestError{StatusCode: resp.StatusCode, Detail: err.Error(), Err: err}
	}
	return parsed, nil
}

func chatCompletionsURL(baseURL string) string {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if strings.HasSuffix(baseURL, "/v1") {
		return baseURL + "/chat/completions"
	}
	return baseURL + "/v1/chat/completions"
}
