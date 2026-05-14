澄清：我在当前目录查看到 `README.md` 明确说明业务优先级顺序由 `PRIORITY_ORDER` 定义：`urgent -> high -> normal -> low`。

假设/默认：
- “按业务优先级排”理解为优先使用 `PRIORITY_ORDER` 的数值从小到大排序。
- 同一优先级内默认按 `created_at` 从小到大稳定排序。
- 未知 priority 默认排在已知优先级之后。

实现总结：
- 已修改 `src/list_sort.py`
- `sort_items()` 现在按业务优先级排序，而不是按 `created_at` 倒序排序。

验证：
- 先尝试使用 `python`，当前环境不存在该命令。
- 使用 `python3` 执行了快速验证，输出为：

```text
['urgent-old', 'urgent-new', 'high', 'normal', 'low']
```

结果符合 `urgent -> high -> normal -> low` 的业务优先级顺序。
