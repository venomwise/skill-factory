---
name: db-explorer
description: >
  Read-only exploration of PostgreSQL, MySQL, and SQLite databases. Use for listing tables,
  inspecting table structure and column types, checking primary keys / indexes / foreign keys,
  sampling table data, running read-only SQL (SELECT / WITH / SHOW / DESCRIBE / EXPLAIN / metadata PRAGMAs),
  verifying that code models match the database schema, and troubleshooting data issues.
  Only trigger when the user is explicitly in a database / SQL / schema / record / column context.
  Do NOT trigger for HTML / Markdown / UI tables or other non-database "table" scenarios.
---

# DB Explorer

Use `scripts/db_query.py` for deterministic, read-only database exploration.
Goal: quickly retrieve trustworthy table, column, constraint, and sample data — not free-form database operations.

## Use this skill when

- The user wants to see which tables exist in a database
- The user wants to inspect a table's columns, types, defaults, primary keys, indexes, or foreign keys
- The user wants to sample a few rows of data
- The user wants to run a read-only SQL query
- The user wants to verify that the database schema matches the ORM / model definitions in code
- The user is troubleshooting things like "what does this record actually look like" or "what is stored in this column"

## Do not use

- Any write operation, DDL, migration, data fix, or backfill
- Bulk data export or heavy data processing
- The word "table" clearly refers to something other than a database table
- You cannot confirm the connection target is safe, and the query may touch sensitive production data

## Capability contract

This skill depends on `scripts/db_query.py`. The only supported commands are:

- `test`: test connection
- `tables`: list tables
- `schema <table>`: show table structure
- `data <table> --limit N`: sample table data
- `query "<sql>"`: run a read-only SQL query
- `--url-env <ENV_VAR>`: read connection info from the specified environment variable

Supported output formats:

- `table` (default)
- `markdown`
- `json`
- `csv`

## Runtime prerequisites

- SQLite uses the Python standard library `sqlite3` — no extra driver needed
- PostgreSQL requires `psycopg2` in the current Python environment
- MySQL requires `mysql-connector-python` in the current Python environment

Execution rules:

- The virtual environment is **fixed** at `<skill-path>/.venv` (inside the db-explorer skill directory)
- If `<skill-path>/.venv` already exists and the required drivers are installed, reuse it — skip creation and installation
- Do **not** create a virtual environment inside the user's project directory — avoid polluting their project
- `pip install` and script execution must use the **same** Python environment

Recommended commands:

```bash
# Check whether the skill-directory venv already exists
if [ ! -d "<skill-path>/.venv" ]; then
    python -m venv "<skill-path>/.venv"
fi
. "<skill-path>/.venv/bin/activate"
python -m pip install -r <skill-path>/requirements.txt
python <skill-path>/scripts/db_query.py --db-type postgres --url "<url>" test
```

If the user is only querying SQLite, do not require PostgreSQL / MySQL drivers.

Do not claim capabilities the script does not directly implement.
For example, "show create table" is only available if you **explicitly write and execute a read-only query** that the target database supports; do not promise `SHOW CREATE TABLE` works on all databases by default — fall back to `schema`.
For SQLite `PRAGMA`: only treat read-only metadata PRAGMAs as allowed; do not treat state-changing PRAGMAs as safe.

## Required inputs

Collect only the minimum information that blocks execution. Infer what you can — do not ask mechanical follow-up questions.

### Project database config file (`.db-explorer.json`)

If `.db-explorer.json` exists in the project root, read connection info from it first. Format:

```json
{
  "default": "dev",
  "envs": {
    "dev": {
      "db-type": "postgres",
      "url": "postgresql://user:pass@dev-host:5432/mydb"
    },
    "test": {
      "db-type": "postgres",
      "url-env": "TEST_DATABASE_URL"
    },
    "pro": {
      "db-type": "mysql",
      "url": "mysql://user:pass@pro-host:3306/mydb"
    }
  }
}
```

Field reference:

| Field | Required | Description |
|-------|----------|-------------|
| `default` | No | Default environment name; used when the user does not specify one |
| `envs` | Yes | Environment configs; each key is an environment name |
| `envs.<name>.db-type` | Yes | Database type: `postgres` / `mysql` / `sqlite` |
| `envs.<name>.url` | One of two | Connection URL (direct) |
| `envs.<name>.url-env` | One of two | Environment variable name holding the URL (avoids storing passwords in plaintext) |

### Information needed

1. **Database type** — `postgres`, `mysql`, or `sqlite`
2. **Connection source** — connection URL, existing environment variable, user-supplied env var name, or SQLite file path
3. **Query target** — table list / table structure / table data / custom SQL / schema-vs-code comparison

Inference rules:

- **Check `.db-explorer.json` in the project root first.** If it exists, read database type and connection info from it:
  - If the user specifies an environment name (e.g., "check the test env", "connect to pro"), use that environment's config
  - If the user does not specify an environment, use the `default` field; if there is no `default`, list all available environments and let the user choose
  - When the config has `url-env`, pass it to the script via `--url-env`
- If the user provides a `.db` / `.sqlite` / `.sqlite3` file path, default to `sqlite` unless context clearly says otherwise
- If the user already said "show me the `users` table structure", do not ask "what do you want to look at"
- When the user provides an env var name, read its value first, then pass it via `--url`; do not print the secret
- When calling the script directly, prefer `--url-env <ENV_VAR_NAME>`
- If `.db-explorer.json` does not exist, suggest the user create one in the project root to simplify future usage; still accept a directly provided URL, env var name, or SQLite file path

## Workflow

### 1. Collect and confirm connection info

- Check whether `.db-explorer.json` exists in the project root
  - If found: read the target environment config; mention "read from `.db-explorer.json` (`<env>` environment)" in the connection summary
  - If not found: suggest "No `.db-explorer.json` found in the project root — consider creating one for faster future connections", then accept manually provided connection info
- Prefer reusing information the user has already provided
- Only ask one critical follow-up question when execution is truly blocked
- Mask passwords in connection summaries: `postgresql://user:***@host:5432/db`
- When connection info comes from a config file or env var, state the source but do not reveal the password

### 2. Verify connection first

Run a connection test before any subsequent queries:

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<connection-url>" test
```

Or from an environment variable:

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url-env DATABASE_URL test
```

If it fails:

- Preserve the key error message as-is
- Suggest a next step
- Do not speculate about table structure or data before the connection is verified

Common failure handling:

- Connection refused: check host / port / database service status
- Authentication failure: check username / password / permissions
- SQLite file not found: check whether the path is relative to the current working directory
- Missing dependency: name the missing Python package and install it in the current venv, then retry
- MySQL C extension error: the script already defaults to `use_pure=True` (pure Python); no extra handling needed

### 3. Choose the minimal operation

Prefer built-in script commands over custom metadata queries.

**List tables**

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" tables
```

**Show table structure**

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" schema <table_name>
```

**Sample table data**

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" data <table_name> --limit 10
```

**Run a read-only SQL query**

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" query "<sql>"
```

The format flag is a global argument — place it before the subcommand:

```bash
python <skill-path>/scripts/db_query.py --db-type sqlite --url "./app.db" --format markdown schema users
```

Selection strategy:

- See what's in the database: `tables`
- See how a table is defined: `schema`
- See what the data looks like: `data --limit 10`
- User gave explicit SQL, or needs JOINs / aggregation / filtering: `query`

### 4. Follow these constraints when writing SQL

- Only write read-only SQL
- For exploratory queries you write on behalf of the user, add `LIMIT 100` by default — unless the user explicitly needs more or a LIMIT is semantically inappropriate
- Do not rewrite SQL the user explicitly provided, unless you are adding an obviously safe `LIMIT` when the user's intent is clearly "just sample a few rows"
- Table names must come from user input or previously listed real table names — do not guess

### 5. Format and present results

Do not dump raw terminal output to the user. Clean it up first.

Default output rules:

- **Table list**
  - If few tables, list them directly
  - If the script returned `row_count`, include it
- **Table structure**
  - Always include: column name, type, nullable, default, primary key
  - If the script printed indexes / foreign keys, organize them into separate subsections
- **Data sample**
  - Show the first 10 rows by default
  - If there are many columns, lead with a one-sentence summary before the table
- **Custom query**
  - Lead with one sentence explaining what was queried
  - Then show the result
  - If the result is large, show only the first few rows and state it was truncated

Format suggestions:

- Small result sets: prefer `--format markdown`
- Wide or long result sets: prefer `table`
- Needs post-processing: use `json` / `csv`

### 6. Compare with code models (optional)

If the user is cross-checking ORM / model definitions:

1. Get the database definition with `schema`
2. Open the corresponding model file
3. Report only key differences:
   - Missing or extra fields
   - Type mismatches
   - Nullable inconsistencies
   - Default value or primary key constraint differences

Do not paste the entire model and entire schema side by side.

### 7. Follow-up suggestions

After a query, only suggest next steps that are directly relevant to the current context, for example:

- "Want to see sample data for this table?"
- "Want me to check the foreign key tables it references?"
- "Want me to compare the database schema with the model in code?"

## Guardrails

- This is a **read-only skill**
- Do not execute INSERT / UPDATE / DELETE / DROP / ALTER / TRUNCATE / CREATE
- Do not expose passwords or full secrets
- Do not treat the script's best-effort read-only validation as a reason to relax boundaries
- `data` defaults to a small limit for sampling; avoid unbounded queries
- When a query fails, report the facts and suggest next steps — do not fabricate results
- When a table, column, or permission does not exist, state the failure point clearly
