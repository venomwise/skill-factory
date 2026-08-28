# Code Review

<!--
固定标题、字段名、ID 和枚举保持原样。叙述内容使用 Review Manifest 声明的 Report language；
命令、路径、代码符号、日志和引用原文不翻译。
-->

## Review Manifest

- **Mode**: Full Review | Closure Review
- **Review kind**: general | spec-backed
- **Report language**: `<用户请求的主要语言及选择依据>`
- **Base ref / commit**: `<ref>` / `<resolved commit>`
- **Head commit**: `<commit>`
- **Merge-base**: `<commit>`
- **Commit list**: `<command and commits>`
- **Committed diff**: `<exact command>`
- **Staged diff**: `<exact command>`
- **Unstaged diff**: `<exact command>`
- **Untracked files**: `<ordered paths and content-read method>`
- **Standards sources**: `<AGENTS.md and other repository guidance>`
- **Spec source**: `<specs/<topic>/design.md>` | None
- **Historical findings**: `<all CR-### IDs and statuses; mark unresolved findings>` | None
- **Validation evidence**:
  - `<command>: <result>`
- **Diff fingerprint**: `sha256:<value>`
- **Scope notes**: `<未纳入范围的内容及原因>`

## Results

- **Standards Result**: Blocked | Reject | Changes Requested | Pass with Notes | Pass
- **Spec Result**: Blocked | Reject | Changes Requested | Pass with Notes | Pass | Not Reviewed
- **Overall Result**: Blocked | Reject | Changes Requested | Pass with Notes | Pass
- **Merge readiness**: No-Go | Conditional | Go
- **Summary**: `<基于证据的简短结论>`

## Findings

### CR-001: <能指出具体问题对象的标题>

- **Axis**: standards | spec
- **Severity**: blocker | major | minor
- **Status**: open | reopened | resolved | rejected | accepted-risk | deferred
- **Origin**: full-review | fix-regression | dependency-unlocked | baseline-miss | context-change
- **Rule / Spec source**: `<AGENTS.md 章节、AC ID、接口契约或运行证据>`
- **Location**: `<path:line 或唯一符号>`
- **Issue**: `<具体错误、缺失或风险条件>`
- **Evidence**: `<代码、命令结果或可复现反例>`
- **Impact**: `<可观察后果>`
- **Recommendation**: `<修复方向；评审阶段不直接修改代码>`
- **Closure test**: `<能够证明问题已解决的命令或反例>`
- **Related findings**: `<CR-###>` | None

## Axis Summary

| Axis | Result | Findings | Key evidence |
|------|--------|----------|--------------|
| Standards | `<result>` | `<CR-### or None>` | `<规则与验证证据>` |
| Spec | `<result or Not Reviewed>` | `<CR-### or None>` | `<AC 与测试证据，或跳过原因>` |

## Accepted / Deferred Risks

| Finding | Status | Decision source | Revisit when |
|---------|--------|-----------------|--------------|
| `<CR-###>` | accepted-risk / deferred | `<用户决定>` | `<重新评估条件>` |

## Closure

| Finding | Axis | Result | Evidence |
|---------|------|--------|----------|
| `<CR-###>` | standards / spec | resolved / reopened | `<本轮闭合证据>` |

## Recommended Next Step

- **Suggested action**: `<处理 findings、补充评审证据或继续交付>`
- **Findings to discuss or fix**: `<CR-### and reason>` | None
- **Decisions needed**: `<需要用户取舍、接受、推迟或提出异议的事项>` | None
- **Expected scope**: `<预计修改的文件、符号和边界>` | None
- **Validation**: `<Closure tests and project checks>` | None
- **After remediation**: `<使用全新 Manifest 和只读子 Agent 自动执行 Closure Review>` | None
