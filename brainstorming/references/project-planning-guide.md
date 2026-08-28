# Project Planning Guide

本文件补充 `SKILL.md` 的项目发现、能力拆分和滚动规划方法。
流程关卡、写入条件和输出边界以 `SKILL.md` 为准。

## Routing signals

| Signal | Route |
|--------|-------|
| 多个可独立产生价值的结果，MVP 或优先级未知 | `brainstorming` |
| 继续项目规划、调整 MVP/优先级、更新 backlog 或对账进度 | `brainstorming` |
| 一个可独立设计、验收和交付的结果 | `clarifying` |
| 已有批准设计，需要评审 | `design-review` |
| 已有批准设计，需要拆实施任务 | `spec-plan` |
| 继续执行已有 `tasks.md` | `spec-exec` |

“开发音乐 APP”是项目级请求；“为音乐 APP 增加本地歌单管理”通常是单个交付目标。
跨越多个代码模块不等于项目级；是否存在多个独立结果才是核心判断。
多个结果即使共享身份、数据或平台前置，仍属于项目级请求；强依赖影响拆分和顺序，
不把整个项目降级为一个 `clarifying` 目标。

## Iteration modes

| Mode | Typical user language | Resume point |
|------|-----------------------|--------------|
| `initialize` | “规划一个音乐 APP”“先拆项目需求” | STEP 2 |
| `resume` | “继续上次的项目规划”“接着处理未决问题” | 第一个阻断状态推进的 Open Question |
| `replan` | “新增这个需求”“调整 MVP”“规划下一阶段” | 最早受影响的 STEP 3-6 |
| `reconcile` | “更新项目进度”“处理 project handoff” | STEP 9 |

后续调用先读取本地项目文档，再解释本轮变化。不要因为会话是新的就重做首次 discovery，
也不要因为用户说“继续”就把规划续接误判成代码实施；动作对象是项目计划时才使用本技能。

## Greenfield discovery

没有现有代码时，优先确认：

- 目标用户及其当前替代方案。
- 最重要的使用场景和用户结果。
- 平台、时间、预算、团队、合规和运营约束。
- 首个版本要验证的最大不确定性。
- 明确不要进入 MVP 的能力。

把无法核实的内容标记为 Assumption，不要伪装成项目事实。

## Capability mapping

从用户旅程或业务结果组织能力：

```text
发现内容 -> 评估内容 -> 完成核心操作 -> 管理结果
```

不要按以下技术层组织项目 backlog：

```text
前端 -> 后端 -> 数据库 -> 部署
```

技术层可以出现在后续设计和 tasks 中，但它们本身通常不能独立产生用户价值。

## Work item decomposition

一个合格 work item 应满足：

- 有单一、清楚的用户或业务 Outcome。
- 可以独立进入 `clarifying` 并形成一个 `design.md`。
- 可以单独验收和交付，或明确说明不可避免的依赖。
- 延期不会让其他无依赖目标失去全部价值。
- 不包含 API、表、类、文件或测试任务等实现拆分。

常见拆分方法：

- 用户旅程：按可完成的用户阶段拆分。
- 业务能力：按独立业务结果拆分。
- 风险优先：先交付用于验证最大假设的纵向切片。
- 发布边界：按能够独立上线和回滚的结果拆分。

## MVP and priority

`Now` 只保留验证核心价值所必需的能力。判断某项是否属于 MVP 时询问：

1. 删除它后，核心用户还能完成最重要的结果吗？能则移到 `Next` 或 `Later`。
2. 它是在验证核心价值，还是只提高完整度、便利性或运营效率？后者通常不属于 `Now`。
3. 它是否解除其他 `Now` 项的真实依赖？不是则不因“以后会需要”提前纳入。

优先级变化属于用户决定。Agent 可以提供成本、风险和依赖证据，但不得自行替用户排序。

## Dependency types

只记录会影响顺序或可交付性的依赖：

- Business：业务规则或能力前置。
- External：供应商、合规、合同、账号或外部数据前置。
- Validation：需要先验证某项假设。
- Technical prerequisite：确实阻止独立设计或交付的技术前置。

不要把团队习惯或技术层顺序伪装成必然依赖。发现循环依赖时重新检查 work item 边界；
无法消除时记录阻塞和需要的项目决策。

## Rolling planning

只让近期目标进入 `ready-for-clarifying`。其他目标保留 Outcome、边界、Priority、Dependencies
和 Open Questions 即可，不提前决定技术方案或正式 AC。

项目事实变化时：

1. 识别受影响的 Goals、MVP、capability 和 work item。
2. 只重新讨论会改变用户决定的部分。
3. 保留稳定 `WI-*`，除非目标语义已经实质改变。
4. 重算受影响的依赖和当前焦点。
5. 呈现差异并等待确认后写回。

## Persistence and recovery

项目规划使用两个状态：

- `draft`：保存已确认内容和未决问题，用于暂停与恢复；不得新交付 work item。
- `ready`：项目方向、MVP 和近期拆分已经闭合，可以交付 `ready-for-clarifying` 目标。

恢复时先读 `project.md` 的 Planning Status，再定位 backlog：

- 有 `backlog.md`：work item、状态、证据和 Current Focus 以该文件为准。
- 无 `backlog.md`：`project.md` 内嵌的唯一 `## Backlog` 承担相同职责。
- 两处都包含 work item：停止并要求消除重复来源，不自行合并。

draft 写入文件前仍需用户确认。只恢复已经确认的内容；未确认方向写入 Open Questions，
不得补成 Goals、MVP 或 work item。

## Status reconciliation

只有 `brainstorming` 修改 backlog 状态。下游技能返回 handoff 或留下真实产物，
本技能按以下证据提出状态变化：

| Target status | Required evidence |
|---------------|-------------------|
| `ready-for-clarifying` | STEP 6 的全部 readiness 关卡通过 |
| `design-ready` | `design.md` 存在且 Project Traceability 匹配；已有 review 时必须 `ready / Go` |
| `planned` | 同目录 `tasks.md` 存在，引用可解析，且可选 review 仍为 `ready / Go` |
| `executing` | `tasks.md` 至少一个必需任务已完成，且仍有必需任务未完成 |
| `done` | 所有必需任务和 Checkpoint 均有完成标记 |
| `blocked` | review 非 `ready / Go`、spec 契约冲突或下游验证失败等明确反向证据 |

任何状态变化都同时写入 `Status evidence`。证据支持更晚状态时可以跨过没有持久证据的
中间阶段，但不得跳过目标状态自身的证据检查。范围、优先级或依赖变化不属于状态对账，
应回到主流程重新确认。

只有证据缺失时才保持原状态；证据明确否定当前状态时，提出 `blocked` 和对应
`Status evidence`，等待用户确认后写回。阻塞解除后重新按目标状态证据检查，不凭口头说明恢复。

## Project-level blind spots

按领域选择最相关的检查项，不机械遍历全部清单：

- 用户获取、激活、留存和退出路径。
- 管理员、运营、客服或内容治理能力。
- 隐私、合规、未成年人保护和内容权利。
- 外部平台政策、配额、合同和供应商依赖。
- 数据来源、生命周期、迁移、导入和导出。
- 上线、观测、支持、回滚和停止投入的条件。
- 跨目标共享的身份、权限、计费或核心数据能力。

这些检查用于发现项目级目标和风险，不用于提前设计实现机制。
