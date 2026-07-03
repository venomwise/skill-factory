---
name: design-refine
description: >
  针对 design-review 评审结果进行讨论分析并给出解决方案供决策的专用技能。
  Make sure to use this skill whenever the user mentions 优化设计、处理评审意见、
  根据 review 改设计、design-review 之后怎么改、review findings 怎么落地、
  设计盲点怎么处理、设计文档需要补充完善、基于反馈更新设计文档——即使他们没有
  明确说"refine"，只要涉及根据评审结果改进设计文档，都应该使用这个 skill。
  在对话中逐项讨论决策，确认后直接写入 design.md 的 Decision Record。
---

# Design Refine

## 核心理念

Use this skill after `design-review` has produced a `review.md`. Turn review findings into explicit design decisions through conversation, then write them into `design.md` — Decision Record first, then affected design sections. The design document itself is the source of truth; the conversation is the discussion record.

**Match the user's language for all output.** 如果用户用中文提问，所有输出都使用中文。

## When to use

- 用户已有 `design.md` 和对应的 `review.md`（由 `design-review` 产出），希望根据评审报告优化设计
- 需要从架构师和技术负责人视角逐项分析设计合理性、可维护性、演进成本和实现风险
- 需要比较多个技术方案，并把决策权交给架构师/技术负责人确认
- 需要在设计文档中补全 Decision Record

## When not to use

- 只是评审设计文档并找问题 → `design-review`
- 从零开始构建设计方案 → `brainstorming`
- 已经有明确设计，需要拆任务做实现计划 → `spec-plan`
- 评审的是 PRD 而非设计文档 → `prd-refine`
- 用户只要求直接修改一个小段落，且不涉及方案讨论

## Inputs / Outputs

**Inputs:**
- `design.md` 文档路径（通常在 `specs/<topic>/design.md`）
- 同目录下对应的 `review.md`
- 可选：架构师/技术负责人的补充说明、技术约束、性能要求、成本限制

**Outputs:**
- **修改 `design.md`**：Decision Record 章节追加决策（叙述体），设计章节更新以反映决策

## Workflow

### Step 1: Load context & mandatory project exploration

**This step is not optional.** Do not skip exploration even if the decisions seem straightforward — assumptions about the project are the most common source of bad architectural decisions.

**1a. Read inputs:**

Read `design.md` and `review.md`. If `design.md` already has a Decision Record section with existing entries, identify which findings have been addressed and which remain. Resume from the first unaddressed item; do not re-discuss settled decisions.

**1b. Explore project identity** (in this order):

| # | What to check | How |
|---|---|---|
| 1 | Project overview | README, CLAUDE.md, AGENTS.md |
| 2 | Tech stack & dependencies | package.json / go.mod / pom.xml / Cargo.toml / pyproject.toml |
| 3 | Entry points & module structure | main files, top-level directory layout, middleware/config |
| 4 | Recent architectural changes | `git log --oneline -10` |

Stop when you can describe the project's purpose, tech stack, architecture style, and module boundaries in 2-3 sentences.

**1c. Explore data layer when the design involves data:**

If any review finding touches data models, storage, schemas, or persistence → invoke the `db-explorer` skill to inspect the existing database. Understand tables, relationships, constraints, and existing migration patterns before proposing data-related decisions.

### Step 2: Build the decision queue

Convert review findings into a prioritized decision queue. **Group related findings** — don't mirror every review bullet.

**例如：** review 说"缺少错误处理章节"（D1）、"未考虑服务不可用的降级策略"（D5）、"可以统一错误码格式"（D7）→ 归并为 1 个"错误处理策略"决策。

**Not every finding needs a decision.** Some findings may be wrong, irrelevant to current scope, or intentionally deferred. Flag these explicitly in the queue as `[驳回]` or `[延后]` with a one-line reason, so the user can confirm. Do not silently drop findings.

Prioritize:

1. `[P0-阻断]` — 核心架构决策未确定、与现有系统冲突、缺少关键组件定义
2. `[P1-需确认]` — 组件职责划分、接口定义、数据模型设计、错误处理策略
3. `[P2-优化]` — 命名规范、配置管理、非关键性能优化

Present the queue to the user. **Wait for the user to confirm, reorder, or skip items before proceeding to Step 3.**

### Step 3: Analyze one decision at a time

**逐项讨论，一项一项来。** 每项决策在提出方案之前，必须先做外部调研：

**3a. Web search for existing solutions** (before proposing options):

For each decision, search the web for mature frameworks, libraries, or established patterns that solve the problem. Ask: "Is there already a well-maintained solution that fits?" Prioritize:

1. **Direct fit** — 已有成熟框架/库可直接引用 → 直接纳入方案，标注版本和许可
2. **Partial fit** — 有相关方案但不完全匹配 → 总结经验，结合项目实际调整
3. **No existing solution** — 确实需要自定义 → 说明为什么现有方案不适用，再提案

**搜索深度:** 每个决策至少尝试 2 组不同的关键词。第一组用技术术语（如 `"golang saga pattern rabbitmq"`），第二组用问题描述（如 `"microservice inventory reservation pattern"`）。两组都无结果才算 "no existing solution"。

This is the architect's "don't reinvent the wheel" reflex. Custom implementation is the last resort, not the default.

**3b. Present analysis:**

- **问题来源与现状**（来自 review 哪条，涉及 design.md 哪些章节，Step 1 探索得到的项目现状证据）
- **2-3 个可行方案**，优先引用外部成熟方案。每个包含：
  - 优点、缺点、风险
  - 实现成本（人日、是否需要新依赖、许可兼容性）
  - 架构影响（对现有系统的影响范围）
  - 可维护性（长期维护和技术债务）
  - 团队适配（学习成本）
- **推荐方案及理由**
- **确认问题**（让用户决策，不要替用户决定）

**Backtracking:** 如果用户在讨论中意识到之前的决策需要调整，回到那个决策项重新分析，同时更新对应的 Decision Record 和设计章节。不做线性流程的奴隶。

**分析示例（展示 3a → 3b 完整流程）：**

<details>
<summary>示例：服务间通信方式选择</summary>

**D1: 订单服务与库存服务通信 - 同步 RPC 还是异步消息？**

**问题来源：**
- Review 报告 B2（D4 Project Fit）："设计使用同步 HTTP 调用库存服务，但项目已有 RabbitMQ 基础设施"
- Review 报告 M3（D5 Blind Spots）："未考虑库存服务不可用时的降级策略"

**项目现状（来自 Step 1 探索）：**
- 项目: Go 单体服务，`clean-architecture` 目录结构，Gin + GORM 技术栈
- 消息基础设施: `go.mod` 已有 `amqp091-go`，`infra/mq/` 有完整的连接管理和消费者注册
- 现有模式: 日志、通知等横切关注点已使用 RabbitMQ 事件驱动
- 数据层: 订单表状态字段 `created/paid/shipped/completed/cancelled`，无中间状态
- 最近提交: 3 天前 `feat: add inventory audit consumer`，已注册 `inventory.decrement` 消费者

**3a. 外部调研：**

搜索 "golang saga pattern rabbitmq" → 找到 `github.com/dtm-labs/dtm`（Go Saga 框架，3.5k stars），以及 `watermill`（Go 事件驱动库，7k stars，原生支持 RabbitMQ）。两个都是 MIT 许可，与项目兼容。

**3b. 方案分析：**

方案 A: 同步 HTTP + 重试 + 熔断
- 优点：逻辑简单，事务一致性好（1 人日）
- 缺点：库存服务不可用时订单创建全部失败
- 架构影响：与现有消息驱动风格不一致
- 风险：中（可用性风险）

方案 B: 异步消息 + Watermill 库 + Saga 补偿
- 优点：解耦彻底，Watermill 封装了 RabbitMQ 的 publish/subscribe/retry 样板代码（4 人日）
- 缺点：最终一致性模型增加复杂度，需幂等处理
- 架构影响：统一为消息驱动架构，Watermill 与现有 `infra/mq/` 可共存
- 外部依赖：`github.com/ThreeDotsLabs/watermill` v1.3+, MIT 许可

**推荐方案：** 方案 B。理由：①复用已有 RabbitMQ + Watermill 减少样板代码；②订单创建天然适合异步；③方案 A 在库存不可用时完全不可用（review 指出的盲点）；④4 人日换来架构统一性。

**需要你确认：**
- 你是否同意用异步消息 + Watermill + Saga 补偿（方案 B）？
- 最终一致性的用户体验是否可接受？
- 是否要评估 dtm 作为替代（更重但 Saga 编排更完整）？

</details>

### Step 4: Write decision directly into design.md

User confirms → **immediately write into `design.md`** in two parts:

**Part A — Decision Record（叙述体，写"为什么"）：**

```markdown
### Decision: <标题>

**来源**: Review <B#/M#/m#>（<D#>），涉及 §<章节>

**决定**: <一句话：选了哪个方案>

**理由**: <2-4 句：为什么选这个、为什么拒绝其他方案>

**约束**: <决策的前提条件或接受的取舍>
```

**Part B — 设计章节更新（写"怎么做"）：**

将最终方案的具体设计写入对应章节。以 Decision Record 为索引，逐项落地：

| 决策涉及 | 更新章节 | 写入内容 |
|----------|---------|---------|
| 架构模式、通信方式、技术选型 | Architecture | 新的架构关系、技术栈变更 |
| 组件拆分、职责、接口 | Components | 组件的职责描述、输入输出接口 |
| 数据流向、处理步骤 | Data Flow | 更新后的主路径流程 |
| 错误场景、异常处理 | Error Handling | 新增/修改的错误处理策略 |
| 测试范围、测试策略 | Testing | 新增的测试要求或用例 |

写完一个决策对应的章节后，在 Decision Record 中追加一行 `> 已更新 §<章节名>`，形成交叉引用。

**格式约束（保护下游 spec-plan）：**

`design.md` 的下游消费者是 `spec-plan`——它会读取设计文档生成任务。被拒绝方案的实现细节会污染 spec-plan 的输出。

| ❌ 危险（污染下游） | ✅ 安全 |
|---|---|
| "方案 A: Redis Cluster 3 节点，key `cache:user:{id}`" | "方案 A: Redis 缓存（被拒绝——V1 不需要，运维成本高）" |
| "方案 B: LRU map + concurrentHashMap，max 10000" | "决定: 内存缓存。V1 单机部署，重启失效可接受。" |
| 被拒绝方案的库名、端点、配置值、表结构 | 被拒绝方案只说"是什么"和"为什么被拒" |

**规则：**
- Decision Record 用**叙述段落**，不用对比表格
- 被拒绝方案**不给实现级细节**
- 最终方案的实现细节**写在 Architecture/Components/Data Flow 章节**，不在 Decision Record 展开
- 无对比时只需一行："决定: X。理由: ..."

**示例：**

```markdown
### Decision: 服务间通信

**来源**: Review B2 (D4), M3 (D5)

**决定**: 使用异步消息 + Saga 补偿模式，复用项目已有的 RabbitMQ 基础设施。

**理由**: 同步 HTTP 方案（被拒绝）与项目现有消息驱动架构不一致，且库存服务不可用时订单创建全部失败。异步方案统一了架构风格，4 人日投入可接受。

**约束**: 消息需幂等处理；接受最终一致性。
```

Then move to the next undecided item.

### Step 5: Final consistency check

All decisions confirmed → do a final pass over `design.md`:

- 每条 Decision Record 在其引用的设计章节中都有对应内容（交叉引用完整）
- 新内容不与已有 Decision Record 冲突

Summarize what was changed and where. Recommend next step: if the design is now complete, proceed to `spec-plan`.

## Guardrails

- **Not a checklist.** 评审说"缺错误处理"不等于复制粘贴错误码列表。思考策略和范围。
- **User decides.** 提供分析建议，最终决策权在用户。
- **Don't reinvent the wheel.** 提案前先搜索外部成熟方案。自定义实现是最后手段，不是默认选择。引用外部方案时标注许可兼容性。
- **Grounded in reality.** 方案基于 Step 1 探索得到的项目现状，不发明未来需求。
- **Respect existing architecture.** 优先与现有系统一致，而非技术新颖性。
- **Ask when uncertain.** 性能目标、成本约束、部署要求等只有用户知道的信息，必须问。
- **Narrative Decision Record.** 被拒绝方案不给实现细节。实现细节只出现在最终方案对应的设计章节中。

## Verification

- [ ] 每个 review 发现都已处理（决策/驳回/合并讨论）
- [ ] Decision Record 叙述体，被拒绝方案无实现细节
- [ ] 每条决策对应的设计变更已在 Architecture/Components/Data Flow 等章节落地
- [ ] 各设计章节之间无矛盾（组件名、接口名、数据格式一致）
- [ ] **下游安全检查**：用 Grep 搜索 Decision Record 中是否出现了库名、API 端点、配置值等实现细节——如有，确认它们只在最终采纳方案的章节中，不在被拒绝方案的描述中

## References

- [示例与反例](references/examples.md) — 完整正例（数据模型决策 + db-explorer + 外部调研）、反例（被拒绝方案污染 spec-plan 的实际后果）、边界案例（驳回不适用发现的处理方式）
