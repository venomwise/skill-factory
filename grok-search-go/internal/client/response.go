package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
)

var urlPattern = regexp.MustCompile(`https?://[^\s)\]}>"']+`)

// ChatResponse is the subset of OpenAI-compatible chat completions response used by the tool.
type ChatResponse struct {
	ID      string         `json:"id,omitempty"`
	Model   string         `json:"model,omitempty"`
	Choices []Choice       `json:"choices"`
	Usage   map[string]any `json:"usage,omitempty"`
}

// Choice is one chat completion choice.
type Choice struct {
	Message Message `json:"message,omitempty"`
	Delta   Message `json:"delta,omitempty"`
}

// Source is a normalized source citation.
type Source struct {
	URL     string `json:"url"`
	Title   string `json:"title,omitempty"`
	Snippet string `json:"snippet,omitempty"`
}

// AssistantContent is the normalized assistant message payload.
type AssistantContent struct {
	Content string
	Sources []Source
	Raw     string
}

// ResearchResponse is the parsed response from one successful request.
type ResearchResponse struct {
	Model            string
	AssistantContent AssistantContent
	Usage            map[string]any
}

// RequestError captures HTTP or transport failure details.
type RequestError struct {
	StatusCode int
	Detail     string
	Err        error
}

func (e *RequestError) Error() string {
	if e.Detail != "" {
		return e.Detail
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "request failed"
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

// ParseChatResponse parses a standard JSON chat completions response.
func ParseChatResponse(raw []byte) (*ResearchResponse, error) {
	var chatResp ChatResponse
	if err := json.Unmarshal(raw, &chatResp); err != nil {
		return nil, err
	}

	message := ""
	if len(chatResp.Choices) > 0 {
		message = chatResp.Choices[0].Message.Content
	}
	return &ResearchResponse{
		Model:            chatResp.Model,
		AssistantContent: ParseAssistantContent(message),
		Usage:            chatResp.Usage,
	}, nil
}

// ParseSSELikeResponse parses data-prefixed chunks returned by some compatible endpoints.
func ParseSSELikeResponse(raw []byte) (*ResearchResponse, error) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	var content strings.Builder
	var model string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" || data == "[DONE]" {
			continue
		}
		var chunk ChatResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if model == "" {
			model = chunk.Model
		}
		if len(chunk.Choices) > 0 {
			content.WriteString(chunk.Choices[0].Delta.Content)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &ResearchResponse{
		Model:            model,
		AssistantContent: ParseAssistantContent(content.String()),
	}, nil
}

// ParseAssistantContent normalizes assistant content into content, sources, and raw fields.
func ParseAssistantContent(content string) AssistantContent {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return AssistantContent{}
	}

	var parsed struct {
		Content string   `json:"content"`
		Sources []Source `json:"sources"`
	}
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
			sources := normalizeSources(parsed.Sources)
			if len(sources) == 0 {
				sources = extractURLs(parsed.Content)
			}
			return AssistantContent{Content: parsed.Content, Sources: sources}
		}
	}

	return AssistantContent{Raw: content, Sources: extractURLs(content)}
}

func normalizeSources(sources []Source) []Source {
	out := make([]Source, 0, len(sources))
	seen := map[string]struct{}{}
	for _, source := range sources {
		url := strings.TrimSpace(source.URL)
		if url == "" {
			continue
		}
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		source.URL = url
		out = append(out, source)
	}
	return out
}

func extractURLs(text string) []Source {
	matches := urlPattern.FindAllString(text, -1)
	out := make([]Source, 0, len(matches))
	seen := map[string]struct{}{}
	for _, match := range matches {
		url := strings.TrimRight(match, ".,;:!?'")
		if url == "" {
			continue
		}
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		out = append(out, Source{URL: url})
	}
	return out
}
