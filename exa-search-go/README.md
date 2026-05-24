# Exa Search Go

A command-line tool for neural web search using the Exa API, written in Go.

## Features

- **Multiple search modes**: general search, docs search, research mode, and similar pages
- **Multi-profile failover**: automatic retry with backup API keys on rate limits
- **Flexible configuration**: TOML files, environment variables, or CLI flags
- **Multiple output formats**: JSON, plain text, or URLs-only
- **Zero dependencies**: statically compiled binary with no runtime requirements

## Installation

### From Source

```bash
go build -o exa-search cmd/exa-search/main.go
```

### From Releases

Download pre-compiled binaries from [GitHub Releases](https://github.com/venomwise/skill-factory/releases).

## Configuration

The tool loads configuration from multiple sources with the following priority:

1. CLI flags (highest priority)
2. Environment variables
3. TOML config file (lowest priority)

### Config File

Create `~/.config/ai-skills/exa-search.toml`:

```toml
# API key profiles for failover support
[[profiles]]
id = "main"
api_key = "your-key-here"

[[profiles]]
id = "backup"
api_key = "your-backup-key"

# Global settings
timeout = 30  # Request timeout in seconds
```

The config file is auto-created with a template on first run.

### Environment Variables

```bash
export EXA_API_KEY="your-key"              # Single key
export EXA_API_KEYS="key1,key2,key3"       # Multiple keys for failover
export EXA_BASE_URL="https://api.exa.ai"   # Optional: override base URL
export EXA_TIMEOUT="30"                    # Optional: timeout in seconds
```

### CLI Flags

```bash
exa-search --api-key YOUR_KEY search --query "golang testing"
```

## Usage

### General Search

```bash
exa-search search --query "golang best practices" --num 10
```

### Documentation Search

Defaults to `docs.openclaw.ai` domain:

```bash
exa-search docs --query "authentication"
```

### Research Mode

Automatically extracts full page text:

```bash
exa-search research --query "machine learning papers" --num 5
```

### Find Similar Pages

```bash
exa-search similar --url "https://example.com/article"
```

### Output Formats

```bash
# JSON (default)
exa-search search --query "test"

# Plain text
exa-search search --query "test" --plain

# URLs only
exa-search search --query "test" --urls
```

### Debug Mode

```bash
exa-search search --query "test" --debug
```

Shows configuration resolution, HTTP requests/responses, and failover decisions.

## Commands

- `search` - General neural search
- `docs` - Search official documentation (defaults to docs.openclaw.ai)
- `research` - Deep research with text extraction
- `similar` - Find pages similar to a given URL
- `version` - Show version information

## Development

### Build

```bash
go build -o exa-search cmd/exa-search/main.go
```

### Build with Version Info

```bash
go build -ldflags "\
  -X main.version=1.0.0 \
  -X main.commit=$(git rev-parse --short HEAD) \
  -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X main.goVersion=$(go version | awk '{print $3}')" \
  -o exa-search cmd/exa-search/main.go
```

### Run Tests

```bash
go test ./...
```

## Architecture

The project follows a four-layer architecture:

1. **CLI Layer** (`cmd/`) - Command parsing and routing using cobra
2. **Config Layer** (`internal/config/`) - Configuration loading and merging using viper
3. **Client Layer** (`internal/client/`) - HTTP client with failover logic
4. **Output Layer** (`internal/output/`) - Formatted output rendering

## License

See the main repository for license information.
