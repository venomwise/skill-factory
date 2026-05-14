澄清：我理解“按业务优先级排”是指使用项目中已定义的 `PRIORITY_ORDER` 排序：`urgent -> high -> normal -> low`。

假设/默认：
- 同一优先级内仍按 `created_at` 倒序排列，保持原先“新创建的在前”的展示习惯。
- `priority` 字段应是 `PRIORITY_ORDER` 中已有值。

实现总结：
- 修改了 `src/list_sort.py`
- `sort_items()` 现在按 `(业务优先级, created_at 倒序)` 排序，而不是只按 `created_at` 倒序排序。

验证：
- 使用 `python3` 执行了一个临时样例验证。
- 结果为：`urgent-new, urgent-old, high, low-new`
- 符合业务优先级优先、同优先级按时间倒序的预期。
