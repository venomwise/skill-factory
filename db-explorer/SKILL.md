---
name: db-explorer
description: >
  Read-only CLI for exploring PostgreSQL, MySQL, and SQLite databases. Use this skill whenever a task
  needs reliable database facts, including schema details, sample records, read-only query results, or
  code-model verification against actual database structure. Do not use for writes, migrations, data
  fixes, bulk exports, or non-database tables.
---

# DB Explorer

Use the bundled `db-explorer` Go binary only. Default to JSON. Run the smallest safe command, parse the JSON envelope, then summarize results.

Hard rules:

- Read-only only; never execute or suggest write SQL, DDL, migrations, data fixes, or backfills.
- Do not install dependencies, use helper scripts, or bypass the bundled binary with `psql`, `mysql`, `sqlite3`, or ad hoc code.
- Do not print passwords, tokens, or full secret-bearing URLs.

## Binary Selection

Select the matching precompiled binary in `bin/`:

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

## Capability contract

Databases: `sqlite`, `postgres`, `mysql`.

Commands:

- `test`: test connection
- `schemas`: list schemas / namespaces
- `tables`: list tables with low-cost row estimate metadata when available
- `views`: list views
- `schema <table>`: show columns, primary keys, indexes, and foreign keys
- `data <table> --limit N`: sample rows with bounded output
- `query "<sql>"`: run one read-only SQL statement
- `version`: show binary version metadata

Connection flags:

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

Defaults: JSON output, 30 second timeout, `data --limit 10`.

JSON responses contain `schema_version: "1"`, `ok`, `command`, `data` or `error`, and `meta`.

Only `query` and `data` support `--format table|markdown|csv`; other commands require JSON. Use non-JSON formats only when the user explicitly asks for rendered rows.

## Connection

Connection priority:

1. CLI flags
2. Project config: `.db-explorer.toml`
3. Global config: `~/.config/ai-skills/db-explorer.toml`
4. Environment fallback: `DATABASE_URL`, `DB_URL`, `POSTGRES_URL`, `MYSQL_URL`

Prefer project `.db-explorer.toml`. If the user provides a `.db`, `.sqlite`, or `.sqlite3` path, default to `--db sqlite` unless context says otherwise. Use `--url-env <ENV_VAR>` for credentials and do not print the env var value.

See `references/configuration.md` for details.

## Workflow

### 1. Verify connection

Ask only for missing connection information that blocks execution. Test before schema or data commands:

```bash
bin/db-explorer-<platform> test --profile local
bin/db-explorer-<platform> test --db sqlite --url ./app.db
bin/db-explorer-<platform> test --db postgres --url-env DATABASE_URL
```

If `test` fails, report the JSON `error.code` and masked `error.message`. Do not speculate about schema or data.

### 2. Run the smallest useful command

```bash
bin/db-explorer-<platform> schemas --profile local
bin/db-explorer-<platform> tables --profile local
bin/db-explorer-<platform> views --profile local
bin/db-explorer-<platform> schema users --profile local
bin/db-explorer-<platform> schema public.users --db postgres --url-env DATABASE_URL
bin/db-explorer-<platform> data users --limit 10 --profile local
bin/db-explorer-<platform> query "SELECT id, email FROM users LIMIT 10" --profile local
```

- Need orientation: `schemas`, `tables`, `views`
- Need definition: `schema <table>`
- Need example records: `data <table> --limit 10`
- Need filtering, joins, or aggregation: `query "<sql>"`

Use built-in commands before custom SQL when they answer the question.

### 3. SQL safety

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

### 4. Sensitive data

For production or sensitive-looking contexts, inspect schema first. Do not sample secret, token, password, payment, personal, or authentication columns unless the user explicitly asks and the connection is safe. Prefer counts, schema, or non-sensitive columns over raw row samples.

### 5. Present results

Do not dump raw terminal output from JSON runs. Parse the envelope and summarize:

- Tables/views: list relation names, schemas, and row estimate metadata when present
- Schema: show columns first, then indexes and foreign keys
- Data/query: show concise rows; mention if `meta.truncated` is true
- Errors: report `error.code` and the masked `error.message`

If a command fails:

- Config/profile/env error: check the selected profile, config file, or env var name.
- Relation missing: run `tables` or `views`.
- Permission error: report the missing read permission.
- SQL rejected: use a built-in command or rewrite as one safe read-only statement.

### 6. Compare code models

When cross-checking ORM/model definitions:

1. Run `schema <table>`.
2. Open the model file.
3. Report only key differences: missing/extra fields, type mismatches, nullability differences, default differences, primary key/index/foreign key differences.

## Additional Resources

- Configuration: `references/configuration.md`
