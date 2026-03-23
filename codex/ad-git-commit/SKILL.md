---
name: ad-git-commit
description: Generate git commit messages and submit commits. Use when preparing commit messages, deciding split commits, or running git commit.
# metadata:
#   short-description: Git commit message and submit
---

# Git Commit Skill

## When to use

- Generate a clear commit message and submit a commit.
- Decide whether to split changes into multiple commits.
- Align commit messages with existing repo history.

## When not to use

- Resolve merge conflicts or perform rebase/squash workflows.
- You are not inside a git repository or lack commit permissions.

## Inputs

- Working tree and index diffs: `git status --short`, `git diff --stat`, `git diff --cached`
- Recent history: `git log -5 --oneline`

## Outputs

- One or more commits with consistent, readable messages.

## Workflow

1. Collect context (PowerShell):
   `git status --short`
   `git diff --stat`
   `git diff --cached`
   `git log -5 --oneline`
   `git rev-parse --verify HEAD`

2. Detect initial commit:

- If `git rev-parse --verify HEAD` fails, treat this as the first commit.
- For the first commit, do not split; commit everything together and skip split confirmation.

3. Evaluate split vs single commit (by functional relevance):

- If code and docs serve the same feature, one commit is fine and proceed without confirmation.
- If changes are clearly split by purpose or scope, proceed with split commits without confirmation.
- Only when `git diff` analysis is uncertain (e.g., mixed new features and refactors) ask the user to confirm single vs split.
- After the user confirms single vs split, proceed with the workflow without further confirmation prompts.

4. Staging policy:

- Do not run `git add -A` automatically.
- Stage only the confirmed group: `git add <paths>`.
- Verify staged content: `git diff --cached`.

5. Generate commit message:

- Use Chinese for the main descriptive text, but allow English technical terms or code snippets when needed.
- Header format must be: `[<emoji>] <type>(<scope>): <subject>`.
- `<type>(<scope>): <subject>` is required and must not be omitted.
- Emoji rule: if recent commit messages include emoji, include emoji; if recent commit messages do not include emoji, do not add emoji; if there is no history, include emoji.
- After the header, leave a blank line, then add `<body>`; after `<body>`, leave a blank line, then add `<footer>`.
- Conventions for emoji, type, scope, subject, body, and footer are in `references/commit-convention.md`.

6. Commit:

- Show the full message but do not ask for confirmation.
- Commit with one or more `-m` flags as needed:
  `git commit -m "<header>" -m "<body>" -m "<footer>"`

7. Verify:
   `git log -1 --stat`
   `git status`

## Verification

- `git log -1 --stat` shows the intended files and message.
- `git status` is clean or only includes intentionally uncommitted files.

## Safety & guardrails

- Never commit secrets, tokens, or personal data.
- Do not commit with unresolved conflicts.
- Only request confirmation when single vs split is uncertain.
- Do not request additional commit confirmations beyond the single vs split decision.

## References

- [Commit conventions](references/commit-convention.md)
