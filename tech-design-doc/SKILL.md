---
name: tech-design-doc
description: Generate a structured Technical Design Document (TDD). Use this skill whenever the user wants to write a technical design doc, system design spec, or technical proposal — including phrases like "write a tech design", "generate a TDD", "create a design doc", "document the technical approach", "write up the design", "帮我写个技术设计", "生成技术方案", "写个设计文档", "出一份 TDD", "整理一下技术方案", "写技术设计文档". Supports two input sources: a design.md produced by the brainstorming skill, or git commit history. Outputs a standardized doc with five sections: Background, Goals, Solution Design, Design Details (Function/API/Database), and Key Points.
---

# Technical Design Document Generator

Generate a structured Technical Design Document from existing inputs (design.md or git history) to align the team on a technical approach before implementation.

## Input Sources

Two supported inputs, in priority order:

1. **design.md** — produced by the `brainstorming` skill, typically at `specs/<topic>/design.md` or `.codex/specs/<topic>/design.md`
2. **git commit history** — analyze recent relevant commits via `git log` to extract change intent and technical context

## Output

- Technical design document, default path: `docs/design/<topic>.md`
- Follow existing conventions if the project already has a `docs/` or `design/` directory

## Workflow

### 1. Identify the input source

Ask the user which input source to use:
- **design.md** from the brainstorming skill
- **git commit history**

If the user chooses design.md, probe for it automatically: check `specs/` and `.codex/specs/` for an existing file. If found, use it directly. If not found, ask the user for the path.

### 2. Gather context

**Source A: design.md**

Read the file and extract:
- Feature goals and background
- Technical approach and architecture decisions
- Interface and data model details
- Confirmed constraints and assumptions

**Source B: git commit history**

```bash
git log --oneline -50
```

Filter to relevant commits by grepping the user's keyword against commit messages:

```bash
git log --oneline -50 | grep -i "<keyword>"
```

For each matched commit, inspect what changed:

```bash
git show <commit-hash> --stat
git diff <commit-hash>^ <commit-hash> -- <relevant-files>
```

Extract:
- Which modules/files changed
- Business intent from commit messages
- New or modified interfaces and data structures

Only analyze commits that match the keyword — ignore unrelated history.

### 3. Confirm missing information

Based on gathered context, identify what's still unclear and ask the user to fill gaps — one question at a time. Typical gaps to check:

- Background: is the business or technical context sufficient for a reader with no prior knowledge?
- Scope: are there explicit non-goals worth calling out?
- API: does this involve new or modified endpoints? (determines whether the API Design section is needed)
- Database: does this involve schema changes? (determines whether the Database Design section is needed)
- Key points: any implementation gotchas or constraints the reader must know?

Don't ask about things already clear from the gathered context.

### 4. Generate the document

Use `assets/tech-design-template.md` to produce the document.

Section guidelines:
- **Background**: Describe the technical and business context — the goal is that a reader with no prior knowledge can understand *why* this design exists and what drove it.
- **Goals**: List the problems being solved or features being added in concise bullet points. Include explicit non-goals when scope boundaries matter — they prevent misunderstandings later.
- **Solution Design**: Summarize the overall approach so reviewers can evaluate the direction before reading details. When the flow involves multiple systems or non-obvious steps, add a Mermaid diagram — a picture resolves ambiguity faster than prose.
- **Design Details — Functional Design**: Walk through the processing flow and working logic at the level a developer needs to implement it, without dropping to code. Think: what happens, in what order, under what conditions.
- **Design Details — API Design** *(only when API changes are involved)*: List endpoints, HTTP methods, and request/response structures. Enough detail for a consumer to integrate without reading the code.
- **Design Details — Database Design** *(only when schema changes are involved)*: List table names, field changes, and index changes. Include the reason for structural decisions when non-obvious.
- **Key Points**: Surface implementation gotchas, constraints, and decisions that would surprise a developer reading the code later. If it's obvious from the design, skip it.

### 5. User review

Present the draft to the user:
- Wording or detail changes → edit in place
- Approach or scope changes → return to step 3
- Missing context → return to step 2

Write the file to the target path only after the user approves.

## Verification

- [ ] Document contains all required sections (Background, Goals, Solution Design, Functional Design, Key Points)
- [ ] Optional sections (API Design, Database Design) appear only when actually applicable
- [ ] Solution Design includes a diagram if the flow involves multiple systems or non-obvious steps
- [ ] Document written to target path and approved by user

## References

- [Document template](assets/tech-design-template.md)

