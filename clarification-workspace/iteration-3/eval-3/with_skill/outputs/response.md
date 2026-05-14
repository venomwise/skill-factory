## 澄清
我检查了当前目录的 `README.md` 和 `src/list_sort.py`。README 明确说明业务优先级顺序定义在 `PRIORITY_ORDER` 中：`urgent > high > normal > low`。

## 假设/默认
- “按业务优先级排”指优先使用 `PRIORITY_ORDER` 排序。
- 同一优先级内，保留原有逻辑：按 `created_at` 倒序，即最新的排前面。
- 未知 `priority` 默认排在所有已知优先级之后。

## 实现总结
已修改 `src/list_sort.py`：

- 原逻辑：仅按 `created_at` 倒序排序。
- 新逻辑：先按 `PRIORITY_ORDER` 升序排序，再按 `created_at` 倒序排序。

核心行为现在是：

```python
urgent -> high -> normal -> low -> unknown
```

## 验证
已运行：

```bash
python3 -m py_compile src/list_sort.py
```

并用示例数据验证排序结果：

```text
['urgent-new', 'urgent-old', 'high-new', 'normal-old', 'low']
```

验证通过。
