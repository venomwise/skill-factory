# Design Review Rubric

Authoritative checklists for reviewing a `design.md`. Self-contained — do not depend on other skills' files. Apply the dimensions, tag each finding with a severity, then compute the verdict.

## Contents

- Canonical section set (what a complete design.md contains)
- The 7 dimensions (D1–D7) with concrete probes
- Blind-spot checklists (by domain)
- Over-engineering signals
- Severity definitions
- Verdict & spec-plan go/no-go rule
- Reviewing discipline

## Canonical section set

A design.md produced by the brainstorming skill follows this structure. Use it as the reference for completeness (D1) and conformance (D3). Headings are in English; section content may be in any language.

| Section | Required? | Purpose |
|---------|-----------|---------|
| Summary | Yes | One paragraph: what is being built and why |
| Goals | Yes | Primary outcomes + measurable success indicators |
| Primary Users / Roles | Yes | Who uses this and their key goals |
| Non-Goals | Yes | Explicitly out of scope |
| Context | Yes | Existing behavior, constraints, related modules |
| Discovery (Key Discoveries / Scope Decisions) | Recommended for non-trivial/ambiguous requests | Assumptions tested, blind spots surfaced, what was in/out and why |
| Decision Record (Options Considered + Decision & Rationale) | Yes | Alternatives weighed and why the chosen one won (or why only one was viable) |
| Proposed Solution → Architecture | Yes | Key building blocks and relationships |
| Proposed Solution → Components | Yes | Components with responsibilities and interfaces |
| Proposed Solution → Data Flow | Yes | Step-by-step flow for the primary path |
| Error Handling | Yes | Top failure modes and how they are handled |
| Testing | Yes | Critical test cases and where they run |
| Open Questions | Yes | Only unresolved questions already surfaced; or "None" |

A missing **required** section is at least a Major finding, and a Blocker when the missing section is load-bearing for implementation (e.g., no Components, no Data Flow, no Error Handling).

## The 7 dimensions

For each finding, record: location (`design.md §section` and/or project file), evidence, recommendation, severity, and the dimension it belongs to.

### D1 — Completeness (完整性)

- Every required section above is present and **substantive** (not a placeholder or a restated heading).
- Architecture, Components, Data Flow, Error Handling, and Testing are all addressed.
- Every entry in Goals maps to something in the Proposed Solution; flag goals with no design support.
- Every Component has a stated responsibility and a defined interface (inputs/outputs, not just a name).
- Data Flow covers the primary happy path; for non-trivial designs, at least the key failure scenarios.
- Decision Record actually records the options weighed and the rationale — not empty, not a single line that skips the trade-off.
- Data models / schemas are defined where the solution depends on them.

### D2 — Usability / Actionability (可用性)

- Could `spec-plan` turn this into requirements and tasks without guessing? If a key decision is deferred to "implementation will decide", that is a finding.
- Interfaces, inputs/outputs, and data shapes are concrete enough to implement unambiguously.
- Data Flow is traceable end to end — no hand-wavy "then the system processes it" steps.
- Terminology is consistent and defined; no undefined jargon or acronyms.
- Success criteria in Goals are measurable, not aspirational ("fast", "scalable" without a number or condition).

### D3 — Document Conformance (规范性)

- Follows the canonical section set and keeps structural **headings in English**.
- Spec lives at a conventional path (`specs/<topic>/` or `.codex/specs/<topic>/`); `<topic>` is kebab-case.
- **Open Questions** contains only questions actually surfaced during discovery, or explicitly states none remain — not a dumping ground for unfinished design.
- **Decision Record** records alternatives that were genuinely considered; fabricated or boilerplate alternatives are a conformance finding, not a strength.
- Discovery section present when the request was non-trivial or ambiguous.

### D4 — Project Fit (符合项目规范)

Grounded in the project exploration (workflow step 3).

- Consistent with the project's tech stack, frameworks, and language conventions (CLAUDE.md / AGENTS.md rules, existing code patterns).
- Respects existing architecture and module boundaries; new components sit where the project would put them.
- **Reuses existing seams** instead of reinventing capability the project already has — flag duplication of an existing utility, service, or pattern.
- Does not silently conflict with or break existing behavior; if it changes existing behavior, the design says so.
- Naming, file layout, and data conventions match the project.

### D5 — Blind Spots (盲点)

Pick the checklists relevant to the design's domain (below) and surface gaps the design does not address. Frame each as "worth checking", and verify Open Questions captured the genuinely open unknowns rather than leaving them implicit.

### D6 — Over-Engineering (过度设计 / YAGNI)

Measured against the design's **own** Goals and Non-Goals.

- Features or components not traceable to any stated Goal — speculative "future use" work.
- Scope creep that contradicts a stated Non-Goal.
- Premature abstraction: layers, plugin systems, or generic frameworks for a single concrete use.
- Defensive bloat: handling for logically impossible cases, configurability nobody asked for.
- Complexity disproportionate to the problem — a simpler structure would meet every Goal.
- For each, point to the simpler alternative that still satisfies the Goals.

### D7 — Optimization (优化点)

- A materially simpler approach reaches the same Goals.
- More reuse of existing project components reduces new surface area.
- Coupling can be reduced without changing unrelated behavior.
- A pattern already proven in the project (or the domain) fits better than the proposed one.
- These are improvements, not blockers — usually Minor unless they overlap a real risk.

## Blind-spot checklists

Select by domain; skip items that clearly do not apply. Do not walk every list mechanically.

**Any user-facing feature**
- Empty / loading / error states
- Permissions and access control
- Undo or rollback
- Offline behavior
- Accessibility
- Internationalization / localization
- Mobile or responsive behavior

**Any data feature**
- Data migration from existing state
- Consistency and conflict resolution
- Concurrency and idempotency
- Retention and cleanup policies
- Privacy and compliance (GDPR, PII handling)
- Backup and recovery

**Any integration**
- Rate limits and quotas
- Authentication and credential rotation
- Failure modes and fallback / degradation
- Versioning and backwards compatibility
- Timeouts, retries, and replay/duplicate handling
- Monitoring and alerting

## Over-engineering signals (quick scan)

- A component whose only justification is "we might need it later".
- An abstraction with exactly one implementation and no second caller in sight.
- Configuration, flags, or extension points with no requirement behind them.
- Generic machinery (queues, caches, schedulers) where a direct call would do.
- A Non-Goal being quietly designed for anyway.

## Severity definitions

| Severity | Meaning | Typical examples |
|----------|---------|------------------|
| **Blocker** | Stops implementation, is internally contradictory, omits a load-bearing required section, or conflicts with existing project behavior/constraints. Must be fixed before `spec-plan`. | No Components/Data Flow/Error Handling; a Goal with no design at all; a design that contradicts how the existing system works |
| **Major** | A real gap or risk on a critical path that should be resolved, but does not by itself block all progress. | Unaddressed blind spot on a primary flow; ambiguous interface; over-engineering that materially inflates scope; empty Decision Record |
| **Minor** | Polish or optional improvement. | Wording, a measurable success criterion to tighten, a non-critical optimization or reuse opportunity |

## Verdict & spec-plan go/no-go rule

- **Reject** → No-Go: one or more **Blocker** findings. Return to `brainstorming` to revise the design before planning.
- **Revise** → Conditional: no Blockers, but one or more **Major** findings. The author should address the Majors; proceeding to `spec-plan` is at their discretion and risk.
- **Pass** → Go: no Blockers and no Majors (Minors allowed). Ready for `spec-plan`.

## Reviewing discipline

- Cite evidence for every finding; if you cannot point to a location, it is not a finding.
- Distinguish a **defect** (something is wrong) from a **blind spot** (something is unaddressed and worth checking). Frame the latter as a consideration.
- Judge the design against its own stated Goals and Non-Goals, not against a different solution you prefer.
- Reserve Blocker for true blockers — a long list of inflated severities is less useful than a short, accurate one.
- Scale the review to the design's size; a small design gets a short report.
