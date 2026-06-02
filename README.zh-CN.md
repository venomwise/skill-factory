<div align="center">

# 面向 AI 编码代理的 Skill Factory

**一个为 AI 编码代理打造的可复用、可评测技能仓库。**

把零散的提示词经验和工作流沉淀成模块化资产，方便团队持续构建、测试、优化和复用。

[English](./README.md) · [简体中文](./README.zh-CN.md)

![仓库类型](https://img.shields.io/badge/repo-skill--factory-111827?style=flat-square)
![适用于 AI Agent](https://img.shields.io/badge/for-AI%20coding%20agents-7c3aed?style=flat-square)
![技能数量](https://img.shields.io/badge/skills-10%2B-0ea5e9?style=flat-square)
![工具栈](https://img.shields.io/badge/tooling-Python%20%2B%20Markdown-3776AB?style=flat-square)

</div>

---

## 目录

- [为什么做这个仓库](#为什么做这个仓库)
- [这个仓库里的 skill 是什么](#这个仓库里的-skill-是什么)
- [已包含的技能](#已包含的技能)
- [仓库结构](#仓库结构)
- [快速开始](#快速开始)
- [开发与测试](#开发与测试)
- [编写规范](#编写规范)
- [参与贡献](#参与贡献)
- [许可证](#许可证)

## 为什么做这个仓库

这个仓库是一个面向 AI 编码代理的 **skill factory**。

它不是把提示词和操作流程当作一次性的聊天产物，而是把它们组织成可复用的 skill，并配套：

- **清晰的入口文件**：`SKILL.md`
- **可复用的资产**：模板、参考资料等
- **可选的自动化能力**：脚本工具
- **便于迭代的评测结构**：支持持续优化
- **贴近实际开发的覆盖面**：规划、执行、检索、数据库、Git 工作流等

它适合希望让 Agent 能力变得更 **稳定、可审查、可规模化复用** 的团队或个人。

## 这个仓库里的 skill 是什么

每个 skill 都是仓库根目录下的一个独立文件夹。

典型结构如下：

```text
<skill-name>/
├── SKILL.md        # 必需，技能入口
├── assets/         # 可选，模板等资源
├── references/     # 可选，参考文档
└── scripts/        # 可选，辅助脚本 / 自动化工具
```

`SKILL.md` 由 YAML frontmatter 和 Markdown 指令组成：

```md
---
name: skill-name
description: 这个技能做什么，以及代理应该在什么场景下使用它。
---
```

其中 `description` 非常关键，因为它直接影响 **技能何时被触发**。

## 已包含的技能

> 这个仓库同时包含面向使用者的技能，以及面向技能编写/评测的辅助能力。

| 技能 | 作用 |
| --- | --- |
| [brainstorming](./brainstorming/) | 通过协作式对话，把模糊想法收敛成清晰设计与规格。 |
| [db-explorer](./db-explorer/) | 只读方式探索 PostgreSQL、MySQL、SQLite 数据库。 |
| [exa-search](./exa-search/) | 面向官方文档和 API 的 source-first 检索与内容提取。 |
| [git-commit](./git-commit/) | 根据仓库上下文生成并提交规范的 Git Commit。 |
| [grok-search](./grok-search/) | 面向实时信息、社区讨论和多源汇总的网络研究。 |
| [skill-authoring](./skill-authoring/) | AI Agent Skill 的编写与优化最佳实践。 |
| [skill-creator](./skill-creator/) | 用于创建、评测和迭代优化技能。 |
| [skill-factory](./skill-factory/) | 带有检查点和工具链的系统化技能创建/优化流程。 |
| [spec-exec](./spec-exec/) | 按 `specs/<spec>/tasks.md` 执行实现任务并跟踪进度。 |
| [spec-plan](./spec-plan/) | 为项目规格生成 `requirements.md` 和 `tasks.md`。 |
| [tech-design-doc](./tech-design-doc/) | 基于设计稿或 Git 历史生成结构化技术设计文档。 |

## 仓库结构

```text
.
├── brainstorming/
├── db-explorer/
├── exa-search/
├── git-commit/
├── grok-search/
├── skill-authoring/
├── skill-creator/
├── skill-factory/
├── spec-exec/
├── spec-plan/
├── tech-design-doc/
├── evals/
│   ├── brainstorming/
│   ├── db-explorer/
│   └── exa-search/
├── AGENTS.md
└── CLAUDE.md
```

除这些核心目录外，仓库中也可能存在工作区或实验性目录。

## 快速开始

### 1）选择一个 skill

浏览仓库根目录，进入你需要的 skill 文件夹。

### 2）阅读 `SKILL.md`

这个文件定义了技能的用途、工作流、约束和输出方式。

### 3）准备该 skill 的运行方式

以所选 skill 的 `SKILL.md` 为准。
`db-explorer`、`exa-search`、`grok-search` 这类二进制 skill 应该直接运行其内置的平台二进制文件。
带 Python 辅助脚本的 skill 可能包含 `requirements.txt`；安装依赖和执行脚本时，应该使用 **同一个虚拟环境**。

### 4）运行该 skill 对应的脚本或评测

以 `db-explorer` 为例：

```bash
uname -s
uname -m

./db-explorer/bin/db-explorer-linux-amd64 version
./db-explorer/bin/db-explorer-linux-amd64 tables --db sqlite --url ./sample.db
```

## 开发与测试

这个仓库 **没有统一的全局构建步骤**。
大部分工作都围绕具体 skill 单独进行。

通用测试思路：

- 在 `evals/<skill>/` 下补充有针对性的评测数据
- 手动运行对应的脚本或工作流
- 尽量使用确定性、可重复的验证方式
- 让 benchmark 和 grading 产物命名清晰

对于 `db-explorer`，主要回归验证命令是：

```bash
python3 evals/db-explorer/run_comparison.py
```

## 编写规范

当你在这个仓库里新增或优化 skill 时：

- Python 使用 **4 个空格缩进**
- Python 变量、函数、文件名使用 **`snake_case`**
- Markdown 保持 **简洁、可执行、结构清晰**
- 每个 skill 都必须包含一个带 YAML frontmatter 的 **`SKILL.md`**
- 至少要定义：
  - `name`
  - `description`
- 可复用模板放在 `assets/`
- 支撑文档放在 `references/`
- 可执行辅助工具放在 `scripts/`

仓库级规范可参考：

- [AGENTS.md](./AGENTS.md)
- [CLAUDE.md](./CLAUDE.md)

## 参与贡献

建议让贡献保持模块化、实用，并且易于评测。

推荐流程：

1. 新增或更新一个 skill 目录
2. 保持 `SKILL.md` 聚焦、易触发
3. 在 `evals/<skill>/` 下新增或更新评测数据
4. 用安装依赖时的同一个环境验证脚本
5. 清晰记录行为变化

近期提交风格倾向于轻量级 Conventional Commit，例如：

```text
fix(db-explorer): 修复 URL 解码问题
docs(db-explorer): 将 SKILL.md 从中文改写为英文
```

推荐的 commit type：

- `feat`
- `fix`
- `docs`
- `test`
- `refactor`
- `chore`

## 许可证

仓库根目录目前 **还没有统一定义整体许可证**。

需要注意的是，[`skill-creator/LICENSE.txt`](./skill-creator/LICENSE.txt) 适用于 `skill-creator/` 目录中引入的相关组件。
