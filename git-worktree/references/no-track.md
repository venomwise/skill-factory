# Upstream and `--no-track`

## Why the new branch inherits the source's remote by default

When you branch from a remote-tracking ref, git decides whether to set an
upstream based on `branch.autoSetupMerge` (default `true`):

```
git worktree add -b feature <path> origin/main
# feature.remote = origin
# feature.merge  = refs/heads/main   <-- upstream is origin/main
```

Consequences of this default:

- `git status` in the new worktree reports ahead/behind against `origin/main`.
- With `push.default = upstream`, a bare `git push` targets `origin/main` — it
  pushes your feature work straight onto the source branch. This is the failure
  mode this skill prevents.

Branching from a *local* branch does not set an upstream, so `--no-track` is a
no-op there but harmless — always pass it for consistency.

## The fix

```
git worktree add --no-track -b feature <path> origin/main
```

`--no-track` overrides `branch.autoSetupMerge` for this branch only. The new
branch starts at `origin/main`'s commit but has no `remote`/`merge` config.

## Establishing the branch's own remote later

From inside the worktree, the first push creates a matching remote branch and
sets upstream to *itself*, not the source:

```
git push -u origin feature   # creates origin/feature, sets feature -> origin/feature
```

## Repairing a branch that already tracks the wrong remote

```
git branch --unset-upstream feature          # drop the inherited upstream
# or point it at the correct remote branch:
git branch --set-upstream-to=origin/feature feature
```

## Verifying no upstream is set

Any of these confirm the branch is untracked:

```
git config branch.feature.remote   # -> empty, exit 1
git config branch.feature.merge    # -> empty, exit 1

git rev-parse --abbrev-ref --symbolic-full-name 'feature@{upstream}'
# -> fatal: no upstream configured for branch 'feature'
```

In PowerShell, `@{...}` is a hashtable literal — the `feature@{upstream}` ref
MUST be single-quoted, or git only receives `feature@` and reports a misleading
`ambiguous argument 'feature@'`. The `git config` checks have no such pitfall,
so prefer them.
