---
name: spec-plan
description: Create requirements.md and tasks.md for a project spec, with traceability and status tracking. Use when the user asks for a spec, requirements, task plan, or execution breakdown (e.g. "做个 spec", "拆任务", "出执行计划", "write a spec", "plan tasks"). Requires an approved specs/<topic>/design.md as input; if missing, route to the brainstorming skill first.
# metadata:
#   short-description: Spec plan (requirements + tasks)
---

# Spec Plan Skill

## When to use

- 从已批准的 `design.md` 创建或刷新 `requirements.md` 和 `tasks.md`
- 将设计转化为准确、可测试的验收标准和可执行的任务
- 建立从任务到需求、从需求回溯到已批准设计的可追溯性

## When not to use

- 你只需要一个快速 TODO 列表或单个文件编辑
- 工作已经在现有 spec 中捕获，只需要实现
- 还没有已批准的 `design.md`；应先创建或请求设计，而非基于假设生成 requirements/tasks

## Inputs

- 项目名称和目标目录
- 已批准的设计文档：`specs/<topic>/design.md`
- 可选的澄清说明：only when `design.md` 存在歧义、不完整或内部不一致时

## Outputs

- `requirements.md`
- `tasks.md`
- 最终响应包含推荐的下一步：让用户对生成的 `specs/<topic>/tasks.md` 运行 `spec-exec`；do not 自动开始实现

## Workflow

1. 解析目标设计文档。
   - If 用户提供了 `specs/<topic>/design.md`，使用它。
   - If 用户 did not provide 设计路径，精确回复：`请指定 design 文件路径（例如 specs/<topic>/design.md）。` Then 结束流程。Do not 搜索、列举或推断设计。
   - If 用户提供的路径 does not exist，结束流程并推荐先运行 **`brainstorming`** skill 以在 `specs/<topic>/design.md` 产出已批准的设计。
2. 从选定的 `specs/<topic>/design.md` 确认目标目录和项目名称。
3. 打开选定的 `design.md` 并从散文内容推断其正文语言，忽略 Markdown 标题、代码、文件路径、标识符和引用/模板标签。If 语言混合，使用用户撰写的说明性内容的主导语言。Do not 仅为确定语言而询问。
4. 使用 `assets/requirements.template.md` 起草 `requirements.md`。
   遵循模板中的 HTML 注释以获取内容深度和覆盖范围指导。
   HTML 注释是创作指令 - Do NOT 将它们包含在最终输出中。
   精确保留模板的英文结构：Markdown 标题、固定的 section/schema 标签、`Requirement N`、`User Story`、`Acceptance Criteria`，以及验收标准控制词如 `WHEN`、`THEN`、`IF` 和 `SHALL`。用推断的 `design.md` 正文语言编写生成的内容，包括引言、术语表定义、需求标题、用户故事文本和验收标准条件/行为文本，unless 使用代码名称、产品名称或固定标识符。
   将 `design.md` 视为需求的权威来源。将已批准的设计转化为准确、可测试的验收标准，without adding、omitting 或 changing 预期行为。在适用的地方反映已批准的行为和约束，but do not 将设计理由、示例、备选方案或未来想法转化为硬性需求，unless 设计明确要求它们。
5. 使用 `assets/tasks.template.md` 起草 `tasks.md` 并链接到 `requirements.md`。
   精确保留模板的英文结构：Markdown 标题、checkbox/task 编号语法、`Phase`、`Checkpoint`、可选的 `*` 标记语法和 `_Requirements:` 标签。用推断的 `design.md` 正文语言编写生成的内容，包括概述、固定的 `Phase N:` 前缀后的阶段标题、任务标题、任务详细列表、检查点验证详情和注释，unless 使用代码名称、产品名称、文件路径、命令名称或固定标识符。
   为每个功能阶段包含测试任务，并在关键里程碑添加 Checkpoint 阶段。Checkpoints 是执行 agent 的验证任务，not 用户批准门；编写具体的验证步骤和阻塞条件，instead of 询问是否继续。
   用 `- [ ]*`（星号紧跟在右括号后）标记可选的 phases/tasks。非必要步骤如测试任务、验证任务、总结/文档收尾任务和 nice-to-have 功能 MUST 使用此标记 — never 使用文本标签如 "可选" 或 "Optional:" 代替。
6. 确保每个任务引用一个或多个需求 ID（用于可追溯性）。
7. 在最终确定之前使用链式验证：验证 `requirements.md` 准确实现 `design.md` without drift，然后验证 `tasks.md` 覆盖并引用 `requirements.md`。
8. If `design.md` 存在歧义、不完整或内部不一致，在最终确定前询问用户。Do not 通过发明需求来填补空白。

## Verification (self-check before finalizing)

Before presenting the final output，scan both files and confirm each item:

- [ ] `design.md` exists and has been read before drafting `requirements.md` or `tasks.md`。
- [ ] The body language for `requirements.md` and `tasks.md` was inferred from `design.md` prose，not from template headings。
- [ ] Template English structure is preserved while generated content uses the inferred `design.md` body language。
- [ ] `requirements.md` contains Introduction, Glossary, and numbered Requirements sections。
- [ ] Every requirement is traceable to approved behavior or constraints in `design.md`。
- [ ] `requirements.md` accurately implements `design.md`: no added behavior, omitted required behavior, changed semantics, or contradictory architecture。
- [ ] Important approved behaviors and constraints from `design.md` are reflected as testable acceptance criteria where applicable。
- [ ] Design rationale, examples, alternatives, and future ideas are not converted into hard requirements unless explicitly required by `design.md`。
- [ ] `tasks.md` links to `requirements.md` and every task includes a `_Requirements: ..._` line。
- [ ] Requirement IDs referenced in `tasks.md` exist in `requirements.md`。
- [ ] Each requirement includes acceptance criteria covering normal flow, error flow, and boundary conditions。
- [ ] Tasks include specific file paths and function/class names。
- [ ] At least one Checkpoint task exists between major phases and describes concrete validation steps，not user approval。
- [ ] Phase headings use `- [ ] N. Phase N:` checkbox format，not markdown headings (`###`)。
- [ ] Task descriptions are indented bullet points under the task title line。
- [ ] Every non-essential step uses `- [ ]*` marker，including test tasks, verification tasks, summary/documentation wrap-up tasks, and nice-to-have features。No task uses text labels like "可选", "Optional", or "(Optional)" as a substitute。
- [ ] Sub-tasks under an optional Phase inherit optionality and do not need their own `*`。
- [ ] Requirements references use `N.M` format (e.g., `_Requirements: 1.1, 2.3_`)，not `RN` format。

## Safety & guardrails

- Never generate `requirements.md` or `tasks.md` without an approved `design.md`。
- Never auto-select a `design.md`。The path must come from the user。
- Do not invent requirements, broaden scope, or reinterpret design intent。Ask for clarification when the design is ambiguous or incomplete。
- Keep requirements testable and phrased as acceptance criteria。
- Do not mark tasks as complete unless the work is done。

## References

- [Requirements template](assets/requirements.template.md)
- [Tasks template](assets/tasks.template.md)
