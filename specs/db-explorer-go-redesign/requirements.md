# Requirements Document: db-explorer Go Redesign

## Introduction

`db-explorer` will be rebuilt as an Agent-first Go binary for read-only database exploration. It replaces the current Python helper script and removes Python runtime, virtualenv, and pip dependency requirements while keeping the skill focused on SQLite, PostgreSQL, and MySQL.

The new tool exposes a redesigned CLI, TOML-based connection discovery, stable JSON output with schema versioning, richer database metadata, and read-only query guardrails. It is distributed like the existing Go skills in this repository: source code lives under `db-explorer-go/`, while the skill directory contains instructions, references, and platform-specific binaries under `db-explorer/bin/`.

## Glossary

- **Agent-first**: Optimized primarily for AI coding agents that need deterministic, machine-readable command results.
- **Dialect adapter**: Database-specific implementation for SQLite, PostgreSQL, or MySQL metadata queries and connection behavior.
- **JSON envelope**: The top-level JSON response wrapper containing `schema_version`, `ok`, command metadata, `data` or `error`, and `meta`.
- **Profile**: A named connection entry in `.db-explorer.toml` or `~/.config/ai-skills/db-explorer.toml`.
- **Project config**: A `.db-explorer.toml` file in the current project root.
- **Global config**: `~/.config/ai-skills/db-explorer.toml`.
- **Row estimate**: A low-cost approximate or unknown row count returned during table listing; not an exact `COUNT(*)`.
- **Read-only SQL**: SQL accepted by the custom query command after rejecting dangerous keywords, multiple statements, and unsafe PRAGMA statements.

## Requirements

### Requirement 1: Go Binary Project Structure and Distribution

**User Story:** As a repository maintainer, I want `db-explorer` implemented as a Go binary skill, so that agents can use it without Python runtime setup.

#### Acceptance Criteria

1. WHEN the repository is checked out, THEN the system SHALL contain a `db-explorer-go/` Go module for the CLI source code.
2. WHEN the Go CLI is built, THEN the output binary SHALL be named `db-explorer` or `db-explorer.exe` on Windows.
3. WHEN skill binaries are updated, THEN platform-specific binaries SHALL be placed under `db-explorer/bin/` using the naming pattern `db-explorer-<os>-<arch>`.
4. WHEN binaries are generated, THEN `db-explorer/bin/SHA256SUMS` SHALL include checksums for the generated binaries.
5. WHEN implementing SQLite support, THEN the system SHALL use a pure Go SQLite driver so builds can run with `CGO_ENABLED=0`.
6. IF Python-only files are no longer required by the skill, THEN the skill documentation SHALL stop instructing agents to create virtualenvs or install pip dependencies.

### Requirement 2: Redesigned Agent-Oriented CLI

**User Story:** As an AI coding agent, I want a clear command-oriented CLI, so that I can inspect databases with low ambiguity.

#### Acceptance Criteria

1. WHEN the agent runs `db-explorer test`, THEN the system SHALL test the resolved database connection.
2. WHEN the agent runs `db-explorer schemas`, THEN the system SHALL list supported schemas or namespaces for the target database.
3. WHEN the agent runs `db-explorer tables`, THEN the system SHALL list tables for the target database.
4. WHEN the agent runs `db-explorer views`, THEN the system SHALL list views for the target database.
5. WHEN the agent runs `db-explorer schema <table>`, THEN the system SHALL return metadata for the named relation.
6. WHEN the agent runs `db-explorer data <table> --limit N`, THEN the system SHALL return at most `N` sampled rows from the named relation.
7. WHEN the agent runs `db-explorer query "<sql>"`, THEN the system SHALL execute the SQL only after read-only validation succeeds.
8. WHEN global flags are provided, THEN the CLI SHALL support `--profile`, `--db`, `--url`, `--url-env`, `--config`, `--format`, `--timeout`, and `--debug`.
9. WHEN no output format is specified, THEN the CLI SHALL default to JSON.
10. WHEN no timeout is specified, THEN the CLI SHALL default to 30 seconds.

### Requirement 3: Connection Configuration Resolution

**User Story:** As an AI coding agent, I want deterministic connection discovery, so that I can connect using explicit flags, project config, global config, or environment fallback.

#### Acceptance Criteria

1. WHEN `--db` and `--url` are provided, THEN the system SHALL use those CLI values as the connection source.
2. WHEN `--url-env ENV_NAME` is provided, THEN the system SHALL read the connection URL from that environment variable without printing its secret value.
3. WHEN `--profile <id>` is provided, THEN the system SHALL resolve the profile from project or global TOML config.
4. WHEN no explicit profile is provided and a config file has `default_profile`, THEN the system SHALL use that profile.
5. WHEN both project and global config define usable profiles, THEN project `.db-explorer.toml` SHALL take precedence over global `~/.config/ai-skills/db-explorer.toml`.
6. WHEN CLI flags conflict with config values, THEN CLI flags SHALL take precedence.
7. WHEN no CLI or config connection is available, THEN the system SHALL check `DATABASE_URL`, `DB_URL`, `POSTGRES_URL`, and `MYSQL_URL` as environment fallback.
8. IF a referenced profile does not exist, THEN the system SHALL return `PROFILE_NOT_FOUND`.
9. IF a referenced environment variable is not set, THEN the system SHALL return `ENV_NOT_SET`.
10. IF no connection can be resolved, THEN the system SHALL return `MISSING_CONNECTION`.

### Requirement 4: Stable JSON Output Contract

**User Story:** As an AI coding agent, I want stable JSON responses, so that I can parse command results reliably across executions.

#### Acceptance Criteria

1. WHEN a command succeeds with JSON output, THEN stdout SHALL contain a JSON envelope with `schema_version`, `ok`, `command`, `db`, `data`, and `meta` fields.
2. WHEN a command fails with JSON output, THEN stdout SHALL contain a JSON envelope with `schema_version`, `ok`, `command`, `error`, and `meta` fields.
3. WHEN any JSON response is emitted, THEN `schema_version` SHALL equal `"1"` for v1.
4. WHEN a command succeeds, THEN `ok` SHALL be `true` and `error` SHALL be absent.
5. WHEN a command fails, THEN `ok` SHALL be `false` and `error.code` plus `error.message` SHALL be present.
6. WHEN command execution duration is known, THEN `meta.duration_ms` SHALL be included.
7. WHEN output is bounded or truncated, THEN `meta.truncated` SHALL indicate the truncation state.
8. WHEN optional `table`, `markdown`, or `csv` formats are requested, THEN those formats SHALL NOT change the default JSON contract.
9. IF diagnostics are emitted, THEN they SHALL NOT corrupt stdout JSON.

### Requirement 5: Database Metadata Introspection

**User Story:** As an AI coding agent, I want complete relation metadata, so that I can compare database schema with application code and reason about relationships.

#### Acceptance Criteria

1. WHEN inspecting SQLite, PostgreSQL, or MySQL, THEN the system SHALL support metadata commands for that database type.
2. WHEN listing schemas, THEN the system SHALL return database-supported schemas or namespaces where available.
3. WHEN listing tables, THEN each relation SHALL include name, type, schema or namespace where applicable, and row estimate metadata when cheaply available.
4. WHEN listing views, THEN each view SHALL include name, type, and schema or namespace where applicable.
5. WHEN inspecting a relation schema, THEN the system SHALL return columns with name, database type, nullability, default value, and primary key status.
6. WHEN inspecting a relation schema, THEN the system SHALL return indexes where the database exposes them.
7. WHEN inspecting a relation schema, THEN the system SHALL return foreign keys where the database exposes them.
8. WHEN row information is not cheaply available, THEN the system SHALL return unknown or null row estimate metadata rather than running exact unbounded `COUNT(*)`.
9. IF metadata retrieval fails, THEN the system SHALL return a structured error instead of fabricating empty metadata.

### Requirement 6: Data Sampling and Read-Only Query Execution

**User Story:** As an AI coding agent, I want bounded samples and read-only SQL execution, so that I can inspect data without mutating the database.

#### Acceptance Criteria

1. WHEN `data <table>` is executed without `--limit`, THEN the system SHALL use a safe default row limit.
2. WHEN `data <table> --limit N` is executed, THEN the system SHALL return no more than `N` rows.
3. WHEN `query "<sql>"` is executed with valid read-only SQL, THEN the system SHALL return columns and rows in the JSON `data` field.
4. WHEN a query result is capped by an output limit, THEN `meta.truncated` SHALL be `true`.
5. WHEN a query returns no rows, THEN the system SHALL return an empty rows array with column metadata when available.
6. IF the target table name is invalid or unsupported, THEN the system SHALL return a structured error.
7. IF query execution exceeds the configured timeout, THEN the system SHALL return `QUERY_TIMEOUT` where timeout detection is available, or `QUERY_FAILED` with a masked database error otherwise.

### Requirement 7: Read-Only SQL Safety and Error Handling

**User Story:** As a repository maintainer, I want custom SQL guarded and errors structured, so that agents do not mutate databases and can handle failures deterministically.

#### Acceptance Criteria

1. WHEN custom SQL contains multiple statements, THEN the system SHALL reject it with `SQL_MULTIPLE_STATEMENTS`.
2. WHEN custom SQL contains dangerous write or DDL keywords, THEN the system SHALL reject it with `SQL_NOT_READONLY`.
3. WHEN custom SQL starts with an unsupported SQL command, THEN the system SHALL reject it with `SQL_NOT_READONLY`.
4. WHEN SQLite PRAGMA is used, THEN the system SHALL allow only approved read-only metadata PRAGMA names.
5. WHEN unsafe PRAGMA is used, THEN the system SHALL reject it with `UNSAFE_PRAGMA`.
6. WHEN an unsupported database type is requested, THEN the system SHALL return `UNSUPPORTED_DB`.
7. WHEN a connection attempt fails, THEN the system SHALL return `CONNECTION_FAILED` without exposing full secrets.
8. WHEN database execution fails, THEN the system SHALL return `QUERY_FAILED` without fabricating results.
9. WHEN any error includes a URL or credential-like value, THEN the system SHALL mask passwords and tokens before output.

### Requirement 8: Skill Documentation, Evaluations, and CI Workflows

**User Story:** As a repository maintainer, I want docs, evals, and workflows updated with the new binary contract, so that the skill can be tested, released, and used consistently.

#### Acceptance Criteria

1. WHEN `db-explorer/SKILL.md` is updated, THEN it SHALL instruct agents to detect platform and invoke `db-explorer/bin/db-explorer-<platform>` directly.
2. WHEN `db-explorer/SKILL.md` is updated, THEN it SHALL describe the new CLI shape, JSON default, config model, and read-only guardrails.
3. WHEN reference docs are added, THEN `db-explorer/references/configuration.md` SHALL describe TOML config and precedence.
4. WHEN reference docs are added, THEN `db-explorer/references/migration-from-python.md` SHALL state that old Python CLI and `.db-explorer.json` are unsupported.
5. WHEN evals are updated, THEN `evals/db-explorer/` SHALL invoke the Go binary and validate JSON-first behavior.
6. WHEN evals are updated, THEN they SHALL preserve negative safety cases for dangerous SQL, unsafe PRAGMA, and multiple statements.
7. WHEN CI is updated, THEN `db-explorer-test.yml` SHALL run Go tests, formatting check, vet, and a build for `db-explorer-go/`.
8. WHEN release automation is updated, THEN `db-explorer-release.yml` SHALL build release archives for the supported platforms on `db-explorer-v*` tags or manual dispatch.
9. WHEN binary update automation is updated, THEN `db-explorer-update-skill.yml` SHALL build binaries into `db-explorer/bin/`, generate checksums, and commit changes when needed.
