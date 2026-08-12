# <Topic> Design Review

> 按 7 维度 rubric 对 `<path/to/design.md>` 的评审。结构标题为英文；发现内容使用设计当前的语言。

## Verdict

- **Overall**: <Pass | Revise | Reject>
- **spec-plan readiness**: <Go | Conditional | No-Go>
- **Findings**: <N> Blocker · <N> Major · <N> Minor
- **Summary**: <一句话——作者应该知道的最重要的事>

## Scope Reviewed

- **Reviewed file**: `<path/to/design.md>`
- **Project context consulted**: <探索过的关键文件/目录——例如 `CLAUDE.md`、`src/...`、近期提交>
- **Rubric**: 7 dimensions (D1 Completeness · D2 Usability · D3 Document Conformance · D4 Project Fit · D5 Blind Spots · D6 Over-Engineering · D7 Optimization)

## Findings

发现携带 ID（`B#` Blocker、`M#` Major、`m#` Minor），以便 Dimension Summary 引用。按严重性分组；空分组省略。

### Blockers

#### [B1] <简短标题> — <维度，例如 D1 Completeness>

- **Location**: `design.md §<section>` <和/或项目文件路径>
- **Issue**: <什么是错误/缺失的，以及支持它的证据>
- **Recommendation**: <具体修复，或作者必须解决的问题>

### Major

#### [M1] <简短标题> — <维度>

- **Location**: `design.md §<section>` <和/或项目文件>
- **Issue**: <基于证据的描述；盲点表述为"值得检查">
- **Recommendation**: <具体修复或考虑事项>

### Minor

#### [m1] <简短标题> — <维度>

- **Location**: `design.md §<section>`
- **Issue**: <描述>
- **Recommendation**: <建议改进>

## Dimension Summary

状态：✓ 无问题 · △ 有发现待处理 · ✗ 存在阻塞问题。

| Dimension | Status | Finding refs |
|-----------|--------|--------------|
| D1 Completeness (完整性) | <✓/△/✗> | <B1, M2 / —> |
| D2 Usability (可用性) | <✓/△/✗> | <—> |
| D3 Document Conformance (规范性) | <✓/△/✗> | <—> |
| D4 Project Fit (符合项目规范) | <✓/△/✗> | <—> |
| D5 Blind Spots (盲点) | <✓/△/✗> | <—> |
| D6 Over-Engineering (过度设计) | <✓/△/✗> | <—> |
| D7 Optimization (优化点) | <✓/△/✗> | <—> |

## Recommended Next Step

<从发现和裁决推导：
- 有发现 → 用这份 `review.md` 和被评审的 `design.md` 运行 `/design-refine` 来解决它们。
- Reject → refine 之后重新运行 `/design-review`。
- 无发现 → 进入 `spec-plan`。>

## Accepted / Deferred Risks

<可选。评审中标记出来、但作者可能故意延后的项目，附理由。若没有则省略本节。>
