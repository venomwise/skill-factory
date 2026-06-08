# Tech Design Doc Skill - 修复日志

## 2026-06-08: 修复忠实性转换问题

### 问题描述

Skill 在将 design.md 转换为技术设计文档时，出现了严重的语义篡改问题：

1. **臆造内容**：生成了原始设计中不存在的备选方案对比表
2. **改变接口语义**：修改了 API 字段名、类型、默认值
3. **删除关键信息**：丢失了 Discovery、Scope Decisions、Internal API Contracts 等章节
4. **简化技术细节**：大幅简化配置项、错误处理、测试策略等内容

**根本原因**：Skill 强制要求按模板"创作"文档，而不是"忠实转换"原始内容。

### 修复内容

#### 1. 新增"忠实性原则"章节

在 `SKILL.md` 开头增加了 **Fidelity Principles**：

```markdown
## Fidelity Principles

**CRITICAL: This skill performs faithful transformation, not creative rewriting.**

- **Preserve all design decisions**: Do not add alternatives, trade-offs, or rationale not present in the source
- **Preserve all technical details**: Do not simplify API contracts, database schemas, or implementation specifics
- **Preserve all semantic content**: Do not change field names, types, default values, or interface semantics
- **Preserve unique sections**: If the source has sections like Discovery, Scope Decisions, or Internal API Contracts, keep them
- **Do not invent content**: If the source lacks alternatives comparison, do not fabricate alternative solutions
- **Do not delete key information**: Configuration values, error handling rules, testing strategies must be preserved

**Transformation = reorganization + formatting, NOT reinterpretation + simplification.**
```

#### 2. 改造 Workflow 第 4 步

从"按模板生成"改为"忠实映射"，增加了三个子步骤：

**Step 4.1: Analyze source structure** - 分析源文档包含哪些章节和内容类型

**Step 4.2: Map to template structure** - 映射源内容到模板结构，关键原则：
- 保留源文档特有的章节（Discovery, Scope Decisions, Internal API Contracts, Data Flow, Error Handling, Testing, Configuration）
- **ONLY include alternatives if source discusses alternatives** - 只有源文档讨论了备选方案才生成对比表
- **Preserve all API contracts exactly** - 保持 API 契约完全一致，不改字段名/类型/默认值
- **Preserve all technical details** - 保留所有技术细节

**Step 4.3: Fidelity checklist** - 生成前验证忠实性，包含 8 项检查点

#### 3. 更新模板文件

`assets/tech-design-template.md` 增加了以下可选章节：

- **3. 关键发现（可选）** - 保留源文档的 Discovery 章节
- **4. 设计决策（可选）** - 保留源文档的 Scope Decisions 章节
- **6.4 内部接口契约（可选）** - 保留源文档的 Internal API Contracts 章节
- **6.5 数据流（可选）** - 保留源文档的 Data Flow 章节
- **6.6 配置项（可选）** - 保留源文档的 Configuration 章节
- **7. 错误处理（可选）** - 保留源文档的 Error Handling 章节
- **8. 测试策略（可选）** - 保留源文档的 Testing 章节

所有可选章节都明确标注：**仅当源文档包含该章节时才保留**。

#### 4. 更新验证清单

Verification 部分增加了"忠实性验证"子清单：

```markdown
- [ ] **Fidelity verification**:
  - [ ] All API contracts match source exactly (no field name/type/default changes)
  - [ ] All configuration values preserved from source
  - [ ] All error handling rules preserved from source
  - [ ] All testing strategies preserved from source
  - [ ] Unique source sections preserved (Discovery, Scope Decisions, etc.)
  - [ ] No alternative solutions invented
  - [ ] Technical details not simplified or omitted
```

### 修复效果

修复后，Skill 将：

✅ **保留原始设计的所有语义**
- 不改变 API 字段名、类型、默认值
- 不删除配置项、错误处理规则、测试策略
- 不简化技术细节

✅ **保留原始文档的特有章节**
- Discovery / Key Discoveries
- Scope Decisions
- Internal API Contracts
- Data Flow
- Error Handling
- Testing
- Configuration

✅ **只在源文档存在时才生成对应内容**
- 只有源文档讨论了备选方案，才生成方案对比表
- 不臆造源文档中不存在的内容

✅ **转换 = 重组 + 格式化，而非重新诠释 + 简化**
- 保持源文档的所有技术决策和理由
- 只调整章节顺序和格式，不改变内容语义

### 验证建议

使用修复后的 skill 重新处理之前的 design.md，对比生成结果：

1. 检查 API 契约是否完全一致（字段名、类型、默认值）
2. 检查是否保留了所有关键章节（Discovery、Scope Decisions、Internal API Contracts 等）
3. 检查配置项、错误处理规则、测试策略是否完整保留
4. 检查是否臆造了源文档中不存在的备选方案
5. 检查技术细节是否被简化或删除
