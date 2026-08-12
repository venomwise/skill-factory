---
name: design-review
description: 对 design.md（来自 brainstorming 技能或手写）进行批判性评审，并在其旁边写入 review.md 裁决报告。当用户想评审设计文档、为规划做质量关卡，或检查完整性、可用性、规范性、盲点、过度设计、项目契合度时调用。触发词包括"评审设计"、"review design"、"design review"、"检查设计文档"、"设计盲点"、"design.md review"。不评审代码或 PR；不编写设计（归 brainstorming）；不根据评审结果修订设计（归 design-refine）。
---

# Design Review

对 `design.md` 的独立质量关卡。读取设计及其周围项目，按 7 维度评估标准评估，然后写入独立的 `review.md`，包含裁决和按严重性分级的发现。本技能是批评者：指出问题供作者修复——不自己填补空白，也从不编辑设计。

## When to use

- 存在 `design.md`（来自 `brainstorming` 技能或手写），并且想在规划前对其进行评审。
- 作为 `brainstorming` 与 `spec-plan` 之间的关卡——在问题变成任务之前捕获它们。
- 用户要求检查设计的盲点、过度设计、或与项目约定的契合度。

## When not to use

- 代码或 pull request 级别的评审不在本技能范围内，应交回相应的代码评审工具。
- 从头编写设计 → 使用 `brainstorming`。
- 根据评审结果修订设计 → 使用 `design-refine`。
- 没有设计文档的琐碎更改、bug 修复、或仅修改拼写错误。

## Inputs

- 要评审的 `design.md` 路径（或自动探测 `specs/<topic>/`）。
- 设计所属的项目（用于契合性和适配性探索）。

## Outputs

- 写入**与被评审的 `design.md` 相同目录**的 `review.md`，使用 `assets/review-template.md`。
- 整体裁决（Pass / Revise / Reject）、`spec-plan` go/no-go、以及标记为 Blocker / Major / Minor 的发现。
- 从发现和裁决派生的推荐下一步。
- `design.md` 和项目文件一律不修改——`review.md` 是唯一写入的产物。

## Workflow

1. **定位 design.md。** 用户给了路径就直接用。否则探测 `specs/` 中的 `design.md`：恰好一个匹配就直接使用；多个或没有匹配时向用户询问路径。有歧义时先确认再评审。
2. **完整阅读 design.md。** 记住它预期遵循的规范章节集（Summary、Goals、Primary Users / Roles、Non-Goals、Context、Discovery、Decision Record、Proposed Solution 含 Architecture / Components / Data Flow、Error Handling、Testing、Open Questions）——权威清单见 [the rubric](references/review-rubric.md)。
3. **按优先级顺序探索目标项目**：README / CLAUDE.md / AGENTS.md、项目配置（package.json / pyproject.toml / Cargo.toml / go.mod / pom.xml）、入口点、最近的提交（最多 10 个）、然后是任何现有的 specs 或设计文档。STOP when 你能将规范性、项目契合度、优化点和过度设计检查都建立在真实证据上。当设计依赖存储的数据或 schema 时使用 `db-explorer`。不要假设——设计中的声明只要能对照代码库核实，就去核实它。
4. **按 [references/review-rubric.md](references/review-rubric.md) 中的 7 维度评估标准进行评估。** 每条发现都记录：具体的**位置**（`design.md §章节` 和/或项目文件路径）、**证据**、具体的**建议**、以及**严重性**。每条发现都要有依据——不接受模糊的批评。过度设计要对照设计**自己**声明的 Goals 和 Non-Goals 判断，而不是你的品味。盲点表述为"值得检查"，而不是作者肯定遗漏的缺陷。
5. **从严重性计数计算裁决和 `spec-plan` go/no-go**（精确规则见 rubric）。
6. **写入 review.md** 到 `design.md` 同级目录，使用 `assets/review-template.md`。保留模板的英文结构标题；发现内容使用设计当前的语言书写（从 `design.md` 和对话中推断；不要仅为确定语言而询问）。
7. **总结并推荐。** 用几行给出裁决、发现计数和最重要的 Blocker。只要存在任何 Blocker、Major 或 Minor 发现，就推荐使用被评审的 `design.md` 和生成的 `review.md` 运行 `/design-refine` 来讨论发现并更新设计；裁决为 Reject 时，同时推荐在 refine 之后重新运行 `/design-review`。只有在没有任何发现时，才推荐直接进入 `spec-plan`。不要修改 `design.md`，也不要自行调用任何下游技能。

## Review dimensions

完整清单见 [references/review-rubric.md](references/review-rubric.md)。摘要：

| # | Dimension | 核心问题 |
|---|-----------|----------|
| D1 | Completeness (完整性) | 所有必需章节是否都存在且有实质内容，包括 architecture、components、data flow、error handling、testing 和真实的 Decision Record？每个 goal 是否都映射到方案元素？ |
| D2 | Usability / Actionability (可用性) | 是否足够具体、无歧义，`spec-plan` 能直接消费——接口、数据形态、可追溯的流程、可衡量的成功标准、一致的术语？ |
| D3 | Document Conformance (规范性) | 是否遵循 design-doc 模板结构和文档约定（英文标题、spec 路径/命名、Open Questions 与 Decision Record 规则）？ |
| D4 | Project Fit (符合项目规范) | 是否与项目的技术栈、架构和模块边界一致？是否复用现有接缝，而不是重复或与现有行为冲突？ |
| D5 | Blind Spots (盲点) | 有哪些没有覆盖——状态、权限、迁移、一致性、留存、隐私、限流、认证、失败/回退、版本化、监控、并发、幂等？Open Questions 是否捕获了正确的未知项？ |
| D6 | Over-Engineering (过度设计 / YAGNI) | 是否有超出 Goals 的投机性功能、违背 Non-Goals 的范围蔓延、过早抽象、防御性膨胀、或与问题不成比例的复杂度？ |
| D7 | Optimization (优化点) | 是否有更简单的方案、更多的复用、更少的耦合、或更契合的现有模式？ |

严重性（Blocker / Major / Minor）和裁决（Pass / Revise / Reject）的定义见 rubric。

## Verification

- [ ] `review.md` 存在于与被评审的 `design.md` 相同的目录中
- [ ] 报告陈述了整体裁决、`spec-plan` go/no-go、以及 Blocker / Major / Minor 计数
- [ ] 7 个维度全部出现在 Dimension Summary 中
- [ ] 每条发现都引用具体位置、给出证据和建议、并携带严重性标签
- [ ] 过度设计的发现是对照设计自己的 Goals / Non-Goals 判断的
- [ ] 盲点以考虑事项的方式表述，并基于领域检查清单
- [ ] `design.md` 和所有项目文件未被修改
- [ ] 结构标题为英文；发现内容使用设计当前的语言
- [ ] 推荐的下一步与发现和裁决一致：有发现时指向 `/design-refine`，仅当无发现时指向 `spec-plan`

## Safety & guardrails

- 对设计和项目保持只读。唯一写入的文件是 `review.md`；不编辑 `design.md` 或任何项目文件。
- 证据优先于观点。每条发现都指向具体位置并陈述支持它的内容；不做无依据的批评。
- 不要编造。不要虚构不适用的替代方案、需求或盲点；只评审存在的内容，并对照设计自己声明的范围。
- 严重性纪律。Blocker 只留给真正阻碍实现或与项目矛盾的问题；不要为了显得周全而夸大。
- 尊重设计的范围。过度设计和优化点都对照设计的 Goals 和 Non-Goals 衡量，而不是你更偏好的另一种方案。
- 根据复杂度调整规模。小设计得到短评审；不要用仪式填充。
- 无实现、无规划、无修订。本技能只评审和推荐——不写代码、不修改设计、不调用 `design-refine` 或 `spec-plan`。

## References

- [Design review rubric](references/review-rubric.md) — 7 维度检查清单、严重性与裁决定义、规范章节集、盲点与过度设计信号
- [Review report template](assets/review-template.md)
- [Design refine skill](../design-refine/SKILL.md) — 与用户讨论评审发现，并将确认的决策写回 `design.md`
