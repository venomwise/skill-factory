---
name: spec-exec
description: Implement code tasks from specs/<topic>/tasks.md and track progress by updating checkboxes. Use when implementing a spec plan, resuming spec execution, or checking spec progress.
# metadata:
#   short-description: Spec exec (tasks.md)
---

# Spec Exec Skill

## When to use

- 实现 spec `tasks.md` 中 `## Tasks` 下列出的任务
- 通过更新任务 checkbox 状态跟踪执行进度

## When not to use

- 你只需要编辑单个任务行而不执行工作
- Spec 没有 `tasks.md` 或没有 `## Tasks` 段落

## Inputs

- Spec root: `specs/`
- Target: `specs/<topic>/tasks.md`
- Acceptance criteria: `specs/<topic>/requirements.md`

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
- 任务行是带数字前缀的 checkbox 列表项：Phase 用 `N.`，sub-task 用 `N.M`。已完成 (`[✅]`) 和可选 (`*`) 变体遵循相同结构。
- 缩进的描述行和 `_Requirements:` 行是元数据，not tasks。
- Phase 标题上的 `_Requirements:` 行是可追溯性摘要；使用各个 sub-task 上的行作为可执行指导。
- 可选 Phase 下的 sub-tasks 继承 Phase 的可选状态，在 MVP 模式下与它一起跳过。

## Outputs

- 更新的 `tasks.md`，带有完成标记：
  - Normal task: `- [✅]`
  - Optional task: `- [✅]*`

## Execution stance

在 spec 执行期间，作为自主实现 agent 行动。

将 `tasks.md` 视为执行计划，将 `requirements.md` 视为验收真相来源。If 当前任务和引用的 requirements 清晰且可执行，proceed without asking the user for confirmation。

When `requirements.md` 对所需细节存在歧义或保持沉默，consult `specs/<topic>/design.md`（if present）作为背景上下文。Never 使用 `design.md` 作为 `requirements.md` 中验收标准的替代品；if 标准本身缺失，按 *Blocker escalation* 升级（type: underspecified task）。

一次性设置选择（例如 Workflow step 6 中的 MVP vs. Full 模式）not 常规确认，may 在运行开始时询问一次。

## Checkpoints

Checkpoint tasks 是验证任务，not 用户批准门。

对于 checkpoint，验证其范围内已完成任务引用的 requirements 是否正确实现。使用 `tasks.md` 识别相关的 requirement IDs，使用 `requirements.md` 作为验收来源。

If validation passes，标记 checkpoint complete and continue。Stop only if validation fails, required resources are unavailable, or the spec is inconsistent。

## Evidence-based validation

Before marking a task or checkpoint complete，尽可能使用具体证据验证：

- Run 任务中列出的显式验证命令。
- Run 相关测试、类型检查、linters 或冒烟测试（if available）。
- Inspect 修改的文件以确保请求的行为存在。
- Compare 实现与引用的验收标准。
- If 验证命令 cannot be run，explain why and use 最强的可用替代检查。

Do not 仅基于未经支持的假设就标记任务完成。

## Blocker escalation

Do not 询问用户常规确认、实现偏好或继续的权限。

Stop and ask the user only if 执行真正被阻塞，such as:

- `tasks.md` 和 `requirements.md` 相互冲突。
- 下一个任务规格不足，cannot 从 `tasks.md` 和 `requirements.md` 解决。
- 针对引用 requirements 的验证失败，且失败 cannot 在当前任务内安全修复。
- 完成任务 would require 更改已批准的 requirements。
- 任务需要破坏性或不可逆操作，such as 删除用户数据、重写历史、删除数据库表或移除大量不相关代码。
- Required credentials, services, files, or environment dependencies are unavailable。

When blocked，do not 询问模糊问题如 "Should I continue?"

Instead，使用此结构化模板报告：

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

1. 解析目标 tasks 文件。
   - If 用户提供了 `specs/<topic>/tasks.md`，使用它。
   - If 用户 did not provide tasks 路径，精确回复：`请指定 tasks.md 文件路径（例如 specs/<topic>/tasks.md）。` Then 结束流程。Do not 搜索、列举或推断 spec。
   - If 用户提供的路径 does not exist，结束流程并推荐先运行 **`spec-plan`** skill 以产出 `specs/<topic>/tasks.md`。
2. 打开 `specs/<topic>/tasks.md` 并定位 `## Tasks` 段落。
   Also 打开 `specs/<topic>/requirements.md`（通常在 Overview 中链接）以获取验收标准。
   `_Requirements: N.M_` 行引用 `requirements.md` 中的验收标准，must be met。
3. 扫描进度并恢复：
   - 计数已完成 `[✅]` 和剩余 `[ ]` 任务。
   - 识别第一个未完成的任务并从那里恢复。
   - If all tasks are complete，报告并停止。
4. 通过搜索包含 `[ ]*` 的 checkbox 行（at any indentation level — both Phase and sub-task lines）检测 `## Tasks` 内的可选任务。
5. If any optional tasks exist，询问：
   `当前任务列表将部分任务（如：单元测试、文档编写）标记为可选，以便集中精力优先实现核心功能。A. 保留可选任务 (MVP) B. 执行所有任务`
6. **Triage** — For each task in `## Tasks` order，确定其处置：
   > **REMINDER**: After completing each task you MUST update its checkbox in `tasks.md` before starting the next one. Never accumulate multiple completed tasks without writing them back.
   - Skip tasks already marked `- [✅]`。
   - If MVP mode was chosen，skip tasks marked with `- [ ]*`。
   - If an optional Phase (`- [ ]*`) is skipped，skip all nested sub-tasks under that Phase。
   - Identify checkpoint/verification tasks by keywords such as **"Checkpoint"**, **"Verify"**, or **"检查点"**。Handle them per the *Checkpoints* section above。
   - Otherwise，proceed to step 7 to implement the task。
7. **Implement & Validate** — For the current task:
   - Read 标题下方的缩进描述行，并使用 `_Requirements: ..._` 行作为显式指导。
   - Review referenced files/modules before changes to understand current behavior and constraints。
   - Implement the task in the codebase following project conventions。
   - Validate the result per the *Evidence-based validation* section。If validation fails，follow *Blocker escalation* (type: validation_failure)。
   - When marking a checkpoint or validation-only task complete，briefly record the validation evidence (command run + key result) in your reply，so the audit trail survives interruption。
8. **Mark completion** — **CRITICAL: Update `tasks.md` NOW, before doing anything else.**
   > This is the most important step in the loop. You MUST write the checkbox change to `tasks.md` for the task you just completed BEFORE moving on to the next task. Failing to do so means progress is lost on interruption.
   - Normal task: change `- [ ]` to `- [✅]`
   - Optional task: change `- [ ]*` to `- [✅]*`
   - Do not mark tasks that failed, were interrupted, or were skipped。
   - **ONE task, ONE write.** Never accumulate multiple completed tasks into a single `tasks.md` update。
   - Then return to step 6 for the next task。
9. When all sub-tasks under a Phase are completed，mark the Phase line as `- [✅]`。
    - In MVP mode，a Phase is complete when all non-optional sub-tasks are done。
    - If a Phase has no sub-tasks (e.g., a Checkpoint)，mark it only after it is completed。
10. After all required tasks are complete，perform repository guidance sync:
    - In MVP mode，treat execution as complete when all non-optional tasks are finished；unchecked optional tasks do not block this step。
    - Check whether the project root contains `AGENTS.md`。
    - If `AGENTS.md` exists，review the work completed in this run and update the file to reflect any changed contributor guidance，such as project structure, development commands, verification flow, or repository conventions introduced by the implementation。
    - Keep the update scoped to guidance affected by the completed work；do not rewrite unrelated sections。

## Verification

- Before marking a task as `[✅]`，perform the validation described in step 7 or confirm the code is runnable。
- Validation tasks (e.g., checkpoints, verify tasks, or manual smoke tests) must be executed。
- Checkpoint tasks are completed by evidence-based validation，not by user confirmation。
- Do not ask the user to confirm successful checkpoints unless execution is blocked。
- Only items under `## Tasks` are modified。
- Optional tasks remain unchecked when MVP mode is chosen。
- Phase items are marked only after all their sub-tasks are completed (step 9)。
- After all required tasks are finished，if the project root contains `AGENTS.md`，it has been reviewed and updated to match the completed work。In MVP mode，this check happens after non-optional tasks are done。

## Safety & guardrails

- Never auto-select a `tasks.md`。The path must come from the user。
- Never mark tasks as done unless execution completed successfully。
- Do not alter task numbering, titles, or descriptions。
- Do not use user interaction as a substitute for reading `tasks.md` and referenced requirements。
- Do not introduce requirement changes during execution。If a requirement change seems necessary，stop and report it as a blocker。
- Do not ask for permission to continue after successful validation。
- Stop and ask the user only if execution is blocked by a blocker defined above。
- If a task fails，keep it as `- [ ]` and escalate via the *Blocker escalation* template (type: validation_failure)，offering at minimum these options: (a) fix in place now, (b) defer and continue to the next task, (c) abort the run。
- If a failed task produces artifacts required by later tasks，include this cascade risk in the blocker `risk` field so the user can choose accordingly。

## References

- [Spec Plan Skill](../spec-plan/SKILL.md) — tasks.md 格式规范
- [Tasks template](../spec-plan/assets/tasks.template.md) — tasks.md 结构模板
