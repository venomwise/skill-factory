---
name: clarification
description: >
  Use this workflow skill to clarify underspecified small code changes before editing.
  Add a lightweight confirmation checkpoint for fixes, tweaks, validations, script changes,
  UI adjustments, and localized behavior changes. Route larger design or architecture work
  to brainstorming/spec planning.
---

# Clarification

Use this skill to prevent rushed implementation of underspecified small changes. The goal is not a full spec; it is just enough clarification to safely code or to recognize that the task is actually larger.

Match the user's language.

## Boundaries

Use for code-change requests that look small but have missing or high-impact details, such as:

- Adding a parameter, flag, validation, display rule, or config default
- Adjusting localized behavior in a script, API, UI, test, or component
- Fixing a bug where the expected behavior is not fully stated
- Requests like "quickly change", "just add", "support", "compatible with", or "don't affect existing behavior"

Route clearly large features, architecture changes, migrations, permission systems, cross-system behavior changes, or explicit design/spec requests to the appropriate planning workflow.

## Workflow

1. Restate the target in one sentence and identify the likely touched behavior, files, commands, or user path.
2. Decide whether this is still a small inline clarification.
   - Usually small: touches a few files with localized behavior, no new architecture, no migration, clear verification path.
   - Likely not small: spans multiple subsystems, introduces new product rules, changes data models, affects permissions/security policy, requires migrations, redesigns public APIs, or has several competing approaches.
3. Inspect minimal project context before asking. Prefer README/AGENTS/CLAUDE, relevant files, call sites, tests, existing patterns, and targeted tool output. Do not ask about facts discoverable from the codebase. Do not perform broad exploration unless needed.
   - If the uncertainty depends on actual database schema or stored records, use `db-explorer` first to inspect read-only facts, then present findings as proposed assumptions before coding.
4. Surface only blind spots that could change implementation:
   - Desired new behavior and preserved old behavior
   - Inputs, outputs, empty/null/error cases, and boundary values
   - Backward compatibility for API/config/CLI/data formats
   - Security, permissions, privacy, performance, or migration risk
   - UI copy, docs, examples, tests, and verification commands
   - Non-goals: what should intentionally not change
5. Ask only blocking or high-impact questions. Use reasonable defaults for low-risk details, but state them as proposed assumptions. Do not block on low-risk formatting, naming, copy, or test-placement choices when existing project patterns make a safe default obvious.
6. Before implementation, pause for user confirmation. Small does not mean confirmed. Summarize the proposed requirement, assumptions/defaults, likely touched areas, and validation plan; implement only after the user confirms. If the task grows beyond a small change, stop and recommend `brainstorming` or spec planning with the specific reason.

## Output style

Keep clarification lightweight:

```markdown
My understanding: <one sentence>
Relevant blind spots: <only the blocking/high-impact ones>
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
Please confirm, and I will start implementation.
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

- `Add --dry-run to scripts/export.py`: inspect current writes/output, ask only whether dry-run skips external calls if unclear, then wait for confirmation.
- `Sort by business priority`: inspect current sort logic first; ask only if priority is not discoverable.
- `Add role-based access`: route to `brainstorming` because permissions are not a local tweak.

## Verification checklist

- [ ] Discoverable facts were inspected before asking.
- [ ] Only blocking/high-impact questions remain.
- [ ] Low-risk defaults are stated as proposed assumptions.
- [ ] The user confirmed before implementation.
- [ ] Relevant tests/checks were run, or the reason they were not run is stated.
