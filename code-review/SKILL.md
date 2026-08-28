---
name: code-review
description: >
  对代码改动、提交或本地 pull request diff 执行独立、证据化的双轴评审，并在 spec
  场景写入 code-review.md。当用户说“代码评审”“review code”“检查实现质量”
  “检查是否符合 AGENTS.md”“review PR”，或 spec-exec 完成高风险实现后需要独立
  质量关卡时调用。默认先只读评审；报告完成且用户确认统一修复计划后，可处理 findings，
  集中修改代码，并用全新 Manifest 和子 Agent 执行 Closure Review。
---

# Code Review

对同一份冻结 diff 分别执行 Standards 和 Spec 评审。评审阶段由主 Agent 固定证据、启动
互不污染上下文的子 Agent、核验候选项并统一裁决；子 Agent 不修改代码，也不分配
`CR-###`。报告完成后，只有用户明确批准统一修复计划，主 Agent 才进入隔离的修复阶段。

本 skill 使用严格模式。运行环境无法启动子 Agent 时，不在主 Agent 中模拟双轴评审，
直接报告 `Blocked / No-Go`。

## When to use

- 评审工作区改动、commit range、分支差异或本地 pull request diff。
- 检查实现是否满足 `AGENTS.md`、architecture 约束、Acceptance Criteria 和测试要求。
- `spec-exec` 完成跨层、公共契约、迁移、并发、安全或多模块耦合改动后执行独立质量关卡。
- 讨论并处理本 skill 已报告的 `CR-###` findings。
- 修复 `CR-###` findings 后复核是否真正闭合。

## When not to use

- 评审 `design.md`，应使用 `design-review`。
- 没有现有评审报告、只要求直接开发或修改代码；本 skill 不替代常规实现流程。
- 没有代码差异的仓库级 architecture 盘点。

## Inputs / Outputs

**Inputs:**

- 用户提供的 diff、commit range、base ref、工作区改动或分支。
- 已加载或仓库内可读取的 `AGENTS.md`、`CLAUDE.md`、README 和模块级指导。
- 可选的 `specs/<topic>/design.md`、`tasks.md` 和历史 `code-review.md`。
- 构建、测试、lint、静态分析和运行日志等验证证据。

**Outputs:**

- Spec 场景：同目录 `specs/<topic>/code-review.md`，使用
  [code review template](assets/code-review-template.md)。
- 非 spec 场景：默认在回复中输出评审；仅在用户指定路径时写文件。
- 评审阶段不修改代码；修复阶段只修改统一计划中已批准的代码。
- 始终不修改 `design.md`、`tasks.md` 或任务 checkbox。

## Source contracts

- 实现规范以用户明确要求和适用的仓库指导为准。
- Spec 场景中的功能行为与设计决策以 `design.md` 为准。
- `tasks.md` 只说明计划、追踪关系和进度，不替代 Acceptance Criteria。
- 实际 diff、调用关系和运行证据决定代码当前事实；checkbox 和提交说明不构成通过证据。
- Repository guidance 优先于通用 code smell；smell 只能补充启发式检查，不能覆盖项目约定。

IF 来源冲突会改变正确结论，THEN 将 Overall Result 标记为 `Blocked`，不要自行选择一个版本。

## Workflow

### STEP 1: Resolve intent and review range

先用 2-3 句话确认评审对象、目标和输出位置。按以下顺序解析 review range：

1. 用户提供明确 diff、commit range 或 base ref：直接使用。
2. 工作区存在 staged、unstaged 或 untracked 改动：评审相对 `HEAD` 的完整工作区变化。
3. 工作区干净且当前分支有 upstream：使用 merge base 到 `HEAD`。
4. 无 upstream 但仓库存在唯一默认分支：使用默认分支 merge base 到 `HEAD`。
5. 仍有多个合理范围：给出候选范围和倾向，每轮只问一个问题并等待回答；不要替用户选择。

这是提问节奏，不是总量上限；持续到范围唯一。范围必须包含 untracked 文件。IF 没有任何差异，
THEN 报告“没有可评审改动”并 STOP。

同时确定 `Report language`：新报告使用用户当前请求的主要语言；用户明确指定语言时以指定
语言为准。Closure Review 默认沿用历史报告声明的语言；历史报告没有该字段时使用当前请求
的主要语言。只有用户明确要求切换语言时，才调整历史 finding 的叙述字段。

固定结构标题、字段名、`CR-###`、Axis、Severity、Status、Origin、Result 和 readiness 枚举保留
模板定义；finding 标题、Issue、Evidence、Impact、Recommendation、Closure test、Summary、
Scope notes 和下一步说明使用 `Report language`。命令、路径、代码符号、日志和引用原文不翻译。

### STEP 2: Discover authoritative sources

完整读取 diff 和所有新增文件，再按风险读取被改符号、调用方、被调用方、公共接口、数据模型、
配置、相关测试和必要 git history。定位所有改动路径适用的仓库指导。

按以下顺序查找 Spec source，命中后停止：

1. 用户明确提供的 issue、spec 或 `design.md`。
2. 历史 `code-review.md` 中记录的 Spec source。
3. 当前请求、`spec-exec` handoff 或 `tasks.md` 指向的同目录 `design.md`。
4. 与 branch 或 topic **完全同名**的 `specs/<topic>/design.md`。
5. commit message 中的 issue 引用，但仅在当前环境已有 issue tracker 读取能力时使用。
6. 仍未找到时只问一次来源；不要模糊匹配或从多个候选中自行挑选。

用户请求来自 spec 工作流、明确引用 AC/`tasks.md`，或历史报告记录过 spec 时，属于
spec-backed review。此类评审缺失 Spec source 时报告 `Blocked / No-Go`。一般评审经用户确认
确实没有 spec 时，允许跳过 Spec Agent，并将 Spec Result 设为 `Not Reviewed`。

IF 历史 `code-review.md` 存在，THEN 完整读取其中的 review range、所有 `CR-###`、Closure
和 accepted/deferred risks。不要凭记忆补写历史 finding。

### STEP 3: Select review mode

- 不存在历史 `code-review.md`：Full Review。
- 用户明确要求完整重审：Full Review。
- review range、公共契约、数据模型、权限、安全、迁移或主要 architecture 边界发生变化：Full Review。
- 仅修复历史 findings，且影响范围没有触发上述变化：Closure Review。
- 历史报告为 `Pass / Go` 且 diff fingerprint 没有变化：说明无需重复评审并 STOP。

Closure Review 只检查历史未闭合 finding、修复 diff、受影响调用链和必要回归测试。
新 finding 必须注明 `fix-regression`、`dependency-unlocked`、`baseline-miss` 或
`context-change`，并提供因果证据；无法归因的自由优化建议不进入 Findings。

### STEP 4: Validate once and freeze the Review Manifest

主 Agent 在启动子 Agent 前运行项目提供的安全构建、测试、lint、静态分析和
`git diff --check`。只使用非修复模式；昂贵验证只运行一次。无法运行时记录原因、替代证据
及其对结论的影响。先记录验证前的 review range；IF 命令改变 tracked、staged 或目标 untracked
内容，THEN 报告 `Blocked / No-Go`，不要清理或把副作用静默纳入评审。验证完成后再冻结
Review Manifest；被 git 忽略的普通构建产物不进入 fingerprint。

Manifest 必须记录：

- mode、general/spec-backed 分类，以及所有历史 finding 的 ID 与 Status，并标出未闭合项。
- `Report language` 及其选择依据。
- base ref 与解析后的 commit、`HEAD` commit、merge-base 和 commit list。
- committed、staged、unstaged diff 的精确命令。
- 排序后的 untracked file 列表，以及读取其完整内容的方法。
- Standards sources、Spec source 和验证命令结果。
- diff fingerprint：对上述 refs、commit list、三份 diff 输出、untracked 路径与内容按清单顺序
  串联后计算 SHA-256；报告中写明算法和值。

在冻结前确定唯一的报告输出路径，并把它作为 review artifact 单独记录。历史报告仍须按
STEP 2 完整读取，但该路径的 committed、staged、unstaged 和 untracked 内容不参与实现 diff
或 fingerprint；否则写入 `code-review.md` 会被误判为实现变化。除这一已记录路径外，不得
排除任何 review range 内的文件。

两个子 Agent 必须使用 Manifest 中相同的命令和来源，不得自行扩大或更换 review range。
冻结后主 Agent 和子 Agent 都不得执行会改变工作区、index、refs 或构建产物的命令。

启动前执行以下关卡：

1. 三类 diff 和 untracked 文件是否全部可重放？IF 否，THEN 补全命令后重建 Manifest。
2. 每个改动路径的 Standards source 是否已定位？IF 否，THEN 检查根、父级和模块级指导。
3. Spec source 是否满足 STEP 2？IF 否且属于 spec-backed，THEN 报告 `Blocked / No-Go`。
4. 历史 finding 状态能否从原报告恢复？IF 否，THEN 执行 Full Review migration 并保留迁移说明。
5. 当前运行环境能否启动只读子 Agent？IF 否，THEN 报告 `Blocked / No-Go`，不得降级为单 Agent 评审。

### STEP 5: Run isolated axis agents

Full Review 中，只要 Spec axis 可评审，就必须**同时**启动 Standards Agent 和 Spec Agent；
不得先后执行。一般评审确认无 spec 时仅启动 Standards Agent。Closure Review 只启动存在历史
未闭合 finding 或受修复影响的 axis；两个 axis 都适用时仍须同时启动。

两个 Agent 均只读项目文件和 Manifest，不修改任何文件，不重复运行构建、测试、lint 或静态分析。
允许使用只读命令检查必要调用关系和上下文。候选项的叙述字段使用 `Report language`，不分配
最终 `CR-###`，每项使用：

- **Axis**: standards | spec
- **Rule / Spec source**: `<文件、规则、AC 或契约>`
- **Location**: `<path:line 或唯一符号>`
- **Issue**: `<问题>`
- **Evidence**: `<代码、Manifest 证据或可复现反例>`
- **Impact**: `<可观察后果>`
- **Suggested severity**: blocker | major | minor
- **Closure test**: `<证明修复的命令或反例>`
- **Related candidate**: `<同 axis 的标题或位置，可选>`

**Standards Agent 职责：**

- 检查 repository conformance、architecture、reliability、maintainability、test quality 和 compatibility。
- Full Review 时必须读取 [code smell baseline](references/code-smells.md)；Closure Review 不重新执行
  开放式 smell 扫描，只复核相关历史项和修复影响。
- Repository guidance 覆盖 smell。Smell 是启发式判断，默认最高为 `minor`；只有存在可观察的
  behavior、reliability 或交付影响时，才可建议更高严重性并给出因果证据。
- 工具已覆盖的规则仍以 Manifest 中的实际结果为证据；工具通过时不重复创建 smell finding，
  工具失败时不得因“可由工具检查”而忽略。

**Spec Agent 职责：**

- 检查遗漏、部分实现或错误实现的 requirement 与 AC。
- 检查 AC coverage、scope creep、spec compatibility，以及测试能否作为 AC 的有效证据。
- 每个候选项引用具体 spec 文本或 AC；不能从 spec 推导的问题不进入此 axis。

任一必需 Agent 未启动、失败、返回范围不一致或没有提供可核验证据时，THEN 将对应 Axis Result
和 Overall Result 标记为 `Blocked`，Merge readiness 为 `No-Go`。

### STEP 6: Adjudicate without collapsing axes

子 Agent 返回后，主 Agent 先按 STEP 4 的同一算法重算 diff fingerprint。IF 值发生变化，
THEN 丢弃本轮候选裁决，报告 `Blocked / No-Go` 并要求基于新快照重跑。

主 Agent 独立读取候选证据并执行以下裁决：

1. 候选项是否指向 Manifest 范围内的代码？IF 否，THEN 删除或改为 scope note。
2. 依据、位置、影响和 Closure test 是否可核验？IF 否，THEN 补证据或删除。
3. 建议严重性是否与实际影响匹配？IF 否，THEN 校准后记录理由。
4. 同一 axis 内是否属于同一根因？IF 是，THEN 合并为一个 finding。
5. 跨 axis 候选是否根因相同？IF 是，THEN 仍保留两个 findings，并用 `Related findings` 互链。

通过裁决后才分配稳定 `CR-###`。新 ID 从历史最大编号递增，不复用已关闭编号；每个 finding
保留 `Axis: standards | spec`。Full Review 新 finding 为 `open`。Closure Review 中原反例仍成立时
沿用原 ID 并设为 `reopened`，反例不再成立时设为 `resolved`；同一根因不得换新 ID。

严重性使用独立字段，不编码进 ID：

- `blocker`：代码不能安全交付，例如明确的数据损坏、安全漏洞、不可恢复迁移或构建失败。
- `major`：可观察行为、批准设计、项目强制规范或关键测试存在实质缺陷。
- `minor`：不阻断交付，但有明确证据的局部质量问题。

`accepted-risk` 和 `deferred` 只能来自用户明确决策。`rejected` 只用于用户提出异议后，重新
核验的证据能够证明原 finding 不成立的情形；用户偏好本身不能推翻证据。accepted/deferred
risks 必须以用户决策和历史报告原文为准，不得编造。

### STEP 7: Compute results and write report

每个 axis 独立计算 Result：必需证据不可得为 `Blocked`；`open`、`reopened`、`accepted-risk`
和 `deferred` 参与计算，存在 blocker 为 `Reject`，无 blocker 但存在 major 为
`Changes Requested`，仅 minor 为 `Pass with Notes`；`resolved` 和 `rejected` 不参与计算。
没有参与计算的 finding 时为 `Pass`。一般评审确认无 spec 时，Spec Result 为 `Not Reviewed`。

Overall Result 取已评审 axis 中更严格的结果，顺序为：
`Blocked > Reject > Changes Requested > Pass with Notes > Pass`。`Not Reviewed` 不参与降级。
Merge readiness：`Blocked` 或 `Reject` 为 `No-Go`，`Changes Requested` 为 `Conditional`，
其余为 `Go`。

Spec 场景使用模板写入 `code-review.md`；非 spec 场景按相同结构在回复中输出。Closure Review
在原 finding 上更新 Status 并刷新紧凑 Closure 表，不删除原 Issue、Evidence、Impact 或
Closure test，也不复制全部历史 finding。

小 diff 允许生成短报告，不为凑维度制造 finding。没有 finding 时仍列出 Manifest、每个 axis
核对过的规则与关键证据。

写完后重新读取实际报告并执行文档语言关卡：

1. 所有叙述字段是否使用 `Report language`？IF 否，THEN 只改对应叙述，不改固定字段和枚举。
2. finding 是否写清具体问题、证据、影响和修复方向，没有模板化空话？IF 否，THEN 回到
   STEP 6 核实后重写，不编造证据。
3. 数字、事实、责任主体、规则来源、Severity、Status、Origin、Result、readiness、命令、路径、
   代码符号、日志和引用是否保持不变？IF 否，THEN 恢复原意并缩小改写范围。
4. 语言调整是否新增、删除、弱化或扩大 finding 和 Closure 结论？IF 是，THEN 撤销该调整，
   依据裁决结果重新表述。
5. 调整后 Results、Axis Summary、Findings、Closure 和 Recommended Next Step 是否仍相互一致？
   IF 否，THEN 修正报告后再进入 STEP 8。

### STEP 8: Recommend and pause

- `Blocked`：说明缺少的 review range、规则、spec 来源、子 Agent 能力或必需证据；补齐前不进入修复。
- 存在 `open` 或 `reopened` findings：在 Recommended Next Step 中列出建议处理的 `CR-###`、
  需要用户决定的事项、预计修改范围、验证方式，以及修复后会自动执行 Closure Review。
- 没有未闭合 finding：说明当前 diff fingerprint 已通过，或列出仍然生效的 accepted/deferred risks。

IF 存在 `open` 或 `reopened` findings，THEN 报告写完后询问用户是要“讨论并处理 findings”，
还是“只保留报告”。提出问题后 STOP 并等待回答；即使用户最初笼统要求“评审并修复”，也
不能代替这次基于实际 findings 的确认。IF 没有未闭合 finding，THEN 交付报告后 STOP，不追加
无意义的处理询问。

### STEP 9: Discuss finding dispositions

仅当用户在 STEP 8 选择继续处理时进入本阶段。先完整重读 `code-review.md` 或回复中的评审，
并按 STEP 4 的原算法重算实现 diff fingerprint。IF 值已变化，THEN 原报告可能过期，回到
STEP 2 基于新快照重新评审；不要沿用旧证据直接修复。

按根因组织讨论，同一根因的多个 `CR-###` 可以作为一个讨论项，但每个 finding 保留独立处置：

1. `fix`：只有一种合理修法、且不涉及用户取舍的直接修复可以批量提出并确认。
2. `skip`：必须取得用户理由；当前明确接受记为 `accepted-risk`，推迟处理记为 `deferred`，
   两者都记录 `Decision source` 和可判断的 `Revisit when`。
3. `dispute`：根据用户提供的新信息重新核验证据。现有证据不能裁决时，启动对应 axis 的全新
   只读 Agent；只有证据证明 finding 不成立时才记为 `rejected`，否则仍为未闭合 finding。

涉及实现取舍、`skip` 或 `dispute` 的事项每轮只讨论一个并等待回答。这是对话节奏，不是总量
上限；持续到每个未闭合 finding 都有处置。IF 用户认为 spec finding 对应的需求不再需要，
THEN 暂停代码修复并进入适用的设计修改流程；未更新 `design.md` 或 AC 前不能将其简单记为
`skip`。

讨论期间只整理决定，不修改代码。全部处置明确后，展示一份统一修复计划，其中逐项列出
根因、关联 `CR-###`、处置、修改文件或符号、修改内容、验证命令和剩余风险。询问用户是否
批准整份计划；提出问题后 STOP。用户反馈改变方案时先更新计划，再重新请求统一确认。

### STEP 10: Apply the approved plan

只有用户明确批准 STEP 9 展示的整份计划后才进入本阶段。修改前重算 fingerprint；IF 与
STEP 9 的快照不同，THEN 停止并回到 STEP 2，不把外部改动静默纳入修复。

先依据已批准的处置更新报告中的 `accepted-risk`、`deferred` 或 `rejected`，再集中修改所有
`fix` 项。只修改计划明确列出的范围，不修改 `design.md`、`tasks.md` 或任务 checkbox，也不在
修复阶段自行把 finding 标为 `resolved`。修改每个文件前重新读取实际内容和适用的 repository
guidance，再执行每条 Closure test 和适用的安全项目验证；验证失败或出现新的实现取舍时停止
扩展修改，回到 STEP 9 讨论增量计划。

修复计划的批准不授权破坏性操作、生产操作、真实凭证使用或数据库写入；这些动作仍需用户
单独明确授权。IF 计划中没有 `fix` 项，THEN 更新报告和 Results 后结束，不制造无变化的
Closure Review。

### STEP 11: Run a fresh Closure Review

第一处代码修改会自动使初审 Manifest 和子 Agent 结论失效。完成已批准修复后，无需再次询问，
从 STEP 3 选择 Closure Review，重新执行安全验证，生成全新 Manifest 和 diff fingerprint，
并为适用 axis 启动全新的只读子 Agent。不得复用初审或争议裁决 Agent 的上下文、候选项或结论。

只有本轮 Closure Review 可以把已修复 finding 更新为 `resolved` 或 `reopened`。按照 STEP 6
裁决修复回归和依赖解锁问题，再按 STEP 7 刷新 Results、Axis Summary、Closure、风险表和
Merge readiness。IF 无法启动必需的新子 Agent，THEN 报告 `Blocked / No-Go`，不得由执行修复
的主 Agent 自行宣告闭合。

## Review dimensions

| Axis | Dimension | Focus |
|------|-----------|-------|
| Standards | Conformance | `AGENTS.md`、代码规范、框架和工具链约束 |
| Standards | Architecture & Reliability | 分层、依赖、错误处理、事务、并发、安全和资源生命周期 |
| Standards | Maintainability & Testing | 复杂度、重复、命名、测试质量、compatibility 和 code smells |
| Spec | Correctness & Coverage | requirement、AC、边界和错误行为是否完整正确实现 |
| Spec | Scope & Compatibility | 无关行为、scope creep、spec compatibility 和迁移影响 |
| Spec | Test Evidence | 测试是否能在 AC 缺失或错误时失败 |

## Verification

- [ ] 每轮评审的 Review Manifest 固定了 refs、三类 diff、untracked、来源、验证结果和
      SHA-256 fingerprint。
- [ ] 报告输出路径已作为 review artifact 单独记录，未污染实现 diff fingerprint。
- [ ] 已读取全部 diff、适用 repository guidance 和必要代码上下文。
- [ ] Spec-backed review 已解析 `design.md`、AC、Testing 和 `tasks.md`。
- [ ] Full / Closure 模式选择正确，Closure 新 finding 具有允许的 Origin。
- [ ] 必需 axis 由隔离的只读子 Agent 执行；两个 axis 适用时并行启动。
- [ ] 主 Agent 在聚合前重算 fingerprint，并独立核验全部候选项。
- [ ] 每条 finding 有 Axis、稳定 ID、位置、依据、证据、影响、建议和 Closure test。
- [ ] 同 axis 只按根因合并；跨 axis 只关联，不合并。
- [ ] Standards、Spec 和 Overall Result 与 findings 或 blocked 原因一致。
- [ ] 初审严格只读；报告完成后已询问是否处理 findings，并在得到回答前停止。
- [ ] 修复前已逐项确定处置并获得统一计划确认；修改仅覆盖批准的 `fix` 项。
- [ ] 修复后使用全新 Manifest、fingerprint 和只读子 Agent 执行 Closure Review。
- [ ] 未修改 `design.md`、`tasks.md` 或任务 checkbox。
- [ ] 报告声明了 `Report language`，所有叙述字段使用该语言。
- [ ] 文档语言调整没有改变 finding、证据、Severity、Status、Result 或 readiness。

## Safety & guardrails

- 证据优先于观点；无法定位和证明影响的内容不是 finding。
- 评审阶段严格只读；不得在 finding 尚未裁决或修复计划尚未统一确认时修改代码。
- 不运行未获单独授权的破坏性命令、数据写入、生产操作或需要真实凭证的验证。
- 不把生成代码、第三方代码或超出 review range 的旧问题混入当前 findings，除非改动使其可达。
- 不因实现者声明、checkbox、commit message 或测试通过而跳过独立核查。
- 写入 `code-review.md` 前必须重新读取实际 diff、项目文件、命令结果和历史报告并据此填写，
  不凭记忆补写或编造。
- 与用户交流及报告叙述使用 `Report language`，表达直接、自然；固定协议字段与技术内容保持
  原样。简单场景允许压缩报告，不堆砌模板内容。

## References

- [Code review template](assets/code-review-template.md) — Spec 场景的 `code-review.md` 输出骨架。
- [Code smell baseline](references/code-smells.md) — Standards Agent 仅在 Full Review 中读取。
