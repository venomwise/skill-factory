# Shadow Review Contract

`clarifying` 在交付中等或复杂设计前执行的影子评审契约。影子评审用与 `design-review`
相同的 rubric 提前发现高信号问题，让真实首轮评审从发现轮变成确认轮。它不是独立评审，
不写入 `review.md`，不分配 `F-###`，也不替代后续任何轮次。

## When to run

- 设计分类为中等或复杂，`design.md` 已写入，并通过文档语言关卡和结构校验器。
- 运行环境能启动隔离的只读子代理。无法启动时跳过，并在最终响应说明；
  不由主流程自我评审替代，因为同一上下文自查会系统性漏检自己刚写下的内容。
- 简单设计没有设计文档，不执行。

## Launch contract

子代理启动指令包含：

- `design.md` 的实际路径。
- canonical template 路径：本技能的 `assets/design-doc-template.md`。
- review rubric 路径：`../design-review/references/review-rubric.md`。
- 项目根目录；子代理按需只读项目代码、配置和数据库工具输出，核实设计中的项目事实。

子代理约束：

- 严格只读，不修改任何文件，不运行构建、测试或任何改变工作区的命令。
- 每条候选项必须能指向设计位置或项目证据；无法指出证据的内容不输出。
- 不计算 Verdict，不写报告，不运行 `design-review` 的 lifecycle 协议；只返回候选项列表。

## Review focus

按 rubric 维度做压缩检查，优先输出会改变评审轮次的问题：

1. Goal -> Solution -> AC 覆盖：核心 Goal 是否都有对应组件和可测试 AC，
   正常流、错误流和关键边界是否覆盖。
2. 契约完整性：跨边界调用是否都有 Contract ID 与请求、响应、错误定义；
   preview/commit、create/update 等成对契约的身份、容量、校验和事务约束是否一致。
3. Data Model 与 Interfaces：Change Summary 是否覆盖全部变化，迁移、回填和回滚是否闭环。
4. D5 盲点：设计是否覆盖适用的容量、并发、权限、失败和运营风险；
   声明的既有列、字段、错误码和行为是否与项目事实一致。
5. 过度设计：是否存在无法追溯到 Goal 的组件、配置或机制。

## Candidate format

每条候选项包含：

- **Location**: `<design.md 章节或项目文件位置>`
- **Issue**: `<在什么条件下，哪个设计对象存在什么错误、缺失或风险>`
- **Evidence**: `<可核实的设计位置或项目证据>`
- **Suggested severity**: blocker | major | minor

## Triage by the main agent

主代理逐条核验候选项后分类处理：

1. 证据不足，或与设计中已确认的边界冲突：丢弃。
2. 改变意图、范围、方案选择或需要用户取舍：呈现给用户并等待确认；
   确认后回到 `SKILL.md` 的 STEP 4 或 STEP 5。
3. 一致性缺失、传播缺口、格式或模板问题：作为直接修复应用。

全部处理完成后重新执行文档语言关卡和结构校验器，并在最终响应中汇总保留、修复和丢弃的
候选项。影子评审的结果不写入 `design.md` 或 `review.md`。
