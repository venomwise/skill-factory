# Brainstorming Skill Test Cases

Use these prompts to manually validate that the `brainstorming` skill diagnoses ambiguity correctly, controls scope, and converges before handing off to `spec-plan`.

## Test 1 — Fuzzy request: notifications

**Prompt**

> 我想给我的项目加个通知系统，用户能收到各种提醒

**Expected behavior**

- Diagnose the request as **Problem unclear**.
- Clarify what reminders exist, who receives them, what triggers them, and which delivery channels matter.
- Ask questions one at a time instead of dumping a checklist.
- Surface at least one assumption and one blind spot before proposing a solution.
- Converge to a scoped problem statement before comparing approaches.

**Primary coverage**

- Fuzziness diagnosis
- Clarification flow
- Assumption challenging
- Blind-spot scanning

## Test 2 — Large multi-domain request: user system

**Prompt**

> 我要做一个完整的用户系统，包括注册登录、权限管理、用户画像分析和社交关系

**Expected behavior**

- Recognize that the request spans multiple independent subsystems.
- Propose a decomposition before brainstorming details.
- Suggest a reasonable split such as:
  - authentication and session management
  - authorization and roles
  - profile / analytics
  - social graph / relationships
- Ask the user to confirm which sub-project to tackle first.
- Brainstorm only the first confirmed sub-project and keep the others out of scope for that design cycle.

**Primary coverage**

- Decomposition heuristics
- Scope control
- Single-spec-cycle discipline

## Test 3 — Technical approach unclear: realtime updates

**Prompt**

> 我需要给现有的 REST API 加实时推送能力，数据变更时客户端能立即收到更新

**Expected behavior**

- Diagnose the request as **Solution unclear**.
- Skip broad ideation and move into approach comparison.
- Compare 1-3 viable options such as WebSocket, SSE, and polling.
- Lead with a recommendation and explain trade-offs.
- Ask the user to confirm the refined problem summary and the chosen approach before moving on.

**Primary coverage**

- Approach comparison
- Recommendation with trade-offs
- Convergence before design writing

## Test 4 — Large fuzzy request with scale pressure: collaborative docs

**Prompt**

> 我要做一个支持十万人同时在线的协同文档

**Expected behavior**

- Diagnose the request as **Problem unclear** before suggesting architecture.
- Clarify core boundaries such as:
  - 这十万人是读多写少，还是读写并发都很高？
  - 服务器预算和成本约束是多少？
  - 团队对技术栈、基础设施、部署方式有什么偏好？
- Avoid jumping straight to CRDT / OT / database / queue choices before the workload shape is known.
- After constraints become clearer, decide whether the request should be decomposed into smaller sub-projects.

**Primary coverage**

- Problem clarification under scale-heavy language
- Constraint discovery
- Delay architecture decisions until the shape of the problem is known
- Decomposition after clarification when needed

## Coverage Matrix

| Test | Main ambiguity type | Main behavior being validated |
|------|----------------------|-------------------------------|
| 1 | Problem unclear | Clarify goals before solutioning |
| 2 | Boundaries unclear + oversized scope | Split into sub-projects before design |
| 3 | Solution unclear | Compare approaches and recommend |
| 4 | Problem unclear + oversized scope | Clarify scale constraints before architecture |

## Pass Criteria

A successful run of the skill should consistently:

- identify the dominant ambiguity type early
- ask focused questions one at a time
- avoid premature architecture decisions
- reduce scope when the request is too large
- converge to a design-worthy problem statement before invoking `spec-plan`
