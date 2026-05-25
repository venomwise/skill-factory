# Query Recipes

Use the selected platform binary:

```bash
bin/grok-search-<platform> <command> --query "..."
```

## 1. Breaking update check

```bash
bin/grok-search-<platform> news --query "What changed in OpenClaw recently?"
```

Use this for fresh updates, dates, and what changed.

## 2. X/Twitter and discourse scan

```bash
bin/grok-search-<platform> social --query "What are people saying about OpenClaw on X?"
```

Use this when the task is about current social/community discussion.

## 3. Fast multi-source synthesis

```bash
bin/grok-search-<platform> research --query "Summarize recent discussion around OpenClaw model failover"
```

Use this for broad situational awareness across live sources.

## 4. Official docs vs community discussion

```bash
bin/grok-search-<platform> docs-compare --query "Compare OpenClaw official docs and community discussion on Telegram streaming"
```

Expected structure in `content`:

- `Official docs:`
- `Community interpretation:`
- `Agreement/conflict:`
- `Bottom line:`

## 5. Force a single key/profile while debugging

```bash
bin/grok-search-<platform> --profile main research --query "OpenClaw Telegram streaming" --plain
```

## 6. URLs only

```bash
bin/grok-search-<platform> research --query "OpenClaw Telegram streaming" --urls
```

## 7. Override cooldown on purpose

```bash
bin/grok-search-<platform> --ignore-cooldown news --query "OpenClaw updates" --plain
```

## 8. Override endpoint temporarily

```bash
bin/grok-search-<platform> --base-url "https://api.x.ai" --model "grok-4.1-fast" research --query "latest model release discussion"
```

## 9. When to switch to Exa instead

Use `exa-search` instead when the user asks for:

- 官方文档
- 官网
- API 文档
- 参数说明
- 价格页
- canonical source retrieval
- direct page text extraction
