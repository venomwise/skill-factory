---
name: design-refine
description: >
  针对 design-review 评审结果进行讨论分析并给出解决方案供决策的专用技能。
  当用户提到优化设计、处理评审意见、根据 review 改设计、design-review 之后怎么改、
  review findings 怎么落地、设计盲点怎么处理、设计文档需要补充完善、基于反馈更新设计文档时，
  必须使用此技能——即使他们没有明确说 "refine"，只要涉及根据评审结果改进设计文档，都应该使用这个技能。
  在对话中逐项讨论决策，确认后直接写入 design.md 的 Decision Record。
---

# Design Refine

在 `design-review` 生成 `review.md` 后使用此技能，将 review 结果转化为明确的设计决策，并写入 `design.md`。
先记录决策，再更新受影响的设计章节。

**输出语言跟随用户。** IF 用户用中文提问，THEN 所有输出都使用中文。

## When to use

- 用户已有 `design.md` 和对应的 `review.md`（由 `design-review` 产出），希望根据评审报告优化设计
- 需要从架构师和技术负责人视角逐项分析设计合理性、可维护性、演进成本和实现风险
- 需要比较多个技术方案，并把决策权交给架构师/技术负责人确认
- 需要在设计文档中补全 Decision Record

## When not to use

- 只是评审设计文档并找问题 → `design-review`
- 从零开始构建设计方案 → `brainstorming`
- 已经有明确设计，需要拆任务做实现计划 → `spec-plan`
- 评审的是 PRD 而不是设计文档 → `prd-refine`
- 用户只要求直接修改一个小段落，且不涉及方案讨论

## Inputs / Outputs

**Inputs:**
- `design.md` 文档路径（通常在 `specs/<topic>/design.md`）
- 同目录下对应的 `review.md`
- 可选：架构师/技术负责人的补充说明、技术约束、性能要求、成本限制

**Outputs:**
- **修改 `design.md`**：Decision Record 章节追加决策（叙述体），设计章节更新以反映决策

## Workflow

### Interaction protocol

**每个需要用户输入的回合只问一个原子问题，然后 STOP 等待回答。** 原子问题只包含一个独立决策变量；可以为同一决策变量提供互斥选项，但不能同时询问多个前提、多个决策项或多个方案维度。除非用户明确要求批量问卷，否则不要列出问题清单，也不要预告后续问题。

每次收到回答后，重新评估决策队列、当前决策和剩余前提。优先询问位于依赖链最上游、答案最可能改变后续问题或淘汰方案的问题。IF 当前回答使后续问题失效，THEN 删除或改写这些问题，不要按预先生成的清单继续提问。

队列概览、现状证据和方案对比可以一次展示；本规则限制的是要求用户回答的独立问题数量，不限制只读信息的数量。

### Step 1: Load context & mandatory project exploration

IF `brainstorming` 或 `design-review` 与当前的 `design-refine` 都在同一个上下文中，或者当前已知的上下文已经十分清晰，THEN 可以暂时省略该步骤。

**1a. Read inputs:**

查阅 `design.md` 和 `review.md` 两份文档。

IF `design.md` 中已包含 Decision Record 部分并且存在决策记录，THEN 必须确认哪些已被处理，哪些尚未处理。从首个未处理的记录开始继续，不要重复讨论已确定的决策。

**1b. Explore project identity** (按此顺序):

| # | 检查项 | 方法 |
|---|---|---|
| 1 | 项目概况 | README, CLAUDE.md, AGENTS.md |
| 2 | 技术栈与依赖项 | package.json / go.mod / pom.xml / Cargo.toml / pyproject.toml |
| 3 | 入口点与模块结构 | 主入口文件、顶层目录布局、中间件/配置 |
| 4 | 近期架构变更 | `git log --oneline -10` |

STOP when 你能用 2-3 句话描述项目目的、技术栈、架构风格和模块边界。

**1c. Explore data layer when the design involves data:**

IF 任何 review 发现涉及数据模型、存储、schema 或持久化，THEN 调用 `db-explorer` skill 检查现有数据库。在提出数据相关决策之前，理解表结构、关系、约束和现有的 migration 模式。

**1d. Explore existing integration patterns when the design touches external systems:**

IF 决策涉及外部 API、第三方服务、跨系统协议、消息中间件等边界，THEN 先在项目内查找已有的集成模式（HTTP client 封装、重试/熔断策略、消息消费者注册、认证方式等），理解现有约定。避免提出与既有集成风格冲突的新方案。

### Step 2: Build the decision queue

将 review 的结果转化为优先排序的决策队列。

**合并相关发现** —— 不要逐条照搬每个 review 项。

**Example:** review 说"缺少错误处理章节"（D1-完整性）、"未考虑服务不可用的降级策略"（D5-盲点）、"可统一错误码格式"（D7-优化）→ 归并为 1 个"错误处理策略"决策。

决策项用**标题 + 来源发现编号**标识，不使用单独的序号。例如：`错误处理策略（来源: B3+M2+m5）`——标题用于指代，来源编号保持对 review 的溯源。

**并非每个发现都需要决策。** 有些发现可能是错误的、与当前范围无关、或故意延后的。在队列中明确标记为 `[驳回]` 或 `[延后]` 并附上一行理由，让用户确认。不要默不作声地忽略发现！

驳回理由必须**引用具体来源**（PRD 原文、用户既定约束、设计目标等），不要只写"不需要"。用户可能不同意驳回并要求恢复讨论——遇到复杂场景（合并多项发现、驳回后恢复）参见 [references/examples.md](references/examples.md)。

Prioritize:

1. `[P0-阻断]` — 核心架构决策未确定、与现有系统冲突、缺少关键组件定义
2. `[P1-需确认]` — 组件职责划分、接口定义、数据模型设计、错误处理策略
3. `[P2-优化]` — 命名规范、配置管理、非关键性能优化

**P0/P1 逐项讨论。P2 只打包展示，不打包决策：** 可以一次列出所有 P2 项作为只读摘要，但不要让用户在同一轮逐项勾选。P0/P1 完成后，只问一个范围问题：是否将全部 P2 延后。IF 用户不同意全部延后，THEN 从首个 P2 开始逐项处理。最终延后的 P2 写入 `review.md` 末尾的 `## Deferred` 段落（标注 `[design-refine deferred]` 与一句原因），不进入 Decision Record。

把队列标记为**临时队列**并呈现给用户。展示队列不代表要求用户一次性确认全部条目：

- IF 队列的合并、顺序、驳回或延后依赖只有用户知道的信息，THEN 选择依赖链最上游的一个队列塑形问题，询问后 STOP
- 收到回答后重新构建并展示受影响的队列部分；不要继续使用回答前准备的问题清单
- IF 队列没有未明确的关键前提，THEN 直接进入首个 P0/P1 决策，无需额外要求用户同时确认、调整和剔除整个队列

### Step 3: Analyze one decision at a time

**逐项讨论，一项一项来。** 每项决策按此顺序：

**3a. Identify and clarify premises:**

每个决策开始前，识别其依赖的前提条件。**前提问题不依赖外部调研结果，而是由决策本身的性质决定。**

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
2. **design.md 已写明哪些？** Goals、Context、Non-Goals 等章节
3. **还有哪些只有用户知道？** → 形成未决前提集合，但不要一次展示全部问题

**如果所有前提都已明确（来自 Step 1 或 design.md），直接进入 3b。**

**如果存在未明确的关键前提，按依赖关系循环澄清：**

1. 根据最新上下文重新计算未决前提集合
2. 选择答案最可能改变后续问题或淘汰方案的上游前提
3. 只问这一个原子问题，说明它影响的方案方向，然后 STOP
4. 收到回答后回到第 1 步；不要默认其余问题仍然成立

提问格式：

> <一个具体问题>（影响 <方案方向 A vs B>）

所有关键前提明确后，再进入 3b/3c。**不要在前提未明确时就提出方案，然后附带提问**——那会导致方案方向错误。

**3b. Web search for existing solutions** (内部决策可跳过):

对每个决策，在网上搜索成熟的框架、库或既定模式来解决问题。核心问题：**"是否已有成熟且维护良好的方案可以直接适用？"** 优先级：

1. **Direct fit** — 已有成熟框架/库可直接引用 → 直接纳入方案，注明版本和许可
2. **Partial fit** — 有相关方案但不完全匹配 → 总结经验，结合项目实际调整
3. **No existing solution** — 确实需要自定义 → 解释为什么现有方案不适用，然后再提案

**搜索深度:** 每个决策至少尝试 2 组不同的关键词。第一组用技术术语（如 `"golang saga pattern rabbitmq"`），第二组用问题描述（如 `"microservice inventory reservation pattern"`）。两组都无结果才算 **"no existing solution"**。

**豁免条件（跳过 3b 直接进入 3c）：** 当决策**完全在项目内部约定范围内**——无外部技术选型、无跨系统协议、无第三方集成——例如：
- 内部模块/包命名与目录结构
- 既有技术栈内的组件职责边界与接口划分
- 项目自有数据模型的字段命名和取值规范
- 已确定框架内部的配置组织方式

这类决策的答案来自项目自身约定而非外部生态，外部搜索没有信号。跳过时**必须在 3c 显式声明"内部决策，跳过外部调研"**，避免与"忘了搜"混淆。

这是架构师的**"不要重复造轮子"**反射。**自定义实现是最后手段，不是默认选择。**

**3c. Present analysis:**

- **来源与现状**（来自 review 哪条，涉及 design.md 哪些章节，Step 1 探索得到的项目现状证据）
- **可行方案（通常 2-3 个）**，优先引用外部成熟方案。简单决策可能只有 1 个成熟方案，复杂决策可适当增加对比项。每个方案包括：
  - 优点、缺点、风险
  - 实现成本（是否需要新依赖、兼容性）
  - 架构影响（对现有系统的影响范围）
  - 可维护性（长期维护和技术债务）
  - 团队适配（学习成本）
- **推荐方案及理由**
- **确认问题**：只要求用户选择或修订当前方案，不同时追问其他前提或下一项决策

IF 用户选择非推荐方案但没有说明理由，THEN 下一回合只追问该选择背后的约束或偏好。获得理由后再进入 Step 4，不要把理由追问与下一项决策合并。

**Backtracking:** IF 用户在讨论中意识到之前的决策需要调整，THEN 回到那个决策项重新分析。修订已确认的决策时不删除原记录，而是追加 Revised 块，格式与实操见 [references/examples.md](references/examples.md)。

**Analysis skeleton (3a → 3b → 3c minimal structure):**

```
<决策标题>
Source: Review <B#/M#/m#> (D# 维度名), §<章节>:<行数>
Current state: <Step 1 探索证据>
3a Premises: <已明确的前提> 或 <向用户提问>
3b External research: <关键词 → 找到的方案 + 许可>（内部决策则声明 "Skip: internal convention"）
3c Options:
  Option A: <一句话概括> — 优点/缺点/成本/架构影响
  Option B: <一句话概括> — 优点/缺点/成本/架构影响
Recommendation: <A 或 B> — <理由>
```

此骨架覆盖标准决策场景。**复杂路径**（合并多项发现、回滚已确认决策）见 [references/examples.md](references/examples.md)。

### Step 4: Write decision directly into design.md

用户确认 → **立即写入 `design.md`**，分为两部分：

**Part A — Decision Record (narrative, write "why"):**

```markdown
### Decision: <标题>

**Source**: Review <B#/M#/m#> (<D#>), affects §<章节>

**Decision**: <一句话：选了哪个方案>

**Rationale**: <为什么选这个、为什么拒绝其他方案>

**Constraints**: <决策的前提条件或接受的取舍>
```

**当用户的决定与推荐不同时，Rationale 部分必须记录用户陈述的约束或偏好（最好逐字引用），而不是"用户选择了 A"这类空泛表述。** 例如：

- ❌ "**Rationale**: User chose Option A."
- ✅ "**Rationale**: 团队目前无 K8s 运维经验，用户明确要求 V1 用最简部署方式（原话："先跑起来再说，别一上来就 K8s"）。方案 B 的容器编排优势在 V2 有运维团队后再评估。"

用户偏好本身是决策依据的一部分，Decision Record 必须能让不在讨论现场的读者理解"为什么最终选了这个方案"。

**Part B — Update design sections (write "how"):**

将最终方案的具体设计写入对应章节。以 Decision Record 为索引，逐项落实：

| 决策涉及 | 更新章节 | 写入内容 |
|----------|---------|---------|
| 架构模式、通信方式、技术选型 | Architecture | 新的架构关系、技术栈变更 |
| 组件拆分、职责、接口 | Components | 组件的职责描述、输入输出接口 |
| 数据流向、处理步骤 | Data Flow | 更新后的主路径流程 |
| 错误场景、异常处理 | Error Handling | 新增/修改的错误处理策略 |
| 测试范围、测试策略 | Testing | 新增的测试要求或用例 |

**章节不匹配时的降级路径：** 表格假设 design.md 已有对应章节。IF 目标章节不存在（例如项目使用不同的文档结构），THEN 不要悄悄新建章节，先向用户确认：

1. 现有 design.md 有哪些顶层章节？（列出）
2. 目标内容应该：(a) 新增独立章节 `## <名称>`；(b) 并入既有相近章节（指名）；(c) 用户指定其他位置？

用户确认后再写入。避免因文档结构不一致导致 spec-plan 读不到关键决策。

写完一个决策对应的章节后，在 Decision Record 中追加一行 `> 已更新 §<章节名>`，形成交叉引用。

**Format constraints (protecting downstream spec-plan):**

`design.md` 的下游消费者是 `spec-plan` —— 它会读取设计文档生成任务。**被拒绝方案的实现细节会污染 spec-plan 的输出**，这是本 skill 最常见也最难察觉的错误模式。下面用一个真实的污染案例来校准直觉。

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

**spec-plan 读完这份文档后实际输出的任务：**

```markdown
### Requirement 3: 缓存层
- 3.1 部署 Redis Cluster（3 节点）
- 3.2 配置 go-redis/v9 客户端连接池
- 3.3 设置 Redis Sentinel 高可用
```

spec-plan 无法按 ✓/✗ 标记过滤——它读到 "Redis Cluster 3 nodes" 就当成设计的一部分，生成了本应被拒绝方案的实现任务。**这就是污染。**

**✅ Correct approach:**

```markdown
### Decision: 缓存策略

**Source**: Review M4 (D6)

**Decision**: V1 使用内存缓存（sync.Map + LRU），不引入外部缓存。

**Rationale**: Redis 方案被拒绝——V1 单机部署且数据量 < 1 万条，引入 Redis Cluster 的运维成本远超收益。V2 如需横向扩展再评估。

**Constraints**: 缓存重启即失效；LRU 最大 10000 条；无持久化。
```

**关键差异：** `Redis` 这个词只作为标识出现一次，没有 `Cluster`、`go-redis`、`Sentinel`、`docker-compose` 等实现细节。spec-plan 读完只会知道 "Redis 被拒绝，不要实现"，不会产生 Redis 相关任务。

**Rules:**
- Decision Record 使用**叙述段落**，不用对比表格
- 被拒绝方案**只写"是什么"和"为什么被拒"** —— 不给库名、版本、端点、配置值、表结构
- 最终方案的实现细节**写入 Architecture/Components/Data Flow 等设计章节**，不在 Decision Record 中展开
- 无需对比时一行即可："Decision: X. Rationale: ..."

然后进入下一个未决事项。

### Step 5: Final consistency check

所有决策已确认 → 对 `design.md` 做最终检查：

- 每条 Decision Record 在其引用的设计章节中都有对应内容（交叉引用完整）
- 新内容不与已有 Decision Record 冲突

总结改动了什么、改在哪里。**推荐下一步：** IF 本轮 refine 由 Reject 裁决触发，THEN 推荐在规划前重新运行 `/design-review`；否则设计完备时进入 `spec-plan`。

## Guardrails

- **不当清单执行器。** 评审说"缺错误处理"不等于复制粘贴错误码列表。思考策略和范围。
- **用户决策，且理由重要。** 提供分析建议，最终决策权在用户。用户选择偏离推荐时，Decision Record 必须记录用户的理由/约束/原话，不要只写"用户选择 X"。
- **不要重复造轮子。** 提案前必须搜索外部成熟方案。自定义实现是最后手段，不是默认选择。引用外部方案时注明许可兼容性。
- **立足现实。** 方案基于 Step 1 探索得到的项目现状，不要发明未来需求。
- **尊重现有架构。** 优先与现有系统一致，而不是技术新颖性。
- **不确定就问。** 性能目标、成本约束、部署信息等只有用户知道的内容，必须问。
- **单问单答。** 每个需要用户输入的回合只问一个原子问题并等待回答；回答后重新规划，不执行预制问题清单。
- **叙述体 Decision Record。** 被拒绝方案不给实现细节。实现细节只出现在最终方案对应的设计章节中。

## Verification

- [ ] 每条 review 发现都已处理（决策/驳回/合并讨论）
- [ ] Decision Record 叙述体，被拒绝方案无实现细节
- [ ] 每条决策对应的设计变更已在 Architecture/Components/Data Flow 等章节落地
- [ ] 所有设计章节之间无矛盾（组件名、接口名、数据格式一致）
- [ ] 每个需要用户输入的回合只包含一个独立决策变量，没有批量问题清单或跨决策勾选
- [ ] 后续问题基于用户最新回答重新生成；已失效的问题没有继续询问
- [ ] **下游安全检查**: 用 Grep 搜索 Decision Record 中是否出现了库名、API 端点、配置值等实现细节——如有，确认它们只出现在最终采纳方案的章节中，而不在被拒绝方案的描述中。参考搜索模式（在 Decision Record 章节内运行）：
  - 端口/主机：`:[0-9]{2,5}\b`、`localhost|127\.0\.0\.1`
  - 表/字段名：`\b(CREATE TABLE|table|column)\b`、`\w+_id\b`
  - 库名与版本：`github\.com/\S+`、`npm:|@\w+/\w+`、`v[0-9]+\.[0-9]+`
  - URL 与端点：`https?://`、`/api/v[0-9]`、`\.\w+\.(com|io|internal)`
  - 配置常量：全大写常量 `\b[A-Z]{3,}_[A-Z_]+\b`、`env\.` / `process\.env`

  命中不代表出错——检查上下文：IF 命中项位于"被拒绝方案"描述中，THEN 删除该细节，只保留"是什么 + 为什么被拒"。

## References

- [边界案例集](references/examples.md) — 只覆盖罕见/复杂路径：驳回不适用发现、多项发现合并讨论、决策回滚。**默认无需阅读** —— SKILL.md 主体已覆盖标准流程与最常见错误模式（Step 4 污染反例）。仅当 Step 2 遇到复杂驳回、或 Step 3 需要合并/回退多项决策时按需查阅。
