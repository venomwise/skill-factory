---
name: spec-exec
description: >
  按照 specs/<topic>/tasks.md 和 design.md 实现任务，并通过更新 checkbox 跟踪进度。
  当用户提到执行 spec、实施任务计划、继续执行 spec、恢复执行、检查 spec 进度时，
  必须使用此技能。
---

# Spec Exec Skill

## When to use

- 实现 `tasks.md` 中 `## Tasks` 下列出的任务
- 依据 `design.md` 的 Acceptance Criteria 验证行为
- 通过更新任务 checkbox 跟踪进度

## When not to use

- 只需要编辑任务文本而不执行工作
- Spec 没有 `design.md`、`tasks.md` 或 `## Tasks`

## Inputs / Outputs

**Inputs:**

- 执行计划：`specs/<topic>/tasks.md`
- 权威设计与验收标准：`specs/<topic>/design.md`

**Outputs:**

- 完成的实现和验证证据
- 更新后的 `tasks.md`：普通任务使用 `- [✅]`，可选任务使用 `- [✅]*`

## Expected tasks.md format

```markdown
- [ ] 1. Phase 1: Title
  - [ ] 1.1 Implement behavior
    - Description line
    - _Acceptance: AC-config-precedence_
    - _Design: Proposed Solution / Components_
  - [ ]* 1.2 Optional test task
    - _Design: Testing_
- [ ] 2. Checkpoint - Verify scope
  - _Acceptance: AC-config-precedence_
```

关键规则：

- Phase 使用 `N.`，sub-task 使用 `N.M`；已完成和可选变体遵循相同结构。
- 缩进描述、`_Acceptance:` 和 `_Design:` 是任务元数据，不是任务。
- `_Acceptance:` 引用 `design.md` 中的 AC；`_Design:` 引用实际存在的设计章节。
- 每个可执行 sub-task 和独立 Checkpoint 至少包含一种追踪元数据；可同时包含两种。包含子任务的 Phase 行只是容器，不需要元数据。
- 可选 Phase 下的子任务继承可选状态。

## Execution stance

将 `tasks.md` 视为执行计划，将 `design.md` 视为已批准行为和设计方案的唯一权威来源。IF 当前任务及其引用清晰且可执行，THEN 自主推进，不请求常规确认。

Acceptance Criteria 只定义可观察结果；验证命令来自任务、Checkpoint、`Testing` 章节和项目工具链。不要把缺失的行为或 AC 当作实现偏好自行补齐。

## Strict contract gate

开始执行前必须验证：

1. `design.md` 包含 `## Acceptance Criteria`，AC ID 唯一且符合 `AC-<domain>-<behavior>`。
2. `tasks.md` 不包含 `_Requirements:`；命中即视为不受支持的旧格式。
3. 每个未完成的可执行 sub-task 和独立 Checkpoint 至少包含 `_Acceptance:` 或 `_Design:`；Phase 容器除外。
4. 所有 `_Acceptance:` ID 和 `_Design:` 章节路径都能在 `design.md` 中解析。
5. 任务内容与引用的 AC、设计章节没有冲突。

IF 任一检查失败，THEN 在修改代码前按 *Blocker escalation* 报告并停止。不要迁移旧 spec，也不要补写 AC。

## Checkpoints

Checkpoint 是验证任务，不是用户批准门。执行其中的命令，使用 `_Acceptance:` 定位必须验证的 AC，并以可观察证据判断通过条件。

IF 验证通过，THEN 标记完成并继续。仅当验证失败、必需资源不可用或 spec 内部不一致时停止。

## Evidence-based validation

在标记任务或 Checkpoint 完成前，尽可能：

- 执行任务或 Checkpoint 中的显式验证命令。
- 运行相关测试、类型检查、lint 或冒烟测试。
- 检查修改文件，确认引用的设计内容已实现。
- 将结果与 `_Acceptance:` 引用的 AC 逐条对比。
- IF 命令无法运行，THEN 说明原因并使用最强的替代检查。

不要仅凭任务文本或未经验证的假设标记完成。

## Blocker escalation

仅在真正阻塞时停下，例如：

- Spec 使用旧 `_Requirements:` 格式或缺少 Acceptance Criteria。
- 引用的 AC 或设计章节不存在、含糊或互相冲突。
- 任务无法从 `tasks.md` 和引用的 `design.md` 内容确定。
- 验证失败且无法在当前任务范围内安全修复。
- 完成任务需要改变已批准 AC 或设计。
- 需要破坏性操作，或缺少凭证、服务、文件、环境依赖。

使用以下结构报告：

```yaml
blocker:
  task: "<N.M> <title>"
  type: <legacy_spec | conflict | underspecified | validation_failure | scope_change | destructive_op | missing_dependency>
  context:
    task_excerpt: "<relevant lines from tasks.md>"
    acceptance: "<referenced AC IDs and criteria>"
    design: "<referenced design sections>"
  tried:
    - "<what you already attempted>"
  risk: "<why proceeding would violate the spec>"
  options:
    - "<option A>"
    - "<option B>"
  needed_from_user: "<minimum decision or input>"
```

## Workflow

1. **定位 tasks.md。**
   - IF 用户提供了路径，THEN 使用它。
   - IF 未提供，THEN 精确回复：`请指定 tasks.md 文件路径（例如 specs/<topic>/tasks.md）。`然后结束流程。不要搜索、列举或推断 spec。
   - IF 路径不存在，THEN 推荐先运行 `spec-plan`。
2. **读取完整 spec。** 打开 `tasks.md` 和同目录 `design.md`，定位 `## Tasks`、`## Acceptance Criteria`、`## Testing` 及任务引用的设计章节。
3. **执行严格关卡。** 按 *Strict contract gate* 检查；失败时在修改代码前停止。
4. **扫描并恢复进度。** 计数 `[✅]` 与 `[ ]`，从首个未完成任务继续；全部完成时报告并停止。
5. **选择执行模式。** IF 存在 `[ ]*`，THEN 仅询问一次：`当前任务列表包含可选任务。A. 保留可选任务（MVP） B. 执行所有任务`。
6. **Triage。** 按顺序处理任务：
   - 跳过 `[✅]`。
   - MVP 模式跳过 `[ ]*`；可选 Phase 被跳过时连同其子任务一起跳过。
   - 包含 sub-task 的 Phase 是容器，不直接实施；按顺序处理其子任务，并在步骤 9 汇总状态。
   - 通过 `Checkpoint`、`Verify`、`检查点` 识别验证任务。
   - 其他任务进入步骤 7。
7. **Implement & Validate。**
   - 阅读任务描述和 `_Acceptance:` / `_Design:` 元数据。
   - 修改前阅读相关文件与模块，理解现有行为。
   - 按项目约定实施，并按 *Evidence-based validation* 验证。
   - IF 失败，THEN 按 *Blocker escalation* 处理。
8. **立即写回进度。** 每完成一个任务，先把它改为 `[✅]` 或 `[✅]*`，再做任何其他动作。一个任务一次写入；失败、跳过或中断的任务不得标记完成。
9. **更新 Phase。** 一个 Phase 的适用子任务全部完成后标记 Phase；MVP 模式下未执行的可选子任务不阻塞 Phase 完成。无子任务的 Checkpoint 只在验证通过后标记。
10. **同步仓库指导。** 所有必需任务完成后，IF 根目录存在 `AGENTS.md`，THEN 只在本次实现改变项目结构、命令、验证流程或仓库约定时更新相关内容，不重写无关章节。

## Verification

- [ ] 执行前通过严格 contract gate
- [ ] 每个可执行 sub-task 和独立 Checkpoint 依据可解析的 `_Acceptance:` 或 `_Design:` 执行
- [ ] 标记完成前有测试、命令、文件检查或其他具体证据
- [ ] Checkpoint 实际执行且按引用的 AC 验证
- [ ] 每完成一个任务立即写回 checkbox
- [ ] MVP 模式下可选任务保持未勾选
- [ ] Phase 仅在其适用子任务完成后标记
- [ ] 必需任务完成后按需同步 `AGENTS.md`

## Safety & guardrails

- 不自动选择 `tasks.md`。
- 不支持旧 `_Requirements:` 协议，不读取 `requirements.md`，也不迁移旧 spec。
- 不更改任务编号、标题、描述、AC 或设计内容。
- 不用用户交互替代阅读 `tasks.md`、`design.md` 和引用内容。
- 执行期间不引入需求或设计变更；需要变更时作为 blocker 报告。
- 验证成功后不请求继续许可。
- 失败任务保持 `[ ]`，并至少提供立即修复、延后继续、中止运行三个选项。
- IF 失败产物阻塞后续任务，THEN 在 blocker 的 `risk` 中说明级联影响。
