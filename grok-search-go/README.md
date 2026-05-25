# Grok Search Go

A command-line tool for real-time web research using an OpenAI-compatible Grok endpoint, written in Go.

## Features

- Research modes: news, social, research, and docs-compare
- Multi-profile failover with cooldown support
- TOML configuration, environment variables, and CLI overrides
- JSON, plain text, and URLs-only output
- Statically compiled binaries with no Python runtime dependency

## Installation

### From Source

```bash
go build -o grok-search ./cmd/grok-search
```

### From Releases

Download pre-compiled binaries from the repository releases.

## Configuration

The default config path is:

```text
~/.config/ai-skills/grok-search.toml
```

The tool defaults to:

```toml
base_url = "https://api.x.ai"
model = "grok-4.1-fast"
```

Profiles can override `base_url` and `model` for compatible proxy endpoints.

## Usage

```bash
grok-search news --query "latest updates"
grok-search social --query "what are people saying now?"
grok-search research --query "summarize recent discussion"
grok-search docs-compare --query "compare official docs and community interpretation"
grok-search version
```

## Development

```bash
go test ./...
go vet ./...
go build -o grok-search ./cmd/grok-search
```

## Architecture

The project follows the same delivery model as `exa-search-go`: CLI, config, client, cooldown, output, prompts, and debug packages.
