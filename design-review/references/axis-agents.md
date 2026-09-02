# Axis Agents Contract

Full Review 的维度子代理启动契约。并行的只读子代理压缩证据探索和维度检查的墙钟时间；
主代理保留模式选择、确定性 preflight、候选裁决、严重性校准和 finding 冻结。
Closure Review 在主线程执行，不启动子代理。

## Capability gate and fallback

1. Full Review 开始时直接尝试并行启动两个子代理，不先自评运行时是否支持。可用机制包括
   运行时提供的代理工具，以及任何能接收提示词并返回文本报告的嵌套只读进程。
2. 只有实际尝试失败后才允许串行降级：主代理按 `SKILL.md` STEP 3 的探索优先级和 rubric 的
   全部维度串行执行，并在响应中记录尝试过的机制和具体失败原因。子代理是提速手段，不是
   评审准入条件；启动失败不构成 `Blocked`，也不得降低检查范围。
3. 任一子代理失败、超时或只返回无证据内容时，主代理串行补跑该维度组。

本契约不绑定具体运行时的工具名；任何满足「启动只读代理并接收文本报告」的机制都适用。
「工具清单里没有名字显眼的代理工具」不是失败证据，未经实际尝试不得走串行路径。

## Division of labor

**项目证据代理（D4/D5）：**

- 范围：项目契合度与盲点。核实设计的技术栈、依赖版本、模块边界、命名和数据约定是否与
  项目事实一致；核实设计引用的既有列、字段、错误码和既有行为；按 rubric 的 D5 清单
  选择适用项，并说明为什么适用于当前 Goals 或变更边界。
- 输入：`design.md` 路径、review rubric 路径、项目根目录，以及主代理已解析的
  Project Traceability 信息（如存在）。
- 允许动作：只读项目文件、配置和最近 10 个提交；设计依赖真实数据库时运行只读
  `db-explorer`。
- 禁止动作：修改任何文件；运行改变工作区的命令；分配 `F-###`；给出 Verdict 或
  readiness 结论。

**设计一致性代理（D1/D2/D6/D7）：**

- 范围：完整性、可用性、过度设计与优化。建立 Goal -> Solution -> AC 覆盖关系；
  检查 Data Model 与 Interfaces 的 Change Summary、详细子章节和成对契约共享约束；
  找出无法追溯到 Goals 的机制，以及可复用更多或更简单的具体路径。
- 输入：`design.md` 路径、canonical template 路径、review rubric 路径。
- 允许动作：只读设计文档、模板、项目文件和设计引用的现有 specs。
- 禁止动作：同上。

两个子代理必须同时启动，不先后执行；各自只使用主代理给定的输入路径，不自行扩大
评审范围。

## Candidate format

每条候选项包含：

- **Dimension**: <D1-D7，可多选>
- **Location**: `<design.md 章节和/或项目文件位置>`
- **Issue**: <在什么条件下，哪个设计对象存在什么错误、缺失或风险>
- **Evidence**: <可核实的设计位置、代码、数据库或契约证据>
- **Suggested severity**: blocker | major | minor
- **Related candidate**: <同根因的其他候选项，可选>

候选项叙述使用评审报告的主要语言；命令、路径、代码符号和引用原文不翻译。

## Adjudication by the main agent

子代理不冻结 findings。主代理独立读取候选证据并执行：

1. 候选是否指向评审范围且证据可核验？IF 否，THEN 补证据或丢弃。
2. 建议严重性是否与 rubric 的实际影响匹配？IF 否，THEN 校准并记录理由。
3. 多个候选是否同一根因？IF 是，THEN 合并为一个，跨维度引用。
4. 候选是否与有效 Decision 或用户已确认的边界冲突？IF 是，THEN 丢弃并说明。

裁决通过后分配稳定 `F-###`，再进入 `SKILL.md` 的 STEP 6。
