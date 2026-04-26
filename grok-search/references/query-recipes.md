# Query Recipes

## 1. Breaking update check

```bash
python3 scripts/grok_search.py --mode news --query "What changed in OpenClaw recently?"
```

## 2. X/Twitter and discourse scan

```bash
python3 scripts/grok_search.py --mode social --query "What are people saying about OpenClaw on X?"
```

## 3. Fast multi-source synthesis

```bash
python3 scripts/grok_search.py --mode research --query "Summarize recent discussion around OpenClaw model failover"
```

## 4. Official docs vs community discussion

```bash
python3 scripts/grok_search.py --mode docs-compare --query "Compare OpenClaw official docs and community discussion on Telegram streaming"
```

Expected structure in `content`:
- `Official docs:`
- `Community interpretation:`
- `Agreement/conflict:`
- `Bottom line:`

## 5. Force a single key/profile while debugging

```bash
python3 scripts/grok_search.py --query "OpenClaw Telegram streaming" --profile main --plain
```

## 6. URLs only

```bash
python3 scripts/grok_search.py --mode research --query "OpenClaw Telegram streaming" --urls
```

## 7. Override cooldown on purpose

```bash
python3 scripts/grok_search.py --mode news --query "OpenClaw updates" --ignore-cooldown --plain
```

## 8. When to switch to Exa instead

Use `exa-search` instead when the user says:
- 官方文档
- 官网
- API 文档
- 参数说明
- 价格页
- 只要源头链接，不要泛化综述