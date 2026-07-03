# Design Refine 示例与反例

完整的正例、反例和边界案例，展示决策分析的深度和格式规范。

## 目录

- [正例：数据模型设计决策](#正例数据模型设计决策)
- [反例：被拒绝方案污染下游](#反例被拒绝方案污染下游)
- [边界案例：驳回不适用发现](#边界案例驳回不适用发现)

---

## 正例：数据模型设计决策

展示从 review 发现 → 项目探索 → db-explorer → 外部调研 → 方案分析 → Decision Record → 设计章节更新的完整流程。

### 输入

**Review 发现：**
- B1（D1 Completeness）："方案提到了用户表，但未定义表结构和字段"
- M2（D5 Blind Spots）："未考虑用户数据与现有订单表的关联方式"
- m3（D7 Optimization）："用户偏好设置可以独立成表，避免主表字段膨胀"

**design.md 涉及章节：** §Components（UserService）、§Data Flow（注册/登录流程）

### Step 1 探索结果

**项目身份：**
- Go + Gin + GORM，`cmd/server/main.go` 入口，`internal/` 下按 domain 分包
- GORM AutoMigrate 管理 schema，migration 文件在 `migrations/`

**db-explorer 输出：**
- 现有 `users` 表：`id, email, password_hash, role, created_at, updated_at`
- `role` 是 `ENUM('customer','merchant','admin')`，无默认值
- 无唯一索引在 `email` 字段上（潜在 bug）
- `orders` 表通过 `user_id BIGINT` 外键关联

### 3a. 外部调研

搜索 `"golang user preference storage pattern"` → 找到 `gorm.io/gorm` 的 JSON field 支持可直接存偏好为 JSON 列。

搜索 `"user profile vs user table database design"` → 社区共识：偏好/ profile 应独立表，避免 `users` 表字段膨胀到 50+ 列。

### 3b. 方案分析

**方案 A: 所有字段放 users 表**
- 优点：查询简单，无 JOIN（0.5 人日）
- 缺点：用户表字段膨胀，`role` ENUM 扩展需 migration
- 架构影响：users 表变成"上帝表"，后续加字段越来越重
- 风险：低（当前字段少，短期没问题）

**方案 B: users 核心字段 + user_profiles 扩展表（1:1）**
- 优点：主表轻量，扩展字段隔离，符合归一化（1.5 人日）
- 缺点：查询需要 JOIN，但 GORM Preload 可解决
- 架构影响：遵循项目现有的 domain 分包模式
- 外部参考：Rails Devise + Profile 模式，成熟方案

**推荐：** 方案 B。理由：①项目已按 domain 分包，扩展表放在 `internal/user/` 下自然；②`role` ENUM 改为 VARCHAR + 应用层校验（顺便修 email 唯一索引缺失的 bug）；③社区共识明确——主表字段超过 15 列就该拆分。

### Step 4 写入

**Decision Record（叙述体，只写"为什么"）：**

```markdown
### Decision: 用户数据模型设计

**来源**: Review B1 (D1), M2 (D5), m3 (D7)

**决定**: users 表保留核心鉴权字段（id, email, password_hash, role），
扩展信息放入 user_profiles 表（1:1 关联）。

**理由**: 社区共识——主表超过 15 列应拆分，避免"上帝表"。
GORM JSON field 可替代独立偏好表，但权衡后独立表更符合项目 domain 分包模式。
同步修复 email 无唯一索引的 bug。

**约束**: role 从 ENUM 改为 VARCHAR(32) + 应用层校验；
GORM Preload 预加载 profile 避免 N+1；
email 加 UNIQUE 索引。
```

**设计章节更新：**

§Components → UserService 下新增 `UserProfile` 结构体定义和 GORM 关系：
```go
type User struct {
    ID           uint   `gorm:"primaryKey"`
    Email        string `gorm:"uniqueIndex;not null"`
    PasswordHash string `gorm:"not null"`
    Role         string `gorm:"type:varchar(32);default:'customer'"`
    Profile      UserProfile `gorm:"foreignKey:UserID"`
}

type UserProfile struct {
    ID        uint   `gorm:"primaryKey"`
    UserID    uint   `gorm:"uniqueIndex"`
    AvatarURL string
    Nickname  string
    Bio       string `gorm:"type:text"`
}
```

> 已更新 §Components、§Data Flow

---

## 反例：被拒绝方案污染下游

展示一个**写坏了**的 Decision Record，以及它对 `spec-plan` 的实际影响。

### ❌ 坏 Decision Record

```markdown
### Decision: 缓存策略

方案 A: Redis 缓存
- Redis Cluster 3 节点，docker-compose 部署
- key 格式: `cache:user:{id}`，TTL 3600s
- 使用 go-redis/v9 客户端，连接池 10
- Redis Sentinel 做高可用

方案 B: 内存缓存（选择 ✅）
- sync.Map 实现 LRU，max 10000 条
- 重启即失效，单机部署无高可用

决定: 选方案 B，V1 先简单。
```

### spec-plan 读到了什么

spec-plan 会读到 "Redis Cluster 3 节点"、"docker-compose 部署"、"go-redis/v9 客户端"、"Redis Sentinel 高可用"——这些都是**被拒绝**的方案 A 的实现细节，但模型无法按 ✓/✗ 标记过滤。

**实际 spec-plan 输出片段：**

```markdown
### Requirement 3: 缓存层
- 3.1 部署 Redis Cluster（3 节点）
- 3.2 配置 go-redis/v9 客户端连接池...
- 3.3 设置 Redis Sentinel 高可用...
```

这就是污染——spec-plan 把被拒绝方案的细节当成了设计的一部分。

### ✅ 正确写法

```markdown
### Decision: 缓存策略

**来源**: Review M4 (D6)

**决定**: V1 使用内存缓存（sync.Map + LRU），不引入外部缓存。

**理由**: Redis 方案被拒绝——V1 单机部署且数据量 < 1 万条，
引入 Redis Cluster 的运维成本远超收益。V2 如需要横向扩展再评估。

**约束**: 缓存重启即失效；LRU 最大 10000 条；不做持久化。
```

关键：被拒绝方案只说"是什么"和"为什么被拒"——`Redis` 这个词只出现一次作为标识，没有 "Cluster"、"go-redis"、"Sentinel" 等实现细节。spec-plan 读了只会知道 "Redis 被拒绝了，不做"，不会产生 Redis 任务。

---

## 边界案例：驳回不适用发现

展示如何正确处理 review 中不适用或误判的发现。

### 场景

Review 报告 M5（D5 Blind Spots）："未考虑国际化/i18n——用户界面需要支持多语言"

### 处理方式

在 Step 2 决策队列中标记为 `[驳回]`，不进入 Step 3 分析：

```
- [x] D1: 缓存策略 → 方案 B（已确认）
- [ ] D2: 错误处理策略
- [~] M5: 国际化支持 [驳回 — PRD 明确写 V1 仅支持中文，目标用户全部在中国大陆]
- [ ] D3: API 版本策略
```

**驳回理由必须具体**——不能只写 "不需要"，要引用来源（PRD、用户原话、设计目标等）。

### 如果用户不同意驳回

用户说 "等等，虽然 PRD 没写，但我觉得还是应该预留 i18n 能力" → 恢复为正常决策项，进入 Step 3 分析方案（如"前端文案抽离到 i18n 文件 vs 硬编码中文，成本差异不大但预留成本很低"）。
