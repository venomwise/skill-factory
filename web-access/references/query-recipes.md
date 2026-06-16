# Query Recipes

Practical command patterns for common web access tasks.

## Source-First Commands (Exa)

### Official Documentation Lookup

```bash
# Search OpenClaw docs (default domain filter)
web-access docs --query "authentication api"

# Search specific documentation sites
web-access docs --query "react hooks" --include-domains "react.dev,legacy.reactjs.org"

# Get full text from docs
web-access docs --query "kubernetes networking" --text
```

### General Web Search

```bash
# Neural search without domain restrictions
web-access search --query "rust async patterns"

# Keyword search with date filter
web-access search --query "machine learning trends" --type keyword --start-date 2026-01-01

# Search with category filter
web-access search --query "climate change" --category news
```

### Text Extraction

```bash
# Extract full text (default behavior when neither --text nor --highlights is set)
web-access extract --query "python async tutorial"

# Extract with highlights only
web-access extract --query "golang best practices" --highlights

# Extract both text and highlights
web-access extract --query "database indexing" --text --highlights
```

### Similar Pages

```bash
# Find similar pages to a reference URL
web-access similar --url "https://example.com/article"

# Find more similar pages
web-access similar --url "https://blog.example.com/post" --num 10
```

## Live Synthesis Commands (Grok)

### Fresh News

```bash
# Get recent developments
web-access news --query "AI developments June 2026"

# Breaking news
web-access news --query "tech industry layoffs"
```

### Social Discourse

```bash
# Analyze community discussions
web-access social --query "developer opinions on rust vs go"

# Social media trends
web-access social --query "reactions to GPT-5 release"
```

### Broad Research

```bash
# Comprehensive research synthesis
web-access research --query "quantum computing practical applications"

# Multi-faceted topic exploration
web-access research --query "remote work productivity studies"
```

### Documentation Comparison

```bash
# Compare official docs with community insights
web-access docs-compare --query "react hooks official vs community best practices"

# Official vs real-world usage
web-access docs-compare --query "kubernetes official recommendations vs production patterns"
```

## Output Format Examples

### JSON Output (default)

```bash
web-access docs --query "example"
# Returns structured JSON with metadata, results/content, sources, usage
```

### Plain Text Output

```bash
web-access docs --query "example" --plain
# Human-readable format for terminal viewing
```

### URLs Only

```bash
web-access docs --query "example" --urls
# Just the source URLs, one per line
```

## Advanced Patterns

### Domain Filtering

```bash
# Include specific domains
web-access search --query "react tutorial" --include-domains "react.dev,github.com"

# Exclude domains
web-access search --query "javascript guides" --exclude-domains "w3schools.com"
```

### Profile Selection

```bash
# Use specific profile for failover
web-access docs --query "example" --profile backup
```

### Custom Grok Model

```bash
# Override Grok model
web-access news --query "AI news" --grok-model "grok-2"
```

### Debug Mode

```bash
# Enable debug logging
web-access docs --query "example" --debug
```

## Tips

- **`docs`** defaults to `docs.openclaw.ai` — override with `--include-domains` for other doc sites
- **`extract`** defaults to text extraction when neither `--text` nor `--highlights` is explicitly set
- **`search`** has no default domain filtering — use `--include-domains` to restrict
- **Grok commands** require `--query` and use mode-specific prompts internally
- Use `--urls` output when you only need source links for further processing
