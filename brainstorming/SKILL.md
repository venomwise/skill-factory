---
name: brainstorming
description: Turn ideas into validated designs and specs through collaborative dialogue. Use before creating features, components, or behavior changes when requirements need scoping and trade-offs. Not for bug fixes, typo-only changes, or clear single-step execution tasks.
---

# Brainstorming Ideas Into Designs

## When to use

- Creating new features or components
- Modifying existing system behavior
- Requirements are unclear and need scoping
- Multiple implementation approaches are possible

## When not to use

- Bug fixes with clear root cause and steps
- Typo or formatting-only changes
- Clear, single-step execution tasks
- The user explicitly declines the design process

## Inputs

- User's idea or goal, possibly vague
- Existing project context (files, docs, recent commits)

## Outputs

- Validated design doc at `specs/<topic>/design.md`
- Confirmed next step to use `spec-plan`

## Workflow

1. Explore project context in priority order: README/CLAUDE/AGENTS.md, project config (package.json / pyproject.toml / Cargo.toml / go.mod / pom.xml), entry points, then recent commits (up to 10 if needed). Stop when you can describe the project's purpose, tech stack, and structure in 2-3 sentences. If the project is complex or unfamiliar, consult [the guide](references/brainstorming-guide.md) for exploration priorities.
2. Validate request fit against the current project before accepting the requested shape. Inspect the existing behavior, data, and interfaces the request would touch, grounding the check in real evidence rather than assumptions — use available tools such as `db-explorer` when the request depends on stored data or schema. If the request looks redundant, conflicts with what exists, solves the wrong problem, or has a simpler alternative, pause, tell the user the finding and your recommended path, then ask whether to proceed.
3. Assess scope. If the request spans multiple independent subsystems, propose a decomposition with sub-projects, dependencies, and build order, then ask the user to confirm before proceeding. Brainstorm only the first confirmed sub-project; each gets its own spec cycle.
4. Brainstorm with the user. First diagnose the fuzziness level of the request (see the guide for criteria), then apply matched techniques:
   - Problem unclear -> reframe: help the user articulate what problem they are actually solving, not just what feature they want.
   - Direction unclear -> explore: sketch a few representative directions before narrowing, using techniques like "what if" scenarios and perspective switching.
   - Boundaries unclear -> scan: use the domain-relevant blind-spot checklist from the guide to surface potential gaps for the user to consider.
   - Solution unclear -> skip to step 5 for approach comparison.
   For non-trivial or ambiguous requests, surface at least one assumption worth confirming and one potential blind spot for the user to consider.
   Ask questions one at a time. Prefer multiple-choice when it helps the user choose between meaningful options.
   Stop when the user can clearly state what they want, why they want it, and what they explicitly do not want, and when key constraints and success criteria are known or explicitly recorded as assumptions.
5. Converge and propose. Summarize what brainstorming revealed: the refined problem statement, request-fit findings, challenged assumptions, discovered blind spots, and trimmed scope. Then propose 1-3 approaches with trade-offs. Lead with your recommendation. If only one approach is viable, explain why alternatives were ruled out. Ask the user to confirm the problem summary and select an approach before continuing.
6. Present the design scaled to complexity. For simple projects, present the full design at once and ask for approval. For moderate or complex projects, present by section and ask for approval after each. Cover architecture, components, data flow, error handling, and testing.
7. Write the design doc to `specs/<topic>/design.md` using the template in `assets/design-doc-template.md`. Name `<topic>` using kebab-case derived from the project or feature name (e.g., `user-auth`, `payment-integration`). Confirm the path with the user if ambiguous. Open Questions must contain only unresolved questions already surfaced to the user; if none remain, write that no open questions remain.
8. User review gate. Ask the user to review the written doc. On feedback:
   - Wording or detail changes: edit the doc and re-confirm.
   - Scope or approach changes: return to step 5.
   - Missing context: return to step 4.
   Proceed only after the user approves the written document.
9. Invoke `spec-plan` as the only next step. Pass the following context: project name (`<topic>`), target directory (`specs/<topic>/`), scope summary (Summary + Non-Goals + Discovery / Scope Decisions if present), constraints (Context + Discovery / Key Discoveries when they contain confirmed constraints, risks, or assumptions), and primary users and goals (Primary Users / Roles + Goals). If a Discovery section exists, treat it as input context for requirements and task planning, especially for confirmed assumptions, surfaced risks, and explicit scope decisions. Do not invoke any implementation skill.

## Verification

- [ ] `specs/<topic>/design.md` exists and the user approved it
- [ ] The design doc includes the template headings from `assets/design-doc-template.md`
- [ ] The design covers architecture, components, data flow, error handling, and testing
- [ ] If brainstorming surfaced notable discoveries, they are recorded in the design doc's Discovery section
- [ ] Open Questions contains only surfaced unresolved questions, or explicitly says none remain
- [ ] `spec-plan` has been invoked as the terminal action

## Safety & guardrails

- No implementation before approval. Do not write code, scaffold projects, or invoke any implementation skill until the design is approved.
- Even simple projects require a design; keep it short when the scope is small.
- Scale the process to complexity. Simple projects may complete steps 4-6 in a few exchanges; do not pad the process with unnecessary ceremony.
- One question at a time. Do not overwhelm the user with multiple questions.
- YAGNI at the right time. During brainstorming (step 4), explore freely and do not dismiss ideas prematurely. During convergence and design (steps 5+), cut ruthlessly - remove anything that is not essential to the core goal.
- Design for isolation. Break the system into units with one clear purpose and well-defined interfaces. See [the guide](references/brainstorming-guide.md) for details.

## References

- [Detailed brainstorming guide](references/brainstorming-guide.md)
- [Design doc template](assets/design-doc-template.md)
