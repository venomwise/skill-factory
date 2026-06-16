---
name: design-review
description: Critically review a design.md (from the brainstorming skill or hand-written) before spec planning, and write a review.md verdict report next to it. Use when the user wants to review, critique, or quality-gate a design doc for completeness, usability, conformance, blind spots, over-engineering, or project-standards fit. Triggers include "评审设计", "review design", "design review", "检查设计文档", "设计盲点", "design.md review". Not for code or PR review, and not for writing or revising the design itself.
---

# Design Review

Independent quality gate for a `design.md`. Read the design and its surrounding project, evaluate against a 7-dimension rubric, then write a standalone `review.md` with a verdict and severity-graded findings. This skill is a critic: it points out problems for the author to fix — it does not fill the gaps itself and never edits the design.

## When to use

- A `design.md` exists (from the `brainstorming` skill or hand-written) and you want it reviewed before planning.
- As a gate between `brainstorming` and `spec-plan` — catch gaps before they become tasks.
- The user asks to check a design for blind spots, over-engineering, or fit with the project's conventions.

## When not to use

- Code or pull-request review → route to `code-review` / `review`.
- Writing or revising the design itself → route to `brainstorming`.
- Trivial changes with no design doc, bug fixes, or typo-only edits.

## Inputs

- Path to the `design.md` to review (or auto-probe `specs/<topic>/` and `.codex/specs/<topic>/`).
- The project the design belongs to (explored for conformance and fit).

## Outputs

- A `review.md` written to the **same directory as the reviewed `design.md`**, using `assets/review-template.md`.
- An overall verdict (Pass / Revise / Reject), a `spec-plan` go/no-go, and findings tagged Blocker / Major / Minor.
- A recommended next step derived from the verdict.
- The `design.md` and project files are never modified — `review.md` is the only artifact written.

## Workflow

1. **Locate the design.md.** If the user gave a path, use it. Otherwise probe `specs/` and `.codex/specs/` for a `design.md`; if exactly one matches, use it; if several or none match, ask the user for the path. Confirm before reviewing when ambiguous.
2. **Read the design.md fully.** Hold the canonical section set it is expected to follow (Summary, Goals, Primary Users / Roles, Non-Goals, Context, Discovery, Decision Record, Proposed Solution with Architecture / Components / Data Flow, Error Handling, Testing, Open Questions) — see [the rubric](references/review-rubric.md) for the authoritative list.
3. **Explore the target project** in priority order: README / CLAUDE.md / AGENTS.md, project config (package.json / pyproject.toml / Cargo.toml / go.mod / pom.xml), entry points, recent commits (up to 10), then any existing specs or design docs. Stop when you can ground the conformance, project-fit, optimization, and over-engineering checks in real evidence. Use `db-explorer` when the design depends on stored data or schema. Do not assume — when a claim in the design can be checked against the codebase, check it.
4. **Evaluate against the 7-dimension rubric** in [references/review-rubric.md](references/review-rubric.md). For every finding, record: a concrete **location** (`design.md §section` and/or a project file path), the **evidence**, a concrete **recommendation**, and a **severity**. Ground every finding — no vague criticism. Judge over-engineering against the design's **own** stated Goals and Non-Goals, not your taste. Frame blind spots as "worth checking", not as defects the author definitely missed.
5. **Compute the verdict and the `spec-plan` go/no-go** from the severity counts (see the rubric for the exact rule).
6. **Write review.md** to the `design.md` sibling directory using `assets/review-template.md`. Preserve the template's English structural headings; write the finding content in the design's current language (infer it from the `design.md` and the conversation; do not ask solely to determine it).
7. **Summarize and recommend.** Give the user the verdict, the finding counts, and the top blockers in a few lines, then recommend the next step from the verdict: return to `brainstorming` to resolve blockers, or proceed to `spec-plan` when clean. Do not modify `design.md`; do not invoke `spec-plan` yourself.

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

- [ ] `review.md` exists in the same directory as the reviewed `design.md`
- [ ] The report states an overall verdict, a `spec-plan` go/no-go, and Blocker / Major / Minor counts
- [ ] All 7 dimensions appear in the Dimension Summary
- [ ] Every finding cites a concrete location, gives evidence and a recommendation, and carries a severity tag
- [ ] Over-engineering findings are judged against the design's own Goals / Non-Goals
- [ ] Blind spots are framed as considerations, grounded in the domain checklists
- [ ] The `design.md` and all project files are unchanged
- [ ] Structural headings are in English; finding content is in the design's current language
- [ ] A recommended next step is given and follows from the verdict

## Safety & guardrails

- Read-only on the design and the project. The only file written is `review.md`; never edit `design.md` or any project file.
- Evidence over opinion. Every finding points to a concrete location and states what supports it. No unsubstantiated criticism.
- Do not fabricate. Do not invent alternatives, requirements, or blind spots that do not apply; review what is there against the design's own stated scope.
- Severity discipline. Reserve Blocker for issues that truly stop implementation or contradict the project; do not inflate to look thorough.
- Respect the design's scope. Over-engineering and optimization are measured against the design's Goals and Non-Goals, not a different solution you would have preferred.
- Scale to complexity. A small design gets a short review; do not pad with ceremony.
- No implementation, no planning. This skill reviews and recommends only — it does not write code, scaffold, or invoke `spec-plan`.

## References

- [Design review rubric](references/review-rubric.md) — 7-dimension checklists, severity and verdict definitions, canonical section set, blind-spot and over-engineering signals
- [Review report template](assets/review-template.md)
