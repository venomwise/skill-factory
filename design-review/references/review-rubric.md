# Design Review Rubric

对 `design.md` 评审的权威检查清单。自包含——不依赖其他技能的文件。逐维度评估，为每条发现标记严重性，然后计算裁决。

## Contents

- 规范章节集（完整的 design.md 包含什么）
- 7 个维度（D1–D7）及具体探查项
- 盲点检查清单（按领域）
- 过度设计信号
- 严重性定义
- 裁决与 spec-plan go/no-go 规则
- 评审纪律

## Canonical section set

由 brainstorming 技能产出的 design.md 遵循此结构。将其作为完整性（D1）和规范性（D3）的参照。标题为英文；章节内容可为任意语言。

| Section | 必需？ | 用途 |
|---------|--------|------|
| Summary | 是 | 一段话：正在构建什么以及为什么 |
| Goals | 是 | 主要成果 + 可衡量的成功指标 |
| Primary Users / Roles | 是 | 谁使用这个以及他们的关键目标 |
| Non-Goals | 是 | 明确超出范围的内容 |
| Context | 是 | 现有行为、约束、相关模块 |
| Discovery (Key Discoveries / Scope Decisions) | 非平凡或模糊请求建议包含 | 测试过的假设、浮现的盲点、范围内/外的内容及原因 |
| Decision Record (Options Considered + Decision & Rationale) | 是 | 权衡过的备选方案以及选择它的理由（或为什么只有一个可行） |
| Proposed Solution → Architecture | 是 | 关键构建块及其关系 |
| Proposed Solution → Components | 是 | 具有职责和接口的组件 |
| Proposed Solution → Data Flow | 是 | 主路径的逐步流程 |
| Error Handling | 是 | 主要失败模式及其处理方式 |
| Testing | 是 | 关键测试用例及其运行位置 |
| Open Questions | 是 | 仅已浮现的未解决问题；或明确写明没有遗留 |

缺失**必需**章节至少是 Major 发现；当缺失的章节对实现至关重要时（如缺少 Components、Data Flow 或 Error Handling）则是 Blocker。

## The 7 dimensions

每条发现记录：位置（`design.md §章节` 和/或项目文件）、证据、建议、严重性、以及所属维度。

### D1 — Completeness (完整性)

- 上表所有必需章节都存在且**有实质内容**（不是占位符，也不是复述标题）。
- Architecture、Components、Data Flow、Error Handling、Testing 全部覆盖。
- Goals 中每一项都映射到 Proposed Solution 中的内容；标记没有设计支撑的 goal。
- 每个 Component 都陈述了职责并定义了接口（inputs/outputs，而不只是一个名字）。
- Data Flow 覆盖主路径；非平凡设计至少覆盖关键失败场景。
- Decision Record 真实记录了权衡过的选项和理由——不是空的，也不是一行跳过权衡的套话。
- 方案依赖的数据模型 / schema 有定义。

### D2 — Usability / Actionability (可用性)

- `spec-plan` 能否不靠猜测就把这份设计变成需求和任务？如果关键决策被推迟到"实现时再定"，这是一条发现。
- 接口、inputs/outputs、数据形态足够具体，可以无歧义地实现。
- Data Flow 端到端可追溯——没有"然后系统处理它"式的含糊步骤。
- 术语一致且有定义；没有未定义的行话或缩写。
- Goals 中的成功标准可衡量，不是愿景式表述（没有数字或条件的"快"、"可扩展"）。

### D3 — Document Conformance (规范性)

- 遵循规范章节集，且结构**标题保持英文**。
- spec 位于约定路径（`specs/<topic>/`）；`<topic>` 为 kebab-case。
- **Open Questions** 仅包含 discovery 期间真实浮现的问题，或明确写明没有遗留——不是未完成设计的垃圾场。
- **Decision Record** 记录真实权衡过的备选方案；编造或套话式的备选方案是规范性问题，不是亮点。
- 非平凡或模糊请求存在 Discovery 章节。

### D4 — Project Fit (符合项目规范)

基于项目探索（workflow 第 3 步）。

- 与项目的技术栈、框架和语言约定一致（CLAUDE.md / AGENTS.md 规则、既有代码模式）。
- 尊重既有架构和模块边界；新组件放在项目本来会放的位置。
- **复用现有接缝**，而不是重新发明项目已有的能力——标记对既有工具、服务或模式的重复。
- 不悄悄冲突或破坏既有行为；如果改变了既有行为，设计要明说。
- 命名、文件布局和数据约定与项目一致。

### D5 — Blind Spots (盲点)

选择与设计领域相关的检查清单（见下文），浮现设计未覆盖的缺口。每条表述为"值得检查"，并核对 Open Questions 是否捕获了真正悬而未决的未知项，而不是让它们隐没。

### D6 — Over-Engineering (过度设计 / YAGNI)

对照设计**自己**的 Goals 和 Non-Goals 衡量。

- 无法追溯到任何已声明 Goal 的功能或组件——投机性的"将来可能用得上"工作。
- 违背已声明 Non-Goal 的范围蔓延。
- 过早抽象：为单一具体用途引入层、插件系统或通用框架。
- 防御性膨胀：处理逻辑上不可能发生的情况、没人要求的可配置性。
- 与问题不成比例的复杂度——更简单的结构就能满足所有 Goal。
- 每一条都指出仍能满足 Goals 的更简单替代方案。

### D7 — Optimization (优化点)

- 存在实质性更简单的方案能达成同样的 Goals。
- 更多复用项目既有组件可以减少新增面。
- 可以在不改变无关行为的前提下降低耦合。
- 项目（或领域）中已验证的模式比所提方案更契合。
- 这些是改进而非阻塞——通常是 Minor，除非与真实风险重叠。

## Blind-spot checklists

按领域选择；跳过明显不适用的项。不要机械地逐条遍历所有清单。

**任何面向用户的功能**
- 空 / 加载 / 错误状态
- 权限与访问控制
- 撤销或回滚
- 离线行为
- 无障碍访问
- 国际化 / 本地化
- 移动端或响应式行为

**任何数据功能**
- 从既有状态迁移数据
- 一致性与冲突解决
- 并发与幂等
- 保留与清理策略
- 隐私与合规（GDPR、PII 处理）
- 备份与恢复

**任何集成**
- 速率限制与配额
- 认证与凭证轮换
- 失败模式与回退 / 降级
- 版本化与向后兼容
- 超时、重试与重放 / 去重处理
- 监控与告警

## Over-engineering signals (quick scan)

- 唯一理由是"我们将来可能需要它"的组件。
- 只有一个实现、也看不到第二个调用方的抽象。
- 背后没有需求的配置、开关或扩展点。
- 直接调用就能解决，却引入通用机制（队列、缓存、调度器）。
- 一个 Non-Goal 被悄悄设计进去了。

## Severity definitions

| 严重性 | 含义 | 典型示例 |
|--------|------|----------|
| **Blocker** | 阻碍实现、内部自相矛盾、遗漏承重的必需章节、或与既有项目行为/约束冲突。必须在 `spec-plan` 前修复。 | 缺 Components / Data Flow / Error Handling；某个 Goal 完全没有设计；设计与既有系统的工作方式矛盾 |
| **Major** | 关键路径上的真实缺口或风险，应当解决，但单条不足以阻塞全部进展。 | 主流程上未覆盖的盲点；含糊的接口；实质性膨胀范围的过度设计；空的 Decision Record |
| **Minor** | 打磨或可选改进。 | 措辞；可以收紧的可衡量成功标准；非关键的优化或复用机会 |

## Verdict & spec-plan go/no-go rule

- **Reject** → No-Go：存在一条或多条 **Blocker** 发现。运行 `/design-refine` 处理发现，然后在规划前重新运行 `/design-review`。
- **Revise** → Conditional：无 Blocker，但存在一条或多条 **Major** 发现。规划前先运行 `/design-refine` 处理发现。
- **Pass** → Go：无 Blocker 且无 Major（允许 Minor）。有 Minor 发现时运行 `/design-refine`；无任何发现时直接进入 `spec-plan`。

## Reviewing discipline

- 每条发现都引用证据；指不出位置就不是发现。
- 区分**缺陷**（有东西错了）与**盲点**（未覆盖、值得检查）。后者表述为考虑事项。
- 对照设计自己声明的 Goals 和 Non-Goals 判断，而不是你偏好的另一种方案。
- Blocker 只留给真正的阻塞——一份夸大严重性的长清单，不如一份准确的短清单有用。
- 评审规模与设计体量匹配；小设计得到短报告。
