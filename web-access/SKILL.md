# Web Access Skill

Unified web access for AI coding agents — source-first search via Exa and live synthesis via Grok.

## When to Use

Use `web-access` when you need:

- **Official documentation lookup** (`docs`) — Find authoritative docs, API references, and official guides
- **Source-first search** (`search`) — Neural web search without domain restrictions
- **Text extraction** (`extract`) — Extract full text or highlights from search results
- **Similar pages** (`similar`) — Find pages similar to a reference URL
- **Fresh news** (`news`) — Get recent developments and breaking news
- **Social discourse** (`social`) — Analyze community discussions and social media trends
- **Broad research** (`research`) — Conduct comprehensive live research and synthesis
- **Docs comparison** (`docs-compare`) — Compare official docs with community interpretations

## Commands

### Exa Provider (Source-First)

- `web-access docs --query "openclaw documentation"` — Search official docs (defaults to docs.openclaw.ai)
- `web-access search --query "rust async patterns"` — Neural web search
- `web-access extract --query "machine learning" --text` — Extract full text content
- `web-access similar --url "https://example.com/article"` — Find similar pages

### Grok Provider (Live Synthesis)

- `web-access news --query "AI developments 2026"` — Fresh news and recent events
- `web-access social --query "developer opinions on rust"` — Social discourse analysis
- `web-access research --query "quantum computing applications"` — Broad live research
- `web-access docs-compare --query "react hooks official vs community"` — Compare docs with community insights

## Output Formats

- **JSON** (default): Structured output with metadata
  ```bash
  web-access docs --query "example"
  ```

- **Plain text**: Human-readable format
  ```bash
  web-access docs --query "example" --plain
  ```

- **URLs only**: Just the source URLs
  ```bash
  web-access docs --query "example" --urls
  ```

## Configuration

Default config: `~/.config/ai-skills/web-access.toml`

See `config.example.toml` and `references/configuration.md` for setup instructions.

## Platform Binaries

The `bin/` directory contains pre-built binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)  
- Windows (amd64)

Select the binary matching your platform. SHA256 checksums are in `bin/SHA256SUMS`.

## References

- [Configuration Guide](references/configuration.md)
- [Query Recipes](references/query-recipes.md)
