---
name: spec-exec
description: Implement code tasks from specs/<spec>/tasks.md and track progress by updating checkboxes. Use when implementing a spec plan, resuming spec execution, or checking spec progress.
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

- Spec root: `specs/`
- Target: `specs/<spec>/tasks.md`
- Acceptance criteria: `specs/<spec>/requirements.md`

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

## Execution stance

During spec execution, act as an autonomous implementation agent.

Treat `tasks.md` as the execution plan and `requirements.md` as the acceptance source of truth. If the current task and referenced requirements are clear and executable, proceed without asking the user for confirmation.

When `requirements.md` is ambiguous or silent on a needed detail, consult `specs/<spec>/design.md` (if present) as background context. Never use `design.md` as a substitute for an acceptance criterion in `requirements.md`; if the criterion itself is missing, escalate per *Blocker escalation* (type: underspecified task).

One-time setup choices (e.g., MVP vs. Full mode in Workflow step 6) are not routine confirmations and may be asked once at the start of a run.

## Checkpoints

Checkpoint tasks are validation tasks, not user approval gates.

For a checkpoint, verify that the requirements referenced by the completed tasks in its scope are correctly implemented. Use `tasks.md` to identify the relevant requirement IDs and `requirements.md` as the acceptance source.

If validation passes, mark the checkpoint complete and continue. Stop only if validation fails, required resources are unavailable, or the spec is inconsistent.

## Evidence-based validation

Before marking a task or checkpoint complete, validate it using concrete evidence whenever possible:

- Run explicit validation commands listed in the task.
- Run relevant tests, type checks, linters, or smoke tests if available.
- Inspect modified files to ensure the requested behavior exists.
- Compare implementation against referenced acceptance criteria.
- If a validation command cannot be run, explain why and use the strongest available alternative check.

Do not mark a task complete based only on an unsupported assumption.

## Blocker escalation

Do not ask the user for routine confirmation, implementation preferences, or permission to continue.

Stop and ask the user only if execution is genuinely blocked, such as:

- `tasks.md` and `requirements.md` conflict with each other.
- The next task is underspecified and cannot be resolved from `tasks.md` and `requirements.md`.
- Validation against referenced requirements fails and the failure cannot be safely fixed within the current task.
- Completing the task would require changing approved requirements.
- The task requires destructive or irreversible operations, such as deleting user data, rewriting history, dropping database tables, or removing large unrelated code.
- Required credentials, services, files, or environment dependencies are unavailable.

When blocked, do not ask a vague question like "Should I continue?"

Instead, report using this structured template:

```yaml
blocker:
  task: "<N.M> <title>"
  type: <conflict | underspecified | validation_failure | scope_change | destructive_op | missing_dependency>
  context:
    task_excerpt: "<relevant lines from tasks.md>"
    requirements: "<referenced requirement IDs and their criteria>"
  tried:
    - "<what you already attempted>"
  risk: "<why proceeding would violate the spec>"
  options:
    - "<option A the user can pick>"
    - "<option B the user can pick>"
  needed_from_user: "<minimum decision or input>"
```

## Workflow

1. List available specs (PowerShell):
   `Get-ChildItem .codex/specs -Directory | Select-Object -ExpandProperty Name`
2. Ask the user to choose:
   `选择需要执行的规格：1. <spec>`
   If multiple specs exist, list all options in order.
3. Open `specs/<spec>/tasks.md` and locate the `## Tasks` section.
   Also open `specs/<spec>/requirements.md` (usually linked in Overview) for acceptance criteria.
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
   - Identify checkpoint/verification tasks by keywords such as **"Checkpoint"**, **"Verify"**, or **"检查点"**. Handle them per the *Checkpoints* section above.
   - Otherwise, proceed to step 8 to implement the task.
8. **Implement & Validate** — For the current task:
   - Read the indented description lines beneath the title and use `_Requirements: ..._` lines as explicit guidance.
   - Review referenced files/modules before changes to understand current behavior and constraints.
   - Implement the task in the codebase following project conventions.
   - Validate the result per the *Evidence-based validation* section. If validation fails, follow *Blocker escalation* (type: validation_failure).
   - When marking a checkpoint or validation-only task complete, briefly record the validation evidence (command run + key result) in your reply, so the audit trail survives interruption.
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
11. After all required tasks are complete, perform repository guidance sync:
    - In MVP mode, treat execution as complete when all non-optional tasks are finished; unchecked optional tasks do not block this step.
    - Check whether the project root contains `AGENTS.md`.
    - If `AGENTS.md` exists, review the work completed in this run and update the file to reflect any changed contributor guidance, such as project structure, development commands, verification flow, or repository conventions introduced by the implementation.
    - Keep the update scoped to guidance affected by the completed work; do not rewrite unrelated sections.

## Verification

- Before marking a task as `[✅]`, perform the validation described in step 8 or confirm the code is runnable.
- Validation tasks (e.g., checkpoints, verify tasks, or manual smoke tests) must be executed.
- Checkpoint tasks are completed by evidence-based validation, not by user confirmation.
- Do not ask the user to confirm successful checkpoints unless execution is blocked.
- Only items under `## Tasks` are modified.
- Optional tasks remain unchecked when MVP mode is chosen.
- Phase items are marked only after all their sub-tasks are completed (step 10).
- After all required tasks are finished, if the project root contains `AGENTS.md`, it has been reviewed and updated to match the completed work. In MVP mode, this check happens after non-optional tasks are done.

## Safety & guardrails

- Never mark tasks as done unless execution completed successfully.
- Do not alter task numbering, titles, or descriptions.
- Do not use user interaction as a substitute for reading `tasks.md` and referenced requirements.
- Do not introduce requirement changes during execution. If a requirement change seems necessary, stop and report it as a blocker.
- Do not ask for permission to continue after successful validation.
- Stop and ask the user only if execution is blocked by a blocker defined above.
- If a task fails, keep it as `- [ ]` and escalate via the *Blocker escalation* template (type: validation_failure), offering at minimum these options: (a) fix in place now, (b) defer and continue to the next task, (c) abort the run.
- If a failed task produces artifacts required by later tasks, include this cascade risk in the blocker `risk` field so the user can choose accordingly.

## References

- [Spec Plan Skill](../spec-plan/SKILL.md) — tasks.md 格式规范
- [Tasks template](../spec-plan/assets/tasks.template.md) — tasks.md 结构模板
