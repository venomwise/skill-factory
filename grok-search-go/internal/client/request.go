package client

// Message is an OpenAI-compatible chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is an OpenAI-compatible chat completions request body.
type ChatRequest struct {
	Model       string         `json:"model"`
	Messages    []Message      `json:"messages"`
	Temperature float64        `json:"temperature"`
	Stream      bool           `json:"stream"`
	Extra       map[string]any `json:"-"`
}

// ResearchRequest contains all inputs needed for a Grok research request.
type ResearchRequest struct {
	BaseURL       string
	APIKey        string
	Model         string
	Mode          string
	Query         string
	SystemPrompt  string
	Timeout       int
	ExtraBody     map[string]any
	ExtraHeaders  map[string]any
	ProfileID     string
	ProfileSource string
}

// BuildChatRequest creates the base non-streaming chat completions body.
func BuildChatRequest(model, systemPrompt, query string, extra map[string]any) map[string]any {
	body := map[string]any{
		"model": model,
		"messages": []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: query},
		},
		"temperature": 0.2,
		"stream":      false,
	}
	for key, value := range extra {
		body[key] = value
	}
	return body
}
