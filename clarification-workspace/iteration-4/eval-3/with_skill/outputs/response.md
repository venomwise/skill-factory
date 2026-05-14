澄清：你提到“按业务优先级排”，我检查了当前目录代码和 README，发现业务优先级已由 `PRIORITY_ORDER` 定义。

假设/默认：
- 优先级顺序为：`urgent` → `high` → `normal` → `low`
- 同一优先级内，继续按 `created_at` 从新到旧排序，保留原先“新创建优先”的展示逻辑

实现总结：
- 修改了 `src/list_sort.py`
- `sort_items()` 现在按 `(业务优先级, 创建时间倒序)` 排序
- 原来只按 `created_at` 倒序排序，导致低优先级的新数据可能排在高优先级前面

验证：
- 运行了一个本地 Python 验证脚本
- 验证结果符合预期：`urgent` 优先于 `high/normal/low`，且 `urgent` 内部按时间倒序排列
- 注意：当前环境没有 `python` 命令，使用 `python3` 验证通过。
