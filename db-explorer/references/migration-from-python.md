# Migration from Python Script

`db-explorer` now ships as prebuilt Go binaries. The old Python script and old JSON project config are no longer supported interfaces.

## Unsupported old interfaces

The following old entry points are unsupported:

```bash
python db-explorer/scripts/db_query.py ...
```

```text
.db-explorer.json
```

Do not create a Python virtual environment, install `psycopg2`, install `mysql-connector-python`, or call `scripts/db_query.py` for the new skill workflow.

## New binary shape

Detect platform and call the matching binary directly:

```bash
bin/db-explorer-<platform> <command> [flags]
```

JSON is the default output format.

## Command migration examples

### Test connection

Old:

```bash
python scripts/db_query.py --db-type sqlite --url ./app.db test
```

New:

```bash
bin/db-explorer-<platform> test --db sqlite --url ./app.db
```

### List tables

Old:

```bash
python scripts/db_query.py --db-type postgres --url-env DATABASE_URL tables
```

New:

```bash
bin/db-explorer-<platform> tables --db postgres --url-env DATABASE_URL
```

### Inspect schema

Old:

```bash
python scripts/db_query.py --db-type sqlite --url ./app.db schema users
```

New:

```bash
bin/db-explorer-<platform> schema users --db sqlite --url ./app.db
```

PostgreSQL schema-qualified names are supported:

```bash
bin/db-explorer-<platform> schema public.users --db postgres --url-env DATABASE_URL
```

### Sample data

Old:

```bash
python scripts/db_query.py --db-type mysql --url-env MYSQL_URL data users --limit 10
```

New:

```bash
bin/db-explorer-<platform> data users --limit 10 --db mysql --url-env MYSQL_URL
```

### Run read-only SQL

Old:

```bash
python scripts/db_query.py --db-type sqlite --url ./app.db query "SELECT id FROM users LIMIT 10"
```

New:

```bash
bin/db-explorer-<platform> query "SELECT id FROM users LIMIT 10" --db sqlite --url ./app.db
```

## Config migration

Old `.db-explorer.json`:

```json
{
  "default": "dev",
  "envs": {
    "dev": {
      "db-type": "postgres",
      "url-env": "DATABASE_URL"
    }
  }
}
```

New `.db-explorer.toml`:

```toml
default_profile = "dev"

[[profiles]]
id = "dev"
db = "postgres"
url_env = "DATABASE_URL"
```

## Important behavior changes

- Default output is JSON, not a human table.
- Every JSON response includes `schema_version: "1"`.
- `tables` returns low-cost `row_estimate` metadata when available; it does not run exact unbounded `COUNT(*)` by default.
- SQLite uses a pure Go driver; no Python or C toolchain setup is required.
- All structured metadata is written to stdout JSON. Diagnostics must not corrupt JSON output.
- Custom SQL remains read-only and rejects multiple statements, write/DDL keywords, and unsafe PRAGMA statements.
