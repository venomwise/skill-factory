# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Purpose

This is a **skill factory** — a collection of reusable AI agent skills for Claude Code (`claude-code/`) and Codex (`codex/`). Each skill is a self-contained directory with a `SKILL.md` file that defines how an AI agent should behave when invoked.

## Structure

```
skill-factory/
├── claude-code/          # Skills for Claude Code (claude.ai/code)
│   ├── spec-plan/        # Generate requirements.md + tasks.md
│   └── spec-exec/        # Execute tasks from tasks.md
└── codex/                # Skills for Codex (OpenAI)
    ├── ad-brainstorming/ # Brainstorm ideas into design docs
    ├── ad-git-commit/    # Generate and submit git commits
    ├── ad-spec-exec/     # Execute spec tasks (Codex variant)
    └── ad-spec-plan/     # Generate spec (Codex variant)
```

Each skill directory follows this layout:
- `SKILL.md` — the skill definition (frontmatter + workflow)
- `assets/` — templates referenced by the skill
- `references/` — supporting reference docs

## Skill File Format

`SKILL.md` files use YAML frontmatter followed by Markdown:

```markdown
---
name: skill-name
description: One-line description used for skill routing/matching.
---

# Skill Title
## When to use
## When not to use
## Inputs / Outputs
## Workflow
## Verification
## Safety & guardrails
## References
```

The `description` field in frontmatter is critical — it determines when the skill is invoked.

## Skill Pairing Pattern

Skills come in complementary pairs:
- **spec-plan** produces `requirements.md` + `tasks.md`
- **spec-exec** consumes `tasks.md`, implements tasks, updates checkboxes as each completes
- **ad-brainstorming** precedes spec-plan; produces `design.md` at `.codex/specs/<topic>/design.md`

## Spec Output Conventions

- Default spec location: `.codex/specs/<project-name>/` (also `specs/` or `docs/specs/`)
- `requirements.md` uses `Requirement N` headings; criteria referenced as `N.M` in tasks
- `tasks.md` uses checkbox list items for phases: `- [ ] Phase N: Title` (never `###` headings)
- Optional tasks: `[optional]` suffix (claude-code) or `- [ ]*` marker (codex)
- Checkpoints: `- [ ] Checkpoint: Verify <scope>` — pause points between phases
- Completed tasks: `- [x]` (claude-code) or `- [✅]` / `- [✅]*` (codex)

## Commit Convention

Format: `[<emoji>] <type>(<scope>): <subject>`

- Main text in **Chinese**; English allowed for technical terms
- Emoji: follow repo history (include if recent commits have emoji, omit if not, include if no history)
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Scope from path: `mcp`, `common`, `doc`, `tests`, `repo` (mixed/unclear)
- Body: `- ` bullets wrapped at ~72 chars explaining why/what
- Footer: only for `BREAKING CHANGE:` or `Closes #N`

Full emoji list and examples: `codex/ad-git-commit/references/commit-convention.md`

## Adding a New Skill

1. Create `<platform>/<skill-name>/SKILL.md` with frontmatter `name` and `description`.
2. Add `assets/` templates and `references/` docs if needed.
3. Mirror across platforms (`claude-code/` and `codex/`) if applicable, adjusting platform-specific syntax (PowerShell vs bash, checkbox markers).
