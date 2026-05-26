# db-explorer Go Redesign Design

## Summary

Rebuild `db-explorer` from a Python helper script into an Agent-first Go binary skill. The new version intentionally does not preserve the old Python CLI or `.db-explorer.json` contract. It will ship prebuilt binaries under `db-explorer/bin/`, default to stable JSON output, support SQLite/PostgreSQL/MySQL, and expose richer database metadata for reliable agent parsing.

## Goals

- Remove Python runtime, virtualenv, and pip dependency requirements from the skill.
- Add a new `db-explorer-go/` Go source project following the existing Go skill pattern used by `exa-search-go/` and `grok-search-go/`.
- Distribute platform-specific binaries in `db-explorer/bin/` for Linux, macOS, and Windows.
- Make JSON the default output format for all commands, with a stable response envelope and `schema_version`.
- Support Agent-friendly connection discovery via CLI flags, project config, global config, and environment fallback.
- Continue supporting SQLite, PostgreSQL, and MySQL only.
- Improve metadata coverage for schemas/namespaces, tables, views, columns, primary keys, indexes, foreign keys, defaults, nullability, and low-cost row estimates.
- Preserve the core read-only guardrails for custom SQL queries.

## Primary Users / Roles

- **AI coding agents**: inspect local or configured databases, compare database schema with code models, run bounded read-only queries, and consume deterministic JSON results.
- **Repository maintainers**: build, test, release, and update the skill binaries with isolated Go workflows.

Human CLI readability is secondary. `table`, `markdown`, and `csv` output may exist, but they must not drive the primary command contract.

## Non-Goals

- No compatibility with `db-explorer/scripts/db_query.py` CLI shape.
- No compatibility with old `.db-explorer.json` project config.
- No new database engines such as DuckDB, SQL Server, ClickHouse, Oracle, or MongoDB.
- No write operations, DDL, migrations, data fixes, backfills, or bulk exports.
- No advanced production safety workflow, audit logging, approval prompts, or permission system in v1.
- No human-first interactive shell or TUI.
- No expensive exact row counts by default.

## Context

Current `db-explorer` is a Python script at `db-explorer/scripts/db_query.py`. It supports `test`, `tables`, `schema`, `data`, and `query` for SQLite/PostgreSQL/MySQL with table/markdown/json/csv output. PostgreSQL and MySQL require Python packages from `db-explorer/requirements.txt`, which forces virtualenv management inside the skill.

`exa-search` and `grok-search` have already moved to a hybrid repository shape: Go source lives in `<skill>-go/`, while the skill directory contains `SKILL.md`, references, and prebuilt binaries in `bin/`. `db-explorer` should follow that pattern.

Existing evals under `evals/db-explorer/` target the Python script and must be redesigned for the new binary and JSON-first output contract.

## Discovery

### Key Discoveries

- The migration is intentionally a redesign, not a compatibility-preserving port.
- The CLI is primarily for agents, so stable machine-readable JSON matters more than human-friendly defaults.
- The old script prints some schema details, such as indexes and foreign keys, to stderr. This is unsuitable for JSON consumers; the new design must put all structured metadata in stdout JSON.
- SQLite driver choice affects binary distribution. A CGO SQLite driver complicates cross-platform builds; a pure Go driver keeps the release workflow aligned with `exa-search` and `grok-search`.
- `tables` must not run unbounded `COUNT(*)` by default because large tables can make basic exploration slow or disruptive.
- Output schema versioning is useful from v1 because agents and evals will depend on the JSON structure.

### Scope Decisions

- Use a pure Go SQLite driver, preferably `modernc.org/sqlite`, so binaries can be built with `CGO_ENABLED=0`.
- Keep database support limited to SQLite, PostgreSQL, and MySQL.
- Replace old `.db-explorer.json` with TOML config.
- Support both project config and global config, but optimize behavior for agent discovery rather than human profile management.
- Include `schema_version: "1"` in every JSON response envelope.
- Return low-cost `row_estimate` / unknown row metadata by default; do not perform exact `COUNT(*)` as part of `tables`.

## Proposed Solution

Build a new Go CLI named `db-explorer` with dialect-specific metadata adapters and a shared JSON response contract. The skill directory will become a binary distribution and usage-instruction layer. The implementation should favor explicit, versioned JSON models over terminal-oriented output.

### Architecture

```text
db-explorer-go/
  cmd/
    db-explorer/
      main.go
    root.go
    test.go
    tables.go
    views.go
    schemas.go
    schema.go
    data.go
    query.go
    version.go
  internal/
    config/
    db/
    dialect/
      sqlite/
      postgres/
      mysql/
    introspect/
    output/
    safety/
```

```text
db-explorer/
  SKILL.md
  bin/
    db-explorer-darwin-amd64
    db-explorer-darwin-arm64
    db-explorer-linux-amd64
    db-explorer-linux-arm64
    db-explorer-windows-amd64.exe
    SHA256SUMS
  references/
    configuration.md
    migration-from-python.md
```

Primary command shape:

```bash
db-explorer test --profile local
db-explorer schemas --profile local
db-explorer tables --profile local
db-explorer views --profile local
db-explorer schema users --profile local
db-explorer data users --limit 10 --profile local
db-explorer query "SELECT id, email FROM users LIMIT 10" --profile local
```

Direct connection shape:

```bash
db-explorer tables --db sqlite --url ./app.db
db-explorer schema users --db postgres --url-env DATABASE_URL
```

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

```text
--format json
--timeout 30
```

### Components

- **`cmd/`**: Defines Cobra commands, binds flags, invokes config resolution, calls the relevant database/dialect operation, and sends results to the output layer.
- **`internal/config`**: Loads and merges connection settings from CLI flags, project `.db-explorer.toml`, global `~/.config/ai-skills/db-explorer.toml`, and environment fallback. Masks secrets in all summaries and errors.
- **`internal/db`**: Opens database connections, applies timeouts where supported, exposes query helpers, and closes resources.
- **`internal/dialect`**: Contains SQLite, PostgreSQL, and MySQL adapters. Each adapter knows how to list schemas, tables, views, columns, indexes, primary keys, foreign keys, and row estimates for its database.
- **`internal/introspect`**: Defines shared metadata models such as `Schema`, `Relation`, `Column`, `Index`, `ForeignKey`, and `QueryResult`.
- **`internal/safety`**: Validates custom SQL before execution. It rejects multiple statements, dangerous write/DDL keywords, and unsafe SQLite PRAGMA statements.
- **`internal/output`**: Emits JSON response envelopes by default and optional table/markdown/csv renderings for human readability or downstream tooling.

Config files use TOML.

Project config path:

```text
.db-explorer.toml
```

Global config path:

```text
~/.config/ai-skills/db-explorer.toml
```

Example:

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

Resolution priority:

1. CLI flags.
2. Project `.db-explorer.toml`.
3. Global `~/.config/ai-skills/db-explorer.toml`.
4. Environment fallback: `DATABASE_URL`, `DB_URL`, `POSTGRES_URL`, `MYSQL_URL`.

JSON success envelope:

```json
{
  "schema_version": "1",
  "ok": true,
  "command": "tables",
  "db": "postgres",
  "profile": "dev",
  "data": {},
  "meta": {
    "duration_ms": 12,
    "truncated": false
  }
}
```

JSON error envelope:

```json
{
  "schema_version": "1",
  "ok": false,
  "command": "query",
  "error": {
    "code": "SQL_NOT_READONLY",
    "message": "Only read-only SQL is allowed",
    "details": {}
  },
  "meta": {
    "duration_ms": 1
  }
}
```

Example `tables` data:

```json
{
  "tables": [
    {
      "schema": "public",
      "name": "users",
      "type": "table",
      "row_estimate": 120,
      "row_estimate_kind": "estimate"
    }
  ]
}
```

Example `schema` data:

```json
{
  "table": {
    "schema": "public",
    "name": "users",
    "type": "table"
  },
  "columns": [
    {
      "name": "id",
      "type": "integer",
      "nullable": false,
      "default": null,
      "primary_key": true
    }
  ],
  "indexes": [],
  "foreign_keys": []
}
```

### Data Flow

Metadata command flow:

1. Agent selects the platform binary from `db-explorer/bin/`.
2. Agent runs a command such as `db-explorer schema users --profile dev`.
3. CLI parses command and flags.
4. Config loader resolves the connection profile and masks secrets for diagnostics.
5. Database connector opens the target database with the configured timeout.
6. Dialect adapter executes metadata queries for the target database.
7. Introspection layer maps database-specific rows into shared metadata structs.
8. Output layer writes a versioned JSON envelope to stdout.

Custom query flow:

1. Agent runs `db-explorer query "SELECT id FROM users LIMIT 10" --profile dev`.
2. Safety layer strips comments, validates single-statement read-only SQL, and rejects unsafe PRAGMA or dangerous keywords.
3. Database connector opens the connection and applies timeout settings where supported.
4. Query executes and returns columns plus rows.
5. Output layer emits JSON with row data and `meta.truncated` if a limit or output cap is applied.

## Error Handling

Errors should be returned as JSON by default, even when the command fails. Exit codes should be non-zero for failed operations.

Initial error codes:

- `CONFIG_NOT_FOUND`
- `PROFILE_NOT_FOUND`
- `INVALID_CONFIG`
- `MISSING_CONNECTION`
- `ENV_NOT_SET`
- `UNSUPPORTED_DB`
- `CONNECTION_FAILED`
- `SQL_NOT_READONLY`
- `SQL_MULTIPLE_STATEMENTS`
- `UNSAFE_PRAGMA`
- `QUERY_TIMEOUT`
- `QUERY_FAILED`
- `TABLE_NOT_FOUND`
- `OUTPUT_FAILED`

Handling rules:

- Never print full passwords, tokens, or unmasked connection URLs.
- Preserve useful database error messages after secret masking.
- Do not fabricate empty metadata when metadata queries fail.
- Do not run exact row counts during `tables` unless a future explicit command or flag requests it.
- Apply bounded defaults to `data` and query result output to avoid excessive terminal output.
- Keep stderr for debug diagnostics only; structured command results belong on stdout.

## Testing

Go unit tests:

- Config file loading and precedence.
- Profile resolution.
- Environment variable resolution.
- Secret masking.
- SQL read-only validation.
- Multiple-statement detection.
- Unsafe PRAGMA rejection.
- JSON envelope rendering.
- Optional table/markdown/csv formatting.

SQLite integration tests run by default in CI:

- `test` succeeds on a temporary SQLite database.
- `schemas` returns SQLite-compatible schema metadata.
- `tables` lists tables without exact count requirements.
- `views` lists views.
- `schema users` returns columns, primary key, indexes, and foreign keys in JSON.
- `data users --limit 10` returns bounded rows.
- `query` accepts read-only SELECT/WITH queries.
- Dangerous SQL such as `VACUUM`, `DROP`, and multi-statement input is rejected.

Optional live database tests:

- Enabled only when `DBX_POSTGRES_URL` and/or `DBX_MYSQL_URL` are set.
- Cover PostgreSQL schema namespace behavior, views, indexes, foreign keys, row estimates, and read-only queries.
- Cover MySQL tables, views, indexes, foreign keys, and read-only queries.

Skill eval updates:

- Update `evals/db-explorer/` to invoke the Go binary instead of the Python script.
- Expect JSON by default.
- Add metadata cases for views, indexes, foreign keys, and non-default schemas where practical.
- Preserve negative safety cases for dangerous SQL, unsafe PRAGMA, and multiple statements.

CI/release workflows:

```text
.github/workflows/db-explorer-test.yml
.github/workflows/db-explorer-release.yml
.github/workflows/db-explorer-update-skill.yml
```

The workflows should mirror the isolated Go skill pattern:

- Test workflow runs `go test ./...`, gofmt check, `go vet ./...`, and a local build.
- Release workflow builds archives for Linux/macOS/Windows on `db-explorer-v*` tags or manual dispatch.
- Update-skill workflow builds platform binaries into `db-explorer/bin/`, generates `SHA256SUMS`, and commits changed binaries.

## Open Questions

None for the current design phase. The following decisions are fixed for v1:

- SQLite uses a pure Go driver, preferably `modernc.org/sqlite`, to keep `CGO_ENABLED=0` builds.
- `tables` returns low-cost `row_estimate` metadata when available and does not run exact `COUNT(*)` by default.
- Every JSON response includes `schema_version: "1"`.
