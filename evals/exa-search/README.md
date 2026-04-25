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

1. Set up the virtual environment:
```bash
python -m venv codex/exa-search/.venv
codex/exa-search/.venv/bin/pip install -r codex/exa-search/requirements.txt
```

2. Configure API key in `codex/exa-search/config.local.json`:
```json
{
  "profiles": [
    { "id": "main", "api_key": "YOUR_EXA_API_KEY" }
  ]
}
```

### Manual Testing

Run individual test cases manually:

```bash
# Test: docs-basic
codex/exa-search/.venv/bin/python codex/exa-search/scripts/exa_search.py docs \
  --query "telegram streaming openclaw" --num 3

# Test: docs-with-text
codex/exa-search/.venv/bin/python codex/exa-search/scripts/exa_search.py docs \
  --query "model failover openclaw" --num 2 --text

# Test: search-domain-filter
codex/exa-search/.venv/bin/python codex/exa-search/scripts/exa_search.py search \
  --query "OpenClaw pricing API parameters" \
  --include-domains docs.openclaw.ai,openclaw.ai \
  --text --num 3
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
- API key failover
- Error handling

## Expected Behavior

Each test case in `test_cases.json` includes:
- `id`: Unique test identifier
- `description`: What the test validates
- `command`: Which command to run (docs, search, research, similar)
- `args`: Command arguments
- `expected_behavior`: List of expected outcomes

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
