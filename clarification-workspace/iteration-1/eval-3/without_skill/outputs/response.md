澄清：我查看了当前目录下的说明和代码，`README.md` 中明确写到业务优先级顺序定义在 `PRIORITY_ORDER` 中。

假设/默认：
- 业务优先级排序为：`urgent` → `high` → `normal` → `low`
- 同一优先级内，继续保持原来的“创建时间倒序 / 最新优先”
- 未知优先级默认排在已知优先级之后

实现总结：
- 修改了 `src/list_sort.py`
- `sort_items()` 现在先按 `created_at` 倒序稳定排序，再按 `PRIORITY_ORDER` 稳定排序
- 这样可以保证优先级是主排序，同优先级时仍然最新在前

验证：
- 尝试使用 `python` 执行验证时，系统中没有 `python` 命令
- 使用 `python3` 执行了手动验证
- 验证结果符合预期：`urgent` 项排在最前，同为 `urgent` 时更新的排在前面，未知优先级排在最后
