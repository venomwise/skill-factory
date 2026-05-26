# db-explorer Configuration

`db-explorer` resolves connection settings from CLI flags, project config, global config, and environment variables. The CLI is Agent-first: prefer explicit, machine-readable settings and avoid printing secret-bearing URLs.

## Precedence

Connection resolution priority:

1. CLI flags
2. Project config: `.db-explorer.toml` in the current project root
3. Global config: `~/.config/ai-skills/db-explorer.toml`
4. Environment fallback: `DATABASE_URL`, `DB_URL`, `POSTGRES_URL`, `MYSQL_URL`

CLI flags override config values. Project config takes precedence over global config when both define usable profiles.

## CLI flags

```bash
bin/db-explorer-<platform> tables --db sqlite --url ./app.db
bin/db-explorer-<platform> tables --db postgres --url-env DATABASE_URL
bin/db-explorer-<platform> tables --profile local
bin/db-explorer-<platform> tables --config ./custom-db-explorer.toml --profile dev
```

Supported connection flags:

| Flag | Description |
| --- | --- |
| `--profile <id>` | Use a named TOML profile |
| `--db sqlite\|postgres\|mysql` | Database type for direct connections |
| `--url <value>` | Direct connection URL or SQLite file path |
| `--url-env <ENV_VAR>` | Environment variable containing the URL |
| `--config <path>` | Explicit TOML config path |

Other global flags:

| Flag | Default | Description |
| --- | --- | --- |
| `--format json\|table\|markdown\|csv` | `json` | Output format |
| `--timeout <seconds>` | `30` | Connection/query timeout |
| `--debug` | `false` | Enable debug diagnostics |

## TOML format

```toml
default_profile = "local"

[[profiles]]
id = "local"
db = "sqlite"
url = "./app.db"

[[profiles]]
id = "dev"
db = "postgres"
url_env = "DATABASE_URL"

[[profiles]]
id = "mysql-dev"
db = "mysql"
url_env = "MYSQL_URL"
```

Fields:

| Field | Required | Description |
| --- | --- | --- |
| `default_profile` | No | Profile used when `--profile` is omitted |
| `profiles[].id` | Yes | Profile name |
| `profiles[].db` | Yes | `sqlite`, `postgres`, or `mysql` |
| `profiles[].url` | One of `url` / `url_env` | Direct URL or SQLite path |
| `profiles[].url_env` | One of `url` / `url_env` | Environment variable holding the URL |

Prefer `url_env` for credentials.

## SQLite examples

Project config:

```toml
default_profile = "local"

[[profiles]]
id = "local"
db = "sqlite"
url = "./data/app.db"
```

Direct command:

```bash
bin/db-explorer-<platform> tables --db sqlite --url ./data/app.db
```

SQLite URLs are also accepted:

```bash
bin/db-explorer-<platform> tables --db sqlite --url sqlite:///absolute/path/app.db
```

## PostgreSQL examples

Project config using an environment variable:

```toml
default_profile = "dev"

[[profiles]]
id = "dev"
db = "postgres"
url_env = "DATABASE_URL"
```

Command:

```bash
bin/db-explorer-<platform> schema public.users --profile dev
```

Direct environment variable command:

```bash
bin/db-explorer-<platform> tables --db postgres --url-env DATABASE_URL
```

## MySQL examples

Project config:

```toml
[[profiles]]
id = "mysql-dev"
db = "mysql"
url_env = "MYSQL_URL"
```

Command:

```bash
bin/db-explorer-<platform> schema users --profile mysql-dev
```

Direct command:

```bash
bin/db-explorer-<platform> tables --db mysql --url-env MYSQL_URL
```

## Environment fallback

If no CLI or config connection is available, `db-explorer` checks these variables:

1. `POSTGRES_URL` → PostgreSQL
2. `MYSQL_URL` → MySQL
3. `DATABASE_URL` → inferred from URL scheme
4. `DB_URL` → inferred from URL scheme

Inference examples:

| Value | Inferred DB |
| --- | --- |
| `postgres://user:pass@host/db` | `postgres` |
| `postgresql://user:pass@host/db` | `postgres` |
| `mysql://user:pass@host/db` | `mysql` |
| `sqlite:///tmp/app.db` | `sqlite` |
| `./app.sqlite3` | `sqlite` |

## Secret handling

- Do not print full connection URLs when they may include credentials.
- Prefer `--url-env` / `url_env` over `--url` for PostgreSQL and MySQL.
- Error messages should mask passwords, tokens, API keys, and secret query parameters.

## JSON contract

JSON is the default output format. Every JSON response includes:

```json
{
  "schema_version": "1",
  "ok": true,
  "command": "tables",
  "data": {},
  "meta": {
    "truncated": false
  }
}
```

Failures include an error object:

```json
{
  "schema_version": "1",
  "ok": false,
  "command": "query",
  "error": {
    "code": "SQL_NOT_READONLY",
    "message": "Only read-only SQL is allowed"
  },
  "meta": {
    "truncated": false
  }
}
```
