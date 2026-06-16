# Requirements Document: Web Access

## Introduction

`web-access` 是面向 AI coding agent 的统一 Web 访问 skill。它提供一个 `web-access/` skill 目录、一个 `web-access-go/` Go CLI、一个 TOML 配置文件和一组项目作用域发布 workflow，让 agent 可以通过同一个入口完成 source-first 检索和实时综合分析。

系统边界包括新的 `web-access/` skill package、新的 `web-access-go/` Go module、Exa-backed source-first 命令、Grok-backed live synthesis 命令、统一配置解析、标准化输出、failover、release automation 和 eval 覆盖。第一版不修改、不删除、不重路由现有 `exa-search` 或 `grok-search` skill，也不提供旧 skill 的生命周期变更设计。

## Glossary

- **Web Access**: 新的统一 Web 访问 skill 和 CLI，包含 `web-access/` 与 `web-access-go/`。
- **Skill directory**: `web-access/` 目录，供 AI agent 读取 `SKILL.md`、references、配置示例和平台二进制。
- **Source project**: `web-access-go/` Go module，负责构建 `web-access` CLI。
- **Provider**: 后端能力来源。第一版包含 Exa 和 Grok 两类 provider。
- **Exa provider**: 用于官方文档、权威来源、通用 source-first 搜索、文本提取和相似页面发现的 provider。
- **Grok provider**: 用于最新新闻、社交讨论、广泛实时综合分析和官方信息与社区解读对比的 provider。
- **Profile**: 命名 API key 配置，可用于 failover，并可配合 provider-specific base URL、timeout 或 model 设置。
- **Failover**: 当当前 profile 遇到可重试失败时，记录 attempt 并尝试下一个可用 profile。
- **Cooldown**: 对失败 profile 的临时抑制机制。`web-access` 第一版不实现或持久化 cooldown 状态。
- **TOML**: `web-access` 的配置文件格式，默认路径为 `~/.config/ai-skills/web-access.toml`。
- **URLs-only output**: 只输出来源 URL 的格式，供 agent 收集来源链接。

## Requirements

### Requirement 1: Skill Package and Isolation

**User Story:** 作为 AI coding agent，我希望有一个独立的 `web-access` skill 入口，以便在不影响旧 skill 的情况下使用统一 Web 访问能力。

#### Acceptance Criteria

1. WHEN the skill package is created, THEN the system SHALL create `web-access/` with `SKILL.md`, `config.example.toml`, `references/configuration.md`, `references/query-recipes.md`, and `bin/`.
2. WHEN the source project is created, THEN the system SHALL create `web-access-go/` as a separate Go module.
3. WHEN supported platform binaries are delivered, THEN `web-access/bin/` SHALL contain Linux amd64, Linux arm64, macOS amd64, macOS arm64, and Windows amd64 binaries plus `SHA256SUMS`.
4. WHEN `web-access/SKILL.md` is written, THEN it SHALL describe when to use source-first commands versus live synthesis commands.
5. WHEN the first release is implemented, THEN the system SHALL NOT modify `exa-search/SKILL.md` or `grok-search/SKILL.md`.
6. WHEN the first release is implemented, THEN the system SHALL NOT delete or rename `exa-search/`, `grok-search/`, `exa-search-go/`, or `grok-search-go/`.
7. WHEN repository routing metadata is reviewed, THEN the system SHALL NOT change old skill frontmatter descriptions or old skill routing behavior.

### Requirement 2: Unified CLI Commands and Contracts

**User Story:** 作为 AI coding agent，我希望 `web-access` 提供清晰的命令族，以便根据任务语义选择 Exa 或 Grok 能力。

#### Acceptance Criteria

1. WHEN the binary runs `docs --query <text>`, THEN the system SHALL route the command to the Exa provider.
2. WHEN the binary runs `search --query <text>`, THEN the system SHALL route the command to the Exa provider.
3. WHEN the binary runs `extract --query <text>`, THEN the system SHALL route the command to the Exa provider and default to text extraction when neither `--text` nor `--highlights` is explicitly set.
4. WHEN the binary runs `similar --url <url>`, THEN the system SHALL route the command to the Exa provider.
5. WHEN the binary runs `news --query <text>`, THEN the system SHALL route the command to the Grok provider with the `news` prompt mode.
6. WHEN the binary runs `social --query <text>`, THEN the system SHALL route the command to the Grok provider with the `social` prompt mode.
7. WHEN the binary runs `research --query <text>`, THEN the system SHALL route the command to the Grok provider with the `research` prompt mode.
8. WHEN the binary runs `docs-compare --query <text>`, THEN the system SHALL route the command to the Grok provider with the `docs-compare` prompt mode.
9. WHEN the binary runs `version`, THEN the system SHALL print version metadata without requiring provider configuration.
10. IF a command is missing its required `--query` or `--url` input, THEN the system SHALL return a CLI validation error before making network requests.

### Requirement 3: Command Flags and Defaults

**User Story:** 作为 skill 维护者，我希望命令 flags 和默认值明确，以便实现和文档都能保持一致。

#### Acceptance Criteria

1. WHEN an Exa search-style command runs, THEN the system SHALL support `--num`, `--type`, `--text`, `--highlights`, `--start-date`, `--include-domains`, `--exclude-domains`, `--category`, and `--no-autoprompt`.
2. WHEN `docs` runs without `--include-domains`, THEN the system SHALL default included domains to `docs.openclaw.ai`.
3. WHEN `search` runs without `--include-domains`, THEN the system SHALL NOT apply a default domain filter.
4. WHEN `docs`, `search`, or `extract` runs without `--num`, THEN the system SHALL use `5`.
5. WHEN `docs`, `search`, or `extract` runs without `--type`, THEN the system SHALL use `neural`.
6. WHEN `similar` runs without `--num`, THEN the system SHALL use `5`.
7. WHEN a Grok command runs, THEN the system SHALL use only global Grok flags beyond its required `--query`.
8. IF an unsupported command/provider combination is requested through any explicit provider override, THEN the system SHALL reject it before making network requests.

### Requirement 4: Unified Configuration and Profile Resolution

**User Story:** 作为用户，我希望通过一个 TOML 配置统一管理 Exa 和 Grok 凭据，以便 agent 能稳定解析配置和覆盖项。

#### Acceptance Criteria

1. WHEN no `--config` flag is provided, THEN the system SHALL use `~/.config/ai-skills/web-access.toml` as the default config path.
2. WHEN the default config path is missing, THEN the system SHALL create a template config before returning any missing credential error.
3. WHEN resolving settings, THEN the system SHALL apply precedence from highest to lowest: CLI flags, `WEB_ACCESS_*` environment variables, provider-specific environment variables, TOML config, built-in defaults.
4. WHEN Exa provider config is resolved, THEN the system SHALL support base URL, timeout, profiles, `EXA_API_KEY`, `EXA_API_KEYS`, `EXA_BASE_URL`, `EXA_TIMEOUT`, and `WEB_ACCESS_EXA_API_KEY`.
5. WHEN Grok provider config is resolved, THEN the system SHALL support base URL, model, timeout, profiles, `GROK_API_KEY`, `GROK_API_KEYS`, `GROK_BASE_URL`, `GROK_MODEL`, `GROK_TIMEOUT`, and `WEB_ACCESS_GROK_API_KEY`.
6. WHEN `WEB_ACCESS_CONFIG` is set, THEN the system SHALL use it as the unified config path unless `--config` is provided.
7. WHEN placeholder API keys such as `YOUR_EXA_API_KEY` or `YOUR_GROK_API_KEY` are encountered, THEN the system SHALL ignore them as invalid credentials.
8. IF the TOML config contains invalid syntax, THEN the system SHALL return `config_parse_error` with config path and parse detail.

### Requirement 5: Exa Source-First Provider

**User Story:** 作为 AI coding agent，我希望 source-first 命令由 Exa 支撑，以便获取官方文档、权威来源、正文内容和相似页面。

#### Acceptance Criteria

1. WHEN `docs`, `search`, or `extract` runs, THEN the system SHALL build an Exa search request with query, result count, search type, autoprompt setting, optional text extraction, optional highlights, date filter, domain filters, and category.
2. WHEN `similar` runs, THEN the system SHALL call Exa's similar-page endpoint with URL and result count.
3. WHEN Exa returns successful results, THEN the system SHALL preserve provider-specific `results`, `resolvedSearchType`, `requestId`, `searchTime`, and `costDollars` where available.
4. WHEN an Exa profile fails with rate limit, quota, or auth failure and another Exa profile is available, THEN the system SHALL record the failed attempt and try the next Exa profile.
5. IF no usable Exa API key is resolved for an Exa command, THEN the system SHALL return `missing_api_key` for provider `exa`.
6. IF the Exa endpoint is unreachable or times out, THEN the system SHALL return `request_failed` with elapsed time, provider, profile attempts, and timeout value.

### Requirement 6: Grok Live Synthesis Provider

**User Story:** 作为 AI coding agent，我希望 live synthesis 命令由 Grok 支撑，以便获取实时新闻、社区讨论和综合分析。

#### Acceptance Criteria

1. WHEN a Grok command runs, THEN the system SHALL construct an OpenAI-compatible chat completion request.
2. WHEN the mode is `news`, `social`, `research`, or `docs-compare`, THEN the system SHALL use the corresponding mode-specific prompt.
3. WHEN Grok config contains model, base URL, timeout, extra body, or extra headers, THEN the system SHALL apply those settings to the request.
4. WHEN Grok returns successful content, THEN the system SHALL preserve `content`, `sources`, `usage`, and `raw` where available.
5. WHEN a Grok profile fails with rate limit, quota, or auth failure and another Grok profile is available, THEN the system SHALL record the failed attempt and try the next Grok profile.
6. WHEN Grok profile failover occurs, THEN the system SHALL NOT write cooldown state or require `--ignore-cooldown`.
7. IF no usable Grok API key is resolved for a Grok command, THEN the system SHALL return `missing_api_key` for provider `grok`.
8. IF the Grok endpoint is unreachable or times out, THEN the system SHALL return `request_failed` with elapsed time, provider, profile attempts, and timeout value.

### Requirement 7: Normalized Output and Error Handling

**User Story:** 作为 agent 或人工用户，我希望输出格式稳定且错误结构清晰，以便程序消费和终端诊断都可靠。

#### Acceptance Criteria

1. WHEN no output format flag is provided, THEN the system SHALL render JSON by default.
2. WHEN `--plain` is provided, THEN the system SHALL render human-readable plain text.
3. WHEN `--urls` is provided, THEN the system SHALL render only source URLs.
4. WHEN successful output is rendered, THEN it SHALL include common metadata `ok`, `mode`, `provider`, `query`, `url`, `profileId`, `profileSource`, `attempts`, and `elapsedMS` where applicable.
5. WHEN Exa output is rendered, THEN the system SHALL preserve Exa `results` without forcing lossy conversion into Grok fields.
6. WHEN Grok output is rendered, THEN the system SHALL preserve Grok `content`, `sources`, `usage`, and `raw` without forcing lossy conversion into Exa fields.
7. WHEN provider response is empty or malformed, THEN the system SHALL return `response_parse_error` and include raw response only when debug mode is enabled or already part of safe output.
8. WHEN all profiles for a provider fail, THEN the system SHALL return `all_profiles_failed` with clear retry guidance and attempt details.

### Requirement 8: Project-Scoped Release Automation

**User Story:** 作为 repository maintainer，我希望 `web-access` 发布流程与其他 Go skill 隔离，以便发布一个 skill 不会影响其他 skill。

#### Acceptance Criteria

1. WHEN the test workflow is added, THEN `.github/workflows/web-access-test.yml` SHALL run only for PR changes under `web-access/**`, `web-access-go/**`, or its workflow file.
2. WHEN the release workflow is added, THEN `.github/workflows/web-access-release.yml` SHALL trigger only on `web-access-v*` tags.
3. WHEN release validation runs, THEN it SHALL run matrix tests on Ubuntu, macOS, and Windows before building release artifacts.
4. WHEN release workflow jobs are defined, THEN the build job SHALL depend on test and the release job SHALL depend on build.
5. WHEN binaries are built for release, THEN the workflow SHALL use `CGO_ENABLED=0` for Linux amd64/arm64, macOS amd64/arm64, and Windows amd64.
6. WHEN the update-skill workflow is added, THEN `.github/workflows/web-access-update-skill.yml` SHALL run after successful web-access release workflow completion.
7. WHEN the update-skill workflow commits binaries, THEN it SHALL update only `web-access/bin/**` and checksums.
8. WHEN non-`web-access-v*` tags are pushed, THEN web-access release workflows SHALL NOT run.

### Requirement 9: Validation Coverage

**User Story:** 作为 maintainer，我希望有自动化测试和 eval 覆盖，以便确认统一 CLI 行为与设计一致。

#### Acceptance Criteria

1. WHEN Go unit tests run, THEN they SHALL cover config precedence across CLI flags, `WEB_ACCESS_*`, provider-specific environment variables, TOML config, and defaults.
2. WHEN Go unit tests run, THEN they SHALL cover command-to-provider routing.
3. WHEN Go unit tests run, THEN they SHALL cover Exa request construction for `docs`, `search`, `extract`, and `similar`.
4. WHEN Go unit tests run, THEN they SHALL cover Grok prompt selection and request construction for `news`, `social`, `research`, and `docs-compare`.
5. WHEN Go unit tests run, THEN they SHALL cover output JSON shape for Exa and Grok provider families.
6. WHEN Go unit tests run, THEN they SHALL cover Exa and Grok failover behavior without cooldown state writes.
7. WHEN command tests run, THEN they SHALL cover `version`, help output, `extract` default text extraction, and `docs` default domain filtering.
8. WHEN eval cases are added, THEN `evals/web-access/` SHALL cover representative agent routing decisions and command examples.
9. WHEN local validation runs, THEN `go test`, formatting check, `go vet`, build, and `version` execution SHALL pass for `web-access-go`.
