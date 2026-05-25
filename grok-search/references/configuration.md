# Configuration Guide

## Default Config

`grok-search` uses TOML configuration at:

```text
~/.config/ai-skills/grok-search.toml
```

The binary creates a template at this path on first run when it does not exist.

Default endpoint settings:

```toml
base_url = "https://api.x.ai"
model = "grok-4.1-fast"
timeout = 120
```

## Configuration Priority

Highest to lowest:

1. CLI flags
2. Environment variables
3. TOML config
4. Built-in defaults

## Single Key Setup

```toml
base_url = "https://api.x.ai"
model = "grok-4.1-fast"
timeout = 120

[[profiles]]
id = "main"
api_key = "YOUR_GROK_API_KEY"
```

## Multiple Keys with Auto Failover

Profiles are tried in order:

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
```

## Profile-Level Endpoint Overrides

Use profile overrides for OpenAI-compatible proxies or custom gateways:

```toml
base_url = "https://api.x.ai"
model = "grok-4.1-fast"

[[profiles]]
id = "main"
api_key = "XAI_KEY"

[[profiles]]
id = "proxy"
api_key = "PROXY_KEY"
base_url = "https://your-compatible-endpoint.example"
model = "grok-custom-model"
```

A profile-level `base_url` or `model` overrides the global value only for that profile.

## Environment Variables

```bash
export GROK_API_KEY="key"
export GROK_API_KEYS="key1,key2,key3"
export GROK_BASE_URL="https://api.x.ai"
export GROK_MODEL="grok-4.1-fast"
export GROK_TIMEOUT="120"
export GROK_CONFIG="/path/to/grok-search.toml"
```

Extra request body and header overrides accept JSON objects:

```bash
export GROK_EXTRA_BODY_JSON='{"search_parameters": {}}'
export GROK_EXTRA_HEADERS_JSON='{"X-Custom-Header": "value"}'
```

CLI equivalents:

```bash
bin/grok-search-<platform> --api-key KEY research --query "test"
bin/grok-search-<platform> --base-url https://api.x.ai --model grok-4.1-fast research --query "test"
bin/grok-search-<platform> --extra-body-json '{"temperature":0.1}' research --query "test"
```

## Failover Behavior

Failover triggers on:

- HTTP 401
- HTTP 403
- HTTP 429
- Error text containing rate limit, quota, credits, billing, exhausted, unauthorized, forbidden, invalid API key, or token-unavailable indicators

Output includes `profileId` and `attempts` so agents can see which profile succeeded or failed.

Force one configured profile:

```bash
bin/grok-search-<platform> --profile main research --query "test query" --plain
```

## Cooldown Behavior

Failover-worthy failures place the profile into temporary cooldown. Later runs skip cooling profiles unless `--ignore-cooldown` is set.

Default cooldown config:

```toml
[cooldown]
enabled = true
state_file = "runtime/cooldowns.json"
default_minutes = 15
rate_limit_minutes = 20
quota_minutes = 60
auth_minutes = 360
```

Cooldown durations:

- Rate limit: 20 minutes
- Quota/billing: 60 minutes
- Auth errors: 360 minutes
- Other failover-worthy errors: 15 minutes

Override cooldown intentionally:

```bash
bin/grok-search-<platform> --ignore-cooldown news --query "latest updates" --plain
```

## Full TOML Example

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
# search_parameters = {}

[extra_headers]
# X-Custom-Header = "value"

[cooldown]
enabled = true
state_file = "runtime/cooldowns.json"
default_minutes = 15
rate_limit_minutes = 20
quota_minutes = 60
auth_minutes = 360
```

## Troubleshooting

### `missing_api_key`

Add a profile to `~/.config/ai-skills/grok-search.toml`, set `GROK_API_KEY` / `GROK_API_KEYS`, or pass `--api-key`.

### `invalid_config`

Check TOML syntax. Strings must be quoted, and profiles use `[[profiles]]`.

### `invalid_json`

`--extra-body-json`, `--extra-headers-json`, `GROK_EXTRA_BODY_JSON`, and `GROK_EXTRA_HEADERS_JSON` must be JSON objects.

### `all_profiles_in_cooldown`

Wait for cooldown expiry or retry with `--ignore-cooldown`.
