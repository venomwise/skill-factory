# 库存管理系统 PRD

## 1. 背景
电商平台需要一个库存管理系统,支持商品库存的实时更新。

## 2. 功能需求

### 2.1 库存查询
- 查询商品当前库存
- 支持批量查询

### 2.2 库存扣减
- 用户下单时扣减库存
- 扣减失败时返回库存不足提示

### 2.3 库存回补
- 订单取消时回补库存
- 订单超时未支付时回补库存

### 2.4 库存预警
- 库存低于预警阈值时,在运营后台显示预警消息
- 预警阈值: 库存 < 10 件
- 定时任务每 30 分钟扫描一次
- 避免重复预警: 同一商品 24 小时内只预警一次

## 3. 数据模型

### 3.1 商品库存表
- product_id: 商品ID(主键)
- stock: 当前库存
- reserved_stock: 预占库存
- updated_at: 更新时间
- **索引:** PRIMARY KEY (product_id)
- **隔离级别要求:** REPEATABLE READ

### 3.2 订单表
- order_id: 订单ID(主键)
- product_id: 商品ID
- status: 订单状态(待支付/已支付/已取消)
- created_at: 创建时间
- timeout_at: 超时时间(created_at + 30分钟)
- **索引:** INDEX idx_timeout (status, timeout_at)

### 3.3 库存预警历史表
- alert_id: 预警ID(主键)
- product_id: 商品ID
- stock: 库存数量
- alert_time: 预警时间
- **索引:** INDEX idx_product_time (product_id, alert_time)

## 4. 业务规则

### 4.1 库存扣减规则
- 下单时先预占库存(reserved_stock +1, stock -1)
- 支付成功后确认扣减(reserved_stock -1)
- 订单取消时释放预占(reserved_stock -1, stock +1)

#### 4.1.1 并发控制
- **并发控制方案:** 使用原子化 UPDATE 语句保证并发安全
  ```sql
  UPDATE product_inventory 
  SET stock = stock - 1, reserved_stock = reserved_stock + 1
  WHERE product_id = ? AND stock > 0;
  ```
- **一致性保证:** 依赖 MySQL InnoDB 行锁机制,支持分布式多实例部署
- **隔离级别要求:** REPEATABLE READ(MySQL 默认级别)
- **性能限制:** 单商品并发 QPS < 500

#### 4.1.2 超时释放机制
- 超时时间: 30 分钟
- 检测方式: 定时任务每分钟扫描超时订单
- 扫描条件: status='待支付' AND timeout_at < NOW()
- 释放操作: reserved_stock -1, stock +1
- 冗余保障: 支付接口调用时也检查订单是否超时

### 4.2 超卖防护
- 库存不足时拒绝下单
- stock >= 0
- **分布式一致性:** MySQL 行锁天然支持多实例部署,无需额外分布式锁
- **适用场景:** 单商品并发 QPS < 500

### 4.3 库存预警规则
- 预警阈值: stock < 10
- 检测频率: 定时任务每 30 分钟扫描一次
- 通知方式: 运营后台站内消息
- 防重通知: 同一商品 24 小时内只预警一次

## 5. 非功能需求

### 5.1 性能要求
- 支持 1000 QPS(总量,单商品并发预期 < 500 QPS)
- 库存数据准确率 100%

### 5.2 技术架构
- **数据存储:** 纯 MySQL 方案,不引入 Redis 或消息队列
- **数据库要求:** 
  - MySQL 8.0+
  - InnoDB 存储引擎
  - 隔离级别 REPEATABLE READ
- **扩展性:** 支持多实例水平扩展,依赖 MySQL 行锁保证一致性

## 6. 未来演进

### 6.1 性能优化
- **秒杀场景:** 如出现单商品 QPS > 1000 的秒杀场景,可引入 Redis 预扣减方案
- **高并发优化:** 如总 QPS 超过 5000,可考虑读写分离或分库分表

### 6.2 功能增强
- **预警增强:** 可增加邮件、短信等多渠道通知方式
- **超时优化:** 如订单量超过每日 10 万单,可升级为延迟消息队列方案
