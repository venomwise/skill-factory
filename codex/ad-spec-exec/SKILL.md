---
name: ad-spec-exec
description: Implement code tasks from .codex/specs/<spec>/tasks.md and track progress by updating checkboxes. Use when implementing a spec plan, resuming spec execution, or checking spec progress.
# metadata:
#   short-description: Spec exec (tasks.md)
---

# Spec Exec Skill

## When to use

- Implement tasks listed under `## Tasks` in a spec `tasks.md`.
- Track progress by updating task checkbox status after execution.

## When not to use

- You only need to edit a single task line without executing work.
- The spec has no `tasks.md` or no `## Tasks` section.

## Inputs

- Spec root: `.codex/specs/`
- Target: `.codex/specs/<spec>/tasks.md`
- Acceptance criteria: `.codex/specs/<spec>/requirements.md`

## Expected tasks.md format

```
- [ ] 1. Phase 1: Title
  - [ ] 1.1 Sub-task title
    - Description line
    - _Requirements: 1.1, 1.2_
  - [ ]* 1.2 Optional sub-task
- [ ] 2. Checkpoint - Verify scope
- [ ]* 3. Optional Phase: Title
  - [ ] 3.1 Sub-task under optional phase
```

Key rules:
- Task lines are checkbox list items with numeric prefixes: `N.` for phases and `N.M` for sub-tasks. Completed (`[✅]`) and optional (`*`) variants follow the same structure.
- Indented description lines and `_Requirements:` lines are metadata, not tasks.
- `_Requirements:` lines on Phase headings are summaries for traceability; use the lines on individual sub-tasks as actionable guidance.
- Sub-tasks under an optional Phase inherit the Phase's optional status and are skipped with it in MVP mode.

## Outputs

- Updated `tasks.md` with completion markers:
  - Normal task: `- [✅]`
  - Optional task: `- [✅]*`

## Workflow

1. List available specs (PowerShell):
   `Get-ChildItem .codex/specs -Directory | Select-Object -ExpandProperty Name`
2. Ask the user to choose:
   `选择需要执行的规格：1. <spec>`
   If multiple specs exist, list all options in order.
3. Open `.codex/specs/<spec>/tasks.md` and locate the `## Tasks` section.
   Also open `.codex/specs/<spec>/requirements.md` (usually linked in Overview) for acceptance criteria.
   `_Requirements: N.M_` lines reference acceptance criteria in `requirements.md` and must be met.
4. Scan progress and resume:
   - Count completed `[✅]` and remaining `[ ]` tasks.
   - Identify the first incomplete task and resume from there.
   - If all tasks are complete, report and stop.
5. Detect optional tasks inside `## Tasks` by searching for checkbox lines containing `[ ]*` (at any indentation level — both Phase and sub-task lines).
6. If any optional tasks exist, ask:
   `当前任务列表将部分任务（如：单元测试、文档编写）标记为可选，以便集中精力优先实现核心功能。A. 保留可选任务 (MVP) B. 执行所有任务`
7. **Triage** — For each task in `## Tasks` order, determine its disposition:
   > **REMINDER**: After completing each task you MUST update its checkbox in `tasks.md` before starting the next one. Never accumulate multiple completed tasks without writing them back.
   - Skip tasks already marked `- [✅]`.
   - If MVP mode was chosen, skip tasks marked with `- [ ]*`.
   - If an optional Phase (`- [ ]*`) is skipped, skip all nested sub-tasks under that Phase.
   - Identify checkpoints by the keyword **"Checkpoint"** or **"检查点"** in the task title. If the task is a checkpoint, pause and summarize progress, then ask the user to confirm before continuing.
   - Otherwise, proceed to step 8 to implement the task.
8. **Implement & Validate** — For the current task:
   - Read the indented description lines beneath the title and use `_Requirements: ..._` lines as explicit guidance.
   - Review referenced files/modules before changes to understand current behavior and constraints.
   - Implement the task in the codebase following project conventions.
   - Validate the task:
     - If the task includes explicit validation steps, execute them.
     - If the task itself is a validation step (e.g., manual smoke test), perform it.
     - Look up referenced requirement IDs in `requirements.md` and verify the implementation satisfies each acceptance criterion.
9. **Mark completion** — **CRITICAL: Update `tasks.md` NOW, before doing anything else.**
   > This is the most important step in the loop. You MUST write the checkbox change to `tasks.md` for the task you just completed BEFORE moving on to the next task. Failing to do so means progress is lost on interruption.
   - Normal task: change `- [ ]` to `- [✅]`
   - Optional task: change `- [ ]*` to `- [✅]*`
   - Do not mark tasks that failed, were interrupted, or were skipped.
   - **ONE task, ONE write.** Never accumulate multiple completed tasks into a single `tasks.md` update.
   - Then return to step 7 for the next task.
10. When all sub-tasks under a Phase are completed, mark the Phase line as `- [✅]`.
    - In MVP mode, a Phase is complete when all non-optional sub-tasks are done.
    - If a Phase has no sub-tasks (e.g., a Checkpoint), mark it only after it is completed.

## Verification

- Before marking a task as `[✅]`, perform the validation described in step 8 or confirm the code is runnable.
- Validation tasks (e.g., manual smoke tests) must be executed.
- Only items under `## Tasks` are modified.
- Optional tasks remain unchecked when MVP mode is chosen.
- Phase items are marked only after all their sub-tasks are completed (step 10).

## Safety & guardrails

- Never mark tasks as done unless execution completed successfully.
- Do not alter task numbering, titles, or descriptions.
- Stop and ask the user if execution is blocked or unclear.
- If a task fails, keep it as `- [ ]`, explain the failure, and ask whether to fix it first or proceed to the next task.
- If a failed task produces artifacts required by later tasks, warn the user that skipping may cause cascading failures.

## References

- [Spec Plan Skill](../ad-spec-plan/SKILL.md) — tasks.md 格式规范
- [Tasks template](../ad-spec-plan/assets/tasks.template.md) — tasks.md 结构模板
