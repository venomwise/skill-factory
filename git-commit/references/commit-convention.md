# Commit Conventions

## Emoji list

- `:art:` Improve structure or formatting
- `:zap:` Improve performance
- `:fire:` Remove code or files
- `:bug:` Fix a bug
- `:ambulance:` Critical hotfix
- `:sparkles:` Introduce new features
- `:memo:` Add or update documentation
- `:rocket:` Deploy something
- `:lipstick:` Add or update UI and style files
- `:tada:` Begin a project
- `:white_check_mark:` Add, update, or pass tests
- `:lock:` Fix security issues
- `:closed_lock_with_key:` Add or update secrets
- `:bookmark:` Release or version tag
- `:rotating_light:` Fix compiler or linter warnings
- `:construction:` Work in progress
- `:green_heart:` Fix CI build
- `:arrow_down:` Downgrade dependencies
- `:arrow_up:` Upgrade dependencies
- `:pushpin:` Pin dependencies to specific versions
- `:construction_worker:` Add or update CI build system
- `:chart_with_upwards_trend:` Add or update analytics or tracking
- `:recycle:` Refactor code
- `:heavy_plus_sign:` Add dependencies
- `:heavy_minus_sign:` Remove dependencies
- `:wrench:` Add or update configuration files
- `:hammer:` Add or update development scripts
- `:globe_with_meridians:` Internationalization and localization
- `:pencil2:` Fix typos
- `:poop:` Write bad code that needs improvement
- `:rewind:` Revert changes
- `:twisted_rightwards_arrows:` Merge branches
- `:package:` Add or update compiled files or packages
- `:alien:` Update code due to external API changes
- `:truck:` Move or rename resources
- `:page_facing_up:` Add or update licenses
- `:boom:` Introduce breaking changes
- `:bento:` Add or update assets

## Type list

- `feat` New feature
- `fix` Bug fix
- `docs` Documentation only
- `style` Formatting only, no logic change
- `refactor` Refactor, no new feature and no bug fix
- `test` Add or update tests
- `chore` Build process or auxiliary tooling

## Scope rules

- `src/mcp/**` -> `mcp`
- `src/common/**` -> `common`
- `doc/**` -> `doc`
- `tests/**` -> `tests`
- Mixed or unclear -> `repo`

## Subject rules

- Start with a verb in present tense
- Keep within 50 characters
- Do not end with a period

## Body rules

- Explain why and what changed
- Wrap at about 72 characters per line
- Use `- ` bullets, one key point per line

## Footer rules

- Only use the footer for two cases: breaking changes or closing issues
- For breaking changes, start with `BREAKING CHANGE:` and include the change,
  the reason, and a migration guide:

```text
BREAKING CHANGE: 隔离作用域绑定定义已更改

要迁移代码，请参考以下示例：

之前：

scope: {
  myAttr: 'attribute',
}

之后：

scope: {
  myAttr: '@',
}

删除的 `inject` 对于指令来说通常没有太大用处，因此应该没有代码在使用它
```

- To close issues, use `Closes #123`
- You can close multiple issues at once:

```text
Closes #123 #234 #345
```

## Good Examples

With emoji:

```text
:sparkles: feat(mcp): add local server health endpoint

- add /health endpoint for readiness checks
- return minimal status payload

Closes #123
```

Without emoji:

```text
fix(common): handle empty log path

- guard when log directory is missing
- keep default console output

Closes #456
```
