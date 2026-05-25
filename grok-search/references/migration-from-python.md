# Migration from Python Script

`grok-search` now ships as prebuilt Go binaries. The old Python script and JSON config are no longer supported interfaces.

## Command Migration

Old Python command shape:

```bash
python3 scripts/grok_search.py --mode news --query "Latest updates"
```

New binary command shape:

```bash
bin/grok-search-<platform> news --query "Latest updates"
```

### Mode Mapping

| Old `--mode` | New subcommand |
| --- | --- |
| `--mode news` | `news` |
| `--mode social` | `social` |
| `--mode research` | `research` |
| `--mode docs-compare` | `docs-compare` |

Examples:

```bash
# Old
python3 scripts/grok_search.py --mode social --query "What are people saying?"

# New
bin/grok-search-<platform> social --query "What are people saying?"
```

```bash
# Old
python3 scripts/grok_search.py --mode docs-compare --query "Compare docs and community discussion"

# New
bin/grok-search-<platform> docs-compare --query "Compare docs and community discussion"
```

## Flag Migration

Most non-mode flags keep the same meaning:

| Old flag | New flag |
| --- | --- |
| `--query` | `--query` |
| `--api-key` | `--api-key` |
| `--base-url` | `--base-url` |
| `--model` | `--model` |
| `--profile` | `--profile` |
| `--plain` | `--plain` |
| `--urls` | `--urls` |
| `--ignore-cooldown` | `--ignore-cooldown` |
| `--extra-body-json` | `--extra-body-json` |
| `--extra-headers-json` | `--extra-headers-json` |

## Config Migration

Old JSON config shape:

```json
{
  "base_url": "https://your-grok-endpoint.example",
  "model": "grok-4.1-fast",
  "timeout_seconds": 120,
  "profiles": [
    { "id": "main", "api_key": "YOUR_GROK_API_KEY" },
    { "id": "backup", "api_key": "YOUR_BACKUP_KEY" }
  ],
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

New TOML config at `~/.config/ai-skills/grok-search.toml`:

```toml
base_url = "https://api.x.ai"
model = "grok-4.1-fast"
timeout = 120

[[profiles]]
id = "main"
api_key = "YOUR_GROK_API_KEY"

[[profiles]]
id = "backup"
api_key = "YOUR_BACKUP_KEY"

[extra_body]

[extra_headers]

[cooldown]
enabled = true
state_file = "runtime/cooldowns.json"
default_minutes = 15
rate_limit_minutes = 20
quota_minutes = 60
auth_minutes = 360
```

## Endpoint Defaults

The new binary defaults to:

```toml
base_url = "https://api.x.ai"
model = "grok-4.1-fast"
```

For custom OpenAI-compatible endpoints, set a global `base_url` or use profile-level overrides:

```toml
[[profiles]]
id = "proxy"
api_key = "YOUR_PROXY_KEY"
base_url = "https://your-compatible-endpoint.example"
model = "grok-custom-model"
```
