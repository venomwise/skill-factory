<div align="center">

# Skill Factory for AI Coding Agents

**A curated repository of reusable, eval-friendly skills for AI coding agents.**

Turn hard-won prompting and workflow knowledge into modular assets your team can build, test, refine, and reuse.

[English](./README.md) · [简体中文](./README.zh-CN.md)

![Repo Type](https://img.shields.io/badge/repo-skill--factory-111827?style=flat-square)
![For AI Agents](https://img.shields.io/badge/for-AI%20coding%20agents-7c3aed?style=flat-square)
![Skills](https://img.shields.io/badge/skills-10%2B-0ea5e9?style=flat-square)
![Tooling](https://img.shields.io/badge/tooling-Python%20%2B%20Markdown-3776AB?style=flat-square)

</div>

---

## Table of Contents

- [Why this repository](#why-this-repository)
- [What is a skill in this repo?](#what-is-a-skill-in-this-repo)
- [Included skills](#included-skills)
- [Repository layout](#repository-layout)
- [Quick start](#quick-start)
- [Development and testing](#development-and-testing)
- [Authoring conventions](#authoring-conventions)
- [Contributing](#contributing)
- [License](#license)

## Why this repository

This repository is a **skill factory** for AI coding agents.

Instead of treating prompts and workflows as one-off chat artifacts, it organizes them as reusable skills with:

- **Clear entry points** via `SKILL.md`
- **Reusable assets** such as templates and references
- **Optional automation** through scripts
- **Eval-ready structure** for iterative improvement
- **Practical coverage** across planning, execution, research, database exploration, and developer workflows

It is designed for teams who want to make agent behavior more **repeatable, reviewable, and scalable**.

## What is a skill in this repo?

Each skill is a self-contained folder at the repository root.

Typical structure:

```text
<skill-name>/
├── SKILL.md        # required entry point
├── assets/         # optional templates
├── references/     # optional supporting docs
└── scripts/        # optional helpers / automation
```

A `SKILL.md` file uses YAML frontmatter plus Markdown instructions:

```md
---
name: skill-name
description: What the skill does and when the agent should use it.
---
```

The `description` is especially important because it helps determine **when the skill should trigger**.

## Included skills

> This repository contains both end-user skills and skill-authoring / evaluation utilities.

| Skill | Purpose |
| --- | --- |
| [brainstorming](./brainstorming/) | Turn rough ideas into scoped designs and specs through collaborative dialogue. |
| [db-explorer](./db-explorer/) | Read-only exploration of PostgreSQL, MySQL, and SQLite databases. |
| [exa-search](./exa-search/) | Source-first web research for docs, APIs, and structured content retrieval. |
| [git-commit](./git-commit/) | Generate and submit git commits aligned with repo history and conventions. |
| [grok-search](./grok-search/) | Real-time web research for fresh updates, discourse, and multi-source synthesis. |
| [skill-authoring](./skill-authoring/) | Best practices for creating and improving AI agent skills. |
| [skill-creator](./skill-creator/) | Create, benchmark, and iteratively improve skills. |
| [skill-factory](./skill-factory/) | A systematic workflow for creating or optimizing skills with checkpoints and tooling. |
| [spec-exec](./spec-exec/) | Execute implementation tasks from `specs/<spec>/tasks.md` and track progress. |
| [spec-plan](./spec-plan/) | Generate `requirements.md` and `tasks.md` for a project spec. |
| [tech-design-doc](./tech-design-doc/) | Generate structured technical design documents from design files or git history. |

## Repository layout

```text
.
├── brainstorming/
├── db-explorer/
├── exa-search/
├── git-commit/
├── grok-search/
├── skill-authoring/
├── skill-creator/
├── skill-factory/
├── spec-exec/
├── spec-plan/
├── tech-design-doc/
├── evals/
│   ├── brainstorming/
│   ├── db-explorer/
│   └── exa-search/
├── AGENTS.md
└── CLAUDE.md
```

Additional workspace or experimental directories may exist alongside these core folders.

## Quick start

### 1) Pick a skill

Browse the repository root and open the skill directory you need.

### 2) Read `SKILL.md`

That file defines the skill's purpose, workflow, guardrails, and expected outputs.

### 3) Install dependencies if the skill includes scripts

Some skills ship with Python helpers and `requirements.txt`.
Use the **same virtual environment** for installation and execution.

### 4) Run skill-specific scripts or evals

Example for `db-explorer`:

```powershell
python -m venv .\db-explorer\.venv
.\db-explorer\.venv\Scripts\pip install -r .\db-explorer\requirements.txt
.\db-explorer\.venv\Scripts\python .\evals\db-explorer\run_comparison.py
python .\db-explorer\scripts\db_query.py --db-type sqlite --url .\sample.db tables
```

## Development and testing

There is **no monolithic build step** for the repository.
Most work is skill-specific.

General testing approach:

- Add focused evaluation data under `evals/<skill>/`
- Run the relevant script or workflow manually
- Prefer deterministic checks when possible
- Keep benchmark and grading artifacts descriptive

For `db-explorer`, the main regression check is:

```powershell
.\db-explorer\.venv\Scripts\python .\evals\db-explorer\run_comparison.py
```

## Authoring conventions

When adding or improving skills in this repository:

- Use **4-space indentation** in Python
- Keep Python names and filenames in **`snake_case`**
- Keep Markdown **concise, instructional, and well-structured**
- Every skill must include a **`SKILL.md`** with YAML frontmatter
- Always define at least:
  - `name`
  - `description`
- Put reusable templates in `assets/`
- Put supporting documents in `references/`
- Put executable helpers in `scripts/`

For repo-specific guidance, see:

- [AGENTS.md](./AGENTS.md)
- [CLAUDE.md](./CLAUDE.md)

## Contributing

Contributions should keep the repository modular, practical, and easy to evaluate.

Recommended workflow:

1. Add or update a skill directory
2. Keep `SKILL.md` focused and triggerable
3. Add or refresh eval data under `evals/<skill>/`
4. Verify any scripts with the same environment used for dependency installation
5. Document behavior changes clearly

Commit style in recent history follows lightweight Conventional Commit patterns, for example:

```text
fix(db-explorer): 修复 URL 解码问题
docs(db-explorer): 将 SKILL.md 从中文改写为英文
```

Preferred commit types:

- `feat`
- `fix`
- `docs`
- `test`
- `refactor`
- `chore`

## License

A repository-wide license is **not currently defined** at the root.

Note that [`skill-creator/LICENSE.txt`](./skill-creator/LICENSE.txt) applies to components included in `skill-creator/`.
