---
name: brainstorming
description: >
  识别开发需求并通过逐步澄清、方案比较和范围收敛生成设计。
  当用户要开发功能、改变系统行为、修复尚无明确根因或方案的复杂问题，
  或处理存在多种实现方向的优化时调用。触发词包括“帮我设计”“先梳理需求”
  “这个功能怎么实现”“方案还不明确”。
---

# Brainstorming Ideas Into Designs

## When to use

- 创建新功能或组件。
- 修改现有系统行为。
- 需求、边界或技术方案尚不明确。
- 存在多种实现方式，需要比较取舍。

## When not to use

- 根因和修复步骤都已明确的 bug。
- 仅涉及拼写、格式或明确的单步修改。
- 用户明确拒绝设计流程。

## Inputs / Outputs

**Inputs:**

- 用户的想法、问题或目标。
- 现有项目上下文和约束。

**Outputs:**

- 简单设计：经过确认的简短方案，不创建设计文档。
- 中等或复杂设计：`specs/<topic>/design.md`。
- 下一步：简单设计询问是否直接实施；中等或复杂设计推荐 `design-review`。

## Workflow

### STEP 1: Explore project context

按优先级读取 README、CLAUDE.md、AGENTS.md、项目配置、入口点和相关模块。
只有理解近期变化确有必要时才查看最近 10 个提交。

STOP when 你能用 2-3 句话描述项目目的、技术栈、结构和本次变更所在边界。
项目复杂或不熟悉时，使用
[brainstorming guide](references/brainstorming-guide.md) 的探索优先级。

### STEP 2: Check requirement fit

根据真实项目证据检查需求是否与现有能力重复、冲突，是否会破坏既有行为，
以及是否存在明显更简单的路径。

IF 发现冲突或更简单替代方案，THEN PAUSE，向用户说明证据、影响和推荐方向，
等待用户确认后继续。

### STEP 3: Check scope

IF 需求跨越多个独立子系统或可以分别交付的能力，THEN 提出拆解方案，说明子项目、
依赖和推荐顺序，并等待用户确认。

多个模块上的同一个小改动不自动视为大需求。

### STEP 4: Clarify intent and boundaries

先诊断主要缺口：

- `Problem unclear`：重新界定实际要解决的问题。
- `Direction unclear`：探索几个有代表性的方向。
- `Boundaries unclear`：检查与领域相关的盲点。
- `Solution unclear`：确认意图和约束后再比较技术方案。

每个需要用户输入的回合只问一个原子问题并等待回答。这是节奏控制，
不是总问题数上限。优先提供有意义的选项和你的倾向，让用户容易判断取舍。

不要因为请求中包含大量技术细节就假设意图已经确认。用户的优先级、偏好、
成功标准和明确不要的内容只能通过提问获得。

持续提问，直到以下条件全部满足：

1. 能说明用户想要什么、为什么需要、明确不要什么。If no，继续本步骤。
2. 关键约束和成功标准已知。If no，继续本步骤。
3. 适用的关键假设和盲点已确认或明确排除。If no，继续本步骤。

### STEP 5: Converge and propose

进入方案比较前执行关卡：

1. 本次会话是否至少问过一个澄清问题？If no，回到 STEP 4。
2. 意图、优先级、范围和成功标准是否明确？If no，回到 STEP 4。
3. 是否能区分项目事实与仍需用户决定的偏好？If no，先完成探索或提问。

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

### STEP 7: Route and write

**简单设计：**

说明低风险依据，询问用户是否直接实施。只有用户同意后才能进入实施。
反馈改变意图或方案时回到 STEP 4 或 STEP 5。

**中等或复杂设计：**

说明风险和核心逻辑，询问是否写入设计文档。用户同意后：

1. 使用唯一权威
   [design document template](assets/design-doc-template.md)。
2. 写入 `specs/<topic>/design.md`；topic 使用 kebab-case。
3. 同名目录或路径有歧义时，先确认重命名还是沿用。
4. 只从实际讨论和项目证据恢复内容；不要凭印象重建或编造选项。
5. 只把持久且非显然的选择写入稳定 `DR-<semantic-topic>`。
6. 用户明确驳回且可能被重复提出的边界，写入带 `Rejected concern` 和
   `Revisit when` 的 Decision。
7. 事实修正、格式处理和章节传播不写入 Decision Record。
8. 数据或契约发生变化时，按模板填写 Data Model / Interfaces Change Summary。
9. 在 Acceptance Criteria 中写入用户确认的可观察行为。
10. 运行只读格式校验器，修复所有非豁免超长行。

格式校验命令：

```bash
node <brainstorming-skill>/scripts/check-markdown-lines.mjs specs/<topic>/design.md
```

设计写盘后推荐运行 `design-review`。用户明确要求进入 `spec-plan` 或实施时，
仍以用户指令为准，但不得绕过未解决的设计问题。

反馈使意图或边界重新模糊时回到 STEP 4；只需重新选择方案时回到 STEP 5。

## Verification

仅针对落地的 `design.md`：

- [ ] 文件位于 `specs/<topic>/design.md`，topic 为 kebab-case。
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
- [ ] Open Questions 不包含仍会改变行为、范围或验收结果的未决项。
- [ ] 行宽校验通过，无非豁免的 120 字符超长行。

## Safety & guardrails

- 未经确认不得写设计或开始实施。
- 简单明确的请求允许在少数几轮内完成，不为凑流程追加问题。
- 发散阶段允许探索；收敛和写设计时删除不服务当前 Goals 的内容。
- `design.md` 是当前有效行为和设计的唯一权威，不兼任完整讨论日志。
- 长寿命内容必须从实际对话和项目证据恢复，不凭印象补写。
- 输出语气自然、直接，像与同事讨论，不使用客服或公文式话术。

## References

- [Brainstorming guide](references/brainstorming-guide.md)
- [Canonical design document template](assets/design-doc-template.md)
- [Markdown line-width validator](scripts/check-markdown-lines.mjs)
