---
name: spec-plan
description: Generate a structured requirements spec and implementation plan for any software project, feature, module, or integration. Use this skill whenever the user wants to scope out work, turn a vague idea into concrete requirements and tasks, create a project spec, plan an implementation, break down a feature into actionable steps, or needs a requirements document with traceability to an execution plan. Also trigger when users say things like "spec this out", "plan this feature", "write requirements for", "break this down into tasks", "scope this project", or "what would the implementation plan look like".
---

# Spec Plan

Generate a requirements document and implementation plan that work together — requirements define *what* to build, the plan defines *how* to build it, and every task traces back to at least one requirement.

## When to use

- Scoping a new feature, module, service, or integration
- Turning a vague idea into clear requirements and actionable tasks
- Creating traceability between requirements and implementation steps
- Planning work that spans multiple files or components

## When NOT to use

- Quick one-off edits or single-file changes — just do the work
- The spec already exists and only needs implementation
- Pure research or investigation tasks

## Workflow

### 1. Gather context

Before writing anything, make sure you understand:

- **What** is being built (core capability, problem it solves)
- **Who** it's for (user roles, personas)
- **Scope boundaries** (what's in, what's explicitly out)
- **Constraints** (platforms, dependencies, performance targets, timelines)
- **Scale** — is this a small feature (2-5 tasks) or a large system (20+ tasks)?
- **Language** — write the spec in the same language the user uses. If the user writes in Chinese, the entire spec should be in Chinese. If English, use English. Match the user's language naturally.

If any of these are unclear, ask. Don't invent requirements from thin air — mark assumptions explicitly and confirm them.

### 2. Choose output location

Ask the user where to put the files. If they don't have a preference, suggest a sensible default based on the project:

- If a `.codex/specs/` or `specs/` directory exists, use `specs/<project-name>/`
- Otherwise suggest `docs/specs/<project-name>/` or just the current directory

### 3. Write requirements.md

Use the template at `assets/requirements.template.md` as your structural guide. The template contains HTML comments with authoring instructions — follow them for guidance but never include them in the final output.

**Calibrate depth to scope:**

| Project scale | Requirements | Criteria per req | Guidance |
|---|---|---|---|
| Small (single feature) | 3–6 | 2–4 | Keep it tight, skip glossary if terms are obvious |
| Medium (module/service) | 5–12 | 3–6 | Include glossary, cover error handling |
| Large (system/platform) | 10–20 | 4–8 | Full glossary, architecture context, cross-cutting concerns |

**Quality bar for each requirement:**
- Testable — someone could write a test for each acceptance criterion
- Specific — no weasel words ("should handle errors appropriately" → say *how*)
- Independent — each requirement stands on its own where possible
- Covers three dimensions: normal flow, error flow, boundary conditions

### 4. Write tasks.md

Use the template at `assets/tasks.template.md` as your structural guide.

**Key principles:**
- Every task references one or more requirement IDs (format: `N.M` where N = requirement number, M = criterion number)
- Tasks specify concrete file paths and function/class names when the codebase is known
- Group tasks into phases that can be completed and verified independently
- Include checkpoint tasks between major phases for the user to pause and validate
- Mark optional tasks (tests, nice-to-haves, stretch goals) with `[optional]` after the task title

**Task granularity:** Each task should be completable in a single focused session. If a task feels like it needs sub-phases, it's too big — split it.

### 5. Self-check before presenting

Scan both files and verify:

- [ ] Every requirement has acceptance criteria covering normal, error, and boundary cases
- [ ] Every task in tasks.md references at least one requirement ID that exists in requirements.md
- [ ] No orphan requirements — each requirement is referenced by at least one task
- [ ] Tasks include concrete file paths / function names (when codebase context is available)
- [ ] At least one checkpoint exists between major phases
- [ ] Optional tasks are clearly marked
- [ ] The scope matches what the user asked for — not inflated, not missing pieces

### 6. Present and iterate

Show both files to the user. Explicitly call out:
- Any assumptions you made
- Areas where you need clarification
- Suggested priorities if not everything can be done

## Format conventions

### requirements.md structure

```
# Requirements: <Project Name>

## Introduction
(1-2 paragraphs: what, why, for whom, scope boundary)

## Glossary
(skip for small/obvious projects)

## Requirements

### Requirement 1: <Capability Title>
**User Story:** As a <role>, I want <goal>, so that <benefit>.

#### Acceptance Criteria
1. WHEN <condition>, THEN the system SHALL <behavior>.
2. ...
```

### tasks.md structure

```
# Implementation Plan: <Project Name>

## Overview
(Link to requirements.md, summarize phases and approach)

## Tasks

- [ ] Phase 1: <Title>
  - [ ] 1.1 <Task title>
    - <Concrete steps, file paths, function names>
    - _Requirements: 1.1, 1.2_
  - [ ] 1.2 <Task title> [optional]
    - ...
    - _Requirements: 1.3_

- [ ] Checkpoint: Verify <what to validate>

- [ ] Phase 2: <Title>
  ...
```

## Safety and guardrails

- Never invent requirements without user confirmation — mark assumptions with "**Assumption:**" and ask
- Keep requirements testable — if you can't imagine a test for it, rewrite it
- Don't mark tasks as complete unless the work is actually done
- Don't inflate scope — match the depth to what the user actually needs
