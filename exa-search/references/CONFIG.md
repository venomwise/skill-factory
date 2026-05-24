# Configuration Guide

## Configuration Priority

The tool resolves configuration in this order (first match wins):

1. `--api-key` command-line flag
2. `EXA_API_KEY` environment variable
3. `EXA_API_KEYS` environment variable (comma-separated)
4. `~/.config/ai-skills/exa-search.toml` (TOML config file)

## Single API Key

### Using TOML Config File

The config file is auto-created at `~/.config/ai-skills/exa-search.toml` on first run.

Edit it to add your API key:

```toml
[[profiles]]
id = "main"
api_key = "YOUR_EXA_API_KEY"

# Global settings
timeout = 30  # Request timeout in seconds
```

### Using Environment Variable

```bash
export EXA_API_KEY="your_api_key"
```

### Using CLI Flag

```bash
exa-search --api-key YOUR_KEY search --query "test"
```

## Multiple Keys with Auto-Failover

For production use or high-volume queries, configure multiple API keys in `~/.config/ai-skills/exa-search.toml`:

```toml
[[profiles]]
id = "main"
api_key = "KEY_1"

[[profiles]]
id = "backup-1"
api_key = "KEY_2"

[[profiles]]
id = "backup-2"
api_key = "KEY_3"

# Global settings
timeout = 30
```

Or use environment variable:

```bash
export EXA_API_KEYS="key1,key2,key3"
```

### Failover Behavior

The tool automatically tries the next profile when:
- HTTP 401 (Unauthorized)
- HTTP 403 (Forbidden)
- HTTP 429 (Too Many Requests)
- Response contains: "rate limit", "quota", "credits"

The output includes `profileId` and `attempts` so you can see which key was used and which failed.

## Force Specific Profile

Use `--profile` to force a specific configured profile:

```bash
exa-search --profile backup-1 docs --query "telegram streaming"
```

## Configuration File Format

### Full TOML Schema

```toml
# API key profiles for failover support
[[profiles]]
id = "main"
api_key = "YOUR_API_KEY"
# base_url = "https://api.exa.ai"  # Optional: override base URL for this profile

[[profiles]]
id = "backup"
api_key = "YOUR_BACKUP_KEY"

# Global settings
base_url = "https://api.exa.ai"  # Default API endpoint
timeout = 30  # Request timeout in seconds
```

### Profile Fields

- `id` (string): Profile identifier for `--profile` flag
- `api_key` (string, required): Exa API key
- `base_url` (string, optional): Override base URL for this profile

### Global Settings

- `base_url` (string, optional): Default API endpoint (default: https://api.exa.ai)
- `timeout` (int, optional): Request timeout in seconds (default: 30)

## Environment Variables

### Single key
```bash
export EXA_API_KEY="your_api_key"
```

### Multiple keys
```bash
export EXA_API_KEYS="key1,key2,key3"
```

### Custom base URL
```bash
export EXA_BASE_URL="https://custom.api.endpoint"
```

### Custom timeout
```bash
export EXA_TIMEOUT="60"
```

## Best Practices

1. **Use `~/.config/ai-skills/exa-search.toml` for shared config**: This is the recommended location for API keys that should be shared across all AI agents
2. **Multiple keys for production**: Configure 2-3 keys for automatic failover
3. **Monitor usage**: Check the `attempts` field in output to see failover activity
4. **Set appropriate timeout**: Increase `timeout` for slow networks or large extractions
5. **Profile naming**: Use descriptive profile IDs like "main", "backup", "high-quota"

## Configuration Locations

### Recommended
`~/.config/ai-skills/exa-search.toml` - TOML config file (auto-created on first run)

### Alternative
- `--config /path/to/config.toml` - Custom config file location
- Environment variables - For CI/CD or temporary overrides

## Troubleshooting

### "missing_api_key" error

No valid API key found. Solutions:
- Pass `--api-key YOUR_KEY`
- Set `EXA_API_KEY` environment variable
- Edit `~/.config/ai-skills/exa-search.toml` and add your API key

### "all_profiles_failed" error

All configured keys failed with failover-triggering errors. Check:
- API key validity
- Account quota/credits
- Rate limits
- Network connectivity

### Profile not found

When using `--profile`, ensure the profile ID matches exactly:
```bash
# Check configured profiles with debug mode
exa-search --debug search --query "test"
```

### Invalid TOML syntax

If you see a config parse error, check your TOML syntax:
- Strings must be quoted: `api_key = "YOUR_KEY"`
- Arrays use double brackets: `[[profiles]]`
- Comments start with `#`

## Debug Mode

Enable debug logging to see configuration resolution:

```bash
exa-search --debug search --query "test"
```

This shows:
- Which config file was loaded
- How many profiles were found
- Which profile is being used for each request
- API key redaction (first 8 chars visible)
- HTTP request/response details
- Failover decisions
