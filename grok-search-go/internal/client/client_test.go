package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoResearchSendsRequestAndParsesJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("Authorization = %q", got)
		}
		if got := r.Header.Get("X-Test"); got != "yes" {
			t.Fatalf("X-Test = %q", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Decode body error = %v", err)
		}
		if body["model"] != "grok-test" {
			t.Fatalf("model = %v", body["model"])
		}
		if body["temperature"] != 0.2 {
			t.Fatalf("temperature = %v", body["temperature"])
		}
		if body["stream"] != false {
			t.Fatalf("stream = %v", body["stream"])
		}
		if body["custom"] != "value" {
			t.Fatalf("custom extra body = %v", body["custom"])
		}
		messages, ok := body["messages"].([]any)
		if !ok || len(messages) != 2 {
			t.Fatalf("messages = %#v", body["messages"])
		}
		first := messages[0].(map[string]any)
		second := messages[1].(map[string]any)
		if first["role"] != "system" || first["content"] != "system prompt" {
			t.Fatalf("system message = %#v", first)
		}
		if second["role"] != "user" || second["content"] != "user query" {
			t.Fatalf("user message = %#v", second)
		}

		assistantContent := `{"content":"answer","sources":[{"url":"https://example.com","title":"Example","snippet":"Snippet"}]}`
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model": "grok-test-response",
			"choices": []map[string]any{{
				"message": map[string]any{"role": "assistant", "content": assistantContent},
			}},
			"usage": map[string]any{"total_tokens": 12},
		})
	}))
	defer server.Close()

	resp, err := New(5).DoResearch(context.Background(), ResearchRequest{
		BaseURL:      server.URL,
		APIKey:       "test-key",
		Model:        "grok-test",
		Query:        "user query",
		SystemPrompt: "system prompt",
		ExtraBody:    map[string]any{"custom": "value"},
		ExtraHeaders: map[string]any{"X-Test": "yes"},
	})
	if err != nil {
		t.Fatalf("DoResearch() error = %v", err)
	}
	if resp.Model != "grok-test-response" {
		t.Fatalf("Model = %q", resp.Model)
	}
	if resp.AssistantContent.Content != "answer" {
		t.Fatalf("Content = %q", resp.AssistantContent.Content)
	}
	if len(resp.AssistantContent.Sources) != 1 || resp.AssistantContent.Sources[0].URL != "https://example.com" {
		t.Fatalf("Sources = %+v", resp.AssistantContent.Sources)
	}
	if resp.Usage["total_tokens"] != float64(12) {
		t.Fatalf("Usage = %+v", resp.Usage)
	}
}

func TestDoResearchParsesSSELikeResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"model\":\"grok-sse\",\"choices\":[{\"delta\":{\"content\":\"hello \"}}]}\n"))
		_, _ = w.Write([]byte("data: not-json\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"world https://example.com\"}}]}\n"))
		_, _ = w.Write([]byte("data: [DONE]\n"))
	}))
	defer server.Close()

	resp, err := New(5).DoResearch(context.Background(), ResearchRequest{
		BaseURL:      server.URL,
		APIKey:       "test-key",
		Model:        "grok-test",
		Query:        "query",
		SystemPrompt: "prompt",
	})
	if err != nil {
		t.Fatalf("DoResearch() error = %v", err)
	}
	if resp.Model != "grok-sse" {
		t.Fatalf("Model = %q", resp.Model)
	}
	if resp.AssistantContent.Raw != "hello world https://example.com" {
		t.Fatalf("Raw = %q", resp.AssistantContent.Raw)
	}
	if len(resp.AssistantContent.Sources) != 1 || resp.AssistantContent.Sources[0].URL != "https://example.com" {
		t.Fatalf("Sources = %+v", resp.AssistantContent.Sources)
	}
}

func TestDoResearchRequestFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := New(5).DoResearch(ctx, ResearchRequest{
		BaseURL:      "http://127.0.0.1:1",
		APIKey:       "test-key",
		Model:        "grok-test",
		Query:        "query",
		SystemPrompt: "prompt",
	})
	if err == nil {
		t.Fatalf("expected request error")
	}
	var reqErr *RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected RequestError, got %T", err)
	}
}

func TestParseAssistantContentPlainTextURLs(t *testing.T) {
	parsed := ParseAssistantContent("See https://a.example/path. Also https://b.example/x) and https://a.example/path")
	if parsed.Raw == "" {
		t.Fatalf("expected raw plain text")
	}
	if len(parsed.Sources) != 2 {
		t.Fatalf("Sources = %+v", parsed.Sources)
	}
	if parsed.Sources[0].URL != "https://a.example/path" || parsed.Sources[1].URL != "https://b.example/x" {
		t.Fatalf("Sources = %+v", parsed.Sources)
	}
}
