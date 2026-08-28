# <Project> Project Plan

<!--
本地项目方向以本文件为准。替换指导文字，删除不适用的可选内容。
普通 Markdown 行不超过 120 个 Unicode 字符。
不要在本文件中写技术设计、正式 Acceptance Criteria 或实施任务。
中文正文使用自然、具体的说明，写清对象、问题、结果和条件，不用口号代替项目内容。
-->

## Summary

用一段话说明目标用户当前遇到什么问题、项目准备改变什么，以及为什么现在需要处理。

## Planning Status

- **Status**: <draft/ready>。
- **Reason**: <当前状态依据；draft 说明主要未决问题，ready 说明通过的项目级关卡>。

## Target Users and Outcomes

| User / Role | Current problem | Desired outcome |
|-------------|-----------------|-----------------|
| <用户或角色> | <当前问题> | <期望结果> |

## Goals

- 每项写清用户或业务对象以及期望结果，不写“提升能力”一类无法验证的口号。

## Non-Goals

- 明确当前项目不解决的方向，防止范围蔓延。

## Success Metrics

| Metric | Baseline | Target | Evidence source |
|--------|----------|--------|-----------------|
| <指标> | <现状或 Unknown> | <可判断是否达成的目标> | <可核实的数据来源> |

## MVP Boundary

### Now

- 列出构成首个可验证 MVP 的能力，并引用对应 `WI-*`。

### Next

- 列出 MVP 得到验证后的优先扩展方向。

### Later

- 列出保留但暂不进入详细澄清的方向。

## Capability Map

| Capability | User outcome | MVP | Work items |
|------------|--------------|-----|------------|
| <能力> | <结果> | <Now/Next/Later> | <WI-...> |

## Constraints and Assumptions

### Confirmed Constraints

- 记录已经核实的业务、资源、合规、平台或时间约束。

### Assumptions to Validate

- 记录仍需验证的假设、验证方式和影响范围。

## Project Risks

| Risk | Impact | Mitigation / next action | Owner |
|------|--------|--------------------------|-------|
| <风险及触发条件> | <受影响对象和具体后果> | <缓解或下一步> | <责任角色> |

## Project Decision Record

只记录持久且非显然的项目范围、MVP、优先级或拆分选择。技术选择放入对应 `design.md`。

### PDR-<semantic-topic>

**Decision**: <谁在什么范围内采用什么选择>。

**Rationale**: <该选择如何对应项目目标和约束；被拒方向只写名称和拒绝理由>。

**Revisit when**: <哪些项目事实变化时需要重新评估>。

## Backlog

标准模式只链接 `[backlog.md](backlog.md)`，不要复制内容。
仅在降级模式下，把 `backlog-template.md` 的 work item 内容放在这里，并且不创建 `backlog.md`。

## Open Questions

- `draft`：记录已经浮现的未决问题、影响范围和下一项需要用户决定的内容。
- `ready`：不得保留会改变项目目标、MVP、拆分或当前焦点的未决问题。
