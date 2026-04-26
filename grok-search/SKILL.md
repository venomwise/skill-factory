---
name: grok-search
description: Real-time web research and live synthesis with sources. Use when the question depends on freshness, community chatter, X/Twitter dynamics, breaking updates, or broad multi-source summaries. Prefer this over source-first search when you need fast situational awareness. Prefer exa-search instead for official docs, API references, pricing pages, canonical source retrieval, or direct page-text extraction.
---

# Grok Search

Use Grok Search for **real-time research**.

Prefer it when the task is about:
- breaking news
- X/Twitter chatter
- fast-moving narratives
- “what are people saying now?”
- quick multi-source synthesis
- comparing official claims vs community discussion

Do **not** default to Grok Search for:
- official docs lookup
- API/reference pages
- pricing / plan details
- direct page text extraction
- canonical source retrieval

For those, prefer `exa-search`.

## Setup

Requires Python 3 (no external dependencies).

**Cross-agent installation**: This skill is installed in `~/.AI-Skills/grok-search/` and linked to both pi and codex.

**API key configuration** (choose one):
- Shared across agents: `~/.config/ai-skills/grok-search.json`
- Skill-specific: `~/.AI-Skills/grok-search/config.local.json`

See `config.example.json` for template.

## Workflow

1. Use `--mode news` for fresh updates
2. Use `--mode social` for X/Twitter and discourse-heavy prompts
3. Use `--mode research` for broad multi-source synthesis
4. Use `--mode docs-compare` for official claims plus community interpretation
5. Use `--plain` for human-readable terminal output

## Config

Key resolution order:
1. `--api-key` → 2. `GROK_API_KEY` → 3. `GROK_API_KEYS` → 4. `~/.config/ai-skills/grok-search.json` → 5. `config.local.json` → 6. `config.json`

**Recommended**: Use `~/.config/ai-skills/grok-search.json` for cross-agent shared keys, or `config.local.json` for skill-specific keys.

**Multi-key failover and cooldown**: See [references/configuration.md](references/configuration.md)

## Quick Start

```bash
python3 scripts/grok_search.py --mode news --query "Latest updates on X"
```

**More examples and query patterns**: See [references/query-recipes.md](references/query-recipes.md)

## Notes

- Optimized for real-time research and breadth, not canonical-source purity
- `docs-compare` mode separates official facts from community interpretation
- For official docs, API references, or pricing pages, use `exa-search` instead