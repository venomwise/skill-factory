# Grok Search Binary Design

## Summary

Convert `grok-search` from a Python-script skill into a Go-based multi-platform binary skill aligned with the completed `exa-search` architecture. Source code will live in `grok-search-go/`; the skill directory `grok-search/` will carry prebuilt binaries, documentation, and TOML configuration examples. The new CLI intentionally adopts subcommands and TOML config, so the old Python CLI and JSON config are not compatibility targets.

## Goals

- Add a new Go CLI project at `grok-search-go/`.
- Add prebuilt multi-platform binaries under `grok-search/bin/`:
  - Linux amd64 / arm64
  - macOS amd64 / arm64
  - Windows amd64
- Preserve the current Grok search capabilities:
  - `news`, `social`, `research`, and `docs-compare` modes
  - JSON, plain text, and URLs-only output
  - multi-profile failover
  - cooldown state for failing profiles
  - extra request body and header overrides
  - profile selection
  - debug and version commands
- Replace JSON config with TOML config at `~/.config/ai-skills/grok-search.toml`.
- Default to the official xAI OpenAI-compatible endpoint:
  - `base_url = "https://api.x.ai"`
  - `model = "grok-4.1-fast"`
- Allow profile-level `base_url` and `model` overrides for proxy or compatible endpoints.
- Add project-scoped GitHub Actions for test, release, and skill binary updates.
- Update `grok-search/SKILL.md` and references to document the binary workflow.

## Primary Users / Roles

- AI coding agents that invoke `grok-search` for real-time research, breaking updates, X/Twitter discourse, and broad live synthesis.
- Skill maintainers who need a dependency-free cross-platform delivery model.
- Users who configure Grok or OpenAI-compatible endpoints through shared skill configuration.

## Non-Goals

- Do not preserve `python3 scripts/grok_search.py --mode ...` as a supported interface.
- Do not support old JSON config as the official configuration format.
- Do not refactor `exa-search-go`.
- Do not create a shared Go library between `exa-search-go` and `grok-search-go`.
- Do not introduce a global monorepo release workflow.
- Do not implement streaming terminal output; requests remain non-streaming, with best-effort parsing for SSE-like responses if an endpoint returns them anyway.

## Context

`exa-search` has already been converted to a hybrid delivery model: Go source in `exa-search-go/`, prebuilt platform binaries in `exa-search/bin/`, and GitHub Actions for testing, releasing, and updating skill binaries. `grok-search` currently remains a Python script at `grok-search/scripts/grok_search.py` with JSON config files and `--mode`-based invocation.

The current Grok script includes useful behavior that should be preserved in the Go rewrite: mode-specific system prompts, OpenAI-compatible chat completions, JSON/plain/URLs output, multi-key failover, cooldown state, extra request body/header merging, profile selection, and structured error output.

## Discovery

### Key Discoveries

- The desired scope is a full conversion matching `exa-search`, not a minimal wrapper.
- Breaking changes are acceptable: the new CLI should use subcommands and TOML config rather than preserving the Python script interface.
- The existing Grok docs imply an official default endpoint, but the Python implementation effectively requires `base_url`; the Go rewrite should make this explicit by defaulting to `https://api.x.ai`.
- GitHub Actions must remain project-scoped so future Go skills do not all rebuild on every release tag.

### Scope Decisions

- Use `grok-search-go/` as a separate Go project, structurally similar to `exa-search-go/`.
- Keep Grok-specific concepts such as research modes and prompt templates; do not mechanically copy Exa's search/docs/similar command model.
- Use project-specific release tags such as `grok-search-v1.0.0`.
- Remove the old Python script from the documented path and plan to delete it during implementation.
- Replace `config.example.json` with `config.example.toml`.
- Preserve the `raw` field for non-JSON assistant responses to aid debugging.

## Proposed Solution

Build a new Go CLI for `grok-search` using the `exa-search-go` project as a structural reference. The CLI exposes Grok research modes as subcommands, loads TOML configuration with CLI/env overrides, calls OpenAI-compatible chat completions, performs failover and cooldown handling, and renders normalized output.

### Architecture

```text
grok-search-go/
  go.mod
  README.md
  cmd/
    grok-search/
      main.go
    root.go
    news.go
    social.go
    research.go
    docs_compare.go
    version.go
  internal/
    client/
      client.go
      request.go
      response.go
      failover.go
    config/
      config.go
      loader.go
      template.go
    cooldown/
      cooldown.go
      store.go
    output/
      output.go
      json.go
      plain.go
      urls.go
      errors.go
    prompts/
      prompts.go
    debug/
      logger.go
```

`exa-search-go` is a reference for project layout, Cobra-based CLI style, output organization, and release automation, but `grok-search-go` keeps its own domain model.

### Components

#### CLI layer: `cmd/`

Responsibilities:

- Define global flags:
  - `--config`
  - `--api-key`
  - `--base-url`
  - `--model`
  - `--timeout`
  - `--profile`
  - `--ignore-cooldown`
  - `--extra-body-json`
  - `--extra-headers-json`
  - `--plain`
  - `--urls`
  - `--json`
  - `--debug`
- Define commands:
  - `news`
  - `social`
  - `research`
  - `docs-compare`
  - `version`
- Route each research command to the shared request pipeline with the matching prompt.

Example CLI:

```bash
grok-search news --query "latest OpenAI updates" --plain
grok-search social --query "what are people saying about X now?"
grok-search research --query "summarize recent model release discussion"
grok-search docs-compare --query "compare official docs and community interpretation"
grok-search version
```

#### Prompt layer: `internal/prompts`

Responsibilities:

- Store mode-specific system prompts equivalent to current Python `MODE_SYSTEM_PROMPTS`.
- Return the correct prompt for each command.

#### Config layer: `internal/config`

Responsibilities:

- Load TOML config from `~/.config/ai-skills/grok-search.toml` by default.
- Auto-create a template config if missing.
- Merge configuration from CLI flags, environment variables, TOML, and built-in defaults.
- Resolve profiles from:
  - `--api-key`
  - `GROK_API_KEY`
  - `GROK_API_KEYS`
  - TOML `[[profiles]]`
- Filter placeholder API keys.
- Support profile-level `base_url` and `model` overrides.
- Parse extra body/header JSON from CLI or environment variables.

Environment variables:

```bash
GROK_API_KEY
GROK_API_KEYS
GROK_BASE_URL
GROK_MODEL
GROK_TIMEOUT
GROK_CONFIG
GROK_EXTRA_BODY_JSON
GROK_EXTRA_HEADERS_JSON
```

TOML schema:

```toml
base_url = "https://api.x.ai"
model = "grok-4.1-fast"
timeout = 120

[[profiles]]
id = "main"
api_key = "YOUR_GROK_API_KEY"

[[profiles]]
id = "proxy"
api_key = "YOUR_PROXY_KEY"
base_url = "https://your-compatible-endpoint.example"
model = "grok-custom-model"

[extra_body]
# search_parameters = {}

[extra_headers]
# X-Custom-Header = "value"

[cooldown]
enabled = true
state_file = "runtime/cooldowns.json"
default_minutes = 15
rate_limit_minutes = 20
quota_minutes = 60
auth_minutes = 360
```

Configuration priority, highest to lowest:

1. CLI flags
2. Environment variables
3. TOML config
4. Built-in defaults

#### Client layer: `internal/client`

Responsibilities:

- Send POST requests to `{base_url}/v1/chat/completions`.
- Build OpenAI-compatible request bodies:

```json
{
  "model": "grok-4.1-fast",
  "messages": [
    {"role": "system", "content": "..."},
    {"role": "user", "content": "..."}
  ],
  "temperature": 0.2,
  "stream": false
}
```

- Merge `extra_body` into the request body.
- Merge `extra_headers` into request headers.
- Parse standard JSON responses.
- Best-effort parse SSE-like `data: ...` responses if an endpoint returns them despite `stream=false`.
- Detect failover-worthy errors.

Failover is triggered by:

- HTTP 401
- HTTP 403
- HTTP 429
- Error text containing rate-limit, quota, credits, billing, exhausted, unauthorized, forbidden, or token-unavailable indicators.

#### Cooldown layer: `internal/cooldown`

Responsibilities:

- Load and save cooldown state.
- Skip profiles currently in cooldown unless `--ignore-cooldown` is set.
- Prune expired cooldown entries.
- Set cooldown duration by failure type:
  - auth errors: 360 minutes
  - quota/billing errors: 60 minutes
  - rate-limit errors: 20 minutes
  - other failover errors: 15 minutes

Default state file:

```text
runtime/cooldowns.json
```

Relative cooldown paths resolve under the skill/tool runtime directory chosen by the implementation.

#### Output layer: `internal/output`

Responsibilities:

- Render JSON output by default.
- Render human-readable output with `--plain`.
- Render only source URLs with `--urls`.
- Return structured errors consistently.

Successful JSON should include useful fields from the Python version:

```json
{
  "ok": true,
  "mode": "news",
  "query": "...",
  "profileId": "main",
  "profileSource": "config.profiles",
  "attempts": [],
  "config_path": "...",
  "config_paths": [],
  "base_url": "https://api.x.ai",
  "model": "grok-4.1-fast",
  "content": "...",
  "sources": [],
  "raw": "",
  "usage": {},
  "elapsed_ms": 123
}
```

### Data Flow

#### Successful request

1. User runs:

   ```bash
   grok-search news --query "..."
   ```

2. CLI maps `news` to the news system prompt.
3. Config loader resolves:
   - config path
   - base URL
   - model
   - profiles
   - timeout
   - cooldown config
   - extra body/header
4. Cooldown store filters unavailable profiles.
5. Client tries profiles in order.
6. First successful response is parsed.
7. Assistant message is normalized into:
   - `content`
   - `sources`
   - fallback `raw`
8. Output renderer prints JSON/plain/URLs.
9. Process exits with code `0`.

#### Failover request

1. Profile A fails with a failover-worthy status or error message.
2. Client records the attempt.
3. Cooldown store marks Profile A unavailable.
4. Client tries Profile B.
5. If Profile B succeeds, output includes both attempts.
6. Process exits with code `0`.

#### All profiles fail

Output includes:

```json
{
  "ok": false,
  "error": "all_profiles_failed",
  "attempts": []
}
```

Process exits with a non-zero status.

## Error Handling

### Missing API key

Return structured error:

```json
{
  "ok": false,
  "error": "missing_api_key",
  "detail": "Pass --api-key, set GROK_API_KEY/GROK_API_KEYS, or configure ~/.config/ai-skills/grok-search.toml"
}
```

### Missing config

Auto-create a template config, then return `missing_api_key` if no key exists.

### Invalid TOML

Return `invalid_config` with the parse detail.

### Invalid JSON overrides

Return `invalid_json` for invalid `--extra-body-json`, `--extra-headers-json`, `GROK_EXTRA_BODY_JSON`, or `GROK_EXTRA_HEADERS_JSON`.

### HTTP errors

- 401, 403, 429, quota, billing, and rate-limit errors trigger failover and cooldown.
- Non-failover HTTP errors stop the request and return a structured error.

### Malformed model response

- If assistant content is valid JSON with `content` and `sources`, normalize those fields.
- If assistant content is not JSON, keep the full assistant message in `raw`, extract URLs into `sources`, and leave `content` empty or best-effort populated according to the simplest implementation.

## GitHub Actions

Use a project-scoped workflow convention so Grok releases do not trigger Exa builds and future Go skills remain isolated.

### Tag convention

Use project-specific release tags:

```text
grok-search-v1.0.0
```

Do not use global release tags such as:

```text
v1.0.0
```

### New workflows

```text
.github/workflows/grok-search-test.yml
.github/workflows/grok-search-release.yml
.github/workflows/grok-search-update-skill.yml
```

### Test workflow

Only runs when Grok Go source or its test workflow changes:

```yaml
on:
  push:
    branches:
      - main
      - develop
    paths:
      - 'grok-search-go/**'
      - '.github/workflows/grok-search-test.yml'
  pull_request:
    paths:
      - 'grok-search-go/**'
      - '.github/workflows/grok-search-test.yml'
```

### Release workflow

Only runs for Grok release tags:

```yaml
on:
  push:
    tags:
      - 'grok-search-v*'
```

### Update skill workflow

Only responds to the Grok release workflow:

```yaml
on:
  workflow_run:
    workflows: ["Build and Release grok-search"]
    types:
      - completed
```

### Workflow isolation rules

- Use `working-directory: grok-search-go`.
- Use `cache-dependency-path: grok-search-go/go.sum`.
- Update only `grok-search/bin/**`.
- Never rebuild or update `exa-search/bin/**` from Grok workflows.

## Skill Directory Changes

Current structure:

```text
grok-search/
  SKILL.md
  config.example.json
  config.json
  references/
  scripts/
    grok_search.py
```

Target structure:

```text
grok-search/
  SKILL.md
  config.example.toml
  bin/
    grok-search-linux-amd64
    grok-search-linux-arm64
    grok-search-darwin-amd64
    grok-search-darwin-arm64
    grok-search-windows-amd64.exe
    README.md
    SHA256SUMS
  references/
    configuration.md
    query-recipes.md
    migration-from-python.md
```

Implementation should remove the old Python script from the documented path and delete it if no longer needed. `config.example.json` should be replaced by `config.example.toml`.

## Testing

### Go unit tests

Cover:

- config loading and precedence
- profile extraction
- placeholder key filtering
- profile-level base URL/model override
- cooldown duration mapping
- cooldown state read/write
- failover decision rules
- response parsing:
  - normal JSON response
  - assistant JSON content
  - plain assistant text with URLs
  - SSE-like `data:` response
- output formatting

### Integration-style tests

Use `httptest.Server` to simulate:

- successful OpenAI-compatible response
- 429 then success with backup profile
- all profiles in cooldown
- invalid response payload
- extra headers/body being passed through

### Workflow validation

Run:

```bash
go test ./...
go vet ./...
gofmt
go build
```

Cross-platform builds run through GitHub Actions.

## Documentation Updates

### `grok-search/SKILL.md`

Update examples to binary usage:

```bash
bin/grok-search-<platform> news --query "..."
```

Add platform detection instructions consistent with `exa-search`.

### `grok-search/references/configuration.md`

Replace JSON configuration docs with TOML configuration docs.

### `grok-search/references/query-recipes.md`

Replace Python examples with binary examples.

### `grok-search/references/migration-from-python.md`

Document command migration:

```bash
python3 scripts/grok_search.py --mode news --query "..."
```

becomes:

```bash
bin/grok-search-<platform> news --query "..."
```

Document config migration from JSON:

```json
{
  "base_url": "...",
  "model": "...",
  "profiles": []
}
```

To TOML:

```toml
base_url = "..."
model = "..."

[[profiles]]
id = "main"
api_key = "..."
```

## Open Questions

No blocking open questions. During implementation, use the recommended scope decisions:

- Delete or de-document the old Python script.
- Replace `config.example.json` with `config.example.toml`.
- Preserve `raw` for non-JSON assistant responses.
