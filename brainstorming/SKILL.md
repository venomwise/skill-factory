---
name: brainstorming
description: >
  把大型、抽象或包含多个独立能力的软件项目整理为本地项目方向和滚动 backlog，
  通过用户与问题发现、MVP 收敛、能力地图、目标拆分和依赖排序生成
  `projects/<project>/project.md` 与 `backlog.md`。当用户说“开发一个音乐 APP”
  “做一个图书管理系统”“帮我规划整个项目”“需求太大先拆一下”时调用；已有项目持续
  迭代时，用户说“继续上次的项目规划”“根据新需求更新项目”“调整 MVP 或优先级”
  “更新 backlog 或项目进度”也应调用。单个可独立交付的需求应使用 `clarifying`。
---

# Brainstorming Projects Into Deliverable Goals

从模糊项目想法中发现用户需求，确认 MVP，并拆成可以逐个进入 `clarifying` 的交付目标。
本技能维护项目规划基线，不替代技术设计、任务规划或开发执行。

## When to use

- 从零规划一个产品、应用或跨能力项目。
- 请求包含多个可以独立交付和验收的用户或业务结果。
- 尚未明确目标用户、核心问题、MVP、能力优先级或交付顺序。
- 继续一个 `draft` 项目，从上次未决问题恢复规划。
- 新需求、约束或项目事实变化，需要调整范围、MVP、优先级、依赖或 backlog。
- 项目已进入交付，需要继续细化 `Next / Later` 或规划下一阶段。
- 下游返回 handoff 或产生新证据，需要对账项目进度和 Current Focus。
- 用户说“继续项目规划”“接着梳理需求”“重新排一下优先级”“更新 backlog”或
  “同步一下项目进度”。

## When not to use

- 请求已经聚焦为一个可独立设计、验收和交付的目标；应使用 `clarifying`。
- 已有批准的 `design.md`，只需评审、拆任务或实施。
- 用户说“继续开发”“继续执行 tasks.md”或“恢复实施”；应使用 `spec-exec`。
- 只需同步 Jira 或 GitHub Projects；外部同步应由独立能力处理。

## Inputs / Outputs

**Inputs:**

- 用户的产品想法、业务目标、迭代变更或下游 project handoff。
- 可选的现有项目代码、文档、PRD、约束和外部项目事实。
- 更新模式下已有的 `projects/<project>/project.md` 与 `backlog.md`。

**Outputs:**

- 标准产物：`projects/<project>/project.md` 和 `projects/<project>/backlog.md`。
- 小项目降级：只生成 `project.md`，并在其中保留唯一的 `## Backlog`。
- 项目规划状态：`draft` 或 `ready`。
- 下一步：`ready` 推荐一个 `ready-for-clarifying` work item；`draft` 记录下一项待决问题。
- 不生成 `design.md`、Acceptance Criteria、`tasks.md` 或代码，不写入外部项目管理工具。

## Workflow

### STEP 1: Select mode and explore context

先根据用户表达和现有项目产物选择模式：

- `initialize`：首次规划，尚不存在对应 `project.md`。
- `resume`：继续未完成的 `draft`，从已记录的未决问题恢复。
- `replan`：新需求、约束或事实影响已有规划，或继续细化 `Next / Later` 和下一阶段。
- `reconcile`：只根据下游 handoff 和真实产物对账状态，不改变产品范围。

IF 用户要求继续或更新但未给路径，THEN 先结合当前对话和项目名检索
`projects/*/project.md`：恰好一个匹配时使用；没有或存在多个匹配时，只询问目标路径并
STOP，不猜测项目。

`resume`、`replan` 和 `reconcile` 必须完整读取 `project.md`，再按实际存储模式恢复 backlog：

- 存在 `backlog.md`：完整读取它，并确认 `project.md` 的 `## Backlog` 只保留链接。
- 不存在 `backlog.md`：读取 `project.md` 中唯一的嵌入式 `## Backlog`。
- 两处都包含 work item 或两处都没有：PAUSE，报告项目规划来源冲突或缺失，等待用户处理。

新建模式按优先级读取 README、CLAUDE.md、AGENTS.md、现有 PRD、相关代码、项目配置和
已有规划。只有理解近期变化确有必要时才查看最近 10 个提交。

按模式继续：

- `initialize`：进入 STEP 2，完成首次发现与拆分。
- `resume`：读取 Planning Status、PDR、Open Questions 和 backlog，从第一项阻止状态推进的
  未决问题继续；不重新询问已确认内容。
- `replan`：比较新信息与已确认的项目文档，建立受影响的 Goals、MVP、capability、work item、
  PDR 和 Current Focus 清单，只回到最早受影响的 STEP 3-6。
- `reconcile`：确认没有范围或优先级变化后直接进入 STEP 9；存在变化时改为 `replan`。

已有 PDR 只有在 `Revisit when` 已触发或用户明确改变决定时才重新讨论。后续调用以本地
项目文档为准，不依赖 Agent 对上次对话的记忆重建已确认内容。

能从代码、文档、数据库或历史核实的事实主动核实，不把可检索信息交给用户回答。
绿地项目没有代码事实时，明确区分用户约束、待验证假设和开放问题。

STOP when 你能说明本次模式、目标项目、已确认基线和本轮变化或待决问题。
需要拆分、恢复、状态对账方法或领域盲点时，读取
[project planning guide](references/project-planning-guide.md)。

### STEP 2: Check project granularity

执行入口关卡：

1. 请求是否包含多个可以独立产生价值的结果？不满足时，推荐使用 `clarifying`。
2. 当前主要缺口是否是项目愿景、MVP、能力范围或优先级？不满足时，推荐使用
   `clarifying`。

多个技术模块围绕同一个交付结果发生变化，不自动构成项目级需求。需要路由到
`clarifying` 时，说明证据并等待用户确认；确认后 STOP。多个项目目标存在强依赖时仍留在
本技能，依赖用于后续拆分和排序，不作为降级到单项需求的理由。

### STEP 3: Discover intent and outcomes

依次确认以下项目级信息：

- 谁会使用，当前问题或机会是什么。
- 用户和业务分别希望获得什么结果。
- 为什么现在做，哪些结果最重要。
- 明确不要做什么，现实约束是什么。
- 怎样判断项目或 MVP 成功。

每个需要用户输入的回合只问一个原子问题并等待回答。这是节奏控制，
不是总问题数上限。优先提供有代表性的选项和你的倾向，让用户容易判断取舍。

不要因为用户列出了大量功能就假设目标和优先级已经确认。项目愿景、取舍偏好、
MVP 边界和成功标准只能由用户决定。

### STEP 4: Define direction and MVP

总结已确认的问题、目标用户、期望结果、Non-Goals、约束和成功指标。
方向仍不明确时提出 1-3 个有代表性的产品方向，说明价值、成本和取舍，并给出推荐。

将候选能力划分为 `Now / Next / Later`：

- `Now`：构成首个可验证 MVP 的最小能力集合。
- `Next`：MVP 得到验证后优先扩展的能力。
- `Later`：保留方向，但当前不投入详细澄清。

向用户呈现项目方向和 MVP 边界并等待确认。用户改变目标或优先级时，更新总结后再次确认；
不要在提出确认问题后用假设代替用户回答。

### STEP 5: Build capability map and backlog

按用户或业务结果建立 capability map，再拆成 work item。不要按前端、后端、数据库或
团队结构拆分。每个 work item 使用稳定的 `WI-<semantic-name>`，并包含：

- Outcome：独立产生的用户或业务结果。
- Scope / Non-Goals：高层边界，不包含实现方案。
- Success signal：项目层面的成功信号，不是正式 Acceptance Criteria。
- Priority：`Now / Next / Later`。
- Dependencies：业务、外部或必要的技术前置。
- Status：当前流程状态。
- Open questions：留到该目标进入 `clarifying` 前后解决的问题。

work item 必须能够独立进入 `clarifying`、形成一个设计并被验收。拆分结果过大时继续拆分；
拆分结果只是代码任务时向上合并为交付结果。

### STEP 6: Sequence rolling delivery

识别依赖、可并行项和推荐顺序。只把近期准备交付的目标推进到
`ready-for-clarifying`；其余目标保留在 `candidate` 或 `deferred`，不提前澄清技术细节。

使用以下主状态：

```text
candidate -> ready-for-clarifying -> design-ready -> planned -> executing -> done
```

旁路状态使用 `blocked`、`deferred`、`cancelled`。状态变化不改变 work item ID。

进入 `ready-for-clarifying` 前执行关卡：

1. 用户、Outcome 和高层边界是否明确？不满足时，回到 STEP 3 或 STEP 5。
2. Priority 和 Dependencies 是否已确认？不满足时，回到 STEP 4 或本步骤。
3. 是否存在阻止单项设计的项目级决策？存在时，回到 STEP 3 或 STEP 4。
4. Success signal 是否足以指导后续澄清？不满足时，回到 STEP 5。

### STEP 7: Confirm persistence state

写入文件前先选择项目规划状态：

- `draft`：用户需要暂停，或仍有问题会改变项目目标、MVP、拆分或当前焦点。
- `ready`：项目级方向已经明确，可以把近期目标交给 `clarifying`。

**draft 关卡：**

1. 是否至少有一项经过用户确认的事实、范围或决策？不满足时，回到 STEP 3 继续讨论。
2. 未决问题是否包含影响范围和下一步决策？不满足时，补全后再继续。
3. 是否没有新的 `ready-for-clarifying` work item？不满足时，将其改回 `candidate` 或
   `blocked`；已有证据支持的 `design-ready` 及后续状态保持不变。
4. 是否能区分已确认内容、项目事实和待验证假设？不满足时，先探索或提问。

**ready 关卡：**

1. 本次会话是否至少问过一个项目级澄清问题？不满足时，回到 STEP 3。
2. 项目问题、用户、目标、Non-Goals、约束和成功指标是否明确？不满足时，回到 STEP 3。
3. MVP 是否能映射到完整的 `Now` work item 集合？不满足时，回到 STEP 4 或 STEP 5。
4. work item 是否按交付结果而不是技术层拆分？不满足时，回到 STEP 5。
5. 是否至少有一个近期目标为 `ready-for-clarifying`，或全部适用目标已经终态？
   不满足时，回到 STEP 6。
6. 是否仍有问题会改变项目目标、MVP、拆分或当前焦点？存在时，改用 `draft`。
7. 是否能区分项目事实、用户决定和待验证假设？不满足时，先探索或提问。

根据所选状态呈现已确认内容、未决问题、capability map、backlog、依赖和主要风险。
`ready` 同时呈现当前焦点；`draft` 呈现下一项待决问题。等待用户确认后再写文件。
反馈改变项目方向时回到 STEP 3 或 STEP 4；只改变拆分或顺序时回到 STEP 5 或 STEP 6。

### STEP 8: Write or update project plan

用户确认后：

1. 使用 [project template](assets/project-template.md) 和
   [backlog template](assets/backlog-template.md)。
2. 写入 `projects/<project>/`；project 使用 kebab-case。
3. 同名目录、文件或项目名有歧义时，先确认重命名还是沿用。
4. 默认分别写 `project.md` 和 `backlog.md`。只有 backlog 无需独立维护、同步或频繁更新时，
   才允许把唯一的 `## Backlog` 放入 `project.md`。
5. 在 `project.md` 写入 `Planning Status`。`draft` 不得包含 `ready-for-clarifying` work item；
   `ready` 必须通过 STEP 7 的 ready 关卡。
6. 重新读取实际讨论和一手项目证据并据此填写；不要凭记忆补写或编造目标、理由和优先级。
7. 项目级持久选择使用稳定 `PDR-<semantic-topic>`；同一主题变化时原位更新。
8. 更新已有规划时保留稳定 work item ID，只修改受新事实或用户决策影响的部分。
9. 本地项目规划以 `project.md` 与 backlog 为准；外部引用只记录映射，不反向覆盖内容。
10. `Spec` 只链接已经存在且能追溯到对应 `WI-*` 的 `design.md`，不创建空设计。
11. 只有本技能可以修改 backlog 的 `Status`、`Spec`、`Status evidence` 和 `Current Focus`。
12. 重新读取实际写入的 `project.md` 和 backlog，执行下面的文档语言关卡；修改语言后重新
    检查状态、引用和模板结构。

**文档语言关卡：**

1. Summary、Goals、Outcome、Success signal、Risk 和 PDR 是否写清用户或业务对象、当前问题、
   期望结果、条件和理由？IF 否，THEN 改成具体陈述，不用口号或项目管理黑话代替内容。
2. 中文正文是否自然、直接？英文只保留模板字段、技术术语、路径、`WI-*`、PDR 和状态枚举。
   IF 否，THEN 只改对应句子。
3. 是否存在逐词翻译、流程隐喻、模板化开场或空泛总结？IF 是，THEN 删除姿态层，保留实际
   业务关系；正常技术术语不做机械替换。
4. 用户、指标、优先级、依赖、已确认决策、状态、ID、路径和责任主体是否保持不变？IF 否，
   THEN 恢复原意并缩小改写范围。
5. 语言修改是否新增、删除、弱化或扩大了 Goals、Non-Goals、MVP、work item 或 PDR？IF 是，
   THEN 撤销该修改，并依据已确认内容重新表述。

IF 状态为 `draft`，THEN 告知用户下一项待决问题后 STOP，不路由到 `clarifying`。
IF 状态为 `ready`，THEN 推荐用户选择一个 `ready-for-clarifying` work item 运行
`clarifying`。本技能不自动调用下游，也不自动同步 Jira 或 GitHub Projects。

### STEP 9: Reconcile delivery status

用户要求更新进度，或下游返回 project handoff 时，读取
[project planning guide](references/project-planning-guide.md) 的 Status reconciliation：

1. 根据 STEP 1 确定应使用的 backlog，并定位稳定 `WI-*`。
2. 从真实 `design.md`、`review.md` 和 `tasks.md` 核实允许的目标状态，不只检查文件名。
3. 确认 `design.md` 的 Project Traceability 与目标 work item 一致。
4. 生成 `Status`、`Spec`、`Status evidence` 和 `Current Focus` 的最小差异。
5. 向用户呈现证据和差异，等待确认后才写回；反馈改变范围或优先级时回到 STEP 3-6。
6. 写回后重新读取实际项目文档，对本次更新字段及相邻说明执行 STEP 8 的文档语言关卡，
   再次核对状态和证据；不改写无关内容。

缺少目标状态证据时保持原状态并报告缺口。新证据明确否定当前状态时，提出转为
`blocked` 的最小差异并等待确认；不要保留已失效的 ready 状态，也不修改下游产物。

## Verification

- [ ] 文件位于 `projects/<project>/`，project 为 kebab-case。
- [ ] 已正确选择 `initialize`、`resume`、`replan` 或 `reconcile`。
- [ ] 标准模式下 `project.md` 和 `backlog.md` 均存在；降级模式只有一个 Backlog。
- [ ] `Planning Status` 为 `draft` 或 `ready`，且满足对应关卡。
- [ ] `ready` 的项目问题、用户、Goals、Non-Goals、MVP、成功指标和约束有实质内容；
  `draft` 只保存已确认内容并明确标出缺口。
- [ ] 每个 MVP capability 至少映射到一个 `Now` work item。
- [ ] work item 使用唯一、稳定的 `WI-<semantic-name>`，并描述交付结果而非代码层。
- [ ] Priority、Status 和 Dependencies 使用受支持的枚举或稳定 ID。
- [ ] `ready-for-clarifying` work item 通过全部准入关卡。
- [ ] `draft` 项目没有 `ready-for-clarifying` work item。
- [ ] `Spec` 和 External ref 只引用真实存在或真实返回的标识。
- [ ] `Status`、`Spec`、`Status evidence` 和 `Current Focus` 以 backlog 为准。
- [ ] 后续迭代只修改受新事实或用户决定影响的内容，稳定 `WI-*` 和未触发的 PDR 保持不变。
- [ ] 项目规划不包含技术方案、正式 AC、实施任务或未经确认的未来能力。
- [ ] 已重新读取实际项目文档并通过文档语言关卡，项目结果和风险不是空泛口号。
- [ ] 语言修改没有改变指标、优先级、依赖、状态、`WI-*`、PDR 或已确认范围。
- [ ] 所有普通 Markdown 行不超过 120 个 Unicode 字符。

## Safety & guardrails

- 未经确认不得写入或重排项目计划。
- 本地项目规划以项目文档为准；不得让外部平台状态静默覆盖本地内容。
- 只有本技能写入 backlog 的状态和交付证据；下游技能只返回 handoff。
- 不创建 Jira、GitHub Projects 或其他外部对象；外部同步由独立能力负责。
- 不把 brainstorming 生成的 backlog 当成已批准设计或实施任务。
- 不要求一次澄清完整 backlog；只深化近期 work item。
- 后续调用不重启完整 discovery，也不重复询问本地项目文档中已经确认的内容。
- 简单或已聚焦的请求及时路由到 `clarifying`，不为凑流程扩成项目规划。
- 写入 `project.md` 和 `backlog.md` 前必须重新读取实际对话和一手证据并据此填写，
  不凭记忆补写。
- 与用户交流时自然、直接，像与同事规划项目，不使用客服或公文式话术；文件正文按
  STEP 8 的文档语言关卡检查。

## References

- [Project planning guide](references/project-planning-guide.md)
- [Project template](assets/project-template.md)
- [Backlog template](assets/backlog-template.md)
