# Web Access Design

## Summary

构建 `web-access`，作为 AI coding agent 的统一 Web 访问 skill。该 skill 提供一个明确文档化的入口、一个 Go CLI 二进制、一个配置文件和一套发布流水线，同时保留现有 `exa-search` 与 `grok-search` 的不同优势：通过 Exa 做 source-first 检索，通过 Grok 做实时综合分析。

## Goals

- 新增 `web-access/` skill 目录，包含 `SKILL.md`、references、配置示例和预编译二进制。
- 新增 `web-access-go/` Go CLI 项目。
- 为所有支持平台提供单一 `web-access` 二进制：
  - Linux amd64 / arm64
  - macOS amd64 / arm64
  - Windows amd64
- 保留现有 Exa 能力：
  - 文档和权威来源搜索
  - 通用 source-first 搜索
  - 文本提取 / research retrieval
  - 相似页面发现
- 保留现有 Grok 能力：
  - 最新新闻研究
  - 社交和社区讨论研究
  - 广泛实时综合分析
  - 官方文档与社区解读对比
- 使用一个 TOML 配置文件：`~/.config/ai-skills/web-access.toml`。
- 默认输出机器可读 JSON，同时支持 plain text 和 URLs-only 格式。
- 增加项目作用域的 GitHub Actions，用于 PR 测试、tag-only release 构建和 skill 二进制更新。

## Primary Users / Roles

- AI coding agent：需要一个可靠的 Web 访问入口，用于文档查询、来源检索、实时更新和社区综合分析。
- Skill 维护者：需要一个自包含的二进制 skill，具备清晰的 provider 边界和项目作用域发布自动化。
- 用户：通过本地 skill 配置配置 Exa、Grok 或兼容代理 endpoint。

## Non-Goals

- 不把 Exa 和 Grok 合并成隐藏 provider 差异的不透明通用搜索模式。
- 第一版不删除 `exa-search/`、`grok-search/`、`exa-search-go/` 或 `grok-search-go/`。
- 不在 `web-access` 内保留旧二进制名称作为强兼容目标。
- 第一版不修改现有 `exa-search/SKILL.md` 或 `grok-search/SKILL.md`。
- 第一版不涉及旧 skill 的路由变更设计。
- 不创建全局 monorepo release workflow。
- 不增加缓存、本地索引、浏览器自动化或页面渲染。
- 第一版不增加 Grok cooldown 状态持久化。
- 第一版不把 AI 自动路由作为必需行为。

## Context

本仓库是一个 skill factory。每个 skill 都是根目录下的独立目录，包含 `SKILL.md` 入口，并可选包含 `bin/`、`references/`、`assets/` 和 `scripts/`。Go-based skill 使用混合架构：源码位于相邻的 Go 项目中，预编译平台二进制位于 skill 目录下。

`exa-search` 当前提供 source-first Web 检索，命令包括 `docs`、`search`、`research` 和 `similar`。它适合官方文档、API reference、价格页、权威来源检索和提取页面正文。它的 Go CLI 位于 `exa-search-go/`，使用 Cobra、Viper-style 配置加载、Exa API client、failover，以及 JSON/plain/URLs 输出。

`grok-search` 当前提供实时研究，命令包括 `news`、`social`、`research` 和 `docs-compare`。它适合最新更新、X/Twitter 或社区讨论、广泛实时综合分析，以及对比官方说法与社区解读。它的 Go CLI 位于 `grok-search-go/`，使用 Cobra、TOML 配置、OpenAI-compatible chat completions、failover、cooldown、prompt templates，以及 JSON/plain/URLs 输出。`web-access` 会保留 Grok failover 和 prompt 行为，但第一版不包含 cooldown 状态持久化。

## Discovery

### Key Discoveries

- 当前两个 skill 是互补关系，不是重复关系。Exa 是 source-first retrieval backend；Grok 是 live synthesis backend。
- 只做一个文档包装 skill 可以提供统一 skill 名称，但不能提供用户要求的统一二进制、统一配置和统一发布模型。
- 源码级合并比 wrapper skill 更符合本仓库现有 Go skill 模式。
- `grok-search-release.yml` 已经遵循较新的 `test -> build -> release` 发布结构。`exa-search-release.yml` 仍有较旧结构特征，所以 `web-access` 应遵循较新的仓库 workflow 要求，而不是照搬 Exa release 文件。
- Provider 差异应继续体现在文档和命令语义中，让 agent 能选择正确的检索模式。

### Scope Decisions

- 新增一等公民 `web-access` skill 和 `web-access-go` CLI，而不是重命名任一现有 skill。
- 默认通过命令语义选择 provider：
  - source-first 命令路由到 Exa
  - live synthesis 命令路由到 Grok
- 第一版保持 provider override 支持最小化。仅在有用且不含糊的地方添加显式 `--provider`。
- 旧 skill 完全不在本设计交付范围内；`web-access` 不改变它们的文档、路由元数据、源码或发布流程。
- 使用 `web-access-v*` release tag 和项目作用域 workflow，不提供手动 release dispatch。

## Decision Record

### Options Considered

- 只新增文档包装 skill：实现成本最低，但仍然保留两个二进制、两个配置和两套发布路径，无法满足统一操作模型。
- 重命名或替换某一个现有 skill：仓库形态更简单，但会模糊 Exa/Grok provider 边界，并给现有用户带来不必要的兼容风险。
- 新增一等公民 `web-access` skill 和 `web-access-go` CLI：实现成本更高，但可以在保留 provider 优势的同时，为 agent 提供一个入口、一个二进制、一个配置和一套发布流水线。

### Decision & Rationale

选择新增一等公民 `web-access` skill 和 `web-access-go` CLI。该方案将 source-first retrieval 和 live synthesis 保持为显式命令族，同时合并配置、输出处理和发布自动化。

已确认的补充决策：

- 第一版保留旧 skill 目录和源码项目不变。
- 将 Exa-backed 文本提取命令命名为 `extract`，从而把 `research` 名称留给 Grok-backed live synthesis。
- `web-access-release.yml` 使用 tag-only release workflow。
- 第一版省略 Grok cooldown 状态持久化，以保持配置和运行时行为更简单。

## Proposed Solution

创建 `web-access-go`，作为一个 Cobra-based Go CLI，并在内部使用 provider-specific package 分别封装 Exa 和 Grok。创建 `web-access` skill 目录，用于说明何时使用各个命令、发布预编译二进制，并提供配置参考。第一版优先保证命令名称清晰、provider 路由可预测，而不是做宽泛的自动路由。

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
    prompts/
```

CLI 共享命令解析、配置解析、输出渲染和 debug 处理。Provider-specific HTTP 请求构造、响应解析、failover 行为和 prompt 处理保持隔离，分别位于 `internal/providers/exa` 和 `internal/providers/grok`。

### Components

#### Skill 目录：`web-access/`

职责：

- 定义可触发的统一 Web access `SKILL.md`。
- 指导 agent 在 source-first retrieval 和 live synthesis 之间选择正确命令。
- 记录平台二进制选择方式。
- 提供 JSON、plain text 和 URLs-only 输出的简洁示例。
- 链接到配置和查询 recipe。

#### CLI 层：`web-access-go/cmd`

职责：

- 定义全局 flags：
  - `--config`
  - `--profile`
  - `--timeout`
  - `--plain`
  - `--urls`
  - `--json`
  - `--debug`
- 在需要时定义 provider-specific 全局 flags：
  - `--exa-api-key`
  - `--grok-api-key`
  - `--grok-model`
  - `--extra-body-json`
  - `--extra-headers-json`
- 定义命令：
  - `docs`
  - `search`
  - `extract`
  - `similar`
  - `news`
  - `social`
  - `research`
  - `docs-compare`
  - `version`

默认 provider 路由：

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

命令合约：

| 命令 | Provider | 必填输入 | 命令专属 flags | 默认值 | 输出 payload | 说明 |
|------|----------|----------|----------------|--------|--------------|------|
| `docs` | Exa | `--query` | `--num`, `--type`, `--text`, `--highlights`, `--start-date`, `--include-domains`, `--exclude-domains`, `--category`, `--no-autoprompt` | `--num 5`，`--type neural`，启用 autoprompt；省略 `--include-domains` 时默认包含 `docs.openclaw.ai`；text/highlights 默认关闭 | Exa `results`、`resolvedSearchType`、`requestId`、`searchTime`、`costDollars` | 面向官方文档和权威来源检索。 |
| `search` | Exa | `--query` | `--num`, `--type`, `--text`, `--highlights`, `--start-date`, `--include-domains`, `--exclude-domains`, `--category`, `--no-autoprompt` | `--num 5`，`--type neural`，启用 autoprompt；无默认 domain filter；text/highlights 默认关闭 | Exa `results`、`resolvedSearchType`、`requestId`、`searchTime`、`costDollars` | 面向通用 source-first 搜索。 |
| `extract` | Exa | `--query` | `--num`, `--type`, `--text`, `--highlights`, `--start-date`, `--include-domains`, `--exclude-domains`, `--category`, `--no-autoprompt` | `--num 5`，`--type neural`，启用 autoprompt；当 `--text` 和 `--highlights` 都未显式设置时，默认启用文本提取 | Exa `results`、`resolvedSearchType`、`requestId`、`searchTime`、`costDollars` | 面向需要页面正文或更深来源内容的检索。 |
| `similar` | Exa | `--url` | `--num` | `--num 5` | Exa `results`、`requestId`、`searchTime`、`costDollars` | 面向相似页面发现。 |
| `news` | Grok | `--query` | 仅使用全局 Grok flags | 使用 `news` prompt mode | Grok `content`、`sources`、`usage`、`raw` | 面向新鲜新闻和突发更新。 |
| `social` | Grok | `--query` | 仅使用全局 Grok flags | 使用 `social` prompt mode | Grok `content`、`sources`、`usage`、`raw` | 面向社交和社区讨论。 |
| `research` | Grok | `--query` | 仅使用全局 Grok flags | 使用 `research` prompt mode | Grok `content`、`sources`、`usage`、`raw` | 面向广泛实时综合分析。 |
| `docs-compare` | Grok | `--query` | 仅使用全局 Grok flags | 使用 `docs-compare` prompt mode | Grok `content`、`sources`、`usage`、`raw` | 面向官方信息与社区解读对比。 |
| `version` | local | 无 | 无 | 输出二进制版本元数据 | Plain command metadata | 本地版本信息命令。 |

#### 配置层：`internal/config`

职责：

- 默认加载 `~/.config/ai-skills/web-access.toml`。
- 默认路径配置不存在时自动创建模板配置。
- 合并 CLI flags、环境变量、TOML 配置和内置默认值。
- 解析独立的 Exa 和 Grok provider 配置。
- 支持 provider-specific profiles、base URLs、timeouts 和 Grok model overrides。
- 尽量保留现有环境变量兼容性：
  - `EXA_API_KEY`
  - `EXA_API_KEYS`
  - `EXA_BASE_URL`
  - `EXA_TIMEOUT`
  - `GROK_API_KEY`
  - `GROK_API_KEYS`
  - `GROK_BASE_URL`
  - `GROK_MODEL`
  - `GROK_TIMEOUT`
- 增加显式统一环境变量：
  - `WEB_ACCESS_EXA_API_KEY`
  - `WEB_ACCESS_GROK_API_KEY`
  - `WEB_ACCESS_CONFIG`

配置形状：

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
```

配置优先级从高到低：

1. CLI flags
2. `WEB_ACCESS_*` 环境变量
3. 现有 provider-specific 环境变量
4. TOML 配置
5. 内置默认值

#### Exa Provider：`internal/providers/exa`

职责：

- 实现由 Exa API 支撑的 source-first 命令。
- 支持当前 `exa-search-go` CLI 的 search request options：
  - 结果数量
  - search type
  - 文本提取
  - highlights
  - published date filters
  - include 和 exclude domains
  - category
  - autoprompt toggle
- 通过 Exa similar-page endpoint 实现 `similar`。
- 保留 profile failover 行为和标准化 attempt metadata。

#### Grok Provider：`internal/providers/grok`

职责：

- 实现由 OpenAI-compatible Grok chat completions 支撑的 live synthesis 命令。
- 保留 `news`、`social`、`research` 和 `docs-compare` 的 mode-specific prompts。
- 支持 model、base URL、timeout、extra body 和 extra header overrides。
- 保留 multi-profile failover 行为，但不写入 cooldown 状态。
- 在 provider 响应可用时保留 source extraction。

#### 输出层：`internal/output`

职责：

- 默认渲染 JSON。
- 为人工终端查看渲染 plain text。
- 为来源收集渲染 URLs-only 输出。
- 标准化 provider 之间的通用 metadata：
  - `ok`
  - `mode`
  - `provider`
  - `query`
  - `url`
  - `profileId`
  - `profileSource`
  - `attempts`
  - `elapsedMS`
- 保留 provider-specific result payload，不做有损转换：
  - Exa 命令返回 `results`
  - Grok 命令返回 `content`、`sources`、`usage` 和 `raw`

#### 发布自动化

新增：

- `.github/workflows/web-access-test.yml`
- `.github/workflows/web-access-release.yml`
- `.github/workflows/web-access-update-skill.yml`

Workflow 要求：

- Test workflow 只在 PR 触发，并使用 `web-access/**`、`web-access-go/**` 和 workflow 路径过滤。
- Release workflow 只在 `web-access-v*` tags 触发。
- Release workflow 在 build 前运行 Ubuntu、macOS 和 Windows matrix tests。
- Build 依赖 test；release 依赖 build。
- 使用 `CGO_ENABLED=0` 构建 Linux amd64/arm64、macOS amd64/arm64 和 Windows amd64。
- Update-skill workflow 在 release workflow 成功完成后运行，只更新 `web-access/bin/**` 和 checksums。

### Data Flow

#### Source-First Retrieval Path

1. Agent 调用 source-first 命令，例如：
   ```bash
   web-access docs --query "OpenClaw streaming API" --text --num 3
   ```
2. Cobra 解析命令并将其映射到 Exa provider。
3. Config 根据 CLI flags、环境变量、TOML 和默认值解析 Exa section。
4. Exa provider 构造 search 或 similar-page request。
5. Provider 使用配置的 profiles 发起请求，并在 rate limit、auth 或 quota 失败时 fail over。
6. Output 渲染标准化 JSON，包含 `provider: "exa"`、attempts 和 Exa result data。

#### Live Synthesis Path

1. Agent 调用 live 命令，例如：
   ```bash
   web-access news --query "latest model release updates" --plain
   ```
2. Cobra 解析命令并将其映射到 Grok provider。
3. Config 根据 CLI flags、环境变量、TOML 和默认值解析 Grok section。
4. Grok provider 选择 mode-specific prompt，并构造 OpenAI-compatible chat completion request。
5. Provider 使用配置的 profiles 发起请求，并在可重试的 profile 失败上 fail over。
6. Output 渲染标准化 JSON 或 plain text，包含 `provider: "grok"`、content、sources、usage 和 attempts。

#### Agent Selection Flow

1. 如果任务需要官方文档、API references、价格页、权威检索或提取页面正文，skill 指导 agent 使用 `docs`、`search`、`extract` 或 `similar`。
2. 如果任务依赖时效性、公开讨论、突发更新、X/Twitter 讨论或跨来源综合分析，skill 指导 agent 使用 `news`、`social`、`research` 或 `docs-compare`。
3. 如果两类能力都需要，skill 指导 agent 先运行 source-first 命令，再运行 live synthesis 命令进行解读或新鲜度补充。

## Error Handling

- 缺少 provider API key：
  - 返回结构化 JSON，包含 `ok: false`、`provider`、`error: "missing_api_key"` 和可执行的配置指引。
- 配置无效：
  - 返回 `config_parse_error`，包含配置路径和解析详情。
- 不支持的 command/provider 组合：
  - 在发起网络请求前返回 CLI validation error。
- Exa rate limit、quota 或 auth 失败：
  - Fail over 到下一个 Exa profile，并在输出中包含所有 attempts。
- Grok rate limit、quota 或 auth 失败：
  - Fail over 到下一个 Grok profile，并在输出中包含所有 attempts。不写入 cooldown 状态。
- 网络超时：
  - 返回 `request_failed`，包含 elapsed time、provider、profile attempts 和 timeout 值。
- Provider 响应为空或格式异常：
  - 返回 `response_parse_error`；仅在 debug mode 启用或已有安全输出中包含 raw response。
- 所有 profiles 不可用：
  - 返回 `all_profiles_failed`，并给出清晰的重试指引。

## Testing

- 单元测试配置优先级：
  - CLI flags
  - `WEB_ACCESS_*` 环境变量
  - provider-specific 环境变量
  - TOML 配置
  - 默认值
- 单元测试 command-to-provider 路由。
- 单元测试 `docs`、`search`、`extract` 和 `similar` 的 Exa request construction。
- 单元测试 `news`、`social`、`research` 和 `docs-compare` 的 Grok prompt selection 和 request construction。
- 单元测试两类 provider family 的 JSON 输出形状。
- 单元测试 Exa 和 Grok failover 行为。
- 单元测试命令专属 flag 默认值，尤其是 `extract` 默认文本提取和 `docs` 默认 domain filtering。
- 增加 `version` 和 help output 的 command smoke tests。
- 增加 release workflow 验证：
  - `go test`
  - formatting check
  - `go vet`
  - 在所有 release-test OS 上 build 并执行 `version`
- 在 `evals/web-access/` 下新增或更新 eval cases，覆盖代表性的 agent 路由决策和命令示例。

## Open Questions

无。关键决策已在 discovery 阶段确认。
