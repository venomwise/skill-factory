# Configuration Guide

## Configuration Priority

The script resolves configuration in this order (first match wins):

1. `--api-key` command-line argument
2. `EXA_API_KEY` environment variable
3. `EXA_API_KEYS` environment variable (comma-separated)
4. `config.local.json` (in skill directory)
5. `config.json` (in skill directory)
6. `~/.codex/config/exa-search.json` (global fallback)

## Single API Key

Create `config.local.json`:

```json
{
  "profiles": [
    { "id": "main", "api_key": "YOUR_EXA_API_KEY" }
  ],
  "base_url": "https://api.exa.ai",
  "timeout_seconds": 30
}
```

## Multiple Keys with Auto-Failover

For production use or high-volume queries, configure multiple API keys:

```json
{
  "profiles": [
    { "id": "main", "api_key": "KEY_1" },
    { "id": "backup-1", "api_key": "KEY_2" },
    { "id": "backup-2", "api_key": "KEY_3" }
  ],
  "base_url": "https://api.exa.ai",
  "timeout_seconds": 30
}
```

### Failover Behavior

The script automatically tries the next profile when:
- HTTP 401 (Unauthorized)
- HTTP 403 (Forbidden)
- HTTP 429 (Too Many Requests)
- Response contains: "rate limit", "quota", "credits", "insufficient", "billing", etc.

The output includes `profileId` and `attempts` so you can see which key was used and which failed.

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
export EXA_TIMEOUT_SECONDS="60"
```

## Force Specific Profile

Use `--profile` to force a specific configured profile:

```bash
python .venv/bin/python scripts/exa_search.py docs \
  --query "telegram streaming" \
  --profile backup-1
```

## Configuration File Format

### Full Schema

```json
{
  "profiles": [
    {
      "id": "profile-name",
      "api_key": "YOUR_API_KEY",
      "base_url": "https://api.exa.ai",
      "enabled": true
    }
  ],
  "base_url": "https://api.exa.ai",
  "timeout_seconds": 30
}
```

### Profile Fields

- `id` (string): Profile identifier for `--profile` flag
- `api_key` (string, required): Exa API key
- `base_url` (string, optional): Override base URL for this profile
- `enabled` (boolean, optional): Set to `false` to disable a profile

### Legacy Formats

The script also supports legacy configuration formats for backward compatibility:

**Simple key list:**
```json
{
  "api_keys": ["KEY_1", "KEY_2", "KEY_3"]
}
```

**Single key:**
```json
{
  "api_key": "YOUR_API_KEY",
  "profile_id": "main"
}
```

## Best Practices

1. **Use `config.local.json`**: Keep API keys in `config.local.json` (gitignored) rather than `config.json`
2. **Multiple keys for production**: Configure 2-3 keys for automatic failover
3. **Monitor usage**: Check the `attempts` field in output to see failover activity
4. **Set appropriate timeout**: Increase `timeout_seconds` for slow networks or large extractions
5. **Profile naming**: Use descriptive profile IDs like "main", "backup", "high-quota"

## Troubleshooting

### "missing_api_key" error

No valid API key found. Solutions:
- Pass `--api-key YOUR_KEY`
- Set `EXA_API_KEY` environment variable
- Create `config.local.json` with a profile

### "all_profiles_failed" error

All configured keys failed with failover-triggering errors. Check:
- API key validity
- Account quota/credits
- Rate limits
- Network connectivity

### Profile not found

When using `--profile`, ensure the profile ID matches exactly:
```bash
# Check configured profiles in output
python .venv/bin/python scripts/exa_search.py docs --query "test"
```
