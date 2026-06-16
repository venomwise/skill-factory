package exa

// SearchRequest represents an Exa search request
type SearchRequest struct {
	Query              string
	NumResults         int
	Type               string // neural, keyword, auto
	UseAutoprompt      bool
	IncludeText        bool
	IncludeHighlights  bool
	StartPublishedDate string
	IncludeDomains     []string
	ExcludeDomains     []string
	Category           string
}

// SimilarRequest represents an Exa find similar request
type SimilarRequest struct {
	URL        string
	NumResults int
}

// SearchResponse represents an Exa API response
type SearchResponse struct {
	Results            []Result `json:"results"`
	ResolvedSearchType string   `json:"resolvedSearchType,omitempty"`
	RequestID          string   `json:"requestId,omitempty"`
	SearchTime         float64  `json:"searchTime,omitempty"`
	CostDollars        float64  `json:"costDollars,omitempty"`
}

// Result represents a single Exa search result
type Result struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	PublishedDate string   `json:"publishedDate,omitempty"`
	Author        string   `json:"author,omitempty"`
	Score         float64  `json:"score,omitempty"`
	Text          string   `json:"text,omitempty"`
	Highlights    []string `json:"highlights,omitempty"`
	Image         string   `json:"image,omitempty"`
	Favicon       string   `json:"favicon,omitempty"`
}
