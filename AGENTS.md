# Repository Guidelines

## Project Structure & Module Organization
This repository is a skill factory for AI coding agents. Each skill has its own folder at the root level with a `SKILL.md` entry point plus optional `assets/`, `references/`, and `scripts/`. Current evaluation cases live in `evals/`, organized by skill name. Use `CLAUDE.md` as the repo-level reference for skill pairing, spec output conventions, and commit style.

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
- `<skill>-test.yml`: Runs only for changes under `<skill>-go/**` or the workflow file.
- `<skill>-release.yml`: Builds release archives only for `<skill>-v*` tags or manual dispatch.
- `<skill>-update-skill.yml`: Updates only `<skill>/bin/**` after that skill's release workflow succeeds.

Current Go workflow sets:
- `db-explorer-test.yml` / `db-explorer-release.yml` / `db-explorer-update-skill.yml`
- `exa-search-release.yml` / `exa-search-update-skill.yml`
- `grok-search-test.yml` / `grok-search-release.yml` / `grok-search-update-skill.yml`

See `.github/workflows/QUICKSTART.md` for release workflow details.

## Coding Style & Naming Conventions
Use 4-space indentation in Python and keep functions, variables, and files in `snake_case`. Keep Markdown concise, instructional, and structured with clear headings. Every skill must expose a `SKILL.md` with YAML frontmatter, especially `name` and `description`. Place reusable templates in `assets/`, supporting docs in `references/`, and executable helpers in `scripts/`.

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
Recent history follows lightweight Conventional Commit patterns, usually with a scope and Chinese subject line, for example `fix(db-explorer): 修复 URL 解码问题` or `docs(db-explorer): 将 SKILL.md 从中文改写为英文`. Use `feat`, `fix`, `docs`, `test`, `refactor`, or `chore`, and keep each commit focused on one change. PRs should summarize the affected skill(s), explain behavior changes, list verification commands, and include sample outputs when a prompt flow, script result, or eval expectation changes.
