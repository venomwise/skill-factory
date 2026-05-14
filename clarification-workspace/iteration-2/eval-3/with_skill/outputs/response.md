已在当前目录完成修改。

## 澄清
- 我根据 `README.md` 中说明确认：业务优先级顺序定义在 `PRIORITY_ORDER`。
- 当前优先级顺序为：`urgent` → `high` → `normal` → `low`。

## 假设 / 默认
- 列表应优先按 `priority` 的业务优先级排序。
- 同一优先级内，保留原来的按 `created_at` 倒序排序，也就是较新的排在前面。
- 未知 `priority` 默认排在已知业务优先级之后。

## 实现总结
修改了 `src/list_sort.py`：

- 原来只按 `created_at` 倒序排序。
- 现在排序 key 改为：
  1. `PRIORITY_ORDER` 中定义的业务优先级
  2. 同优先级下按 `created_at` 倒序

## 验证
执行了一个临时 Python 验证脚本：

- `urgent` 被排在 `normal`、`low` 之前
- 同为 `urgent` 时，较新的 `created_at` 排在前面
- 未知优先级排在最后

验证通过。首次尝试 `python` 命令不可用，随后使用 `python3` 验证成功。
