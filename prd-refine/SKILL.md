---
name: prd-refine
description: >
  优化和完善 PRD 文档的专用技能。当用户提到"根据评审报告优化 PRD"、
  "prd-review 之后怎么改"、"评审发现的问题怎么处理"、"PRD 需要补充完善"、
  "review 报告里的建议怎么落地"、"基于反馈更新需求文档"、"处理评审意见"时,
  必须使用此技能。将评审发现转化为逐项决策,记录方案对比和选择理由,在用户确认后
  才写回 PRD,确保需求演进过程可追溯。即使用户没有明确说"refine",只要涉及 
  PRD 优化、评审反馈处理、需求文档完善、评审建议落地,都应该使用这个 skill。
---

# PRD Refine

## 核心理念

Use this skill after `prd-review` has produced a review report. The goal is not to patch the PRD line by line, but to turn review findings into explicit decisions, evaluate options against the real project context, record the decision path, and only update the PRD after the user confirms the final writeback.

**为什么需要这个流程?** 直接修改 PRD 会丢失决策上下文。当未来需要调整需求时,团队无法回溯"当初为什么这样设计"。通过记录决策过程,我们保留了方案对比、取舍理由和项目约束,让需求演进有据可查。

**Match the user's language for all output.** 如果用户用中文提问,所有输出都使用中文。

## When to use

- 用户已有 PRD 和对应的 `<prd-basename>-review.md`，希望根据评审报告优化 PRD
- 需要从 PM 和架构师视角逐项分析需求合理性、可维护性、演进成本和最终效果
- 需要比较多个方案，并把决策权交给 PM/架构师确认
- 需要保留完整决策路径，最后再统一询问是否写回 PRD

## When not to use

- 只是评审 PRD 并找问题（用 `prd-review`）
- 已经有明确设计，需要拆任务（用 `spec-plan`）
- 需要从零规划大型产品方向或拆分多个交付目标（用 `brainstorming`）
- 需要从零编写单项需求的技术方案（用 `clarifying`）
- 用户只要求直接编辑一个小段落，且不涉及评审报告或决策过程

## Inputs

- PRD 文档路径
- 同目录下对应的 `<prd-basename>-review.md`
- 可选：PM/架构师补充说明、优先级、成本约束、上线时间或业务目标

## Outputs

- `<prd-basename>-refine.md`：写入 PRD 同级目录，记录决策队列、方案分析、用户确认、最终写回计划
- 可选 PRD 更新：只有在所有决策总结后，且用户明确确认写回 PRD 时，才修改原 PRD

## Workflow

### Step 1: Load context

Read the PRD, matching review report, and any existing `<prd-basename>-refine.md`.

If the refine file exists, resume from the first undecided item. Do not restart the process or discard previous decisions.

**Explore project context when a decision depends on current implementation.** 不要凭空假设项目现状,基于实际代码和文档做决策。

**何时需要调研项目上下文?**
- 决策涉及现有代码改动 → 搜索相关函数、类、API
- 需要了解数据模型 → 查看数据库 schema、ORM 模型定义
- 涉及第三方集成 → 检查现有集成代码、配置文件、API 文档
- 权限和安全相关 → 查看现有的认证授权实现
- 性能相关决策 → 了解当前数据量、并发情况、瓶颈点
- 技术栈选型 → 检查项目依赖(package.json、pom.xml、requirements.txt)

**如何调研?**
- 使用 Grep 搜索关键代码和配置
- 读取相关的设计文档(docs/、CLAUDE.md、README.md、ADR)
- 检查数据库 migration 文件或 schema 定义
- 查看配置文件(application.yml、.env、config.*)
- 阅读现有的单元测试(了解预期行为)

Use repository docs, code search, database exploration, and relevant references as needed. Do not invent project facts.

### Step 2: Build or update the decision queue

Convert review findings into decision items. **为什么需要决策队列?** 评审报告通常包含几十条发现,如果逐条处理会导致决策碎片化。将相关发现归并为决策项,能让 PM/架构师从产品或技术视角整体思考。

Group related findings by product or architecture decision instead of mirroring every review bullet mechanically. **例如:** 如果评审报告有 3 条关于权限的发现,不要创建 3 个决策项,而是归并为 1 个"权限模型设计"决策。

Prioritize:

1. `[P0-阻断]` decisions that block requirement closure (例如:核心功能需求不明确、技术可行性未验证)
2. `[P1-需确认]` decisions that affect product behavior or architecture direction (例如:交互流程设计、数据模型设计)
3. `[P2-优化]` decisions that improve maintainability, usability, or delivery quality (例如:错误提示优化、代码结构调整)

Each decision item should include:

- ID, title, status
- Source review finding and PRD reference (让决策可追溯到原始问题)
- Why this is a decision rather than a simple text edit (如果是简单文字修正,直接改 PRD 即可,不需要决策流程)
- Affected PRD sections

### Step 3: Analyze one decision item at a time

**逐项分析决策,而不是一次处理所有问题。** 为什么?因为每个决策都可能涉及多个方案对比、成本收益权衡和项目上下文调研,一次处理多项会让对话失焦,用户难以做出高质量决策。逐项处理能确保每个决策都得到充分讨论和确认。

Work on exactly one undecided item per interaction. For that item, provide:

- **Current issue and source** (问题是什么,来自评审报告哪一条)
- **Project-context evidence** (为什么需要证据?避免凭空假设项目现状,决策必须基于实际代码和文档)
- **2-3 viable options** (为什么是 2-3 个?太少没得选,太多增加决策负担。排除明显不可行的方案,只列出真正需要权衡的选项)
- **Pros, cons, and risks for each option** (帮助用户理解每个选项的取舍)
- **Impact analysis:**
  - User experience impact (用户会感知到什么变化?)
  - Delivery cost (开发工作量、时间成本)
  - Maintainability (未来维护和扩展的难度)
  - Extensibility (是否为未来需求预留空间?)
  - Final product effect (对最终产品质量的影响)
- **Recommended option with reasoning** (你的建议是什么,为什么?)
- **Focused confirmation question for PM/architect** (让用户做决策,不要替用户决定)

Ask the user to choose or revise the option. Do not modify the PRD at this stage.

**决策分析示例:**

<details>
<summary>示例:用户权限模型设计决策</summary>

**D1: 用户权限模型 - 是否支持动态角色?**

**问题来源:** Review 报告 P1-需确认 #3: "权限需求描述模糊,未明确角色定义和权限边界"

**项目现状调研:**
- 当前系统使用静态角色(admin/user/guest)
- 用户表已有 `role` 字段,类型为 `ENUM('admin','user','guest')`
- 权限判断硬编码在各个 Controller 的 `@PreAuthorize` 注解中
- 业务方在需求讨论会上提到"未来可能有区域经理、大区经理等多层级角色"

**方案分析:**

**方案 A: 保持静态角色 + 明确权限清单**
- 优点: 实现简单,2 人日可完成,代码易理解
- 缺点: 后续新增角色需要改表结构和代码
- 风险: 如果 3 个月内需要新增角色,会产生二次开发成本
- 用户体验: 无影响(用户无感知)
- 可维护性: 权限逻辑分散在代码中,修改需要改多处

**方案 B: 动态角色 + RBAC 权限管理**
- 优点: 扩展性好,运营可通过管理后台自助配置角色和权限
- 缺点: 开发成本增加到 8 人日,需要额外的角色管理、权限配置界面
- 风险: 当前版本不需要这么复杂的权限体系,可能过度设计
- 用户体验: 管理员可以自定义角色(如果需要的话)
- 可维护性: 权限配置化,后续调整无需改代码

**方案 C: 静态角色 + 预留扩展字段(JSON extra_permissions)**
- 优点: 当前简单实现,未来有一定扩展余地(3 人日)
- 缺点: JSON 字段不好查询,权限判断逻辑会变复杂
- 风险: 如果未来真需要复杂权限,还是要重构为方案 B
- 用户体验: 无影响
- 可维护性: 中等,扩展字段增加了理解成本

**推荐方案:** 方案 A

**推荐理由:**
1. 当前明确需求只有 3 个角色(admin/user/guest),PRD 中未提及角色配置需求
2. "未来可能"不是明确需求,过早优化会拖慢 V1 交付
3. 如果 V2 真需要动态角色,到时再重构成本可控(数据迁移简单)
4. 建议在 PRD 的"未来演进"章节补充:"V2 可能需要支持动态角色配置"

**需要你确认:**
- 你是否同意方案 A?
- 当前版本是否有"运营自主配置角色"的需求?(如果有,需要选方案 B)
- 业务方说的"区域经理、大区经理"是近期(3 个月内)要上线的功能吗?

</details>

### Step 4: Record the decision path

After the user confirms an option, update `<prd-basename>-refine.md` with:

- Options considered (为什么要记录被拒绝的方案?未来可能需要重新评估,保留当时的分析能节省重复调研成本)
- Recommendation
- Final decision (用户最终选择了什么,是否与推荐一致)
- Decision rationale (用户为什么选择这个方案?有哪些关键考量因素?)
- Confirmed constraints or assumptions (决策依赖的前提条件,例如"假设 3 个月内不新增角色")
- Expected PRD impact (这个决策会导致 PRD 哪些章节需要修改)

**为什么要详细记录?** 6 个月后当需求变化时,团队需要回溯"当初为什么这样设计",有了完整记录,就能快速理解决策上下文,避免重复讨论。

Then move to the next undecided item.

### Step 5: Final summary and writeback preview

After all decision items are confirmed, summarize:

- All decisions made (列出所有决策项的标题和最终选择)
- Major tradeoffs accepted (重要的取舍,例如"为了快速上线,暂不支持动态角色")
- PRD sections that will change (哪些章节需要修改)
- Requirement conflicts or unresolved risks, if any (是否有决策之间的冲突?是否有风险需要在 PRD 中标注?)
- Proposed PRD update preview (展示具体会怎么改 PRD,让用户预览)

**为什么需要预览?** 避免用户批准决策后,发现实际写入 PRD 的内容不符合预期。预览阶段是最后的检查点。

Ask the user whether to write the summarized decisions back to the PRD. **The question must be explicit and must not assume approval.** 

**正确的询问方式:**
- "我准备将上述决策写入 PRD,你确认要写回吗?"
- "上面是所有修改预览,我现在可以写回 PRD 吗?"

**错误的询问方式(不要这样问):**
- "我现在开始写回 PRD..." (没有给用户拒绝的机会)
- "应该可以写了吧?" (语气不够明确)

### Step 6: Update PRD only after explicit approval

If the user confirms writeback, update the PRD with final requirement text only. **Keep analysis, rejected options, and detailed tradeoffs in `<prd-basename>-refine.md`.** 

**为什么 PRD 和 refine 文件要分离?**
- PRD 是面向开发团队的需求文档,需要简洁清晰
- refine 文件是决策记录,面向未来的需求变更和复盘
- 混在一起会让 PRD 变得冗长,开发者只想知道"做什么",不需要看"为什么不选方案 B"

After writing, update the refine file with:

- Writeback time (记录写回时间,便于追溯)
- Changed PRD sections (具体修改了哪些章节)
- Summary of edits (修改摘要,例如"新增了权限矩阵表格,明确了 3 个角色的权限边界")
- Final status (标记为"已完成"或"已写回")

If the user declines writeback, leave the PRD unchanged and keep the refine file as the decision record. **这也是合理的结果** - 用户可能想等所有决策做完后再一起写回,或者发现有些决策需要重新讨论。

## Refine file structure

Use this structure for `<prd-basename>-refine.md`:

```markdown
# PRD 优化决策记录

> **PRD 文档：** ...
> **Review 报告：** ...
> **优化状态：** 进行中 / 待写回 / 已写回 / 已暂停

## 决策队列

- [x] D1：...
- [ ] D2：...

## D1：决策标题

### 问题来源
### 项目现状
### 方案分析
### 推荐方案
### 用户确认
### 对 PRD 的影响

## 最终决策总结

## PRD 更新预览

## 写回记录
```

## Guardrails

这些原则确保 skill 按照预期方式工作,避免常见错误:

- **Do not update the PRD before all relevant decisions are recorded and the user explicitly approves writeback.** 为什么?过早写回会导致决策不完整,用户事后发现问题时 PRD 已经被改了。
- **Do not treat review findings as a checklist of direct text additions.** 评审报告说"缺少错误码定义"不等于直接复制粘贴一段错误码列表到 PRD。要思考:为什么需要错误码?当前有没有?用什么格式定义?这才是决策。
- **Do not replace PM or architect judgment.** 你的角色是提供分析和建议,最终决策权在用户手里。即使你认为方案 A 明显更好,用户选择方案 B 也要尊重(但可以确认一次"你是否考虑了 XX 风险?")。
- **Do not invent future features.** "未来可能需要"不能作为当前版本的设计依据,除非用户明确说"3 个月内会上这个功能"。Options must be grounded in PRD goals, review findings, and current project context.
- **Keep PRD text clean and final. Keep decision history in the refine file.** PRD 里不要出现"我们考虑过方案 A/B/C,最后选了 A"这种内容,这些放在 refine 文件里。
- **If critical business priority, cost tolerance, or acceptance criteria exist only in the user's head, ask before deciding.** 不要假设"用户肯定希望快速上线",问清楚"这个版本的上线时间有硬性要求吗?可以接受多少人日的开发成本?"

## Verification

Before finishing, verify:

- [ ] `<prd-basename>-refine.md` exists next to the PRD (决策记录必须保存,不能只存在对话中)
- [ ] Every decision cites a review finding, PRD section, or project-context source (决策必须有依据,不能凭空想象)
- [ ] Each confirmed decision records considered options, tradeoffs, final choice, and rationale (完整的决策记录包含这 4 个要素)
- [ ] The PRD remains unchanged until explicit final writeback approval (没有用户明确同意,绝对不能改 PRD)
- [ ] The final summary lists all accepted decisions and affected PRD sections (让用户清楚知道会改什么)
- [ ] If PRD was updated, the refine file contains a writeback record (追溯性要求:能看出什么时候写回的、改了哪些地方)
