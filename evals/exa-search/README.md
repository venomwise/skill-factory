# Exa Search Evaluations

Test cases for validating the exa-search skill functionality.

## Test Cases

The `test_cases.json` file contains 15 test scenarios covering:

- Basic documentation search
- Text and highlight extraction
- Domain filtering (include/exclude)
- Research mode
- Similar page discovery
- Date filtering
- Different search types (neural, keyword, magic)
- Category filtering
- Output formats (JSON, plain, URLs-only)
- Error handling
- API key failover

## Running Tests

### Prerequisites

1. Build or download the exa-search binary (see main README)

2. Configure API key in `~/.config/ai-skills/exa-search.toml`:
```toml
[[profiles]]
id = "main"
api_key = "YOUR_EXA_API_KEY"
```

Or use environment variable:
```bash
export EXA_API_KEY="YOUR_EXA_API_KEY"
```

### Manual Testing

Run individual test cases manually:

```bash
# Test: docs-basic
exa-search docs --query "telegram streaming openclaw" --num 3

# Test: docs-with-text
exa-search docs --query "model failover openclaw" --num 2 --text

# Test: search-domain-filter
exa-search search \
  --query "OpenClaw pricing API parameters" \
  --include-domains docs.openclaw.ai,openclaw.ai \
  --text --num 3

# Test: research-mode
exa-search research --query "Exa AI company overview" --num 3

# Test: similar-pages
exa-search similar --url "https://docs.openclaw.ai/channels/telegram" --num 5

# Test: date-filter
exa-search search --query "OpenClaw releases" --start-date 2026-01-01 --num 5

# Test: highlights
exa-search docs --query "authentication methods" --highlights --num 3

# Test: search-type-keyword
exa-search search --query "telegram streaming API" --type keyword --num 5

# Test: category-filter
exa-search search --query "machine learning" --category "research paper" --num 5

# Test: exclude-domains
exa-search search \
  --query "Python async programming" \
  --exclude-domains stackoverflow.com,reddit.com \
  --num 5

# Test: plain-output
exa-search docs --query "telegram streaming" --plain

# Test: urls-output
exa-search docs --query "telegram streaming" --urls

# Test: no-autoprompt
exa-search search --query "exact query match" --no-autoprompt --num 3

# Test: debug-mode
exa-search --debug search --query "test" --num 1
```

### Validation Checklist

For each test case, verify:

- [ ] Command executes without errors
- [ ] Response format matches expected structure
- [ ] Results are relevant to the query
- [ ] Special features work as expected (text extraction, domain filtering, etc.)
- [ ] Error handling is appropriate

## Test Categories

### Core Functionality (tests 1-5)
- Basic search operations
- Text extraction
- Domain filtering
- Research mode
- Similar page discovery

### Advanced Features (tests 6-10)
- Date filtering
- Highlight extraction
- Search type selection
- Category filtering
- Domain exclusion

### Output Formats (tests 11-12)
- Plain text output
- URLs-only output

### Edge Cases (tests 13-15)
- Autoprompt control
- Debug mode
- Error handling

## Expected Behavior

Each test case in `test_cases.json` includes:
- `id`: Unique test identifier
- `description`: What the test validates
- `command`: Which command to run (docs, search, research, similar)
- `args`: Command arguments
- `expected_behavior`: List of expected outcomes

## Failover Testing

To test multi-profile failover, configure multiple keys:

```toml
[[profiles]]
id = "main"
api_key = "KEY_1"

[[profiles]]
id = "backup"
api_key = "KEY_2"
```

Then run with debug mode to see failover in action:

```bash
exa-search --debug search --query "test"
```

## Notes

- Some tests require valid Exa API key(s)
- Failover tests require multiple API keys configured
- Date filter tests may need adjustment based on current date
- Similar pages test requires the reference URL to exist
- Category filter availability depends on Exa API version

## Future Improvements

- [ ] Automated test runner script
- [ ] Baseline comparison (like db-explorer)
- [ ] Performance benchmarks
- [ ] Coverage metrics
