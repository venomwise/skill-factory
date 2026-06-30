# 使用场景示例

**场景 1: 首次审查支付功能 PRD**
```
用户：/prd-review docs/payment-feature.md
→ 本地没有支付相关代码
→ web_search "支付系统幂等性设计" "支付对账流程"
→ 生成 docs/payment-feature-review.md
→ 发现 5 个逻辑漏洞：
  - [P0-阻断] 重复支付如何处理？PRD 未说明
  - [P0-阻断] 支付失败后订单状态变更逻辑缺失
  - [P1-需确认] 退款流程未描述
  - [P1-需确认] 对账周期和异常处理未明确
  - [P2-优化] 支付超时时间建议明确
→ Todo List 8 个问题
```

**场景 2: PM 补充说明后，第二次审查**
```
用户：/prd-review docs/payment-feature.md
      PM 说：重复支付通过订单号去重，退款走异步队列，对账每日凌晨执行
→ Read docs/payment-feature-review.md（已有的分析）
→ Edit 更新报告：
  - 标记"重复支付"问题为 [已解决]
  - 标记"退款流程"问题为 [已解决]
  - 标记"对账周期"问题为 [已解决]
→ 新增发现：异步队列的失败重试机制未说明
→ 更新 Todo List
```

**场景 3: 所有疑问解决**
```
用户：/prd-review docs/payment-feature.md
→ Read docs/payment-feature-review.md
→ 检查 PRD 是否有更新（用户未说明，假设无更新）
→ 确认所有 Todo 已勾选，无阻断性风险
→ 更新评审状态为 🟢 已通过
→ 输出："需求澄清完成，所有逻辑闭环和技术可行性问题已确认，可以开始设计技术方案。"
```

**场景 4: 涉及数据库的需求**
```
用户：/prd-review docs/user-report.md
      PRD 要求：生成用户行为报表，按时间范围查询，支持导出
→ 读取 PRD
→ 检索现有报表相关代码
→ 通过 subagent 执行 db-explorer skill 检查：
  - users 表结构
  - user_behaviors 表是否有时间索引
  - 数据量级（如果有百万级数据，全表扫描不可行）
→ 发现：user_behaviors 表无 created_at 索引，大数据量查询会超时
→ 生成 docs/user-report-review.md，标记 [P0-阻断] 性能风险
```
