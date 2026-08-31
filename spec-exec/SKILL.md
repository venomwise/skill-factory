---
name: spec-exec
description: >
  按照 specs/<topic>/tasks.md、design.md 和项目仓库指导实现任务，并通过更新 checkbox
  跟踪进度和验证质量。
  当用户提到执行 spec、实施任务计划、继续执行 spec、恢复执行、检查 spec 进度时，
  必须使用此技能。
---

# Spec Exec Skill

## When to use

- 实现 `tasks.md` 中 `## Tasks` 下列出的任务
- 依据 `design.md` 的 Acceptance Criteria 验证行为
- 将已加载的仓库指导转化为当前任务可验证的 Implementation Contract
- 通过更新任务 checkbox 跟踪进度

## When not to use

- 只需要编辑任务文本而不执行工作
- Spec 没有 `design.md`、`tasks.md` 或 `## Tasks`

## Inputs / Outputs

**Inputs:**

- 执行计划：`specs/<topic>/tasks.md`
- 已批准的设计与验收标准：`specs/<topic>/design.md`
- 适用的仓库指导：上下文中已加载或仓库内可读取的 `AGENTS.md`、`CLAUDE.md`、
  README 和模块级约束

**Outputs:**

- 完成的实现和验证证据
- 更新后的 `tasks.md`：普通任务使用 `- [✅]`，可选任务使用 `- [✅]*`
- 最终响应中的命令结果、Implementation Contract 检查结果和未验证项

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
- checkbox 状态只由执行流程的主代理写回；工作包模式下的子代理不修改 `tasks.md`。

## Execution stance

将 `tasks.md` 视为执行计划；已批准的行为和设计方案以 `design.md` 为准，实现约束以适用的
仓库指导为准。两类来源互不替代；IF 它们冲突，THEN
在修改代码前按 *Blocker escalation* 报告。

IF 当前任务、引用和实现约束清晰且可执行，THEN 自主推进，不请求常规确认。

Acceptance Criteria 只定义可观察结果；验证命令来自任务、Checkpoint、`Testing` 章节和项目工具链。不要把缺失的行为或 AC 当作实现偏好自行补齐。

## Implementation Contract

仓库指导即使已经出现在上下文中，也必须在实施前转化为当前运行可检查的
Implementation Contract。上下文存在只证明规则可见，不证明代码已经合规。

只提取本次改动适用的约束，不复制整份指导文档：

- 指令来源及作用域，包括根目录和改动文件路径下更具体的指导。
- 执行基线，包括开始时的 git 状态、目标 diff 范围和需要保留的既有用户改动。
- architecture、依赖方向、公共接口和数据边界。
- 命名、注释、格式、错误处理和框架使用规则。
- 构建、测试、lint、静态分析和生成代码约定。
- 无法自动检查、必须通过 diff 人工核对的规则。

开始修改代码前执行以下关卡：

1. 是否已定位所有改动路径适用的仓库指导？IF 否，THEN 检查仓库根目录、父目录和模块目录后重试。
2. 是否已记录执行基线并区分既有用户改动？IF 否，THEN 检查 git 状态和目标 diff 后重试。
3. 每条与当前任务相关的约束是否有命令或明确的 diff 检查方法？IF 否，THEN 先补充检查方法。
4. 无法自动化的规则是否已进入人工检查清单？IF 否，THEN 不得开始实施。
5. 仓库指导与 `design.md` 是否存在冲突？IF 是，THEN 按 `guidance_conflict` blocker 停止。

每次进入新 Phase、发生上下文压缩或改动跨越新的模块边界时，重新核对 Contract。
简单且局部的任务允许使用短清单，不为凑流程扩写无关规则。

## Strict contract gate

开始执行前必须验证：

1. `design.md` 包含 `## Acceptance Criteria`，AC ID 唯一且符合 `AC-<domain>-<behavior>`。
2. `tasks.md` 不包含 `_Requirements:`；命中即视为不受支持的旧格式。
3. 每个未完成的可执行 sub-task 和独立 Checkpoint 至少包含 `_Acceptance:` 或 `_Design:`；Phase 容器除外。
4. 所有 `_Acceptance:` ID 和 `_Design:` 章节路径都能在 `design.md` 中解析。
5. 任务内容与引用的 AC、设计章节没有冲突。
6. Implementation Contract 已通过关卡，且每条适用约束有验证方法。

IF 任一检查失败，THEN 在修改代码前按 *Blocker escalation* 报告并停止。不要迁移旧 spec，也不要补写 AC。

## Checkpoints

Checkpoint 是验证任务，不是用户批准门。执行其中的命令，使用 `_Acceptance:` 定位必须验证的 AC，
并同时核对适用的 Implementation Contract，以可观察证据判断通过条件。

IF 验证通过，THEN 标记完成并继续。仅当验证失败、必需资源不可用或 spec 内部不一致时停止。

## Evidence-based validation

在标记任务或 Checkpoint 完成前，必须：

- 执行任务或 Checkpoint 中的显式验证命令。
- 运行项目提供的相关构建、测试、类型检查、lint、静态分析或冒烟测试。
- 在 git 仓库中运行 `git diff --check`，并检查未跟踪文件没有遗漏。
- 检查实际 diff，将结果与引用的 AC、设计章节和 Implementation Contract 逐条对比。
- 检查新增或修改的测试是否对关键行为包含有效断言，而不只证明代码可编译或调用不抛异常。
- IF 某条命令无法运行，THEN 说明原因并执行最强的替代检查；IF 必需约束仍无法验证，THEN 按验证失败停止。

构建或测试通过不能替代无法自动化的架构、注释和代码规范检查。不要仅凭任务文本、
文件存在、编译成功或未经验证的假设标记完成。

## Completion quality gates

**Task gate：** 标记单个任务前检查：

1. 引用的行为和设计是否有直接证据？IF 否，THEN 返回当前任务补充实现或测试。
2. Implementation Contract 的适用规则是否全部通过？IF 否，THEN 修复后重新验证。
3. 改动是否引入任务范围外行为或无关重构？IF 是，THEN 移除无关改动或按 blocker 报告范围变化。
4. 验证命令和关键结果是否已记录在本轮回复中？IF 否，THEN 先记录证据再更新 checkbox。

**Phase gate：** 标记 Phase 前审查该阶段累计 diff，检查 architecture、依赖方向、错误处理、
重复实现、测试质量和仓库规范。IF 发现问题，THEN 保持 Phase 未完成，并重新打开对应任务。

**Final gate：** 所有必需任务完成后，相对执行基线重新检查本次完整 diff，运行适用的最高层级验证命令，
确认各 Phase 的局部正确性组合后仍满足 AC 和 Implementation Contract。IF 失败，THEN
重新打开对应任务或 Checkpoint，不得报告交付完成。

## Blocker escalation

仅在真正阻塞时停下，例如：

- Spec 使用旧 `_Requirements:` 格式或缺少 Acceptance Criteria。
- 引用的 AC 或设计章节不存在、含糊或互相冲突。
- 仓库指导与已批准设计冲突，无法同时满足。
- 任务无法从 `tasks.md` 和引用的 `design.md` 内容确定。
- 验证失败且无法在当前任务范围内安全修复。
- 完成任务需要改变已批准 AC 或设计。
- 需要破坏性操作，或缺少凭证、服务、文件、环境依赖。

使用以下结构报告：

```yaml
blocker:
  task: "<N.M> <title>"
  type: >-
    <legacy_spec | conflict | guidance_conflict | underspecified |
    validation_failure | scope_change | destructive_op | missing_dependency>
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

## Work package dispatch mode

用户选择工作包模式后，主代理不直接实现任务，只负责分包、派发、验证和状态记账；
实现由子代理完成。目标是上下文隔离和失败隔离，不是必然提速。分派指令的完整结构、
报告格式和禁止事项见 [dispatch contract](references/dispatch-contract.md)。

**分包规则：**

1. 一个工作包是文件集互不相交、相互无顺序依赖的一组相邻任务；包内任务由同一个子代理
   按原顺序连续实现，保留包内的连续改动和共享测试夹具。
2. 数据迁移、事务边界、公共契约或共享状态相关的任务不进工作包，留在主线程顺序执行。
3. Checkpoint 由主线程执行，不放进工作包。
4. 每个包必须映射到可验证的 AC 子集；无法映射的任务留在主线程。
5. 分包结果（包编号、任务清单、文件集、引用的 AC）展示给用户，确认后才开始派发。
6. IF 分包后只剩一个可用包，THEN 建议用户回到主线程模式，不为单包派发。

**派发与验证：**

1. 子代理按分派契约接收任务文本、追踪元数据、适用的 Implementation Contract 条目、
   执行基线和报告格式；子代理不读取也不修改 `tasks.md` 的 checkbox。
2. 包之间默认串行派发；仅当两个包的实现与验证都确认互不影响时，才允许同时派发。
3. 收到报告后，主代理重读实际 diff 验证，不以子代理声明作为完成证据：对照包引用的
   AC 和 Implementation Contract 逐条核对，执行 *Evidence-based validation* 中适用的
   命令和 Task gate。
4. 验证通过后由主代理立即写回该包任务的 checkbox，再派发或验证下一个包。
5. 验证失败时，把具体缺口反馈给子代理重做一次；再次失败则由主代理接管该包，
   或按 *Blocker escalation* 报告。
6. 子代理报告包外改动、设计缺口或需要用户决策的事项时，主代理按 *Blocker escalation*
   处理，不把报告直接当作完成。

派发模式下 Phase gate、Final gate 和 `AGENTS.md` 同步仍由主代理执行；
*Completion quality gates* 不因派发模式放宽。

## Workflow

1. **定位 tasks.md。**
   - IF 用户提供了路径，THEN 使用它。
   - IF 未提供，THEN 精确回复：`请指定 tasks.md 文件路径（例如 specs/<topic>/tasks.md）。`然后结束流程。不要搜索、列举或推断 spec。
   - IF 路径不存在，THEN 推荐先运行 `spec-plan`。
2. **读取完整 spec。** 打开 `tasks.md` 和同目录 `design.md`，定位 `## Tasks`、`## Acceptance Criteria`、`## Testing` 及任务引用的设计章节。
3. **建立 Implementation Contract。** 记录执行基线和既有用户改动；从已加载和仓库可读取的指导中提取本次改动适用的规则，按对应关卡确认验证方法。
4. **执行严格关卡。** 按 *Strict contract gate* 检查；失败时在修改代码前停止。
5. **扫描并恢复进度。** 计数 `[✅]` 与 `[ ]`，从首个未完成任务继续；全部完成时仍执行 Final gate，再报告状态。
6. **选择执行模式。**
   - IF 存在 `[ ]*`，THEN 仅询问一次：`当前任务列表包含可选任务。A. 保留可选任务（MVP） B. 执行所有任务`。
   - 执行模式默认**主线程顺序执行**。IF 运行环境能启动实现子代理且剩余未完成任务可拆出
     至少两个满足分包条件的工作包，THEN 允许一次性询问：`A. 主线程顺序执行（默认）
     B. 工作包派发模式`。无子代理能力或只有一个可用包时不提供该选项，不阻塞流程。
     选择工作包模式后按 *Work package dispatch mode* 执行。
7. **Triage。** 按顺序处理任务：
   - 跳过 `[✅]`。
   - MVP 模式跳过 `[ ]*`；可选 Phase 被跳过时连同其子任务一起跳过。
   - 包含 sub-task 的 Phase 是容器，不直接实施；按顺序处理其子任务，并在步骤 10 汇总状态。
   - 通过 `Checkpoint`、`Verify`、`检查点` 识别验证任务。
   - 工作包模式下，进入工作包的任务按 *Work package dispatch mode* 派发和验证，
     其余任务仍走步骤 8。
   - 其他任务进入步骤 8。
8. **Implement & Validate。**
   - 阅读任务描述和 `_Acceptance:` / `_Design:` 元数据。
   - 修改前阅读相关文件与模块，理解现有行为。
   - 按 Implementation Contract 实施，并执行 *Evidence-based validation* 和 Task gate。
   - IF 失败，THEN 按 *Blocker escalation* 处理。
   - 记录命令、关键结果和人工 Contract 检查结论，使验证证据在任务间可恢复。
9. **立即写回进度。** 每完成一个任务，先把它改为 `[✅]` 或 `[✅]*`，再做任何其他动作。一个任务一次写入；失败、跳过或中断的任务不得标记完成。工作包模式下，主代理在一个包的验证通过后集中写回该包任务，再处理下一个包。
10. **更新 Phase。** 一个 Phase 的适用子任务全部完成后执行 Phase gate；通过后标记 Phase。MVP 模式下未执行的可选子任务不阻塞 Phase 完成。无子任务的 Checkpoint 只在验证通过后标记。
11. **执行最终质量关卡。** 所有必需任务完成后执行 Final gate。IF 根目录存在 `AGENTS.md`，THEN 只在本次实现改变项目结构、命令、验证流程或仓库约定时更新相关内容，不重写无关章节。
12. **推荐独立评审。** IF 改动涉及跨层调用、共享或公共契约、数据迁移、并发、安全、权限，
    或多个相互耦合模块，THEN 推荐运行 `code-review`；不要自动调用下游 skill。
    简单局部改动通过 Final gate 后可直接结束。

## Verification

- [ ] 执行前通过严格 contract gate
- [ ] 已从仓库指导建立任务相关的 Implementation Contract
- [ ] 每个可执行 sub-task 和独立 Checkpoint 依据可解析的 `_Acceptance:` 或 `_Design:` 执行
- [ ] 标记完成前同时有行为证据和仓库规范证据
- [ ] Checkpoint 实际执行且按引用的 AC 和 Implementation Contract 验证
- [ ] 每个 Phase 已通过累计 diff 的 Phase gate
- [ ] 全部必需任务已通过完整 diff 的 Final gate
- [ ] 每完成一个任务立即写回 checkbox
- [ ] MVP 模式下可选任务保持未勾选
- [ ] Phase 仅在其适用子任务完成后标记
- [ ] 必需任务完成后按需同步 `AGENTS.md`
- [ ] 工作包模式经过用户一次性确认，分包方案展示并确认后才派发
- [ ] 工作包之间文件集不相交、无顺序依赖，共享边界任务和 Checkpoint 在主线程执行
- [ ] 子代理报告均通过主代理重读实际 diff 的独立验证后才写回 checkbox

## Safety & guardrails

- 不自动选择 `tasks.md`。
- 不支持旧 `_Requirements:` 协议，不读取 `requirements.md`，也不迁移旧 spec。
- 不更改任务编号、标题、描述、AC 或设计内容。
- 不用用户交互替代阅读 `tasks.md`、`design.md` 和引用内容。
- 不把仓库指导已经出现在上下文中当作合规证据，也不为本次任务复制整份指导文档。
- 不覆盖、回退或混入执行基线中已有的用户改动。
- 执行期间不引入需求或设计变更；需要变更时作为 blocker 报告。
- 验证成功后不请求继续许可。
- 失败任务保持 `[ ]`，并至少提供立即修复、延后继续、中止运行三个选项。
- IF 失败产物阻塞后续任务，THEN 在 blocker 的 `risk` 中说明级联影响。
- 工作包模式必须经用户确认分包方案后才派发；无子代理能力时保持主线程模式，不阻塞。
- `tasks.md` 的 checkbox 只由主代理更新；子代理报告不作为完成证据，验证失败最多重试一次，
  再失败由主代理接管或按 blocker 报告。
- 输出语气直接、自然，不使用模板化客服话术。

## References

- [Dispatch contract](references/dispatch-contract.md)
