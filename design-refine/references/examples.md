# Design Refine Examples

本文件只覆盖主流程之外的边界场景。标准分类和写入规则以 `SKILL.md` 为准。

## Direct repair batch

Review 有以下 findings：

- F-004：Testing 漏掉一个已存在 AC 的验证映射。
- F-006：Interfaces 中新增 PATCH 被错误标为 MODIFY。
- F-009：Data Model 字段长度与已核实 SQL 不一致。

三项都能从现有 Decision、AC 和 SQL 唯一推导，不需要新设计选择。

向用户展示：

```markdown
以下 3 项属于 Direct repair：

1. F-004：§Testing 补充 AC 引用。
2. F-006：§Interfaces 把 metadata PATCH 标为 ADD。
3. F-009：§Data Model 按 SQL 修正字段长度。

是否统一应用这 3 项修复？
```

确认后更新章节，三个 finding 分别指向实际章节，不创建 Decision。

## Explicit boundary rejection

### Scenario

F-012 建议为无状态 preview 增加 Redis session、claim 和聚合配额。
用户确认 preview 只处理当前请求，commit 不消费 preview 结果，并明确拒绝该建议。

### Persisted decision

```markdown
### DR-preview-remains-stateless

**Decision**: preview 不保存服务端会话，commit 不消费 preview 结果。

**Rationale**: preview 只承担当前请求的解析和校验；跨请求状态不解决当前业务问题。

**Constraints**: commit 必须根据自身完整请求独立校验。

**Rejected concern**: 不为 preview 增加 session、claim、租约或聚合配额。

**Revisit when**: preview 结果需要被 commit 消费，或服务端开始跨请求保存状态。
```

F-012 状态设为 `rejected`，Resolution ref 指向 `DR-preview-remains-stateless`。
后续 reviewer 在 Revisit 条件未触发时不得重新提出同一 concern。

## Reusing an existing boundary

下一轮 F-018 再次建议增加 preview TTL，但 `DR-preview-remains-stateless` 仍有效。

处理方式：

1. 不创建新 Decision。
2. 如果 F-018 与 F-012 是同一根因，沿用 F-012，而不是保留 F-018。
3. Closure 说明已有 Decision 覆盖且 Revisit 条件未触发。
4. 必要时原位补强现有 Decision 的措辞。

## Merging findings into one decision

以下 findings 指向同一个设计选择：

- F-021：错误响应格式未确定。
- F-022：批量错误是否逐项返回未确定。
- F-024：客户端重试需要稳定错误分类。

队列合并为：

```text
统一错误契约（来源: F-021, F-022, F-024）
```

用户只需确认一次错误契约。三个 finding 保留独立状态，共同指向同一个
`DR-import-error-contract` 和对应 Interfaces / Error Handling 章节。

## Finding Closure Proof

F-014 指出两个连续接口共享同一批数据，但只定义了第一个接口的容量边界。refine 补充第二个接口后，
不能仅凭章节已修改就标记闭合。pre-closure audit 应记录：

| Finding | Closure test | Resolution evidence | Counterexample check | Result |
|---------|--------------|---------------------|----------------------|--------|
| F-014 | 两个接口的容量责任均明确，差异可观察 | DR、两个接口、Error Handling、AC、Testing | 第一个接口接受而第二个因隐式上限失败时，契约已定义错误和恢复方式 | passed |

如果第二个接口仍没有上限、错误或恢复语义，反例仍成立，应重新打开 F-014，不能写 `passed`。

## Next review mode

refine 只补充遗漏的 Error Handling 和 Testing 时，`Next review mode` 为 `Closure`，按 changed sections
执行范围化语义预检。

refine 改变 Goals、Non-Goals 或公共容量契约时，`Next review mode` 为 `Full`。pre-closure audit 必须
先执行全维度 dry-run，再交给 `design-review` 独立 Full Review；不能因为原 finding 已终态而缩小范围。

## Updating an existing decision

已有：

```markdown
### DR-cache-location

**Decision**: 使用进程内缓存。

**Rationale**: 当前为单实例部署。

**Constraints**: 不支持多副本共享缓存。
```

后续项目事实变为多副本部署时，原位更新同一个 ID：

```markdown
### DR-cache-location

**Decision**: 使用共享 Redis 缓存。

**Rationale**: 多副本必须读取一致缓存；原进程内方案的单实例前提已失效。

**Constraints**: Redis 不可用时按 Error Handling 的降级规则处理。
```

不要保留旧正文后再追加 `Revised`。旧选择及变化原因由 review Closure 和版本历史追踪。

## Closure origins

### refine-regression

F-030 的修复新增状态 `RECONCILING`，但 State Machines 没有该状态。
新 finding 可以使用 `refine-regression`，证据是 changed section 引入了未传播状态。

### dependency-unlocked

F-031 先确定数据以哪个存储为准。只有该决策完成后，才能判断回滚是否需要 backfill 版本。
后续 finding 可以使用 `dependency-unlocked`，并说明依赖关系。

### baseline-miss

接口字段在 R1 前已经含糊，但 R1 漏检。Closure 必须明确标记 `baseline-miss`，
不能把 reviewer 漏检包装成 refine 新问题。

### context-change

上轮后项目从单实例改为多副本。与缓存一致性相关的新 finding 使用
`context-change`，并引用项目配置或提交证据。

无法归入以上 Origin 的自由优化建议不进入 Closure Findings。
