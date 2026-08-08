---
name: design-review
description: Critically review a design.md (from the brainstorming skill or hand-written) before spec planning, and write a review.md verdict report next to it. Use when the user wants to review, critique, or quality-gate a design doc for completeness, usability, conformance, blind spots, over-engineering, or project-standards fit. Triggers include "评审设计", "review design", "design review", "检查设计文档", "设计盲点", "design.md review". Not for code or PR review, and not for writing or revising the design itself.
---

# Design Review

对 `design.md` 的独立质量关卡。读取设计及其周围项目，根据 7 维度评估标准进行评估，然后编写独立的 `review.md`，包含裁决和按严重性分级的发现。This skill is a critic：它指出问题供作者修复——它 not 自己填补空白，never 编辑设计。

## When to use

- 存在 `design.md`（来自 `brainstorming` skill 或手写），并且你想在规划前对其进行评审。
- As a gate between `brainstorming` and `spec-plan` —— 在它们变成任务之前捕获空白。
- 用户要求检查设计的盲点、过度设计、或与项目约定的契合度。

## When not to use

- Code or pull-request review → route to `code-review` / `review`。
- 从头编写设计 → route to `brainstorming`。
- 根据评审结果修订设计 → route to `design-refine`。
- 没有设计文档的琐碎更改、bug 修复、或仅修改拼写错误。

## Inputs

- 要评审的 `design.md` 路径（或自动探测 `specs/<topic>/`）。
- 设计所属的项目（为契合性和适配性探索）。

## Outputs

- 写入**与被评审的 `design.md` 相同目录**的 `review.md`，使用 `assets/review-template.md`。
- 整体裁决（Pass / Revise / Reject）、`spec-plan` go/no-go、以及标记为 Blocker / Major / Minor 的发现。
- 从发现和裁决派生的推荐下一步。
- `design.md` 和项目文件 never 被修改——`review.md` is the only artifact written。

## Workflow

1. **定位 design.md。** If 用户给了路径，使用它。Otherwise 探测 `specs/` 中的 `design.md`；if exactly 一个匹配，使用它；if 多个或没有匹配，向用户询问路径。When ambiguous 在评审前确认。
2. **完整阅读 design.md。** Hold 它预期遵循的规范章节集（Summary, Goals, Primary Users / Roles, Non-Goals, Context, Discovery, Decision Record, Proposed Solution with Architecture / Components / Data Flow, Error Handling, Testing, Open Questions）—— see [the rubric](references/review-rubric.md) for the authoritative list。
3. **按优先级顺序探索目标项目**：README / CLAUDE.md / AGENTS.md、项目配置（package.json / pyproject.toml / Cargo.toml / go.mod / pom.xml）、入口点、最近的 commits（最多 10 个）、然后任何现有的 specs 或设计文档。Stop when 你能将 conformance、project-fit、optimization、and over-engineering 检查建立在真实证据上。When the design 依赖存储的数据或 schema 时使用 `db-explorer`。Do not assume —— when 设计中的声明可以根据代码库检查，check it。
4. **根据 [references/review-rubric.md](references/review-rubric.md) 中的 7 维度评估标准进行评估。** For every finding，记录：具体的 **location**（`design.md §section` and/or 项目文件路径）、**evidence**、具体的 **recommendation**、and a **severity**。Ground every finding —— no vague criticism。Judge over-engineering against the design's **own** stated Goals and Non-Goals，not 你的品味。Frame blind spots as "值得检查"，not as 作者肯定遗漏的缺陷。
5. **从严重性计数计算裁决和 `spec-plan` go/no-go**（see the rubric for the exact rule）。
6. **写入 review.md** 到 `design.md` 同级目录，使用 `assets/review-template.md`。Preserve the template's English structural headings；write the finding content in the design's current language（从 `design.md` 和对话中推断；do not ask solely to determine it）。
7. **总结并推荐。** Give the user 裁决、发现计数、and the top blockers in a few lines。If any Blocker, Major, or Minor finding exists，推荐使用被评审的 `design.md` 和生成的 `review.md` 运行 `/design-refine` 来讨论发现并更新设计；for a Reject verdict，also 推荐在 refinement 后重新运行 `/design-review`。Only recommend proceeding directly to `spec-plan` when 没有发现。Do not modify `design.md` or invoke any downstream skill yourself。

## Review dimensions

Full checklists live in [references/review-rubric.md](references/review-rubric.md). Summary:

| # | Dimension | Core question |
|---|-----------|---------------|
| D1 | Completeness (完整性) | Are all required sections present and substantive, with architecture, components, data flow, error handling, testing, and a real Decision Record? Does every goal map to a solution element? |
| D2 | Usability / Actionability (可用性) | Is it concrete and unambiguous enough for `spec-plan` to consume — interfaces, data shapes, traceable flow, measurable success criteria, consistent terms? |
| D3 | Document Conformance (规范性) | Does it follow the design-doc template structure and the doc conventions (English headings, spec path/naming, Open Questions and Decision Record rules)? |
| D4 | Project Fit (符合项目规范) | Is it consistent with the project's stack, architecture, and module boundaries? Does it reuse existing seams instead of duplicating or conflicting with existing behavior? |
| D5 | Blind Spots (盲点) | What did it not address — states, permissions, migration, consistency, retention, privacy, rate limits, auth, failure/fallback, versioning, monitoring, concurrency, idempotency? Did Open Questions capture the right unknowns? |
| D6 | Over-Engineering (过度设计 / YAGNI) | Any speculative features beyond the Goals, scope creep against Non-Goals, premature abstraction, defensive bloat, or complexity out of proportion to the problem? |
| D7 | Optimization (优化点) | Is there a simpler approach, more reuse, less coupling, or a better-fit existing pattern? |

Severity (Blocker / Major / Minor) and verdict (Pass / Revise / Reject) definitions are in the rubric.

## Verification

- [ ] `review.md` 存在于与被评审的 `design.md` 相同的目录中
- [ ] 报告陈述了整体裁决、`spec-plan` go/no-go、以及 Blocker / Major / Minor 计数
- [ ] All 7 dimensions 出现在 Dimension Summary 中
- [ ] Every finding 引用具体位置、给出证据和推荐、并携带严重性标签
- [ ] Over-engineering findings are judged against the design's own Goals / Non-Goals
- [ ] Blind spots are framed as considerations，grounded in the domain checklists
- [ ] `design.md` 和所有项目文件 are unchanged
- [ ] Structural headings are in English；finding content is in the design's current language
- [ ] The recommended next step follows from the findings and verdict：`/design-refine` when findings exist，`spec-plan` only when none exist

## Safety & guardrails

- 对设计和项目保持只读。The only file written is `review.md`；never edit `design.md` or any project file。
- 证据优先于观点。Every finding 指向具体位置并陈述支持它的内容。No unsubstantiated criticism。
- Do not fabricate。Do not invent alternatives, requirements, or blind spots that do not apply；review 存在的内容 against the design's own stated scope。
- 严重性纪律。Reserve Blocker for issues that truly stop implementation or contradict the project；do not inflate to look thorough。
- 尊重设计的范围。Over-engineering and optimization are measured against the design's Goals and Non-Goals，not a different solution you would have preferred。
- 根据复杂度调整规模。A small design gets a short review；do not pad with ceremony。
- 无实现、规划或修订。This skill reviews and recommends only —— it does not write code, edit the design, or invoke `design-refine` or `spec-plan`。

## References

- [Design review rubric](references/review-rubric.md) — 7-dimension checklists, severity and verdict definitions, canonical section set, blind-spot and over-engineering signals
- [Review report template](assets/review-template.md)
- [Design refine skill](../design-refine/SKILL.md) — discusses review findings with the user and writes confirmed decisions back to `design.md`
