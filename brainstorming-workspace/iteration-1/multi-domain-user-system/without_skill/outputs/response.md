# 完整用户系统设计方案

很好的需求！一个完整的用户系统涵盖四个核心模块：**注册登录**、**权限管理**、**用户画像分析**和**社交关系**。下面我会从系统架构、数据库设计、核心功能到技术选型，给出一个完整的方案。

---

## 一、系统整体架构

```
┌─────────────────────────────────────────────────────────┐
│                      API Gateway                         │
│              (认证拦截 / 限流 / 路由)                    │
└──────────┬──────────┬─────────────┬──────────────────────┘
           │          │             │             │
    ┌──────┴───┐ ┌────┴────┐ ┌─────┴─────┐ ┌────┴──────┐
    │ 认证服务  │ │权限服务  │ │画像分析服务│ │社交关系服务│
    │  Auth    │ │  RBAC   │ │ Profile   │ │  Social   │
    └──────────┘ └─────────┘ └───────────┘ └───────────┘
           │          │             │             │
    ┌──────┴──────────┴─────────────┴─────────────┴───────┐
    │                    数据存储层                         │
    │   MySQL(主数据) + Redis(缓存) + ClickHouse(分析)     │
    └─────────────────────────────────────────────────────┘
```

---

## 二、模块一：注册登录系统

### 2.1 支持的登录方式

| 方式 | 说明 |
|------|------|
| 用户名/密码 | 基础登录，密码 bcrypt 加密存储 |
| 手机号 + 验证码 | 短信 OTP，有效期 5 分钟 |
| 邮箱 + 验证码 | 邮件 OTP，有效期 15 分钟 |
| OAuth 第三方 | 微信、Google、GitHub 等 |
| SSO 单点登录 | 企业场景，支持 SAML / OIDC |

### 2.2 Token 方案

采用 **双 Token 机制**：
- **Access Token**：JWT，有效期 15 分钟，无状态验证
- **Refresh Token**：随机字符串，有效期 7 天，存入 Redis + 数据库，支持轮换（Rotation）

```
登录流程：
  用户凭证 → 验证 → 生成 Access Token + Refresh Token
                      ↓
               Refresh Token 存 Redis（支持主动吊销）

续签流程：
  Access Token 过期 → 客户端用 Refresh Token 换新 Token 对
                      ↓
               旧 Refresh Token 作废，新 Refresh Token 返回
```

### 2.3 安全措施

- **密码策略**：最小长度 8 位，强制复杂度，bcrypt cost=12
- **登录保护**：连续失败 5 次锁定账号 15 分钟（Redis 计数器）
- **多设备管理**：记录设备指纹，支持查看并踢出指定设备
- **异地登录提醒**：IP 地理位置检测，异常登录发邮件/短信提醒
- **MFA 多因素认证**：可选启用 TOTP（Google Authenticator 兼容）

### 2.4 核心数据表

```sql
-- 用户主表
CREATE TABLE users (
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid        VARCHAR(36) UNIQUE NOT NULL,     -- 对外暴露的 ID
    username    VARCHAR(50) UNIQUE,
    email       VARCHAR(255) UNIQUE,
    phone       VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255),
    status      ENUM('active','disabled','locked') DEFAULT 'active',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME ON UPDATE CURRENT_TIMESTAMP
);

-- OAuth 绑定表
CREATE TABLE user_oauth (
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id     BIGINT NOT NULL,
    provider    VARCHAR(20) NOT NULL,            -- 'wechat','google','github'
    open_id     VARCHAR(255) NOT NULL,
    union_id    VARCHAR(255),
    UNIQUE KEY (provider, open_id)
);

-- 登录会话表
CREATE TABLE user_sessions (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id         BIGINT NOT NULL,
    refresh_token   VARCHAR(255) UNIQUE NOT NULL,
    device_info     JSON,                        -- 设备指纹信息
    ip_address      VARCHAR(45),
    last_active_at  DATETIME,
    expires_at      DATETIME NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 三、模块二：权限管理系统（RBAC）

### 3.1 权限模型选择

采用 **RBAC（Role-Based Access Control）+ 资源策略** 的混合模型：

```
用户 (User)
  └── 拥有多个角色 (Role)
        └── 角色包含多个权限 (Permission)
              └── 权限对应资源 + 操作 (Resource:Action)

例：
  用户 Alice → 角色「内容编辑」→ 权限「article:create, article:edit」
  用户 Bob  → 角色「管理员」  → 权限「*:*」（所有权限）
```

### 3.2 权限粒度设计

| 级别 | 示例 |
|------|------|
| 菜单权限 | 能否看到某个菜单项 |
| 操作权限 | 能否点击「导出」「删除」按钮 |
| 数据权限 | 只能看自己部门的数据 |
| 字段权限 | 某些字段对特定角色隐藏（如手机号脱敏） |

### 3.3 核心数据表

```sql
-- 角色表
CREATE TABLE roles (
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(255),
    is_system   BOOLEAN DEFAULT FALSE,           -- 系统内置角色不可删除
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 权限表
CREATE TABLE permissions (
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource    VARCHAR(100) NOT NULL,           -- 资源，如 'article'
    action      VARCHAR(50) NOT NULL,            -- 操作，如 'create'
    description VARCHAR(255),
    UNIQUE KEY (resource, action)
);

-- 角色-权限关联
CREATE TABLE role_permissions (
    role_id       BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    PRIMARY KEY (role_id, permission_id)
);

-- 用户-角色关联（支持多租户）
CREATE TABLE user_roles (
    user_id     BIGINT NOT NULL,
    role_id     BIGINT NOT NULL,
    tenant_id   BIGINT DEFAULT 0,               -- 0 表示全局角色
    granted_by  BIGINT,
    granted_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at  DATETIME,                        -- 支持临时授权
    PRIMARY KEY (user_id, role_id, tenant_id)
);
```

### 3.4 权限校验流程

```
请求到达 → API Gateway 验证 Access Token
           ↓ Token 有效
        从 Redis 缓存获取用户权限集合（缓存 5 分钟）
           ↓ 缓存未命中
        查询 user_roles → role_permissions → permissions 表
           ↓ 组装权限集合
        存入 Redis 缓存
           ↓
        判断当前接口所需权限是否在集合中
           ↓ 有权限          ↓ 无权限
        放行请求           返回 403
```

### 3.5 动态权限刷新

管理员修改用户角色后，通过消息队列（MQ）发布「权限变更」事件，订阅方主动删除对应用户的 Redis 权限缓存，确保权限实时生效。

---

## 四、模块三：用户画像分析

### 4.1 画像数据来源

| 来源 | 数据类型 |
|------|----------|
| 注册信息 | 年龄、性别、地区、职业（显性属性） |
| 行为日志 | 浏览、点击、搜索、购买、停留时长 |
| 社交行为 | 关注的人、加入的群组、互动内容 |
| 设备信息 | 手机型号、操作系统、常用时段 |
| 偏好设置 | 用户主动设置的兴趣标签 |

### 4.2 画像标签体系

```
用户标签
├── 基础属性
│   ├── 人口统计：性别、年龄段、地区
│   └── 账户属性：注册时长、活跃度等级
├── 兴趣偏好
│   ├── 内容类别：科技、体育、娱乐...
│   └── 消费偏好：价格敏感度、品类偏好
├── 行为特征
│   ├── 活跃时段：早鸟型、夜猫型
│   ├── 使用频率：高频/中频/低频
│   └── 流失风险：近 N 天未登录
└── 价值分层
    └── RFM 模型：最近购买、购买频率、购买金额
```

### 4.3 数据架构

```
行为事件上报（埋点 SDK）
    ↓
Kafka 消息队列（削峰 / 解耦）
    ↓
Flink 实时流处理
    ├── 实时标签更新（Redis）
    └── 写入 ClickHouse（历史明细）
         ↓
    定时批处理（T+1 离线任务）
         ↓
    更新 MySQL 用户画像宽表
```

### 4.4 核心数据表

```sql
-- 用户画像宽表（MySQL）
CREATE TABLE user_profiles (
    user_id         BIGINT PRIMARY KEY,
    age_group       VARCHAR(10),                 -- '18-24','25-34'...
    gender          TINYINT,                     -- 0未知 1男 2女
    city            VARCHAR(50),
    active_level    TINYINT,                     -- 1低 2中 3高
    last_active_at  DATETIME,
    register_days   INT,                         -- 注册天数
    interest_tags   JSON,                        -- ['tech','sports'...]
    rfm_score       DECIMAL(5,2),
    risk_churn      TINYINT,                     -- 流失风险 0-100
    updated_at      DATETIME ON UPDATE CURRENT_TIMESTAMP
);

-- 用户行为事件表（ClickHouse，适合大数据量）
CREATE TABLE user_events (
    event_id    String,
    user_id     Int64,
    event_type  String,                          -- 'view','click','purchase'
    target_id   String,                          -- 目标资源 ID
    properties  String,                          -- JSON 扩展属性
    client_ip   String,
    device_type String,
    event_time  DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(event_time)
ORDER BY (user_id,