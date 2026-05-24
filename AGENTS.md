# Repository Guidelines

## Project Structure & Module Organization
This repository is a skill factory for AI coding agents. Each skill has its own folder at the root level with a `SKILL.md` entry point plus optional `assets/`, `references/`, and `scripts/`. Current evaluation cases live in `evals/`, organized by skill name. Use `CLAUDE.md` as the repo-level reference for skill pairing, spec output conventions, and commit style.

## Build, Test, and Development Commands
There is no monolithic build step; work is usually skill-specific.

### Python-based skills (db-explorer)
```powershell
python -m venv .\db-explorer\.venv
.\db-explorer\.venv\Scripts\pip install -r .\db-explorer\requirements.txt
.\db-explorer\.venv\Scripts\python .\evals\db-explorer\run_comparison.py
python .\db-explorer\scripts\db_query.py --db-type sqlite --url .\sample.db tables
```

Use the same Python environment for dependency installation and script execution. `run_comparison.py` is the main regression check for `db-explorer`; it compares the current script against the baseline and refreshes grading output.

### Go-based skills (exa-search)
```bash
cd exa-search-go
go build -o exa-search cmd/exa-search/main.go
./exa-search version
./exa-search search --query "test" --api-key <key>
```

The `exa-search` binary is a statically-compiled Go application with zero runtime dependencies. Source code lives in `exa-search-go/` (independent Go project), while the skill definition is in `exa-search/` (lightweight, documentation-focused). Build with `go build` and run the resulting binary directly. Configuration is loaded from `~/.config/ai-skills/exa-search.toml` (auto-created on first run), environment variables (`EXA_API_KEY`, `EXA_API_KEYS`), or CLI flags (`--api-key`).

### Automated Releases (exa-search)
GitHub Actions workflows automate building and releasing multi-platform binaries:

```bash
# Trigger a release manually via GitHub Actions UI:
# 1. Go to Actions → "Build and Release exa-search"
# 2. Click "Run workflow"
# 3. Enter version (e.g., v1.0.0)

# Or create a git tag:
git tag -a exa-search-v1.0.0 -m "Release v1.0.0"
git push origin exa-search-v1.0.0
```

The workflow builds for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64), generates SHA256 checksums, and creates a GitHub Release with installation instructions. See `.github/workflows/QUICKSTART.md` for details.

## Coding Style & Naming Conventions
Use 4-space indentation in Python and keep functions, variables, and files in `snake_case`. Keep Markdown concise, instructional, and structured with clear headings. Every skill must expose a `SKILL.md` with YAML frontmatter, especially `name` and `description`. Place reusable templates in `assets/`, supporting docs in `references/`, and executable helpers in `scripts/`.

## Testing Guidelines
This repo does not currently enforce a global coverage threshold. Instead, add focused eval data under `evals/<skill>/` and run the relevant script or workflow manually. For database work, prefer deterministic checks through `evals/db-explorer/run_comparison.py`. Name new eval artifacts descriptively, following existing patterns such as `evals.json`, `grades.json`, and `benchmark.md`.

## Commit & Pull Request Guidelines
Recent history follows lightweight Conventional Commit patterns, usually with a scope and Chinese subject line, for example `fix(db-explorer): 修复 URL 解码问题` or `docs(db-explorer): 将 SKILL.md 从中文改写为英文`. Use `feat`, `fix`, `docs`, `test`, `refactor`, or `chore`, and keep each commit focused on one change. PRs should summarize the affected skill(s), explain behavior changes, list verification commands, and include sample outputs when a prompt flow, script result, or eval expectation changes.
