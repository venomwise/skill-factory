澄清：你提到“按业务优先级排”，我在项目说明中看到业务优先级已由 `PRIORITY_ORDER` 定义：`urgent > high > normal > low`。

假设/默认：
- 默认列表展示应按 `PRIORITY_ORDER` 升序排序，即数值越小优先级越高。
- 同优先级内暂按 `created_at` 升序作为稳定的次级排序字段。
- `priority` 字段必须存在且值在 `PRIORITY_ORDER` 中。

实现总结：
- 修改了 `src/list_sort.py`
- 将原来的按 `created_at` 倒序排序，改为按业务优先级排序：

```python
key=lambda item: (PRIORITY_ORDER[item['priority']], item['created_at'])
```

验证：
- 用 `python3` 手动验证了排序结果。
- 输入包含 `normal / urgent / high / low` 的列表，输出顺序为：

```text
['urgent', 'high', 'normal', 'low']
```

符合业务优先级顺序。
