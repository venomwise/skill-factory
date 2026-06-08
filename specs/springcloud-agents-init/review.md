# Spring Cloud Agents Init Skill 设计评审

评审对象:`specs/springcloud-agents-init/design.md`
评审依据:`skill-best-practice/skill-authoring-best-practices.md` 与本仓库现有 skill(db-explorer 等)约定
视角:Prompt Engineering + AI Agent Engineering

## 总体结论

设计扎实、覆盖完整。evidence-first 与 read-before-write 两条主线把握准确,渐进披露结构(`SKILL.md + references/detection.md + templates.md`)、非目标、风险表均符合规范。可进入实现,但以下几处会直接影响实际运行质量,建议先调整。

优化点按优先级排序:P0 强烈建议在 v1 落地,P1 建议落地,P2 可选增强。

## P0 优化点

### P0-1 加入 plan-validate-execute 检查点

当前工作流是"检测 → 直接写所有文档"。在 30+ 服务的 monorepo 里,这等于一次性生成几十个文件 —— 中等 blast radius 的动作,且检测可能误判。规范专门给出 plan-validate-execute 模式来在动手前拦截错误。

建议:在 Step 6 与 Step 7 之间插入确认检查点,先产出"计划清单"再批量写:

- 微服务清单 + 每个的判定证据
- 共享模块清单 + 判定证据
- `Needs verification` 清单
- 计划写入/更新的文件路径列表

用户确认后再批量写入。单点收益最大,同时缓解"误判"和"一次性生成几十个文件"两个风险。

### P0-2 模板从"必填骨架"改为"可裁剪默认值"

规范要求按任务脆弱性/可变性匹配自由度。微服务文档生成是高可变任务,但设计给的是大量固定表格 + 成功标准里"没证据就填 Unknown"。两者叠加会让 Agent 机械填满每个单元格的 `Unknown`,产出低价值文档,违背设计 11.3 的"最小有效文档"。

建议:

- 模板降级为"sensible default, 按需裁剪"(规范的 flexible guidance 写法)
- 规则改为:有证据才写,无证据的小节直接省略
- 所有不确定项收敛到唯一的 `Open Questions`,不散落成满屏 Unknown

### P0-3 复跑/更新加显式合并标记

设计反复强调"保留人工说明",但没说清机制:第二次运行时 Agent 如何区分"上次自己生成的"与"人手写的"?没有标记就只能靠语义猜,容易要么覆盖人工内容、要么重复追加。

建议:

- 为生成区块加显式标记,如 `<!-- managed:springcloud-agents-init -->` … `<!-- /managed -->`
- 更新时只重写标记区内内容,标记外一律保留
- 把该机制写进 Step 2 和 11.1

## P1 优化点

### P1-1 templates.md 补一个"填好的正例 + 一个反例"

规范明确:examples 比 description 更能让模型对齐质量与详略度。当前 templates.md 全是空骨架。

建议:

- 放一个填好的 mini 服务文档正例(含真实证据引用)
- 配一个"类清单式臃肿文档"反例
- 把 11.3 的"不要变成类清单"从抽象规则变成可对照样例

### P1-2 大仓采用 breadth-first 工作流

Step 流程是线性的,没回答"50 个服务怎么办"。建议显式改为广度优先:

- 先建全量模块清单(只读 pom)
- 再只对确认的微服务做深挖
- 避免对 shared module 做昂贵的全量扫描

与 P0-1 的检查点天然配合。

### P1-3 Step 9 验证清单做成可复制 checklist

规范推荐复杂工作流用可勾选清单。Step 9 目前是散文式描述,建议改成 Agent 可贴进回复逐项打勾的 checklist,降低跳步概率。

## P2 可选增强

### P2-1 v1 引入"只枚举不判定"的辅助脚本

设计 12 节选 v1 纯靠 Agent 手工 `find`/`rg`/`read`。理由(避免隐藏判断过程)成立,但在大仓里逐个读 N 份 `pom.xml` + N 份 `application.yml` + 全局扫主类,会让上下文急剧膨胀且结果不稳定。规范指出 utility script 更可靠、省 token、保证一致性。

建议折中:

- 写一个只做枚举、不做判定的脚本,遍历 `pom.xml` 输出结构化 JSON(path / artifactId / packaging / 依赖列表)
- "是不是微服务"的判定仍留给 Agent,保留判断透明度
- 把最耗 token、最易出错的机械扫描下沉到脚本
- 判定逻辑不稳再迭代,符合设计原有的迭代思路

定为 P2:若不加,P1-2 的 breadth-first 也能缓解上下文膨胀。

## name / description 评价

- description 三人称、含关键词与中文触发词("总分式"),符合规范,under-triggering 风险低
- `springcloud-agents-init` 非规范偏好的动名词形式,但 `init` 属 action-oriented 的"可接受替代",无需改

## 对 15 节 7 个待确认问题的建议

| # | 问题 | 建议 |
|---|---|---|
| 1 | 名称 `springcloud-agents-init` | 可用 |
| 2 | 文件选择策略 | 合理,接受 |
| 3 | v1 只支持 Maven | 接受,Gradle 列增强 |
| 4 | Git 规则只提取已有、不编模板 | 接受,符合 git-safety 与防编造 |
| 5 | 服务文档放微服务 pom.xml 同级 | 接受,定位明确 |
| 6 | 关键 shared module 是否单独建文档 | 不单独建,只在根文档表格列出 + 标注影响范围,避免文件爆炸 |
| 7 | 是否加脚本 | 见 P2-1:v1 可加枚举脚本,不加判定脚本 |

## 落地优先级小结

- P0(强烈建议 v1):检查点、模板自由度、合并标记
- P1(建议 v1):正反例、breadth-first、checklist
- P2(可选):枚举脚本
