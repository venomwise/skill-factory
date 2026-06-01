# Web Access Design

## Summary

Build `web-access` as the unified web access skill for AI coding agents. The skill will provide one documented entry point, one Go CLI binary, one configuration file, and one release pipeline while preserving the distinct strengths of the existing `exa-search` and `grok-search` tools: source-first retrieval through Exa and live synthesis through Grok.

## Goals

- Create a new `web-access/` skill directory with `SKILL.md`, references, config example, and prebuilt binaries.
- Create a new `web-access-go/` Go CLI project.
- Expose a single `web-access` binary for all supported platforms:
  - Linux amd64 / arm64
  - macOS amd64 / arm64
  - Windows amd64
- Preserve the existing Exa-backed capabilities:
  - documentation and canonical source search
  - general source-first search
  - text extraction / research retrieval
  - similar page discovery
- Preserve the existing Grok-backed capabilities:
  - fresh news research
  - social and community discourse research
  - broad live synthesis
  - official-docs-vs-community comparison
- Use one TOML config at `~/.config/ai-skills/web-access.toml`.
- Keep output machine-readable by default with JSON, plus plain text and URLs-only formats.
- Add project-scoped GitHub Actions for PR tests, tag-only release builds, and skill binary updates.
- Deprecate the old `exa-search` and `grok-search` skills in documentation without deleting them in the first release.

## Primary Users / Roles

- AI coding agents that need one reliable web access entry point for documentation lookup, source retrieval, real-time updates, and community synthesis.
- Skill maintainers who need a self-contained binary skill with clear provider boundaries and project-scoped release automation.
- Users who configure Exa, Grok, or compatible proxy endpoints through local skill config.

## Non-Goals

- Do not collapse Exa and Grok into an opaque generic search mode that hides provider differences.
- Do not remove `exa-search/`, `grok-search/`, `exa-search-go/`, or `grok-search-go/` in the first release.
- Do not preserve old binary names as hard compatibility targets inside `web-access`.
- Do not create a global monorepo release workflow.
- Do not add caching, local indexing, browser automation, or page rendering.
- Do not implement AI-based automatic routing as a required first-version behavior.

## Context

This repository is a skill factory where each skill is a root-level directory with a `SKILL.md` entry point and optional `bin/`, `references/`, `assets/`, and `scripts/`. Go-based skills use a hybrid architecture: source code lives in a sibling Go project, and prebuilt platform binaries live under the skill directory.

`exa-search` currently provides source-first web retrieval with commands such as `docs`, `search`, `research`, and `similar`. It is best suited for official documentation, API references, pricing pages, canonical source retrieval, and extracted page content. Its Go CLI lives in `exa-search-go/` and uses Cobra, Viper-style config loading, Exa API clients, failover, and JSON/plain/URLs output.

`grok-search` currently provides real-time research with commands such as `news`, `social`, `research`, and `docs-compare`. It is best suited for fresh updates, X/Twitter or community discourse, broad multi-source live synthesis, and comparing official claims with community interpretation. Its Go CLI lives in `grok-search-go/` and uses Cobra, TOML config, OpenAI-compatible chat completions, failover, cooldowns, prompt templates, and JSON/plain/URLs output.

## Discovery

### Key Discoveries

- The current two skills are complementary, not redundant. Exa is a source-first retrieval backend; Grok is a live synthesis backend.
- A simple documentation wrapper would provide a unified skill name but would not deliver the requested unified binary, config, or release model.
- A source-level merge fits the repository's existing Go skill pattern better than a wrapper skill.
- `grok-search-release.yml` already follows the newer `test -> build -> release` release structure. `exa-search-release.yml` still has older structure traits, so `web-access` should follow the newer repository workflow requirements rather than copying Exa's release file.
- Provider differences should remain visible in docs and command semantics so agents choose the correct retrieval mode.

### Scope Decisions

- Implement a new first-class `web-access` skill and `web-access-go` CLI instead of renaming either existing skill.
- Use command semantics to choose providers by default:
  - source-first commands route to Exa
  - live synthesis commands route to Grok
- Keep provider override support minimal in the first version. Add explicit `--provider` only where it is useful and unambiguous.
- Keep old skills as deprecated compatibility documentation for one release cycle.
- Use `web-access-v*` release tags and project-scoped workflows.

## Proposed Solution

Create `web-access-go` as a single Cobra-based Go CLI with provider-specific internal packages for Exa and Grok. Create `web-access` as the skill directory that documents when to use each command, ships prebuilt binaries, and contains configuration references. The first version should favor clear command names and predictable provider routing over broad automatic routing.

### Architecture

```text
web-access/
  SKILL.md
  config.example.toml
  bin/
    web-access-linux-amd64
    web-access-linux-arm64
    web-access-darwin-amd64
    web-access-darwin-arm64
    web-access-windows-amd64.exe
    SHA256SUMS
  references/
    configuration.md
    query-recipes.md
    migration-from-exa-grok.md

web-access-go/
  go.mod
  README.md
  cmd/
    web-access/
      main.go
    root.go
    docs.go
    search.go
    extract.go
    similar.go
    news.go
    social.go
    research.go
    docs_compare.go
    version.go
  internal/
    config/
    output/
    debug/
    providers/
      exa/
      grok/
    cooldown/
    prompts/
```

The CLI has shared command parsing, config resolution, output rendering, and debug handling. Provider-specific HTTP request construction, response parsing, failover behavior, and prompt handling stay isolated under `internal/providers/exa` and `internal/providers/grok`.

### Components

#### Skill Directory: `web-access/`

Responsibilities:

- Define the triggerable `SKILL.md` for unified web access.
- Teach agents which command to use for source-first retrieval versus live synthesis.
- Document platform binary selection.
- Include concise examples for JSON, plain text, and URLs-only output.
- Link to configuration, query recipes, and migration guidance.

#### CLI Layer: `web-access-go/cmd`

Responsibilities:

- Define global flags:
  - `--config`
  - `--profile`
  - `--timeout`
  - `--plain`
  - `--urls`
  - `--json`
  - `--debug`
- Define provider-specific global flags where needed:
  - `--exa-api-key`
  - `--grok-api-key`
  - `--grok-model`
  - `--ignore-cooldown`
  - `--extra-body-json`
  - `--extra-headers-json`
- Define commands:
  - `docs`
  - `search`
  - `extract`
  - `similar`
  - `news`
  - `social`
  - `research`
  - `docs-compare`
  - `version`

Default provider routing:

```text
docs         -> exa
search       -> exa
extract      -> exa
similar      -> exa
news         -> grok
social       -> grok
research     -> grok
docs-compare -> grok
```

#### Config Layer: `internal/config`

Responsibilities:

- Load `~/.config/ai-skills/web-access.toml` by default.
- Auto-create a template config if the default path is missing.
- Merge CLI flags, environment variables, TOML config, and built-in defaults.
- Resolve separate Exa and Grok provider configs.
- Support provider-specific profiles, base URLs, timeouts, and Grok model overrides.
- Preserve existing environment variable compatibility where practical:
  - `EXA_API_KEY`
  - `EXA_API_KEYS`
  - `EXA_BASE_URL`
  - `EXA_TIMEOUT`
  - `GROK_API_KEY`
  - `GROK_API_KEYS`
  - `GROK_BASE_URL`
  - `GROK_MODEL`
  - `GROK_TIMEOUT`
- Add explicit unified environment variables:
  - `WEB_ACCESS_EXA_API_KEY`
  - `WEB_ACCESS_GROK_API_KEY`
  - `WEB_ACCESS_CONFIG`

Config shape:

```toml
[exa]
base_url = "https://api.exa.ai"
timeout = 30

[[exa.profiles]]
id = "main"
api_key = "YOUR_EXA_API_KEY"

[grok]
base_url = "https://api.x.ai"
model = "grok-4.1-fast"
timeout = 120

[[grok.profiles]]
id = "main"
api_key = "YOUR_GROK_API_KEY"

[grok.cooldown]
enabled = true
state_file = "runtime/web-access-grok-cooldowns.json"
default_minutes = 15
rate_limit_minutes = 20
quota_minutes = 60
auth_minutes = 360
```

Configuration priority, highest to lowest:

1. CLI flags
2. `WEB_ACCESS_*` environment variables
3. existing provider-specific environment variables
4. TOML config
5. built-in defaults

#### Exa Provider: `internal/providers/exa`

Responsibilities:

- Implement source-first commands backed by the Exa API.
- Support search request options from the current `exa-search-go` CLI:
  - result count
  - search type
  - text extraction
  - highlights
  - published date filters
  - include and exclude domains
  - category
  - autoprompt toggle
- Implement `similar` through Exa's similar-page endpoint.
- Preserve profile failover behavior and normalized attempt metadata.

#### Grok Provider: `internal/providers/grok`

Responsibilities:

- Implement live synthesis commands backed by OpenAI-compatible Grok chat completions.
- Preserve mode-specific prompts for `news`, `social`, `research`, and `docs-compare`.
- Support model, base URL, timeout, extra body, and extra header overrides.
- Preserve multi-profile failover and cooldown behavior.
- Preserve source extraction from model responses where available.

#### Output Layer: `internal/output`

Responsibilities:

- Render JSON by default.
- Render plain text for human terminal review.
- Render URLs-only output for source collection.
- Normalize common metadata across providers:
  - `ok`
  - `mode`
  - `provider`
  - `query`
  - `url`
  - `profileId`
  - `profileSource`
  - `attempts`
  - `elapsedMS`
- Preserve provider-specific result payloads without forcing lossy conversion:
  - Exa commands return `results`
  - Grok commands return `content`, `sources`, `usage`, and `raw`

#### Release Automation

Add:

- `.github/workflows/web-access-test.yml`
- `.github/workflows/web-access-release.yml`
- `.github/workflows/web-access-update-skill.yml`

Workflow requirements:

- Test workflow is PR-only with `web-access/**`, `web-access-go/**`, and workflow path filters.
- Release workflow triggers only on `web-access-v*` tags plus optional manual dispatch.
- Release workflow runs matrix tests on Ubuntu, macOS, and Windows before build.
- Build depends on test; release depends on build.
- Builds use `CGO_ENABLED=0` for Linux amd64/arm64, macOS amd64/arm64, and Windows amd64.
- Update-skill workflow runs on successful release workflow completion and optional manual dispatch, updating only `web-access/bin/**` and checksums.

#### Migration Documentation

Add `web-access/references/migration-from-exa-grok.md` with command mapping:

```text
exa-search docs      -> web-access docs
exa-search search    -> web-access search
exa-search research  -> web-access extract
exa-search similar   -> web-access similar
grok-search news     -> web-access news
grok-search social   -> web-access social
grok-search research -> web-access research
grok-search docs-compare -> web-access docs-compare
```

Update `exa-search/SKILL.md` and `grok-search/SKILL.md` to state that new work should prefer `web-access`, while keeping their existing commands documented for compatibility.

### Data Flow

#### Source-First Retrieval Path

1. Agent invokes a source-first command, for example:
   ```bash
   web-access docs --query "OpenClaw streaming API" --text --num 3
   ```
2. Cobra parses the command and maps it to the Exa provider.
3. Config resolves the Exa section from CLI flags, environment variables, TOML, and defaults.
4. The Exa provider builds the search or similar-page request.
5. The provider attempts the request with configured profiles, failing over on rate limit, auth, or quota failures.
6. Output renders normalized JSON with `provider: "exa"`, attempts, and Exa result data.

#### Live Synthesis Path

1. Agent invokes a live command, for example:
   ```bash
   web-access news --query "latest model release updates" --plain
   ```
2. Cobra parses the command and maps it to the Grok provider.
3. Config resolves the Grok section from CLI flags, environment variables, TOML, and defaults.
4. The Grok provider selects the mode-specific prompt and constructs an OpenAI-compatible chat completion request.
5. The provider attempts configured profiles, applying cooldowns for repeat failures.
6. Output renders normalized JSON or plain text with `provider: "grok"`, content, sources, usage, and attempts.

#### Agent Selection Flow

1. If the task requires official docs, API references, pricing pages, canonical retrieval, or extracted page text, the skill directs the agent to `docs`, `search`, `extract`, or `similar`.
2. If the task depends on freshness, public discussion, breaking updates, X/Twitter discourse, or synthesis across sources, the skill directs the agent to `news`, `social`, `research`, or `docs-compare`.
3. If both are needed, the skill directs the agent to run a source-first command first and then a live synthesis command for interpretation or freshness.

## Error Handling

- Missing provider API key:
  - Return structured JSON with `ok: false`, `provider`, `error: "missing_api_key"`, and actionable config guidance.
- Invalid config:
  - Return `config_parse_error` with the config path and parse detail.
- Unsupported command/provider combination:
  - Return a CLI validation error before making network requests.
- Exa rate limit, quota, or auth failure:
  - Fail over to the next Exa profile and include all attempts in output.
- Grok rate limit, quota, or auth failure:
  - Fail over to the next Grok profile, write cooldown state when configured, and include cooldown metadata in output.
- Network timeout:
  - Return `request_failed` with elapsed time, provider, profile attempts, and timeout value.
- Empty or malformed provider response:
  - Return `response_parse_error`; include raw response only when debug mode is enabled or when already part of existing safe output.
- All profiles unavailable:
  - Return `all_profiles_failed` or `all_profiles_in_cooldown` with clear retry guidance.

## Testing

- Unit test config precedence for:
  - CLI flags
  - `WEB_ACCESS_*` environment variables
  - provider-specific environment variables
  - TOML config
  - defaults
- Unit test command-to-provider routing.
- Unit test Exa request construction for `docs`, `search`, `extract`, and `similar`.
- Unit test Grok prompt selection and request construction for `news`, `social`, `research`, and `docs-compare`.
- Unit test output JSON shape for both provider families.
- Unit test failover and cooldown behavior.
- Add command smoke tests for `version` and help output.
- Add release workflow verification for:
  - `go test`
  - formatting check
  - `go vet`
  - build plus `version` execution on all release-test OSes.
- Add or update eval cases under `evals/web-access/` for representative agent routing decisions and command examples.

## Open Questions

None. Key decisions were confirmed during discovery.
