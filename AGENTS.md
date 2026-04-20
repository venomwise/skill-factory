# Repository Guidelines

## Project Structure & Module Organization
This repository is a skill factory for AI coding agents. Primary source material lives in `claude-code/` and `codex/`, where each skill has its own folder with a `SKILL.md` entry point plus optional `assets/`, `references/`, and `scripts/`. Current evaluation cases live in `evals/`, and exploratory outputs are stored under `brainstorming-workspace/`. Use `CLAUDE.md` as the repo-level reference for skill pairing, spec output conventions, and commit style.

## Build, Test, and Development Commands
There is no monolithic build step; work is usually skill-specific.

```powershell
python -m venv .\claude-code\db-explorer\.venv
.\claude-code\db-explorer\.venv\Scripts\pip install -r .\claude-code\db-explorer\requirements.txt
.\claude-code\db-explorer\.venv\Scripts\python .\evals\db-explorer\run_comparison.py
python .\claude-code\db-explorer\scripts\db_query.py --db-type sqlite --url .\sample.db tables
```

Use the same Python environment for dependency installation and script execution. `run_comparison.py` is the main regression check for `db-explorer`; it compares the current script against the baseline and refreshes grading output.

## Coding Style & Naming Conventions
Use 4-space indentation in Python and keep functions, variables, and files in `snake_case`. Keep Markdown concise, instructional, and structured with clear headings. Every skill must expose a `SKILL.md` with YAML frontmatter, especially `name` and `description`. Place reusable templates in `assets/`, supporting docs in `references/`, and executable helpers in `scripts/`. When a capability exists for both platforms, keep `claude-code/` and `codex/` variants aligned unless there is a documented platform-specific reason to diverge.

## Testing Guidelines
This repo does not currently enforce a global coverage threshold. Instead, add focused eval data under `evals/<skill>/` and run the relevant script or workflow manually. For database work, prefer deterministic checks through `evals/db-explorer/run_comparison.py`. Name new eval artifacts descriptively, following existing patterns such as `evals.json`, `grades.json`, and `benchmark.md`.

## Commit & Pull Request Guidelines
Recent history follows lightweight Conventional Commit patterns, usually with a scope and Chinese subject line, for example `fix(db-explorer): 修复 URL 解码问题` or `docs(db-explorer): 将 SKILL.md 从中文改写为英文`. Use `feat`, `fix`, `docs`, `test`, `refactor`, or `chore`, and keep each commit focused on one change. PRs should summarize the affected skill(s), explain behavior changes, list verification commands, and include sample outputs when a prompt flow, script result, or eval expectation changes.
