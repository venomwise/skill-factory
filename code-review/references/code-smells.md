# Code Smell Baseline

Standards Agent 仅在 Full Review 中读取本文件。以下条目源自 Martin Fowler《Refactoring》
第 3 章，是跨项目的启发式检查，不是 repository guidance。

适用规则：

- Repository guidance 优先；项目明确认可的写法不得因本清单产生 finding。
- Smell 必须标为启发式判断，默认最高为 `minor`。只有存在可观察的 behavior、reliability
  或交付影响时，才可建议更高严重性并提供因果证据。
- 工具通过时，不为同一规则重复创建 smell finding；工具失败时，将实际失败作为证据。
- 只检查本次 diff 及其必要调用上下文，不把未受改动影响的历史问题带入评审。

## Smells

- **Mysterious Name**：函数、变量或类型的名字无法说明职责或数据含义。优先改成准确名称；
  无法准确命名通常说明设计边界仍不清楚。
- **Duplicated Code**：本次改动的多个 hunk 或文件出现相同逻辑形态。提取共享逻辑并复用。
- **Feature Envy**：方法访问其他对象的数据多于自身数据。考虑把行为移到它依赖的数据一侧。
- **Data Clumps**：同一组字段或参数反复一起传递。考虑将它们组合成表达领域含义的类型。
- **Primitive Obsession**：primitive 或 string 代替了应有独立约束的领域概念。考虑建立小型领域类型。
- **Repeated Switches**：针对同一类型的 `switch` 或 `if` 分支在多处重复。考虑使用 polymorphism，
  或集中为多处共享的一份映射。
- **Shotgun Surgery**：一个逻辑变化迫使多个文件分散修改。考虑把一起变化的职责收拢到同一模块。
- **Divergent Change**：一个文件或模块因多个无关原因同时变化。考虑按变化原因拆分职责。
- **Speculative Generality**：新增 spec 未要求的 abstraction、参数或 hook。删除无现实需求的扩展点，
  回到当前需求所需的最小实现。
- **Message Chains**：调用方依赖较长的 `a.b().c().d()` 导航链。考虑由第一个对象封装内部导航。
- **Middle Man**：class 或函数主要把调用转发给下一个对象。考虑移除中间层并直接调用真实目标。
- **Refused Bequest**：subclass 或 implementer 忽略或覆盖大部分继承能力。考虑放弃 inheritance，
  改用 composition。
