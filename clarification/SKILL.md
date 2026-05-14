---
name: clarification
description: >
  Clarify requirements before implementing small code changes. Use when the user asks for a seemingly small fix, tweak, option, validation, script change, UI adjustment, or localized behavior change, but important details may be missing. Surface blind spots, make low-risk reasonable defaults explicit, then either implement directly if the change is small or recommend brainstorming/spec work if the scope is larger than it first appears.
---

# Clarification

Use this skill to prevent rushed implementation of underspecified small changes. The goal is not a full spec; it is just enough clarification to safely code or to recognize that the task is actually larger.

Match the user's language.

## When to use

Use for code-change requests that look small but need requirement clarification, such as:

- Adding a parameter, flag, validation, display rule, or config default
- Adjusting localized behavior in a script, API, UI, test, or component
- Fixing a bug where the expected behavior is not fully stated
- Requests like "quickly change", "just add", "support", "compatible with", or "don't affect existing behavior"

## Do not use

- Purely mechanical tasks with no meaningful ambiguity, such as typo fixes or formatting-only edits
- Clearly large features, architecture changes, migrations, permission systems, or cross-system behavior changes; recommend `brainstorming` instead
- Requests where the user explicitly asks for design/spec planning; use the appropriate planning workflow

## Workflow

1. Restate the target in one sentence. Identify the likely touched behavior, files, commands, or user path.
2. Judge whether the task is small enough to clarify inline.
   - Usually small: 1-3 files, localized behavior, no new architecture, no migration, clear verification path.
   - Likely not small: multiple subsystems, new product rules, data model changes, permissions/security policy, migrations, public API redesign, or several competing approaches.
3. If scope is unclear, inspect minimal project context before asking: README/AGENTS/CLAUDE, relevant files, call sites, tests, or existing patterns. Do not perform broad exploration unless needed.
4. Scan for blind spots that could change the implementation:
   - Desired new behavior and preserved old behavior
   - Inputs, outputs, empty/null/error cases, and boundary values
   - Backward compatibility for API/config/CLI/data formats
   - Security, permissions, privacy, performance, or migration risk
   - UI copy, docs, examples, tests, and verification commands
   - Non-goals: what should intentionally not change
5. Ask only blocking or high-impact questions. Continue until the implementation is safe enough, not until every theoretical question is answered.
6. Use reasonable defaults for low-risk details. State them explicitly, for example: "I will assume X; if that's wrong, tell me before I continue." If a default could affect users, data, security, or public contracts, ask instead of assuming.
7. Decide the path:
   - If still a small change: summarize the confirmed requirement, assumptions/defaults, and start implementation.
   - If it has grown into a larger change: stop before coding and recommend `brainstorming` for design/spec clarification. Explain the specific reason.
   - If blocked by missing information: ask the smallest useful question, preferably with options.
8. After implementation, report concisely: what changed, assumptions used, validation run or not run, and any remaining risk.

## Output style

Keep clarification lightweight. Prefer:

```markdown
My understanding: <one sentence>
Potential blind spots: <only the relevant ones>
Assumptions/defaults: <if any>
Question(s): <blocking questions, or say you can proceed>
```

When proceeding to code:

```markdown
Confirmed scope: <short summary>
I will proceed with: <implementation plan in 1-3 bullets>
```

When recommending `brainstorming`:

```markdown
This looks larger than a small clarification because <specific reasons>.
I recommend using `brainstorming` to define the design/spec before implementation.
```

## Examples

### Small change, clarify then implement

User: "Add a --dry-run option to scripts/export.py."

Good response:
- Confirm whether dry-run should skip all writes and external calls.
- Default to preserving existing output format and exit codes unless contradicted.
- Implement once the behavior is clear.

### Hidden large change, route away

User: "Just add role-based access so admins and users see different pages, and APIs are protected too."

Good response:
- Identify roles, authorization boundaries, frontend/backend enforcement, tests, and security implications.
- Do not start coding.
- Recommend `brainstorming` because it is a permission model, not a local tweak.

### Ambiguous bug, inspect first

User: "The list sorting is wrong; sort by business priority."

Good response:
- Inspect current sort logic and any existing priority fields/constants.
- Ask what defines business priority if not discoverable.
- If a clear existing priority field exists, state the assumption and implement.

## Verification checklist

Before coding:

- [ ] The desired behavior is clear enough to implement
- [ ] Important blind spots were considered
- [ ] Assumptions/defaults were stated
- [ ] The task still qualifies as a small change, or `brainstorming` was recommended

After coding:

- [ ] Relevant tests/checks were run, or the reason they were not run is stated
- [ ] Summary includes changed files/behavior and remaining risks
