# Requirements Document: Grok Search Binary

## Introduction

Grok Search Binary converts the existing `grok-search` skill from a Python-script workflow into a Go-based, dependency-free, multi-platform command-line tool. The tool provides real-time research modes for breaking news, social discourse, broad synthesis, and official-docs-vs-community comparison through an OpenAI-compatible Grok chat completions endpoint.

The system boundary includes the new `grok-search-go/` source project, the binary-oriented `grok-search/` skill package, TOML configuration, normalized outputs, failover/cooldown behavior, project-scoped GitHub Actions, and updated skill documentation. Backward compatibility with the old Python CLI and JSON config is explicitly out of scope.

## Glossary

- **Grok Search Binary**: The new Go CLI tool delivered as platform-specific binaries under `grok-search/bin/`.
- **Skill directory**: The self-contained `grok-search/` directory used by AI agents, containing `SKILL.md`, references, config examples, and binaries.
- **Source project**: The `grok-search-go/` Go module that builds the Grok Search Binary.
- **OpenAI-compatible endpoint**: An HTTP API exposing `/v1/chat/completions` with OpenAI-style request and response shapes.
- **Profile**: A named API key configuration that may optionally override `base_url` and `model`.
- **Failover**: Trying the next available profile after a failover-worthy error.
- **Cooldown**: Temporary suppression of a failed profile to avoid repeatedly using a key that is rate-limited, exhausted, or invalid.
- **TOML**: The configuration file format used by the new tool at `~/.config/ai-skills/grok-search.toml`.
- **SSE**: Server-Sent Events; the tool only performs best-effort parsing of SSE-like `data:` responses and does not implement streaming terminal output.

## Requirements

### Requirement 1: Go CLI Structure and Commands

**User Story:** As an AI coding agent, I want a platform-native Grok CLI with clear subcommands, so that I can invoke real-time research modes without a Python runtime.

#### Acceptance Criteria

1. WHEN the source project is created, THEN the system SHALL place the Go module under `grok-search-go/`.
2. WHEN the binary is executed with `news --query <text>`, THEN the system SHALL run the news research mode.
3. WHEN the binary is executed with `social --query <text>`, THEN the system SHALL run the social/discourse research mode.
4. WHEN the binary is executed with `research --query <text>`, THEN the system SHALL run the broad synthesis research mode.
5. WHEN the binary is executed with `docs-compare --query <text>`, THEN the system SHALL run the official-docs-vs-community comparison mode.
6. WHEN the binary is executed with `version`, THEN the system SHALL print version metadata without requiring API configuration.
7. IF a research command is missing `--query`, THEN the system SHALL return a command-line error and non-zero exit status.
8. WHEN the old Python command shape `python3 scripts/grok_search.py --mode ...` is absent from documentation, THEN the system SHALL document the new subcommand form instead.

### Requirement 2: TOML Configuration and Precedence

**User Story:** As a user configuring Grok Search, I want one shared TOML config with clear override rules, so that agents can use the same credentials and endpoint settings consistently.

#### Acceptance Criteria

1. WHEN no `--config` flag is provided, THEN the system SHALL use `~/.config/ai-skills/grok-search.toml` as the default config path.
2. WHEN the default config file is missing, THEN the system SHALL create a template TOML config before reporting any missing credential error.
3. WHEN resolving settings, THEN the system SHALL apply precedence from highest to lowest: CLI flags, environment variables, TOML config, built-in defaults.
4. WHEN no base URL is provided by CLI, environment, or TOML, THEN the system SHALL default to `https://api.x.ai`.
5. WHEN no model is provided by CLI, environment, or TOML, THEN the system SHALL default to `grok-4.1-fast`.
6. WHEN the TOML config contains invalid syntax, THEN the system SHALL return an `invalid_config` error with parse detail.
7. WHEN `--extra-body-json`, `--extra-headers-json`, `GROK_EXTRA_BODY_JSON`, or `GROK_EXTRA_HEADERS_JSON` contains invalid JSON, THEN the system SHALL return an `invalid_json` error.

### Requirement 3: Profile Resolution and Overrides

**User Story:** As a user with one or more API keys, I want predictable profile selection and overrides, so that I can use primary, backup, or proxy endpoints safely.

#### Acceptance Criteria

1. WHEN `--api-key` is provided, THEN the system SHALL use it as a single CLI profile with highest credential priority.
2. WHEN `GROK_API_KEYS` is provided, THEN the system SHALL resolve comma-separated keys as ordered environment profiles.
3. WHEN `GROK_API_KEY` is provided and `GROK_API_KEYS` is not used, THEN the system SHALL resolve it as a single environment profile.
4. WHEN TOML `[[profiles]]` are configured and no higher-priority credentials are present, THEN the system SHALL use the configured profiles in file order.
5. WHEN a profile defines `base_url` or `model`, THEN the system SHALL use those values for that profile instead of global values.
6. WHEN `--profile <id>` is provided, THEN the system SHALL restrict execution to the matching resolved profile.
7. IF no non-placeholder API key is resolved, THEN the system SHALL return a `missing_api_key` error.
8. WHEN placeholder key values such as `YOUR_GROK_API_KEY` are encountered, THEN the system SHALL ignore them as invalid credentials.

### Requirement 4: OpenAI-Compatible Request Pipeline and Prompt Modes

**User Story:** As an agent invoking Grok Search, I want each command to send the correct prompt and request shape, so that responses are optimized for the selected research mode.

#### Acceptance Criteria

1. WHEN a research command runs, THEN the system SHALL send a POST request to `{base_url}/v1/chat/completions`.
2. WHEN building the request body, THEN the system SHALL include `model`, system message, user message, `temperature: 0.2`, and `stream: false`.
3. WHEN the mode is `news`, `social`, `research`, or `docs-compare`, THEN the system SHALL use the corresponding mode-specific system prompt.
4. WHEN `extra_body` is configured, THEN the system SHALL merge it into the request body.
5. WHEN `extra_headers` is configured, THEN the system SHALL merge it into the request headers.
6. WHEN a standard JSON chat completion response is returned, THEN the system SHALL parse the first assistant message content.
7. WHEN an SSE-like `data:` response is returned despite non-streaming mode, THEN the system SHALL best-effort reconstruct assistant content from chunks.
8. IF the endpoint is unreachable or times out, THEN the system SHALL return a structured request failure error.

### Requirement 5: Failover and Cooldown

**User Story:** As a user with multiple keys, I want failing keys to be skipped temporarily and backups tried automatically, so that rate limits or quota issues do not block research when alternatives exist.

#### Acceptance Criteria

1. WHEN a profile fails with HTTP 401, 403, or 429, THEN the system SHALL treat the failure as failover-worthy.
2. WHEN an error body contains rate-limit, quota, credits, billing, exhausted, unauthorized, forbidden, or token-unavailable indicators, THEN the system SHALL treat the failure as failover-worthy.
3. WHEN a failover-worthy error occurs and another profile is available, THEN the system SHALL record the failed attempt and try the next profile.
4. WHEN a failover-worthy error occurs, THEN the system SHALL write a cooldown entry for the failed profile if cooldown is enabled.
5. WHEN a profile is in cooldown, THEN the system SHALL skip it unless `--ignore-cooldown` is set.
6. WHEN cooldown entries expire, THEN the system SHALL prune them before profile selection completes.
7. WHEN all profiles are cooling down, THEN the system SHALL return `all_profiles_in_cooldown` with attempt details.
8. WHEN all profiles fail, THEN the system SHALL return `all_profiles_failed` with attempt details and a non-zero exit status.

### Requirement 6: Normalized Output Formats

**User Story:** As an AI agent or human user, I want consistent output formats, so that results can be consumed programmatically or read in a terminal.

#### Acceptance Criteria

1. WHEN no output format flag is provided, THEN the system SHALL emit pretty JSON output.
2. WHEN `--plain` is provided, THEN the system SHALL emit human-readable terminal output.
3. WHEN `--urls` is provided, THEN the system SHALL emit only source URLs, one per line.
4. WHEN assistant content is valid JSON containing `content` and `sources`, THEN the system SHALL normalize those fields into the output object.
5. WHEN assistant content is not valid JSON, THEN the system SHALL preserve the full assistant message in `raw` and extract any URLs into `sources`.
6. WHEN output is successful JSON, THEN it SHALL include `ok`, `mode`, `query`, `profileId`, `profileSource`, `attempts`, `config_path`, `config_paths`, `base_url`, `model`, `content`, `sources`, `raw`, `usage`, and `elapsed_ms` where available.
7. WHEN an error occurs, THEN the system SHALL emit a structured error object unless command-line parsing fails before output rendering.

### Requirement 7: Skill Packaging and Binary Delivery

**User Story:** As a skill maintainer, I want Grok Search packaged like Exa Search, so that agents can run it on major platforms without installing Python dependencies.

#### Acceptance Criteria

1. WHEN the skill is packaged, THEN `grok-search/bin/` SHALL contain binaries for Linux amd64, Linux arm64, macOS amd64, macOS arm64, and Windows amd64.
2. WHEN binaries are generated, THEN `grok-search/bin/SHA256SUMS` SHALL contain checksums for the generated binary files.
3. WHEN the binary directory is documented, THEN `grok-search/bin/README.md` SHALL describe supported platforms, file names, automatic updates, and manual builds.
4. WHEN `grok-search/SKILL.md` is updated, THEN it SHALL instruct agents to detect OS and architecture and select the matching binary.
5. WHEN JSON config examples are replaced, THEN the skill SHALL provide `grok-search/config.example.toml` instead of `config.example.json`.
6. WHEN the Python script is removed or de-documented, THEN all official usage examples SHALL point to `bin/grok-search-<platform>` commands.

### Requirement 8: Project-Scoped GitHub Actions

**User Story:** As a repository maintainer, I want Grok Search workflows isolated by project, so that releasing one Go skill does not rebuild or update unrelated skills.

#### Acceptance Criteria

1. WHEN Grok tests are configured, THEN `.github/workflows/grok-search-test.yml` SHALL run only for changes under `grok-search-go/**` or its own workflow file.
2. WHEN Grok release is configured, THEN `.github/workflows/grok-search-release.yml` SHALL trigger only on `grok-search-v*` tags or manual dispatch.
3. WHEN Grok skill binary updates are configured, THEN `.github/workflows/grok-search-update-skill.yml` SHALL respond only to the Grok release workflow or manual dispatch.
4. WHEN workflows set up Go caching, THEN they SHALL use `grok-search-go/go.sum` as the cache dependency path.
5. WHEN workflows build or test code, THEN they SHALL use `grok-search-go` as the working directory.
6. WHEN the update workflow commits binaries, THEN it SHALL update only `grok-search/bin/**`.
7. WHEN an `exa-search-v*` tag is pushed, THEN Grok release workflows SHALL NOT run.

### Requirement 9: Documentation and Migration Guidance

**User Story:** As a user migrating from the Python version, I want clear documentation for the new binary CLI and TOML config, so that I can update my usage without reading source code.

#### Acceptance Criteria

1. WHEN `grok-search/SKILL.md` is updated, THEN it SHALL describe when to use Grok Search and show binary-based examples.
2. WHEN `grok-search/references/configuration.md` is updated, THEN it SHALL document TOML config, profile overrides, environment variables, precedence, failover, and cooldown.
3. WHEN `grok-search/references/query-recipes.md` is updated, THEN all command examples SHALL use the new binary subcommands.
4. WHEN migration docs are added, THEN `grok-search/references/migration-from-python.md` SHALL map old `--mode` usage to new subcommands.
5. WHEN migration docs are added, THEN they SHALL map old JSON config shape to the new TOML shape.
6. WHEN docs mention official endpoints, THEN they SHALL state the default `https://api.x.ai` base URL and profile-level override behavior.

### Requirement 10: Validation Coverage

**User Story:** As a maintainer, I want automated checks for the rewrite, so that the binary behavior remains correct across platforms and future changes.

#### Acceptance Criteria

1. WHEN Go unit tests are run, THEN they SHALL cover config precedence, profile extraction, placeholder filtering, and profile-level overrides.
2. WHEN Go unit tests are run, THEN they SHALL cover cooldown duration mapping, cooldown state read/write, and failover decision rules.
3. WHEN Go unit tests are run, THEN they SHALL cover response parsing for standard JSON, assistant JSON content, plain assistant text with URLs, and SSE-like `data:` responses.
4. WHEN integration-style tests run, THEN they SHALL use `httptest.Server` to simulate success, failover after 429, all profiles in cooldown, invalid payloads, and extra body/header propagation.
5. WHEN local validation runs, THEN `go test ./...`, `go vet ./...`, `gofmt`, and `go build` SHALL pass for `grok-search-go`.
6. WHEN release validation runs, THEN GitHub Actions SHALL build all supported platform binaries.
