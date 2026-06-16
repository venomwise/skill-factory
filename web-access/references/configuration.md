# Web Access Configuration

## Configuration File

Default path: `~/.config/ai-skills/web-access.toml`

Override with `--config` flag or `WEB_ACCESS_CONFIG` environment variable.

## Template Creation

If the default config file doesn't exist, `web-access` creates a template automatically on first run.

## Configuration Precedence

Settings are resolved in the following order (highest to lowest priority):

1. **CLI flags** — `--exa-api-key`, `--grok-api-key`, `--grok-model`, `--timeout`, `--profile`
2. **`WEB_ACCESS_*` environment variables** — `WEB_ACCESS_EXA_API_KEY`, `WEB_ACCESS_GROK_API_KEY`, `WEB_ACCESS_CONFIG`
3. **Provider-specific environment variables** — `EXA_API_KEY`, `EXA_API_KEYS`, `GROK_API_KEY`, `GROK_API_KEYS`
4. **TOML config file** — `~/.config/ai-skills/web-access.toml`
5. **Built-in defaults**

## TOML Structure

```toml
[exa]
base_url = "https://api.exa.ai"
timeout = 30

[[exa.profiles]]
id = "primary"
api_key = "YOUR_EXA_API_KEY"

[[exa.profiles]]
id = "backup"
api_key = "YOUR_BACKUP_EXA_KEY"

[grok]
base_url = "https://api.x.ai/v1"
model = "grok-beta"
timeout = 60

[[grok.profiles]]
id = "primary"
api_key = "YOUR_GROK_API_KEY"

[[grok.profiles]]
id = "backup"
api_key = "YOUR_BACKUP_GROK_KEY"
```

## Profile Resolution

Profiles allow multiple API keys for failover:

- **Exa profiles**: Resolved from `EXA_API_KEY`, `EXA_API_KEYS`, `WEB_ACCESS_EXA_API_KEY`, or TOML `[[exa.profiles]]`
- **Grok profiles**: Resolved from `GROK_API_KEY`, `GROK_API_KEYS`, `WEB_ACCESS_GROK_API_KEY`, or TOML `[[grok.profiles]]`

Placeholder keys like `YOUR_EXA_API_KEY` or `YOUR_GROK_API_KEY` are automatically filtered out.

Use `--profile <id>` to select a specific profile by ID.

## Failover Behavior

When a profile fails with rate limit, quota, or auth errors, `web-access` automatically tries the next available profile.

**No cooldown state** is persisted — each run starts fresh with all configured profiles.

## Extra Body and Headers (Grok)

For Grok requests, you can merge additional JSON into the request:

- `--extra-body-json '{"temperature": 0.7}'`
- `--extra-headers-json '{"X-Custom": "value"}'`
- Environment: `GROK_EXTRA_BODY_JSON`, `GROK_EXTRA_HEADERS_JSON`
- TOML: `extra_body` and `extra_headers` under `[grok]` or `[[grok.profiles]]`

## Error Handling

- **`missing_api_key`**: No valid API key found for the provider
- **`config_parse_error`**: Invalid TOML syntax in config file
- **`request_failed`**: Network error or timeout
- **`all_profiles_failed`**: All configured profiles failed

See command output for detailed retry guidance.
