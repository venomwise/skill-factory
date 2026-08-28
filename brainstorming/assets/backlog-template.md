# <Project> Backlog

<!--
项目交付目标、顺序和状态以本文件为准。
work item 描述用户或业务结果，不描述文件、函数、API、表结构或实现任务。
正式 Acceptance Criteria 和技术方案只进入对应 `specs/<topic>/design.md`。
只有 `brainstorming` 可以修改 Status、Spec、Status evidence 和 Current Focus。
中文正文写清用户或业务对象、结果和条件，不用抽象项目术语代替实际内容。
-->

## Status Model

主状态：

```text
candidate -> ready-for-clarifying -> design-ready -> planned -> executing -> done
```

旁路状态：`blocked`、`deferred`、`cancelled`。

## Current Focus

- `<WI-...>` — <为什么现在推进，以及接下来由谁做什么>。

## Work Items

| ID | Outcome | Priority | Status | Dependencies |
|----|---------|----------|--------|--------------|
| `WI-<semantic-name>` | <谁在什么情况下获得什么结果> | <Now/Next/Later> | <status> | <WI-* 或 None> |

### WI-<semantic-name>: <Title>

- **Outcome**: <谁在什么情况下获得什么可判断的结果>。
- **Scope**: <高层能力边界>。
- **Non-Goals**: <本目标明确不包含的能力>。
- **Success signal**: <能判断该结果是否出现的项目级信号和证据来源>。
- **Priority**: <Now/Next/Later>。
- **Status**: <受支持状态>。
- **Dependencies**: <WI-*、外部前置或 None>。
- **Readiness**: <已经满足的条件，或具体阻塞和受影响范围>。
- **Spec**: `<真实 design.md 路径或 Not created>`。
- **Status evidence**: <支持当前状态的路径和可观察事实，或 None>。
- **External ref**: <真实 Jira/GitHub 标识或 None>。

#### Open Questions

- <留到本目标进入 `clarifying` 前后解决的问题；没有则写“无”>。

## Dependency Notes

- 说明无法仅通过 `Dependencies` 字段表达的业务、外部或交付顺序约束。

## Deferred / Cancelled

- 保留稳定 work item ID、状态和简短理由，不删除曾被引用的目标。
