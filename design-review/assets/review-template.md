# <Topic> Design Review

> 对 `<path/to/design.md>` 的设计评审。Verdict 是评审时快照；当前流程状态见 Current Readiness。

<!--
中文正文写清条件、受影响对象、证据和后果。不要用“存在一定风险”“不够完善”或
“进一步优化”代替具体 finding。固定字段、ID、路径、Severity、Origin 和状态枚举保持不变。
-->

## Verdict

- **Review round**: <R1 / R2 / ...>
- **Review mode**: <Full / Closure>
- **Mode reason**: <触发当前模式的具体 lifecycle 条件>
- **Overall at review**: <Pass / Revise / Reject>
- **spec-plan readiness at review**: <Go / Conditional / No-Go>
- **Findings**: <N Blocker / N Major / N Minor>
- **Summary**: <最重要的结论、直接原因和评审当时是否可进入下一步>

## Scope Reviewed

- **Reviewed file**: `<path/to/design.md>`
- **Project baseline**: <commit、数据库事实或其他基线>
- **Project context consulted**: <关键文件和证据>
- **Changed sections**: <Closure 模式填写；Full 模式写 All>
- **Acceptance Criteria checked**: <数量、Goal 覆盖和可测试性结论>
- **Canonical contract**: `clarifying/assets/design-doc-template.md`
- **Rubric**: D1-D7

## Findings

按 Blocker、Major、Minor 排序。空分组不输出。finding ID 与严重性无关，跨轮保持稳定。

### [F-001] <能指出具体问题对象的简短标题>

- **Severity**: <Blocker / Major / Minor>
- **Introduced in**: <R1 / R2 / ...>
- **Origin**: <initial-review / refine-regression / dependency-unlocked / baseline-miss / context-change>
- **Location**: `design.md §<section>` <和/或项目文件>
- **Issue**: <在什么条件下，哪个设计对象存在什么错误、缺失或风险>
- **Evidence**: <可核实的设计位置、代码、数据库或契约证据，以及它与期望行为的差异>
- **Recommendation**: <最小修改面：需要改动的章节或契约、预期需要同步的传播章节，或必须由用户做出的具体决策>

## Dimension Summary

状态：`OK` 无问题，`Attention` 有 finding，`Blocked` 存在 Blocker。

| Dimension | Status | Finding refs |
|-----------|--------|--------------|
| D1 Completeness | <OK/Attention/Blocked> | <F-001 / None> |
| D2 Usability | <OK/Attention/Blocked> | <refs> |
| D3 Conformance | <OK/Attention/Blocked> | <refs> |
| D4 Project Fit | <OK/Attention/Blocked> | <refs> |
| D5 Blind Spots | <OK/Attention/Blocked> | <refs> |
| D6 Over-Engineering | <OK/Attention/Blocked> | <refs> |
| D7 Optimization | <OK/Attention/Blocked> | <refs> |

## Current Readiness

本节表示 refine 后的当前状态，可由 `design-refine` 更新。上方 Review Snapshot 保持不变。

- **Review round**: <R1 / R2 / ...>
- **Refine run**: <None / R1-refine-1 / ...>
- **Overall**: <not-started / in-progress / blocked / ready-for-closure / ready>
- **spec-plan readiness**: <Go / Conditional / No-Go>
- **Next step**: <当前应执行的动作>
- **Pre-closure audit**: <Not run / Passed / Failed；仅在 refine 后填写>
- **Next review mode**: <Not determined / Full / Closure>
- **Mode trigger**: <命中的 Full Review 条件 / None>
- **Impact matrix**: <changed sections -> 预先列出的关联章节；仅在 refine 后填写>

| Finding | Severity | Status | Resolution ref | Changed sections | Note |
|---------|----------|--------|----------------|------------------|------|
| F-001 | <severity> | <status> | <DR-id / §section / risk ref / None> | <sections / None> | <证据或阻塞摘要> |

### Finding Closure Proof

<!-- pre-closure audit 未执行时省略。不得只根据上表的 status 宣称 finding 已闭合。 -->

| Finding | Closure test | Resolution evidence | Counterexample check | Result |
|---------|--------------|---------------------|----------------------|--------|
| F-001 | <由原 Issue/Evidence 和期望结果提炼> | <DR / section / AC / test / risk> | <原反例当前结果> | <passed / failed> |

## Closure

<!-- 首轮省略。复审时保留所有历史 finding 的紧凑闭合结果。 -->

| Finding | Title | Result | Resolution ref | Summary |
|---------|-------|--------|----------------|---------|
| F-001 | <title> | <closed / continued> | <evidence> | <一句结果> |

## Accepted / Deferred Risks

<!-- 没有则省略。Minor 延期和非持久接受风险记录在这里。 -->

- **F-###**: <accepted-risk / deferred>，<受影响对象、接受或延期理由，以及重新评估条件>。

## Recommended Next Step

<!--
本节由 Current Readiness 派生：
- Blocker/Major -> design-refine，完成后按 Next review mode 执行 Full 或 Closure Review。
- Pass with Minor -> 修复或延期。
- ready / Go -> spec-plan。
-->

<根据 Current Readiness 写明下一位处理者、需要执行的动作和完成条件>。
