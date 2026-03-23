---
name: ad-spec-plan
description: Create requirements.md and tasks.md for a project spec. Use when scoping a feature, module, service, or integration and you need requirements plus an execution plan with status tracking.
# metadata:
#   short-description: Spec plan (requirements + tasks)
---

# Spec Plan Skill

## When to use

- Create or refresh a project spec with requirements and an execution plan.
- Turn a vague request into clear requirements and actionable tasks.
- Establish traceability from tasks to requirements.

## When not to use

- You only need a quick TODO list or a single file edit.
- The work is already captured in an existing spec and only needs implementation.

## Inputs

- Project name and target directory.
- Scope summary (what is in and out of scope).
- Primary user roles and key goals.
- Constraints (platforms, dependencies, timelines).

## Outputs

- `requirements.md`
- `tasks.md`

## Workflow

1. Confirm the target directory and project name. If the user did not specify a path, suggest `.codex/specs/<project-name>/` and confirm.
2. Draft `requirements.md` using `assets/requirements.template.md`.
   Follow the HTML comments in the template for content depth and coverage guidance.
   HTML comments are authoring instructions - do NOT include them in the final output.
3. Draft `tasks.md` using `assets/tasks.template.md` and link to `requirements.md`.
   Include test tasks for each functional phase and add Checkpoint stages at key milestones.
   Mark optional phases/tasks with `- [ ]*` (asterisk immediately after closing bracket). Test tasks, verification tasks, and nice-to-have features MUST use this marker — never use text labels like "可选" or "Optional:" instead.
4. Ensure every task references one or more requirement IDs (for traceability).
5. If any requirements are ambiguous or missing, ask the user before finalizing.

## Verification (self-check before finalizing)

Before presenting the final output, scan both files and confirm each item:

- [ ] `requirements.md` contains Introduction, Glossary, and numbered Requirements sections.
- [ ] `tasks.md` links to `requirements.md` and every task includes a `_Requirements: ..._` line.
- [ ] Requirement IDs referenced in `tasks.md` exist in `requirements.md`.
- [ ] Each requirement includes acceptance criteria covering normal flow, error flow, and boundary conditions.
- [ ] Tasks include specific file paths and function/class names.
- [ ] At least one Checkpoint task exists between major phases.
- [ ] Phase headings use `- [ ] N. Phase N:` checkbox format, not markdown headings (`###`).
- [ ] Task descriptions are indented bullet points under the task title line.
- [ ] Every test task, verification task, and nice-to-have feature uses `- [ ]*` marker. No task uses text labels like "可选", "Optional", or "(Optional)" as a substitute.
- [ ] Sub-tasks under an optional Phase inherit optionality and do not need their own `*`.
- [ ] Requirements references use `N.M` format (e.g., `_Requirements: 1.1, 2.3_`), not `RN` format.

## Safety & guardrails

- Do not invent requirements without confirmation; mark assumptions explicitly.
- Keep requirements testable and phrased as acceptance criteria.
- Do not mark tasks as complete unless the work is done.

## References

- [Requirements template](assets/requirements.template.md)
- [Tasks template](assets/tasks.template.md)
