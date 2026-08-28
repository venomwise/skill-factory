---
name: design-review
description: >
  对 design.md 进行独立、证据化的设计评审，并在同目录写入 review.md。
  当用户要求评审设计、检查完整性、可实施性、项目契合度、盲点或过度设计时调用。
  触发词包括“评审设计”“review design”“design review”“检查 design.md”
  “设计盲点”。不修改设计；根据评审结果修订设计应使用 design-refine。
---

# Design Review

对 `design.md` 执行独立质量关卡。首轮执行 Full Review；完成 refine 后默认执行
Closure Review。唯一写入产物是同目录的 `review.md`。

## When to use

- 规划或实施前评审 `design.md`。
- 检查设计完整性、可实施性、项目契合度、盲点和过度设计。
- refine 后复核 finding 是否真正闭合。

## When not to use

- 代码或 pull request 评审。
- 从零构建单项需求设计，应使用 `clarifying`。
- 根据 finding 修改设计，应使用 `design-refine`。
- 没有设计文档的琐碎修改。

## Inputs / Outputs

**Inputs:**

- `design.md` 路径或可唯一定位的 `specs/<topic>/`。
- 同目录可选的历史 `review.md`。
- 设计所属项目的代码、配置、数据库和约束。

**Outputs:**

- 同目录 `review.md`，使用 [review template](assets/review-template.md)。
- 评审时 Verdict、当前 readiness、证据化 findings 和推荐下一步。
- 不修改 `design.md` 或项目文件。

## Authoritative contracts

开始评审前必须读取：

- [Canonical design document template](../clarifying/assets/design-doc-template.md)
- [Review lifecycle contract](references/review-lifecycle.md)
- [Review rubric](references/review-rubric.md)

章节、Decision、Data Model、Interfaces、AC 和格式要求以前者为准。
Full/Closure、finding ID、Origin、状态和 readiness 以后两者为准。

## Workflow

### STEP 1: Locate and read inputs

用户给出路径时直接使用。否则在 `specs/` 中定位 `design.md`：

- 恰好一个匹配：直接使用。
- 没有或存在多个匹配：只询问目标路径，然后 STOP。

完整读取 `design.md`。同目录存在 `review.md` 时完整读取其中的所有 finding ID、
Closure、Current Readiness、changed sections、Finding Closure Proof 和已接受或延期风险。

IF `design.md` 包含 Project Traceability，THEN 解析 Project、Backlog 和 Work item，
完整读取链接的项目规划并核对稳定 `WI-*`。路径缺失、ID 不存在或目标边界不一致时，
记录候选 finding；不要向用户询问能够从这些文件核实的事实。

IF 历史 review 缺少 `## Current Readiness`，或只使用 severity-coded finding ID，THEN
将其标记为 legacy review：保留一行迁移说明，但不猜测旧状态；本轮执行 Full Review，
按当前协议重新分配 `F-###` 并生成双视图报告。

### STEP 2: Select review mode

根据 review lifecycle 选择模式：

- 没有历史 review：Full Review。
- 历史 review 属于 legacy protocol：Full Review migration。
- 用户明确要求完整评审：Full Review。
- 自动复审时，Current Readiness 为 `not-started`、`in-progress` 或 `blocked`：
  停止并建议继续 `design-refine`，不得评审处理中设计。
- 自动复审时，Current Readiness 已为 `ready / Go`：说明无需再次评审。
- 自动复审时，Current Readiness 为 `ready-for-closure` 但 pre-closure audit 缺失或失败：
  停止并返回 `design-refine`，不得开始独立 review。
- audit 通过后，Goals、Non-Goals、核心架构、数据所有权、公共契约、安全、权限、外部集成、
  迁移边界或相关项目事实变化：升级为 Full Review。
- audit 通过且未触发 Full 条件：Closure Review。

`Next review mode` 是 refine 的预判，不是评审结论。独立核对实际 changed sections 和 lifecycle 触发条件；
不一致时按实际条件选择模式，并在 Mode reason 记录差异。不要把每次复审默认当作新的全量审计。

### STEP 3: Explore project evidence

Full Review 按以下优先级探索：

1. README、CLAUDE.md、AGENTS.md。
2. 项目配置和依赖。
3. 与设计相关的入口、组件、接口和持久化代码。
4. 需要理解近期变化时，查看最近 10 个提交。
5. 相关现有 specs 和机器可读契约。
6. Project Traceability 指向的项目规划和 work item。

Closure Review 复用上一轮证据，只重新核实 changed sections 影响的项目区域，
以及 `context-change` 所指向的新事实。

refine 后的 review 开始前读取 pre-closure audit、影响矩阵和 Finding Closure Proof。
它们用于定位核查范围，但不构成评审结论；reviewer 必须从原 Issue、Evidence 和期望结果
重新执行 Closure test 和反例检查，不强制采用原 Recommendation 的具体修法。

设计依赖真实数据库时使用 `db-explorer`。能通过代码、配置、数据库或历史核实的声明，
必须主动核实，不把可检索事实交给用户回答。

STOP when 每个候选 finding 都能指向设计位置、项目证据或明确的契约缺口。

### STEP 4: Run deterministic preflight

全量检查以下规则：

- canonical template 的章节、稳定 Decision ID 和 AC ID。
- 可选 Project Traceability 的路径、`WI-*` 和设计范围映射。
- Data Model / Interfaces Change Summary 与详细子章节的映射。
- Data Flow 中 Contract ID 的引用完整性。
- Open Questions 是否仍含行为相关未决项。
- Markdown 行宽。

行宽校验命令：

```bash
node <clarifying-skill>/scripts/check-markdown-lines.mjs <path/to/design.md>
```

把同类格式问题合并为一个 D3 finding，并列出代表性行号；不要为每一行创建 finding。

### STEP 5: Perform semantic review

**Full Review:**

1. refine 后升级为 Full 时，先独立复测 Finding Closure Proof；原反例仍成立时沿用原 finding ID。
2. 按 rubric 的 7 个维度建立候选 finding ledger。
3. 建立 Goal -> Solution -> AC 覆盖关系。
4. 单独检查 Data Model、Interfaces、迁移、并发和错误契约。
5. 检查共享模型、身份、容量、错误或事务约束在成对契约间是否一致；有意差异是否有 Decision 和 AC。
6. 检查项目契合度、适用盲点和 YAGNI。
7. 按根因合并跨维度候选项，校准严重性。
8. 对最终集合再做一次交叉检查，然后冻结本轮 findings。

**Closure Review:**

1. 从原 Issue、Evidence 和期望结果独立执行每条 Closure test，不以当前 status 作为闭合证据。
2. 复测 Finding Closure Proof 的原反例；仍成立时沿用原 finding ID 并设为 `reopened`。
3. 核对影响矩阵、changed sections 和其中预先列出的关联章节。
4. 检查成对契约的共享约束，以及有意差异对应的 Decision、Error Handling、AC 和 Testing。
5. 验证相关 Decision 是否仍有效，`Revisit when` 是否触发。
6. 只在满足 Origin gate 时创建新 finding。
7. 无法归入允许 Origin 的自由优化建议不得进入 Findings。

同一根因沿用原 `F-###`。新 ID 从历史最大编号递增，不复用已关闭编号。
Closure 新 finding 必须写明 Origin 和因果证据。

### STEP 6: Compute snapshot and readiness

根据 rubric 计算评审时 Verdict：

- Blocker 存在：Reject。
- 无 Blocker、存在 Major：Revise。
- 无 Blocker/Major：Pass。

根据 lifecycle 初始化 `Current Readiness`：

- 无 finding：`ready / Go`。
- 仅 Minor：`not-started / Conditional`，等待用户修复或延期选择。
- 存在 Major：`not-started / Conditional`。
- 存在 Blocker：`not-started / No-Go`。

Full Review 的新 finding 初始化为 `pending`。Closure Review 中继续成立的历史 finding
保留原 ID 并初始化为 `reopened`；满足 Origin gate 的新 finding 初始化为 `pending`。
Closure 已确认关闭的历史 finding 只进入紧凑 Closure。

### STEP 7: Write and validate review.md

使用 review template 写入同目录 `review.md`：

- Review Snapshot 使用设计当前语言，结构标题和固定枚举保持英文。
- 每条 finding 包含 Severity、Introduced in、Origin、Location、Issue、Evidence 和 Recommendation。
- Current Readiness 包含状态、Resolution ref、audit、下一轮模式、影响矩阵、Finding Closure Proof 和当前下一步。
- Closure 不复制上一轮完整 finding 文本。
- Accepted / Deferred Risks 原样携带并根据当前设计重新核实。

写完后重新读取实际 `review.md`，执行文档语言关卡：

1. 中文正文是否自然、直接？英文只保留模板字段、技术术语、路径、ID、Origin 和固定枚举。
   IF 否，THEN 只改对应句子。
2. 每条 finding 的 Issue 是否说明具体错误、缺失或风险条件，Evidence 是否给出可核实依据，
   Recommendation 是否指出需要修改的内容或必须做出的决策？IF 否，THEN 回到 STEP 3 或
   STEP 5 核实后重写，不编造证据；不用“存在一定风险”“不够完善”“进一步优化”代替证据
   和动作。
3. Summary、Note 和风险说明是否存在逐词翻译、商业黑话、流程隐喻、模板化开场或空泛总结？
   IF 是，THEN 删除姿态层，保留评审结论和证据关系。
4. Review mode、Verdict、Severity、finding ID、Origin、状态、readiness、证据和责任主体是否
   保持不变？IF 否，THEN 恢复原意并缩小改写范围；不得为了语气温和而降低严重性或确定性。
5. 语言修改是否增加了没有证据的 finding，或改变了已有 finding 的范围和 Closure 结果？
   IF 是，THEN 撤销该修改，并依据实际证据重新表述。

语言检查完成后重新核对模板结构和 Current Readiness，再运行行宽校验器。review 自身存在
非豁免超长行时，先换行再结束。

### STEP 8: Recommend next step

- Reject/Revise：运行 `design-refine`，完成后按 `Next review mode` 运行 Full 或 Closure Review。
- Pass 且有 Minor：让用户选择修复或延期；全部终态后进入 `spec-plan`。
- 无 finding：进入 `spec-plan`。

只推荐，不自动调用下游 skill。

## Review dimensions

| ID | Dimension | Focus |
|----|-----------|-------|
| D1 | Completeness | Goals、对应设计、AC 和必要章节是否完整且相互对应 |
| D2 | Usability | `spec-plan` 是否无需发明行为或契约 |
| D3 | Conformance | 是否符合 canonical template 和格式契约 |
| D4 | Project Fit | 是否符合真实技术栈、模块边界和已有模式 |
| D5 | Blind Spots | 适用的状态、迁移、权限、失败和运营风险 |
| D6 | Over-Engineering | 是否存在无法追溯到 Goals 的复杂度 |
| D7 | Optimization | 是否有实质性更简单或复用更多的方案 |

详细检查项和严重性规则见 [review rubric](references/review-rubric.md)。

## Verification

- [ ] 已读取 canonical design template、review lifecycle 和 rubric。
- [ ] 已正确选择 Full 或 Closure，并记录升级原因。
- [ ] Legacy review 已通过 Full Review migration 转换，未猜测旧状态。
- [ ] refine 后自动复审已有通过的 pre-closure audit；用户明确要求 Full 的例外已记录。
- [ ] Closure 只在 Current Readiness 为 ready-for-closure 时执行。
- [ ] Closure 开始前已存在通过的 pre-closure audit；该审计未被当作独立评审结论。
- [ ] Next review mode 已按实际 changed sections 独立核对，模式差异已记录。
- [ ] Finding Closure Proof 已从原 finding 独立复测，未使用当前 status 代替证据。
- [ ] 成对契约的共享约束已检查，有意差异存在 Decision 和 AC。
- [ ] Full Review 已在冻结前完成全部维度和根因合并。
- [ ] Closure Review 只检查旧 finding、changed sections 和 Impact matrix 中预先列出的关联章节。
- [ ] Closure 新 finding 有允许的 Origin 和因果证据。
- [ ] 后续 Full Review 的新 finding 使用 context-change 或 baseline-miss。
- [ ] finding 使用稳定 `F-###`，严重性不编码进 ID。
- [ ] 同一根因延续时沿用原 ID。
- [ ] 每条 finding 有位置、问题、证据和具体建议。
- [ ] Data Model 和 Interfaces 按 canonical schema 检查。
- [ ] 存在 Project Traceability 时，Project、Backlog 和 Work item 均可解析且互相一致。
- [ ] Review Snapshot 与 Current Readiness 职责分离。
- [ ] Closure 紧凑，不复制完整历史 finding。
- [ ] `design.md` 和项目文件未修改。
- [ ] 已重新读取实际 `review.md` 并通过文档语言关卡，每条 finding 都是具体判断而非空泛评价。
- [ ] 语言修改没有改变 Verdict、Severity、finding ID、Origin、readiness 或证据关系。
- [ ] `review.md` 行宽校验通过。

## Safety & guardrails

- 证据优先于观点；指不出位置和证据就不是 finding。
- Blocker 只用于真正阻止实现或与项目事实冲突的问题。
- Blind Spot 表述为需要检查的事项，不假定作者一定遗漏。
- Decision 仍有效且 `Revisit when` 未触发时，不重复提出相同问题。
- 小设计得到短报告，不用仪式填充。
- 本 skill 不修设计、不规划、不实施，也不调用下游 skill。
- 与用户交流时直接、自然，不使用模板化客服话术；文件正文按 STEP 7 的文档语言关卡检查。

## References

- [Canonical design document template](../clarifying/assets/design-doc-template.md)
- [Review lifecycle contract](references/review-lifecycle.md)
- [Review rubric](references/review-rubric.md)
- [Review template](assets/review-template.md)
