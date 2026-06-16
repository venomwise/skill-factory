package exa

// Response wraps the raw Exa API response data
type Response struct {
	*SearchResponse
	Raw map[string]interface{} `json:"raw,omitempty"`
}
