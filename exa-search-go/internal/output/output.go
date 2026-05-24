package output

import "github.com/venomwise/skill-factory/exa-search/internal/client"

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatJSON  OutputFormat = "json"
	FormatPlain OutputFormat = "plain"
	FormatURLs  OutputFormat = "urls"
)

// OutputData represents the complete output structure
type OutputData struct {
	OK            bool              `json:"ok"`
	Mode          string            `json:"mode,omitempty"`
	Query         string            `json:"query,omitempty"`
	URL           string            `json:"url,omitempty"`
	ProfileID     string            `json:"profileId,omitempty"`
	ProfileSource string            `json:"profileSource,omitempty"`
	Attempts      []Attempt         `json:"attempts,omitempty"`
	ConfigPath    string            `json:"configPath,omitempty"`
	BaseURL       string            `json:"baseURL,omitempty"`
	Results       []client.Result   `json:"results,omitempty"`
	Error         string            `json:"error,omitempty"`
	Detail        string            `json:"detail,omitempty"`
	ElapsedMS     int64             `json:"elapsedMS"`
	
	// Additional fields from API response
	ResolvedSearchType string  `json:"resolvedSearchType,omitempty"`
	RequestID          string  `json:"requestId,omitempty"`
	SearchTime         float64 `json:"searchTime,omitempty"`
	CostDollars        float64 `json:"costDollars,omitempty"`
}

// Attempt represents a single API request attempt
type Attempt struct {
	ProfileID string `json:"profileId"`
	OK        bool   `json:"ok"`
	Status    int    `json:"status,omitempty"`
	Failover  bool   `json:"failover,omitempty"`
	Detail    string `json:"detail,omitempty"`
}

// ConvertAttempts converts client.Attempt to output.Attempt
func ConvertAttempts(clientAttempts []client.Attempt) []Attempt {
	attempts := make([]Attempt, len(clientAttempts))
	for i, ca := range clientAttempts {
		attempts[i] = Attempt{
			ProfileID: ca.ProfileID,
			OK:        ca.OK,
			Status:    ca.Status,
			Failover:  ca.Failover,
			Detail:    ca.Detail,
		}
	}
	return attempts
}
