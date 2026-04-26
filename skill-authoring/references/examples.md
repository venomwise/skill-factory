# Real-World Skill Examples

Examples of well-structured skills at different complexity levels.

## Example 1: Simple Skill (Instructions Only)

A skill that provides guidance without additional files or scripts.

### File: `commit-message-helper/SKILL.md`

```markdown
---
name: commit-message-helper
description: >
  Generate descriptive commit messages following conventional commit format.
  Use when the user asks for help writing commit messages, reviewing staged changes,
  or wants to follow commit message conventions.
---

# Commit Message Helper

Generate commit messages following the Conventional Commits specification.

## Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

## Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, no logic change)
- **refactor**: Code refactoring
- **test**: Adding or updating tests
- **chore**: Maintenance tasks

## Examples

**Example 1: New feature**
```
feat(auth): implement JWT-based authentication

Add login endpoint with JWT token generation.
Include token validation middleware for protected routes.
```

**Example 2: Bug fix**
```
fix(reports): correct date formatting in timezone conversion

Dates were displaying in local time instead of UTC.
Now consistently use UTC timestamps across all reports.
```

**Example 3: Multiple changes**
```
chore: update dependencies and refactor error handling

- Upgrade lodash to 4.17.21 for security patch
- Standardize error response format across all endpoints
- Add error logging middleware
```

## Guidelines

1. **Subject line**: Imperative mood, no period, max 50 chars
2. **Body**: Explain what and why, not how. Wrap at 72 chars
3. **Scope**: Optional, indicates affected component
4. **Breaking changes**: Add `BREAKING CHANGE:` in footer

## Workflow

1. Review the changes (git diff)
2. Identify the primary type (feat, fix, etc.)
3. Determine the scope (optional)
4. Write a clear subject line
5. Add body if changes need explanation
6. Add footer for breaking changes or issue references
```

**Why this works:**
- Concise and focused
- Provides clear examples
- Includes a simple workflow
- No unnecessary files

---

## Example 2: Medium Complexity (With References)

A skill that uses progressive disclosure with reference files.

### File: `api-integration/SKILL.md`

```markdown
---
name: api-integration
description: >
  Integrate with REST APIs including authentication, request handling, error recovery,
  and rate limiting. Use when working with external APIs, HTTP requests, or API client code.
---

# API Integration

Guide for integrating with REST APIs.

## Quick Start

Basic API request:

```python
import requests

response = requests.get(
    "https://api.example.com/endpoint",
    headers={"Authorization": f"Bearer {token}"}
)

if response.status_code == 200:
    data = response.json()
```

## Common Patterns

**Authentication**: See [AUTH.md](AUTH.md) for OAuth, API keys, JWT
**Error Handling**: See [ERRORS.md](ERRORS.md) for retry logic and recovery
**Rate Limiting**: See [RATE_LIMITS.md](RATE_LIMITS.md) for handling rate limits
**Testing**: See [TESTING.md](TESTING.md) for mocking and test patterns

## Best Practices

1. **Always handle errors explicitly**
   - Check status codes
   - Parse error responses
   - Implement retry logic

2. **Use timeouts**
   ```python
   response = requests.get(url, timeout=30)
   ```

3. **Log requests and responses**
   - Log request URL and method
   - Log response status and time
   - Redact sensitive data (tokens, passwords)

4. **Respect rate limits**
   - Check response headers for rate limit info
   - Implement exponential backoff
   - Cache responses when appropriate
```

### File: `api-integration/AUTH.md`

```markdown
# API Authentication Patterns

## OAuth 2.0

```python
from requests_oauthlib import OAuth2Session

client_id = "your_client_id"
client_secret = "your_client_secret"
token_url = "https://api.example.com/oauth/token"

oauth = OAuth2Session(client_id)
token = oauth.fetch_token(
    token_url,
    client_secret=client_secret
)

# Use token for requests
response = oauth.get("https://api.example.com/data")
```

## API Key

```python
import requests

headers = {
    "X-API-Key": "your_api_key"
}

response = requests.get(
    "https://api.example.com/endpoint",
    headers=headers
)
```

## JWT Bearer Token

```python
import requests

headers = {
    "Authorization": f"Bearer {jwt_token}"
}

response = requests.get(
    "https://api.example.com/endpoint",
    headers=headers
)
```

## Token Refresh Pattern

```python
def make_request_with_refresh(url, token, refresh_token):
    response = requests.get(
        url,
        headers={"Authorization": f"Bearer {token}"}
    )
    
    if response.status_code == 401:
        # Token expired, refresh it
        new_token = refresh_access_token(refresh_token)
        response = requests.get(
            url,
            headers={"Authorization": f"Bearer {new_token}"}
        )
    
    return response
```
```

**Why this works:**
- Main file provides overview and quick start
- Detailed topics in separate files
- One level of references (no deep nesting)
- Each reference file is focused and complete

---

## Example 3: Complex Skill (With Scripts)

A skill that includes utility scripts for deterministic operations.

### File: `database-schema-validator/SKILL.md`

```markdown
---
name: database-schema-validator
description: >
  Validate database schemas against expected structure, check for schema drift,
  and verify migrations. Use when working with database schemas, migrations,
  or ORM model definitions.
---

# Database Schema Validator

Validate database schemas and detect drift from expected structure.

## Quick Start

Validate current schema:

```bash
python scripts/validate_schema.py --db postgresql://localhost/mydb
```

## Workflow

Copy this checklist:

```
Schema Validation:
- [ ] Step 1: Export current schema
- [ ] Step 2: Compare with expected schema
- [ ] Step 3: Review differences
- [ ] Step 4: Generate migration (if needed)
- [ ] Step 5: Apply migration
- [ ] Step 6: Verify schema matches
```

**Step 1: Export current schema**

```bash
python scripts/export_schema.py --db postgresql://localhost/mydb --output current_schema.json
```

This creates a JSON file with all tables, columns, indexes, and constraints.

**Step 2: Compare with expected schema**

```bash
python scripts/compare_schemas.py expected_schema.json current_schema.json
```

Output shows:
- Missing tables
- Missing columns
- Type mismatches
- Missing indexes
- Missing constraints

**Step 3: Review differences**

Examine the comparison output. Determine if differences are:
- **Expected**: Schema intentionally changed
- **Drift**: Unintended changes that need correction
- **Migration needed**: Changes require a migration

**Step 4: Generate migration (if needed)**

```bash
python scripts/generate_migration.py \
  --from current_schema.json \
  --to expected_schema.json \
  --output migration.sql
```

Review `migration.sql` before applying.

**Step 5: Apply migration**

```bash
python scripts/apply_migration.py \
  --db postgresql://localhost/mydb \
  --migration migration.sql \
  --backup
```

The `--backup` flag creates a backup before applying changes.

**Step 6: Verify schema matches**

Run validation again:

```bash
python scripts/validate_schema.py --db postgresql://localhost/mydb
```

Should report: "Schema validation passed. No differences found."

## Utility Scripts

All scripts are in the `scripts/` directory:

**export_schema.py**: Export database schema to JSON
**compare_schemas.py**: Compare two schema files
**generate_migration.py**: Generate SQL migration from differences
**apply_migration.py**: Apply migration with backup
**validate_schema.py**: Validate schema against expected structure

## Error Handling

If validation fails:

1. Review the error message
2. Check if the database is accessible
3. Verify connection string format
4. Ensure user has read permissions

If migration fails:

1. Restore from backup: `python scripts/restore_backup.py`
2. Review migration SQL for syntax errors
3. Check for conflicting constraints
4. Verify data compatibility with new schema
```

### File: `database-schema-validator/scripts/validate_schema.py`

```python
#!/usr/bin/env python3
"""
Validate database schema against expected structure.
"""
import argparse
import json
import sys
from sqlalchemy import create_engine, inspect

def load_expected_schema(path):
    """Load expected schema from JSON file."""
    with open(path) as f:
        return json.load(f)

def get_current_schema(db_url):
    """Extract current schema from database."""
    engine = create_engine(db_url)
    inspector = inspect(engine)
    
    schema = {
        "tables": {}
    }
    
    for table_name in inspector.get_table_names():
        columns = inspector.get_columns(table_name)
        indexes = inspector.get_indexes(table_name)
        
        schema["tables"][table_name] = {
            "columns": {col["name"]: str(col["type"]) for col in columns},
            "indexes": [idx["name"] for idx in indexes]
        }
    
    return schema

def compare_schemas(expected, current):
    """Compare expected and current schemas."""
    differences = []
    
    # Check for missing tables
    for table in expected["tables"]:
        if table not in current["tables"]:
            differences.append(f"Missing table: {table}")
    
    # Check for missing columns
    for table in expected["tables"]:
        if table in current["tables"]:
            for column in expected["tables"][table]["columns"]:
                if column not in current["tables"][table]["columns"]:
                    differences.append(f"Missing column: {table}.{column}")
    
    return differences

def main():
    parser = argparse.ArgumentParser(description="Validate database schema")
    parser.add_argument("--db", required=True, help="Database URL")
    parser.add_argument("--expected", default="expected_schema.json", 
                       help="Expected schema file")
    args = parser.parse_args()
    
    try:
        expected = load_expected_schema(args.expected)
        current = get_current_schema(args.db)
        differences = compare_schemas(expected, current)
        
        if differences:
            print("Schema validation FAILED. Differences found:")
            for diff in differences:
                print(f"  - {diff}")
            sys.exit(1)
        else:
            print("Schema validation PASSED. No differences found.")
            sys.exit(0)
            
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(2)

if __name__ == "__main__":
    main()
```

**Why this works:**
- Clear workflow with checklist
- Scripts handle complex operations deterministically
- Scripts provide reliable output
- Error handling is explicit
- Each script has a single, clear purpose

---

## Example 4: Skill with Domain-Specific References

A skill organized by domain for efficient context loading.

### File: `company-analytics/SKILL.md`

```markdown
---
name: company-analytics
description: >
  Query and analyze company data across sales, finance, product, and marketing domains.
  Use when analyzing business metrics, generating reports, or querying company databases.
---

# Company Analytics

Access company data across multiple domains.

## Available Domains

**Sales**: Opportunities, pipeline, forecasts, accounts
→ See [reference/sales.md](reference/sales.md)

**Finance**: Revenue, expenses, billing, subscriptions
→ See [reference/finance.md](reference/finance.md)

**Product**: Usage metrics, features, adoption, performance
→ See [reference/product.md](reference/product.md)

**Marketing**: Campaigns, attribution, conversions, email
→ See [reference/marketing.md](reference/marketing.md)

## Quick Search

Find specific metrics:

```bash
# Search for revenue metrics
grep -i "revenue" reference/finance.md

# Search for conversion metrics
grep -i "conversion" reference/marketing.md

# Search for usage metrics
grep -i "usage" reference/product.md
```

## Common Queries

**Monthly Recurring Revenue (MRR)**:
See [reference/finance.md](reference/finance.md) → MRR Calculation

**Sales Pipeline Value**:
See [reference/sales.md](reference/sales.md) → Pipeline Metrics

**Feature Adoption Rate**:
See [reference/product.md](reference/product.md) → Adoption Metrics

**Campaign ROI**:
See [reference/marketing.md](reference/marketing.md) → ROI Calculation
```

### File: `company-analytics/reference/sales.md`

```markdown
# Sales Domain

## Database Tables

**opportunities**: Sales opportunities and deals
**accounts**: Customer accounts
**contacts**: Contact information
**pipeline_stages**: Sales pipeline stages

## Key Metrics

### Pipeline Value

Total value of all open opportunities:

```sql
SELECT SUM(amount) as pipeline_value
FROM opportunities
WHERE stage != 'Closed Won' AND stage != 'Closed Lost'
```

### Win Rate

Percentage of opportunities won:

```sql
SELECT 
  COUNT(CASE WHEN stage = 'Closed Won' THEN 1 END) * 100.0 / COUNT(*) as win_rate
FROM opportunities
WHERE stage IN ('Closed Won', 'Closed Lost')
```

### Average Deal Size

```sql
SELECT AVG(amount) as avg_deal_size
FROM opportunities
WHERE stage = 'Closed Won'
```

## Common Filters

**Exclude test accounts**:
```sql
WHERE account_id NOT IN (SELECT id FROM accounts WHERE name LIKE '%Test%')
```

**Current quarter**:
```sql
WHERE created_at >= DATE_TRUNC('quarter', CURRENT_DATE)
```

**By sales rep**:
```sql
WHERE owner_id = 'user_id'
```
```

**Why this works:**
- Domain-specific organization prevents loading irrelevant context
- Each domain file is self-contained
- Quick search helps find specific metrics
- Common queries provide starting points

---

## Summary

Effective skills at different complexity levels:

1. **Simple**: Single file with instructions and examples
2. **Medium**: Main file + focused reference files
3. **Complex**: Main file + references + utility scripts
4. **Domain-specific**: Organized by domain for efficient loading

Choose the appropriate complexity for your use case. Start simple and add complexity only when needed.
