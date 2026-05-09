---
name: spec-plan
description: Create requirements.md and tasks.md for a project spec, with traceability and status tracking. Use when the user asks for a spec, requirements, task plan, or execution breakdown (e.g. "做个 spec", "拆任务", "出执行计划", "write a spec", "plan tasks"). Requires an approved specs/<topic>/design.md as input; if missing, route to the brainstorming skill first.
# metadata:
#   short-description: Spec plan (requirements + tasks)
---

# Spec Plan Skill

## When to use

- Create or refresh `requirements.md` and `tasks.md` from an approved `design.md`.
- Convert the design into accurate, testable acceptance criteria and actionable tasks.
- Establish traceability from tasks to requirements, and from requirements back to the approved design.

## When not to use

- You only need a quick TODO list or a single file edit.
- The work is already captured in an existing spec and only needs implementation.
- No approved `design.md` exists yet; create or request the design first instead of generating requirements/tasks from assumptions.

## Inputs

- Project name and target directory.
- Approved design document: `specs/<topic>/design.md`.
- Optional clarifications only when `design.md` is ambiguous, incomplete, or internally inconsistent.

## Outputs

- `requirements.md`
- `tasks.md`
- Final response includes a recommended next step for the user to run `spec-exec` on the generated `specs/<topic>/tasks.md`; do not start implementation automatically.

## Workflow

1. Resolve the target design document.
   - If the user provided `specs/<topic>/design.md`, use it.
   - If the user did not provide a design path, reply exactly: `请指定 design 文件路径（例如 specs/<topic>/design.md）。` Then end the workflow. Do not search, list, or infer a design.
   - If the user-provided path does not exist, end the workflow and recommend running the **`brainstorming`** skill first to produce an approved design at `specs/<topic>/design.md`.
2. Confirm the target directory and project name from the selected `specs/<topic>/design.md`.
3. Open the selected `design.md`.
4. Draft `requirements.md` using `assets/requirements.template.md`.
   Follow the HTML comments in the template for content depth and coverage guidance.
   HTML comments are authoring instructions - do NOT include them in the final output.
   Treat `design.md` as the authoritative source of requirements. Translate the approved design into accurate, testable acceptance criteria without adding, omitting, or changing intended behavior. Reflect approved behavior and constraints where applicable, but do not convert design rationale, examples, alternatives, or future ideas into hard requirements unless the design explicitly requires them.
5. Draft `tasks.md` using `assets/tasks.template.md` and link to `requirements.md`.
   Include test tasks for each functional phase and add Checkpoint stages at key milestones. Checkpoints are validation tasks for the execution agent, not user approval gates; write concrete verification steps and blocker conditions instead of asking whether to continue.
   Mark optional phases/tasks with `- [ ]*` (asterisk immediately after closing bracket). Non-essential steps such as test tasks, verification tasks, summary/documentation wrap-up tasks, and nice-to-have features MUST use this marker — never use text labels like "可选" or "Optional:" instead.
6. Ensure every task references one or more requirement IDs (for traceability).
7. Use chain validation before finalizing: verify `requirements.md` accurately implements `design.md` without drift, then verify `tasks.md` covers and references `requirements.md`.
8. If `design.md` is ambiguous, incomplete, or internally inconsistent, ask the user before finalizing. Do not fill gaps by inventing requirements.

## Verification (self-check before finalizing)

Before presenting the final output, scan both files and confirm each item:

- [ ] `design.md` exists and has been read before drafting `requirements.md` or `tasks.md`.
- [ ] `requirements.md` contains Introduction, Glossary, and numbered Requirements sections.
- [ ] Every requirement is traceable to approved behavior or constraints in `design.md`.
- [ ] `requirements.md` accurately implements `design.md`: no added behavior, omitted required behavior, changed semantics, or contradictory architecture.
- [ ] Important approved behaviors and constraints from `design.md` are reflected as testable acceptance criteria where applicable.
- [ ] Design rationale, examples, alternatives, and future ideas are not converted into hard requirements unless explicitly required by `design.md`.
- [ ] `tasks.md` links to `requirements.md` and every task includes a `_Requirements: ..._` line.
- [ ] Requirement IDs referenced in `tasks.md` exist in `requirements.md`.
- [ ] Each requirement includes acceptance criteria covering normal flow, error flow, and boundary conditions.
- [ ] Tasks include specific file paths and function/class names.
- [ ] At least one Checkpoint task exists between major phases and describes concrete validation steps, not user approval.
- [ ] Phase headings use `- [ ] N. Phase N:` checkbox format, not markdown headings (`###`).
- [ ] Task descriptions are indented bullet points under the task title line.
- [ ] Every non-essential step uses `- [ ]*` marker, including test tasks, verification tasks, summary/documentation wrap-up tasks, and nice-to-have features. No task uses text labels like "可选", "Optional", or "(Optional)" as a substitute.
- [ ] Sub-tasks under an optional Phase inherit optionality and do not need their own `*`.
- [ ] Requirements references use `N.M` format (e.g., `_Requirements: 1.1, 2.3_`), not `RN` format.

## Safety & guardrails

- Never generate `requirements.md` or `tasks.md` without an approved `design.md`.
- Never auto-select a `design.md`. The path must come from the user.
- Do not invent requirements, broaden scope, or reinterpret design intent. Ask for clarification when the design is ambiguous or incomplete.
- Keep requirements testable and phrased as acceptance criteria.
- Do not mark tasks as complete unless the work is done.

## References

- [Requirements template](assets/requirements.template.md)
- [Tasks template](assets/tasks.template.md)
