---
name: clarification
description: >
  Use this workflow skill to clarify underspecified small code changes before editing.
  Add a lightweight confirmation checkpoint for fixes, tweaks, validations, script changes,
  UI adjustments, and localized behavior changes. Route larger design or architecture work
  to brainstorming/spec planning.
---

# Clarification

Use this skill as a pre-edit triage for code-affecting requests. The goal is not a full spec or a confirmation loop; it is just enough clarification to code safely, proceed quickly when the request is already clear, or recognize that the task is actually larger.

Match the user's language.

## Boundaries

Trigger broadly for any request that may modify:

- Source code, tests, scripts, configs, build files, generated source, or docs that affect behavior
- Runtime behavior in a script, API, UI, CLI, test, job, or component
- Parameters, flags, validation, display rules, defaults, compatibility behavior, or error handling
- Requests like "change", "modify", "add", "remove", "fix", "support", "make compatible", "quickly change", or "don't affect existing behavior"

Route clearly large features, architecture changes, migrations, permission systems, cross-system behavior changes, or explicit design/spec requests to the appropriate planning workflow.

Do not use as an extra preflight inside `spec-exec`. When the user asks to implement, resume, or check an approved `specs/<topic>/tasks.md`, let `spec-exec` own clarification through its blocker escalation rules.

## Workflow

1. Restate the target in one sentence and identify the likely touched behavior, files, commands, or user path.
2. Triage the request.
   - Proceed directly when it is clear, local, low-risk, follows existing patterns, and has an obvious verification path.
   - Clarify first when missing details could change behavior, compatibility, data handling, security, permissions, UX, API/config/CLI contracts, or test expectations.
   - Route to planning when it spans multiple subsystems, introduces new product rules, changes data models, affects permissions/security policy, requires migrations, redesigns public APIs, or has several competing approaches.
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
6. Before implementation, choose the matching path:
   - Clear and low-risk: state the understood scope and proceed without waiting for extra confirmation.
   - Blocking ambiguity or high-impact defaults: summarize the proposed requirement, assumptions/defaults, likely touched areas, and validation plan; implement only after the user confirms.
   - Larger than a local change: stop and recommend `brainstorming` or spec planning with the specific reason.

## Output style

Keep clarification lightweight:

When the request is clear and low-risk:

```markdown
My understanding: <one sentence>
Scope: <what will change / what will not change>
Proceeding with: <implementation plan in 1-3 bullets>
Validation: <checks/tests to run>
```

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

- `Rename a button label in one component`: inspect the component, state the local scope, then proceed.
- `Add --dry-run to scripts/export.py`: inspect current writes/output, ask only whether dry-run skips external calls if unclear.
- `Sort by business priority`: inspect current sort logic first; ask only if priority is not discoverable.
- `Add role-based access`: route to `brainstorming` because permissions are not a local tweak.
- `Run spec-exec on specs/auth/tasks.md`: do not use this skill; `spec-exec` handles blockers.

## Verification checklist

- [ ] Discoverable facts were inspected before asking.
- [ ] Approved `spec-exec` runs were left to `spec-exec` blocker escalation instead of re-clarified here.
- [ ] Clear, local, low-risk requests proceeded without unnecessary confirmation.
- [ ] Only blocking/high-impact questions remain.
- [ ] Low-risk defaults are stated as proposed assumptions.
- [ ] The user confirmed before implementation when blocking ambiguity or high-impact defaults existed.
- [ ] Relevant tests/checks were run, or the reason they were not run is stated.
