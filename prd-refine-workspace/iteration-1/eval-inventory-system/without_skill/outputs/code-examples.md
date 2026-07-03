# 代码实现示例

## 方案一：MySQL 悲观锁完整实现

### 1. 实体类

```java
/**
 * @Description: 商品库存实体
 * @Author: wenlong.chen
 * @Date: 2026/07/01 12:00:00
 */
@Entity
@Table(name = "product_stock")
@Data
public class ProductStock {

    /**
     * 商品ID
     */
    @Id
    private Long productId;

    /**
     * 当前库存
     */
    private Integer stock;

    /**
     * 预占库存
     */
    private Integer reservedStock;

    /**
     * 更新时间
     */
    private LocalDateTime updatedAt;
}
```

### 2. Repository 层

```java
/**
 * @Description: 商品库存仓储
 * @Author: wenlong.chen
 * @Date: 2026/07/01 12:00:00
 */
public interface ProductStockRepository extends JpaRepository<ProductStock, Long> {

    /**
     * 查询库存并加行锁
     *
     * @param productId 商品ID
     * @return {@link ProductStock} 商品库存
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Query("SELECT s FROM ProductStock s WHERE s.productId = :productId FOR UPDATE")
    ProductStock findByIdForUpdate(@Param("productId") Long productId);
}
```

### 3. Service 层

```java
/**
 * @Description: 库存服务实现类
 * @Author: wenlong.chen
 * @Date: 2026/07/01 12:00:00
 */
@Service
@RequiredArgsConstructor
@Slf4j
public class StockServiceImpl implements StockService {

    /**
     * 商品库存仓储
     */
    private final ProductStockRepository stockRepository;

    /**
     * 扣减库存（悲观锁）
     *
     * @param productId 商品ID
     * @param quantity 扣减数量
     * @return true-扣减成功 false-库存不足
     * @throws StockException 库存操作异常
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Override
    @Transactional(rollbackFor = Exception.class)
    public boolean deductStock(Long productId, Integer quantity) {
        Objects.requireNonNull(productId, "商品ID不能为空");
        Objects.requireNonNull(quantity, "扣减数量不能为空");

        // 加行锁查询
        ProductStock stock = stockRepository.findByIdForUpdate(productId);
        if (stock == null) {
            throw new StockException("商品库存不存在");
        }

        // 检查库存是否足够
        if (stock.getStock() < quantity) {
            log.warn("库存不足: productId={}, 当前库存={}, 需要={}", 
                productId, stock.getStock(), quantity);
            return false;
        }

        // 扣减库存，增加预占
        stock.setStock(stock.getStock() - quantity);
        stock.setReservedStock(stock.getReservedStock() + quantity);
        stock.setUpdatedAt(LocalDateTime.now());

        stockRepository.save(stock);
        log.info("库存扣减成功: productId={}, quantity={}", productId, quantity);
        return true;
    }

    /**
     * 确认扣减（支付成功）
     *
     * @param productId 商品ID
     * @param quantity 确认数量
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Override
    @Transactional(rollbackFor = Exception.class)
    public void confirmDeduct(Long productId, Integer quantity) {
        ProductStock stock = stockRepository.findByIdForUpdate(productId);
        if (stock == null) {
            throw new StockException("商品库存不存在");
        }

        // 减少预占库存
        stock.setReservedStock(stock.getReservedStock() - quantity);
        stock.setUpdatedAt(LocalDateTime.now());

        stockRepository.save(stock);
        log.info("确认扣减成功: productId={}, quantity={}", productId, quantity);
    }

    /**
     * 释放预占库存（订单取消）
     *
     * @param productId 商品ID
     * @param quantity 释放数量
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Override
    @Transactional(rollbackFor = Exception.class)
    public void releaseStock(Long productId, Integer quantity) {
        ProductStock stock = stockRepository.findByIdForUpdate(productId);
        if (stock == null) {
            throw new StockException("商品库存不存在");
        }

        // 释放预占，回补库存
        stock.setStock(stock.getStock() + quantity);
        stock.setReservedStock(stock.getReservedStock() - quantity);
        stock.setUpdatedAt(LocalDateTime.now());

        stockRepository.save(stock);
        log.info("释放库存成功: productId={}, quantity={}", productId, quantity);
    }
}
```

### 4. 定时任务（处理超时订单）

```java
/**
 * @Description: 订单超时处理定时任务
 * @Author: wenlong.chen
 * @Date: 2026/07/01 12:00:00
 */
@Component
@RequiredArgsConstructor
@Slf4j
public class OrderTimeoutTask {

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
    public void processTimeoutOrders() {
        LocalDateTime expireTime = LocalDateTime.now().minusMinutes(ORDER_TIMEOUT_MINUTES);

        // 查询超时未支付订单
        List<Order> timeoutOrders = orderRepository
            .findByStatusAndCreatedTimeBefore("PENDING", expireTime);

        log.info("开始处理超时订单，数量: {}", timeoutOrders.size());

        for (Order order : timeoutOrders) {
            try {
                // 释放预占库存
                stockService.releaseStock(order.getProductId(), order.getQuantity());

                // 更新订单状态
                order.setStatus("TIMEOUT");
                order.setUpdatedAt(LocalDateTime.now());
                orderRepository.save(order);

                log.info("超时订单处理成功: orderId={}", order.getId());

            } catch (Exception e) {
                log.error("超时订单处理失败: orderId={}", order.getId(), e);
            }
        }
    }
}
```

### 5. 数据库配置优化

```properties
# application.properties

# HikariCP 连接池配置
spring.datasource.hikari.maximum-pool-size=50
spring.datasource.hikari.minimum-idle=10
spring.datasource.hikari.connection-timeout=30000
spring.datasource.hikari.idle-timeout=600000
spring.datasource.hikari.max-lifetime=1800000

# JPA 配置
spring.jpa.show-sql=false
spring.jpa.properties.hibernate.format_sql=false
spring.jpa.properties.hibernate.jdbc.batch_size=20
```

### 6. 数据库表结构

```sql
CREATE TABLE product_stock (
    product_id BIGINT PRIMARY KEY COMMENT '商品ID',
    stock INT NOT NULL DEFAULT 0 COMMENT '当前库存',
    reserved_stock INT NOT NULL DEFAULT 0 COMMENT '预占库存',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品库存表';

CREATE TABLE `order` (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '订单ID',
    product_id BIGINT NOT NULL COMMENT '商品ID',
    quantity INT NOT NULL COMMENT '购买数量',
    status VARCHAR(20) NOT NULL COMMENT '订单状态：PENDING-待支付，PAID-已支付，TIMEOUT-超时，CANCELLED-已取消',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status_created (status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';
```

---

## 方案四：Redis 分布式锁实现

### 1. 添加依赖

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

### 2. Redisson 配置

```java
/**
 * @Description: Redisson 配置类
 * @Author: wenlong.chen
 * @Date: 2026/07/01 12:00:00
 */
@Configuration
public class RedissonConfig {

    /**
     * Redisson 客户端
     *
     * @return {@link RedissonClient} Redisson 客户端实例
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Bean
    public RedissonClient redissonClient() {
        Config config = new Config();
        config.useSingleServer()
            .setAddress("redis://127.0.0.1:6379")
            .setDatabase(0)
            .setConnectionPoolSize(50)
            .setConnectionMinimumIdleSize(10);
        return Redisson.create(config);
    }
}
```

### 3. 分布式锁 Service

```java
/**
 * @Description: 库存服务实现类（分布式锁版本）
 * @Author: wenlong.chen
 * @Date: 2026/07/01 12:00:00
 */
@Service
@RequiredArgsConstructor
@Slf4j
public class StockServiceWithLock implements StockService {

    /**
     * 商品库存仓储
     */
    private final ProductStockRepository stockRepository;

    /**
     * Redisson 客户端
     */
    private final RedissonClient redissonClient;

    /**
     * 锁等待时间（秒）
     */
    private static final long LOCK_WAIT_TIME = 3;

    /**
     * 锁持有时间（秒）
     */
    private static final long LOCK_LEASE_TIME = 10;

    /**
     * 扣减库存（使用分布式锁）
     *
     * @param productId 商品ID
     * @param quantity 扣减数量
     * @return true-扣减成功 false-库存不足或获取锁失败
     * @throws StockException 库存操作异常
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Override
    @Transactional(rollbackFor = Exception.class)
    public boolean deductStock(Long productId, Integer quantity) {
        Objects.requireNonNull(productId, "商品ID不能为空");
        Objects.requireNonNull(quantity, "扣减数量不能为空");

        String lockKey = "lock:stock:" + productId;
        RLock lock = redissonClient.getLock(lockKey);

        try {
            // 尝试获取锁
            boolean locked = lock.tryLock(LOCK_WAIT_TIME, LOCK_LEASE_TIME, TimeUnit.SECONDS);
            if (!locked) {
                log.warn("获取锁失败: productId={}", productId);
                return false;
            }

            // 查询库存
            ProductStock stock = stockRepository.findById(productId)
                .orElseThrow(() -> new StockException("商品库存不存在"));

            // 检查库存
            if (stock.getStock() < quantity) {
                log.warn("库存不足: productId={}, 当前库存={}, 需要={}", 
                    productId, stock.getStock(), quantity);
                return false;
            }

            // 扣减库存
            stock.setStock(stock.getStock() - quantity);
            stock.setReservedStock(stock.getReservedStock() + quantity);
            stock.setUpdatedAt(LocalDateTime.now());

            stockRepository.save(stock);
            log.info("库存扣减成功: productId={}, quantity={}", productId, quantity);
            return true;

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            log.error("获取锁被中断: productId={}", productId, e);
            return false;
        } finally {
            // 释放锁
            if (lock.isHeldByCurrentThread()) {
                lock.unlock();
            }
        }
    }
}
```

---

## 方案三：Redis 缓存实现

### 1. Lua 脚本配置

```java
/**
 * @Description: Redis Lua 脚本配置
 * @Author: wenlong.chen
 * @Date: 2026/07/01 12:00:00
 */
@Configuration
public class RedisLuaScriptConfig {

    /**
     * 库存扣减 Lua 脚本
     *
     * @return {@link DefaultRedisScript} Lua 脚本对象
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Bean
    public DefaultRedisScript<Long> stockDeductScript() {
        DefaultRedisScript<Long> script = new DefaultRedisScript<>();
        script.setScriptText(
            "local stock = redis.call('get', KEYS[1]) " +
            "if not stock or tonumber(stock) < tonumber(ARGV[1]) then " +
            "    return 0 " +
            "end " +
            "redis.call('decrby', KEYS[1], ARGV[1]) " +
            "return 1"
        );
        script.setResultType(Long.class);
        return script;
    }
}
```

### 2. Redis 缓存 Service

```java
/**
 * @Description: 库存服务实现类（Redis 缓存版本）
 * @Author: wenlong.chen
 * @Date: 2026/07/01 12:00:00
 */
@Service
@RequiredArgsConstructor
@Slf4j
public class StockServiceWithRedis implements StockService {

    /**
     * Redis 模板
     */
    private final StringRedisTemplate redisTemplate;

    /**
     * 库存扣减 Lua 脚本
     */
    private final DefaultRedisScript<Long> stockDeductScript;

    /**
     * 库存仓储
     */
    private final ProductStockRepository stockRepository;

    /**
     * 扣减库存（使用 Redis 缓存）
     *
     * @param productId 商品ID
     * @param quantity 扣减数量
     * @return true-扣减成功 false-库存不足
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Override
    public boolean deductStock(Long productId, Integer quantity) {
        String stockKey = "stock:" + productId;

        // 执行 Lua 脚本扣减 Redis 库存
        Long result = redisTemplate.execute(
            stockDeductScript,
            Collections.singletonList(stockKey),
            quantity.toString()
        );

        if (result == null || result == 0) {
            log.warn("Redis 库存不足: productId={}, quantity={}", productId, quantity);
            return false;
        }

        // 异步更新数据库
        asyncUpdateDatabase(productId, quantity);

        log.info("Redis 库存扣减成功: productId={}, quantity={}", productId, quantity);
        return true;
    }

    /**
     * 异步更新数据库
     *
     * @param productId 商品ID
     * @param quantity 扣减数量
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @Async
    public void asyncUpdateDatabase(Long productId, Integer quantity) {
        try {
            ProductStock stock = stockRepository.findById(productId)
                .orElseThrow(() -> new StockException("商品库存不存在"));

            stock.setStock(stock.getStock() - quantity);
            stock.setReservedStock(stock.getReservedStock() + quantity);
            stock.setUpdatedAt(LocalDateTime.now());

            stockRepository.save(stock);
            log.info("数据库库存更新成功: productId={}, quantity={}", productId, quantity);

        } catch (Exception e) {
            log.error("数据库库存更新失败: productId={}, quantity={}", productId, quantity, e);
            // 补偿机制：Redis 回滚
            rollbackRedisStock(productId, quantity);
        }
    }

    /**
     * Redis 库存回滚
     *
     * @param productId 商品ID
     * @param quantity 回滚数量
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    private void rollbackRedisStock(Long productId, Integer quantity) {
        String stockKey = "stock:" + productId;
        redisTemplate.opsForValue().increment(stockKey, quantity);
        log.info("Redis 库存回滚: productId={}, quantity={}", productId, quantity);
    }

    /**
     * 预热库存到 Redis
     *
     * @param productId 商品ID
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    public void preloadStock(Long productId) {
        ProductStock stock = stockRepository.findById(productId)
            .orElseThrow(() -> new StockException("商品库存不存在"));

        String stockKey = "stock:" + productId;
        redisTemplate.opsForValue().set(stockKey, stock.getStock().toString());

        log.info("库存预热成功: productId={}, stock={}", productId, stock.getStock());
    }
}
```

---

## 性能优化建议

### 1. 数据库索引优化

```sql
-- 订单状态和创建时间联合索引（用于超时订单扫描）
CREATE INDEX idx_status_created ON `order`(status, created_at);

-- 库存更新时间索引
CREATE INDEX idx_updated_at ON product_stock(updated_at);
```

### 2. 连接池调优

```properties
# HikariCP 配置
spring.datasource.hikari.maximum-pool-size=50
spring.datasource.hikari.minimum-idle=10
spring.datasource.hikari.connection-timeout=30000
```

### 3. 批量操作优化

```java
/**
 * 批量释放库存
 *
 * @param releaseList 释放列表
 * @author wenlong.chen
 * @date 2026/07/01 12:00:00
 */
@Transactional(rollbackFor = Exception.class)
public void batchReleaseStock(List<StockReleaseDTO> releaseList) {
    for (StockReleaseDTO release : releaseList) {
        ProductStock stock = stockRepository.findById(release.getProductId())
            .orElse(null);
        if (stock != null) {
            stock.setStock(stock.getStock() + release.getQuantity());
            stock.setReservedStock(stock.getReservedStock() - release.getQuantity());
        }
    }
    // 批量保存
    stockRepository.flush();
}
```
