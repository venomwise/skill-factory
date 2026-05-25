# Exa Search Binaries

This directory contains pre-compiled binaries for all supported platforms.

## Supported Platforms

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Files

After the first build, this directory will contain:

```
bin/
├── exa-search-linux-amd64
├── exa-search-linux-arm64
├── exa-search-darwin-amd64
├── exa-search-darwin-arm64
├── exa-search-windows-amd64.exe
└── SHA256SUMS
```

## Automatic Updates

These binaries are automatically built and updated by GitHub Actions when a new version is tagged.

See `.github/workflows/exa-search-update-skill.yml` for the automation workflow.

## Manual Build

To manually build and update binaries:

```bash
cd ../exa-search-go

# Build for all platforms
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o ../exa-search/bin/exa-search-linux-amd64 cmd/exa-search/main.go
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o ../exa-search/bin/exa-search-linux-arm64 cmd/exa-search/main.go
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o ../exa-search/bin/exa-search-darwin-amd64 cmd/exa-search/main.go
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o ../exa-search/bin/exa-search-darwin-arm64 cmd/exa-search/main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o ../exa-search/bin/exa-search-windows-amd64.exe cmd/exa-search/main.go

# Generate checksums
cd ../exa-search/bin
sha256sum exa-search-* > SHA256SUMS
```

## Size

Total size: ~50MB (all binaries combined)
- Each binary: ~10MB (statically compiled, stripped)
