---
name: brainstorming
description: Turn ideas into validated design docs through collaborative dialogue before implementation. Use this skill whenever the user wants to brainstorm, design, or scope a new feature, component, or behavior change — especially when requirements are fuzzy, multiple approaches exist, or the user says things like "I want to build...", "how should we design...", "let's think through...", "I have an idea for...", "help me plan...", "我想做一个...", "帮我规划...", "一起想想...", "这个该怎么设计...", or "先别写代码，先把方案想清楚". Also trigger when the user describes a goal without a clear path, asks for architecture advice on something new, or wants to explore trade-offs before committing. Not for bug fixes, typo-only changes, or clear single-step tasks.
---

# Brainstorming Ideas Into Designs

Help the user go from a rough idea to a validated design document, then hand off to `spec-plan` for requirements and task breakdown. The goal is to reduce ambiguity and surface blind spots *before* any code gets written — because changing a design doc is cheap, changing an implementation is not.

## When to use

- Creating new features or components
- Modifying existing system behavior in non-trivial ways
- Requirements are unclear and need scoping
- Multiple implementation approaches are possible
- The user wants to think through trade-offs before committing

## When not to use

- Bug fixes with clear root cause
- Typo or formatting-only changes
- Clear, single-step execution tasks
- The user explicitly declines the design process

## Inputs

- User's idea or goal (possibly vague)
- Existing project context (codebase, docs, recent commits)

## Outputs

- Validated design doc at `specs/<topic>/design.md` (or project-specific location)
- Handoff to `spec-plan` with structured context

## Workflow

### 1. Explore project context

Understand the project before asking the user anything — it shows respect for their time and lets you ask smarter questions.

Use tools in this priority order, stopping as soon as you can describe the project's purpose, tech stack, and structure in 2-3 sentences:

1. README, CLAUDE.md, AGENTS.md — project-level docs
2. Project config (package.json / pyproject.toml / Cargo.toml / go.mod / pom.xml)
3. Entry points and directory structure (use Glob to scan key patterns)
4. Recent commits (up to 10) for current momentum and conventions

For complex or unfamiliar projects, consult [the brainstorming guide](references/brainstorming-guide.md) for deeper exploration tactics.

### 2. Assess scope

If the request spans multiple independent subsystems, propose a decomposition:
- List sub-projects with dependencies and build order
- Ask the user to confirm before proceeding
- Brainstorm only the first confirmed sub-project; each gets its own spec cycle

This prevents scope creep and keeps each design focused.

### 3. Brainstorm with the user

First, diagnose the fuzziness level of the request — this determines your approach:

- **Problem unclear** (user describes symptoms, not goals) → Help them articulate what problem they're actually solving, not just what feature they want.
- **Direction unclear** (user has a goal but no sense of approach) → Sketch a few representative directions before narrowing. Use "what if" scenarios, perspective switching, constraint removal.
- **Boundaries unclear** (user knows what they want but not the edges) → Use the domain-relevant blind-spot checklist from [the guide](references/brainstorming-guide.md) to surface gaps.
- **Solution unclear** (user knows what and scope, needs technical approach) → Skip to step 4 for approach comparison.

Guidelines for this conversation:
- Ask questions one at a time. Prefer multiple-choice when it helps the user choose between meaningful options.
- For non-trivial requests, surface at least one assumption worth confirming and one potential blind spot.
- Explore freely here — don't dismiss ideas prematurely. YAGNI comes later in convergence.

Stop when: the user can clearly state what they want, why they want it, and what they explicitly do not want, and key constraints and success criteria are known.

### 4. Converge and propose

Summarize what brainstorming revealed:
- Refined problem statement
- Challenged assumptions and discovered blind spots
- Trimmed scope

Then propose 1-3 approaches with trade-offs. Lead with your recommendation. If only one approach is viable, explain why alternatives were ruled out.

Ask the user to confirm the problem summary and select an approach before continuing.

### 5. Present the design

Scale the presentation to complexity:
- **Simple projects**: present the full design at once, ask for approval.
- **Moderate/complex projects**: present by section, ask for approval after each.

Write and present the design in the same language the user uses unless they explicitly ask for another language.

Cover: architecture, components, data flow, error handling, and testing.

### 6. Write the design doc

Write to `specs/<topic>/design.md` using the template in `assets/design-doc-template.md`.

Write the document in the same language the user uses unless they explicitly ask for another language.

Name `<topic>` using kebab-case derived from the project or feature name (e.g., `user-auth`, `payment-integration`). If the project already has a specs directory (check for `specs/`, `docs/specs/`, `.codex/specs/`), use the existing convention. Confirm the path with the user if ambiguous.

### 7. User review gate

Ask the user to review the written doc. On feedback:
- Wording or detail changes → edit the doc and re-confirm
- Scope or approach changes → return to step 4
- Missing context → return to step 3

Proceed only after the user approves the written document.

### 8. Hand off to spec-plan

Invoke `spec-plan` as the only next step. Pass this context:
- **Project name**: `<topic>`
- **Target directory**: `specs/<topic>/` (or wherever the design doc lives)
- **Scope summary**: Summary + Non-Goals + Discovery/Scope Decisions (if present)
- **Constraints**: Context section + any confirmed constraints, risks, or assumptions from Discovery
- **Primary users and goals**: Primary Users/Roles + Goals

If a Discovery section exists, treat it as input context for requirements and task planning — especially confirmed assumptions, surfaced risks, and explicit scope decisions.

Do not invoke any implementation skill.

## Verification

- [ ] The design doc exists at the agreed target path and the user approved it
- [ ] The design doc includes the required headings from `assets/design-doc-template.md` and uses the optional Discovery section only when needed
- [ ] The design covers architecture, components, data flow, error handling, and testing
- [ ] If brainstorming surfaced notable discoveries, they are recorded in the Discovery section
- [ ] `spec-plan` has been invoked as the terminal action

## Safety & guardrails

- **No implementation before approval.** Do not write code, scaffold projects, or invoke any implementation skill until the design is approved.
- **Even simple projects get a design** — just keep it short when the scope is small.
- **Scale the process to complexity.** Simple projects may complete steps 3-5 in a few exchanges; don't pad with unnecessary ceremony.
- **One question at a time.** Don't overwhelm the user with multiple questions.
- **YAGNI at the right time.** During brainstorming (step 3), explore freely. During convergence and design (steps 4+), cut ruthlessly — remove anything not essential to the core goal.
- **Design for isolation.** Break the system into units with one clear purpose and well-defined interfaces.

## References

- [Brainstorming guide](references/brainstorming-guide.md) — fuzziness diagnosis, blind-spot checklists, exploration techniques
- [Design doc template](assets/design-doc-template.md) — output format
