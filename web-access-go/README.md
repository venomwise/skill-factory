# web-access-go

Go implementation of the `web-access` unified Web access CLI for AI coding agents.

## Local Development

### Build

```bash
go build -o web-access ./cmd/web-access
```

### Build with version information

```bash
go build -ldflags "-X main.version=1.0.0 -X main.commit=$(git rev-parse HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.goVersion=$(go version | awk '{print $3}')" -o web-access ./cmd/web-access
```

### Test

```bash
go test ./...
```

### Run locally

```bash
./web-access version
./web-access docs --query "openclaw documentation"
./web-access search --query "example query"
./web-access news --query "latest AI developments"
```

### Format and vet

```bash
go fmt ./...
go vet ./...
```

## Configuration

Default config path: `~/.config/ai-skills/web-access.toml`

See `../web-access/config.example.toml` for configuration examples.

## Project Structure

```
web-access-go/
├── cmd/
│   ├── web-access/          # Main entry point
│   │   └── main.go
│   ├── root.go              # Root command and global flags
│   ├── docs.go              # Exa: official docs lookup
│   ├── search.go            # Exa: source-first search
│   ├── extract.go           # Exa: text extraction
│   ├── similar.go           # Exa: similar pages
│   ├── news.go              # Grok: fresh news
│   ├── social.go            # Grok: social discourse
│   ├── research.go          # Grok: broad research
│   ├── docs_compare.go      # Grok: docs comparison
│   └── version.go           # Version command
├── internal/
│   ├── config/              # Configuration and profile resolution
│   ├── providers/
│   │   ├── exa/            # Exa provider implementation
│   │   └── grok/           # Grok provider implementation
│   ├── prompts/            # Grok prompt registry
│   ├── output/             # Output renderers and error helpers
│   └── debug/              # Debug logger
└── go.mod
```
