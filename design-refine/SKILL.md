---
name: design-refine
description: >
  针对 design-review 评审结果进行讨论分析并给出解决方案供决策的专用技能。
  Make sure to use this skill whenever the user mentions 优化设计、处理评审意见、
  根据 review 改设计、design-review 之后怎么改、review findings 怎么落地、
  设计盲点怎么处理、设计文档 needs 补充完善、基于反馈更新设计文档——即使他们没有
  明确说 "refine"，只要涉及根据评审结果改进设计文档，都 should 使用这个 skill。
  在对话中逐项讨论决策，确认后直接写入 design.md 的 Decision Record。
---

# Design Refine

在 `design-review` 生成 `review.md` 后使用此技能，将 review 结果转化为明确的设计决策，并写入 `design.md`。
先记录决策，再更新受影响的设计章节。

**Match the user's language for all output.** If 用户用中文提问，all 输出都使用中文。

## When to use

- 用户已有 `design.md` 和对应的 `review.md`（由 `design-review` 产出），希望根据评审报告优化设计
- Need 从架构师和技术负责人视角逐项分析设计合理性、可维护性、演进成本和实现风险
- Need 比较多个技术方案，并把决策权交给架构师/技术负责人确认
- Need 在设计文档中补全 Decision Record

## When not to use

- 只是评审设计文档并找问题 → `design-review`
- 从零开始构建设计方案 → `brainstorming`
- 已经有明确设计，need 拆任务做实现计划 → `spec-plan`
- 评审的是 PRD not 设计文档 → `prd-refine`
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

If `brainstorming` 或者 `design-review` 和当前的 `design-refine` 都在同一个上下文中，或者当前已知的上下文已经十分清晰的话 can 暂时省略该步骤。

**1a. Read inputs:**

查阅 `design.md` 和 `review.md` 两份文档。

If `design.md` 中已包含 `Decision Record` 部分并且存在决策记录，must 确认哪些已被处理，哪些尚未处理。从首个未处理的记录开始继续，don't 重复讨论已确定的决策。

**1b. Explore project identity** (按此顺序):

| # | What to check | How |
|---|---|---|
| 1 | 项目概况 | README, CLAUDE.md, AGENTS.md |
| 2 | 技术栈与依赖项 | package.json / go.mod / pom.xml / Cargo.toml / pyproject.toml |
| 3 | 入口点与模块结构 | 主入口文件、顶层目录布局、中间件/配置 |
| 4 | 近期架构变更 | `git log --oneline -10` |

When 你能用 2-3 句话描述项目目的、技术栈、架构风格和模块边界时即可 STOP。

**1c. Explore data layer when the design involves data:**

IF 任何 review 发现涉及数据模型、存储、schema 或持久化，THEN 调用 `db-explorer` skill 检查现有数据库。在提出数据相关决策之前，理解表结构、关系、约束和现有的 migration 模式。

**1d. Explore existing integration patterns when the design touches external systems:**

IF 决策涉及外部 API、第三方服务、跨系统协议、消息中间件等边界，THEN 先在项目内查找已有的集成模式（HTTP client 封装、重试/熔断策略、消息消费者注册、认证方式等），理解现有约定。避免提出与既有集成风格冲突的新方案。

### Step 2: Build the decision queue

将 review 的结果转化为优先排序的决策队列。 

**Group related findings** — don't mirror 每个 review 项。

**Example:** review 说"缺少错误处理章节"（D1-Completeness）、"未考虑服务不可用的降级策略"（D5-Blind Spots）、"可统一错误码格式"（D7-Optimization）→ 归并为 1 个"错误处理策略"决策。

**Not every finding needs a decision.** 有些发现可能是错误的、与当前范围无关、或故意延后的。在队列中明确标记为 `[驳回]` 或 `[延后]` 并附上一行理由，让用户确认。Not 默不作声地忽略发现！

驳回理由必须 **引用具体来源**（PRD 原文、用户既定约束、设计目标等），don't only write "不需要"。用户可能不同意驳回并要求恢复讨论——遇到复杂场景（合并多项发现、驳回后恢复）参见 [references/examples.md](references/examples.md)。

Prioritize:

1. `[P0-阻断]` — 核心架构决策未确定、与现有系统冲突、缺少关键组件定义
2. `[P1-需确认]` — 组件职责划分、接口定义、数据模型设计、错误处理策略
3. `[P2-优化]` — 命名规范、配置管理、非关键性能优化

**Present the queue to the user.**

### Step 3: Analyze one decision at a time

**逐项讨论，一项一项来。** 每项决策按此顺序：

**3a. Identify and clarify premises:**

每个决策开始前，识别其依赖的前提条件。**前提问题不依赖外部调研结果，而是决策本身的性质决定的。**

| 前提类型 | 例子 | 如何识别 |
|---------|------|---------|
| 部署环境 | 单机/多副本/Serverless/边缘节点 | 影响缓存、存储、会话管理方案 |
| 性能量级 | QPS、数据规模、延迟要求、并发数 | 影响同步/异步、单库/分片、连接池大小 |
| 团队能力 | 技术栈熟悉度、运维能力、团队规模 | 影响新技术引入可行性 |
| 业务优先级 | 快速上线 vs 长期维护 vs 成本优先 | 影响 MVP vs 完整架构 |
| 成本约束 | 预算、机器规模、云服务额度 | 影响自建 vs 托管服务 |
| 合规/安全 | 数据主权、审计要求、加密标准 | 影响存储位置、传输方式 |

**检查前提来源（按此顺序）：**
1. **Step 1 探索已确认哪些？** README、git log、package.json、现有架构证据
2. **design.md 已写明哪些？** Goals、Constraints、Background 章节
3. **还有哪些只有用户知道？** → 列出具体问题，先问用户

**如果所有前提都已明确（来自 Step 1 或 design.md），直接进入 3b。**

**如果存在未明确的关键前提，先停下来问：**

> 这个决策依赖以下前提信息：
> 1. <具体问题 1>（影响 <方案方向 A vs B>）
> 2. <具体问题 2>（影响 <成本/复杂度差异>）
> 
> 请明确后我再继续分析方案。

等用户回答后，再进入 3b/3c。**不要在前提未明确时就提出方案，然后附带提问**——那会导致方案方向错误。

**3b. Web search for existing solutions** (内部决策可跳过):

对每个决策，在网上搜索成熟的框架、库或既定模式来解决问题。Ask: **"Is there already a well-maintained solution that fits?"** Priority：

1. **Direct fit** — 已有成熟框架/库可直接引用 → 直接纳入方案，annotate 版本和许可
2. **Partial fit** — 有相关方案但不完全匹配 → summarize 经验，结合项目实际调整
3. **No existing solution** — 确实 need 自定义 → explain 为什么现有方案 not applicable，then 提案

**搜索深度:** 每个决策至少尝试 2 组不同的关键词。第一组用技术术语（如 `"golang saga pattern rabbitmq"`），第二组用问题描述（如 `"microservice inventory reservation pattern"`）。两组都无结果才算 **"no existing solution"**。

**豁免条件（跳过 3b 直接进入 3c）：** When 决策**完全在项目内部约定范围内**——无外部技术选型、无跨系统协议、无第三方集成——for example：
- 内部模块/包命名与目录结构
- 既有技术栈内的组件职责边界与接口划分
- 项目自有数据模型的字段命名和取值规范
- 已确定框架内部的配置组织方式

这类决策的答案来自项目自身约定 not 外部生态，外部搜索无信号。跳过时 **MUST 在 3c 显式声明"内部决策，跳过外部调研"**，avoid 与"忘了搜"混淆。

这是架构师的 **"不要重复造轮子"** 反射。**Custom implementation is the last resort, not the default.**

**3c. Present analysis:**

- **Source and current state**（来自 review 哪条，涉及 design.md 哪些章节，Step 1 探索得到的项目现状证据）
- **可行方案（通常 2-3 个）**，prioritize 引用外部成熟方案。简单决策可能只有 1 个成熟方案，复杂决策可适当增加对比项。Each includes：
  - 优点、缺点、风险
  - 实现成本（是否需要新依赖、兼容性）
  - 架构影响（对现有系统的影响范围）
  - 可维护性（长期维护和技术债务）
  - 团队适配（学习成本）
- **Recommended solution and rationale**

**Backtracking:** IF 用户在讨论中意识到之前的决策需要调整，THEN 回到那个决策项重新分析

**Analysis skeleton (3a → 3b → 3c minimal structure):**

```
D<n>: <决策标题>
Source: Review <B#/M#> ([D1-完整性/D2-可用性/D3-规范性/D4-符合项目规范/D5-盲点/D6-过度设计/D7-优化的]),  §<章节>:<行数>
Current state: <Step 1 探索证据>
3a Premises: <已明确的前提> 或 <向用户提问>
3b External research: <关键词 → 找到的方案 + 许可>（内部决策则声明 "Skip: internal convention"）
3c Options:
  Option A: <一句话概括> — 优点/缺点/成本/架构影响
  Option B: <一句话概括> — 优点/缺点/成本/架构影响
Recommendation: <A 或 B> — <理由>
```

This skeleton covers standard decision scenarios. **For complex paths** (merging multiple findings, rolling back confirmed decisions), see [references/examples.md](references/examples.md).

### Step 4: Write decision directly into design.md

User confirms → **immediately write into `design.md`** in two parts:

**Part A — Decision Record (narrative, write "why"):**

```markdown
### Decision: <标题>

**Source**: Review <B#/M#/m#> (<D#>), affects §<章节>

**Decision**: <一句话：选了哪个方案>

**Rationale**: <为什么选这个、为什么拒绝其他方案>

**Constraints**: <决策的前提条件或接受的取舍>
```

**When user's decision differs from recommendation, the Rationale section must record the user's stated constraints or preferences (ideally verbatim), not generic "user chose A" statements.** For example：

- ❌ "**Rationale**: User chose Option A."
- ✅ "**Rationale**: 团队目前无 K8s 运维经验，用户明确要求 V1 用最简部署方式（原话："先跑起来再说，别一上来就 K8s"）。方案 B 的容器编排优势在 V2 有运维团队后再评估。"

用户偏好本身是决策依据的一部分，Decision Record must 能让不在讨论现场的读者 understand "为什么最终 not 推荐方案"。

**Part B — Update design sections (write "how"):**

Write the final solution's specific design into corresponding sections. Indexed by Decision Record, implement item by item:

| Decision involves | Update section | Write content |
|----------|---------|---------|
| 架构模式、通信方式、技术选型 | Architecture | 新的架构关系、技术栈变更 |
| 组件拆分、职责、接口 | Components | 组件的职责描述、输入输出接口 |
| 数据流向、处理步骤 | Data Flow | 更新后的主路径流程 |
| 错误场景、异常处理 | Error Handling | 新增/修改的错误处理策略 |
| 测试范围、测试策略 | Testing | 新增的测试要求或用例 |

**Fallback when sections don't match:** 表格假设 design.md 已有对应章节。If 目标章节不存在（for example 项目使用不同的文档结构），do not silently create a new section，先向用户确认：

1. 现有 design.md 有哪些顶层章节？（列出）
2. 目标内容 should：(a) 新增独立章节 `## <名称>`；(b) 并入既有相近章节（指名）；(c) 用户指定其他位置？

用户确认后 then 写入。Avoid 因文档结构不一致导致 spec-plan 读不到关键决策。

写完一个决策对应的章节后，在 Decision Record 中追加一行 `> 已更新 §<章节名>`，形成交叉引用。

**Format constraints (protecting downstream spec-plan):**

`design.md` 的下游消费者是 `spec-plan` —— 它 will 读取设计文档生成任务。**Implementation details of rejected options pollute spec-plan's output**，这是本 skill 最常见也最难察觉的错误模式。下面用一个真实的污染案例来校准直觉。

**❌ Bad Decision Record (pollutes downstream):**

```markdown
### Decision: 缓存策略

Option A: Redis 缓存
- Redis Cluster 3 节点，docker-compose 部署
- key 格式: `cache:user:{id}`，TTL 3600s
- 使用 go-redis/v9 客户端，连接池 10
- Redis Sentinel 做高可用

Option B: 内存缓存（选择 ✅）
- sync.Map 实现 LRU，max 10000 条

决定: 选方案 B，V1 先简单。
```

**spec-plan reads this doc and actually outputs these tasks:**

```markdown
### Requirement 3: 缓存层
- 3.1 部署 Redis Cluster（3 节点）
- 3.2 配置 go-redis/v9 客户端连接池
- 3.3 设置 Redis Sentinel 高可用
```

spec-plan cannot 按 ✓/✗ 标记过滤 — 它读到 "Redis Cluster 3 nodes" 就当成设计的一部分，生成了本应被拒绝方案的实现任务。**This is pollution.**

**✅ Correct approach:**

```markdown
### Decision: 缓存策略

**Source**: Review M4 (D6)

**Decision**: V1 使用内存缓存（sync.Map + LRU），不引入外部缓存。

**Rationale**: Redis 方案被拒绝——V1 单机部署且数据量 < 1 万条，引入 Redis Cluster 的运维成本远超收益。V2 如 need 横向扩展再评估。

**Constraints**: 缓存重启即失效；LRU 最大 10000 条；no persistence。
```

**Key difference:** `Redis` 这个词 only 作为标识出现一次，没有 `Cluster`、`go-redis`、`Sentinel`、`docker-compose` 等实现细节。spec-plan 读完 only will 知道 "Redis was rejected, don't implement"，won't 产生 Redis 相关任务。

**Rules:**
- Decision Record uses **narrative paragraphs**, not comparison tables
- Rejected options **only state "what it is" and "why rejected"** — don't 给库名、版本、端点、配置值、表结构
- Final option's implementation details **go in Architecture/Components/Data Flow sections**, not expanded in Decision Record
- When no comparison needed, one line suffices: "Decision: X. Rationale: ..."

Then move to the next undecided item.

### Step 5: Final consistency check

All decisions confirmed → 对 `design.md` 做最终检查：

- 每条 Decision Record 在其引用的设计章节中都有对应内容（交叉引用完整）
- 新内容不与已有 Decision Record 冲突

Summarize what was changed and where. **Recommend next step:** if the design is now complete, proceed to `spec-plan`.

## Guardrails

- **Not a checklist.** 评审说"缺错误处理"不等于复制粘贴错误码列表。思考策略和范围。
- **User decides, and their reasons matter.** Provide 分析建议，最终决策权在用户。用户选择偏离推荐时，Decision Record must record 用户的理由/约束/原话，don't only write "用户选择 X"。
- **Don't reinvent the wheel.** 提案前 must 搜索外部成熟方案。自定义实现是最后手段，not the default。引用外部方案时 annotate 许可兼容性。
- **Grounded in reality.** 方案 based on Step 1 探索得到的项目现状，don't 发明未来需求。
- **Respect existing architecture.** Prioritize 与现有系统一致，not 技术新颖性。
- **Ask when uncertain.** 性能目标、成本约束、部署要求等 only 用户知道的信息，must 问。
- **Narrative Decision Record.** 被拒绝方案 don't 给实现细节。实现细节 only 出现在最终方案对应的设计章节中。

## Verification

- [ ] Each review 发现都已处理（决策/驳回/合并讨论）
- [ ] Decision Record 叙述体，被拒绝方案无实现细节
- [ ] 每条决策对应的设计变更已在 Architecture/Components/Data Flow 等章节落地
- [ ] All 设计章节之间无矛盾（组件名、接口名、数据格式一致）
- [ ] **Downstream safety check**: 用 Grep 搜索 Decision Record 中是否出现了库名、API 端点、配置值等实现细节——如有，confirm 它们 only 在最终采纳方案的章节中，not 在被拒绝方案的描述中。参考搜索模式（在 Decision Record 章节内运行）：
  - 端口/主机：`:[0-9]{2,5}\b`、`localhost|127\.0\.0\.1`
  - 表/字段名：`\b(CREATE TABLE|table|column)\b`、`\w+_id\b`
  - 库名与版本：`github\.com/\S+`、`npm:|@\w+/\w+`、`v[0-9]+\.[0-9]+`
  - URL 与端点：`https?://`、`/api/v[0-9]`、`\.\w+\.(com|io|internal)`
  - 配置常量：全大写常量 `\b[A-Z]{3,}_[A-Z_]+\b`、`env\.` / `process\.env`

  命中不代表出错——检查上下文：If 命中项位于"被拒绝方案"描述中，delete 该细节，only keep "是什么 + 为什么被拒"。

## References

- [边界案例集](references/examples.md) — Only covers rare/complex paths: 驳回不适用发现、多项发现合并讨论、决策回滚。**No need to read by default** — SKILL.md 主体已覆盖标准流程与最常见错误模式（Step 4 污染反例）。Only when Step 2 遇到复杂驳回、或 Step 3 need 合并/回退多项决策时按需查阅。
