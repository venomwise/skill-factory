---
name: prd-review
description: Analyze PRD documents for logic gaps, boundary issues, and technical feasibility before development starts. Focus on business logic closure and architectural reasonability, not implementation details.
---

# PRD Review

Technical feasibility review for Product Requirement Documents. This skill evaluates PRD logic completeness, boundary coverage, and technical risks **before any design or development work begins**. It operates in the requirements clarification phase — identifying what needs to be clarified with the PM, not how to implement it.

Match the user's language for all output.

## When to use

- PM 提供了 PRD 文档，需要技术团队评审
- 需要识别需求中的逻辑漏洞、边界缺失、技术可行性问题
- 生成结构化的澄清问题清单供 PM 回答
- 在进入任何技术设计或开发前，确保需求本身逻辑闭环

## When not to use

- 已有明确技术设计方案，需拆解为开发任务（用 `spec-plan`）
- 代码已实现需要评审（用 `code-review` 或 `review`）
- 简单功能调整或 bug 修复（用 `clarification`）
- 需要编写或修改技术方案本身（用 `brainstorming`）

## Inputs

- PRD 文档路径（用户提供，例如 `/path/to/docs/user-auth-prd.md`）
- 可选：用户补充的 PM 回答或上下文说明

## Outputs

- `<prd-basename>-review.md` 写入 PRD 文档的同级目录（`<prd-basename>` = PRD 文件名去掉扩展名）
  - 例如：PRD 在 `/path/to/docs/payment-prd.md`，则输出 `/path/to/docs/payment-prd-review.md`
  - 同一目录下多个 PRD 各自对应独立评审文件，互不覆盖
- 报告包含：风险点清单、技术可行性分析、需 PM 澄清的 Todo List

## Workflow

### Step 1: 检查迭代状态

检查 PRD 同级目录是否存在对应的 `<prd-basename>-review.md`（与当前 PRD 同名、以 `-review.md` 结尾）：

- **首次分析：** 不存在 → 创建新的完整分析（见 Step 2-4）
- **迭代分析：** 存在 → Read 已有分析，进行增量更新

迭代时的处理：
1. 读取旧的 `<prd-basename>-review.md`，理解已标记的风险点和疑问
2. 对比 PRD 的变化（如果用户说明 PRD 已更新）或用户的补充信息（如 PM 的回答）
3. 更新分析：标记已解决的问题、新增发现的问题、删除不再适用的旧问题
4. 保持历史的可追溯性（建议添加"变更记录"章节记录每次迭代）

迭代报告的具体写法见 [templates/report-template.md](templates/report-template.md) 的"迭代分析"部分。

### Step 2: 探索技术上下文

**主动探索，不要凭空推测。** 目标是搞清楚事实，怎么查由你判断。

#### 2.1 本地项目信息

探索现有架构和技术栈，为可行性评估提供依据。自行选用合适的工具（文件搜索、代码检索、读取文档、符号跳转等），重点搞清楚：

- **技术栈与架构：** 用什么语言/框架/运行时？现有分层与关键模块是什么？
- **相关实现：** PRD 提到的功能在代码库中是否已有实现或近似实现？入口与调用链在哪？
- **项目约定：** 是否有 README / CLAUDE.md / AGENTS.md 等说明文件，记录了相关约束或设计决策？

当 PRD 涉及数据库相关需求时（如数据查询、状态变更、报表统计），通过 subagent 执行 `db-explorer` skill 来查询（只读探索，不修改数据），搞清楚相关表结构是否存在、字段类型是否支持需求的数据范围、索引是否能支撑查询性能要求。

#### 2.2 行业/开源最佳实践

对于本地项目尚未实现的新功能，通过 `web_search` 获取行业标准和常见陷阱。按领域的搜索主题与搜索原则见 [references/best-practices.md](references/best-practices.md)。

### Step 3: 结构化分析

**分析粒度原则：** 聚焦业务逻辑闭环和架构合理性，不纠结实现细节。

按以下三个维度进行结构化分析：

1. **逻辑漏洞与边界缺失** — 正常流程 + 异常分支、边界条件、状态机、多角色交互是否闭环
2. **技术可行性与架构风险** — 现有架构能否支持、性能/安全/合规要求是否可达、第三方依赖
3. **需 PM 澄清的疑问** — 量化指标、业务规则、交互逻辑、前置条件的缺失

完整的"应该质疑 vs 不应该质疑"清单和每个维度的检查点，见 [references/analysis-dimensions.md](references/analysis-dimensions.md)。**务必先读取该文件，确保质疑聚焦逻辑层面而非实现细节。**

### Step 4: 生成/更新报告

风险级别全文统一使用三档：

- `[P0-阻断]` — 高危，逻辑/架构存在阻断性缺陷，需 PM 或架构澄清后才能继续
- `[P1-需确认]` — 中危，需与 PM/团队确认，可在设计阶段处理但有风险
- `[P2-优化]` — 低危，优化建议，不阻塞推进

报告写入 `<prd 同级目录>/<prd-basename>-review.md`。完整的报告结构模板（首次分析 + 迭代分析）见 [templates/report-template.md](templates/report-template.md)。

### Step 5: 完成确认

- 如果仍有未解决的疑问（Todo List 中有未勾选项），告知用户："已生成评审报告，发现 X 个待澄清问题，请与 PM 确认后再次调用此 skill 进行迭代分析。"
- 如果所有疑问已解决（所有 Todo 已勾选，无阻断性风险），告知用户："需求澄清完成，所有逻辑闭环和技术可行性问题已确认，可以开始设计技术方案。"

## Verification

生成报告后自检：

- [ ] `<prd-basename>-review.md` 已写入 PRD 同级目录
- [ ] 每个风险点都引用了 PRD 的具体章节或当前技术现状（文件路径）
- [ ] Todo List 中的问题清晰、具体、可直接转发给 PM（无需二次编辑）
- [ ] 技术可行性分析基于真实代码库状态（已通过工具实地探索代码/数据库验证，而非凭空推测）
- [ ] 风险级别标注一致且有明确判断依据
- [ ] 没有质疑实现细节（如 UI 样式、命名规范、技术选型细节）
- [ ] 迭代分析时保留了历史记录的可追溯性

## Safety & guardrails

- **实事求是：** 每个风险点必须有 PRD 原文引用或技术现状依据，不能无中生有。
- **聚焦大方向：** 质疑业务逻辑闭环和架构合理性，不纠结实现细节。
- **客观专业：** 评审语气客观、严谨、具有建设性，不带主观情绪。
- **主动探索：** 不要等待用户提供上下文，主动选用合适工具（代码检索、读取文档、subagent 执行 db-explorer skill、web_search 等）获取信息。
- **明确边界：** 不代替 PM 做产品决策（如优先级、需求取舍），不代替开发做技术选型（如用什么框架、什么数据库）。
- **增量更新：** 迭代时基于已有分析进行增量更新，保留历史记录，不重复造轮子。

## Examples

四个典型场景（首次审查、PM 补充后迭代、全部解决、涉及数据库的需求）见 [references/examples.md](references/examples.md)。
