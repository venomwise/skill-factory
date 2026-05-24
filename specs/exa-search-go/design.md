# Exa Search Go Rewrite Design

## Summary

Rewrite the `exa-search` skill from Python to Go, producing a static binary that eliminates Python runtime dependencies, simplifies cross-platform distribution, and provides a modern CLI experience. The Go implementation will use cobra + viper for CLI/config management, TOML for configuration, and GitHub Actions for automated multi-platform builds.

## Goals

- **Zero-dependency distribution**: Users download a single binary with no Python/pip/venv setup required
- **Cross-platform consistency**: Support Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64) with identical behavior
- **Feature parity**: Preserve all current functionality (search/docs/research/similar modes, multi-profile failover, three output formats)
- **Modern UX**: Simplified configuration, user-friendly error messages, idiomatic Go CLI patterns
- **Automated releases**: GitHub Actions builds and publishes binaries on tag push

## Primary Users / Roles

- **Technical users**: Developers integrating exa-search into scripts, CI/CD pipelines, or AI agent workflows
- **Non-technical users**: End users who need neural search capabilities without managing Python environments
- **Skill authors**: Contributors maintaining and extending the exa-search skill

## Non-Goals

- Backward compatibility with Python script's configuration format (no existing users to migrate)
- Self-update mechanism (manual download is acceptable for v1)
- GUI or TUI interface (CLI-only)
- Caching or local result storage (stateless execution)
- Support for Python 2.x or legacy systems

## Context

### Current Implementation
- **Language**: Python 3.x with `requests` library
- **Size**: 547 lines in `scripts/exa_search.py`
- **Config**: JSON files at multiple paths with complex merge logic
- **Distribution**: Users must install Python, create venv, install dependencies
- **Pain points**: 
  - Python version fragmentation
  - Dependency installation failures
  - Virtual environment management complexity
  - Cross-platform path/permission issues

### Exa API
- REST API with two main endpoints: `/search` and `/findSimilar`
- Requires API key via `x-api-key` header
- Supports neural/keyword/magic search types
- Optional text extraction and highlights
- Rate limiting and quota enforcement (requires failover logic)

### Skill Integration
- Called from SKILL.md with command-line examples
- Must output structured JSON for programmatic consumption
- Supports plain text and URLs-only output for human use

## Proposed Solution

Rewrite the tool as a Go binary using modern Go ecosystem standards (cobra for CLI, viper for config, TOML for config files). Distribute pre-compiled binaries via GitHub Releases with automated builds for all major platforms.

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  (cobra commands: search, docs, research, similar)          │
└────────────────┬────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────┐
│                      Config Layer                            │
│  (viper: TOML file + env vars + flags, priority merging)   │
└────────────────┬────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────┐
│                      Client Layer                            │
│  (HTTP client with retry/failover, profile management)      │
└────────────────┬────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────┐
│                      Output Layer                            │
│  (JSON/plain/URLs formatter with error normalization)       │
└─────────────────────────────────────────────────────────────┘
```

### Components

#### 1. CLI Layer (`cmd/`)
- **Responsibility**: Parse commands and flags, orchestrate execution
- **Interface**: 
  - Root command with global flags (`--config`, `--api-key`, `--profile`, `--timeout`)
  - Subcommands: `search`, `docs`, `research`, `similar`
  - Output flags: `--plain`, `--urls`, `--json` (default)
- **Implementation**: cobra command tree

#### 2. Config Layer (`internal/config/`)
- **Responsibility**: Load and merge configuration from multiple sources
- **Interface**:
  ```go
  type Config struct {
      Profiles  []Profile
      BaseURL   string
      Timeout   time.Duration
  }
  
  type Profile struct {
      ID      string
      APIKey  string
      BaseURL string  // optional override
  }
  
  func Load(configPath string) (*Config, error)
  func (c *Config) GetProfile(id string) (*Profile, error)
  ```
- **Priority**: CLI flags > env vars > TOML file
- **Config paths**: 
  - Explicit: `--config` flag
  - Default: `~/.config/ai-skills/exa-search.toml`
- **Environment variables**:
  - `EXA_API_KEY`: single key (creates implicit profile)
  - `EXA_API_KEYS`: comma-separated keys (creates multiple profiles)
  - `EXA_BASE_URL`: override base URL
  - `EXA_TIMEOUT`: request timeout in seconds

#### 3. Client Layer (`internal/client/`)
- **Responsibility**: Execute API requests with failover and retry logic
- **Interface**:
  ```go
  type Client struct {
      profiles []Profile
      baseURL  string
      timeout  time.Duration
      httpClient *http.Client
  }
  
  type SearchRequest struct {
      Query            string
      NumResults       int
      Type             string  // neural, keyword, magic
      UseAutoprompt    bool
      IncludeText      bool
      IncludeHighlights bool
      StartPublishedDate string
      IncludeDomains   []string
      ExcludeDomains   []string
      Category         string
  }
  
  type SearchResponse struct {
      Results           []Result
      ResolvedSearchType string
      RequestID         string
      SearchTime        float64
      CostDollars       float64
  }
  
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
  
  func New(profiles []Profile, baseURL string, timeout time.Duration) *Client
  func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error)
  func (c *Client) FindSimilar(ctx context.Context, url string, numResults int) (*SearchResponse, error)
  ```
- **Failover logic**:
  - Iterate through profiles on rate limit (429), auth errors (401/403), or quota errors
  - Track attempts for debugging
  - Return aggregated error if all profiles fail

#### 4. Output Layer (`internal/output/`)
- **Responsibility**: Format results for different output modes
- **Interface**:
  ```go
  type OutputFormat string
  const (
      FormatJSON  OutputFormat = "json"
      FormatPlain OutputFormat = "plain"
      FormatURLs  OutputFormat = "urls"
  )
  
  type OutputData struct {
      OK              bool
      Mode            string
      Query           string
      URL             string
      ProfileID       string
      ProfileSource   string
      Attempts        []Attempt
      ConfigPath      string
      BaseURL         string
      Results         []Result
      Error           string
      Detail          string
      ElapsedMS       int64
  }
  
  type Attempt struct {
      ProfileID string
      OK        bool
      Status    int
      Failover  bool
      Detail    string
  }
  
  func Render(data OutputData, format OutputFormat) error
  ```
- **Error messages**: User-friendly with actionable guidance (see Error Handling section)

### Data Flow

#### Primary Path: Search Command
1. User runs: `exa-search search --query "OpenClaw streaming" --num 5`
2. Cobra parses command and flags
3. Viper loads config:
   - Check `--config` flag → `~/.config/ai-skills/exa-search.toml`
   - Merge environment variables (`EXA_API_KEY`, etc.)
   - Apply CLI flag overrides (`--api-key`, `--profile`)
4. Config layer validates and returns profile list
5. Client layer constructs `SearchRequest` from flags
6. Client attempts request with first profile:
   - Build HTTP request with `x-api-key` header
   - POST to `{baseURL}/search` with JSON payload
   - Parse response into `SearchResponse`
7. On success: Output layer renders results in requested format (JSON default)
8. On failover error (429/401/403): Retry with next profile, track attempt
9. On final failure: Output layer renders error with all attempts

#### Docs Command Specialization
- Inherits search flow
- Defaults `IncludeDomains` to `["docs.openclaw.ai"]` if not specified

#### Research Command Specialization
- Inherits search flow
- Defaults `IncludeText` to `true` if neither `--text` nor `--highlights` specified

#### Similar Command Flow
- Similar to search, but calls `FindSimilar` endpoint with `--url` parameter

## Error Handling

### Top Failure Modes

#### 1. Missing API Key
**Detection**: No profiles after config loading
**Response**:
```
Error: No API key configured

Provide an API key using one of these methods:
  1. Command line:  exa-search --api-key YOUR_KEY search ...
  2. Environment:   export EXA_API_KEY=YOUR_KEY
  3. Config file:   ~/.config/ai-skills/exa-search.toml

Example config file:
  [[profiles]]
  id = "main"
  api_key = "your-key-here"

Get your API key at: https://exa.ai/
```

#### 2. Rate Limit / Quota Exhausted
**Detection**: HTTP 429 or response body contains "rate limit", "quota", "credits"
**Response** (single profile):
```
Error: Rate limit exceeded

Your API key has reached its rate limit or quota.
Check your usage at: https://exa.ai/dashboard

Consider:
  - Waiting before retrying
  - Upgrading your plan
  - Adding a backup API key in config
```

**Response** (multiple profiles, all failed):
```
Error: All API keys exhausted

Tried 3 profiles, all failed:
  - main: rate limit exceeded
  - backup: rate limit exceeded
  - fallback: quota exhausted

Check your usage at: https://exa.ai/dashboard
```

#### 3. Network / Timeout Errors
**Detection**: HTTP client timeout or connection errors
**Response**:
```
Error: Request failed

Could not connect to Exa API (timeout after 30s).
Check your network connection and try again.

If the problem persists, check Exa status: https://status.exa.ai/
```

#### 4. Invalid Configuration
**Detection**: TOML parse error or validation failure
**Response**:
```
Error: Invalid configuration

Failed to parse config file: ~/.config/ai-skills/exa-search.toml
Line 5: invalid TOML syntax

Fix the syntax error or delete the file to use defaults.
```

#### 5. Invalid Command Usage
**Detection**: Missing required flags or invalid values
**Response**: Cobra's built-in usage help with specific error highlighted

### Error Output Format

All errors follow consistent structure in JSON mode:
```json
{
  "ok": false,
  "error": "rate_limit_exceeded",
  "detail": "Your API key has reached its rate limit or quota.",
  "attempts": [
    {"profileId": "main", "ok": false, "status": 429, "failover": true, "detail": "rate limit"}
  ],
  "configPath": "~/.config/ai-skills/exa-search.toml",
  "baseURL": "https://api.exa.ai",
  "elapsedMS": 1234
}
```

## Testing

### Unit Tests
- **Config loading**: Test TOML parsing, env var merging, priority resolution
- **Client failover**: Mock HTTP responses, verify profile iteration logic
- **Output formatting**: Verify JSON/plain/URLs output correctness
- **Error normalization**: Test all error code paths produce expected messages

### Integration Tests
- **Live API calls**: Test against real Exa API with test key (optional, gated by env var)
- **End-to-end**: Run compiled binary with various flag combinations, verify output

### Test Data
- Reuse existing `evals/exa-search/test_cases.json` as reference scenarios
- Add Go-specific test cases for config edge cases

### CI Testing
- GitHub Actions runs tests on Linux/macOS/Windows before building binaries
- Fail build if tests don't pass

## Implementation Plan

### Phase 1: Core Structure (Day 1)
- [ ] Initialize Go module: `go mod init github.com/your-org/skill-factory/exa-search`
- [ ] Set up project structure: `cmd/`, `internal/config/`, `internal/client/`, `internal/output/`
- [ ] Add dependencies: `cobra`, `viper`, `BurntSushi/toml`
- [ ] Implement basic cobra command tree (root + 4 subcommands)

### Phase 2: Config Layer (Day 1-2)
- [ ] Define TOML schema and Go structs
- [ ] Implement viper-based config loading with priority merging
- [ ] Add environment variable support
- [ ] Write unit tests for config scenarios

### Phase 3: Client Layer (Day 2-3)
- [ ] Implement HTTP client with search/findSimilar methods
- [ ] Add profile failover logic with attempt tracking
- [ ] Implement error detection (rate limit, auth, network)
- [ ] Write unit tests with mocked HTTP responses

### Phase 4: Output Layer (Day 3)
- [ ] Implement JSON formatter
- [ ] Implement plain text formatter
- [ ] Implement URLs-only formatter
- [ ] Add user-friendly error messages
- [ ] Write unit tests for all formats

### Phase 5: Integration & Testing (Day 4)
- [ ] Wire all layers together in cobra commands
- [ ] Test against real Exa API
- [ ] Compare output with Python version for parity
- [ ] Fix discrepancies

### Phase 6: Distribution (Day 4-5)
- [ ] Create GitHub Actions workflow for multi-platform builds
- [ ] Test compiled binaries on Linux/macOS/Windows
- [ ] Write installation documentation
- [ ] Update SKILL.md with new binary-based commands

### Phase 7: Documentation & Migration (Day 5)
- [ ] Write comprehensive README with setup instructions
- [ ] Document TOML config format with examples
- [ ] Add troubleshooting guide
- [ ] Create example config file template
- [ ] Update evals to use Go binary

## Decisions

### 1. Binary naming convention
**Decision**: `exa-search` (matches skill name, avoids conflicts)

### 2. GitHub Actions trigger strategy
**Decision**: Manual workflow dispatch (full control over releases)

### 3. Config file auto-creation
**Decision**: Yes, create `~/.config/ai-skills/exa-search.toml` with template on first run if missing
- Improves first-run experience
- Template includes comments explaining each field
- Users can still use env vars exclusively if they delete the file

### 4. Logging/debug mode
**Decision**: Yes, add `--debug` flag for verbose logging
- Logs config resolution, HTTP requests/responses, failover decisions
- Redact API keys in logs (show only first 8 chars: `abcd1234...`)
- Output to stderr to avoid polluting JSON output

### 5. Version command
**Decision**: Yes, add `exa-search version` command
- Shows version, commit hash, build date, Go version
- Injected via `-ldflags` at build time
- Example output:
  ```
  exa-search version 1.0.0
  commit: a1b2c3d
  built: 2026-05-24T10:30:00Z
  go: go1.22.3
  ```
