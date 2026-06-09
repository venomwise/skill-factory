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
4. Brainstorm with the user through iterative rounds. Do not pre-commit to how many questions you will ask or to a single topic to ask about before exploration (steps 1-3) is complete — the number and subject of questions are determined by the gaps and blind spots you actually surface, not estimated up front. First diagnose the fuzziness level of the request (see the guide for criteria), then apply matched techniques:
   - Problem unclear -> reframe: help the user articulate what problem they are actually solving, not just what feature they want.
   - Direction unclear -> explore: sketch a few representative directions before narrowing, using techniques like "what if" scenarios and perspective switching.
   - Boundaries unclear -> scan: use the domain-relevant blind-spot checklist from the guide to surface potential gaps for the user to consider.
   - Solution unclear -> the technical approach still needs comparison, but first confirm intent, priorities, and constraints below before moving to step 5. A request full of specific technical terms does NOT mean intent is confirmed; familiarity with the technology is not knowledge of what the user wants.
   For non-trivial or ambiguous requests, surface at least one assumption worth confirming and one potential blind spot for the user to consider.
   Ask one question per turn and wait for the answer before asking the next — this is about cadence, not a cap on how many questions you ask in total. Prefer multiple-choice when it helps the user choose between meaningful options. Do not ask a question and then answer it yourself with an assumption in the same turn. Treat unknown intent, priorities, trade-off preferences, and acceptance criteria as things to ask about, not to assume — these live only in the user's head and cannot be retrieved by inspecting the project.
   Continue asking until the exit condition is met. Exit condition: the user can clearly state what they want, why they want it, and what they explicitly do not want; key constraints and success criteria are known or explicitly recorded as assumptions; and the problem statement has been stable for the last 2 exchanges. Before leaving this step, verify: can you summarize the user's intent, priorities, and constraints in 2-3 sentences using their own words? If not, ask more questions.
5. Converge and propose. Admission gate — check all three conditions before entering:
   1. Have you asked the user at least one clarifying question in this session? If no, go back to step 4.
   2. Can the user now state what they want, why they want it, and what they explicitly do not want? If no, go back to step 4 and ask more questions.
   3. Are key constraints and success criteria known or explicitly recorded as assumptions? If no, go back to step 4 and ask more questions.
   Note: A request being technically detailed does not satisfy these conditions. Technical detail says nothing about the user's intent, priorities, or trade-off preferences.
   Once all conditions are met, summarize what brainstorming revealed: the refined problem statement, request-fit findings, challenged assumptions, discovered blind spots, and trimmed scope. Then propose 1-3 approaches with trade-offs. Lead with your recommendation. If only one approach is viable, explain why alternatives were ruled out. Ask the user to confirm the problem summary and select an approach before continuing.
6. Present the design scaled to complexity. For simple projects, present the full design at once and ask for approval. For moderate or complex projects, present by section and ask for approval after each. Cover architecture, components, data flow, error handling, and testing.
7. Write the design doc to `specs/<topic>/design.md` using the template in `assets/design-doc-template.md`. Preserve the template's English headings and structural labels exactly, while writing the section content in the user's current language. Infer that language from the conversation; do not ask solely to determine it. Name `<topic>` using kebab-case derived from the project or feature name (e.g., `user-auth`, `payment-integration`). Confirm the path with the user if ambiguous. Fill the Decision Record section from the approach comparison produced in step 5: record each option that was actually weighed with its key trade-offs, then the chosen approach and the concrete reasons it won (or, if only one approach was viable, why the alternatives were ruled out). This comparison was generated live in step 5 and is the main thing worth preserving for later review, so do not drop it; if it has scrolled out of the recent conversation, go back and recover it rather than reconstructing from memory. Do not fabricate alternatives that were never discussed. Open Questions must contain only unresolved questions already surfaced to the user; if none remain, write that no open questions remain.
8. User review gate. Ask the user to review the written doc. On feedback:
   - Wording or detail changes: edit the doc and re-confirm.
   - Scope or approach changes: return to step 5.
   - Missing context: return to step 4.
   Proceed only after the user approves the written document.
9. Invoke `spec-plan` as the only next step. Pass the following context: project name (`<topic>`), target directory (`specs/<topic>/`), scope summary (Summary + Non-Goals + Discovery / Scope Decisions if present), constraints (Context + Discovery / Key Discoveries when they contain confirmed constraints, risks, or assumptions), and primary users and goals (Primary Users / Roles + Goals). If a Discovery section exists, treat it as input context for requirements and task planning, especially for confirmed assumptions, surfaced risks, and explicit scope decisions. Do not invoke any implementation skill.

## Verification

- [ ] `specs/<topic>/design.md` exists and the user approved it
- [ ] The design doc includes the template headings from `assets/design-doc-template.md`
- [ ] The design doc keeps template headings in English and writes section content in the user's current language
- [ ] The design covers architecture, components, data flow, error handling, and testing
- [ ] If brainstorming surfaced notable discoveries, they are recorded in the design doc's Discovery section
- [ ] The Decision Record captures the approaches compared in step 5 and the rationale for the chosen one (or why alternatives were ruled out when only one was viable)
- [ ] Open Questions contains only surfaced unresolved questions, or explicitly says none remain
- [ ] `spec-plan` has been invoked as the terminal action

## Safety & guardrails

- No implementation before approval. Do not write code, scaffold projects, or invoke any implementation skill until the design is approved.
- Even simple projects require a design; keep it short when the scope is small.
- Scale the process to complexity. Simple projects may complete steps 4-6 in a few exchanges; do not pad the process with unnecessary ceremony.
- Pace questions correctly. Ask one question per turn and wait for the answer before the next (cadence), but keep going until intent and constraints are clear — the number of questions is set by the gaps you find, not minimized for its own sake. Full rules live in step 4.
- YAGNI at the right time. While diverging (surfacing options and blind spots), explore freely and do not dismiss ideas prematurely. While converging and designing, cut ruthlessly - remove anything that is not essential to the core goal.
- Design for isolation. Break the system into units with one clear purpose and well-defined interfaces. See [the guide](references/brainstorming-guide.md) for details.

## References

- [Detailed brainstorming guide](references/brainstorming-guide.md)
- [Design doc template](assets/design-doc-template.md)
