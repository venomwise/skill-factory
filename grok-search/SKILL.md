---
name: grok-search
description: Real-time web research and live synthesis with sources. Use when the question depends on freshness, community chatter, X/Twitter dynamics, breaking updates, or broad multi-source summaries. Prefer this over source-first search when you need fast situational awareness. Prefer exa-search instead for official docs, API references, pricing pages, canonical source retrieval, or direct page-text extraction.
---

# Grok Search

Use Grok Search for **real-time research**.

Prefer it when the task is about:
- breaking news
- X/Twitter chatter
- fast-moving narratives
- “what are people saying now?”
- quick multi-source synthesis
- comparing official claims vs community discussion

Do **not** default to Grok Search for:
- official docs lookup
- API/reference pages
- pricing / plan details
- direct page text extraction
- canonical source retrieval

For those, prefer `exa-search`.

## Setup

The skill includes pre-compiled binaries for major platforms in `bin/`. Detect the user's platform and select the matching binary:

- Linux x86_64: `bin/grok-search-linux-amd64`
- Linux ARM64: `bin/grok-search-linux-arm64`
- macOS Intel: `bin/grok-search-darwin-amd64`
- macOS Apple Silicon: `bin/grok-search-darwin-arm64`
- Windows x86_64: `bin/grok-search-windows-amd64.exe`

Platform detection:

```bash
uname -s  # Linux, Darwin, or MINGW64_NT
uname -m  # x86_64, arm64, aarch64
```

## Config

The binary auto-creates a config template at `~/.config/ai-skills/grok-search.toml` on first run.

Recommended config:

```toml
base_url = "https://api.x.ai"
model = "grok-4.1-fast"
timeout = 120

[[profiles]]
id = "main"
api_key = "YOUR_GROK_API_KEY"
```

Profile-specific proxy endpoints are supported:

```toml
[[profiles]]
id = "proxy"
api_key = "YOUR_PROXY_KEY"
base_url = "https://your-compatible-endpoint.example"
model = "grok-custom-model"
```

Configuration priority:
1. CLI flags
2. Environment variables: `GROK_API_KEY`, `GROK_API_KEYS`, `GROK_BASE_URL`, `GROK_MODEL`, `GROK_TIMEOUT`
3. `~/.config/ai-skills/grok-search.toml`
4. Built-in defaults

See [configuration.md](references/configuration.md) for failover, cooldown, and advanced settings.

## Workflow

1. Use `news` for fresh updates and breaking news.
2. Use `social` for X/Twitter and discourse-heavy prompts.
3. Use `research` for broad multi-source synthesis.
4. Use `docs-compare` for official claims plus community interpretation.
5. Use `--plain` for human-readable terminal output.
6. Use `--urls` when only source URLs are needed.

## Commands

### Breaking news / fresh updates

```bash
bin/grok-search-<platform> news --query "Latest updates on X"
```

### Social and community discourse

```bash
bin/grok-search-<platform> social --query "What are people saying about OpenClaw on X?"
```

### Broad live synthesis

```bash
bin/grok-search-<platform> research --query "Summarize recent discussion around OpenClaw model failover"
```

### Official docs vs community interpretation

```bash
bin/grok-search-<platform> docs-compare --query "Compare official docs and community discussion on Telegram streaming"
```

### Output formats

```bash
bin/grok-search-<platform> research --query "OpenClaw Telegram streaming"
bin/grok-search-<platform> research --query "OpenClaw Telegram streaming" --plain
bin/grok-search-<platform> research --query "OpenClaw Telegram streaming" --urls
```

### Debugging

```bash
bin/grok-search-<platform> --debug research --query "test"
bin/grok-search-<platform> --profile main research --query "test" --plain
bin/grok-search-<platform> --ignore-cooldown news --query "OpenClaw updates" --plain
```

## Additional Resources

- **Source code**: See `grok-search-go/` for the Go implementation.
- **Configuration**: See [configuration.md](references/configuration.md).
- **Query examples**: See [query-recipes.md](references/query-recipes.md).
- **Migration**: See [migration-from-python.md](references/migration-from-python.md).

## Notes

- Optimized for real-time research and breadth, not canonical-source purity.
- `docs-compare` separates official facts from community interpretation.
- For official docs, API references, pricing pages, canonical retrieval, or direct page-text extraction, use `exa-search` instead.
- Zero Python runtime dependency: use the selected platform binary directly.
