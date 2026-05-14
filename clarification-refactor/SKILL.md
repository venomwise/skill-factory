---
name: clarification-refactor
description: >
  Route and clarify coding requests before implementation. Use when a user asks for a small code change, ambiguous bug fix, localized tweak, CLI/API/UI adjustment, or "quick/simple" change. First decide whether to implement, ask focused blind-spot questions, inspect context, or hand off to brainstorming. For user-visible behavior changes such as ordering, filtering, validation, output format, errors, permissions, or API behavior, ask before coding. For high-risk or cross-cutting work such as authorization, security, migrations, billing, public contracts, product policy, or multiple subsystems, do not code; explicitly recommend brainstorming before implementation.
---

# Clarification Refactor

This skill is a router before it is an executor. Its first job is to decide whether to:

1. implement a small local change,
2. ask focused blind-spot questions,
3. inspect limited context, or
4. hand off to `brainstorming`.

Match the user's language.

## First response router

Before editing files, classify the request:

### IMPLEMENT

Use when the change is local, low-risk, and no ask-before-code category is unresolved.

- State the understood behavior and low-risk defaults briefly.
- Implement, validate, and summarize.

### ASK

Use when the task may be small, but a user-visible or behavior-changing decision is unresolved.

- Inspect relevant context first if the answer may already be in code/docs/tests.
- Ask the smallest blocking question.
- Stop and wait. Do not edit files in the same turn.

### INSPECT

Use when the answer is likely in code, tests, config, docs, or nearby patterns.

- Inspect only enough context to route correctly.
- Then transition to IMPLEMENT, ASK, or BRAINSTORM.

### BRAINSTORM

Use when the request exceeds a localized code edit.

- Do not edit files.
- Produce the `brainstorming` handoff packet and stop.

## Decision priority

When rules compete, use this priority order:

1. User/project instructions and safety constraints
2. Safety, security, data, privacy, migration, and public-contract risk
3. Scope escalation to `brainstorming`
4. Ask-before-code for user-visible behavior transformations
5. Reasonable defaults for low-risk implementation details
6. Implementation

Reasonable defaults must not override scope escalation or ask-before-code categories.

## Scope escalation → brainstorming

Use `BRAINSTORM` when the request is no longer a localized code edit.

Scope escalation includes:

- New or changed product/business policy
- Security, privacy, permissions, authorization, access boundaries, audit behavior, or compliance
- Data model changes, migrations, billing/payment, irreversible data operations, or ownership rules
- Public API contracts, compatibility guarantees, or cross-version behavior
- Coordinated frontend/backend behavior
- Multiple subsystems, user flows, deployment surfaces, or teams
- Competing implementation approaches or architectural choices
- Missing project structure where implementation would require inventing architecture

When scope escalation is detected:

1. Do not edit files.
2. Do not ask only for missing framework/files as the next step.
3. Use this exact bridge sentence: "I recommend using `brainstorming` before implementation."
4. Provide a handoff packet:
   - Current goal
   - Why this exceeds a small change
   - Key unknowns / decisions
   - Known project context
   - Suggested first brainstorming question

`BRAINSTORM` is terminal for this skill. After the handoff packet, stop and wait.

### Missing information handling

- Missing implementation detail: flag name, output wording, config default, or a single local behavior detail. Ask a focused question or use a low-risk default.
- Missing design context: role model, permission matrix, data ownership, API contract, migration strategy, cross-system behavior, or product policy. Bridge to `brainstorming`.
- Missing code location: if the change is clearly small, ask for file/path. If the request also has scope-escalation signals, bridge to `brainstorming` and include missing code location in the handoff packet.

## Ask-before-code categories

Use `ASK` before coding when the change affects user-visible behavior transformations, including:

- ordering, filtering, grouping, ranking, pagination, or search behavior
- validation, error behavior, empty/null/boundary handling, or output format
- permissions, access, visibility, privacy, or public API behavior
- data deletion/update semantics or irreversible actions

Ask about the smallest relevant subset of:

- precedence: which rule wins when multiple rules apply?
- tie-breakers: what happens when values are equal?
- fallback: what happens for unknown/missing values?
- compatibility: should existing behavior/output/API stay compatible?
- preservation: which existing behavior should intentionally remain unchanged?

Do not implement from assumptions in the same turn unless the user already confirmed the behavior or explicitly authorized assumptions.

## Reasonable defaults

Use defaults only for low-risk implementation details. Prefer defaults that:

- preserve existing behavior
- minimize scope
- follow nearby project patterns
- remain backward compatible
- are easy to validate or revert

State defaults explicitly when proceeding.

## Implementation workflow

Only use this after routing to `IMPLEMENT`.

1. Summarize confirmed behavior and defaults in 1-3 bullets if helpful.
2. Edit the smallest relevant code surface.
3. Run focused validation: tests, syntax checks, targeted command, or explain why not run.
4. Summarize changed behavior, files, validation, and remaining risk.

## Output patterns

### Small local change

```markdown
My understanding: <one sentence>
Assumptions/defaults: <low-risk defaults>
Proceeding with: <1-3 bullets, optional>
```

### Ask-before-code

```markdown
I found <relevant existing behavior/context>. Before editing, one behavior-changing detail needs confirmation:

<focused question>

Suggested default: <default>, because <reason>. I will wait because this changes user-visible behavior.
```

### Brainstorming handoff

```markdown
This is larger than a small clarification because <specific reasons>.

I recommend using `brainstorming` before implementation.

Handoff context:
- Current goal: <summary>
- Why this exceeds a small change: <bullets>
- Key unknowns / decisions: <bullets>
- Known project context: <what was inspected or what is missing>
- Suggested first brainstorming question: <one focused question>
```

## Bridge anti-patterns

Avoid replacing scope escalation with:

- "Which framework should I use?"
- "Please provide the app files."
- "Once you send the repo, I can add permissions."
- "I can scaffold a generic solution."

Prefer a `brainstorming` handoff when the missing information is design context or architecture.

## Examples

### Local CLI flag

User: "Add `--dry-run` to `scripts/export.py`."

Good: inspect the script, state low-risk defaults such as no writes/no output directory creation, implement, validate normal and dry-run behavior.

### Cross-cutting policy/security change

User: "Quickly add role-based access so different roles see different pages and APIs are protected."

Good: treat as scope escalation because role model, enforcement points, frontend/backend consistency, API behavior, and security boundaries need design. Use the `brainstorming` handoff packet. Do not ask only for framework/files.

### User-visible ordering change

User: "Sort this list by business priority."

Good: inspect existing sort logic and priority definitions if available, then ask about precedence, tie-breakers, fallback for unknown values, and which existing ordering behavior should be preserved. Wait before editing.
