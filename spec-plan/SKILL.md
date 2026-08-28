---
name: spec-plan
description: >
  从已批准的 specs/<topic>/design.md 生成可追溯的 tasks.md。
  当用户提到做 spec、拆任务、出执行计划、write a spec、plan tasks 时，
  必须使用此技能。前置条件：design.md 包含完整的 Acceptance Criteria；
  若没有设计，先路由到 clarifying。
---

# Spec Plan Skill

## When to use

- 从已批准的 `design.md` 创建或刷新 `tasks.md`
- 将设计方案和 Acceptance Criteria 转化为可执行、可验证的实施任务
- 建立任务到 AC 或已批准设计章节的可追溯性

## When not to use

- 只需要快速 TODO 列表或单个文件编辑
- 工作已经由现有 `tasks.md` 捕获，只需要实现
- 没有已批准的 `design.md`；应先运行 `clarifying`

## Inputs / Outputs

**Inputs:**

- 已批准的 `specs/<topic>/design.md`
- 同目录下可选的 `review.md`

**Outputs:**

- `specs/<topic>/tasks.md`
- 最终响应推荐用户对生成的 `tasks.md` 运行 `spec-exec`；不自动开始实现

## Acceptance Criteria contract

已批准的行为和设计方案以 `design.md` 为准。`spec-plan` 只读取并使用其中的
`## Acceptance Criteria`，不得创建、补全、改写或重新解释 AC。

有效 AC 必须满足：

- ID 全局唯一，格式为 `AC-<domain>-<behavior>`，各段使用小写 kebab-case
- 每条只描述一个可观察、可测试的结果
- 使用 `WHEN` 或 `IF` 表达条件，并使用 `THEN` 和 `SHALL` 表达结果
- 每个核心 Goal 至少映射到一条 AC
- 正常流、错误流和适用的关键边界均被覆盖
- 不与 Goals、Non-Goals、Proposed Solution、Error Handling 或 Testing 冲突
- `Open Questions` 中没有仍会改变行为、范围或验收结果的未决项

## Task traceability

每个可执行 sub-task 和独立 Checkpoint 必须至少使用一种追踪元数据；包含子任务的 Phase 行只是容器，不需要追踪元数据：

- `_Acceptance: AC-..._`：直接实现或验证可观察行为。一个任务可引用多个 AC。
- `_Design: <section path>_`：脚手架、架构调整或内部重构。路径必须指向 `design.md` 中实际存在的章节。
- 同时使用两者：任务既受行为契约约束，也依赖具体设计方案。

不要把纯内部任务强行绑定到 AC。每条 AC 必须至少被一个实现任务、测试任务或 Checkpoint 的 `_Acceptance:` 覆盖。

## Workflow

1. **定位设计。**
   - IF 用户提供了 `specs/<topic>/design.md`，THEN 使用它。
   - IF 用户未提供路径，THEN 精确回复：`请指定 design 文件路径（例如 specs/<topic>/design.md）。`然后结束流程。不要搜索、列举或推断设计。
   - IF 路径不存在，THEN 结束流程并推荐先运行 `clarifying`。
2. **读取设计并确定语言。** 完整读取 `design.md`；从说明性正文推断生成内容的主导语言，忽略标题、代码、路径和固定标签。不要仅为确定语言而询问。
3. **执行严格 AC 关卡。** 按 *Acceptance Criteria contract* 检查章节、格式、唯一性、可测试性、Goal 覆盖、错误与边界覆盖、Open Questions 以及跨章节一致性。
   - IF 任一检查失败，THEN 停止，不生成或刷新 `tasks.md`；指出具体位置和缺口，
     并建议返回 `clarifying` 或 `design-refine`。
   - IF 行为含糊、缺失或冲突，THEN 不得通过推断补齐。
4. **处理可选评审。** 检查 `design.md` 同目录是否存在 `review.md`。
   - IF 不存在，THEN 继续。
   - IF 存在，THEN 读取 `## Current Readiness`，不得使用评审时 Verdict 代替当前状态。
   - IF `Overall` 为 `ready` 且 `spec-plan readiness` 为 `Go`，THEN 执行以下关卡：
     1. Current Readiness 中的每个 finding 都是终态：`resolved`、`rejected`、
        `accepted-risk` 或 `deferred`。
     2. 每个 `Resolution ref` 都能解析到 `design.md` 中实际存在的 `DR-*`、设计章节，
        或 `review.md` 的 Accepted / Deferred Risks 条目。
     3. `rejected` 必须指向包含 `Rejected concern` 和 `Revisit when` 的持久边界 Decision。
     4. changed sections 与实际更新内容一致；行为或错误语义变化已同步到 AC 和 Testing。
     全部通过后继续；任一失败则停止并报告 finding ID 和断裂引用。
   - IF Current Readiness 不是 `ready / Go`，THEN 停止，报告 `Overall`、readiness 和
     `Next step`，不得根据 Review Snapshot 的旧 Verdict 绕过当前关卡。
   - IF review 存在但缺少 `## Current Readiness`，THEN 停止并建议重新运行
     `design-review` 迁移到当前 review lifecycle；不要解析旧的 severity-coded Source 协议。
5. **建立内部覆盖表。** 将每条 AC 映射到实现、测试或 Checkpoint，
   并将必要的内部设计工作映射到具体设计章节。覆盖表仅用于规划和校验，
   不输出独立需求文档。
6. **生成 tasks.md。** 使用 `assets/tasks.template.md`。保留英文结构标题、
   checkbox/task 编号、`Phase`、`Checkpoint`、可选 `*` 语法、
   `_Acceptance:` 和 `_Design:` 标签；生成内容使用步骤 2 确定的语言。
   - 任务写明具体文件路径和函数、类或交付物。
   - 每个功能阶段包含适用的测试任务。
   - 关键里程碑添加 Checkpoint，列出验证命令、引用的 AC 和明确通过条件。Checkpoint 是执行 Agent 的验证任务，不是用户批准门。
   - 非必要测试、验证、文档收尾和 nice-to-have 功能使用 `- [ ]*`；不要使用“可选”或 `Optional:` 文本标签代替。
7. **验证覆盖。** 检查所有引用都能解析、每个可执行 sub-task 和独立 Checkpoint 至少有一种追踪元数据、每条 AC 至少被覆盖一次，并确认任务没有来自被拒方案、未来想法或 Non-Goals。

## Verification

- [ ] 起草前已完整读取 `design.md`
- [ ] `Acceptance Criteria` 通过严格关卡，无缺失、重复、歧义或冲突
- [ ] 可选 `review.md` 不存在，或其 Current Readiness 为 `ready / Go`，
  所有 finding 均为终态且 Resolution ref 可解析
- [ ] 只生成或刷新 `tasks.md`，未生成独立需求文档
- [ ] 每个可执行 sub-task 和独立 Checkpoint 至少包含 `_Acceptance:` 或 `_Design:`，且引用可解析
- [ ] 每条 AC 至少被一个任务或 Checkpoint 覆盖
- [ ] Checkpoint 包含验证命令、AC 引用和通过条件
- [ ] Phase 使用 `- [ ] N. Phase N:`，任务描述使用缩进 bullets
- [ ] 所有非必要步骤使用 `- [ ]*`，可选 Phase 的子任务继承可选性
- [ ] 任务未引入被拒方案、未来想法、Non-Goals 或设计外行为

## Safety & guardrails

- 没有已批准且 AC 完整的 `design.md`，不生成 `tasks.md`。
- 不自动选择设计路径。
- 不发明、补写或修改 AC；设计有缺口时严格阻断。
- 不生成 `requirements.md`，也不使用旧 `_Requirements:` 协议。
- 工作未实际完成前，不标记任务为完成。

## References

- [Tasks template](assets/tasks.template.md)
