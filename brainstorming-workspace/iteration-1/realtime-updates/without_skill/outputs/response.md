# 为 REST API 添加实时推送能力

## 背景

传统 REST API 是请求-响应模型，客户端需要主动轮询才能获取最新数据。要实现数据变更时客户端能立即收到更新，有以下几种主流方案。

---

## 方案一：Server-Sent Events（SSE）

**适用场景**：服务器单向推送，客户端只需接收（如通知、日志流、数据监控）。

**原理**：基于 HTTP 长连接，服务器持续向客户端发送事件流（text/event-stream）。

### 后端示例（Node.js / Express）

```js
app.get('/api/events', (req, res) => {
  res.setHeader('Content-Type', 'text/event-stream');
  res.setHeader('Cache-Control', 'no-cache');
  res.setHeader('Connection', 'keep-alive');

  const sendEvent = (data) => {
    res.write(`data: ${JSON.stringify(data)}\n\n`);
  };

  // 监听数据变更，推送给客户端
  const unsubscribe = dataChangeEmitter.on('change', sendEvent);

  req.on('close', () => {
    unsubscribe();
  });
});
```

### 前端示例

```js
const source = new EventSource('/api/events');
source.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('收到更新:', data);
};
```

**优点**：
- 实现简单，基于标准 HTTP
- 浏览器原生支持，无需额外库
- 自动重连机制

**缺点**：
- 单向通信（服务器 → 客户端）
- 浏览器有并发连接数限制（HTTP/1.1 下每域名 6 个）

---

## 方案二：WebSocket

**适用场景**：需要双向通信（如聊天、协作编辑、游戏）。

**原理**：通过 HTTP Upgrade 握手建立持久化全双工连接。

### 后端示例（Node.js / ws 库）

```js
const WebSocket = require('ws');
const wss = new WebSocket.Server({ port: 8080 });

wss.on('connection', (ws) => {
  // 监听数据变更，推送给所有连接的客户端
  const unsubscribe = dataChangeEmitter.on('change', (data) => {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(data));
    }
  });

  ws.on('close', () => {
    unsubscribe();
  });

  ws.on('message', (message) => {
    // 处理客户端发来的消息（双向通信）
    console.log('客户端消息:', message);
  });
});
```

### 前端示例

```js
const ws = new WebSocket('ws://localhost:8080');
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('收到更新:', data);
};
```

**优点**：
- 双向实时通信
- 低延迟、低开销（无 HTTP 头部重复传输）

**缺点**：
- 需要处理连接管理、断线重连
- 某些代理/防火墙可能阻断 WebSocket
- 服务端需维护大量长连接，对资源要求较高

---

## 方案三：长轮询（Long Polling）

**适用场景**：对实时性要求不极致，但需要兼容性最好的方案（无需特殊协议支持）。

**原理**：客户端发起请求，服务器挂起直到有新数据或超时，然后返回响应，客户端立即发起下一次请求。

### 后端示例（Node.js / Express）

```js
app.get('/api/poll', (req, res) => {
  const timeout = setTimeout(() => {
    res.json({ type: 'timeout' });
  }, 30000); // 30秒超时

  dataChangeEmitter.once('change', (data) => {
    clearTimeout(timeout);
    res.json({ type: 'update', data });
  });
});
```

### 前端示例

```js
async function poll() {
  try {
    const res = await fetch('/api/poll');
    const result = await res.json();
    if (result.type === 'update') {
      console.log('收到更新:', result.data);
    }
  } finally {
    poll(); // 立即发起下一次请求
  }
}
poll();
```

**优点**：
- 兼容性最好，纯 HTTP
- 适合防火墙严格的环境

**缺点**：
- 延迟较高（每次需要重新建立连接）
- 服务器并发连接数增多

---

## 方案四：使用消息队列 + WebSocket/SSE（生产级推荐）

对于生产环境，通常结合消息队列（如 Redis Pub/Sub、Kafka、RabbitMQ）实现跨实例的实时推送。

```
数据写入 → 数据库 → 触发事件 → 发布到 Redis Pub/Sub
                                         ↓
              WebSocket Server 订阅 → 推送给对应客户端
```

### 示例（Redis Pub/Sub + WebSocket）

```js
const redis = require('redis');
const subscriber = redis.createClient();

subscriber.subscribe('data-changes');

subscriber.on('message', (channel, message) => {
  const data = JSON.parse(message);
  // 找到对应的 WebSocket 连接并推送
  broadcastToClients(data);
});

// 数据写入时发布事件
async function updateData(id, payload) {
  await db.update(id, payload);
  publisher.publish('data-changes', JSON.stringify({ id, payload }));
}
```

**优点**：
- 支持水平扩展（多实例部署）
- 解耦数据层和推送层
- 高可靠性

---

## 方案对比总结

| 方案 | 方向 | 复杂度 | 延迟 | 兼容性 | 适用场景 |
|------|------|--------|------|--------|----------|
| SSE | 单向（服务器→客户端）| 低 | 低 | 好（现代浏览器）| 通知、数据流 |
| WebSocket | 双向 | 中 | 极低 | 好 | 聊天、协作、游戏 |
| 长轮询 | 单向 | 低 | 中 | 最好 | 兼容性要求高的场景 |
| 消息队列 + WS | 双向 | 高 | 极低 | 好 | 生产级、多实例部署 |

---

## 推荐选择

- **快速上手、单向推送**：用 **SSE**，实现最简单，HTTP 原生支持。
- **需要双向通信**：用 **WebSocket**（配合 `socket.io` 库可简化开发）。
- **生产级多实例**：用 **Redis Pub/Sub + WebSocket**，保证跨实例推送。
- **极端兼容性要求**：用**长轮询**作为降级方案。

如果你告诉我当前的技术栈（语言/框架）和具体业务场景，我可以给出更有针对性的实现方案。
