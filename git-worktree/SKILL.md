---
name: git-worktree
description: Create and manage git worktrees for isolated parallel work. Use this skill whenever the user mentions worktrees, wants to work on multiple branches simultaneously, needs to test changes without affecting their main working directory, or asks about creating/managing/removing worktrees. Also trigger when users want to start fresh work without committing current changes, need to review a PR while keeping their current work intact, or mention working on features in parallel. Handles listing, removing, pruning, moving, and locking worktrees.
# metadata:
#   short-description: Create and manage git worktrees
---

# Git Worktree Skill

## When to use

- Create a new worktree: new directory + new branch, checked out from a source
  branch, where the new branch must NOT track the source's remote branch.
- Manage existing worktrees: list, remove, prune stale entries, move, lock/unlock.

## When not to use

- Not inside a git repository.
- Ordinary branch/checkout work in the current worktree (no new working dir).
- User just wants to switch branches (use `git checkout` or `git switch`).
- User wants the new branch to track a remote upstream automatically (worktrees
  created with this skill have no upstream by default; use regular `git checkout -b`
  if tracking is needed).

## Inputs

- Create: `<path>` (new worktree dir), `<new-branch>` (new branch name),
  `<source-branch>` (branch to check out from, local or remote-tracking e.g. `origin/main`).
- Manage: target worktree path (from `git worktree list`).

## Outputs

- A new worktree on a new local branch with **no upstream**, or the result of a
  list/remove/prune/move/lock operation.

## Workflow

All commands below are standard git commands. Adapt them to the user's shell environment (quote paths with spaces appropriately, etc.). For shell-specific considerations, see [Shell-Specific Notes](references/shell-notes.md).

### Create a worktree

1. Confirm you have all three inputs. If `<source-branch>` is a remote branch,
   refresh it first so the new branch starts from the latest commit:
   `git fetch <remote>`   (e.g. `git fetch origin`)

2. Create the worktree and new branch **without an upstream**:
   `git worktree add --no-track -b <new-branch> <path> <source-branch>`

   - `--no-track` is the key flag: it prevents git from auto-setting the new
     branch's upstream to `<source-branch>` (which would happen by default when
     the source is a remote-tracking branch, per `branch.autoSetupMerge`).
   - If `<new-branch>` already exists, `-b` fails; use `-B` only when the user
     explicitly wants to reset it.

3. Verify no upstream was set (see Verification). If one slipped in, unset it:
   `git branch --unset-upstream <new-branch>`

4. Tell the user that the first push must create the branch's own remote:
   `git push -u origin <new-branch>`  (run later, from inside the worktree)

### Manage existing worktrees

- List (always start here to get exact paths):
  `git worktree list`  (add `--porcelain` for machine-readable detail)

- Remove a worktree:
  `git worktree remove <path>`
  - Use `--force` only if it has a dirty tree or submodules and the user confirms.
  - Removing a worktree does NOT delete its branch. Delete the branch only if the
    user asks: `git branch -D <branch>` (run from the main worktree).

- Prune stale entries (directories deleted manually outside git):
  `git worktree prune -v`

- Move a worktree to a new path:
  `git worktree move <path> <new-path>`

- Lock / unlock (prevent auto-prune, e.g. on removable media):
  `git worktree lock <path> --reason "<why>"`
  `git worktree unlock <path>`

## Verification

After create, confirm the branch has **no upstream** using `git config` checks:
```bash
git config branch.<new-branch>.remote   # should be empty, exit 1
git config branch.<new-branch>.merge    # should be empty, exit 1
```

These commands exit non-zero with empty output when the key is absent, which is unambiguous across shells.

Also confirm the worktree registered correctly:
```bash
git worktree list  # should show <path> on <new-branch>
```

After remove/prune:
```bash
git worktree list  # should no longer show the target
```

## Examples

### Example 1: Create a feature worktree

```bash
# 1. Fetch latest from remote
git fetch origin

# 2. Create worktree for new feature (branching from origin/main, no upstream)
git worktree add --no-track -b feature/user-auth ../worktrees/user-auth origin/main

# 3. Verify no upstream was set
git config branch.feature/user-auth.remote  # empty, exit 1
git config branch.feature/user-auth.merge   # empty, exit 1

# 4. List all worktrees
git worktree list

# 5. When ready to push (run from inside the worktree)
cd ../worktrees/user-auth && git push -u origin feature/user-auth
```

### Example 2: Review a PR in isolation

```bash
# Create worktree from PR branch to review without affecting current work
git fetch origin pull/123/head:pr-123
git worktree add --no-track -b review-pr-123 ../worktrees/pr-review pr-123

# After review, remove the worktree
git worktree remove ../worktrees/pr-review
git branch -D review-pr-123  # Clean up branch if done
```

### Example 3: List and clean up worktrees

```bash
# List all worktrees with details
git worktree list --porcelain

# Remove a finished worktree (must have no uncommitted changes)
git worktree remove ../worktrees/completed-feature

# Remove a worktree with uncommitted changes (use with caution!)
git worktree remove --force ../worktrees/abandoned-feature

# Prune stale worktree entries (after manual directory deletion)
git worktree prune -v
```

## Common errors and solutions

### Error: `fatal: '<branch-name>' is already checked out`

**原因**: 尝试创建的分支已在另一个worktree中被检出。

**解决方案**: 使用不同的分支名，或者先移除已有的worktree。

### Error: `fatal: invalid reference: <source-branch>`

**原因**: 源分支不存在或名称错误。

**解决方案**: 
- 检查分支名拼写
- 如果是远程分支，先运行 `git fetch <remote>`
- 使用 `git branch -a` 查看所有可用分支

### Error: Directory already exists

**原因**: 目标路径已存在。

**解决方案**: 选择不同的路径，或删除现有目录后重试。

### Upstream accidentally set

如果发现新分支意外设置了上游：

```powershell
# 取消上游设置
pwsh -Command "git branch --unset-upstream <branch-name>"

# 验证已取消
pwsh -Command "git config branch.<branch-name>.remote"  # should be empty
```

## Integration with Claude Code

当创建worktree后，你可能需要：

1. **切换工作目录**: 使用 `cd` 命令进入新的worktree目录
2. **返回主worktree**: 使用 `cd` 返回原始工作目录
3. **并行工作**: 在不同终端窗口中分别打开不同的worktree目录

注意：Claude Code的当前工作目录会影响后续命令的执行上下文。创建worktree后，询问用户是否需要切换到新worktree目录继续工作。

## Safety & guardrails

- Never delete a branch as a side effect of `worktree remove`; branch deletion is
  a separate, explicit step the user must request.
- Use `--force` (remove) only after confirming with the user; it discards
  uncommitted work in that worktree.
- Do not push automatically. The first push is left to the user so the new remote
  branch is created deliberately (`git push -u`).
- Confirm the target path with `git worktree list` before remove/move/lock.

## References

- [Upstream and --no-track details](references/no-track.md)
