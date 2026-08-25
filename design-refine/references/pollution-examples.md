# Decision Record Pollution Examples

本文件校准 Decision 如何污染下游 `spec-plan`。只在判断被拒方案细节是否泄漏时读取。

## Pollution mechanism

`spec-plan` 会把 design 中出现的实现细节视为当前设计。被拒方案的库名、部署、端点、
配置和表结构进入 Decision 后，可能被错误转化为实施任务。

## Bad example

```markdown
### DR-cache-strategy

Option A: Redis Cluster

- 使用 3 节点和 Sentinel。
- 客户端为 go-redis/v9，连接池为 10。
- key 使用 cache:user:{id}，TTL 为 3600 秒。

Option B: 进程内缓存（选择）

- 使用 sync.Map 和 LRU。
```

下游可能据此生成 Redis Cluster、Sentinel 和连接池任务，即使该方案已被拒绝。

## Correct example

```markdown
### DR-cache-strategy

**Decision**: V1 使用进程内缓存，不引入外部缓存服务。

**Rationale**: Redis 方案被拒绝，因为当前单实例和数据规模不足以抵消额外运维成本。

**Constraints**: 缓存重启后失效；横向扩展时重新评估。
```

被拒方案只保留名称和拒绝理由。最终方案的必要细节写入 Components、Data Model 或
Key Mechanisms，而不是在 Decision 中展开方案对比清单。

## Boundary rejection example

```markdown
### DR-preview-remains-stateless

**Decision**: preview 不保存服务端会话。

**Rationale**: commit 不消费 preview 结果，跨请求状态没有业务作用。

**Constraints**: commit 根据自身完整请求独立校验。

**Rejected concern**: 不增加 Redis session、claim 或租约。

**Revisit when**: preview 结果需要被后续请求消费。
```

`Rejected concern` 描述被排除的能力边界，不包含 Redis key、TTL、Lua 或部署配置。

## Signals to inspect

Decision 中出现以下内容时检查上下文：

- 库名和版本。
- 主机、端口、URL 和 API endpoint。
- 表、字段、DDL 和迁移语句。
- 配置常量、容量值和部署拓扑。
- 被拒方案的详细步骤。

命中不自动表示错误。最终采纳方案的必要约束可以保留；被拒方向只保留名称和理由。
