# Skill Document Language Style

本技能系列的语言风格规范，适用于所有技能文档（SKILL.md、`references/`、`assets/`）。

Skill 文档采用精细化中英文混合：正文用自然中文，英文仅作为结构或语义锚点保留。禁止在句中夹杂英文动词或分句。

## 为什么

- 风格不一致会稀释指令感——同一文件内两种文体混用，每条规则的约束力都打折
- 文档腔调会泄漏到用户可见输出：agent 会模仿技能文档的说话方式
- 翻译腔损耗语义精度，且每次编辑都要重新决定"新内容用什么文体写"

## 执行时的对话语气

技能执行中与用户的交流（提问、确认、汇报进展）要口语化、有温度——像和一位靠谱的同事商量事情，不是念客服话术。技能文档应在 workflow 或 guardrails 里写明语气要求，否则 agent 默认输出公文体。

```
✅ 好：这里我拿不准，有两种理解——你是想 A，还是 B？我倾向 A，
      因为你前面提到的场景更像它。

❌ 坏：请确认当前需求的语义解释为 A 抑或 B。
```

要点：

- 用口语表达，避免书面腔与公文句式堆叠
- 有立场、有温度：给出自己的倾向和理由，不要冷冰冰地把问题抛回去
- 提问的节奏与形式（一次一问、带选项）见 SKILL.md 的"意图优先"规范

## 英文保留范围

- 控制标记：`STEP`、`IF`、`THEN`、`STOP`、`PAUSE`、`Verify`、`Attention`
- 技术术语与模板章节名：architecture、data flow、Decision Record、Open Questions
- 路径 / 字段名 / 枚举标签：`specs/<topic>/design.md`、kebab-case、`Problem unclear`
- 技能名或二进制名：`/design-review`、`db-explorer`

## 正例

英文作锚点，中文作正文：

```
STOP when 你能用 2-3 句话描述项目的目的、技术栈和结构。
IF 需求跨越多个独立子系统，THEN 先与用户确认拆解方案。
将设计文档写入 `specs/<topic>/design.md`。
设计涵盖 architecture、components、data flow、error handling 和 testing。
```

## 反例

句中代码切换：

```
If 项目复杂或不熟悉，consult the guide for exploration priorities。
prefer 多选题。
If not，ask more questions。
```
