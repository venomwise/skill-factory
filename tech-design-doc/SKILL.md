---
name: tech-design-doc
description: >
  Generate a structured Technical Design Document (TDD). Use this skill whenever the user wants to write a technical design doc,
  system design spec, or technical proposal — including phrases like "write a tech design", "generate a TDD", "create a design doc",
  "document the technical approach", "write up the design", "帮我写个技术设计", "生成技术方案", "写个设计文档", "出一份 TDD", "整理一下技术方案",
  "写技术设计文档". Supports two input sources: a design.md produced by the brainstorming skill, or git commit history. Outputs a
  standardized doc with a TL;DR plus sections: Background, Goals, Solution Design (with alternatives), Design Details (Function/API/Database), Impact & Risks, and Key Points.
---

# Technical Design Document Generator

Generate a structured Technical Design Document from existing inputs (design.md or git history) to align the team on a technical approach before implementation.

## Fidelity Principles

**CRITICAL: This skill performs faithful transformation, not creative rewriting.**

When working from an existing design.md:

- **Preserve all design decisions**: Do not add alternatives, trade-offs, or rationale not present in the source
- **Preserve all technical details**: Do not simplify API contracts, database schemas, or implementation specifics
- **Preserve all semantic content**: Do not change field names, types, default values, or interface semantics
- **Preserve unique sections**: If the source has sections like Discovery, Scope Decisions, or Internal API Contracts, keep them
- **Do not invent content**: If the source lacks alternatives comparison, do not fabricate alternative solutions
- **Do not delete key information**: Configuration values, error handling rules, testing strategies must be preserved

**Transformation = reorganization + formatting, NOT reinterpretation + simplification.**

## Writing principles

A good design doc is read, not deciphered. Apply these throughout — they matter more than filling every section:

- **Lead with the conclusion (inverted pyramid)**: Open with a TL;DR — problem, decision, impact. Let the reader decide in 30 seconds whether and where to read deeper.
- **Show the trade-offs, not just the answer**: A decision can't be evaluated without seeing what else was considered and why it was rejected. The alternatives table is the part reviewers question most — but only include it if the source discusses alternatives.
- **A picture beats prose, when it's the right picture**: Topology and ordering are painful to read as text. Use a diagram and pick the right type (see table below).
- **Write for a reader with zero context**: Define terms, state assumptions, name the audience. The author already knows everything; the doc is for everyone else.
- **Concrete over abstract**: Pair every abstract flow with one concrete example — a real request/response, a single end-to-end scenario walkthrough.
- **Make it scannable**: Bold the key conclusions, use tables for comparisons, keep paragraphs short. People skim docs before they read them.

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

- Audience: who needs to read this? (shapes how much background to spell out)
- Background: is the business or technical context sufficient for a reader with no prior knowledge?
- Alternatives: what other approaches were considered, and why was this one chosen? (the alternatives table needs this)
- Scope: are there explicit non-goals worth calling out?
- API: does this involve new or modified endpoints? (determines whether the API Design section is needed)
- Database: does this involve schema changes? (determines whether the Database Design section is needed)
- Impact & risk: which modules or teams does this touch, and how would it roll back?
- Key points: any implementation gotchas or constraints the reader must know?

Don't ask about things already clear from the gathered context.

### 4. Generate the document

**Map source content to template structure while preserving all semantics.**

#### Step 4.1: Analyze source structure

Identify what sections and content types exist in the source:
- Does it have a Discovery or Key Discoveries section?
- Does it have a Scope Decisions section?
- Does it discuss alternative solutions?
- Does it have detailed API contracts (not just examples)?
- Does it have detailed data flow walkthroughs?
- Does it have error handling strategies?
- Does it have testing plans?
- Does it have configuration/constants sections?

#### Step 4.2: Map to template structure

Use `assets/tech-design-template.md` as the base structure, but:

**Preserve source-specific sections:**
- If source has Discovery/Key Discoveries → keep as a dedicated section after Goals
- If source has Scope Decisions → keep as a dedicated section after Goals  
- If source has Internal API Contracts → keep as a dedicated section in Design Details
- If source has Data Flow → keep as a dedicated section in Design Details
- If source has Error Handling → keep as a dedicated section after Design Details
- If source has Testing → keep as a dedicated section after Error Handling
- If source has Configuration → keep as a dedicated section in Design Details

**Section mapping guidelines:**

- **TL;DR**: Extract the core problem, chosen approach, and impact from source. If source lacks a summary, synthesize one from Goals and Proposed Solution sections. Write this last but place it first.

- **Background**: Map from source's Context, Background, or Current State sections. Preserve all technical context and business drivers. Structure as *current state → pain point → why now*.

- **Goals**: Map from source's Goals, Objectives, or Primary Users sections. Preserve all listed goals and non-goals verbatim. Include explicit non-goals when present in source.

- **Solution Design — Overall approach**: Map from source's Proposed Solution, Architecture, or Components sections. Preserve all architectural decisions. Add a Mermaid diagram if topology or ordering is involved (see diagram-type table below).

- **Solution Design — Alternatives**: **ONLY include if source discusses alternatives.** If source compares multiple approaches, create the comparison table with actual content from source. **NEVER invent alternative solutions to fill this section.** If source doesn't discuss alternatives, skip this section entirely.

- **Design Details — Functional Design**: Map from source's workflow, processing logic, or scenario sections. Preserve all implementation details. Include all concrete scenarios from source.

- **Design Details — API Design** *(only when source contains API specs)*: Preserve all API contracts exactly as specified in source — do not change field names, types, or defaults. Include all request/response examples from source.

- **Design Details — Database Design** *(only when source contains schema changes)*: Preserve all table definitions, field types, indexes, and design rationale from source. Include DDL if present in source.

- **Impact & Risks**: Map from source's Impact, Risks, or Rollback sections. Preserve all identified risks and mitigation strategies.

- **Key Points**: Map from source's implementation notes, constraints, or critical decisions. Preserve all gotchas and constraints mentioned in source.

- **Open questions / Glossary** *(optional)*: Include if present in source.

#### Step 4.3: Fidelity checklist

Before finalizing, verify:
- [ ] No design decisions added that weren't in source
- [ ] No alternative solutions invented to fill template
- [ ] All API contracts match source exactly (field names, types, defaults)
- [ ] All configuration values preserved from source
- [ ] All error handling rules preserved from source
- [ ] All testing strategies preserved from source
- [ ] Unique source sections preserved (Discovery, Scope Decisions, etc.)
- [ ] Technical details not simplified or deleted

### Choosing the right Mermaid diagram

| Intent | Diagram type |
|--------|--------------|
| Component/architecture relationships, data flow, decision branches | `flowchart` |
| Interaction ordering across services/actors over time | `sequenceDiagram` |
| Lifecycle / status transitions of an entity | `stateDiagram-v2` |
| Data model, tables and relationships | `erDiagram` |

Skip the diagram when the flow is a trivial linear step — a diagram should reduce ambiguity, not add ceremony.

### 5. User review

Present the draft to the user:
- Wording or detail changes → edit in place
- Approach or scope changes → return to step 3
- Missing context → return to step 2

Write the file to the target path only after the user approves.

## Verification

- [ ] TL;DR at the top states problem, approach, and impact in three bullets
- [ ] Document contains all required sections (Background, Goals, Solution Design, Functional Design, Impact & Risks, Key Points)
- [ ] Solution Design includes an alternatives comparison table **only if source discussed alternatives**
- [ ] Functional Design includes all scenarios from source (not simplified)
- [ ] Optional sections (API Design, Database Design) appear only when present in source
- [ ] Solution Design includes a diagram (correct type) if topology or ordering is involved
- [ ] **Fidelity verification**:
  - [ ] All API contracts match source exactly (no field name/type/default changes)
  - [ ] All configuration values preserved from source
  - [ ] All error handling rules preserved from source
  - [ ] All testing strategies preserved from source
  - [ ] Unique source sections preserved (Discovery, Scope Decisions, Internal API Contracts, Data Flow, Error Handling, Testing, Configuration)
  - [ ] No alternative solutions invented
  - [ ] Technical details not simplified or omitted
- [ ] Document written to target path and approved by user

## References

- [Document template](assets/tech-design-template.md)
