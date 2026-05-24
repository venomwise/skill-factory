# Query Recipes

## Official Documentation

### Basic docs search

Search official documentation (defaults to docs.openclaw.ai):

```bash
exa-search docs --query "telegram streaming openclaw"
```

### Docs with text extraction

Extract full text from documentation pages:

```bash
exa-search docs --query "model failover openclaw" --text --num 2
```

### Specific domain documentation

Search documentation on specific domains:

```bash
exa-search search \
  --query "OpenClaw pricing API parameters" \
  --include-domains docs.openclaw.ai,openclaw.ai \
  --text --num 3
```

## Product Research

### Company overview

Deep extraction mode for product/company research:

```bash
exa-search research --query "Exa AI company overview" --num 3
```

### Find similar pages

Expand from one canonical source page:

```bash
exa-search similar --url "https://docs.openclaw.ai/channels/telegram" --num 5
```

**Warning**: `similar` uses semantic similarity, not official-source-only discovery. For official docs, prefer `docs` search with strict domain restriction.

## Advanced Filtering

### Freshness filter

Filter by publication date:

```bash
exa-search search \
  --query "OpenClaw releases" \
  --start-date 2026-01-01 \
  --num 5
```

### Domain hygiene

Use `--include-domains` when the user mentions:
- 官方文档 (official docs)
- 官网 (official website)
- API 文档 (API documentation)
- 价格页 (pricing page)
- 参数说明 (parameter documentation)

Examples:
```bash
# Official docs only
exa-search search \
  --query "OpenClaw API" \
  --include-domains docs.openclaw.ai

# Official site + docs
exa-search search \
  --query "OpenClaw pricing" \
  --include-domains openclaw.ai,docs.openclaw.ai

# GitHub repositories
exa-search search \
  --query "OpenClaw examples" \
  --include-domains github.com
```

### Exclude domains

Filter out unwanted sources:

```bash
exa-search search \
  --query "Python async programming" \
  --exclude-domains stackoverflow.com,reddit.com \
  --num 5
```

## Output Formats

### JSON output (default)

Structured JSON for programmatic consumption:

```bash
exa-search docs --query "telegram streaming"
```

### Human-readable output

Plain text format:

```bash
exa-search docs --query "telegram streaming" --plain
```

### URLs only

Extract just the URLs:

```bash
exa-search docs --query "telegram streaming" --urls
```

## Search Types

### Neural search (default)

Semantic understanding of query intent:

```bash
exa-search search --query "how to implement real-time streaming" --type neural
```

### Keyword search

Traditional keyword matching:

```bash
exa-search search --query "telegram streaming API" --type keyword
```

### Magic search

Exa's automatic type selection:

```bash
exa-search search --query "OpenClaw documentation" --type magic
```

## Common Patterns

### Official docs with highlights

Get key excerpts from documentation:

```bash
exa-search docs --query "authentication methods" --highlights --num 3
```

### Recent articles only

Find fresh content:

```bash
exa-search search \
  --query "AI coding assistants 2026" \
  --start-date 2026-01-01 \
  --num 10
```

### Category-specific search

Search within specific content categories:

```bash
exa-search search \
  --query "machine learning" \
  --category "research paper" \
  --num 5
```

Available categories: `company`, `research paper`, `news`, `github`, `tweet`, `movie`, `song`, `personal site`, `pdf`

## Debug Mode

Enable verbose logging to see what's happening:

```bash
exa-search --debug search --query "test"
```

Shows:
- Configuration resolution
- HTTP requests/responses
- Failover decisions
- API key redaction (first 8 chars visible)

## Tips

1. **Start with `docs`**: For official documentation, always try `docs` first
2. **Use domain filters**: Restrict to official sources when accuracy matters
3. **Extract text**: Add `--text` when you need content, not just links
4. **Limit results**: Use `--num` to control result count (default: 5)
5. **Check output**: Use `--plain` for quick human review, JSON for automation
6. **Disable autoprompt**: Add `--no-autoprompt` if Exa's query rewriting interferes
