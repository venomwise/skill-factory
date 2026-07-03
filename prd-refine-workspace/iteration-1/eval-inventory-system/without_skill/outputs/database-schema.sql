-- 数据库表结构

-- 1. 商品库存表
CREATE TABLE product_stock (
    product_id BIGINT PRIMARY KEY COMMENT '商品ID',
    stock INT NOT NULL DEFAULT 0 COMMENT '当前库存',
    reserved_stock INT NOT NULL DEFAULT 0 COMMENT '预占库存',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品库存表';

-- 2. 订单表
CREATE TABLE `order` (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '订单ID',
    order_no VARCHAR(64) NOT NULL COMMENT '订单号',
    product_id BIGINT NOT NULL COMMENT '商品ID',
    quantity INT NOT NULL COMMENT '购买数量',
    status VARCHAR(20) NOT NULL COMMENT '订单状态：PENDING-待支付，PAID-已支付，TIMEOUT-超时，CANCELLED-已取消',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_order_no (order_no),
    INDEX idx_product (product_id),
    INDEX idx_status_created (status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 3. 库存操作日志表（用于监控和对账）
CREATE TABLE stock_operation_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    product_id BIGINT NOT NULL COMMENT '商品ID',
    operation_type VARCHAR(20) NOT NULL COMMENT '操作类型：DEDUCT-扣减，CONFIRM-确认，RELEASE-释放',
    quantity INT NOT NULL COMMENT '操作数量',
    before_stock INT NOT NULL COMMENT '操作前库存',
    after_stock INT NOT NULL COMMENT '操作后库存',
    order_id BIGINT COMMENT '关联订单ID',
    result VARCHAR(20) NOT NULL COMMENT '操作结果：SUCCESS-成功，FAIL-失败',
    fail_reason VARCHAR(255) COMMENT '失败原因',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_product_created (product_id, created_at),
    INDEX idx_order (order_id),
    INDEX idx_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存操作日志表';

-- 初始化测试数据
INSERT INTO product_stock (product_id, stock, reserved_stock) VALUES
(1001, 1000, 0),
(1002, 500, 0),
(1003, 100, 0);

-- 查询库存视图（实际可用库存）
CREATE VIEW v_available_stock AS
SELECT
    product_id,
    stock,
    reserved_stock,
    stock - reserved_stock as available_stock,
    updated_at
FROM product_stock;

-- 常用查询脚本

-- 查询超时未支付订单
SELECT
    o.id,
    o.order_no,
    o.product_id,
    o.quantity,
    o.created_at,
    TIMESTAMPDIFF(MINUTE, o.created_at, NOW()) as timeout_minutes
FROM `order` o
WHERE o.status = 'PENDING'
AND o.created_at < DATE_SUB(NOW(), INTERVAL 30 MINUTE)
ORDER BY o.created_at;

-- 查询商品库存情况
SELECT
    ps.product_id,
    ps.stock as total_stock,
    ps.reserved_stock,
    ps.stock - ps.reserved_stock as available_stock,
    (ps.reserved_stock * 100.0 / NULLIF(ps.stock, 0)) as reserved_rate,
    ps.updated_at
FROM product_stock ps
ORDER BY reserved_rate DESC;

-- 查询最近 1 小时库存操作统计
SELECT
    product_id,
    operation_type,
    COUNT(*) as operation_count,
    SUM(quantity) as total_quantity,
    SUM(CASE WHEN result = 'SUCCESS' THEN 1 ELSE 0 END) as success_count,
    SUM(CASE WHEN result = 'FAIL' THEN 1 ELSE 0 END) as fail_count,
    SUM(CASE WHEN result = 'SUCCESS' THEN 1 ELSE 0 END) * 100.0 / COUNT(*) as success_rate
FROM stock_operation_log
WHERE created_at > DATE_SUB(NOW(), INTERVAL 1 HOUR)
GROUP BY product_id, operation_type
ORDER BY product_id, operation_type;

-- 数据修复脚本（紧急情况使用）

-- 修复超时订单的预占库存
UPDATE product_stock ps
INNER JOIN (
    SELECT
        product_id,
        SUM(quantity) as total_quantity
    FROM `order`
    WHERE status = 'PENDING'
    AND created_at < DATE_SUB(NOW(), INTERVAL 30 MINUTE)
    GROUP BY product_id
) timeout_orders ON ps.product_id = timeout_orders.product_id
SET
    ps.stock = ps.stock + timeout_orders.total_quantity,
    ps.reserved_stock = ps.reserved_stock - timeout_orders.total_quantity,
    ps.updated_at = NOW();

-- 更新超时订单状态
UPDATE `order`
SET status = 'TIMEOUT',
    updated_at = NOW()
WHERE status = 'PENDING'
AND created_at < DATE_SUB(NOW(), INTERVAL 30 MINUTE);

-- 性能优化脚本

-- 分析表
ANALYZE TABLE product_stock;
ANALYZE TABLE `order`;
ANALYZE TABLE stock_operation_log;

-- 查看索引使用情况
EXPLAIN SELECT * FROM `order` WHERE status = 'PENDING' AND created_at < DATE_SUB(NOW(), INTERVAL 30 MINUTE);
EXPLAIN SELECT * FROM product_stock WHERE product_id = 1001 FOR UPDATE;

-- 清理历史日志（保留最近 30 天）
DELETE FROM stock_operation_log
WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY)
LIMIT 10000;
