package output

// Source is a normalized source citation in command output.
type Source struct {
	URL     string `json:"url"`
	Title   string `json:"title,omitempty"`
	Snippet string `json:"snippet,omitempty"`
}

// Attempt describes one profile attempt in command output.
type Attempt struct {
	ProfileID        string `json:"profileId"`
	OK               bool   `json:"ok"`
	Status           int    `json:"status,omitempty"`
	Failover         bool   `json:"failover,omitempty"`
	Cooldown         bool   `json:"cooldown,omitempty"`
	CooldownSeconds  int    `json:"cooldownSeconds,omitempty"`
	RemainingSeconds int    `json:"remainingSeconds,omitempty"`
	UntilText        string `json:"untilText,omitempty"`
	Detail           string `json:"detail,omitempty"`
}

// Result is the normalized successful command output.
type Result struct {
	OK            bool           `json:"ok"`
	Mode          string         `json:"mode"`
	Query         string         `json:"query"`
	ProfileID     string         `json:"profileId"`
	ProfileSource string         `json:"profileSource"`
	Attempts      []Attempt      `json:"attempts"`
	ConfigPath    string         `json:"config_path"`
	ConfigPaths   []string       `json:"config_paths"`
	BaseURL       string         `json:"base_url"`
	Model         string         `json:"model"`
	Content       string         `json:"content"`
	Sources       []Source       `json:"sources"`
	Raw           string         `json:"raw"`
	Usage         map[string]any `json:"usage"`
	ElapsedMS     int            `json:"elapsed_ms"`
}

// ErrorResponse is the normalized runtime error output.
type ErrorResponse struct {
	OK                bool      `json:"ok"`
	Error             string    `json:"error"`
	Detail            string    `json:"detail,omitempty"`
	FailoverExhausted bool      `json:"failoverExhausted,omitempty"`
	Attempts          []Attempt `json:"attempts,omitempty"`
	CooldownStateFile string    `json:"cooldownStateFile,omitempty"`
	ConfigPath        string    `json:"config_path,omitempty"`
	ConfigPaths       []string  `json:"config_paths,omitempty"`
	BaseURL           string    `json:"base_url,omitempty"`
	Model             string    `json:"model,omitempty"`
	ElapsedMS         int       `json:"elapsed_ms,omitempty"`
}
