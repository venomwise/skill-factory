# Requirements Document: Exa Search Go Rewrite

## Introduction

The Exa Search Go binary is a command-line tool that provides neural web search capabilities through the Exa API, eliminating Python runtime dependencies and simplifying cross-platform distribution. It serves technical users integrating search into scripts and CI/CD pipelines, as well as non-technical users who need neural search without managing Python environments. The tool uses a four-layer architecture (CLI, Config, Client, Output) built with cobra for command parsing, viper for configuration management, and standard library HTTP for API communication. The system boundary includes search execution, configuration management, multi-profile failover, and formatted output, but excludes result caching, self-update mechanisms, and GUI interfaces.

## Glossary

- **Profile**: A configuration entry containing an API key and optional base URL override, used for failover when rate limits or quotas are exceeded
- **Failover**: Automatic retry logic that switches to the next configured profile when the current one encounters rate limiting (429), authentication errors (401/403), or quota exhaustion
- **cobra**: Go CLI framework used for command parsing and subcommand routing
- **viper**: Go configuration library that merges settings from TOML files, environment variables, and CLI flags
- **TOML**: Configuration file format (Tom's Obvious Minimal Language) used for `~/.config/ai-skills/exa-search.toml`
- **Neural search**: Exa's semantic search type that uses embeddings rather than keyword matching
- **Text extraction**: Optional feature to retrieve full page content in search results
- **Highlights**: Optional feature to retrieve relevant excerpts from search results

## Requirements

### Requirement 1: CLI Command Structure

**User Story:** As a user, I want a consistent command-line interface with subcommands for different search modes, so that I can easily discover and use the tool's capabilities.

#### Acceptance Criteria

1. WHEN the user runs `exa-search` without arguments, THEN the system SHALL display usage help with available subcommands.
2. WHEN the user runs `exa-search search`, THEN the system SHALL execute a general neural search.
3. WHEN the user runs `exa-search docs`, THEN the system SHALL execute a search defaulting to `docs.openclaw.ai` domain.
4. WHEN the user runs `exa-search research`, THEN the system SHALL execute a search with text extraction enabled by default.
5. WHEN the user runs `exa-search similar --url <URL>`, THEN the system SHALL find pages similar to the provided URL.
6. WHEN the user runs `exa-search version`, THEN the system SHALL display version, commit hash, build date, and Go version.
7. WHEN the user provides an invalid subcommand, THEN the system SHALL display an error and show available subcommands.

### Requirement 2: Configuration Loading and Priority

**User Story:** As a user, I want flexible configuration options with clear priority rules, so that I can override settings for different environments without editing files.

#### Acceptance Criteria

1. WHEN the user provides `--api-key` flag, THEN the system SHALL use that key and ignore all other sources.
2. WHEN the user sets `EXA_API_KEYS` environment variable, THEN the system SHALL create multiple profiles from the comma-separated list.
3. WHEN the user sets `EXA_API_KEY` environment variable, THEN the system SHALL create a single profile.
4. WHEN the user provides `--config` flag, THEN the system SHALL load configuration from the specified TOML file.
5. WHEN no `--config` flag is provided, THEN the system SHALL load from `~/.config/ai-skills/exa-search.toml` if it exists.
6. WHEN the configuration file does not exist on first run, THEN the system SHALL create it with a template including comments.
7. WHEN multiple configuration sources are present, THEN the system SHALL apply priority: CLI flags > environment variables > TOML file.
8. WHEN no API key is found from any source, THEN the system SHALL return an error with setup instructions.
9. WHEN the TOML file contains invalid syntax, THEN the system SHALL return an error indicating the file path and line number.

### Requirement 3: TOML Configuration Schema

**User Story:** As a user, I want a simple and well-documented configuration file format, so that I can easily set up multiple API keys and customize behavior.

#### Acceptance Criteria

1. WHEN the configuration file contains a `[[profiles]]` array, THEN the system SHALL load each profile with `id` and `api_key` fields.
2. WHEN a profile includes an optional `base_url` field, THEN the system SHALL use it to override the default Exa API endpoint for that profile.
3. WHEN the configuration file includes a top-level `timeout` field, THEN the system SHALL use it as the default request timeout in seconds.
4. WHEN the auto-generated template is created, THEN it SHALL include comments explaining each field and providing example values.
5. WHEN a profile's `api_key` is empty or contains placeholder text, THEN the system SHALL skip that profile during loading.

### Requirement 4: Multi-Profile Failover

**User Story:** As a user with multiple API keys, I want automatic failover when one key hits rate limits, so that my searches continue without manual intervention.

#### Acceptance Criteria

1. WHEN the first profile returns HTTP 429 (rate limit), THEN the system SHALL retry with the next profile.
2. WHEN the first profile returns HTTP 401 or 403 (authentication error), THEN the system SHALL retry with the next profile.
3. WHEN the first profile returns a response body containing "rate limit", "quota", or "credits", THEN the system SHALL retry with the next profile.
4. WHEN a profile succeeds, THEN the system SHALL stop trying additional profiles and return the result.
5. WHEN all profiles fail with failover-eligible errors, THEN the system SHALL return an error listing all attempts.
6. WHEN a profile fails with a non-failover error (e.g., network timeout), THEN the system SHALL stop and return the error without trying additional profiles.
7. WHEN the user provides `--profile <id>` flag, THEN the system SHALL use only that profile and skip failover.
8. WHEN failover occurs, THEN the system SHALL track each attempt with profile ID, status code, and error detail for debugging.

### Requirement 5: Search API Integration

**User Story:** As a user, I want to execute neural searches with various filters and options, so that I can find relevant documentation and web pages.

#### Acceptance Criteria

1. WHEN the user provides `--query`, THEN the system SHALL send the query to the Exa `/search` endpoint.
2. WHEN the user provides `--num N`, THEN the system SHALL request N results (default: 5).
3. WHEN the user provides `--type`, THEN the system SHALL use the specified search type (neural, keyword, or magic).
4. WHEN the user provides `--text`, THEN the system SHALL request full page text extraction.
5. WHEN the user provides `--highlights`, THEN the system SHALL request relevant excerpts.
6. WHEN the user provides `--include-domains`, THEN the system SHALL filter results to the comma-separated domain list.
7. WHEN the user provides `--exclude-domains`, THEN the system SHALL exclude results from the comma-separated domain list.
8. WHEN the user provides `--start-date`, THEN the system SHALL filter results published after the ISO date.
9. WHEN the user provides `--category`, THEN the system SHALL filter results by category (e.g., company, research paper, news).
10. WHEN the user provides `--no-autoprompt`, THEN the system SHALL disable Exa's automatic query enhancement.
11. WHEN the API returns results, THEN the system SHALL normalize them into a consistent structure with id, title, url, publishedDate, author, score, and optional text/highlights.

### Requirement 6: Similar Pages API Integration

**User Story:** As a user, I want to find pages similar to a canonical URL, so that I can discover related documentation or resources.

#### Acceptance Criteria

1. WHEN the user runs `exa-search similar --url <URL>`, THEN the system SHALL send the URL to the Exa `/findSimilar` endpoint.
2. WHEN the user provides `--num N`, THEN the system SHALL request N similar results (default: 5).
3. WHEN the API returns results, THEN the system SHALL normalize them into the same structure as search results.
4. WHEN the user omits `--url`, THEN the system SHALL return an error indicating the required flag.

### Requirement 7: Output Formatting

**User Story:** As a user, I want multiple output formats for different use cases, so that I can integrate the tool into scripts or read results directly.

#### Acceptance Criteria

1. WHEN no output flag is provided, THEN the system SHALL output structured JSON to stdout.
2. WHEN the user provides `--plain`, THEN the system SHALL output human-readable text with titles, URLs, scores, and text previews.
3. WHEN the user provides `--urls`, THEN the system SHALL output only URLs, one per line.
4. WHEN the output is JSON, THEN it SHALL include `ok`, `mode`, `query`, `profileId`, `results`, `attempts`, `configPath`, `baseURL`, and `elapsedMS` fields.
5. WHEN an error occurs, THEN the JSON output SHALL include `ok: false`, `error` code, `detail` message, and `attempts` array.
6. WHEN the output is plain text, THEN it SHALL include profile information, attempt summary, and formatted results with truncated text previews.
7. WHEN text extraction is enabled, THEN plain text output SHALL show the first 1200 characters with a truncation indicator if longer.

### Requirement 8: Error Handling and User Guidance

**User Story:** As a user, I want clear and actionable error messages when operations fail, so that I can quickly diagnose and fix issues.

#### Acceptance Criteria

1. WHEN no API key is configured, THEN the system SHALL display an error with three setup methods (CLI flag, environment variable, config file) and a link to get an API key.
2. WHEN a single profile hits rate limit, THEN the system SHALL display an error explaining rate limit/quota exhaustion with a link to the dashboard and suggestions (wait, upgrade, add backup key).
3. WHEN all profiles hit rate limit, THEN the system SHALL display an error listing all failed profiles with their specific error reasons.
4. WHEN a network timeout occurs, THEN the system SHALL display an error indicating connection failure with timeout duration and a link to check Exa status.
5. WHEN the configuration file has invalid TOML syntax, THEN the system SHALL display an error with the file path and line number.
6. WHEN a required flag is missing, THEN the system SHALL display cobra's usage help with the specific error highlighted.
7. WHEN debug mode is enabled with `--debug`, THEN the system SHALL log config resolution, HTTP requests/responses, and failover decisions to stderr.
8. WHEN debug mode logs API keys, THEN the system SHALL redact them to show only the first 8 characters followed by ellipsis.

### Requirement 9: Docs Command Specialization

**User Story:** As a user searching official documentation, I want a dedicated command that defaults to documentation domains, so that I get relevant results without manually specifying filters.

#### Acceptance Criteria

1. WHEN the user runs `exa-search docs --query <query>`, THEN the system SHALL default `--include-domains` to `docs.openclaw.ai`.
2. WHEN the user explicitly provides `--include-domains`, THEN the system SHALL use the user-provided value instead of the default.
3. WHEN the user provides other search flags (--num, --text, --type), THEN the system SHALL apply them as with the search command.

### Requirement 10: Research Command Specialization

**User Story:** As a user conducting deep research, I want a dedicated command that automatically extracts full text, so that I can analyze page content without additional flags.

#### Acceptance Criteria

1. WHEN the user runs `exa-search research --query <query>`, THEN the system SHALL default `--text` to true.
2. WHEN the user explicitly provides `--text` or `--highlights`, THEN the system SHALL use the user-provided values.
3. WHEN the user provides other search flags (--num, --type, --include-domains), THEN the system SHALL apply them as with the search command.

### Requirement 11: Build and Distribution

**User Story:** As a maintainer, I want automated multi-platform builds, so that users can download pre-compiled binaries without installing Go.

#### Acceptance Criteria

1. WHEN a maintainer triggers the GitHub Actions workflow manually, THEN the system SHALL compile binaries for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, and windows/amd64.
2. WHEN compilation completes, THEN the system SHALL inject version, commit hash, build date, and Go version via `-ldflags`.
3. WHEN binaries are built, THEN the system SHALL create a GitHub Release with all binaries and checksums.
4. WHEN tests fail during the workflow, THEN the system SHALL abort the build and report the failure.
5. WHEN the binary is executed, THEN it SHALL run without requiring any external dependencies or runtime environments.

### Requirement 12: Debug and Diagnostics

**User Story:** As a user troubleshooting issues, I want verbose logging to understand what the tool is doing, so that I can diagnose configuration or API problems.

#### Acceptance Criteria

1. WHEN the user provides `--debug`, THEN the system SHALL log configuration resolution steps to stderr.
2. WHEN the user provides `--debug`, THEN the system SHALL log HTTP request details (method, URL, headers with redacted API keys) to stderr.
3. WHEN the user provides `--debug`, THEN the system SHALL log HTTP response details (status code, body preview) to stderr.
4. WHEN the user provides `--debug`, THEN the system SHALL log failover decisions (which profile failed, why, which profile is next) to stderr.
5. WHEN debug logs include API keys, THEN the system SHALL show only the first 8 characters followed by `...` (e.g., `abcd1234...`).
6. WHEN debug mode is disabled, THEN the system SHALL output only the final result to stdout with no intermediate logs.
