---
name: db-explorer
description: >
  Read-only exploration of PostgreSQL, MySQL, and SQLite databases. Use for listing schemas,
  tables, and views; inspecting columns, primary keys, indexes, and foreign keys; sampling data;
  running read-only SQL; verifying code models against database schema; and troubleshooting stored records.
  Only trigger when the user is explicitly in a database / SQL / schema / record / column context.
  Do NOT trigger for HTML / Markdown / UI tables or other non-database "table" scenarios.
---

# DB Explorer

Use the prebuilt `db-explorer` Go binary for deterministic, read-only database exploration.
Goal: retrieve trustworthy database metadata and small samples with stable JSON output for agent parsing.

## Use this skill when

- The user wants to list database schemas, tables, or views
- The user wants to inspect columns, types, defaults, primary keys, indexes, or foreign keys
- The user wants to sample a few rows of data
- The user wants to run a read-only SQL query
- The user wants to verify that code models match the database schema
- The user is troubleshooting stored records or column values

## Do not use

- Any write operation, DDL, migration, data fix, backfill, or destructive maintenance operation
- Bulk data export or heavy data processing
- Non-database tables such as HTML, Markdown, UI, or spreadsheet tables
- Unclear or unsafe production connections where the query may expose sensitive data

## Setup

The skill includes precompiled binaries for major platforms in `bin/`. Detect the user's platform and select the matching binary:

- Linux x86_64: `bin/db-explorer-linux-amd64`
- Linux ARM64: `bin/db-explorer-linux-arm64`
- macOS Intel: `bin/db-explorer-darwin-amd64`
- macOS Apple Silicon: `bin/db-explorer-darwin-arm64`
- Windows x86_64: `bin/db-explorer-windows-amd64.exe`

Platform detection:

```bash
uname -s  # Linux, Darwin, or MINGW64_NT
uname -m  # x86_64, arm64, aarch64
```

Invoke the selected binary directly. No Python, virtualenv, pip install, or database driver setup is required.

## Capability contract

Supported databases:

- `sqlite`
- `postgres`
- `mysql`

Supported commands:

- `test`: test connection
- `schemas`: list schemas / namespaces
- `tables`: list tables with low-cost row estimate metadata when available
- `views`: list views
- `schema <table>`: show columns, primary keys, indexes, and foreign keys
- `data <table> --limit N`: sample rows with bounded output
- `query "<sql>"`: run one read-only SQL statement
- `version`: show binary version metadata

Global flags:

```text
--profile <id>
--db sqlite|postgres|mysql
--url <connection-url-or-sqlite-path>
--url-env <ENV_VAR>
--config <path>
--format json|table|markdown|csv
--timeout <seconds>
--debug
```

Defaults:

- Output format: `json`
- Timeout: `30` seconds
- Data sample limit: `10` rows

All JSON responses include `schema_version: "1"`.

## Configuration

Connection resolution priority:

1. CLI flags
2. Project config: `.db-explorer.toml`
3. Global config: `~/.config/ai-skills/db-explorer.toml`
4. Environment fallback: `DATABASE_URL`, `DB_URL`, `POSTGRES_URL`, `MYSQL_URL`

Example `.db-explorer.toml`:

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

Prefer `url_env` for credentials. Do not print full connection URLs that contain passwords or tokens.

See `references/configuration.md` for details.

## Workflow

### 1. Resolve connection info

- Prefer existing `.db-explorer.toml` in the project root.
- If no project config exists, use explicit user-provided flags, global config, or environment fallback.
- If the user provides a `.db`, `.sqlite`, or `.sqlite3` path, default to `--db sqlite` unless context clearly says otherwise.
- When using an env var, pass `--url-env <ENV_VAR>` and do not print its value.
- Ask only for missing information that blocks execution.

### 2. Verify connection first

```bash
bin/db-explorer-<platform> test --profile local
bin/db-explorer-<platform> test --db sqlite --url ./app.db
bin/db-explorer-<platform> test --db postgres --url-env DATABASE_URL
```

If the test fails, report the JSON error code and masked message. Do not speculate about schema or data before connection verification succeeds.

### 3. Choose the minimal operation

```bash
# List schemas, tables, or views
bin/db-explorer-<platform> schemas --profile local
bin/db-explorer-<platform> tables --profile local
bin/db-explorer-<platform> views --profile local

# Inspect structure
bin/db-explorer-<platform> schema users --profile local
bin/db-explorer-<platform> schema public.users --db postgres --url-env DATABASE_URL

# Sample data
bin/db-explorer-<platform> data users --limit 10 --profile local

# Run read-only SQL
bin/db-explorer-<platform> query "SELECT id, email FROM users LIMIT 10" --profile local
```

Selection strategy:

- Need orientation: `schemas`, `tables`, `views`
- Need definition: `schema <table>`
- Need example records: `data <table> --limit 10`
- Need filtering / joins / aggregation: `query "<sql>"`

### 4. SQL safety rules

Only run read-only SQL. The binary rejects:

- Multiple statements
- `INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`, `TRUNCATE`, `CREATE`, `GRANT`, `REVOKE`, `EXEC`, `MERGE`, `CALL`, `VACUUM`, `REINDEX`, `ATTACH`, `DETACH`
- Unsupported command prefixes
- Unsafe state-changing SQLite PRAGMA statements

Allowed custom SQL prefixes:

- `SELECT`
- `WITH`
- `SHOW`
- `DESCRIBE` / `DESC`
- `EXPLAIN`
- approved read-only metadata `PRAGMA`

For exploratory SQL written on behalf of the user, add a reasonable `LIMIT` unless the query semantics make a limit inappropriate.

### 5. Present results

Do not dump raw terminal output. Parse JSON and summarize:

- Tables/views: list relation names, schemas, and row estimate metadata when present
- Schema: show columns first, then indexes and foreign keys
- Data/query: show concise rows; mention if `meta.truncated` is true
- Errors: report `error.code` and the masked `error.message`

### 6. Compare with code models

If cross-checking ORM/model definitions:

1. Run `schema <table>`.
2. Open the model file.
3. Report only key differences: missing/extra fields, type mismatches, nullability differences, default differences, primary key/index/foreign key differences.

## Guardrails

- This is a read-only skill.
- Never execute or suggest write SQL, DDL, migrations, fixes, or backfills.
- Do not expose passwords, tokens, or full secret-bearing URLs.
- Do not rely on exact row counts from `tables`; row metadata is an estimate or unknown unless explicitly queried.
- When a query fails, report the structured error and next troubleshooting step; do not fabricate results.

## Additional Resources

- Source code: `db-explorer-go/`
- Configuration: `references/configuration.md`
- Migration notes: `references/migration-from-python.md`
- Evaluations: `evals/db-explorer/`
