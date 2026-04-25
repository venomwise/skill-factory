---
name: exa-search
description: >
  Neural web search for official documentation, API references, and structured source retrieval with text extraction.
  Use for precise, source-first research when you need high-quality results or extracted page content.
  For breaking news or real-time sentiment, use grok-search instead.
---

# Exa Search

Use Exa for **source-first retrieval**.

Prefer it when the task is about:
- official documentation
- API/reference pages
- pricing/plan details
- product/company pages
- extracting the text of a page instead of just finding the link
- expanding from one canonical page to similar pages

Do **not** default to Exa for:
- breaking news
- X/Twitter chatter
- live sentiment / fast-moving discourse
- broad real-time synthesis across many fresh sources

For those, prefer `grok-search`.

## Setup

**1. Create virtual environment:**

```bash
python -m venv codex/exa-search/.venv
```

**2. Install dependencies:**

```bash
codex/exa-search/.venv/bin/pip install -r codex/exa-search/requirements.txt
```

**3. Configure API key:**

Create `codex/exa-search/config.local.json`:

```json
{
  "profiles": [
    { "id": "main", "api_key": "YOUR_EXA_API_KEY" }
  ]
}
```

For multiple keys, environment variables, or advanced configuration, see [CONFIG.md](references/CONFIG.md).

## Workflow

1. Start with `docs` for official documentation lookups.
2. Use `search --text` or `research` when you need extracted body text.
3. Restrict domains aggressively when the user wants official sources.
4. Use `similar` when you already have the best canonical page and want adjacent sources.
5. For official-doc-only work, prefer `docs` plus domain restriction over `similar`; `similar` is semantic, not source-pure.
6. Return links plus extracted evidence, not just titles.

## Commands

All commands use the virtual environment Python:

### Official docs search
```bash
codex/exa-search/.venv/bin/python codex/exa-search/scripts/exa_search.py docs \
  --query "telegram streaming openclaw"
```

### Official docs with text extraction
```bash
codex/exa-search/.venv/bin/python codex/exa-search/scripts/exa_search.py docs \
  --query "model failover openclaw" --text --num 2
```

### General source-first search
```bash
codex/exa-search/.venv/bin/python codex/exa-search/scripts/exa_search.py search \
  --query "OpenClaw Telegram streaming" --num 5
```

### Deep extraction / research
```bash
codex/exa-search/.venv/bin/python codex/exa-search/scripts/exa_search.py research \
  --query "OpenClaw model failover" --num 3
```

### Find similar pages
```bash
codex/exa-search/.venv/bin/python codex/exa-search/scripts/exa_search.py similar \
  --url "https://docs.openclaw.ai/channels/telegram" --num 5
```

## Additional Resources

- **Query examples**: See [query-recipes.md](references/query-recipes.md) for ready-made patterns
- **Configuration**: See [CONFIG.md](references/CONFIG.md) for advanced setup and failover
- **Evaluations**: See `evals/exa-search/test_cases.json` for test scenarios

## Notes

- `docs` defaults to `includeDomains=docs.openclaw.ai`
- `research` defaults to text extraction
- Output is normalized JSON for reliable consumption
- Script supports automatic failover across multiple API keys