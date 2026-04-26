# Configuration Guide

## Cross-Agent Setup

This skill is installed in `~/.AI-Skills/grok-search/` and linked to multiple agents (pi, codex, etc.).

**Shared configuration**: `~/.config/ai-skills/grok-search.json` (recommended for cross-agent use)
**Skill-specific**: `~/.AI-Skills/grok-search/config.local.json` (if you need different keys per skill instance)

## Key Resolution Order

The script checks for API keys in this order:

1. `--api-key` command-line argument
2. `GROK_API_KEY` environment variable
3. `GROK_API_KEYS` environment variable (comma-separated)
4. `~/.config/ai-skills/grok-search.json` (cross-agent shared)
5. `config.local.json` in skill directory
6. `config.json` in skill directory

**Recommended**: Use `~/.config/ai-skills/grok-search.json` so all agents (pi, codex, etc.) share the same keys.

## Single Key Setup

Create `~/.config/ai-skills/grok-search.json`:

```json
{
  "profiles": [
    { "id": "main", "api_key": "YOUR_GROK_API_KEY" }
  ]
}
```

## Multiple Keys with Auto Failover

Create `config.local.json`:

```json
{
  "profiles": [
    { "id": "main", "api_key": "KEY_1" },
    { "id": "backup-1", "api_key": "KEY_2" },
    { "id": "backup-2", "api_key": "KEY_3" }
  ],
  "base_url": "https://your-grok-endpoint.example",
  "model": "grok-4.1-fast",
  "timeout_seconds": 120,
  "extra_body": {},
  "extra_headers": {},
  "cooldown": {
    "enabled": true,
    "state_file": "runtime/cooldowns.json",
    "default_minutes": 15,
    "rate_limit_minutes": 20,
    "quota_minutes": 60,
    "auth_minutes": 360
  }
}
```

### Failover Behavior

- Profiles are tried in order
- 401 / 403 / 429 errors automatically move to the next key
- Quota, billing, rate-limit, and token-unavailable errors trigger failover
- Output includes `profileId` and `attempts` for debugging

### Cooldown Behavior

- Failover-worthy failures place the profile into temporary cooldown
- Cooldown state is stored in `runtime/cooldowns.json` by default
- Later runs skip cooling profiles instead of retrying the same failing key
- Use `--ignore-cooldown` to force a retry on a cooling profile

### Cooldown Duration by Error Type

- Rate limit: 20 minutes
- Quota exceeded: 60 minutes
- Auth errors: 360 minutes (6 hours)
- Other errors: 15 minutes (default)

## Configuration Fields

### Required

- `profiles`: Array of profile objects with `id` and `api_key`

### Optional

- `base_url`: Custom API endpoint (default: Grok's official endpoint)
- `model`: Model name (default: `grok-4.1-fast`)
- `timeout_seconds`: Request timeout (default: 120)
- `extra_body`: Additional JSON fields to include in API request body
- `extra_headers`: Additional HTTP headers
- `cooldown`: Cooldown configuration object (see above)

## Testing a Specific Profile

Force a specific profile for debugging:

```bash
python3 scripts/grok_search.py --query "test query" --profile main --plain
```

This bypasses failover and uses only the specified profile.
