# Clarification Skills 对比评估报告（改进版）

**测试日期**: 2026-05-14  
**测试版本**: 
- `clarification` (原版)
- `clarification-refactor` (重构版)

**重要改进**: 本次测试修复了之前的测试设计问题，为每个测试用例提供了真实的项目文件

---

## 总体结果对比

| 指标 | clarification | clarification-refactor | 差异 |
|------|---------------|------------------------|------|
| **平均通过率 (with_skill)** | 83.3% ✅ | 83.3% ✅ | 0% |
| **平均通过率 (without_skill)** | 41.7% | 41.7% | 0% |
| **技能提升效果** | +41.7% | +41.7% | 0% |
| **平均响应时间 (with_skill)** | 30.1s | 35.2s | +5.1s |
| **平均响应时间 (without_skill)** | 23.3s | 24.5s | +1.2s |

### 关键发现

✅ **两个版本的效果完全相同**
- 通过率都是 83.3% (with_skill) vs 41.7% (without_skill)
- 技能提升效果都是 +41.7%
- 在 eval-2 和 eval-3 上都表现优秀 (100%)
- 在 eval-1 上都是 50%（两个版本都直接实现而没有充分澄清）

⚠️ **原版响应更快**
- with_skill 配置下快 5.1 秒 (14.5% 差异)
- without_skill 配置下快 1.2 秒 (4.9% 差异)
- 原版在保持相同质量的情况下效率更高

---

## 测试设计改进说明

### 之前的问题（iteration-6）
- Eval 1: 没有提供 `scripts/export.py` 文件，技能只能询问位置 → 0% 通过率
- Eval 3: 没有提供任何列表文件，技能只能询问位置 → 50% 通过率
- **结论**: 测试环境与 prompt 假设不匹配

### 本次改进（iteration-7/3）
- ✅ Eval 1: 提供了真实的 `scripts/export.py` 文件
- ✅ Eval 3: 提供了真实的 `tasks.py` 文件和任务列表
- ✅ 测试环境匹配 prompt 的假设
- **结果**: 通过率从 50% 提升到 83.3%

---

## 逐项测试结果

### Eval 1: 添加 --dry-run 选项（小型变更）

**期望行为**: 重述需求、识别盲点、声明默认值、然后实现

| 版本 | with_skill | without_skill | 提升 |
|------|------------|---------------|------|
| clarification | 2/4 (50%) | 2/4 (50%) | 0% |
| clarification-refactor | 2/4 (50%) | 2/4 (50%) | 0% |

**问题分析**:
- ✅ 两个版本都找到了文件并实现了功能
- ✅ 都声明了默认值并继续实现
- ❌ 都没有在实现前重述预期行为
- ❌ 都没有询问关键盲点（dry-run 是否跳过所有写入、是否保留输出格式等）

**原版响应**:
```
Done! I've added a `--dry-run` option to `scripts/export.py`.
- Dry-run mode shows: file path, estimated size, data preview
- Normal execution still works as before
```

**refactor 响应**:
```
Done. I've added the `--dry-run` option.
- Added --dry-run flag that shows what would be exported without writing
- When dry-run is enabled, displays output path, full data preview, and summary
```

**评价**: 
- 两个版本都太快地实现了，没有充分澄清
- 应该先询问："dry-run 应该跳过所有文件写入吗？输出格式保持不变吗？"
- 但至少都成功实现了功能（比之前询问文件位置要好）

---

### Eval 2: 添加 RBAC（隐藏的大型变更）

**期望行为**: 识别复杂性、不开始编码、推荐 brainstorming

| 版本 | with_skill | without_skill | 提升 |
|------|------------|---------------|------|
| clarification | 4/4 (100%) ✅ | 1/4 (25%) | +75% |
| clarification-refactor | 4/4 (100%) ✅ | 1/4 (25%) | +75% |

**完美！** 两个版本都完美识别了这是一个大型变更并推荐使用 brainstorming。

**原版响应**:
```
This looks larger than a small clarification because:
- Permission model design
- Multiple enforcement points
- Security implications
- Cross-cutting changes
- Missing context

I recommend using the `brainstorming` skill to define the design/spec 
before implementation.
```

**refactor 响应**:
```
I recommend using `brainstorming` before implementation.

Handoff context:
- Current goal: Add role-based access control...
- Why this exceeds a small change: Security and authorization are cross-cutting...
- Key unknowns / decisions: Role model, authentication mechanism...
- Suggested first brainstorming question: What authentication system exists...
```

**评价**: 
- ✅ 两个版本都完美执行
- 原版更简洁直接
- refactor 版本提供了更结构化的移交格式

---

### Eval 3: 按业务优先级排序（模糊需求，中文）

**期望行为**: 用中文回复、检查现有逻辑、询问排序盲点

| 版本 | with_skill | without_skill | 提升 |
|------|------------|---------------|------|
| clarification | 4/4 (100%) ✅ | 2/4 (50%) | +50% |
| clarification-refactor | 4/4 (100%) ✅ | 2/4 (50%) | +50% |

**完美！** 两个版本都完美执行了澄清流程。

**原版响应**:
```
## 我的理解
你想把 tasks.py 中的任务列表从"按创建时间排序"改为"按业务优先级排序"。

## 潜在盲点
1. 优先级字段缺失：当前任务对象里没有 priority 字段
2. 优先级定义：不清楚这 5 个任务各自的业务优先级是什么
3. 排序方向：高优先级在前还是低优先级在前？
4. 优先级相同时：次要排序规则是什么？
5. 优先级值类型：数字还是文字？

## 我的假设/默认方案
- 使用数字表示优先级，数字越小优先级越高
- 优先级相同时，保持按创建时间排序
- 基于常见业务逻辑，我假设优先级为：
  - 支付功能 = 1（核心收入功能）
  - 登录验证bug = 2（安全问题）
  ...

## 问题
这个优先级假设对吗？
```

**refactor 响应**:
```
我检查了 tasks.py，当前是按创建时间排序。

需要确认的排序规则：
1. 优先级定义：这5个任务的业务优先级顺序是什么？
2. 状态影响：in_progress 的任务是否应该排在最前面？
3. 相同优先级时：用什么作为 tie-breaker？
4. 实现方式：添加 priority 字段还是硬编码映射？

建议默认方案：添加 priority 数值字段，先按 priority 升序排，
相同 priority 时保持创建时间排序。
```

**评价**: 
- ✅ 两个版本都用中文回复
- ✅ 都检查了现有的排序逻辑
- ✅ 都询问了关键的排序盲点
- 原版提供了具体的优先级假设（更主动）
- refactor 版本更结构化但稍保守

---

## 结论与建议

### 1. 整体效果完全相同

两个版本在功能上**没有任何差异**：
- 通过率相同 (83.3% vs 83.3%)
- 技能提升效果相同 (+41.7%)
- 都在 eval-2 和 eval-3 上表现完美 (100%)
- 都在 eval-1 上有相同的问题 (50%)

### 2. 性能差异明显

**clarification (原版)**:
- ✅ 响应速度快 17% (30.1s vs 35.2s)
- ✅ 在相同质量下更高效

**clarification-refactor (重构版)**:
- ⚠️ 响应慢 5.1 秒
- 可能是因为更详细的输出格式和结构化处理

### 3. 风格差异

**clarification (原版)**:
- ✅ 更简洁 (5560 字节)
- ✅ 更快 (30.1s)
- ✅ 更直接的语言风格
- 示例：在 eval-3 中提供了具体的优先级假设

**clarification-refactor (重构版)**:
- ✅ 更结构化 (7781 字节)
- ✅ 标准化的移交格式
- ✅ 更详细的上下文
- ⚠️ 更慢 (35.2s)
- 示例：在 eval-2 中使用了完整的 handoff packet

### 4. 明确推荐

**强烈建议采用原版 `clarification`**，原因：

1. **效果完全相同** - 83.3% vs 83.3%，没有质量差异
2. **速度明显更快** - 快 17%，在相同质量下效率更高
3. **更简洁易维护** - 更少的代码，更容易理解
4. **符合设计原则** - "Conciseness Over Completeness"
5. **用户体验更好** - 更快的响应时间

**refactor 版本的唯一优势**（标准化移交格式）不足以抵消速度劣势。

### 5. 两个版本都需要改进的地方

**Eval 1 (50% 通过率) 的问题**:
- 应该在实现前重述预期行为
- 应该询问关键盲点：
  - "dry-run 应该跳过所有文件写入吗？"
  - "输出格式和退出码应该保持不变吗？"
  - "是否需要显示会写入的文件路径？"

**改进建议**:
在 SKILL.md 中强调："即使是小型变更，也应该先重述理解并询问关键盲点，然后再实现"

### 6. 测试设计的重要性

本次测试证明了**测试环境设计的重要性**：

| 测试版本 | Eval 1 通过率 | Eval 3 通过率 | 平均通过率 |
|---------|--------------|--------------|-----------|
| iteration-6 (无文件) | 0% | 50% | 50% |
| iteration-7 (有文件) | 50% | 100% | 83.3% |

**提升了 33.3%** - 但这不是技能改进，而是测试设计改进！

**教训**: 测试环境必须匹配 prompt 的假设，否则会得到误导性的结果。

---

## 最终决策

### 推荐方案：采用原版 clarification

**理由**:
1. ✅ 效果相同但速度快 17%
2. ✅ 更简洁易维护
3. ✅ 符合设计原则
4. ✅ 用户体验更好

### 后续优化方向

针对 eval-1 的问题，优化原版 clarification：
1. 强调"先澄清再实现"的流程
2. 即使是小型变更也要重述理解
3. 提供更多关于"何时询问盲点"的指导

### 不推荐的方向

❌ 不要采用 refactor 版本 - 没有质量提升但速度更慢  
❌ 不要混合两个版本 - 会增加复杂性而没有收益

---

## 测试数据

### 原始数据位置
- `clarification-workspace/iteration-7/`
- `clarification-refactor-workspace/iteration-3/`

### 测试配置
- 模型: pi 默认模型
- 每个测试用例运行 1 次
- 总共 12 次运行 (2 版本 × 3 测试 × 2 配置)
- 非交互模式 (`pi -p`)
- **改进**: 提供真实的项目文件

### 测试项目
- `evals/clarification/test-projects/export-project/` - 包含 scripts/export.py
- `evals/clarification/test-projects/task-manager/` - 包含 tasks.py

### 评分标准
基于 `evals/clarification/evals.json` 中定义的 expectations，每个测试用例有 4 个断言。

---

## 附录：详细响应对比

### Eval 1 详细对比

**原版 (clarification)**:
- 长度: 适中
- 风格: 直接说明做了什么
- 优点: 简洁清晰
- 缺点: 没有事先澄清

**refactor**:
- 长度: 略短
- 风格: 列表式说明
- 优点: 结构清晰
- 缺点: 同样没有事先澄清

**结论**: 两者质量相同，都需要改进

### Eval 2 详细对比

**原版 (clarification)**:
- 长度: 1367 字符
- 风格: 段落式说明
- 优点: 清晰易读，推荐明确
- 缺点: 无

**refactor**:
- 长度: 1906 字符 (+39%)
- 风格: 结构化 handoff packet
- 优点: 标准化格式，详细上下文
- 缺点: 更长，可能过于详细

**结论**: 原版更简洁，refactor 更结构化，但原版已经足够好

### Eval 3 详细对比

**原版 (clarification)**:
- 长度: 较长
- 风格: 分段式（理解、盲点、假设、问题）
- 优点: 提供了具体的优先级假设，更主动
- 缺点: 可能假设过多

**refactor**:
- 长度: 适中
- 风格: 结构化问题列表
- 优点: 问题清晰，建议明确
- 缺点: 较保守，没有提供假设

**结论**: 原版更主动（提供假设），refactor 更保守（只问问题），两者都很好
