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

**1d. Explore existing integration patterns when the design touches external systems:**

如果决策涉及外部 API、第三方服务、跨系统协议、消息中间件等边界 → 先在项目内查找已有的集成模式（HTTP client 封装、重试/熔断策略、消息消费者注册、认证方式等），了解现有约定。避免提出与既有集成风格冲突的新方案。

### Step 2: Build the decision queue

Convert review findings into a prioritized decision queue. **Group related findings** — don't mirror every review bullet.

**例如：** review 说"缺少错误处理章节"（D1）、"未考虑服务不可用的降级策略"（D5）、"可以统一错误码格式"（D7）→ 归并为 1 个"错误处理策略"决策。

**Not every finding needs a decision.** Some findings may be wrong, irrelevant to current scope, or intentionally deferred. Flag these explicitly in the queue as `[驳回]` or `[延后]` with a one-line reason, so the user can confirm. Do not silently drop findings.

驳回理由必须**引用具体来源**（PRD 原文、用户既定约束、设计目标等），不能只写"不需要"。用户可能不同意驳回并要求恢复讨论——遇到复杂场景（合并多项发现、驳回后恢复）参见 [references/examples.md](references/examples.md)。

Prioritize:

1. `[P0-阻断]` — 核心架构决策未确定、与现有系统冲突、缺少关键组件定义
2. `[P1-需确认]` — 组件职责划分、接口定义、数据模型设计、错误处理策略
3. `[P2-优化]` — 命名规范、配置管理、非关键性能优化

**P0/P1 逐项讨论。P2 打包处理：** 一次性列出所有 P2 项（每项一行摘要），让用户勾选哪些进入 Step 3 讨论。未被勾选的 P2 项写入 `review.md` 末尾的 `## Deferred` 段落（标注 `[design-refine deferred]` 与一句原因），不进入 Decision Record。避免 P2 长尾拖长流程。

Present the queue to the user. **Wait for the user to confirm, reorder, or skip items before proceeding to Step 3.**

### Step 3: Analyze one decision at a time

**逐项讨论，一项一项来。** 每项决策在提出方案之前，必须先做外部调研：

**3a. Web search for existing solutions** (before proposing options):

For each decision, search the web for mature frameworks, libraries, or established patterns that solve the problem. Ask: "Is there already a well-maintained solution that fits?" Prioritize:

1. **Direct fit** — 已有成熟框架/库可直接引用 → 直接纳入方案，标注版本和许可
2. **Partial fit** — 有相关方案但不完全匹配 → 总结经验，结合项目实际调整
3. **No existing solution** — 确实需要自定义 → 说明为什么现有方案不适用，再提案

**搜索深度:** 每个决策至少尝试 2 组不同的关键词。第一组用技术术语（如 `"golang saga pattern rabbitmq"`），第二组用问题描述（如 `"microservice inventory reservation pattern"`）。两组都无结果才算 "no existing solution"。

**豁免条件（跳过 3a 直接进入 3b）：** 当决策**完全在项目内部约定范围内**——无外部技术选型、无跨系统协议、无第三方集成——例如：
- 内部模块/包命名与目录结构
- 既有技术栈内的组件职责边界与接口划分
- 项目自有数据模型的字段命名和取值规范
- 已确定框架内部的配置组织方式

这类决策的答案来自项目自身约定而非外部生态，外部搜索无信号。**跳过时须在 3b 显式声明"内部决策，跳过外部调研"**，避免与"忘了搜"混淆。

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

**Backtracking:** 如果用户在讨论中意识到之前的决策需要调整，回到那个决策项重新分析。**revise 时保留原 Decision Record 段落不删**，在其末尾追加：

```markdown
> **Revised** (<YYYY-MM-DD>): <新决定，一句话>
> **修订理由**: <2-3 句：为什么原决定不再适用、新信息是什么>
> （原设计章节已同步更新）
```

同步更新对应的设计章节（Architecture/Components 等）以反映新决定。**不删除历史**——下游 spec-plan 和后续 review 需要看到决策演进轨迹，避免出现"这个字段为什么以前是这样"的追溯断链。

**分析骨架（3a → 3b 最小结构）：**

```
D<n>: <决策标题>
问题来源: Review <B#/M#>（<D#>），涉及 §<章节>
项目现状: <2-3 行 Step 1 探索证据>
3a 外部调研: <关键词 → 找到的方案 + 许可>（内部决策则声明"跳过：内部约定"）
3b 方案对比:
  方案 A: <一句话> — 优点/缺点/成本/架构影响
  方案 B: <一句话> — 优点/缺点/成本/架构影响
推荐: <A 或 B> — <2-3 句理由>
需你确认: <1-3 个具体问题>
```

此骨架足以覆盖标准决策场景。若遇到需要合并多个 review 发现、回退已确认决策等复杂路径，参见 [references/examples.md](references/examples.md)。

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

**当用户决策与推荐不一致时，理由段必须记录用户提出的额外约束或偏好（尽量原话），而不是"用户选择方案 A"这种无信息表述。** 例如：

- ❌ "**理由**: 用户选择方案 A。"
- ✅ "**理由**: 团队目前无 K8s 运维经验，用户明确要求 V1 用最简部署方式（原话："先跑起来再说，别一上来就 K8s"）。方案 B 的容器编排优势在 V2 有运维团队后再评估。"

用户偏好本身是决策依据的一部分，Decision Record 必须能让不在讨论现场的读者理解"为什么最终不是推荐方案"。

**Part B — 设计章节更新（写"怎么做"）：**

将最终方案的具体设计写入对应章节。以 Decision Record 为索引，逐项落地：

| 决策涉及 | 更新章节 | 写入内容 |
|----------|---------|---------|
| 架构模式、通信方式、技术选型 | Architecture | 新的架构关系、技术栈变更 |
| 组件拆分、职责、接口 | Components | 组件的职责描述、输入输出接口 |
| 数据流向、处理步骤 | Data Flow | 更新后的主路径流程 |
| 错误场景、异常处理 | Error Handling | 新增/修改的错误处理策略 |
| 测试范围、测试策略 | Testing | 新增的测试要求或用例 |

**章节不匹配的降级路径：** 表格假设 design.md 已有对应章节。如果目标章节不存在（例如项目使用不同的文档结构），**不要静默创建新章节**，先向用户确认：

1. 现有 design.md 有哪些顶层章节？（列出）
2. 目标内容应该：(a) 新增独立章节 `## <名称>`；(b) 并入既有相近章节（指名）；(c) 用户指定其他位置？

用户确认后再写入。避免因文档结构不一致导致 spec-plan 读不到关键决策。

写完一个决策对应的章节后，在 Decision Record 中追加一行 `> 已更新 §<章节名>`，形成交叉引用。

**格式约束（保护下游 spec-plan）：**

`design.md` 的下游消费者是 `spec-plan`——它会读取设计文档生成任务。**被拒绝方案的实现细节会污染 spec-plan 的输出**，这是本 skill 最常见也最难察觉的错误模式。下面用一个真实的污染案例来校准直觉。

**❌ 坏 Decision Record（污染下游）：**

```markdown
### Decision: 缓存策略

方案 A: Redis 缓存
- Redis Cluster 3 节点，docker-compose 部署
- key 格式: `cache:user:{id}`，TTL 3600s
- 使用 go-redis/v9 客户端，连接池 10
- Redis Sentinel 做高可用

方案 B: 内存缓存（选择 ✅）
- sync.Map 实现 LRU，max 10000 条

决定: 选方案 B，V1 先简单。
```

**spec-plan 读这份文档后实际输出的任务：**

```markdown
### Requirement 3: 缓存层
- 3.1 部署 Redis Cluster（3 节点）
- 3.2 配置 go-redis/v9 客户端连接池
- 3.3 设置 Redis Sentinel 高可用
```

spec-plan 无法按 ✓/✗ 标记过滤——它读到"Redis Cluster 3 节点"就当成设计的一部分，生成了本应被拒绝方案的实现任务。这就是污染。

**✅ 正确写法：**

```markdown
### Decision: 缓存策略

**来源**: Review M4 (D6)

**决定**: V1 使用内存缓存（sync.Map + LRU），不引入外部缓存。

**理由**: Redis 方案被拒绝——V1 单机部署且数据量 < 1 万条，引入 Redis Cluster 的运维成本远超收益。V2 如需要横向扩展再评估。

**约束**: 缓存重启即失效；LRU 最大 10000 条；不做持久化。
```

关键差异：`Redis` 这个词只作为标识出现一次，没有 `Cluster`、`go-redis`、`Sentinel`、`docker-compose` 等实现细节。spec-plan 读完只会知道"Redis 被拒了，不做"，不会产生 Redis 相关任务。

**规则：**
- Decision Record 用**叙述段落**，不用对比表格
- 被拒绝方案**只说"是什么"和"为什么被拒"**——不给库名、版本、端点、配置值、表结构
- 最终方案的实现细节**写在 Architecture/Components/Data Flow 章节**，不在 Decision Record 展开
- 无对比时只需一行："决定: X。理由: ..."

Then move to the next undecided item.

### Step 5: Final consistency check

All decisions confirmed → do a final pass over `design.md`:

- 每条 Decision Record 在其引用的设计章节中都有对应内容（交叉引用完整）
- 新内容不与已有 Decision Record 冲突

Summarize what was changed and where. Recommend next step: if the design is now complete, proceed to `spec-plan`.

## Guardrails

- **Not a checklist.** 评审说"缺错误处理"不等于复制粘贴错误码列表。思考策略和范围。
- **User decides, and their reasons matter.** 提供分析建议，最终决策权在用户。用户选择偏离推荐时，Decision Record 必须记录用户的理由/约束/原话，不能只写"用户选择 X"。
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
- [ ] **下游安全检查**：用 Grep 搜索 Decision Record 中是否出现了库名、API 端点、配置值等实现细节——如有，确认它们只在最终采纳方案的章节中，不在被拒绝方案的描述中。参考搜索模式（在 Decision Record 章节内运行）：
  - 端口/主机：`:[0-9]{2,5}\b`、`localhost|127\.0\.0\.1`
  - 表/字段名：`\b(CREATE TABLE|table|column)\b`、`\w+_id\b`
  - 库名与版本：`github\.com/\S+`、`npm:|@\w+/\w+`、`v[0-9]+\.[0-9]+`
  - URL 与端点：`https?://`、`/api/v[0-9]`、`\.\w+\.(com|io|internal)`
  - 配置常量：全大写常量 `\b[A-Z]{3,}_[A-Z_]+\b`、`env\.` / `process\.env`

  命中不代表出错——检查上下文：若命中项位于"被拒绝方案"描述中，删除该细节，只保留"是什么 + 为什么被拒"。

## References

- [边界案例集](references/examples.md) — 只收罕见/复杂路径：驳回不适用发现、多项发现合并讨论、决策回滚。**默认无需阅读**——SKILL.md 主体已覆盖标准流程与最常见错误模式（Step 4 的污染反例）。仅当 Step 2 遇到复杂驳回、或 Step 3 需要合并/回退多项决策时按需查阅。
