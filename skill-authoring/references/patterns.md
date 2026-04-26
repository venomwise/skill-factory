# Common Skill Patterns

Detailed examples of effective skill patterns.

## Pattern: Multi-Step Workflow with Validation

For complex operations that require validation at each step:

```markdown
## Database Migration Workflow

Copy this checklist:

```
Migration Progress:
- [ ] Step 1: Backup current database
- [ ] Step 2: Run migration script
- [ ] Step 3: Validate schema changes
- [ ] Step 4: Run data integrity checks
- [ ] Step 5: Update application code
- [ ] Step 6: Deploy and monitor
```

**Step 1: Backup current database**

Run: `./scripts/backup.sh production`

Verify backup file exists and has non-zero size.

**Step 2: Run migration script**

Run: `./scripts/migrate.sh --dry-run`

Review the output. If acceptable, run: `./scripts/migrate.sh --execute`

**Step 3: Validate schema changes**

Run: `./scripts/validate_schema.sh`

This checks that all expected tables, columns, and constraints exist.
If validation fails, DO NOT PROCEED. Restore from backup.

**Step 4: Run data integrity checks**

Run: `./scripts/check_integrity.sh`

Verify:
- No orphaned records
- All foreign keys valid
- No null values in required fields

**Step 5: Update application code**

Update model definitions to match new schema.
Run tests: `./scripts/run_tests.sh`

**Step 6: Deploy and monitor**

Deploy application. Monitor logs for 15 minutes.
Watch for database errors or performance issues.
```

## Pattern: Conditional Execution Based on Context

Guide the agent to choose the right approach:

```markdown
## API Request Handling

1. **Determine the API version:**

   Check the endpoint URL:
   - Contains `/v1/` → Use V1 pattern
   - Contains `/v2/` → Use V2 pattern
   - No version → Use latest (V2)

2. **V1 Pattern (Legacy):**

   ```python
   import requests
   
   response = requests.post(
       "https://api.example.com/v1/endpoint",
       headers={"Authorization": f"Bearer {token}"},
       json=payload
   )
   ```

   V1 returns XML. Parse with: `import xml.etree.ElementTree as ET`

3. **V2 Pattern (Current):**

   ```python
   import requests
   
   response = requests.post(
       "https://api.example.com/v2/endpoint",
       headers={
           "Authorization": f"Bearer {token}",
           "Content-Type": "application/json"
       },
       json=payload
   )
   ```

   V2 returns JSON. Parse with: `response.json()`

4. **Error Handling:**

   For both versions:
   - 401: Token expired, refresh and retry
   - 429: Rate limited, wait 60 seconds
   - 500: Server error, retry up to 3 times
```

## Pattern: Domain-Specific Reference Organization

For skills covering multiple domains:

```markdown
# Company Data Analysis

## Quick Reference

**Sales Data**: Opportunities, pipeline, forecasts
→ See [reference/sales.md](reference/sales.md)

**Finance Data**: Revenue, expenses, billing
→ See [reference/finance.md](reference/finance.md)

**Product Data**: Usage metrics, features, adoption
→ See [reference/product.md](reference/product.md)

**Marketing Data**: Campaigns, attribution, conversions
→ See [reference/marketing.md](reference/marketing.md)

## Finding Specific Metrics

Use grep to search across domains:

```bash
# Find revenue-related metrics
grep -i "revenue" reference/*.md

# Find conversion metrics
grep -i "conversion" reference/*.md
```

## Common Queries

**Monthly recurring revenue:**
See [reference/finance.md](reference/finance.md) → MRR section

**Sales pipeline value:**
See [reference/sales.md](reference/sales.md) → Pipeline section

**Feature adoption rate:**
See [reference/product.md](reference/product.md) → Adoption section
```

## Pattern: Template with Variations

Provide a base template with clear variation points:

```markdown
## Report Generation

### Base Template

```markdown
# [Report Title]

**Date:** [YYYY-MM-DD]
**Author:** [Name]
**Status:** [Draft/Final]

## Executive Summary
[2-3 sentences summarizing key findings]

## [Section 1 - Adapt based on report type]
[Content]

## [Section 2 - Adapt based on report type]
[Content]

## Recommendations
1. [Specific, actionable recommendation]
2. [Specific, actionable recommendation]

## Next Steps
- [ ] [Action item with owner]
- [ ] [Action item with owner]
```

### Variations by Report Type

**For Analysis Reports:**
- Section 1: "Data Analysis"
- Section 2: "Key Findings"
- Include charts and visualizations

**For Status Reports:**
- Section 1: "Progress Update"
- Section 2: "Blockers and Risks"
- Include timeline and milestones

**For Incident Reports:**
- Section 1: "Incident Timeline"
- Section 2: "Root Cause Analysis"
- Include impact assessment
```

## Pattern: Feedback Loop with Script Validation

Ensure quality through automated validation:

```markdown
## Document Generation Process

1. **Generate initial draft**

   Create the document following the template in [TEMPLATE.md](TEMPLATE.md)

2. **Validate structure**

   Run: `python scripts/validate_structure.py draft.md`

   This checks:
   - All required sections present
   - Proper heading hierarchy
   - No broken links

3. **If validation fails:**

   Review error messages:
   ```
   ERROR: Missing required section "Recommendations"
   ERROR: Heading level skipped (h1 → h3)
   ERROR: Broken link: [reference](missing.md)
   ```

   Fix each issue and run validation again.

4. **Validate content**

   Run: `python scripts/validate_content.py draft.md`

   This checks:
   - Terminology consistency
   - Example format compliance
   - Code block syntax

5. **If content validation fails:**

   Fix issues and re-validate.
   DO NOT PROCEED until both validations pass.

6. **Finalize**

   Run: `python scripts/finalize.py draft.md output.md`

   This applies final formatting and generates the output.
```

## Pattern: Progressive Complexity

Start simple, reference advanced features:

```markdown
# Data Processing

## Basic Usage

For simple CSV processing:

```python
import pandas as pd

df = pd.read_csv("data.csv")
result = df.groupby("category").sum()
result.to_csv("output.csv")
```

## Advanced Features

**For large files (>1GB):**
See [LARGE_FILES.md](LARGE_FILES.md) for chunked processing

**For complex transformations:**
See [TRANSFORMATIONS.md](TRANSFORMATIONS.md) for pipeline patterns

**For data validation:**
See [VALIDATION.md](VALIDATION.md) for schema checking

**For performance optimization:**
See [PERFORMANCE.md](PERFORMANCE.md) for parallel processing
```

## Pattern: Error Recovery Workflow

Guide the agent through error scenarios:

```markdown
## API Integration Workflow

1. **Attempt API call**

   ```python
   response = requests.post(url, json=payload)
   ```

2. **Check response status**

   - **200-299**: Success → Continue to step 5
   - **400-499**: Client error → Go to step 3
   - **500-599**: Server error → Go to step 4

3. **Handle client errors (400-499)**

   - **400 Bad Request**: Check payload format
   - **401 Unauthorized**: Refresh authentication token
   - **403 Forbidden**: Check permissions
   - **404 Not Found**: Verify endpoint URL
   - **429 Too Many Requests**: Wait 60 seconds, retry

   After fixing, return to step 1.

4. **Handle server errors (500-599)**

   - **500 Internal Server Error**: Retry up to 3 times with exponential backoff
   - **503 Service Unavailable**: Wait 5 minutes, retry

   If all retries fail, log error and notify user.

5. **Process successful response**

   Parse response data and continue with workflow.
```

## Pattern: Context-Aware Instructions

Adapt behavior based on detected context:

```markdown
## Code Review

1. **Detect language and framework**

   Check file extensions and imports:
   - `.py` + `import django` → Django project
   - `.js` + `import React` → React project
   - `.java` + `@SpringBootApplication` → Spring Boot project

2. **Apply language-specific checks**

   **For Python/Django:**
   - Check for proper use of Django ORM
   - Verify migrations are included
   - Check for security issues (SQL injection, XSS)

   **For JavaScript/React:**
   - Check for proper hook usage
   - Verify prop types or TypeScript types
   - Check for accessibility issues

   **For Java/Spring Boot:**
   - Check for proper dependency injection
   - Verify exception handling
   - Check for transaction management

3. **Apply universal checks**

   Regardless of language:
   - Code readability and naming
   - Test coverage
   - Documentation completeness
   - Error handling
```

## Summary

Effective patterns:
- Break complex tasks into validated steps
- Provide clear decision points for conditional logic
- Organize domain-specific content separately
- Include feedback loops with validation
- Start simple, reference advanced features
- Handle errors explicitly
- Adapt to context when appropriate
