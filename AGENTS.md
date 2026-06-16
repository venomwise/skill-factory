# Repository Guidelines

## Project Structure & Module Organization
This repository is a **skill factory** — a collection of reusable AI agent skills. Each skill is a self-contained folder at the root level with a `SKILL.md` entry point that defines how an agent should behave when the skill is invoked.

```
skill-factory/
├── brainstorming/         # Brainstorm ideas into a validated design.md
├── clarification/         # Clarify underspecified small changes before editing
├── db-explorer/           # Read-only database exploration and querying
├── design-review/         # Review a design.md and produce review.md
├── exa-search/            # Neural web search for documentation
├── git-commit/            # Generate and submit git commits
├── grok-search/           # Real-time web research
├── hld-generator/         # High-level (technical) design documentation
├── skill-authoring/       # Guide for creating AI agent skills
├── spec-plan/             # Generate requirements.md + tasks.md from a design
├── spec-exec/             # Execute tasks from tasks.md
├── springcloud-init/      # Spring Cloud project initialization
├── db-explorer-go/        # Go source for the db-explorer binary
├── exa-search-go/         # Go source for the exa-search binary
├── grok-search-go/        # Go source for the grok-search binary
├── evals/                 # Evaluation cases, organized by skill name
└── specs/                 # Spec outputs (design/requirements/tasks per topic)
```

Each top-level directory that contains a `SKILL.md` is a skill. A skill folder follows this layout:
- `SKILL.md` — the skill definition (YAML frontmatter + workflow)
- `assets/` — reusable templates referenced by the skill
- `references/` — supporting reference docs
- `scripts/` — executable helpers (optional)

The Go-based skills (`db-explorer`, `exa-search`, `grok-search`) ship pre-compiled binaries in their own `bin/` directory; the source lives in the matching `*-go/` project. Evaluation cases live in `evals/<skill>/`. This file (`AGENTS.md`) is the single repo-level source of truth for all conventions below.

## Build, Test, and Development Commands
There is no monolithic build step; work is usually skill-specific.

### Go-based skills (db-explorer, exa-search, grok-search)
```bash
# Build db-explorer from source
cd db-explorer-go
go test ./...
go build -o db-explorer ./cmd/db-explorer
./db-explorer version
cd ..
python3 evals/db-explorer/run_comparison.py

# Build exa-search from source
cd exa-search-go
go build -o exa-search cmd/exa-search/main.go
./exa-search version

# Build grok-search from source
cd ../grok-search-go
go test ./...
go build -o grok-search ./cmd/grok-search
./grok-search version
```

The `db-explorer`, `exa-search`, and `grok-search` skills use a hybrid architecture: pre-compiled Go binaries for supported platforms are stored in each skill's `bin/` directory, while source code lives in independent Go projects (`db-explorer-go/`, `exa-search-go/`, `grok-search-go/`). Skill definitions remain self-contained in their skill directories. Configuration is loaded from `~/.config/ai-skills/<skill>.toml`, project config where supported, environment variables, or CLI flags.

For `db-explorer`, invoke the selected platform binary directly, for example `db-explorer/bin/db-explorer-darwin-arm64 tables --db sqlite --url ./sample.db`. It uses project `.db-explorer.toml`, global `~/.config/ai-skills/db-explorer.toml`, environment variables (`DATABASE_URL`, `DB_URL`, `POSTGRES_URL`, `MYSQL_URL`), or CLI flags (`--db`, `--url`, `--url-env`, `--profile`).

For `grok-search`, invoke the selected platform binary directly, for example `grok-search/bin/grok-search-darwin-arm64 news --query "test"`. It uses TOML config at `~/.config/ai-skills/grok-search.toml`, environment variables (`GROK_API_KEY`, `GROK_API_KEYS`, `GROK_BASE_URL`, `GROK_MODEL`, `GROK_TIMEOUT`), or CLI flags (`--api-key`, `--base-url`, `--model`).

### Automated Releases (Go skills)
GitHub Actions workflows automate building, releasing, and updating Go-based skills. Use project-scoped release tags so one skill does not rebuild another:

```bash
# Exa Search
git tag -a exa-search-v1.0.0 -m "Release exa-search v1.0.0"
git push origin exa-search-v1.0.0

# Grok Search
git tag -a grok-search-v1.0.0 -m "Release grok-search v1.0.0"
git push origin grok-search-v1.0.0
```

Each Go skill owns isolated workflows:
- `<skill>-test.yml`: Runs only for pull requests that touch that skill's paths or workflow file.
- `<skill>-release.yml`: Builds release archives for `<skill>-v*` tags; release workflows must serialize `test -> build -> release` and must not run on normal branch pushes.
- `<skill>-update-skill.yml`: Updates only `<skill>/bin/**` after that skill's release workflow succeeds.

#### GitHub Actions Workflow Requirements

When creating or modifying workflows for Go skills, follow these rules exactly:

1. **Use project-scoped triggers only.** Release tags must be `<skill>-v*` (for example `db-explorer-v0.0.1`). Never use broad `v*` tags.
2. **Do not trigger test or release workflows on normal branch pushes.** A push to `main` must not start `<skill>-test.yml` or `<skill>-release.yml`.
3. **Standalone test workflow is PR-only.** `<skill>-test.yml` may use `pull_request` with skill-specific `paths`, but must not include `push: branches: [main]`.
4. **Release workflow is tag-only unless explicitly approved otherwise.** `<skill>-release.yml` should trigger on:
   ```yaml
   on:
     push:
       tags:
         - '<skill>-v*'
   ```
5. **Release workflow must include matrix tests before build.** The release workflow must have a `test` job using:
   ```yaml
   strategy:
     matrix:
       os: [ubuntu-latest, macos-latest, windows-latest]
   ```
   The test job must run dependency download, `go test`, formatting check, `go vet`, and build+`version` verification. Formatting checks may be Linux-only.
6. **Build must depend on test.** The release `build` job must use `needs: test`; the `release` job must use `needs: build`.
7. **Build all supported binaries with `CGO_ENABLED=0`.** Produce Linux amd64/arm64, macOS amd64/arm64, and Windows amd64 binaries unless the skill explicitly documents a narrower support matrix.
8. **Update-skill workflow must not run on source pushes.** `<skill>-update-skill.yml` should run on successful release workflow completion and optional manual dispatch only. It must update only that skill's `bin/` directory and checksums.
9. **Compare before creating.** Before adding a new workflow, read the latest existing workflow for the most similar Go skill and mirror its trigger/build/release structure. Do not invent a new trigger model.
10. **Validate trigger intent before committing.** Inspect the final YAML and confirm: PR-only test, tag-only release, matrix release tests, `needs` chain, and no accidental `main` push trigger.

Current Go workflow sets:
- `db-explorer-test.yml` / `db-explorer-release.yml` / `db-explorer-update-skill.yml`
- `exa-search-release.yml` / `exa-search-update-skill.yml`
- `grok-search-test.yml` / `grok-search-release.yml` / `grok-search-update-skill.yml`

See `.github/workflows/QUICKSTART.md` for release workflow details.

## Coding Style & Naming Conventions
Use 4-space indentation in Python and keep functions, variables, and files in `snake_case`. Keep Markdown concise, instructional, and structured with clear headings. Every skill must expose a `SKILL.md` with YAML frontmatter, especially `name` and `description`. Place reusable templates in `assets/`, supporting docs in `references/`, and executable helpers in `scripts/`.

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

The `description` field is critical — it determines when the skill is invoked.

## Skill Pairing & Spec Conventions
Several skills are designed to chain into a spec pipeline, all rooted at `specs/<topic>/`:
- **brainstorming** turns an idea into a validated `specs/<topic>/design.md`, then hands off to spec-plan.
- **spec-plan** consumes the approved `design.md` and produces `requirements.md` + `tasks.md`.
- **spec-exec** consumes `tasks.md`, implements each task, and updates checkboxes as they complete; it consults `requirements.md` for acceptance criteria and `design.md` only as background context.

Spec output conventions:
- Default spec location: `specs/<topic>/` (`<topic>` in kebab-case, e.g. `user-auth`).
- `requirements.md` uses `Requirement N` headings; acceptance criteria are referenced as `N.M` in tasks.
- `tasks.md` uses checkbox list items for phases: `- [ ] Phase N: Title` (never `###` headings).
- Optional tasks: `[optional]` suffix (claude-code) or `- [ ]*` marker (codex).
- Checkpoints: `- [ ] Checkpoint: Verify <scope>` — pause points between phases.
- Completed tasks: `- [x]` (claude-code) or `- [✅]` / `- [✅]*` (codex).

## Adding a New Skill
1. Create `<skill-name>/SKILL.md` at the root level with frontmatter `name` and `description`.
2. Add `assets/` templates and `references/` docs if needed.
3. Add `scripts/` for executable helpers if the skill requires automation.
4. Create corresponding eval cases under `evals/<skill-name>/` for testing.

## Skill Design Principles

### Conciseness Principle: Every Token Must Justify Its Cost

The context window is shared between system prompt, skill content, and conversation history. Verbose skills reduce available space for actual work.

**Decision Framework:**
Before including any content in a skill, verify:
1. Does this convey information the AI cannot infer?
2. Would removing this prevent the AI from completing the task?
3. Does this example demonstrate a pattern the AI cannot deduce?

If the answer is "no" to all three, remove it.

**Examples:**

Bad (verbose):
```
PDF files are a common document format that contains text and images.
To extract text from PDFs, you need a specialized library. There are
many options available in Python, such as PyPDF2, pdfplumber, and others.
Each has different features and capabilities...
```

Good (concise):
```
Extract text with pdfplumber:
import pdfplumber
with pdfplumber.open("file.pdf") as pdf:
    text = pdf.pages[0].extract_text()
```

### Direct Execution Principle: Let AI Do What AI Does Best

Don't wrap operations that AI can perform directly. Each abstraction layer adds cognitive load and hides implementation details.

**Decision Framework:**
Before adding a wrapper script or helper tool, ask:
1. Can the AI accomplish this with 2-3 direct commands?
2. Does this wrapper hide information the AI needs to understand?
3. Is this abstraction for human convenience or AI capability?

If the answers are "yes"/"yes"/"human", don't add the wrapper.

**Examples:**

Bad (unnecessary wrapper):
```python
# scripts/detect_platform.py
def get_platform():
    system = platform.system()
    machine = platform.machine()
    # ... 50 lines of mapping logic ...
    return binary_name

# Usage: python scripts/detect_platform.py
```

Good (direct execution):
```bash
# AI detects platform directly
uname -s  # Darwin
uname -m  # arm64

# AI selects binary
bin/tool-darwin-arm64 --query "test"
```

**Why Direct is Better:**
- AI sees the actual execution path
- Transparent for debugging
- Truly zero dependencies
- Teaches platform awareness
- Enables adaptation to edge cases

### When to Add Abstraction

Abstraction is appropriate when:
- The operation requires domain-specific knowledge AI lacks
- The task involves complex state management
- Multiple steps must be atomic (transactions)
- The abstraction genuinely simplifies the mental model

Example of good abstraction:
```go
// internal/db - handles connection setup, driver quirks, and error translation
// AI doesn't need to know database driver internals for routine exploration
```

## Testing Guidelines
This repo does not currently enforce a global coverage threshold. Instead, add focused eval data under `evals/<skill>/` and run the relevant script or workflow manually. For database work, prefer deterministic checks through `evals/db-explorer/run_comparison.py`. Name new eval artifacts descriptively, following existing patterns such as `evals.json`, `grades.json`, and `benchmark.md`.

## Commit & Pull Request Guidelines
Recent history follows lightweight Conventional Commit patterns, usually with a scope and Chinese subject line, for example `fix(db-explorer): 修复 URL 解码问题` or `docs(db-explorer): 将 SKILL.md 从中文改写为英文`.

Format: `[<emoji>] <type>(<scope>): <subject>`
- Main subject text in **Chinese**; English allowed for technical terms.
- Emoji: follow repo history (include if recent commits have emoji, omit if not, include if no history).
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.
- Scope from path: skill name (`db-explorer`, `brainstorming`, …), or `mcp`, `common`, `doc`, `tests`, `repo` (mixed/unclear).
- Body: `- ` bullets wrapped at ~72 chars explaining why/what.
- Footer: only for `BREAKING CHANGE:` or `Closes #N`.

Full emoji list and examples: `git-commit/references/commit-convention.md`.

Keep each commit focused on one change. PRs should summarize the affected skill(s), explain behavior changes, list verification commands, and include sample outputs when a prompt flow, script result, or eval expectation changes.
