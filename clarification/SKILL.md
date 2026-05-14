---
name: clarification
description: >
  Clarify requirements before implementing small code changes. Use when the user asks for a seemingly small fix, tweak, option, validation, script change, UI adjustment, or localized behavior change, but important details may be missing. Surface blind spots, make low-risk reasonable defaults explicit, pause for user confirmation before behavior-changing edits, then implement after confirmation or recommend brainstorming/spec work if the scope is larger than it first appears.
---

# Clarification

Use this skill to prevent rushed implementation of underspecified small changes by adding a lightweight confirmation checkpoint before editing code. The goal is not a full spec; it is just enough clarification to safely code or to recognize that the task is actually larger.

Match the user's language.

## When to use

Use for code-change requests that look small but need requirement clarification, such as:

- Adding a parameter, flag, validation, display rule, or config default
- Adjusting localized behavior in a script, API, UI, test, or component
- Fixing a bug where the expected behavior is not fully stated
- Requests like "quickly change", "just add", "support", "compatible with", or "don't affect existing behavior"

## Do not use

- Clearly large features, architecture changes, migrations, permission systems, or cross-system behavior changes; recommend `brainstorming` instead
- Requests where the user explicitly asks for design/spec planning; use the appropriate planning workflow

## Workflow

1. Restate the target in one sentence. Identify the likely touched behavior, files, commands, or user path.
2. Judge whether the task is small enough to clarify inline.
   - Usually small: touches a few files with localized behavior, no new architecture, no migration, clear verification path.
   - Likely not small: spans multiple subsystems, introduces new product rules, changes data models, affects permissions/security policy, requires migrations, redesigns public APIs, or has several competing approaches.
3. If scope is unclear, inspect minimal project context before asking: README/AGENTS/CLAUDE, relevant files, call sites, tests, or existing patterns. Do not perform broad exploration unless needed.
4. Scan for blind spots that could change the implementation:
   - Desired new behavior and preserved old behavior
   - Inputs, outputs, empty/null/error cases, and boundary values
   - Backward compatibility for API/config/CLI/data formats
   - Security, permissions, privacy, performance, or migration risk
   - UI copy, docs, examples, tests, and verification commands
   - Non-goals: what should intentionally not change
5. Ask only blocking or high-impact questions. Continue until the implementation is safe enough, not until every theoretical question is answered.
6. Use reasonable defaults for low-risk details, but treat them as proposed assumptions before coding. State them explicitly, for example: "I propose to assume X; please confirm before I implement." If a default could affect users, data, security, or public contracts, ask instead of assuming.
7. Before implementation, pause for confirmation. Small does not mean confirmed. For behavior-changing edits, summarize the proposed requirement, assumptions/defaults, likely touched areas, and validation plan, then wait for the user to confirm.
8. Decide the path:
   - If still a small change and the user confirms: summarize the confirmed scope and start implementation.
   - If the user rejects or requests changes: clarify concerns and return to step 5.
   - If it has grown into a larger change: stop before coding and recommend `brainstorming` for design/spec clarification. Explain the specific reason.
   - If blocked by missing information: ask the smallest useful question, preferably with options.
9. After implementation, report concisely: what changed, assumptions used, validation run or not run, and any remaining risk.

## Output style

Keep clarification lightweight. Prefer:

```markdown
My understanding: <one sentence>
Potential blind spots: <only the relevant ones>
Assumptions/defaults: <proposed assumptions, if any>
Question(s): <blocking questions, if any>
```

When no blocking questions remain but the user has not confirmed implementation:

```markdown
My understanding: <one sentence>
Scope: <what will change / what will not change>
Assumptions/defaults: <proposed assumptions>
Plan: <implementation plan in 1-3 bullets>
Validation: <checks/tests to run>
Please confirm, and I’ll start implementation.
```

When the user confirms and you proceed to code:

```markdown
Confirmed scope: <short summary>
Proceeding with: <implementation plan in 1-3 bullets>
```

When recommending `brainstorming`:

```markdown
This looks larger than a small clarification because <specific reasons>.
I recommend using `brainstorming` to define the design/spec before implementation.
```

After implementation:

```markdown
Changed: <files and behavior>
Assumptions used: <list if any>
Validation: <tests run or reason not run>
Remaining risks: <if any>
```

## Examples

### Small change, clarify then implement

User: "Add a --dry-run option to scripts/export.py."

Good response:
- Confirm whether dry-run should skip all writes and external calls.
- Propose preserving existing output format and exit codes unless contradicted.
- Pause with a short plan and ask for confirmation before editing.
- Implement after the user confirms.

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
- If a clear existing priority field exists, state it as a proposed assumption, share the short plan, and ask for confirmation before editing.

## Verification checklist

Before coding:

- [ ] The desired behavior is clear enough to implement
- [ ] Important blind spots were considered
- [ ] Assumptions/defaults were stated as proposals
- [ ] A short implementation and validation plan was shared
- [ ] The user confirmed implementation, or explicitly authorized immediate implementation, or the change only affects formatting/comments/typos without any logic change
- [ ] The task still qualifies as a small change, or `brainstorming` was recommended

After coding:

- [ ] Relevant tests/checks were run, or the reason they were not run is stated
- [ ] Summary includes changed files/behavior and remaining risks
