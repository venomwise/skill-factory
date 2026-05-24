# Quick Start: Creating a Release

## Prerequisites

- All tests passing in `exa-search-go/`
- Changes committed and pushed to main branch
- Write access to the repository

## Option 1: Manual Release (Recommended)

1. **Go to GitHub Actions**
   - Navigate to: https://github.com/venomwise/skill-factory/actions
   - Click on "Build and Release exa-search" workflow

2. **Run the workflow**
   - Click "Run workflow" button
   - Select branch: `main`
   - Enter version: `v1.0.0` (must start with `v`)
   - Click "Run workflow"

3. **Wait for completion**
   - Build takes ~5-10 minutes
   - Watch the progress in the Actions tab
   - All 5 platform builds must succeed

4. **Verify the release**
   - Go to: https://github.com/venomwise/skill-factory/releases
   - Find "Exa Search v1.0.0"
   - Check all 10 files are present (5 archives + 5 checksums)

## Option 2: Git Tag Release

```bash
# Create an annotated tag
git tag -a exa-search-v1.0.0 -m "Release exa-search v1.0.0"

# Push the tag
git push origin exa-search-v1.0.0

# The workflow will automatically trigger
```

## What Gets Built

For version `v1.0.0`, the following files are created:

```
exa-search-v1.0.0-linux-amd64.tar.gz       (Linux x86_64)
exa-search-v1.0.0-linux-amd64.tar.gz.sha256
exa-search-v1.0.0-linux-arm64.tar.gz       (Linux ARM64)
exa-search-v1.0.0-linux-arm64.tar.gz.sha256
exa-search-v1.0.0-darwin-amd64.tar.gz      (macOS Intel)
exa-search-v1.0.0-darwin-amd64.tar.gz.sha256
exa-search-v1.0.0-darwin-arm64.tar.gz      (macOS Apple Silicon)
exa-search-v1.0.0-darwin-arm64.tar.gz.sha256
exa-search-v1.0.0-windows-amd64.zip        (Windows x86_64)
exa-search-v1.0.0-windows-amd64.zip.sha256
```

## Testing a Release Locally

Before creating a public release, test the build locally:

```bash
cd exa-search-go

# Build with version info
VERSION="v1.0.0-dev"
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GO_VERSION=$(go version | awk '{print $3}')

go build -ldflags "\
  -s -w \
  -X main.version=${VERSION} \
  -X main.commit=${COMMIT} \
  -X main.date=${DATE} \
  -X main.goVersion=${GO_VERSION}" \
  -o exa-search cmd/exa-search/main.go

# Verify
./exa-search version
```

Expected output:
```
exa-search version v1.0.0-dev
commit: abc1234
built: 2026-05-24T22:15:00Z
go: go1.22.0
```

## Troubleshooting

### Workflow fails on build

Check the logs in GitHub Actions. Common issues:
- Go version mismatch
- Missing dependencies
- Syntax errors in code

### Release already exists

Delete the existing release and tag:
```bash
# Delete remote tag
git push --delete origin exa-search-v1.0.0

# Delete local tag
git tag -d exa-search-v1.0.0

# Delete release on GitHub (via web UI)
```

Then create the release again.

### Binary doesn't work on target platform

Ensure:
- `CGO_ENABLED=0` for static compilation
- Correct `GOOS` and `GOARCH` for target platform
- No platform-specific dependencies

## Version Numbering

Follow semantic versioning:
- `v1.0.0` - Major release (breaking changes)
- `v1.1.0` - Minor release (new features, backward compatible)
- `v1.0.1` - Patch release (bug fixes)

Examples:
- `v1.0.0` - Initial release
- `v1.1.0` - Add new search mode
- `v1.0.1` - Fix config parsing bug
- `v2.0.0` - Change CLI interface (breaking)

## Post-Release Checklist

- [ ] Verify all binaries download correctly
- [ ] Test binary on at least one platform
- [ ] Verify checksums match
- [ ] Update SKILL.md if installation instructions changed
- [ ] Announce in relevant channels
- [ ] Close related issues/PRs

## Next Steps

After your first release:
1. Users can download binaries from GitHub Releases
2. Update `exa-search/SKILL.md` to point to the latest release
3. Consider adding to package managers (Homebrew, Chocolatey, etc.)
