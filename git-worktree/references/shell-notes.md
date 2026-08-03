# Shell-Specific Notes

This document covers shell-specific considerations when running git worktree commands.

## PowerShell

### Hash Table Literal Issue

In PowerShell, `@{...}` is a hash table literal. When using git references that contain `@{upstream}` or similar syntax, you must quote the reference to prevent PowerShell from interpreting it as a hash table.

**Problem:**
```powershell
# This will fail with "ambiguous argument" error
git rev-parse --abbrev-ref --symbolic-full-name feature-branch@{upstream}
```

PowerShell sees `@{upstream}` and tries to parse it as a hash table, so git only receives `feature-branch@` (without the `{upstream}` part).

**Solution:**
```powershell
# Use single quotes to prevent PowerShell interpretation
git rev-parse --abbrev-ref --symbolic-full-name 'feature-branch@{upstream}'
```

This issue affects any git command that uses `@{...}` syntax, including:
- `git rev-parse branch@{upstream}`
- `git log branch@{1.day.ago}`
- `git reflog show HEAD@{2}`

### Path Separators

PowerShell on Windows supports both forward slashes (`/`) and backslashes (`\`) in paths. Git commands prefer forward slashes, but both work:

```powershell
# Both are valid
git worktree add ../worktrees/feature origin/main
git worktree add ..\worktrees\feature origin/main
```

### Quoting Paths with Spaces

When paths contain spaces, use quotes:

```powershell
git worktree add "../my worktrees/feature name" origin/main
```

## Bash / Zsh

### No Special Quoting Needed

In bash and zsh, `@{...}` has no special meaning, so git upstream syntax works directly:

```bash
git rev-parse --abbrev-ref --symbolic-full-name feature-branch@{upstream}
```

### Path Separators

Use forward slashes:

```bash
git worktree add ../worktrees/feature origin/main
```

### Quoting Paths with Spaces

Use quotes or escape spaces:

```bash
git worktree add "../my worktrees/feature name" origin/main
# or
git worktree add ../my\ worktrees/feature\ name origin/main
```

## Windows Command Prompt (cmd.exe)

### Path Separators

Backslashes are standard, but forward slashes also work:

```cmd
git worktree add ..\worktrees\feature origin/main
```

### Quoting Paths with Spaces

Use double quotes:

```cmd
git worktree add "..\my worktrees\feature name" origin/main
```

## General Best Practices

1. **Always quote paths with spaces** regardless of shell
2. **Use relative paths** (like `../worktrees/`) for portability across team members
3. **Test git commands** after creating worktrees to ensure they work in your shell environment
4. **Check shell documentation** if you encounter unexpected parsing errors
