# <Topic> Design Review

> 对 `<path/to/design.md>` 的设计评审。Verdict 是评审时快照；当前流程状态见 Current Readiness。

## Verdict

- **Review round**: <R1 / R2 / ...>
- **Review mode**: <Full / Closure>
- **Mode reason**: <首次评审、refine 后闭环或升级原因>
- **Overall at review**: <Pass / Revise / Reject>
- **spec-plan readiness at review**: <Go / Conditional / No-Go>
- **Findings**: <N Blocker / N Major / N Minor>
- **Summary**: <最重要的结论>

## Scope Reviewed

- **Reviewed file**: `<path/to/design.md>`
- **Project baseline**: <commit、数据库事实或其他基线>
- **Project context consulted**: <关键文件和证据>
- **Changed sections**: <Closure 模式填写；Full 模式写 All>
- **Acceptance Criteria checked**: <数量、Goal 覆盖和可测试性结论>
- **Canonical contract**: `brainstorming/assets/design-doc-template.md`
- **Rubric**: D1-D7

## Findings

按 Blocker、Major、Minor 排序。空分组不输出。finding ID 与严重性无关，跨轮保持稳定。

### [F-001] <标题>

- **Severity**: <Blocker / Major / Minor>
- **Introduced in**: <R1 / R2 / ...>
- **Origin**: <initial-review / refine-regression / dependency-unlocked / baseline-miss / context-change>
- **Location**: `design.md §<section>` <和/或项目文件>
- **Issue**: <具体错误、缺失或需要确认的风险>
- **Evidence**: <设计、代码、数据库或契约证据>
- **Recommendation**: <具体修复或必须做出的决策>

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

| Finding | Severity | Status | Resolution ref | Changed sections | Note |
|---------|----------|--------|----------------|------------------|------|
| F-001 | <severity> | <status> | <DR-id / §section / risk ref / None> | <sections / None> | <short note> |

## Closure

<!-- 首轮省略。复审时保留所有历史 finding 的紧凑闭合结果。 -->

| Finding | Title | Result | Resolution ref | Summary |
|---------|-------|--------|----------------|---------|
| F-001 | <title> | <closed / continued> | <evidence> | <一句结果> |

## Accepted / Deferred Risks

<!-- 没有则省略。Minor 延期和非持久接受风险记录在这里。 -->

- **F-###**: <accepted-risk / deferred>，<理由和重新评估条件>。

## Recommended Next Step

<!--
本节由 Current Readiness 派生：
- Blocker/Major -> design-refine，完成后 Closure Review。
- Pass with Minor -> 修复或延期。
- ready / Go -> spec-plan。
-->

<当前推荐动作>。
