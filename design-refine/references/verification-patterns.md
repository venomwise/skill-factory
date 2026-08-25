# Decision Verification Patterns

在 `design.md` 的 Decision Record 范围内使用以下模式，检查被拒方案是否泄漏实现细节。

## Search patterns

| Category | Pattern | Example |
|----------|---------|---------|
| Host / port | `:[0-9]{2,5}\b` | `:6379` |
| Host / port | `localhost\|127\.0\.0\.1` | `localhost:3306` |
| SQL | `\b(CREATE TABLE\|ALTER TABLE\|column)\b` | `CREATE TABLE users` |
| Field names | `\w+_id\b` | `user_id` |
| Library | `github\.com/\S+` | `github.com/...` |
| Package | `npm:\|@\w+/\w+` | `@scope/package` |
| Version | `v[0-9]+\.[0-9]+` | `v9.0.2` |
| URL | `https?://` | `https://api.example.com` |
| Endpoint | `/api/v[0-9]` | `/api/v2/users` |
| Constant | `\b[A-Z]{3,}_[A-Z_]+\b` | `MAX_POOL_SIZE` |
| Environment | `env\.\|process\.env` | `process.env.REDIS_URL` |

## Interpretation

1. 最终采纳方案的必要约束：可以保留，但优先放入对应设计章节。
2. 被拒方案描述：删除实现细节，只保留方案名称和拒绝理由。
3. `Rejected concern`：只描述排除的能力边界。
4. `Revisit when`：只描述触发条件，不预先设计未来实现。

## Example commands

```bash
rg -n ':[0-9]{2,5}\b|localhost|127\.0\.0\.1' specs/<topic>/design.md
rg -n 'github\.com/\S+|npm:|@\w+/\w+|v[0-9]+\.[0-9]+' specs/<topic>/design.md
rg -n 'https?://|/api/v[0-9]' specs/<topic>/design.md
rg -n '\b[A-Z]{3,}_[A-Z_]+\b|env\.|process\.env' specs/<topic>/design.md
```

对每个命中项读取上下文。不要仅凭正则命中自动删除内容。
