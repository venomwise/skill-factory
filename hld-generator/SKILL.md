---
name: hld-generator
description: >
  Generate High-Level Design (HLD) — concise, visual overview for human review. Transforms detailed 
  design.md (LLD) into scannable HLD with diagrams/tables (300-400 lines, ~50% of source). 
  Use for "write HLD", "概要设计", "高层设计", "generate tech design", "帮我写个技术设计", "生成设计文档".
  Input: design.md (from brainstorming skill) or git commit history.
---

# High-Level Design (HLD) Generator

Generate a High-Level Design document from detailed design.md or git history — a concise, visual overview for human review and architectural alignment.

## HLD vs LLD：概念与对比

**What is HLD and LLD?**

- **High-Level Design (HLD) / 概要设计**：面向人类审查者的架构文档，关注"为什么这样设计"和"核心决策理由"。使用图表和简短说明展示整体架构、技术选型、关键流程。目标是让审查者在 30 分钟内理解设计意图并评估方案合理性。

- **Low-Level Design (LLD) / 详细设计**：面向实现者（AI Agent、工程师）的精确规范，关注"如何实现"和"所有实现细节"。包含完整的 API 契约、数据库 DDL、算法步骤、边界处理逻辑。目标是让实现者能够无歧义地编写代码。

**In this skill:**
- **Input (LLD)** = design.md（详细设计，来自 brainstorming skill 或 git history）
- **Output (HLD)** = 生成的概要文档（面向审查者）

---

**Five-dimension comparison:**

| 维度 | HLD（本 skill 输出） | LLD（本 skill 输入：design.md） |
|------|---------------------|-------------------------------|
| **关注点** | WHY & WHAT（为什么做、做什么） | HOW（怎么做、实现细节） |
| **核心读者** | 开发者、架构师、审查者 | AI Agent、实现工程师 |
| **详细程度** | 高层架构、核心决策、关键流程 | 完整契约、算法步骤、边界处理 |
| **核心交付物** | Mermaid 图 + 表格 + 简短说明 | 完整 DDL/API 契约/配置/错误码 |
| **技术决策边界** | 方案选择 + 理由 | 实现约束 + 所有边界情况 |

**Key distinction**: HLD explains "why this design", LLD specifies "how to implement precisely".

## Transformation Goal

This skill performs **faithful transformation** from LLD (detailed design.md) to HLD (overview document):

**Transformation = Reorganization + Abstraction, NOT Reinterpretation**

- **Input**: design.md (LLD) — complete implementation spec from brainstorming skill or inferred from git history
- **Output**: HLD document — concise, visual overview (300-400 lines, ~50% of source)
- **Readers**: Developers, architects, reviewers (humans, not AI agents)
- **Purpose**: Enable 30-minute architectural review and design approval

## Fidelity Principles

**CRITICAL: This skill performs faithful transformation, not creative rewriting.**

When working from an existing design.md:

- **Preserve all design decisions**: Do not add alternatives, trade-offs, or rationale not present in the source
- **Preserve all technical details**: Do not simplify API contracts, database schemas, or implementation specifics
- **Preserve all semantic content**: Do not change field names, types, default values, or interface semantics
- **Preserve unique sections**: If the source has sections like Discovery, Scope Decisions, or Internal API Contracts, keep them
- **Do not invent content**: If the source lacks alternatives comparison, do not fabricate alternative solutions
- **Do not delete key information**: Configuration values, error handling rules, testing strategies must be preserved

**Transformation = reorganization + formatting, NOT reinterpretation + simplification.**

## Overview Principles

**CRITICAL: Generate an OVERVIEW document for human review, NOT a detailed implementation spec.**

The generated overview document should be:
- **Scannable**: Use diagrams, tables, and bullet points; avoid long prose
- **Concise**: Each paragraph max 5 lines; each section answers ONE question
- **High-level**: Focus on WHY and WHAT, not detailed HOW
- **Visual-first**: Every major section has at least one diagram or table
- **Reference-aware**: Point to design.md for implementation specifics

Transform detailed design.md content to overview level:

**What to KEEP**:
- Business context and pain points (2-3 sentences)
- Goals and non-goals (bullet lists)
- Key architectural decisions + rationale (table format)
- High-level data flow (Mermaid diagram + 2-3 sentences)
- Major components and responsibilities (table format)
- Critical risks and mitigations (table format)

**What to TRANSFORM**:
- Code/SQL snippets → Description + reference to design.md
- Configuration YAML blocks → Summary + reference to design.md
- Step-by-step algorithms → Approach description + key technique
- Exhaustive error codes → Category table + reference to design.md
- Detailed test cases → Test strategy summary + reference to design.md
- Implementation details → High-level mechanism + reference to design.md

**What to PRIORITIZE**:
- Diagrams over prose
- Tables over paragraphs
- Bullet points over sentences
- References over duplication

**Target length**: 300-400 lines (approximately 50% of source design.md length)

## Writing Principles (Overview Document)

A good overview document is scanned first, read second. Apply these principles:

- **Visual-first**: Start sections with diagram or table, then explanation
- **Lead with conclusion**: TL;DR at top (problem | solution | impact in 3 bullets)
- **One idea per paragraph**: Max 5 lines or 400 characters per paragraph
- **Active voice**: "cursor 采用签名防篡改" not "cursor 被设计为采用签名以防篡改"
- **Concrete terms**: "pending event + scanner 补偿" not "可靠消息投递机制"
- **Reference pattern**: Use "详见 design.md §X [Section Name]" for implementation details
- **Show trade-offs**: Alternatives table when source discusses options (never invent alternatives)
- **Zero-context reader**: Define terms, state assumptions, name the audience

**Paragraph break triggers**:
- When listing ≥ 3 items → use bullets
- When comparing ≥ 2 options → use table
- When exceeding 5 lines → split with sub-heading (####)
- When describing flow → use diagram

**Content transformation examples**:

❌ **Detailed version**:
```
cursor 是服务端编码的字符串，内部包含 version、kid、syncMode、enterpriseNum、updateTime、id、initialSyncTime、precision 和 HMAC-SHA256 签名。使用 Base64URL 编码 JSON，签名覆盖除 sig 以外的字段，字段序列化必须稳定。cursor 携带 kid，用于定位签名密钥...（继续 20+ 行）
```

✅ **Overview version**:
```
**游标机制**: 服务端编码字符串，HMAC-SHA256 签名防篡改，支持密钥轮转。

| 机制 | 说明 |
|------|------|
| 签名防篡改 | HMAC-SHA256 + kid 密钥版本标识 |
| 分页稳定性 | (update_time, id) 复合游标避免漏读 |
| 一致性保证 | safety lag 避免事务乱序提交 |
| 模式切换 | 初始同步完成后自动切换增量模式 |

详见 design.md §4.4 Cursor Service 的完整编码规则和边界处理逻辑。
```

**More transformation examples (空泛 vs 具体)**:

| 场景 | ❌ 空泛描述 | ✅ 具体描述 |
|------|-----------|-----------|
| 高可用 | "系统采用高可用架构确保稳定性" | "Redis 主从+哨兵，主节点故障 30s 内自动切换，降级读本地缓存" |
| 异常处理 | "做好异常处理防止系统崩溃" | "Redis 连接超时（50ms）时降级读 JVM 缓存，上报监控告警" |
| 并发控制 | "通过分布式锁防止并发问题" | "Redisson 实现 SKU 级锁（key: `lock:stock:{skuId}`），超时 3s" |
| 容量规划 | "系统支持高并发访问" | "Redis QPS 支持 10w+，1000 SKU × 1KB，预留 10 倍冗余" |

## Input Sources

Two supported inputs, in priority order:

1. **design.md** — produced by the `brainstorming` skill, typically at `specs/<topic>/design.md` or `.codex/specs/<topic>/design.md`
2. **git commit history** — analyze recent relevant commits via `git log` to extract change intent and technical context

## Output

- Technical design document, default path: `docs/design/<topic>.md`
- Follow existing conventions if the project already has a `docs/` or `design/` directory

## Workflow

### 1. Identify the input source

Ask the user which input source to use:
- **design.md** from the brainstorming skill
- **git commit history**

If the user chooses design.md, probe for it automatically: check `specs/` and `.codex/specs/` for an existing file. If found, use it directly. If not found, ask the user for the path.

### 2. Gather context

**Source A: design.md**

Read the file and extract:
- Feature goals and background
- Technical approach and architecture decisions
- Interface and data model details
- Confirmed constraints and assumptions

**Source B: git commit history**

```bash
git log --oneline -50
```

Filter to relevant commits by grepping the user's keyword against commit messages:

```bash
git log --oneline -50 | grep -i "<keyword>"
```

For each matched commit, inspect what changed:

```bash
git show <commit-hash> --stat
git diff <commit-hash>^ <commit-hash> -- <relevant-files>
```

Extract:
- Which modules/files changed
- Business intent from commit messages
- New or modified interfaces and data structures

Only analyze commits that match the keyword — ignore unrelated history.

### 3. Confirm missing information

Based on gathered context, identify what's still unclear and ask the user to fill gaps — one question at a time. Typical gaps to check:

- Audience: who needs to read this? (shapes how much background to spell out)
- Background: is the business or technical context sufficient for a reader with no prior knowledge?
- **Alternatives (only ask if source doesn't discuss)**: what other approaches were considered, and why was this one chosen? **If source already discusses alternatives, skip this question.**
- Scope: are there explicit non-goals worth calling out?
- API: does this involve new or modified endpoints? (determines whether the API Design section is needed)
- Database: does this involve schema changes? (determines whether the Database Design section is needed)
- Impact & risk: which modules or teams does this touch, and how would it roll back?
- Key points: any implementation gotchas or constraints the reader must know?

**Important**: Do NOT ask if the information is already clear in the source. Do NOT ask for alternatives if source already has them.

### 4. Generate the overview document

**Transform source content to overview structure: extract high-level design, use visual elements, reference details.**

#### Step 4.1: Analyze source structure

Identify what sections and content types exist in the source:
- Does it have a Discovery or Key Discoveries section?
- Does it have a Scope Decisions section?
- Does it discuss alternative solutions?
- Does it have detailed API contracts (not just examples)?
- Does it have detailed data flow walkthroughs?
- Does it have error handling strategies?
- Does it have testing plans?
- Does it have configuration/constants sections?

#### Step 4.2: Map to overview structure

**Transformation mindset**: An overview document is not a "shortened detailed doc" — it's "telling the story from a different angle".

- Detailed doc: tells AI "how to write code" (precise down to every field, every error)
- Overview doc: tells humans "why this design" (decision rationale, core approach, key constraints)

Ask yourself when transforming: If I were a reviewer, what information do I need to judge whether this design is sound?

---

Extract **overview-level** content from source and organize into template structure. Apply Overview Principles: use visual elements, keep paragraphs short, reference details.

**Section transformation rules**:

**TL;DR**: 
- Extract 3 bullets (Problem | Solution | Impact)
- Problem: 1 sentence from Background pain points
- Solution: 1 sentence from Proposed Solution core approach
- Impact: 1 sentence from Impact section or infer from scope
- Total: under 150 words

**背景说明**:
- Current state + pain points (2-3 sentences, max 5 lines)
- Why now (1-2 sentences)
- DO NOT include: technical implementation details, schema descriptions, configuration values

**目标阐述**:
- Goals: copy list as-is from source (bullet points)
- Primary Users/Roles: copy list as-is if present
- Non-goals: copy list as-is from source (bullet points)

**关键发现** (only if source has Discovery section):
- Preserve source's Key Discoveries section
- Use bullet points, max 2 levels nesting
- Keep concise: one finding per bullet
- DO NOT expand or explain; just list the discoveries

**设计决策** (only if source has Scope Decisions):
- Preserve source's Scope Decisions section
- Transform to table format if comparing options:
  | 决策点 | 选择 | 理由 |
- Use bullet points for single decisions
- Keep rationale brief (1 sentence per decision)

**方案设计 — 整体思路**:
- Architecture diagram: Mermaid flowchart or sequenceDiagram (max 8-10 nodes)
- 2-3 sentences explaining diagram (max 5 lines)
- DO NOT include: step-by-step implementation, detailed algorithms

**方案设计 — 方案对比** (ONLY if source discusses alternatives):
- Create comparison table:
  | 方案 | 优点 | 缺点 | 是否采用 |
- Include actual alternatives from source
- Add decision rationale (2-3 sentences from source)
- **NEVER invent alternatives to fill this section**
- If source doesn't discuss alternatives, SKIP this section entirely

**设计说明 — 核心组件**:
- Components table:
  | 组件 | 职责 | 依赖 |
- Extract from source's Components or Proposed Solution
- Keep responsibility description to 1 sentence per component

**设计说明 — 数据模型** (only when source contains schema):
- Mermaid erDiagram showing key entities and relationships
- OR table format with key fields:
  | 表名 | 关键字段 | 设计要点 |
- List key design points (indexes, constraints) in 2-3 bullets
- Reference: "详见 design.md §X Database Design 的完整 DDL"
- DO NOT include: complete SQL DDL, all field definitions

**设计说明 — 接口设计** (only when source contains API specs):
- API table:
  | 接口 | 方法 | 用途 | 关键参数 |
- Error handling strategy table:
  | HTTP Status | 类别 | 重试策略 |
- DO NOT include: complete request/response schemas, all error codes
- Reference: "详见 design.md §X API Design 的完整契约"
- **Preserve API field names exactly** as in source (for fidelity)

**设计说明 — 关键流程**:
- Include 1-2 most important flows only
- Use Mermaid sequenceDiagram (max 5 actors)
- Add 2-3 sentences highlighting critical points (max 5 lines)
- Reference: "详见 design.md §X Data Flow 的完整边界处理逻辑"
- DO NOT include: all scenarios, detailed error handling in diagram

**设计说明 — 配置项** (only if source has Configuration):
- Summarize configuration strategy in 2-3 bullets
- DO NOT include: complete YAML blocks, all parameter values
- Reference: "详见 design.md §X Configuration 的完整配置项"

**错误处理** (only if source has Error Handling):
- Error category table:
  | 场景 | HTTP Status | Code | 处理策略 |
- Keep to major categories (5-8 rows max)
- Reference: "详见 design.md §X Error Handling 的完整错误码表"
- DO NOT include: all error codes, detailed retry logic

**测试策略** (only if source has Testing):
- Test category summary (bullet list):
  - Category 1: [brief description]
  - Category 2: [brief description]
- Reference: "详见 design.md §X Testing 的完整测试用例"
- DO NOT include: detailed test cases, all test scenarios

**影响范围与风险**:
- Impact scope: bullet lists (services / data / teams)
- Key risks table:
  | 风险 | 可能性 | 影响 | 缓解措施 |
- Keep to top 3-5 risks
- **Demonstrate deep thinking**:
  * 性能瓶颈：数据库死锁、缓存击穿、热点数据
  * 安全隐患：越权访问、防刷、重放攻击
  * 工程边界：第三方依赖超时、网络抖动、降级策略
- Preserve rollback strategy if mentioned in source

**关键点说明**:
- Extract 3-5 most critical implementation constraints
- Use bullet points, 1 line per point
- Focus on "would surprise the implementer" constraints
- DO NOT include: all implementation details

**待讨论/开放问题** (optional):
- Include if present in source
- Keep brief (bullet list)

**References section** (NEW, always include):
```markdown
---

## 参考文档

- **详细设计规范**: [design.md](path/to/design.md) — 完整的实现规范，供 AI Agent 编码和工程师实现时参考
  - 包含：完整 API 契约、数据库 DDL、配置参数、错误码表、测试用例、边界处理逻辑
- **相关文档**: [if any from source]
```

#### Step 4.3: Apply visualization rules

Every major section MUST have at least ONE visual element:

**Architecture/Flow → Mermaid Diagram**:
- Overall architecture: `flowchart TD` (max 8-10 nodes)
- Interaction flows: `sequenceDiagram` (max 5 actors)
- Data model: `erDiagram` OR table with key fields
- State transitions: `stateDiagram-v2`
- Keep diagrams focused; add 2-3 sentence explanation (not paragraph)

**Comparison/Structure → Table**:
- Alternatives comparison: | 方案 | 优点 | 缺点 | 是否采用 |
- Components: | 组件 | 职责 | 依赖 |
- API endpoints: | 接口 | 方法 | 用途 | 关键参数 |
- Error strategy: | HTTP Status | 类别 | 重试策略 |
- Risks: | 风险 | 可能性 | 影响 | 缓解措施 |
- Configuration: | 类别 | 关键参数 | 说明 |

**Lists → Bullet Points**:
- Goals / Non-goals (max 2 levels nesting)
- Impact scope (services / data / teams)
- Test categories (strategy level only)
- Implementation constraints

**Code/SQL → NEVER in overview**:
- NEVER include code snippets
- NEVER include SQL DDL (use table description instead)
- NEVER include JSON examples (use field list in table instead)
- ALWAYS replace with: brief description + reference to design.md

**Configuration blocks → Summary + Reference**:
- NEVER include complete YAML blocks
- Summarize: "支持配置 safety lag、限流阈值、墓碑保留窗口等参数"
- Reference: "详见 design.md §X Configuration"

#### Step 4.4: Fidelity checklist

Before finalizing, verify:
- [ ] No design decisions added that weren't in source
- [ ] No alternative solutions invented to fill template
- [ ] All API field names match source exactly (for fields that are shown)
- [ ] Unique source sections preserved (Discovery, Scope Decisions, etc.)
- [ ] Technical details not simplified where shown (but relegated to references)

### Choosing the right Mermaid diagram

| Intent | Diagram type | Max elements |
|--------|--------------|--------------|
| Component/architecture relationships, data flow | `flowchart TD` | 8-10 nodes |
| Interaction ordering across services/actors | `sequenceDiagram` | 5 actors |
| Lifecycle / status transitions | `stateDiagram-v2` | 6-8 states |
| Data model, tables and relationships | `erDiagram` | 4-6 entities |

Keep diagrams focused on the core flow. Omit error paths and edge cases from diagrams—mention them in text or reference design.md.

**When NOT to use diagrams**:
- Simple CRUD interface lists → use tables
- Single-direction dependency (A calls B) → text description
- Only 2-3 states → use text or table

**Diagram quality standards**:
- Keep nodes/actors within max limits (avoid information overload)
- Add 2-3 sentence explanation (highlight key points, don't repeat diagram content)
- Error paths and edge cases not in diagram—mention in text or reference design.md

### 5. User review

Present the draft to the user:
- Wording or detail changes → edit in place
- Approach or scope changes → return to step 3
- Missing context → return to step 2

Write the file to the target path only after the user approves.

## Verification

### 概要质量检查
- [ ] 【TL;DR】3 bullets (Problem | Solution | Impact), under 150 words
- [ ] 【视觉化】每个主要章节 ≥1 图表或表格（除 TL;DR、References）
- [ ] 【段落】无超 5 行段落，无代码/SQL/完整配置
- [ ] 【引用】技术细节以"详见 design.md §X"结尾
- [ ] 【长度】300-400 行（约 50% 源文档）

### 忠实性检查
- [ ] 【语义】API 字段名/类型与源文档一致（已展示部分）
- [ ] 【决策】无臆造备选方案（只在源讨论时才包含）
- [ ] 【章节】保留源文档特有章节（Discovery/Scope Decisions等）
- [ ] 【完整性】无添加源文档不存在的设计决策
- [ ] 【准确性】技术细节未扭曲（简化呈现可以，错误信息不行）

### 最终确认
- [ ] 文档已写入目标路径并获用户批准

## References

- [Document template](assets/tech-design-template.md)
- [HLD 编写指南](references/hld-writing-guide.md) — HLD vs LLD 区别、转换标准、深度思考体现
