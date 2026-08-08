# Brainstorming Guide

## Context exploration priorities

1. README, CLAUDE.md, 或等效的项目文档。
2. 项目配置：package.json, pyproject.toml, Cargo.toml, go.mod。
3. 入口点：main files, index files, route definitions。
4. 最近的 commits（最多 10 个）以了解当前动向和约定。
5. Stop when 你能陈述项目的目的、技术栈和目录结构。

## Decomposition heuristics

- Split when 请求跨越多个领域或独立能力。
- Split when 各部分需要不同的数据模型或外部集成。
- Split when 各部分可以由不同团队拥有或单独部署。
- When splitting，列出子项目、依赖关系和推荐顺序。

## Good scope split examples

- "Chat + file storage + billing" 变成三个子项目，各有自己的 spec。
- "Redesign UI + add analytics" 变成一个 UI 子项目和一个 metrics 子项目。

## Question patterns

- Purpose: 谁是用户，这使什么决策成为可能？
- Constraints: latency, scale, security, compliance, 或运营限制。
- Success criteria: 我们如何知道它起作用了？
- Data: inputs, outputs, 以及 source of truth。
- Risks: 要处理的关键失败模式或边缘情况。

## Design depth guidance

- Simple changes 仍使用所有章节，but keep 每个章节只有几句话。
- Moderate changes 包含逐步的 data flow 和关键错误情况。
- Complex changes 包含主要 happy path 加上 2-3 个关键失败场景。

## Existing codebase tactics

- 识别可重用的接缝并尊重本地模式。
- Challenge 在设计之前请求的形状是否适合当前行为。
- Ground the check 在真实证据而非假设中——当请求依赖存储的数据或 schema 时使用可用工具如 `db-explorer`。
- Refactors are allowed only when 它们解除对当前目标的阻塞。
- Prefer 减少耦合而不改变无关行为的变更。

## Post-design complexity examples

从完成的设计分类，not 用户初始提示的大小。Simple 路由故意狭窄。

**Usually simple:** 局部验证规则、一个 UI 状态调整、或遵循现有模式并有针对性测试的小命令行为变更。

**Moderate or complex:** schema 或数据迁移、新的外部服务、身份验证或权限变更、公共 API 契约变更、跨组件工作流、并发或一致性要求、运营推出风险、或任何未解决的设计问题。

When signals are mixed，推荐 `/design-review`。当不正确的设计难以逆转或会影响多个消费者时，review 的成本是合理的。

## Fuzziness diagnosis

在选择策略之前，根据这些级别评估用户的请求：

| Level | Signal | Example | Strategy |
|-------|--------|---------|----------|
| Problem unclear | 用户描述症状，not goals；"感觉不对" | "系统难以使用" | Reframe the problem |
| Direction unclear | 用户有目标但没有方法感觉 | "我想要更好的用户参与度" | Explore possibilities |
| Boundaries unclear | 用户知道他们想要什么但不知道边界 | "我想要用户认证" | Scan for blind spots |
| Solution unclear | 用户知道什么和范围，需要技术方法 | "我想要带 SAML 的 SSO" | Confirm intent, priorities, and constraints first, then compare approaches |

## Assumption challenging

对于非平凡的请求，识别潜在假设并提供给用户确认。Frame them as "值得检查" rather than "你漏了这个"。

1. 识别用户可能隐含做出的 1-3 个假设。
2. 对于每个，ask: "这一定是真的吗？如果不是呢？"
3. 值得检查的常见假设：
   - "用户会按我想象的方式使用此功能"
   - "当前架构可以在不改变的情况下支持这个"
   - "这需要在一个发布中完成"
   - "性能或规模不会是问题"
   - "现有数据模型足够了"

## Blind spot scanning

选择与项目领域最相关的 1-2 个检查清单。Do not 遍历每个列表；按领域选择并跳过明显不适用的项目。

**Any user-facing feature:**
- Empty states, loading states, error states
- Permissions and access control
- Undo or rollback
- Offline behavior
- Accessibility
- Internationalization or localization
- Mobile or responsive behavior

**Any data feature:**
- 从现有状态迁移数据
- Consistency and conflict resolution
- Retention and cleanup policies
- Privacy and compliance (GDPR, etc.)
- Backup and recovery

**Any integration:**
- Rate limits and quotas
- Authentication and credential rotation
- Failure modes and fallback behavior
- Versioning and backwards compatibility
- Monitoring and alerting

## Divergent exploration techniques

在 brainstorming 期间使用这些来扩展问题空间：

- **Constraint removal**: "如果我们没有技术、时间或预算限制，理想情况会是什么样？" Then 逐个添加回约束。
- **Negative brainstorming**: "什么会让这个功能完全失败？" Flip 每个失败为一个需求。
- **Perspective switching**: 从不同角色考虑同一功能（end user, admin, ops engineer, new hire, power user）。
- **Time horizon**: "这在 6 个月后如何演进？一年后？" Identify 哪些决策难以逆转。
- **Priority forcing**: "如果你只能发布 3 件事，是哪 3 件？" This 强制用户揭示真正重要的内容。
- **Analogy**: "类似产品或领域如何解决这个？" Borrow 经过验证的模式。

## When to stop brainstorming

The authoritative exit condition lives in SKILL.md step 4; do not duplicate or override it here. Use 这些作为支持信号，表明条件可能满足：

- 用户能陈述他们想要什么、为什么想要、以及明确不想要什么。
- 关键约束和成功标准已知或明确记录为假设。
- For non-trivial requests，至少一个假设已被确认，一个盲点已被考虑。
- 问题陈述是稳定的，并且在最近 2 次交流中没有改变。
