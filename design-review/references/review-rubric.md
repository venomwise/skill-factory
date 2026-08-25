# Design Review Rubric

本文件定义语义评审维度、严重性和 Verdict。章节结构、Decision、Data Model、Interfaces、
Acceptance Criteria 和格式以
[canonical design document template](../../brainstorming/assets/design-doc-template.md) 为唯一权威。
评审模式、finding 和 readiness 以
[review lifecycle](review-lifecycle.md) 为唯一权威。

## Evaluation discipline

每条 finding 必须记录具体位置、问题、证据、建议、严重性和 Origin。
一个根因影响多个维度时只创建一个 finding，并在 Dimension Summary 中多处引用。

Full Review 在完成所有维度和根因合并后一次性冻结 finding 集合。
Closure Review 不重新自由发散，只检查旧 finding、changed sections 和传播边界。

## D1: Completeness

- canonical template 的必需章节存在且有实质内容。
- 每个核心 Goal 同时映射到 Proposed Solution 和至少一条 AC。
- Architecture、Components、Data Flow、Error Handling 和 Testing 形成闭环。
- Component 有职责、输入和产出，不只是名称。
- Data Flow 覆盖主路径；非平凡设计覆盖关键失败路径。
- 持久且非显然的选择有 Decision；直接修复没有被包装成 Decision。
- 用户明确驳回的持久边界有 `Rejected concern` 和 `Revisit when`。

缺少整个 Acceptance Criteria 或核心 Goal 没有 AC 是 Blocker。
缺少承重的 Components、Data Flow 或 Error Handling 通常是 Blocker。

## D2: Usability

- `spec-plan` 可以在不发明行为、数据或契约的前提下生成任务。
- 输入、输出、数据形态、状态和错误语义无歧义。
- 术语有定义并跨章节一致。
- Goal 的成功条件可观察或可衡量。
- 每条 AC 只表达一个可验证结果。
- 行为相关 Open Question 已解决。

“实现时再决定”的关键选择是一条 finding，不是可接受占位符。

## D3: Conformance

- 章节和按需子章节符合 canonical template。
- spec 位于 `specs/<topic>/`，topic 使用 kebab-case。
- Decision 使用稳定 `DR-<semantic-topic>`，同一主题没有重复或 `Revised` 堆叠。
- AC ID 唯一，符合 `AC-<domain>-<behavior>`。
- AC 使用 `WHEN` 或 `IF`、`THEN` 和 `SHALL`。
- 普通 Markdown 行不超过 120 字符；表格保持短列。
- Open Questions 不充当未完成设计的垃圾场。

## D4: Project Fit

- 设计符合项目技术栈、依赖版本、模块边界和本地约定。
- 复用现有接缝，没有重复已有能力。
- 明确说明被改变的既有行为。
- 命名、数据约定和文件布局与项目一致。
- 数据库和代码声明与已核实事实一致。

## D5: Blind Spots

只选择与当前设计相关的检查项。

**User-facing:**

- 空、加载和错误状态。
- 权限、访问控制、撤销和回滚。
- 国际化、无障碍和客户端兼容。

**Data:**

- 迁移、回填和回滚保护。
- 一致性、并发、幂等和冲突处理。
- 保留、清理、隐私、备份和恢复。

**Integration:**

- 配额、认证和凭证轮换。
- 超时、重试、重放、去重和降级。
- 版本、向后兼容、监控和告警。

Blind Spot 必须说明为什么适用于当前 Goals 或变更边界。纯推测不进入 Findings。

## D6: Over-Engineering

- 每个组件和机制都能追溯到 Goal。
- 不设计 Non-Goals 中明确排除的内容。
- 不为单一用途引入没有第二个调用方的通用框架。
- 不增加无需求支撑的配置、扩展点或防御逻辑。
- 更简单方案满足全部 Goals 时，指出具体替代方案。

判断依据是设计自己的 Goals 和 Non-Goals，不是 reviewer 的技术偏好。

## D7: Optimization

- 是否可以复用更多现有组件。
- 是否可以降低耦合或减少新增面。
- 是否存在项目中已经验证的更简单模式。
- 优化不应改变用户已确认的行为和范围。

D7 通常是 Minor；与真实风险重叠时按风险影响定级。

## Data Model checks

设计包含数据变化时，检查：

- Change Summary 枚举全部 ADD/MODIFY/REMOVE 对象。
- 新对象给出完整目标模型；已有对象只描述变化字段。
- Persistence、Domain/Serialized、Migration 和 Concurrency 内容归属正确。
- 约束、索引、事实源和所有权清晰。
- 可执行 SQL 位于独立文件，design 引用路径且不维护第二份 DDL/DML。
- 正向、回填、回滚和数据保护语义闭合。
- Data Model 没有混入接口请求、资源流程或无关实现叙述。

## Interfaces checks

设计包含契约变化时，检查：

- Change Summary 枚举全部 ADD/MODIFY/REMOVE API、消息和协议。
- 每个变更有稳定 Contract ID。
- HTTP 定义认证、请求、响应、错误、幂等和兼容性。
- Message 定义 channel、参与者、payload、投递、顺序、重试、去重和版本。
- Protocol 定义方向、command、framing、payload、ACK、超时、重试和版本。
- `MODIFY` 说明原契约到新契约变化；`REMOVE` 说明迁移和失败行为。
- 机器可读契约存在时，design 引用它而不复制第二份完整 schema。
- Data Flow 的跨边界调用能映射到 Contract ID。

缺少实现关键的请求、响应或消息语义通常是 Major；导致调用方无法实现时可定为 Blocker。

## Severity

| Severity | Meaning | Examples |
|----------|---------|----------|
| Blocker | 阻止实现或与项目事实冲突 | 核心 Goal 无设计/AC，关键契约自相矛盾 |
| Major | 关键路径的真实缺口或风险 | 含糊接口、迁移不安全、执行链路缺承载者 |
| Minor | 不阻断的改进或打磨 | 命名、可读性、非关键复用机会 |

不要通过夸大严重性让报告显得全面。

## Verdict at review

- **Reject / No-Go**：至少一条 Blocker。
- **Revise / Conditional**：无 Blocker，至少一条 Major。
- **Pass**：无 Blocker/Major；Minor 仍需用户修复或延期。

当前是否可以进入 `spec-plan` 不只看 Verdict，还必须读取 Current Readiness。
只有 `Overall: ready` 且 `spec-plan readiness: Go` 才能进入规划。

## Full Review completion gate

1. 7 个维度是否全部检查？If no，继续评审。
2. Goal -> Solution -> AC 覆盖是否完成？If no，补完。
3. Data Model 和 Interfaces 是否按适用性检查？If no，补完。
4. 候选 finding 是否按根因合并并校准严重性？If no，先合并。
5. 最终集合是否经过交叉检查？If no，不得写报告。

## Closure Review admission gate

每条 Closure 新 finding 必须回答：

1. 它属于哪一个允许 Origin？
2. 哪个 changed section、已解锁依赖或新项目事实使它现在成立？
3. 它为什么不是已有 finding 的同一根因？
4. 是否已有有效 Decision 明确拒绝该 concern？

任一问题无法回答时，不得创建新的 Closure finding。
