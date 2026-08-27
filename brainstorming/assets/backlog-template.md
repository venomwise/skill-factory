# <Project> Backlog

<!--
本文件是项目交付目标、顺序和状态的本地权威来源。
work item 描述用户或业务结果，不描述文件、函数、API、表结构或实现任务。
正式 Acceptance Criteria 和技术方案只进入对应 `specs/<topic>/design.md`。
只有 `brainstorming` 可以修改 Status、Spec、Status evidence 和 Current Focus。
-->

## Status Model

主状态：

```text
candidate -> ready-for-clarifying -> design-ready -> planned -> executing -> done
```

旁路状态：`blocked`、`deferred`、`cancelled`。

## Current Focus

- `<WI-...>` — <为什么现在推进以及下一步>。

## Work Items

| ID | Outcome | Priority | Status | Dependencies |
|----|---------|----------|--------|--------------|
| `WI-<semantic-name>` | <用户或业务结果> | <Now/Next/Later> | <status> | <WI-* 或 None> |

### WI-<semantic-name>: <Title>

- **Outcome**: <独立产生的用户或业务结果>。
- **Scope**: <高层能力边界>。
- **Non-Goals**: <本目标明确不包含的能力>。
- **Success signal**: <进入详细澄清前用于判断方向的项目级信号>。
- **Priority**: <Now/Next/Later>。
- **Status**: <受支持状态>。
- **Dependencies**: <WI-*、外部前置或 None>。
- **Readiness**: <已满足条件或当前阻塞>。
- **Spec**: `<真实 design.md 路径或 Not created>`。
- **Status evidence**: <支持当前状态的路径和可观察事实，或 None>。
- **External ref**: <真实 Jira/GitHub 标识或 None>。

#### Open Questions

- <留到本目标进入 `clarifying` 前后解决的问题；没有则写“无”>。

## Dependency Notes

- 说明无法仅通过 `Dependencies` 字段表达的业务、外部或交付顺序约束。

## Deferred / Cancelled

- 保留稳定 work item ID、状态和简短理由，不删除曾被引用的目标。
