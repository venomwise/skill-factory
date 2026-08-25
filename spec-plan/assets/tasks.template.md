# Implementation Plan: <Project Name>

## Overview

This implementation plan is driven by the approved design and Acceptance Criteria in [design.md](design.md).

<!-- Describe the phase count, execution order, and key implementation decisions. -->

<总结实施策略、执行顺序和关键技术决策。>

## Tasks

<!-- FORMAT RULES:

1. Phase headings are checkbox list items:
   - [ ] 1. Phase 1: Title

2. Task details are indented bullets under the task title.

3. Optional work uses `*` immediately after the checkbox:
   - [ ]* 1.3 Optional task
   Sub-tasks under an optional Phase inherit its optional status.

4. Every executable sub-task and standalone Checkpoint has at least one
   traceability line. Phase containers do not need one:
   - _Acceptance: AC-domain-behavior_
   - _Design: Proposed Solution / Components_
   A task may contain both. Acceptance IDs and design section paths must exist in design.md.
   Design paths contain exact Markdown headings only. Names that appear only in
   prose, lists, tables, or code are not valid path segments.

5. Use `_Acceptance:` for observable behavior and `_Design:` for scaffolding,
   architecture, or internal refactoring. Do not invent an AC for internal work.

6. Checkpoints list concrete commands, AC IDs, and pass conditions. They are
   execution-agent validation tasks, not user approval gates.
-->

- [ ] 1. Phase 1: <Phase Title>
  - [ ] 1.1 <实现可观察行为>
    - Create or modify `src/<module>/<file>` with `<function/class>` implementing <行为>
    - _Acceptance: AC-<domain>-<behavior>_
    - _Design: Proposed Solution / Components_
  - [ ] 1.2 <内部结构或脚手架任务>
    - Create or modify `src/<module>/<file>` according to the approved component design
    - _Design: Proposed Solution / Architecture_
  - [ ]* 1.3 为本阶段编写测试
    - Test normal, error, and boundary behavior covered by the referenced AC
    - _Acceptance: AC-<domain>-<behavior>, AC-<domain>-<failure>_
    - _Design: Testing_

- [ ] 2. Phase 2: <Phase Title>
  - [ ] 2.1 <任务标题>
    - Create or modify `src/<module>/<file>` with `<function/class>` implementing <行为或设计>
    - _Acceptance: AC-<domain>-<behavior>_

- [ ] 3. Checkpoint - Verify <scope>
  - Run `<test or validation command>`
  - Verify `AC-<domain>-<behavior>` and `AC-<domain>-<failure>`
  - Pass when <observable expected result>; stop on failure, spec conflict, or missing dependency
  - _Acceptance: AC-<domain>-<behavior>, AC-<domain>-<failure>_
  - _Design: Testing_

- [ ]* 4. Optional Phase: <Phase Title>
  - [ ] 4.1 <任务标题>
    - <具体步骤或交付物>
    - _Design: <existing design section path>_

## Notes

- Tasks marked with `*` may be skipped in MVP mode.
- Keep task numbering stable after execution begins.
- `design.md` is the source of truth; do not copy full AC text into this file.
