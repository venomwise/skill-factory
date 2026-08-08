---
name: brainstorming
description: Turn ideas into validated designs and specs through collaborative dialogue. Use before creating features, components, or behavior changes when requirements need scoping and trade-offs. Not for bug fixes, typo-only changes, or clear single-step execution tasks.
---

# Brainstorming Ideas Into Designs

## When to use

- 创建新功能或组件
- 修改现有系统行为
- 需求不明确，need scoping
- 存在多种实现方式

## When not to use

- 有明确根因和步骤的 bug 修复
- 仅涉及拼写或格式的修改
- 明确的单步执行任务
- 用户明确拒绝设计流程

## Inputs

- 用户的想法或目标（可能模糊）
- 现有项目上下文（文件、文档、近期提交）

## Outputs

- 已验证的设计文档，位于 `specs/<topic>/design.md`
- 根据复杂度匹配的下一步：简单设计建议直接实现，或推荐中等/复杂设计先运行 `/design-review`

## Workflow

1. 按优先级探索项目上下文：README/CLAUDE/AGENTS.md，项目配置（package.json / pyproject.toml / Cargo.toml / go.mod / pom.xml），入口点，然后是近期提交（需要时最多 10 个）。Stop when 你能用 2-3 句话描述项目的目的、技术栈和结构。If 项目复杂或不熟悉，consult [the guide](references/brainstorming-guide.md) for exploration priorities。
2. 在接受请求的形态之前，先验证请求与当前项目的契合度。检查请求将触及的现有行为、数据和接口，基于真实证据而非假设进行验证——当请求依赖存储的数据或 schema 时，use available tools such as `db-explorer`。If 请求看起来冗余、与现有内容冲突、解决了错误的问题、或有更简单的替代方案，pause，告诉用户你的发现和推荐路径，然后询问是否继续。
3. 评估范围。If 请求跨越多个独立子系统，propose 一个拆解方案，包含子项目、依赖关系和构建顺序，然后请求用户确认后再继续。Only brainstorm 第一个确认的子项目；每个都有自己的 spec 周期。
4. 通过迭代轮次与用户进行 brainstorming。Do not 预先承诺你会问多少个问题或在探索（步骤 1-3）完成之前就决定要问什么主题——问题的数量和主题由你实际发现的空白和盲点决定，not 预先估计。首先诊断请求的模糊程度（见 guide 中的标准），然后应用匹配的技术：
   - Problem unclear -> reframe：帮助用户阐明他们实际要解决的问题，not just 他们想要什么功能。
   - Direction unclear -> explore：在缩小范围之前，先勾勒几个有代表性的方向，使用"假设"场景和视角切换等技术。
   - Boundaries unclear -> scan：使用 guide 中与领域相关的盲点检查清单，帮助用户考虑潜在的空白。
   - Solution unclear -> 技术方案仍需比较，but first 在进入步骤 5 之前确认意图、优先级和约束。请求中充满具体技术术语 does NOT mean 意图已确认；熟悉技术 is not 知道用户想要什么。
   对于非平凡或模糊的请求，至少浮现一个值得确认的假设和一个值得用户考虑的潜在盲点。
   每轮只问一个问题并等待答案——这是关于节奏，not 对你总共问多少个问题的上限。当多选题能帮助用户在有意义的选项间做选择时，prefer 多选题。Don't 在同一轮中问一个问题然后用假设自己回答它。将未知的意图、优先级、权衡偏好和验收标准视为需要询问的内容，not 需要假设的——这些 only 存在于用户的头脑中，cannot 通过检查项目来获取。
   持续提问直到满足退出条件。Exit condition：用户能清楚地陈述他们想要什么、为什么想要、以及他们明确不想要什么；关键约束和成功标准已知或明确记录为假设；问题陈述在最近 2 次交流中保持稳定。离开此步骤前，verify：你能用用户自己的话用 2-3 句话总结用户的意图、优先级和约束吗？If not，ask more questions。
5. 收敛并提议。Admission gate — 进入前 check 所有三个条件：
   1. 在本次会话中，你是否至少向用户提出了一个澄清性问题？If no，回到步骤 4。
   2. 用户现在能陈述他们想要什么、为什么想要、以及明确不想要什么吗？If no，回到步骤 4 并 ask more questions。
   3. 关键约束和成功标准已知或明确记录为假设了吗？If no，回到步骤 4 并 ask more questions。
   Note：请求在技术上很详细 does not satisfy 这些条件。Technical detail 对用户的意图、优先级或权衡偏好没有任何说明。
   Once 所有条件都满足，summarize brainstorming 揭示的内容：refined problem statement、request-fit findings、challenged assumptions、discovered blind spots、以及 trimmed scope。Then propose 1-3 个方案及其权衡。Lead with 你的推荐。If only 一个方案可行，explain 为什么其他方案被排除。Ask 用户确认问题总结并选择一个方案后再继续。
6. Present 与复杂度匹配的设计。For simple projects，一次性 present 完整设计并请求批准。For moderate or complex projects，分段 present 并在每段后请求批准。Cover architecture、components、data flow、error handling、and testing。
7. 将设计文档写入 `specs/<topic>/design.md`，使用 `assets/design-doc-template.md` 中的模板。Preserve 模板的英文标题和结构标签，同时用用户当前的语言编写章节内容。从对话中推断该语言；don't 仅为确定语言而询问。使用从项目或功能名称派生的 kebab-case 命名 `<topic>`（例如 `user-auth`、`payment-integration`）。If ambiguous，与用户确认路径。从步骤 5 中产生的方案比较填充 Decision Record 章节：记录实际权衡过的每个选项及其关键权衡，然后是选定的方案及其获胜的具体原因（或者，if only 一个方案可行，为什么其他方案被排除）。这个比较是在步骤 5 中实时生成的，是值得保留供后续审查的主要内容，so don't drop it；if 它已从最近的对话中滚出，回去恢复它而非从记忆中重建。Do not fabricate 从未讨论过的替代方案。Open Questions must only contain 已向用户提出的未解决问题；if none 保留，写明没有未解决的问题。
8. 根据实现风险对完成的设计进行分类，使用设计本身而非初始请求的表面简单性。**Only** 当以下所有条件都为真时才视为 **simple**：变更局限于一个组件或行为；它遵循已建立的项目模式；它 not introduces 数据迁移、外部集成、安全/权限边界、并发/可靠性问题、或公共契约变更；测试和回滚是直接的；并且 no open question 保留。If any 条件为假或不确定，classify it as **moderate or complex**。See [the guide](references/brainstorming-guide.md) for examples。
9. 根据该分类路由下一步：
   - **Simple:** 说明设计路径并简要解释为什么它是低风险的，然后询问是否开始直接实现已编写的设计。Do not require 单独的文档审查轮次或 `spec-plan`。If 用户同意，该响应既批准了设计又授权了实现；继续实现。If 用户给出反馈，update the doc，当反馈改变意图或方案时返回步骤 4 或 5，并在再次询问前重新评估复杂度。
   - **Moderate or complex:** 说明设计路径和分类背后的具体风险信号，然后建议用户在规划或实现之前对该路径运行 `/design-review`。Do not invoke `design-review`、`spec-plan`、或实现技能 automatically。If 用户直接提供反馈，按简单路径处理，然后重新评估复杂度。

## Verification

- [ ] `specs/<topic>/design.md` 存在
- [ ] 设计文档包含 `assets/design-doc-template.md` 中的模板标题
- [ ] 设计文档保持模板标题为英文，章节内容使用用户当前的语言
- [ ] 设计涵盖 architecture、components、data flow、error handling、and testing
- [ ] If brainstorming 浮现了值得注意的发现，它们被记录在设计文档的 Discovery 章节中
- [ ] Decision Record 捕获了步骤 5 中比较的方案以及选择方案的理由（或当 only 一个方案可行时为什么其他方案被排除）
- [ ] Open Questions only contains 已浮现的未解决问题，或明确说明 none remain
- [ ] 使用步骤 8 的风险信号从完成的设计评估复杂度
- [ ] Simple design 以询问直接实现确认结束；moderate or complex design 以 `/design-review` 推荐结束

## Safety & guardrails

- No implementation before explicit confirmation。Only simple-design 路由可以直接进行，并且 only 在用户同意实现已编写的设计后。
- Moderate and complex designs stop at 推荐 `/design-review`；do not auto-continue 进入 review、planning、or implementation。
- Even simple projects require a design；当范围较小时保持简短。
- Scale the process to complexity。Simple projects 可能在几次交流中完成步骤 4-6；don't 用不必要的仪式填充流程。
- Pace questions correctly。每轮 ask one question 并等待答案（节奏），but keep going 直到意图和约束清楚——问题的数量由你发现的空白决定，not 为了最小化本身而最小化。Full rules live in step 4。
- YAGNI at the right time。While diverging（浮现选项和盲点），freely explore 并 don't 过早地驳回想法。While converging and designing，ruthlessly cut - remove 任何对核心目标 not essential 的内容。
- Design for isolation。将系统分解为具有一个明确目的和明确定义接口的单元。See [the guide](references/brainstorming-guide.md) for details。

## References

- [Detailed brainstorming guide](references/brainstorming-guide.md)
- [Design doc template](assets/design-doc-template.md)
