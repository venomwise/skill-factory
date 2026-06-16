# Design Review Skill Test Cases

Use these to validate that the `design-review` skill reads a `design.md`, evaluates it against the 7-dimension rubric, and writes a `review.md` with an accurate verdict and severity-graded, evidence-grounded findings — without modifying the source.

Each case points the skill at a fixture under `fixtures/<case>/design.md`.

## Test 1 — Clean design: orders pagination

**Prompt**

> 评审这份设计文档：fixtures/clean/design.md

**Expected behavior**

- Write `review.md` into `fixtures/clean/`.
- Verdict **Pass**, spec-plan readiness **Go**.
- Zero Blockers and zero Majors; do not invent issues to look thorough.
- All 7 dimensions marked sound (or with only Minor notes).
- Leave the original `design.md` unchanged.

**Primary coverage**

- Accurate Pass verdict
- No fabricated findings
- Read-only guarantee

## Test 2 — Incomplete design: bulk order export

**Prompt**

> 评审这份设计文档：fixtures/incomplete/design.md

**Expected behavior**

- Verdict **Reject**, spec-plan readiness **No-Go**.
- Blocker(s) for the missing load-bearing sections (Error Handling, Testing).
- Flag the empty/placeholder Decision Record and the missing Open Questions section.
- Flag the in-app notification goal that has no design support.
- Every finding cites a concrete location and gives a recommendation.

**Primary coverage**

- D1 Completeness (missing sections, unsupported goal)
- D3 Document Conformance (empty Decision Record, missing Open Questions)
- Severity discipline (true Blockers)

## Test 3 — Over-engineered design: theme preference

**Prompt**

> 评审这份设计文档：fixtures/over-engineered/design.md

**Expected behavior**

- Acknowledge the doc is structurally complete (all sections present).
- Dominant findings are **D6 Over-Engineering**: plugin registry, multiple storage backends, write-through cache, event bus, and audit subscriber are unjustified by the Goals.
- Point out the design contradicts the Non-Goals (single value, theme only, no custom themes) and that the Decision Record justifies the framework with speculative future needs.
- Recommend the simpler `single column users.theme` already listed in Options Considered.
- Judge over-engineering against the design's own Goals/Non-Goals. Verdict **Revise** (no fabricated Blocker).

**Primary coverage**

- D6 Over-Engineering / YAGNI
- Judging against the doc's own scope
- Recommending the simpler in-doc alternative

## Test 4 — Hidden blind spots: nightly CRM contact sync

**Prompt**

> 评审这份设计文档：fixtures/blind-spots/design.md

**Expected behavior**

- Recognize the design looks complete but ignores integration/data blind spots.
- Surface **D5 Blind Spot** findings, including at least:
  - Acme API rate limits / quotas while paging all contacts.
  - Failure handling beyond "log the error and stop" — no retry, resume, or partial-failure recovery; the table is left half-updated.
  - Upstream deletions never removed by a full-reload upsert (stale rows persist).
  - No conflict resolution for locally edited contacts; no monitoring/alerting on failure.
- Flag that Open Questions says "None" when real unknowns remain (deletion + conflict policy) — a D3 finding.
- Frame blind spots as worth checking rather than certain defects. Verdict **Revise**.

**Primary coverage**

- D5 Blind Spots (integration + data checklists)
- D3 Open Questions accuracy
- "Worth checking" framing

## Coverage Matrix

| Test | Fixture | Dominant signal | Expected verdict |
|------|---------|-----------------|------------------|
| 1 | clean | Sound design, accurate Pass | Pass / Go |
| 2 | incomplete | Missing sections, unsupported goal | Reject / No-Go |
| 3 | over-engineered | YAGNI vs the doc's own Non-Goals | Revise / Conditional |
| 4 | blind-spots | Unaddressed integration/data risks | Revise / Conditional |

## Pass Criteria

A successful run of the skill should consistently:

- write `review.md` to the reviewed `design.md`'s sibling directory and leave the source unmodified
- state an overall verdict, a spec-plan go/no-go, and Blocker/Major/Minor counts
- cover all 7 dimensions in the Dimension Summary
- ground every finding in a concrete location with a recommendation and a severity tag
- judge over-engineering against the design's own Goals/Non-Goals
- frame blind spots as considerations, not as certain defects
- reserve Blocker for genuine blockers and not inflate severities
