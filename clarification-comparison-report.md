# Clarification Skills 对比评估报告

**测试日期**: 2026-05-14  
**测试版本**: 
- `clarification` (原版)
- `clarification-refactor` (重构版)

**测试用例**: 3个场景，每个场景测试 with_skill 和 without_skill 两种配置

---

## 总体结果对比

| 指标 | clarification | clarification-refactor | 差异 |
|------|---------------|------------------------|------|
| **平均通过率 (with_skill)** | 50.0% | 50.0% | 0% |
| **平均通过率 (without_skill)** | 25.0% | 25.0% | 0% |
| **技能提升效果** | +25.0% | +25.0% | 0% |
| **平均响应时间 (with_skill)** | 18.6s | 19.4s | +0.8s |
| **平均响应时间 (without_skill)** | 17.0s | 16.4s | -0.6s |

### 关键发现

✅ **两个版本的整体效果相同**
- 都将通过率从 25% 提升到 50%
- 都在 eval-2 (RBAC) 上表现优秀 (100% vs 25%)
- 都在 eval-1 和 eval-3 上有改进空间

⚠️ **refactor 版本略慢**
- with_skill 配置下慢 0.8 秒 (4.3% 差异)
- 可能是因为更详细的输出格式

---

## 逐项测试结果

### Eval 1: 添加 --dry-run 选项（小型变更）

**期望行为**: 重述需求、识别盲点、声明默认值、然后实现

| 版本 | with_skill | without_skill | 提升 |
|------|------------|---------------|------|
| clarification | 0/4 (0%) | 0/4 (0%) | 0% |
| clarification-refactor | 0/4 (0%) | 0/4 (0%) | 0% |

**问题分析**:
- ❌ 两个版本都询问文件位置而不是假设合理的默认值
- ❌ 都没有重述 dry-run 的预期行为
- ❌ 都没有继续实现

**原版响应**:
```
I don't see a `scripts/export.py` at the root level of this repository...
Could you clarify:
1. Which `export.py` file you're referring to?
2. What should this script export?
```

**refactor 响应**:
```
I found multiple `export.py` files in evaluation workspaces...
Which one would you like me to modify? Or would you like me to:
1. Create a new `scripts/export.py` at the root level
2. Modify one of the existing files...
```

**改进建议**: 
- 应该假设用户指的是项目中的 `scripts/export.py`（即使不存在也可以询问是否创建）
- 应该先重述 dry-run 的预期行为（跳过写入、保留输出格式等）

---

### Eval 2: 添加 RBAC（隐藏的大型变更）

**期望行为**: 识别复杂性、不开始编码、推荐 brainstorming

| 版本 | with_skill | without_skill | 提升 |
|------|------------|---------------|------|
| clarification | 4/4 (100%) ✅ | 1/4 (25%) | +75% |
| clarification-refactor | 4/4 (100%) ✅ | 1/4 (25%) | +75% |

**成功！** 两个版本都完美识别了这是一个大型变更并推荐使用 brainstorming。

**原版响应**:
```
**My understanding**: Add role-based access control...

**Why this is larger than a small clarification**:
This is a **permission model**, not a localized tweak. It involves:
- Security architecture
- Multiple subsystems
- Authorization boundaries
- Migration strategy
...

**Recommendation**: Use the `brainstorming` skill to design the RBAC system properly.
```

**refactor 响应**:
```
This is larger than a small clarification because it involves:
- Authorization and security
- Multiple subsystems
- Missing design decisions
...

I recommend using `brainstorming` before implementation.

Handoff context:
- Current goal: ...
- Why this exceeds a small change: ...
- Key unknowns / decisions: ...
- Suggested first brainstorming question: ...
```

**对比**:
- ✅ 原版：更简洁，清晰推荐 brainstorming
- ✅ refactor：使用标准化的移交格式，提供更详细的上下文
- 两者都符合要求，refactor 版本更结构化

---

### Eval 3: 按业务优先级排序（模糊需求，中文）

**期望行为**: 用中文回复、不立即实现、计划检查现有逻辑、询问排序盲点

| 版本 | with_skill | without_skill | 提升 |
|------|------------|---------------|------|
| clarification | 2/4 (50%) | 2/4 (50%) | 0% |
| clarification-refactor | 2/4 (50%) | 2/4 (50%) | 0% |

**部分成功**: 
- ✅ 都用中文回复
- ✅ 都没有立即实现
- ❌ 都没有计划检查现有排序逻辑
- ❌ 都没有询问排序特定的盲点（tie-breaker、未知优先级等）

**原版响应**:
```
我需要先看看你说的是哪个列表。能告诉我具体是哪个文件或者哪部分内容吗？
或者你可以直接把列表内容发给我。
```

**refactor 响应**:
```
能否告诉我具体是哪个文件中的列表？我看到项目中有多个文档：
- `README.md` / `README.zh-CN.md` - 项目主文档
- `AGENTS.md` / `CLAUDE.md` - 代理指南
- 各个技能目录下的文档

或者你是在看某个特定文件？我可以帮你找到并按业务优先级重新排序。
```

**对比**:
- refactor 版本提供了更多上下文（列出了可能的文件）
- 但两者都缺少对排序盲点的询问

**改进建议**:
- 应该询问："业务优先级"的定义是什么？
- 如果优先级相同，如何排序？
- 是否保留现有的次要排序规则？

---

## 结论与建议

### 1. 整体效果相同
两个版本在功能上**没有显著差异**：
- 通过率相同 (50% vs 50%)
- 技能提升效果相同 (+25%)
- 都在 RBAC 测试上表现优秀
- 都在其他两个测试上有改进空间

### 2. 风格差异

**clarification (原版)**:
- ✅ 更简洁 (5560 字节)
- ✅ 响应速度略快 (18.6s vs 19.4s)
- ✅ 易于理解和维护

**clarification-refactor (重构版)**:
- ✅ 更结构化 (7781 字节)
- ✅ 使用标准化的移交格式
- ✅ 提供更详细的上下文
- ⚠️ 略慢 (+0.8s)
- ⚠️ 更复杂

### 3. 推荐方案

**建议采用原版 `clarification`**，原因：
1. **效果相同但更简洁** - 在相同效果下，简洁性是优势
2. **响应更快** - 虽然差异不大，但原版略快
3. **易于维护** - 更少的代码意味着更容易理解和修改
4. **符合 skill-factory 原则** - "Conciseness Over Completeness"

**但可以从 refactor 版本借鉴**:
- 标准化的 brainstorming 移交格式（如果团队需要统一格式）
- 明确的决策优先级说明

### 4. 两个版本都需要改进的地方

**Eval 1 (dry-run) 问题**:
- 应该假设合理的默认值而不是过度询问
- 应该重述预期行为
- 可以询问"如果文件不存在，是否创建？"而不是"哪个文件？"

**Eval 3 (排序) 问题**:
- 应该询问排序特定的盲点
- 应该计划检查现有排序逻辑
- 可以提供更聚焦的问题（tie-breaker、未知值处理等）

### 5. 下一步行动

**选项 A: 采用原版并优化**
1. 保留 `clarification` 作为主版本
2. 针对 eval-1 和 eval-3 的问题进行优化
3. 运行新一轮测试验证改进

**选项 B: 混合优化**
1. 以原版为基础
2. 加入 refactor 版本的标准化移交格式（可选）
3. 修复两个版本共同的问题

**选项 C: 继续对比测试**
1. 增加更多测试用例
2. 测试更多边缘场景
3. 收集实际使用反馈

---

## 测试数据

### 原始数据位置
- `clarification-workspace/iteration-6/`
- `clarification-refactor-workspace/iteration-2/`

### 测试配置
- 模型: pi 默认模型
- 每个测试用例运行 1 次
- 总共 12 次运行 (2 版本 × 3 测试 × 2 配置)
- 非交互模式 (`pi -p`)

### 评分标准
基于 `evals/clarification/evals.json` 中定义的 expectations，每个测试用例有 4 个断言。
