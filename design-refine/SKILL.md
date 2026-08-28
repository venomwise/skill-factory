---
name: design-refine
description: >
  根据 design-review 的 review.md 与用户逐项确认设计决策、边界驳回和风险，
  批量处理无需新决策的直接修复，并同步更新 design.md 和 Current Readiness。
  当用户说“根据 review 改设计”“处理评审意见”“优化 design.md”
  “review findings 怎么落地”时必须调用。
---

# Design Refine

把 review finding 落实到 `design.md`，并更新可追溯的 readiness。真正的设计选择逐项确认；
能够从已批准设计和项目事实唯一推导的修复批量确认。

## Critical rules

1. 每个需要用户输入的回合只问一个原子问题，然后 STOP。
2. finding 不等于 Decision；先分类，再决定写入位置。
3. Decision 使用稳定语义 `DR-*` 并原位更新，不追加 `Revised`。
4. 用户明确驳回的重要边界写入带 `Rejected concern` 和 `Revisit when` 的 Decision。
5. 直接修复不进入 Decision Record。
6. 只更新 `review.md` 的 Current Readiness、风险和当前下一步，不改写 Review Snapshot。
7. finding 终态必须有 design 章节、Decision 或风险记录作为证据。

## When to use

- 已有 `design.md` 和对应 `review.md`，需要处理 findings。
- 需要比较设计方案或确认持久边界。
- 需要补全 review 指出的设计传播和一致性缺口。

## When not to use

- 只评审设计，应使用 `design-review`。
- 从零构建单项需求设计，应使用 `clarifying`。
- 已批准设计需要拆任务，应使用 `spec-plan`。
- 只是修改不涉及 review 的小段文字。

## Inputs / Outputs

**Inputs:**

- `specs/<topic>/design.md`。
- 同目录 `review.md`。
- 用户补充的业务选择、约束和风险偏好。

**Outputs:**

- 更新后的 `design.md`。
- 更新后的 `review.md` Current Readiness、风险和下一步。
- Review Snapshot 的 Verdict、Findings 和 Dimension Summary 保持不变。

## Authoritative contracts

开始前必须读取：

- [Canonical design document template](../clarifying/assets/design-doc-template.md)
- [Review lifecycle contract](../design-review/references/review-lifecycle.md)

不要在本 skill 中重新定义章节、finding ID、Origin、状态或 readiness。

## Interaction protocol

每轮只询问一个独立决策变量。可以一次展示队列、证据和方案，但不能要求用户同时
回答多个前提或决策。收到回答后重新计算队列，丢弃已经失效的问题。

用户的意图、优先级、取舍和接受风险只能询问；项目代码、数据库、文档和历史事实应主动检索。

## Workflow

### STEP 1: Load and validate context

完整读取 `design.md`、`review.md` 和两个 authoritative contracts。
确认 review 包含 Verdict、Findings、Dimension Summary 和 Current Readiness。
缺少必要结构时，提示先重新运行 `design-review`，不要猜测状态。

从现有文件读取：

- 当前 review round 和 refine run。
- 所有 finding 的稳定 `F-###`、Severity、Origin 和 Status。
- 已有 `DR-*`、风险记录和 changed sections。
- 历史 Closure 中已经终态的 finding。

本次运行编号使用 `R<n>-refine-<n>`。同一会话处理多个事项不重复递增。

### STEP 2: Explore only relevant project facts

围绕当前 findings 和候选 changed sections 读取：

1. 项目说明和约束。
2. 相关组件、接口、模型和测试。
3. 相关数据库 schema；需要时使用 `db-explorer`。
4. 外部 API、消息或协议的现有集成模式。

STOP when 能区分可由证据唯一修复的问题和必须由用户选择的问题。
不要因为 refine 再次全量探索无关项目区域。

### STEP 3: Build the resolution queue

按根因和依赖合并相关 finding，并分类：

| Type | Criteria | Interaction |
|------|----------|-------------|
| `Decision required` | 存在多个合理方案或持久取舍 | 逐项确认 |
| `Boundary rejection` | 用户可能明确驳回长期 concern | 逐项确认 |
| `Risk acceptance` | 用户可能接受未消除风险 | 逐项确认 |
| `Direct repair` | 可由现有 Decision 和项目事实唯一推导 | 批量确认 |
| `Blocked` | 缺少用户输入或外部前提 | 只问最上游前提 |
| `Minor optional` | 不阻断的打磨 | 批量修复或延期 |

进入讨论前，将相关 finding 标为 `in-discussion`。重复运行时跳过所有终态 finding。

多个 finding 指向同一个设计选择时合并为一个队列项，但每个 `F-###` 在 readiness 表中
保留独立状态，并允许共同指向同一个 `DR-*`。

### STEP 4: Process decisions one at a time

对 `Decision required` 按顺序执行：

1. 列出 review 来源、设计位置和项目证据。
2. 区分已知事实、已确认约束和仍需用户决定的前提。
3. 存在未决前提时，只询问最上游的一个问题，然后 STOP。
4. 前提明确后，比较 1-3 个真实可行方案和成本、风险、兼容性。
5. 给出推荐和理由，只要求用户确认当前选择。

只有决策涉及尚未确定的外部技术、第三方服务或协议时才进行外部调研。
内部字段、章节传播和项目已有模式不做无信号的 Web 搜索。

用户选择非推荐方案但没有说明约束时，下一轮只追问该约束，不同时进入下一个事项。

### STEP 5: Process boundary rejection and risk acceptance

用户明确驳回 finding 时，先检查已有 Decision：

- 已有 Decision 足够清晰：复用该 `DR-*`，不新增记录。
- 已有 Decision 表述不足：原位补充 `Rejected concern` 和 `Revisit when`。
- 没有持久边界：创建新的语义 `DR-*`。

边界 Decision 必须直接说明已确认的设计边界，而不是只写“Review F-### 被拒绝”。格式：

```markdown
### DR-<semantic-boundary>

**Decision**: <已确认的设计边界>。

**Rationale**: <用户确认的业务或架构理由>。

**Constraints**: <成立前提>。

**Rejected concern**: <不纳入的 concern，不展开实现细节>。

**Revisit when**: <需要重新评估的触发条件>。
```

用户在当前消息中已经明确驳回并给出理由时，该消息就是确认，不重复询问。
理由不充分时，只追问拒绝依据。完成后把 finding 标为 `rejected`，
Resolution ref 指向对应 `DR-*`。

用户选择接受风险时，必须记录风险、影响、接受理由和重新评估条件：

- 风险形成长期架构或运维边界：创建或原位更新 `DR-*`，在 Constraints 中写明风险。
- 风险只属于当前交付：写入 `Accepted / Deferred Risks`，不创建 Decision。

完成后把 finding 标为 `accepted-risk`，Resolution ref 指向 Decision 或风险条目。
Blocker/Major 的风险接受必须逐项确认，不得通过 Minor 批处理隐式接受。

### STEP 6: Batch direct repairs

所有上游 Decision 完成后，重新计算 Direct repair 列表。一次展示：

- finding ID。
- 修改章节。
- 项目或 Decision 依据。
- 预期变化。

然后只问一个问题：是否统一应用这些修复。

用户不同意全部应用时，只询问存在异议的具体项，并将其重新分类为
`Decision required`、`Boundary rejection`、`Risk acceptance` 或 `deferred`。

Direct repair 只更新对应章节。完成一致性检查后状态设为 `resolved`，
Resolution ref 直接指向章节，不创建 Decision。

### STEP 7: Process Minor items

只含 Minor 的可选项可以打包询问全部修复或全部延期。
用户选择延期时，在 `Accepted / Deferred Risks` 记录 finding、理由和重新评估条件，
状态设为 `deferred`。

Minor 与 Blocker/Major 共用同一个根因时，不得通过延期绕过阻断问题。

### STEP 8: Apply design updates

每个已确认事项按以下顺序更新：

1. 创建或原位更新符合准入条件的 `DR-*`。
2. 更新 canonical template 对应的当前设计章节。
3. 行为、错误语义或边界变化时同步更新 AC。
4. 更新 Testing 对 AC 的覆盖。
5. 更新 Current Readiness 的 Status、Resolution ref 和 Changed sections。

Decision 准入条件：

- 长期影响 Architecture、公共契约、数据所有权、安全、可靠性、迁移或运维边界。
- 存在多个合理方案，或结论非显然，后续实现者需要理解理由。
- 只修改具体章节不足以解释为什么必须这样设计。

不符合准入条件的 finding 不进入 Decision Record。design 中不再创建 `Review Resolution`。
被拒方案只写名称和拒绝理由，不写库版本、端点、配置值、DDL 或其他实现细节。

### STEP 9: Run consistency checks

**General propagation:**

- 每个 Goal 都对应至少一个组件和一条 AC。
- Non-Goals、Discovery 和 Scope Decisions 与当前 Decision 无冲突。
- 状态名、组件名、字段名和 Contract ID 跨章节一致。
- 每个行为变化同步到 Error Handling、AC 和 Testing。

**Data Model:**

- Change Summary 覆盖所有 ADD/MODIFY/REMOVE 数据对象。
- 新对象是完整目标模型，已有对象只描述变化。
- Persistence、Domain、Migration 和 Concurrency 内容归属正确。
- SQL 路径存在，design 没有复制第二份完整 DDL/DML。
- 迁移、回填、回滚和数据保护一致。

**Interfaces:**

- Change Summary 覆盖所有变化 API、消息和协议。
- 每个契约有稳定 Contract ID 和适用的请求、响应、错误、投递或兼容语义。
- `MODIFY` 和 `REMOVE` 明确兼容影响。
- Data Flow 中跨边界调用都引用 Contract ID。

**Readiness:**

- 每个当前 finding 有有效状态和证据引用。
- 任一非终态存在时，Overall 为 `in-progress` 或 `blocked`。
- Blocker/Major 全部终态且 pre-closure audit 通过后，Overall 为 `ready-for-closure`。
- 仅处理 Pass 的 Minor 且全部终态时，可以设置 `ready / Go`。
- 经 Blocker/Major refine 后不得直接设置 `Go`。

### STEP 10: Check language of changed content

在 pre-closure audit 前重新读取实际 `design.md` 和 `review.md`。默认只检查本次 changed sections、
Current Readiness 的本次更新及其相邻段落；用户明确要求清理全文时才扩大范围，不为统一文风
改写其他内容。

执行文档语言关卡：

1. 中文正文是否自然、直接？英文只保留模板字段、技术术语、路径、ID、Origin 和固定枚举。
   IF 否，THEN 只改对应句子。
2. 修改后的设计是否写清参与者、条件、动作、结果和失败行为；readiness 说明是否写清证据和
   下一步？IF 否，THEN 改成具体陈述，不用抽象流程词代替实际关系。
3. 是否存在逐词翻译、商业黑话、流程隐喻、模板化开场或空泛总结？IF 是，THEN 删除姿态层，
   保留设计与评审含义；正常技术术语不做机械替换。
4. 用户决定、事实、数字、Decision、finding ID、Severity、Origin、状态、字段、Contract ID、
   AC、路径、引用和责任主体是否保持不变？IF 否，THEN 恢复原意并缩小改写范围。
5. 语言修改是否新增、删除、弱化或扩大了设计边界、风险、契约、finding 或 Closure 结论？
   IF 是，THEN 撤销该修改，并依据已确认内容重新表述。

任一语言修改完成后重新执行 STEP 9。检查通过后才能进入 pre-closure audit。

### STEP 11: Run pre-closure audit

在设置 `ready-for-closure` 之前，执行一次 pre-closure audit。它是进入独立 review 的准入检查，
不是新的 Verdict，也不能替代 `design-review` 的独立判断。

同一 refine run 只启动这一套 audit 流程；audit 发现的问题回到当前 resolution queue，
不会在中间启动新的 Closure Review。

执行 [review lifecycle](../design-review/references/review-lifecycle.md) 的 Pre-closure Audit Contract：

1. 从 Review Snapshot 的原 Issue、Evidence 和期望结果为每条 finding 提炼 Closure test；Recommendation
   只提供候选修法，不把其中某个实现绑定为唯一通过条件。
   IF 无法写成可验证条件，THEN 保持 `in-progress`，不得用当前 `resolved` 标签代替证明。
2. 根据本次 `Changed sections` 建立影响矩阵，覆盖预先列出的关联章节和共享同一模型、身份、容量、错误或事务约束的
   相关契约。对 preview/commit、create/update、forward/rollback、producer/consumer 等成对契约检查对称性；
   有意差异必须由 Decision、Error Handling、AC 和 Testing 共同说明。
3. 用原 finding 的证据场景或最小反例复测当前设计。反例仍成立时重新打开原 finding，不创建同根新 finding。
4. 在 dry-run 前按 lifecycle 判定下一轮模式。Goals、Non-Goals、核心架构、数据所有权、公共契约、安全、
   权限、外部集成、迁移边界或相关项目事实变化时，记录 `Next review mode: Full` 并执行全维度 dry-run；
   其余情况记录 `Next review mode: Closure`，只检查 finding、changed sections 和
   Impact matrix 中预先列出的关联章节。
5. 任一检查失败时保持 `in-progress`，把缺口放回当前 resolution queue；直接修复批量处理，需要用户决定的事项
   继续逐项确认。修复后从第 1 项重跑，不在中间启动独立 review。
6. 全部通过后，按 lifecycle 的固定字段写入 audit 结果、模式依据、影响矩阵和 Finding Closure Proof，
   再将 Current Readiness 设为 `ready-for-closure`。

`design-review` 必须独立复核这些证明；pre-closure audit 不能替代 Full 或 Closure Review。

### STEP 12: Validate format and route

运行：

```bash
node <clarifying-skill>/scripts/check-markdown-lines.mjs <design.md> <review.md>
```

修复所有非豁免超长行。然后更新 Current Readiness 的 Next step 和报告末尾的
Recommended Next Step：

- `ready-for-closure`：按 `Next review mode` 运行 `design-review` Full 或 Closure Review。
- `ready / Go`：进入 `spec-plan`。
- `blocked`：说明缺少的唯一上游输入。

## Error handling

| Scenario | Action |
|----------|--------|
| review 缺少必要结构 | 停止并要求重新运行 design-review |
| design 不符合 canonical template | 作为 Direct repair 提交批量确认 |
| Decision ID 重复 | 停止写入，先合并为语义唯一 ID |
| finding 无法形成 Closure test | 保持 in-progress，回到原 finding 澄清根因和闭合条件 |
| 成对契约约束不一致 | 保持 in-progress，修复传播或确认有意差异 |
| 写入中断 | 保持非终态，并在 Note 说明中断位置 |
| 用户暂不决策 | 标记 blocked，输出当前唯一阻塞项 |

## Verification

- [ ] 已读取 canonical design template 和 review lifecycle。
- [ ] 所有非终态 finding 已分类。
- [ ] Decision required、Boundary rejection 和 Risk acceptance 已逐项确认。
- [ ] Direct repair 已批量展示并获得一次确认。
- [ ] Decision 使用稳定 `DR-*`，同一主题原位更新。
- [ ] 用户明确驳回的持久边界包含 `Rejected concern` 和 `Revisit when`。
- [ ] Direct repair 没有进入 Decision Record。
- [ ] design 中没有新增 `Review Resolution` 或 `Revised` 块。
- [ ] Data Model 和 Interfaces 符合 canonical schema。
- [ ] AC、Testing 和受影响章节已完成传播。
- [ ] 每条 finding 都有 Closure test、Resolution evidence、Counterexample check 和结果。
- [ ] 成对契约的共享约束一致，或差异已有 Decision、AC 和 Testing。
- [ ] Next review mode、Mode trigger 与 lifecycle 条件一致，对应 dry-run 已通过。
- [ ] Review Snapshot 未被修改。
- [ ] Current Readiness 状态、证据、changed sections 和下一步一致。
- [ ] 已检查本次改动及相邻段落的语言，没有为统一文风修改无关章节。
- [ ] 语言修改没有改变 Decision、finding、状态、契约、风险、AC 或 Closure 结论。
- [ ] design.md 与 review.md 行宽校验通过。

## References

- [Canonical design document template](../clarifying/assets/design-doc-template.md)
- [Review lifecycle contract](../design-review/references/review-lifecycle.md)
- [Boundary and multi-round examples](references/examples.md)
- [Decision 下游内容示例](references/pollution-examples.md)
- [Decision verification patterns](references/verification-patterns.md)
