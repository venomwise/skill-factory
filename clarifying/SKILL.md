---
name: clarifying
description: >
  澄清一个可独立交付的软件需求，通过项目证据、逐步提问、方案比较和范围收敛生成
  `specs/<topic>/design.md` 与 Acceptance Criteria。当用户说“设计这个功能”
  “先把这个需求说清楚”“这个功能怎么实现”“方案还不明确”，或要处理单个功能、
  行为变更、复杂修复和技术取舍时调用。若请求仍是包含多个独立目标的大型项目，
  先使用 `brainstorming`。
---

# Clarifying Requirements Into Designs

把一个可独立交付的需求收敛成已确认的技术设计。项目级方向和目标拆分由
`brainstorming` 负责；实施任务拆分由 `spec-plan` 负责。

## When to use

- 创建或修改一个可独立交付的功能、组件或系统行为。
- 需求目标已聚焦，但边界、成功标准或技术方案尚不明确。
- 复杂修复尚未确定根因、目标行为或安全修复路径。
- 存在多种技术实现方向，需要比较成本、风险和取舍。
- `brainstorming` 已将项目目标标记为 `ready-for-clarifying`。

## When not to use

- 请求仍是包含多个独立目标的大型产品或项目；应使用 `brainstorming`。
- 根因、修复步骤和验证方式都已明确的局部修改。
- 仅涉及拼写、格式或明确的单步调整。
- 已有批准的 `design.md`，只需评审、拆任务或实施。

## Inputs / Outputs

**Inputs:**

- 用户提出的单个交付目标，或 `projects/<project>/backlog.md` 中的 work item。
- 现有项目上下文、约束，以及可选的 `projects/<project>/project.md`。

**Outputs:**

- 简单设计：经过确认的简短方案，不创建设计文档。
- 中等或复杂设计：`specs/<topic>/design.md`。
- 下一步：简单设计询问是否直接实施；中等或复杂设计推荐 `design-review`。
- 项目 work item：在设计中记录 Project Traceability，并返回 project handoff。
- 不创建 `tasks.md`，也不改写项目级 backlog。

## Workflow

### STEP 1: Explore project context

按优先级读取 README、CLAUDE.md、AGENTS.md、项目配置、入口点和相关模块。
输入来自项目规划时，完整读取 `project.md`，再读取独立 `backlog.md` 或其中唯一的
嵌入式 `## Backlog`，定位一个稳定 `WI-*`。这些内容只作为项目目标、边界和追踪来源。
只有理解近期变化确有必要时才查看最近 10 个提交。

STOP when 你能用 2-3 句话描述项目目的、技术栈、结构和本次交付目标所在边界。
项目复杂或不熟悉时，使用
[clarifying guide](references/clarifying-guide.md) 的探索优先级。

### STEP 2: Check requirement fit

根据真实项目证据检查需求是否与现有能力重复、冲突，是否会破坏既有行为，
以及是否存在明显更简单的路径。

IF 发现冲突或更简单替代方案，THEN PAUSE，向用户说明证据、影响和推荐方向，
等待用户确认后继续。

### STEP 3: Check delivery granularity

IF 输入来自项目规划，THEN 先执行 handoff 关卡：

1. `project.md` 的 Planning Status 是否为 `ready`？不满足时，停止并推荐使用
   `brainstorming` 继续项目规划。
2. 是否能唯一定位稳定 `WI-*`，且 Status 为 `ready-for-clarifying`？不满足时，停止并推荐
   使用 `brainstorming` 对账状态。
3. Readiness 是否没有阻止单项设计的项目决策或依赖？不满足时，报告具体阻塞并停止。
4. work item 的 Outcome、Scope、Non-Goals 和 Success signal 是否完整？不满足时，返回
   `brainstorming` 修复 backlog，不在本技能中补造项目目标。

handoff 关卡只验证项目准入，不替代下面的需求与技术澄清。

进入详细澄清前执行关卡：

1. 请求是否围绕一个用户或业务结果？不满足时，说明其中的独立目标并路由到
   `brainstorming`。
2. 该结果是否可以独立设计、验收和交付？不满足时，路由到 `brainstorming`。
3. 当前主要缺口是否属于需求边界或技术方案，而不是项目愿景、MVP 和目标优先级？
   不满足时，路由到 `brainstorming`。

路由前说明判断证据并等待用户确认；确认后 STOP，不在本技能中代替
`brainstorming` 完成项目拆分。多个模块围绕同一交付目标发生变化，不自动视为项目级需求。

### STEP 4: Clarify intent and boundaries

先诊断主要缺口：

- `Problem unclear`：重新界定实际要解决的问题。
- `Direction unclear`：探索几个有代表性的方向。
- `Boundaries unclear`：按 [clarifying guide](references/clarifying-guide.md) 的盲点清单排查领域盲点。
- `Solution unclear`：确认意图和约束后再比较技术方案。

每个需要用户输入的回合只问一个原子问题并等待回答。这是节奏控制，
不是总问题数上限。优先提供有意义的选项和你的倾向，让用户容易判断取舍。

不要因为请求中包含大量技术细节就假设意图已经确认。用户的优先级、偏好、
成功标准和明确不要的内容只能通过提问获得。

持续提问，直到以下条件全部满足：

1. 能说明用户想要什么、为什么需要、明确不要什么。不满足时，继续本步骤。
2. 关键约束和成功标准已知。不满足时，继续本步骤。
3. 适用的关键假设和盲点已确认或明确排除；盲点清单的适用项已逐项处置为设计覆盖、
   写入 Non-Goals 或确认不适用，不允许沉默跳过。不满足时，继续本步骤。

### STEP 5: Converge and propose

进入方案比较前执行关卡：

1. 本次会话是否至少问过一个澄清问题？不满足时，回到 STEP 4。
2. 意图、优先级、范围和成功标准是否明确？不满足时，回到 STEP 4。
3. 是否能区分项目事实与仍需用户决定的偏好？不满足时，先完成探索或提问。

关卡通过后，总结问题、范围、已验证事实、被质疑的假设和关键盲点。
提出 1-3 个可行方案，说明成本、风险和取舍，并优先给出推荐方案。
只有一个方案可行时，说明其他方向被排除的原因。

让用户确认问题总结和方案选择。用户修正意图或方案时，更新总结并再次确认；
不要在同一回合提出确认问题后用假设代替用户回答。

### STEP 6: Classify the design

只有以下条件全部为真时，才归类为简单设计：

- 变更局限于单个组件或行为。
- 遵循项目既有模式。
- 不涉及数据迁移、外部集成、安全、权限、并发或可靠性问题。
- 不改变公共契约。
- 测试和回滚直接。
- 没有未解决问题。

任一条件为 false 或不确定时，归类为中等或复杂。

来自项目 work item 的请求无论复杂度如何都进入设计文档路径，因为 backlog 需要稳定的
`design.md`、Project Traceability 和后续状态证据。简单设计的无文档路径只适用于独立请求。

### STEP 7: Route and write

**简单设计：**

说明低风险依据，询问用户是否直接实施。只有用户同意后才能进入实施。
本路径只适用于不关联项目 work item 的独立请求。
反馈改变意图或方案时回到 STEP 4 或 STEP 5。

**中等或复杂设计：**

说明风险和核心逻辑，询问是否写入设计文档。用户同意后：

1. 按 [design document template](assets/design-doc-template.md) 写入设计。
2. 写入 `specs/<topic>/design.md`；topic 使用 kebab-case。
3. 同名目录或路径有歧义时，先确认重命名还是沿用。
4. 输入来自项目 work item 时，填写 Project Traceability 的 Project、Backlog 和 Work item；
   独立请求删除该可选章节。
5. 重新读取实际讨论、项目文档和项目证据并据此填写；不要凭记忆补写或编造选项。
6. 设计中引用的既有列、字段、错误码、配置和既有行为必须来自本次会话的实际检索结果；
   无法核实的事实写入 Open Questions，不作为设计前提。
7. 枚举 Data Flow 中的每个跨边界调用，为每个调用分配 Contract ID，并写清请求、响应和错误定义；
   缺契约的调用先补齐再写文件。
8. 只把持久且非显然的选择写入稳定 `DR-<semantic-topic>`。
9. 用户明确驳回且可能被重复提出的边界，写入带 `Rejected concern` 和
   `Revisit when` 的 Decision。
10. 事实修正、格式处理和章节传播不写入 Decision Record。
11. 数据或契约发生变化时，按模板填写 Data Model / Interfaces Change Summary。
12. 在 Acceptance Criteria 中写入用户确认的可观察行为。
13. 重新读取实际写入的完整 `design.md`，执行下面的文档语言关卡。
14. 语言修改后重新检查 Goals、Decision、Data Model、Interfaces、Data Flow、Error Handling、
    Acceptance Criteria 和 Testing 的对应关系。
15. 运行结构校验器，修复全部违规项。

**文档语言关卡：**

1. 中文正文是否自然、直接？英文只保留技术术语、模板字段、章节名、路径、ID 和固定枚举。
   IF 否，THEN 只改对应句子。
2. Summary、Context、Proposed Solution、Components 和 Data Flow 是否写清参与者、触发条件、
   输入、动作、结果或副作用？IF 否，THEN 改成具体关系，不用抽象名词代替设计内容。
3. 是否存在逐词翻译、商业黑话、流程隐喻、模板化开场或空泛总结？IF 是，THEN 删除姿态层，
   保留技术含义；生命周期、消息消费等正常技术术语不做机械替换。
4. 数字、事实、责任主体、用户决定、Decision、状态、字段、Contract ID、路径、命令、引用和
   技术术语是否保持不变？IF 否，THEN 恢复原意并缩小改写范围。
5. 语言修改是否新增、删除、弱化或扩大了 Goals、Non-Goals、风险、契约、错误行为或 AC？
   IF 是，THEN 撤销该修改，并依据已确认设计重新表述。

结构校验命令：

```bash
node <clarifying-skill>/scripts/check-design-doc.mjs specs/<topic>/design.md
```

**Shadow review：**

中等或复杂设计通过结构校验后，在推荐 `design-review` 前，启动一个隔离的只读子代理对
`design.md` 执行影子评审，把高信号问题在交付前消灭。启动契约、评审范围和候选项格式见
[shadow review contract](references/shadow-review.md)。

1. 运行环境无法启动只读子代理时，跳过本环节并在最终响应说明；不由主流程自我评审替代。
2. 候选项没有设计位置或项目证据时直接丢弃，不进入修复。
3. 候选项改变意图、范围或需要新的用户决策时，向用户呈现并等待确认；确认后回到 STEP 4 或 STEP 5。
4. 候选项属于一致性缺失、传播缺口或格式问题时，作为直接修复应用，并在最终响应中汇总。
5. 修复后重新执行文档语言关卡和结构校验器。

影子评审只做交付前预检，不替代 `design-review` 的独立评审，不写入 `review.md`，也不生成
finding ID。简单设计没有设计文档，跳过本环节。

设计写入文件后推荐运行 `design-review`。用户明确要求进入 `spec-plan` 或实施时，
仍以用户指令为准，但不得绕过未解决的设计问题。

输入来自项目 work item 时，在最终响应追加以下 handoff，不直接改 backlog：

```yaml
project_handoff:
  work_item: "WI-<semantic-name>"
  design: "specs/<topic>/design.md"
  proposed_status: "design-ready"
  status_evidence: "design.md Project Traceability 与 work item 匹配"
  next_owner: "brainstorming"
```

提醒用户由 `brainstorming` 核实 handoff 和可选 review 后更新 `Status` 与 `Spec`。

反馈使意图或边界重新模糊时回到 STEP 4；只需重新选择方案时回到 STEP 5。

## Verification

仅针对已写入的 `design.md`：

- [ ] 文件位于 `specs/<topic>/design.md`，topic 为 kebab-case。
- [ ] 项目 work item 对应的设计包含可解析且匹配的 Project Traceability。
- [ ] 标题和章节符合 canonical design document template。
- [ ] Goals、Non-Goals、Architecture、Components、Data Flow、Error Handling、
  Acceptance Criteria、Testing 和 Open Questions 有实质内容。
- [ ] Decision 只包含符合准入条件的持久选择，并使用稳定语义 ID。
- [ ] 用户明确驳回的持久边界包含 `Rejected concern` 和 `Revisit when`。
- [ ] 涉及数据变化时，Data Model 使用 Change Summary 和适用的固定子章节。
- [ ] 可执行 DDL/DML 位于独立 SQL 文件，设计只引用路径和迁移语义。
- [ ] 涉及契约变化时，Interfaces 枚举全部 ADD/MODIFY/REMOVE 契约。
- [ ] Data Flow 的跨边界调用均能映射到 Contract ID。
- [ ] 每个核心 Goal 至少有一条 AC，ID 唯一且符合 kebab-case 约定。
- [ ] 每条 AC 只描述一个可观察结果，覆盖正常流、错误流和关键边界。
- [ ] AC 不包含实现步骤、验证命令、被拒方案或未经确认的内容。
- [ ] Open Questions 不包含仍会改变行为、范围或 Acceptance Criteria 的未决项。
- [ ] 项目 work item 的最终响应包含完整 project handoff，backlog 未被本技能修改。
- [ ] 盲点清单的适用项已逐项处置为设计覆盖、Non-Goals 或不适用。
- [ ] 设计中引用的项目事实均有本次会话的检索证据；无法核实项在 Open Questions。
- [ ] 每个跨边界调用有 Contract ID 与请求、响应、错误定义。
- [ ] 中等或复杂设计已执行影子评审，或已说明跳过原因。
- [ ] 已重新读取完整 `design.md` 并通过文档语言关卡，抽象说法能够对应具体设计关系。
- [ ] 语言修改没有改变事实、Decision、字段、Contract ID、错误语义或 AC 行为。
- [ ] 结构校验器通过：章节、AC、DR、Contract ID 和行宽均无违规。

## Safety & guardrails

- 未经确认不得写设计或开始实施。
- 不把项目 backlog 直接复制成设计，也不为 `spec-plan` 预先生成 tasks。
- 不修改 backlog 的 Status、Spec、Status evidence 或 Current Focus；这些字段由
  `brainstorming` 统一写入。
- 简单明确的请求允许在少数几轮内完成，不为凑流程追加问题。
- 发散阶段允许探索；收敛和写设计时删除不服务当前 Goals 的内容。
- 当前行为和设计以 `design.md` 为准；该文件不兼任完整讨论日志。
- 写入 `design.md` 前必须重新读取实际对话和项目证据并据此填写，不凭记忆补写。
- 与用户交流时自然、直接，像与同事讨论，不使用客服或公文式话术；文件正文按 STEP 7
  的文档语言关卡检查。

## References

- [Clarifying guide](references/clarifying-guide.md)
- [Canonical design document template](assets/design-doc-template.md)
- [Shadow review contract](references/shadow-review.md)
- [Design document validator](scripts/check-design-doc.mjs)
