---
name: exa-search
description: >
  Neural web search for official documentation, API references, and structured source retrieval with text extraction.
  Use for precise, source-first research when you need high-quality results or extracted page content.
  For breaking news or real-time sentiment, use grok-search instead.
---

# Exa Search

Use Exa for **source-first retrieval**.

Prefer it when the task is about:
- official documentation
- API/reference pages
- pricing/plan details
- product/company pages
- extracting the text of a page instead of just finding the link
- expanding from one canonical page to similar pages

Do **not** default to Exa for:
- breaking news
- X/Twitter chatter
- live sentiment / fast-moving discourse
- broad real-time synthesis across many fresh sources

For those, prefer `grok-search`.

## Setup

**1. Configure API key**

The tool auto-creates a config template at `~/.config/ai-skills/exa-search.toml` on first run.

Edit the file and add your API key:

```toml
[[profiles]]
id = "main"
api_key = "YOUR_EXA_API_KEY"
```

Alternatively, use environment variables:

```bash
export EXA_API_KEY="your-key-here"
```

Or pass via CLI flag:

```bash
bin/exa-search-<platform> --api-key YOUR_KEY search --query "test"
```

For multiple keys, advanced configuration, and failover setup, see [CONFIG.md](references/CONFIG.md).

**Note:** The skill includes pre-compiled binaries for all major platforms in `bin/`. Detect the user's platform and select the appropriate binary:
- Linux x86_64: `bin/exa-search-linux-amd64`
- Linux ARM64: `bin/exa-search-linux-arm64`
- macOS Intel: `bin/exa-search-darwin-amd64`
- macOS Apple Silicon: `bin/exa-search-darwin-arm64`
- Windows x86_64: `bin/exa-search-windows-amd64.exe`

## Workflow

1. Start with `docs` for official documentation lookups.
2. Use `search --text` or `research` when you need extracted body text.
3. Restrict domains aggressively when the user wants official sources.
4. Use `similar` when you already have the best canonical page and want adjacent sources.
5. For official-doc-only work, prefer `docs` plus domain restriction over `similar`; `similar` is semantic, not source-pure.
6. Return links plus extracted evidence, not just titles.

## Usage Pattern

**Platform Detection:**
Before executing, detect the user's platform:
```bash
# Detect OS
uname -s  # Linux, Darwin, or MINGW64_NT (Windows)

# Detect architecture
uname -m  # x86_64, arm64, aarch64
```

Then select the appropriate binary from `bin/` and execute directly.

## Commands

### Official docs search
```bash
bin/exa-search-<platform> docs --query "telegram streaming openclaw"
```

### Official docs with text extraction
```bash
bin/exa-search-<platform> docs --query "model failover openclaw" --text --num 2
```

### General source-first search
```bash
bin/exa-search-<platform> search --query "OpenClaw Telegram streaming" --num 5
```

### Deep extraction / research
```bash
bin/exa-search-<platform> research --query "OpenClaw model failover" --num 3
```

### Find similar pages
```bash
bin/exa-search-<platform> similar --url "https://docs.openclaw.ai/channels/telegram" --num 5
```

### Output formats

```bash
# JSON (default)
bin/exa-search-<platform> search --query "test"

# Plain text
bin/exa-search-<platform> search --query "test" --plain

# URLs only
bin/exa-search-<platform> search --query "test" --urls
```

### Debug mode

```bash
bin/exa-search-<platform> search --query "test" --debug
```

## Additional Resources

- **Source code**: See `exa-search-go/` for the Go implementation
- **Query examples**: See [query-recipes.md](references/query-recipes.md) for ready-made patterns
- **Configuration**: See [CONFIG.md](references/CONFIG.md) for advanced setup and failover
- **Evaluations**: See `evals/exa-search/test_cases.json` for test scenarios

## Notes

- `docs` defaults to `includeDomains=docs.openclaw.ai`
- `research` defaults to text extraction
- Output is normalized JSON for reliable consumption
- Supports automatic failover across multiple API keys
- Zero runtime dependencies - statically compiled Go binaries included for all platforms
- Detect user's platform (OS + architecture) and select the appropriate binary from `bin/`
- All binaries support the same command-line interface
