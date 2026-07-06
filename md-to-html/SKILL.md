---
name: md-to-html
description: 将 Markdown 转换为独立 HTML 文件（含目录导航、内容摘要、关键信息高亮和文档质量预检）。触发词：转 HTML / 导出 HTML / 生成网页 / convert md to html / render markdown。
---

# Markdown → HTML 增强转换

将 Markdown 文件转换为带智能分析增强的 HTML 文件，保存在源文件同目录，文件名相同仅扩展名改为 `.html`。

## 执行方式

脚本位于本 skill 目录下的 `scripts/convert.mjs`。用相对本 SKILL.md 的路径调用：

```bash
node scripts/convert.mjs "<markdown-file-path>"
```

示例：

```bash
node scripts/convert.mjs "docs/prds/example/02-prd.md"
```

### 前置检查

首次运行前需在 skill 目录内安装依赖：确认 `node_modules/marked` 是否存在，不存在则执行 `npm install`。具体命令按当前平台选择（POSIX：`[ -d node_modules/marked ] || npm install`；PowerShell：`if (-not (Test-Path node_modules/marked)) { npm install }`）。

若依赖缺失，脚本会因 `import { Marked } from "marked"` 抛错退出。

## 功能概要

- **摘要面板**：从 YAML frontmatter 提取元数据；统计功能模块（F0x）、业务规则（BR-xx）、验收标准（AC-xx）、表格、代码块、字数；列出 h2 章节标题与首段摘要。
- **目录导航**：左侧固定侧边栏（270px），列出 h1–h3，点击跳转、滚动高亮当前章节。窄屏（< 1080px）自动隐藏。
- **关键信息高亮**：自动识别 `BR-xx`（红）、`AC-xx`（绿）、`F0x`（蓝）、`Given/When/Then`（紫）。
- **质量预检**：转换时执行 8 类检查（未闭合代码块、标题层级跳跃、空链接、表格列数不一致、编号重复、frontmatter 缺失、中英文间距、连续空行），结果仅打印到控制台，不嵌入 HTML。
- **其他**：内联 CSS、返回顶部按钮、平滑滚动、`@media print` 打印优化。

普通 Markdown 同样可转换，统计项显示为 0，不影响输出。

## 控制台输出示例

```
✅ 转换成功: /path/to/02-prd.html
   📊 功能模块: 13 | 业务规则: 42 | 验收标准: 23 | 表格: 33
   ✅ 文档质量检验: 全部通过
```

## 批量转换

对目录下所有 `.md` 递归转换——按平台选择命令：

- **POSIX（bash / zsh）**：`find docs -name "*.md" -print0 | xargs -0 -n1 node scripts/convert.mjs`
- **PowerShell**：`Get-ChildItem docs -Recurse -Filter *.md | ForEach-Object { node scripts/convert.mjs $_.FullName }`

## 错误处理

| 场景 | 脚本行为 |
|------|----------|
| 未提供文件路径 | 输出用法提示，退出码 1 |
| 文件不存在 | 输出错误路径，退出码 1 |
| 非 `.md` / `.markdown` 文件 | 输出格式警告，退出码 1 |

## 注意事项

- **覆盖策略**：目标目录已有同名 HTML 会被静默覆盖。批量转换前请与用户确认，或先备份。
- **依赖**：仅需 `marked`（已在 `package.json` 声明），Node 版本需支持 ES modules（Node ≥ 18）。
- **分析范围**：BR/AC/F 高亮基于正则模式识别，为 PRD 类文档优化；对普通文档无副作用。
