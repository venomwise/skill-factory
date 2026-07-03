# 支付网关 PRD 评审问题分析与解决方案（续）

## P1-需确认问题

### 问题 3：回调超时处理策略不明确

#### 问题详细分析

**风险等级:** ⚠️ 高（影响用户体验和订单准确性）

**问题根源:**
支付回调可能因为以下原因失败或超时：
1. 网络故障（支付网关 → 我们的服务器）
2. 我们的服务器宕机或重启
3. 我们的服务处理超时（业务逻辑复杂，超过支付网关的超时时间）
4. 支付网关自身异常

**影响:**
- 用户已支付，但订单状态未更新，用户认为支付失败
- 用户可能重复支付
- 客服工作量增加（用户投诉"我付了钱怎么没到账"）

**数据:**
根据行业经验，支付回调失败率约 0.1%-0.5%，对于日均 10 万笔交易的系统，每天会有 100-500 笔回调失败。

#### 解决方案

##### 方案 1：主动查询 + 指数退避重试（推荐）

**核心思路:**
当回调超时后，主动调用支付网关的查询接口，使用指数退避算法进行重试。

**实现步骤:**

1. **设计补偿任务表**

```sql
CREATE TABLE payment_compensation_task (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id VARCHAR(64) NOT NULL,
    payment_no VARCHAR(64) NOT NULL,
    retry_count INT DEFAULT 0 COMMENT '重试次数',
    max_retry INT DEFAULT 10 COMMENT '最大重试次数',
    next_retry_time DATETIME COMMENT '下次重试时间',
    status TINYINT DEFAULT 0 COMMENT '0-待处理 1-成功 2-失败',
    error_msg TEXT COMMENT '错误信息',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_next_retry_time(next_retry_time, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

2. **创建补偿任务（在发起支付时）**

```java
@Service
@RequiredArgsConstructor
public class PaymentService {
    
    private final PaymentGatewayClient gatewayClient;
    private final CompensationTaskMapper compensationTaskMapper;
    
    @Transactional(rollbackFor = Exception.class)
    public PaymentVO createPayment(OrderDTO orderDTO) {
        // 调用支付网关
        PaymentVO payment = gatewayClient.createPayment(orderDTO);
        
        // 创建补偿任务（5 分钟后开始第一次查询）
        CompensationTask task = new CompensationTask();
        task.setOrderId(orderDTO.getOrderId());
        task.setPaymentNo(payment.getPaymentNo());
        task.setNextRetryTime(LocalDateTime.now().plusMinutes(5));
        task.setMaxRetry(10);
        
        compensationTaskMapper.insert(task);
        
        return payment;
    }
}
```

3. **定时任务执行补偿查询（完整实现见完整代码文件）**

**重试时间策略:**

| 重试次数 | 延迟时间 | 累计时间 | 说明 |
|---------|---------|---------|------|
| 1 | 2 分钟 | 2 分钟 | 快速重试，处理网络抖动 |
| 2 | 4 分钟 | 6 分钟 | 用户可能还在支付页面 |
| 3 | 8 分钟 | 14 分钟 | 用户可能切换到其他 App |
| 4 | 16 分钟 | 30 分钟 | 半小时内必须确认 |
| 5 | 32 分钟 | 62 分钟 | 1 小时内确认 |
| 6 | 64 分钟 | 126 分钟 | 2 小时内确认 |
| 7-10 | 120 分钟 | 最多 10 小时 | 最后的兜底 |

**优点:**
- 可靠性高，不依赖回调
- 指数退避避免过度查询
- 有明确的失败处理机制

**缺点:**
- 需要维护补偿任务表
- 增加系统复杂度

**适用场景:** ✅ 推荐用于金融项目

---

##### 方案 2：MQ 延迟队列 + 主动查询

**核心思路:**
使用 RabbitMQ 延迟队列或 RocketMQ 定时消息，定时触发支付状态查询。

**优点:**
- 无需定时任务，利用 MQ 的延迟能力
- 天然支持分布式

**缺点:**
- 依赖 MQ，增加系统复杂度
- 消息可能丢失（需要持久化）

**适用场景:** 已有 RocketMQ 基础设施的项目

---

##### 方案对比与推荐

| 对比维度 | 方案 1（定时任务） | 方案 2（MQ 延迟队列） |
|---------|------------------|---------------------|
| 可靠性 | ⭐⭐⭐⭐⭐ 极高 | ⭐⭐⭐⭐ 高 |
| 实现复杂度 | ⭐⭐⭐ 中等 | ⭐⭐⭐⭐ 复杂 |
| 依赖 | 仅依赖数据库 | 依赖 MQ |
| 可观测性 | ⭐⭐⭐⭐⭐ 极好（数据库可查） | ⭐⭐⭐ 一般 |

**最终推荐:** **方案 1（主动查询 + 指数退避重试）**

**推荐理由:**
1. 实现简单，依赖少
2. 补偿任务可在数据库中查询，便于排查问题
3. 指数退避策略平衡了用户体验和系统负载

---

### 问题 4：并发支付冲突未考虑

#### 问题详细分析

**风险等级:** ⚠️ 高（可能导致重复扣款）

**问题根源:**
用户在短时间内多次点击"支付"按钮，可能导致：
1. 生成多个支付订单（同一个业务订单对应多个支付订单）
2. 用户支付多次，但只收到一次商品
3. 退款流程复杂，增加客服成本

**真实场景:**
- 用户在支付页面点击"确认支付"后，因为网络慢，页面未跳转，用户又点击了一次
- 前端防抖失效（用户禁用了 JavaScript）
- 用户使用爬虫或脚本重复请求

#### 解决方案

##### 方案 1：前端防抖 + 后端幂等令牌（推荐）

**核心思路:**
前端按钮防抖 + 后端使用一次性令牌保证幂等。

**实现步骤:**

1. **前端防抖（第一道防线）**

```javascript
// React 示例
const PaymentButton = () => {
  const [loading, setLoading] = useState(false);
  
  const handlePay = async () => {
    if (loading) return;  // 防止重复点击
    
    setLoading(true);
    try {
      const response = await axios.post('/api/pay/create', orderData);
      window.location.href = response.data.paymentUrl;
    } catch (error) {
      message.error('支付失败，请重试');
      setLoading(false);  // 失败后允许重新点击
    }
    // 注意：成功后不 setLoading(false)，因为会跳转页面
  };
  
  return (
    <Button 
      onClick={handlePay} 
      loading={loading}
      disabled={loading}
    >
      {loading ? '支付中...' : '确认支付'}
    </Button>
  );
};
```

2. **后端幂等令牌（第二道防线）**

```java
@Service
@RequiredArgsConstructor
public class IdempotentTokenService {
    
    private final RedisTemplate<String, String> redisTemplate;
    
    /**
     * 生成幂等令牌
     *
     * @param userId 用户ID
     * @return 令牌
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    public String generateToken(Long userId) {
        String token = UUID.randomUUID().toString();
        String key = "idempotent:token:" + token;
        
        // 令牌有效期 5 分钟
        redisTemplate.opsForValue().set(key, String.valueOf(userId), 5, TimeUnit.MINUTES);
        
        return token;
    }
    
    /**
     * 验证并消费令牌（原子操作）
     *
     * @param token 令牌
     * @param userId 用户ID
     * @return true-验证成功 false-验证失败
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    public boolean validateAndConsumeToken(String token, Long userId) {
        String key = "idempotent:token:" + token;
        
        // 使用 Lua 脚本保证原子性：检查 + 删除
        String luaScript = 
            "if redis.call('get', KEYS[1]) == ARGV[1] then " +
            "    redis.call('del', KEYS[1]) " +
            "    return 1 " +
            "else " +
            "    return 0 " +
            "end";
        
        Long result = redisTemplate.execute(
            new DefaultRedisScript<>(luaScript, Long.class),
            Collections.singletonList(key),
            String.valueOf(userId)
        );
        
        return result != null && result == 1;
    }
}
```

3. **支付接口使用令牌**

```java
@RestController
@RequestMapping("/api/pay")
@RequiredArgsConstructor
public class PaymentController {
    
    private final IdempotentTokenService tokenService;
    private final PaymentService paymentService;
    
    /**
     * 获取幂等令牌（在进入支付页面时调用）
     *
     * @return 令牌
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @GetMapping("/token")
    public Result<String> getIdempotentToken() {
        Long userId = UserContext.getCurrentUserId();
        String token = tokenService.generateToken(userId);
        return Result.success(token);
    }
    
    /**
     * 创建支付订单
     *
     * @param request 支付请求（包含 idempotentToken）
     * @return {@link PaymentVO} 支付信息
     * @author wenlong.chen
     * @date 2026/07/01 12:00:00
     */
    @PostMapping("/create")
    public Result<PaymentVO> createPayment(@Valid @RequestBody PaymentCreateRequest request) {
        Long userId = UserContext.getCurrentUserId();
        
        // 验证并消费令牌
        boolean valid = tokenService.validateAndConsumeToken(request.getIdempotentToken(), userId);
        if (!valid) {
            throw new BusinessException("请勿重复提交");
        }
        
        // 创建支付订单
        PaymentVO payment = paymentService.createPayment(request);
        
        return Result.success(payment);
    }
}
```

**优点:**
- 双重防护，可靠性高
- 令牌一次有效，天然防重
- 实现简单

**缺点:**
- 依赖 Redis
- 用户如果支付失败，需要重新获取令牌

**适用场景:** ✅ 推荐用于所有支付场景

---

##### 方案 2：基于订单状态机 + 数据库锁

**核心思路:**
利用订单状态流转 + 数据库行锁，防止重复支付。

**优点:**
- 不依赖 Redis
- 利用数据库事务保证一致性

**缺点:**
- 行锁会阻塞其他并发请求，性能较差

**适用场景:** 小规模系统，无 Redis 环境

---

##### 方案对比与推荐

| 对比维度 | 方案 1（幂等令牌） | 方案 2（状态机 + 行锁） |
|---------|------------------|----------------------|
| 性能 | ⭐⭐⭐⭐ 高 | ⭐⭐ 低（有锁等待） |
| 可靠性 | ⭐⭐⭐⭐⭐ 极高 | ⭐⭐⭐⭐ 高 |
| 用户体验 | ⭐⭐⭐⭐⭐ 好 | ⭐⭐⭐ 一般（重复请求阻塞） |

**最终推荐:** **方案 1（前端防抖 + 后端幂等令牌）**

