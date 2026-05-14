澄清：当前 README 已说明业务优先级由 `PRIORITY_ORDER` 定义，顺序为 `urgent > high > normal > low`。

假设/默认：
- “按业务优先级排”理解为优先按 `PRIORITY_ORDER` 升序排序。
- 未知 `priority` 默认排在已知优先级之后。
- 同一优先级内使用 `created_at` 升序作为稳定的次级排序依据。

实现总结：
- 修改了 `src/list_sort.py`
- `sort_items()` 现在按 `PRIORITY_ORDER` 排序，而不是按 `created_at` 倒序排序。

验证：
- 尝试运行 `python -m pytest -q`，环境中没有 `python` 命令。
- 尝试运行 `python3 -m pytest -q`，环境中未安装 `pytest`。
- 使用 `python3` 手动验证排序结果：
  - 输入包含 `low / urgent / normal / high`
  - 输出为 `urgent / high / normal / low`
  - 符合业务优先级顺序。
