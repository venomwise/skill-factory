---
name: hackernews-digest
description: Fetch Hacker News stories, comments, and users. Use when the user says "Hacker News", "HN", "Y Combinator", "tech news digest", "top stories on HN", or asks to pull/aggregate tech news from HN.
---

# Hacker News Digest

Fetch and digest Hacker News content via the official Firebase API and the Algolia search API.

## When to use

- User wants the top/best/new stories from HN right now.
- User wants Ask HN / Show HN / job postings.
- User wants details + comments for a specific item.
- User wants a tech-news search across HN (use `search`).
- User wants a quick "HN digest" they can skim on a phone.

## Core tool

A single Python script provides everything. No install needed (uses stdlib only):

```
python3 SKILL_DIR/scripts/hn.py <command>
```

Where `SKILL_DIR` is this skill's directory (the folder containing this `SKILL.md`, e.g. `./hackernews-digest` inside the skill-factory repo).

### Commands

| Command | Purpose | Example |
|---|---|---|
| `categories` | List available categories | `hn.py categories` |
| `top / best / new / ask / show / job` | List top N stories from that feed | `hn.py top --count 10` |
| ... + `--min-score N` | Filter by upvote score | `hn.py top --count 8 --min-score 100` |
| ... + `--summary` | Auto-generated 1-line content/context snippet | `hn.py top --summary` |
| ... + `--vibe` | Comment sentiment analysis (positive/mixed/debate, concerns) | `hn.py top --vibe` |
| ... + `--summary --vibe` | **Full digest mode** — snippet + sentiment + 1-line takeaway | `hn.py top --summary --vibe` |
| `item <id>` | Fetch a single item's details | `hn.py item 36201593` |
| `story <id> [--comments]` | Story detail + top comments | `hn.py story 36201593 --comments` |
| `user <name>` | User profile + recent submissions | `hn.py user pg` |
| `search <query> [--count N]` | Search stories across HN (uses Algolia) | `hn.py search "Rust borrow checker" --count 5` |

### Digest mode output fields

When `--summary --vibe` are active, each item shows:

- **Score + comment count** — raw engagement signal
- **📌 摘要** — short context snippet (story body for Ask/Show, Algolia community text for linked articles)
- **💬 情绪标签** — 😀正面 / 😐褒贬不一 / ⚔️争论激烈 + 评论基调细项 (N正/M争议/K负)
- **关注** — top 3 concern keywords (e.g. privacy, cost, bias, memory)
- **✏️ 一句话** — actionable takeaway sentence (值得点开 / 存在分歧 / 需要正反判断) |

## Workflow guidance

**Fast digest (default):** run `hn.py top --count 5` (or best/new per user taste) for a lightweight ranked list.

**Full decision-friendly digest (recommended):** `hn.py top --summary --vibe --count 10` — this is the mode the user asked for: each item carries a content snippet, one-line takeaway, and comment sentiment so they can pick what to dive into.

**Higher signal:** add `--min-score` to cut noise, e.g. `hn.py top --summary --vibe --count 10 --min-score 150`.

**Deeper dive:** after surfacing items, if user asks for details run `hn.py story <id> --comments` to pull the thread.

**Output formatting rule:** keep it phone-friendly — a numbered list with title, score, comment count, snippet, sentiment label, and the direct URL as a tap-able Markdown link. Do not dump raw JSON.

**Note on 摘要:** for pure linked articles (no HN body text) the "summary" is drawn from the top HN comment reaction as **context**, not the article's own abstract — tag it honestly (📌) and don't present it as the article's official summary.

## Example output (top 3)

```
1. [12345] Some headline | score:615 | comments:69 | by:user
   https://link.to/article
   📌 摘要: short context snippet...
   💬 😀 正面占主导 (8正/7争议/0负) | 关注: bias, memory
   ✏️ 一句话: 社区反响正面，值得点开。
```

## Notes / limits

- Parallel fetching: comments/summaries are fetched concurrently (8 workers); ~10 items with full analysis takes ~1min. For pure `top` (metadata only) it's fast (seconds).
- Deep threads: `story --comments` caps at ~20 top-level comments.
- Sentiment is heuristic keyword-based (positive/debate/negative lexicons). Good for pointing at *which* stories have arguments, not a substitute for reading.
- `new` feed moves fast and is noisy; `top`/`best` give the best digest signal.
