package grok

// Response wraps the Grok API response
type Response struct {
	Content string                 `json:"content"`
	Sources []Source               `json:"sources,omitempty"`
	Usage   *Usage                 `json:"usage,omitempty"`
	Raw     map[string]interface{} `json:"raw,omitempty"`
}

// Source represents a cited source
type Source struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}
