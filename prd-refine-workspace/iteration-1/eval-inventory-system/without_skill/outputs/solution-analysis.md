# 库存扣减并发问题解决方案分析

## 评审报告问题总结

评审报告中提出了以下核心问题：

### P0 级别（阻断问题）
1. **高并发下的库存超卖风险** - 未说明并发控制方案
2. **预占库存的超时释放机制缺失** - 超时订单如何处理

### P1 级别（需确认）
3. **1000 QPS 性能目标可行性** - 直接操作 MySQL 可能成为瓶颈
4. **分布式场景下的库存一致性** - 多实例部署的一致性保证

---

## 当前技术栈

- **数据库**: MySQL (Spring Boot + JPA)
- **缓存**: 无（未使用 Redis）
- **依赖**: Spring Boot 2.7.0, MySQL Connector, Lombok

---

## 解决方案对比分析

### 方案一：纯 MySQL 方案（悲观锁 + 行锁）

#### 实现方式

```java
/**
 * 库存扣减（使用悲观锁）
 *
 * @param productId 商品ID
 * @param quantity 扣减数量
 * @return true-扣减成功 false-库存不足
 * @author wenlong.chen
 * @date 2026/07/01 12:00:00
 */
@Transactional(isolation = Isolation.READ_COMMITTED)
public boolean deductStock(Long productId, Integer quantity) {
    // 使用 SELECT ... FOR UPDATE 加行锁
    ProductStock stock = stockRepository.findByIdForUpdate(productId);
    
    // 检查库存是否足够
    if (stock.getStock() < quantity) {
        return false;
    }
    
    // 扣减库存
    stock.setStock(stock.getStock() - quantity);
    stock.setReservedStock(stock.getReservedStock() + quantity);
    stock.setUpdatedAt(LocalDateTime.now());
    
    stockRepository.save(stock);
    return true;
}

// Repository 层
@Query("SELECT s FROM ProductStock s WHERE s.productId = :productId FOR UPDATE")
ProductStock findByIdForUpdate(@Param("productId") Long productId);
```

#### 优点
✅ **实现简单** - 无需额外中间件，直接使用 MySQL 行锁
✅ **数据强一致性** - 数据库事务保证 ACID
✅ **运维成本低** - 不需要维护 Redis 等额外组件
✅ **适合中小流量** - 对于 100-500 QPS 完全够用

#### 缺点
❌ **性能瓶颈明显** - 1000 QPS 可能吃力，热门商品会成为瓶颈
❌ **数据库压力大** - 所有请求都打到数据库
❌ **锁等待时间长** - 并发高时大量请求排队等锁
❌ **扩展性差** - 无法通过增加应用实例提升性能

#### 适用场景
- 订单量较小的电商系统（日订单量 < 10 万）
- 非秒杀场景
- 预算有限、不想增加技术复杂度
- 团队对 Redis 不熟悉

---

### 方案二：纯 MySQL 方案（乐观锁 + 版本号）

#### 实现方式

```java
/**
 * 商品库存实体
 */
@Entity
@Table(name = "product_stock")
@Data
public class ProductStock {
    @Id
    private Long productId;
    
    private Integer stock;
    
    private Integer reservedStock;
    
    @Version
    private Long version;  // 乐观锁版本号
    
    private LocalDateTime updatedAt;
}

/**
 * 库存扣减（使用乐观锁）
 *
 * @param productId 商品ID
 * @param quantity 扣减数量
 * @return true-扣减成功 false-库存不足或并发冲突
 * @author wenlong.chen
 * @date 2026/07/01 12:00:00
 */
@Transactional
public boolean deductStockWithRetry(Long productId, Integer quantity) {
    int maxRetry = 3;
    int retryCount = 0;
    
    while (retryCount < maxRetry) {
        try {
            ProductStock stock = stockRepository.findById(productId).orElseThrow();
            
            // 检查库存
            if (stock.getStock() < quantity) {
                return false;
            }
            
            // 扣减库存
            stock.setStock(stock.getStock() - quantity);
            stock.setReservedStock(stock.getReservedStock() + quantity);
            stock.setUpdatedAt(LocalDateTime.now());
            
            // JPA 自动使用 version 字段进行乐观锁控制
            stockRepository.save(stock);
            return true;
            
        } catch (OptimisticLockException e) {
            // 版本冲突，重试
            retryCount++;
            if (retryCount >= maxRetry) {
                throw new BusinessException("库存扣减失败，请重试");
            }
        }
    }
    return false;
}
```

#### 优点
✅ **无锁等待** - 不会阻塞其他请求
✅ **吞吐量较高** - 比悲观锁性能好
✅ **实现简单** - JPA 内置 @Version 支持
✅ **适合低冲突场景** - 并发不高时效率很好

#### 缺点
❌ **高并发下重试频繁** - 冲突多时大量请求失败重试
❌ **用户体验差** - 重试失败会返回"请重试"错误
❌ **不适合秒杀** - 热门商品会导致大量失败
❌ **数据库压力仍然大** - 依然所有请求打到数据库

#### 适用场景
- 普通商品销售场景（非秒杀）
- 并发冲突率低（< 10%）
- 不想用悲观锁阻塞请求

---

### 方案三：MySQL + Redis 缓存方案

#### 实现方式

```java
/**
 * 库存扣减（Redis 缓存 + MySQL 持久化）
 *
 * @param productId 商品ID
 * @param quantity 扣减数量
 * @return true-扣减成功 false-库存不足
 * @author wenlong.chen
 * @date 2026/07/01 12:00:00
 */
public boolean deductStockWithRedis(Long productId, Integer quantity) {
    String stockKey = "stock:" + productId;
    
    // Redis Lua 脚本保证原子性
    String luaScript = 
        "local stock = redis.call('get', KEYS[1]) " +
        "if not stock or tonumber(stock) < tonumber(ARGV[1]) then " +
        "    return 0 " +
        "end " +
        "redis.call('decrby', KEYS[1], ARGV[1]) " +
        "return 1";
    
    // 执行 Redis 扣减
    Long result = redisTemplate.execute(
        new DefaultRedisScript<>(luaScript, Long.class),
        Collections.singletonList(stockKey),
        quantity.toString()
    );
    
    if (result == 0) {
        return false;
    }
    
    // 异步持久化到 MySQL
    asyncUpdateDatabase(productId, quantity);
    return true;
}

/**
 * 异步更新数据库
 */
@Async
private void asyncUpdateDatabase(Long productId, Integer quantity) {
    // 批量更新或者使用消息队列
    stockRepository.updateStock(productId, quantity);
}
```

#### 优点
✅ **性能极高** - Redis 内存操作，轻松支持 10000+ QPS
✅ **响应快** - 毫秒级响应
✅ **热点商品友好** - 秒杀场景表现优秀
✅ **可扩展** - 通过 Redis 集群进一步提升性能

#### 缺点
❌ **一致性问题** - Redis 和 MySQL 数据可能不一致
❌ **Redis 故障风险** - Redis 宕机后库存数据丢失
❌ **运维复杂度高** - 需要维护 Redis，增加部署成本
❌ **数据同步复杂** - 需要处理 Redis 和 MySQL 的数据同步
❌ **开发复杂度高** - 需要处理缓存穿透、击穿、雪崩

#### 适用场景
- 高并发秒杀场景
- 1000+ QPS 的性能要求
- 有运维 Redis 的能力
- 可接受最终一致性

---

### 方案四：MySQL + Redis 分布式锁方案

#### 实现方式

```java
/**
 * 库存扣减（Redis 分布式锁 + MySQL）
 *
 * @param productId 商品ID
 * @param quantity 扣减数量
 * @return true-扣减成功 false-库存不足
 * @author wenlong.chen
 * @date 2026/07/01 12:00:00
 */
public boolean deductStockWithDistributedLock(Long productId, Integer quantity) {
    String lockKey = "lock:stock:" + productId;
    String requestId = UUID.randomUUID().toString();
    
    try {
        // 获取分布式锁
        boolean locked = redisTemplate.opsForValue().setIfAbsent(
            lockKey, 
            requestId, 
            10, 
            TimeUnit.SECONDS
        );
        
        if (!locked) {
            return false;
        }
        
        // 查询并扣减库存
        ProductStock stock = stockRepository.findById(productId).orElseThrow();
        
        if (stock.getStock() < quantity) {
            return false;
        }
        
        stock.setStock(stock.getStock() - quantity);
        stock.setReservedStock(stock.getReservedStock() + quantity);
        stockRepository.save(stock);
        
        return true;
        
    } finally {
        // 释放锁（使用 Lua 脚本保证原子性）
        releaseLock(lockKey, requestId);
    }
}
```

#### 优点
✅ **分布式友好** - 支持多实例部署
✅ **强一致性** - 数据库保证数据准确
✅ **实现相对简单** - 基于 Redis 的 SETNX
✅ **性能适中** - 比纯数据库锁好，比纯缓存差

#### 缺点
❌ **锁竞争激烈** - 热门商品依然有瓶颈
❌ **Redis 单点问题** - Redis 故障影响业务
❌ **需要 Redis** - 增加技术栈和运维成本
❌ **锁超时问题** - 需要处理锁续期

#### 适用场景
- 分布式部署场景
- 500-1000 QPS
- 需要强一致性
- 已有 Redis 基础设施

---

## 针对 PRD 评审问题的解决方案

### 问题 1: 高并发下的库存超卖风险

**推荐方案**: 
- **当前阶段（无 Redis）**: 方案一（MySQL 悲观锁）
- **性能不足时**: 升级到方案四（Redis 分布式锁）或方案三（Redis 缓存）

**具体实现**:
```sql
-- MySQL 层面使用行锁
UPDATE product_stock 
SET stock = stock - ?, 
    reserved_stock = reserved_stock + ?,
    updated_at = NOW()
WHERE product_id = ? 
  AND stock >= ?;  -- 防止超卖
```

---

### 问题 2: 预占库存的超时释放机制

**解决方案**: 定时任务 + 消息队列（延迟消息）

#### 方案 A: Spring 定时任务（简单方案）

```java
/**
 * 超时订单库存释放定时任务
 *
 * @author wenlong.chen
 * @date 2026/07/01 12:00:00
 */
@Component
@Slf4j
public class StockReleaseTask {

    /**
     * 订单超时时间（分钟）
     */
    private static final int ORDER_TIMEOUT_MINUTES = 30;

    /**
     * 订单仓储
     */
    private final OrderRepository orderRepository;

    /**
     * 库存服务
     */
    private final StockService stockService;

    /**
     * 每 5 分钟扫描一次超时订单
     *
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Scheduled(cron = "0 */5 * * * ?")
    public void releaseExpiredStock() {
        LocalDateTime expireTime = LocalDateTime.now()
            .minusMinutes(ORDER_TIMEOUT_MINUTES);
        
        // 查询超时未支付订单
        List<Order> expiredOrders = orderRepository
            .findByStatusAndCreatedTimeBefore("PENDING", expireTime);
        
        for (Order order : expiredOrders) {
            try {
                // 释放库存
                stockService.releaseReservedStock(
                    order.getProductId(), 
                    order.getQuantity()
                );
                
                // 更新订单状态
                order.setStatus("CANCELLED");
                orderRepository.save(order);
                
                log.info("释放超时订单库存: orderId={}, productId={}", 
                    order.getId(), order.getProductId());
                    
            } catch (Exception e) {
                log.error("释放库存失败: orderId={}", order.getId(), e);
            }
        }
    }
}
```

**优点**: 实现简单，无需额外组件
**缺点**: 
- 定时任务故障会导致库存永久锁死
- 扫描全表性能差
- 时效性差（最多延迟 5 分钟）

#### 方案 B: 延迟消息队列（推荐）

```java
/**
 * 下单时发送延迟消息
 *
 * @param orderId 订单ID
 * @author wenlong.chen
 * @date 2026/07/01 12:00:00
 */
public void createOrder(Long orderId) {
    // ... 创建订单逻辑
    
    // 发送 30 分钟后的延迟消息
    rabbitTemplate.convertAndSend(
        "order.delay.exchange",
        "order.timeout.check",
        orderId,
        message -> {
            message.getMessageProperties().setDelay(30 * 60 * 1000);
            return message;
        }
    );
}

/**
 * 消费延迟消息，检查订单状态
 *
 * @param orderId 订单ID
 * @author wenlong.chen
 * @date 2026/07/01 12:00:00
 */
@RabbitListener(queues = "order.timeout.queue")
public void checkOrderTimeout(Long orderId) {
    Order order = orderRepository.findById(orderId).orElse(null);
    
    if (order != null && "PENDING".equals(order.getStatus())) {
        // 订单仍未支付，释放库存
        stockService.releaseReservedStock(
            order.getProductId(), 
            order.getQuantity()
        );
        
        order.setStatus("CANCELLED");
        orderRepository.save(order);
    }
}
```

**优点**: 
- 时效性准确
- 不需要扫描全表
- 消息可靠性高

**缺点**: 需要引入 RabbitMQ 或 RocketMQ

---

### 问题 3: 1000 QPS 性能目标可行性

**分析**:

| 方案 | 预估 QPS | 是否满足 1000 QPS |
|------|---------|------------------|
| 方案一（MySQL 悲观锁） | 300-500 | ❌ 不满足 |
| 方案二（MySQL 乐观锁） | 500-800 | ⚠️ 勉强满足 |
| 方案三（Redis 缓存） | 10000+ | ✅ 远超目标 |
| 方案四（Redis 分布式锁） | 800-1200 | ✅ 满足 |

**推荐**:
- 如果 1000 QPS 是硬性要求 → **必须引入 Redis**（方案三或方案四）
- 如果实际流量 < 500 QPS → 方案一足够
- 如果可以接受最终一致性 → 方案三性能最佳

---

### 问题 4: 分布式场景下的库存一致性

**问题根源**: 多个应用实例同时操作同一条库存记录

**解决方案对比**:

| 方案 | 分布式一致性 | 实现复杂度 |
|------|------------|-----------|
| MySQL 悲观锁 | ✅ 强一致 | 简单 |
| MySQL 乐观锁 | ✅ 强一致 | 简单 |
| Redis 分布式锁 | ✅ 强一致 | 中等 |
| Redis 缓存 | ⚠️ 最终一致 | 复杂 |

**推荐**: 
- 方案一或方案四都能保证分布式一致性
- 数据库行锁天然支持分布式场景

---

## 综合建议

### 场景一：预算有限、团队规模小、流量不大

**推荐**: **方案一（MySQL 悲观锁）**

**理由**:
- 实现简单，无需额外中间件
- 500 QPS 以内完全够用
- 运维成本低
- 数据强一致性

**实施步骤**:
1. 使用 `SELECT ... FOR UPDATE` 加行锁
2. 在事务中完成库存扣减
3. 使用 Spring `@Scheduled` 定时任务释放超时库存
4. 数据库连接池调优（HikariCP）

### 场景二：需要支持 1000 QPS，有 Redis 运维能力

**推荐**: **方案四（Redis 分布式锁 + MySQL）**

**理由**:
- 保证强一致性
- 性能满足 1000 QPS
- 支持分布式部署
- 实现复杂度适中

**实施步骤**:
1. 引入 Redis（单机或哨兵模式）
2. 使用 Redisson 实现分布式锁
3. 锁内操作 MySQL 数据库
4. 使用 RabbitMQ 延迟消息处理超时订单

### 场景三：秒杀场景，需要极高性能

**推荐**: **方案三（Redis 缓存 + 异步持久化）**

**理由**:
- 性能极高（10000+ QPS）
- 适合秒杀等瞬时高并发
- 用户体验好

**实施步骤**:
1. 预热库存到 Redis
2. 使用 Lua 脚本保证原子性
3. 异步批量同步到 MySQL
4. 使用 Canal 监听 Binlog 保证数据一致性

---

## 推荐技术方案（针对当前项目）

基于当前项目情况（无 Redis，Spring Boot + MySQL），建议采用**分阶段演进**策略：

### 第一阶段（立即实施）：MySQL 悲观锁方案

```xml
<!-- 无需新增依赖 -->
```

**实施清单**:
- [ ] 实现 `SELECT ... FOR UPDATE` 查询
- [ ] 添加库存扣减事务控制
- [ ] 实现定时任务扫描超时订单
- [ ] 添加数据库索引优化
- [ ] 配置数据库连接池参数

**预期效果**: 支持 300-500 QPS

### 第二阶段（性能不足时）：引入 Redis 分布式锁

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-data-redis</artifactId>
    <version>2.7.0</version>
</dependency>
<dependency>
    <groupId>org.redisson</groupId>
    <artifactId>redisson-spring-boot-starter</artifactId>
    <version>3.17.0</version>
</dependency>
```

**实施清单**:
- [ ] 部署 Redis（单机或哨兵）
- [ ] 集成 Redisson 分布式锁
- [ ] 改造库存扣减逻辑
- [ ] 添加 Redis 监控
- [ ] 制定 Redis 故障降级方案

**预期效果**: 支持 800-1200 QPS

### 第三阶段（秒杀场景）：Redis 缓存方案

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-amqp</artifactId>
    <version>2.7.0</version>
</dependency>
```

**实施清单**:
- [ ] 库存数据预热到 Redis
- [ ] 实现 Lua 脚本原子扣减
- [ ] 部署 RabbitMQ 或 RocketMQ
- [ ] 实现异步持久化逻辑
- [ ] 添加 Canal 监听 Binlog
- [ ] 实现数据对账机制

**预期效果**: 支持 10000+ QPS

---

## 关键代码示例

### 完整的库存服务实现（方案一：MySQL 悲观锁）
