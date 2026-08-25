# <Topic> Design

<!--
Canonical design document contract for the brainstorming -> design-review -> design-refine -> spec-plan chain.
Replace guidance text with project-specific content. Remove unused optional subsections.

Formatting rules:
- Wrap ordinary Markdown lines at 120 Unicode characters.
- Keep tables to 5-6 short columns; move explanations below the table.
- Fenced code blocks and standalone unbreakable URLs/hashes may exceed the limit.
- Validate the completed document with scripts/check-markdown-lines.mjs from the brainstorming skill.
-->

## Summary

用一段话说明正在构建什么、为谁构建以及为什么。

## Goals

- 列出主要成果和可衡量的成功条件。

## Primary Users / Roles

- 列出使用者、调用方或运维角色及其目标。

## Non-Goals

- 明确当前范围之外的能力，防止范围蔓延。

## Context

说明现有行为、项目约束、相关模块和已经核实的事实。

## Discovery

<!-- 非平凡或起初含糊的需求保留本节；简单设计可省略。 -->

### Key Discoveries

- 记录在 brainstorming 中验证过的假设和影响设计的事实。

### Scope Decisions

- 记录明确纳入或排除的范围及其理由。

## Decision Record

只记录持久且非显然的设计选择。缺字段、格式修复、事实对齐和章节传播不属于 Decision。
Decision 使用稳定的语义 ID；同一主题发生变化时原位更新，不追加 `Revised` 副本。

### DR-<semantic-topic>

**Decision**: <最终选择及其有效边界>。

**Rationale**: <为什么该选择优于其他可行方向；被拒方向只写名称和拒绝理由>。

**Constraints**: <选择成立所依赖的前提和接受的取舍>。

<!-- 仅当用户明确驳回一个可能被后续重复提出的问题时保留以下字段。 -->

**Rejected concern**: <明确不纳入设计的问题边界，不展开被拒方案的实现细节>。

**Revisit when**: <哪些需求或项目事实发生变化时需要重新评估>。

## Proposed Solution

说明总体方案以及它为什么满足已确认的目标和约束。

### Architecture

描述关键构建块和关系。契约、数据模型和机制细节放入后续子章节。

### Components

按“触发条件 -> 输入 -> 产出或副作用”描述组件职责。不给出方法签名。

### Data Model

<!--
出现实体、表、字段、序列化模型、迁移或并发控制变化时保留本节。
只生成适用的子章节。可执行 DDL/DML 放在独立 sql/ 文件中，本节只引用路径和总结语义。
-->

#### Change Summary

| Change | Layer | Object | Purpose | Source |
|--------|-------|--------|---------|--------|
| <ADD/MODIFY/REMOVE> | <Database/Domain/Serialized> | <名称> | <简述> | <文件或章节> |

#### Persistence Model

##### [<ADD/MODIFY>] <table-or-document>

<!-- 新对象写完整目标模型；已有对象只写新增、修改和删除字段。 -->

| Field | Type | Null | Default | Constraint |
|-------|------|------|---------|------------|
| `<field>` | `<type>` | <Yes/No> | <value> | <短约束> |

- **Ownership**: <负责读写的组件>。
- **Lifecycle**: <创建、更新、删除和保留语义>。
- **Source of truth**: <权威存储>。

#### Domain / Serialized Models

##### [<ADD/MODIFY>] <model>

- 描述最终字段、嵌套结构、兼容性和事实源。
- 长字段语义使用列表，不塞入表格单元格。

#### Constraints and Indexes

- 描述唯一性、索引、关联、归属和必须由存储层保证的约束。

#### Migration / Backfill / Rollback

- **Forward SQL**: `<path/to/forward.sql>`
- **Rollback SQL**: `<path/to/rollback.sql>`
- **Backfill**: <数据转换或 None>。
- **Rollback guard**: <阻止数据丢失的前置条件>。

#### Concurrency and Consistency

- 描述乐观锁、幂等键、事务边界、一致性要求及其承载字段或机制。

### Interfaces

<!--
出现新增、修改或删除的 API、消息或协议契约时保留本节。
完全未变化的单个接口不进入 Change Summary；没有某类契约变化时明确写 No changes。
-->

#### Change Summary

| Change | Kind | Contract ID | Producer | Consumer | Compatibility |
|--------|------|-------------|----------|----------|---------------|
| <ADD/MODIFY/REMOVE> | <HTTP/Message/Protocol> | <IF-...> | <调用方> | <接收方> | <影响> |

- **Messages**: <No changes 或变更摘要>。
- **Protocols**: <No changes 或变更摘要>。

#### HTTP APIs

##### [<ADD/MODIFY/REMOVE>] <IF-HTTP-semantic-name>

- **Method / Path**: `<METHOD /path>`
- **Auth**: <认证和权限>。
- **Content-Type**: `<content-type>`
- **Idempotency**: <幂等、重复调用或冲突语义>。

###### Request

| Field | Location | Type | Required | Meaning |
|-------|----------|------|----------|---------|
| `<field>` | <path/query/header/body/multipart> | `<type>` | <Yes/No> | <短说明> |

###### Response

- 描述成功状态、响应结构和重要字段语义。

###### Errors

| Condition | Result |
|-----------|--------|
| <条件> | <状态、错误码或安全行为> |

###### Compatibility

说明 `MODIFY` 的原契约到新契约变化，或 `REMOVE` 的弃用和迁移规则。

#### Messages

##### [<ADD/MODIFY/REMOVE>] <IF-MSG-semantic-name>

- **Channel / Topic**: `<name>`
- **Producer / Consumer**: <双方>。
- **Trigger**: <发送时机>。
- **Envelope / Payload**: <字段结构或机器可读契约路径>。
- **Delivery**: <投递保证、顺序、重试、去重和 DLQ>。
- **Versioning**: <版本与兼容规则>。

#### Protocols

##### [<ADD/MODIFY/REMOVE>] <IF-PROTO-semantic-name>

- **Direction**: <发送方 -> 接收方>。
- **Command / Opcode**: `<value>`
- **Framing / Payload**: <帧和载荷结构或契约路径>。
- **ACK / Timeout / Retry**: <确认与失败语义>。
- **Versioning**: <协议版本和兼容规则>。

### State Machines

<!-- 出现显式生命周期或状态转换时保留本节。 -->

| Current state | Trigger / Guard | Target state | Side effects |
|---------------|-----------------|--------------|--------------|
| <state> | <event and guard> | <state> | <effects> |

### Key Mechanisms

<!-- 出现跨组件可靠投递、对账、恢复、协作或其他机制时保留本节。 -->

#### <Mechanism>

描述参与者、触发条件、步骤、不变量和失败恢复。

### Data Flow

按顺序描述主路径。非平凡设计同时覆盖关键失败路径。
Data Flow 中的跨边界调用必须引用 Interfaces 中的 Contract ID。

## Error Handling

列出主要失败模式、检测位置、返回语义、重试或回滚行为。

## Acceptance Criteria

按行为域组织已确认的可观察行为。每条标准使用唯一、稳定的
`AC-<domain>-<behavior>` kebab-case ID，并只描述一个可验证结果。

### <Behavior Domain>

- **AC-<domain>-<behavior>**:
  WHEN <条件>, THEN 系统 SHALL <可观察结果>。
- **AC-<domain>-<failure>**:
  IF <错误或边界条件>, THEN 系统 SHALL <安全行为>。

覆盖每个核心 Goal 的正常流、错误流和关键边界。不要写入架构理由、被拒方案、文件路径、
实现步骤、测试命令或未经确认的未来想法。

## Testing

描述关键测试层级、范围和运行位置。按需引用 AC ID；具体验证命令留给实施计划和 Checkpoint。

## Open Questions

只列出 discovery 中已经浮现、且尚未解决的问题。若没有，明确写“无遗留问题”。
任何会改变行为、范围或 Acceptance Criteria 的问题都必须在进入 `spec-plan` 前解决。
