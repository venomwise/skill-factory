# GitHub Actions Workflows for exa-search

This directory contains automated workflows for building, testing, and releasing the exa-search Go CLI tool.

## Workflows

### 1. Test Workflow (`exa-search-test.yml`)

**Triggers:**
- Push to `main` or `develop` branches (when exa-search-go files change)
- Pull requests (when exa-search-go files change)

**What it does:**
- Runs tests on Linux, macOS, and Windows
- Tests with Go 1.22 and 1.23
- Checks code formatting with `gofmt`
- Runs `go vet` for static analysis
- Runs `golangci-lint` for comprehensive linting
- Generates code coverage reports
- Builds and verifies the binary

**Usage:**
Automatically runs on every push and PR. No manual action needed.

### 2. Build and Release Workflow (`exa-search-release.yml`)

**Triggers:**
- Manual workflow dispatch (with version input)
- Push tags matching `exa-search-v*` (e.g., `exa-search-v1.0.0`)

**What it does:**
- Builds binaries for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
- Injects version information via `-ldflags`
- Creates compressed archives (`.tar.gz` for Unix, `.zip` for Windows)
- Generates SHA256 checksums for verification
- Creates a GitHub Release with all artifacts
- Includes installation instructions in release notes

**Usage:**

#### Option 1: Manual Trigger

1. Go to Actions tab in GitHub
2. Select "Build and Release exa-search"
3. Click "Run workflow"
4. Enter version (e.g., `v1.0.0`)
5. Click "Run workflow"

#### Option 2: Git Tag

```bash
# Create and push a tag
git tag exa-search-v1.0.0
git push origin exa-search-v1.0.0
```

## Release Process

### Step-by-step Guide

1. **Ensure all tests pass**
   ```bash
   cd exa-search-go
   go test ./...
   go vet ./...
   ```

2. **Update version documentation** (if needed)
   - Update CHANGELOG.md
   - Update README.md if there are breaking changes

3. **Create a release** using one of these methods:

   **Method A: Manual workflow dispatch**
   - Go to GitHub Actions
   - Run "Build and Release exa-search" workflow
   - Enter version: `v1.0.0`

   **Method B: Git tag**
   ```bash
   git tag -a exa-search-v1.0.0 -m "Release version 1.0.0"
   git push origin exa-search-v1.0.0
   ```

4. **Verify the release**
   - Check GitHub Releases page
   - Download and test binaries for your platform
   - Verify checksums

5. **Announce the release**
   - Update skill documentation if needed
   - Notify users of new features/fixes

## Build Artifacts

Each release includes:

```
exa-search-v1.0.0-linux-amd64.tar.gz
exa-search-v1.0.0-linux-amd64.tar.gz.sha256
exa-search-v1.0.0-linux-arm64.tar.gz
exa-search-v1.0.0-linux-arm64.tar.gz.sha256
exa-search-v1.0.0-darwin-amd64.tar.gz
exa-search-v1.0.0-darwin-amd64.tar.gz.sha256
exa-search-v1.0.0-darwin-arm64.tar.gz
exa-search-v1.0.0-darwin-arm64.tar.gz.sha256
exa-search-v1.0.0-windows-amd64.zip
exa-search-v1.0.0-windows-amd64.zip.sha256
```

## Version Injection

The build process injects version information at compile time:

```bash
go build -ldflags "\
  -X main.version=v1.0.0 \
  -X main.commit=abc12345 \
  -X main.date=2026-05-24T22:00:00Z \
  -X main.goVersion=go1.22.0"
```

This information is displayed by `exa-search version`:

```
exa-search version v1.0.0
commit: abc12345
built: 2026-05-24T22:00:00Z
go: go1.22.0
```

## Local Testing

To test the build process locally:

```bash
cd exa-search-go

# Build for current platform
go build -o exa-search cmd/exa-search/main.go

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o exa-search-linux-amd64 cmd/exa-search/main.go
GOOS=darwin GOARCH=arm64 go build -o exa-search-darwin-arm64 cmd/exa-search/main.go
GOOS=windows GOARCH=amd64 go build -o exa-search-windows-amd64.exe cmd/exa-search/main.go

# Build with version info
go build -ldflags "\
  -X main.version=v1.0.0-dev \
  -X main.commit=$(git rev-parse --short HEAD) \
  -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X main.goVersion=$(go version | awk '{print $3}')" \
  -o exa-search cmd/exa-search/main.go

# Verify
./exa-search version
```

## Troubleshooting

### Build fails on specific platform

Check the build logs in GitHub Actions. Common issues:
- Missing dependencies in `go.mod`
- Platform-specific code without build tags
- CGO dependencies (should be disabled with `CGO_ENABLED=0`)

### Release creation fails

Ensure:
- You have write permissions to the repository
- The tag follows the pattern `exa-search-v*`
- The version string is valid (e.g., `v1.0.0`, not `1.0.0`)

### Checksums don't match

This usually indicates:
- Download was corrupted
- Wrong file was downloaded
- File was modified after download

Always verify checksums:
```bash
# Linux/macOS
sha256sum -c exa-search-v1.0.0-linux-amd64.tar.gz.sha256

# Windows (PowerShell)
Get-FileHash exa-search-v1.0.0-windows-amd64.zip -Algorithm SHA256
```

## Security

- All binaries are built in GitHub's secure runners
- No secrets are embedded in binaries
- Checksums allow verification of download integrity
- Static compilation reduces supply chain attack surface

## Future Improvements

- [ ] Add automated changelog generation
- [ ] Add Docker image builds
- [ ] Add Homebrew formula auto-update
- [ ] Add Windows installer (MSI)
- [ ] Add Linux package repositories (deb, rpm)
- [ ] Add code signing for macOS and Windows binaries
