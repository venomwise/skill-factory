# Review Lifecycle Contract

`design-review` 是本协议的权威来源。`design-refine` 必须直接引用本文件，不得复制状态定义。

## Ownership

`review.md` 使用单文件双视图：

- `Review Snapshot` 由 `design-review` 写入，表示本轮评审发生时的事实，不由 refine 改写。
- `Current Readiness` 由 `design-review` 初始化，由 `design-refine` 更新当前状态和下一步。
- `Closure` 由下一轮 `design-review` 更新，只保留紧凑的跨轮闭合证据。
- `Accepted / Deferred Risks` 由 `design-refine` 更新。
- `Recommended Next Step` 从 Current Readiness 派生，可由两个 skill 更新。

`design.md` 只表达当前有效设计。评审处理历史留在 `review.md`，不通过重复 Decision 堆入设计。

## Review Modes

### Legacy review migration

历史 `review.md` 缺少 `## Current Readiness`，或只使用 severity-coded finding ID 时，
`design-review` 执行一次 Full Review migration。迁移不猜测旧 finding 状态；它按当前设计
重新生成 `F-###`、Review Snapshot 和 Current Readiness，并在 Closure 保留一行 legacy 说明。
`design-refine` 不直接处理 legacy review。

### Full Review

以下任一条件成立时执行：

- 不存在历史 `review.md`。
- 用户明确要求完整评审。
- Goals 或 Non-Goals 发生变化。
- 核心 Architecture、数据所有权或公共契约发生变化。
- 新增安全、权限、外部集成或迁移边界。
- 与设计相关的代码、数据库或外部事实在上一轮后发生变化。

Full Review 完整执行 7 个维度，合并同根问题后冻结本轮 findings。

### Closure Review

已有 review、Current Readiness 为 `ready-for-closure`，且 refine 未触发 Full Review 条件时执行。
`not-started`、`in-progress` 或 `blocked` 表示 refine 尚未完成，不得启动 Closure Review。
`ready / Go` 表示已经放行，除非用户明确要求，否则不重复评审。

Closure Review 的语义检查范围限于：

- 上一轮未闭合或重新打开的 finding。
- `Current Readiness` 记录的 changed sections。
- 这些变化对 Goals、Non-Goals、Components、Data Flow、Error Handling、
  Acceptance Criteria 和 Testing 的传播。

文档结构、ID、引用完整性和格式等确定性检查仍然全量执行。

## Finding Identity

finding 使用全局单调且与严重性无关的 ID：`F-001`、`F-002`。
同一根因继续存在、标题变化或严重性变化时沿用原 ID。

每条 finding 独立记录：

- `Severity`: Blocker、Major 或 Minor。
- `Introduced in`: 首次发现轮次。
- `Origin`: finding 产生原因。
- `Status`: 当前处理状态。

R1 Full Review 的 finding 使用 `initial-review`。R2 及后续评审的新 finding 使用：

| Origin | 含义 |
|--------|------|
| `refine-regression` | refine 修改直接引入矛盾或缺口 |
| `dependency-unlocked` | 上一轮阻断解除后才具备评估条件 |
| `baseline-miss` | 问题原先存在，但 Full Review 漏检 |
| `context-change` | 项目代码、数据库或外部约束发生变化 |

R2 及后续评审的新 finding 必须说明 Origin 的因果证据。Closure Review 无法分类时不得写入 Findings。
后续 Full Review 中，范围或项目事实变化使用 `context-change`；原先已存在但漏检的问题使用
`baseline-miss`。

## Finding Status

| Status | 含义 | Terminal |
|--------|------|----------|
| `pending` | 尚未开始 | No |
| `in-discussion` | 正在确认决策或边界 | No |
| `blocked` | 缺少用户输入或外部前提 | No |
| `reopened` | 已处理问题因证据变化再次成立 | No |
| `resolved` | 设计或文档已更新并通过一致性检查 | Yes |
| `rejected` | 用户明确驳回，且持久边界已有证据 | Yes |
| `accepted-risk` | 用户明确接受风险，且风险已有证据 | Yes |
| `deferred` | Minor 已明确延期 | Yes |

`rejected` 必须指向现有或新增的持久边界 Decision。若现有 Decision 已充分覆盖，直接复用，
不得追加同义记录。

## Current Readiness

`Current Readiness` 使用以下 Overall 状态：

| Overall | 含义 |
|---------|------|
| `not-started` | 当前 findings 尚未开始处理 |
| `in-progress` | 至少一条 finding 正在处理 |
| `blocked` | 当前处理缺少必要输入 |
| `ready-for-closure` | Blocker/Major 已处理，必须运行 Closure Review |
| `ready` | 当前设计可以进入 `spec-plan` |

`spec-plan readiness` 使用 `Go`、`Conditional`、`No-Go`：

- `Go`: Overall 为 `ready`，所有 finding 已终态，且没有行为相关 Open Question。
- `Conditional`: 仍有 Major、Minor 待处置，或已完成 refine 但等待 Closure Review。
- `No-Go`: 存在 Blocker、状态不完整或 Closure Review 发现新的阻断问题。

Pass 但存在 Minor 时，用户必须选择修复或延期。所有 Minor 终态后才设置 `ready / Go`。
经过 Blocker/Major refine 的设计必须先进入 `ready-for-closure`，不得直接设置 `Go`。

## State Transitions

1. `design-review` 创建 Review Snapshot，并把当前 finding 初始化为 `pending`。
2. `design-refine` 处理 finding，更新 Status、Resolution ref 和 Changed sections。
3. 全部 Blocker/Major 终态后，Current Readiness 进入 `ready-for-closure`。
4. `design-review` 执行 Closure Review，将历史结果写入 Closure。
5. 无未解决 Blocker/Major，且 Minor 已处理时，Current Readiness 进入 `ready / Go`。

Refine run 使用 `R<n>-refine-<n>`，不复用 finding 的 `F` 前缀。

## Closure Rules

Closure 每条只保留：finding ID、标题、closed/continued、Resolution ref 和一句结果。
不要复制上一轮完整 Issue 和 Recommendation。

下一轮遇到相同或类似问题时，先检查已有 Decision：

- Decision 仍有效且 `Revisit when` 未触发：不得重复创建 finding。
- 设计正文与 Decision 冲突：沿用原 ID，状态设为 `reopened`。
- `Revisit when` 已触发：允许以 `context-change` 创建或重新打开 finding。
- 风险与原边界实质不同：允许新 finding，但必须说明差异。
