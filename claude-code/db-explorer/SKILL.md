---
name: db-explorer
description: >
  只读探索 PostgreSQL、MySQL、SQLite 数据库。适用于列出表、查看表结构和字段类型、
  检查主键/索引/外键、采样表数据、执行只读 SQL（SELECT/WITH/SHOW/DESCRIBE/EXPLAIN/元数据 PRAGMA）、
  验证代码中的 model 与数据库 schema 是否一致、排查数据问题。
  仅在用户明确处于数据库/SQL/schema/record/column 语境时触发。
  不要为 HTML/Markdown/UI table 等非数据库“表格”场景触发。
---

# DB Explorer

使用 `scripts/db_query.py` 做确定性的只读数据库探索。
目标是：快速拿到可信的表、字段、约束和样例数据，而不是自由发挥数据库操作。

## Use this skill when

- 用户想看数据库里有哪些表
- 用户想看某张表的字段、类型、默认值、主键、索引、外键
- 用户想抽样看几条数据
- 用户想执行只读 SQL 查询
- 用户想确认数据库 schema 和代码中的 model / ORM 定义是否一致
- 用户在排查“这条记录实际长什么样”“某字段到底存的是什么”

## Do not use

- 任何写操作、DDL、迁移、修复数据、回填数据
- 导出大批量数据或做重型数据处理
- 用户说的“表”明显不是数据库表
- 你无法确认连接目标是否安全，且查询可能触及敏感生产数据

## Capability contract

这个 skill 当前依赖 `scripts/db_query.py`，实际支持的能力只有：

- `test`: 测试连接
- `tables`: 列出表
- `schema <table>`: 查看表结构
- `data <table> --limit N`: 采样表数据
- `query "<sql>"`: 执行只读 SQL
- `--url-env <ENV_VAR>`: 直接从指定环境变量读取连接信息

支持的输出格式：

- `table`（默认）
- `markdown`
- `json`
- `csv`

## Runtime prerequisites

- SQLite 使用 Python 标准库 `sqlite3`，不需要额外安装驱动
- PostgreSQL 依赖当前 Python 环境中的 `psycopg2`
- MySQL 依赖当前 Python 环境中的 `mysql-connector-python`

执行原则：

- 优先使用项目已经存在的虚拟环境
- 如果项目没有可用虚拟环境，再创建一个隔离虚拟环境后安装依赖
- `pip install` 和运行脚本必须使用**同一个 Python 环境**

推荐命令：

```bash
python -m venv .venv
. .venv/bin/activate
python -m pip install -r <skill-path>/requirements.txt
python <skill-path>/scripts/db_query.py --db-type postgres --url "<url>" test
```

如果用户只查 SQLite，不要要求安装 PostgreSQL/MySQL 驱动。

不要声称支持脚本没有直接实现的固定能力。  
例如，“查看建表语句”只能在你**明确写出并执行了某个数据库可用的只读查询**时再说；否则默认提供 `schema`，不要承诺 `SHOW CREATE TABLE` 一类能力在所有数据库都可用。
对于 SQLite 的 `PRAGMA`，仅把只读元数据查询视为允许范围，不要把会改状态的 `PRAGMA` 当成安全操作。

## Required inputs

只收集阻塞执行的最少信息；能推断就推断，不要机械追问。

需要的信息：

1. **数据库类型**
   - `postgres`
   - `mysql`
   - `sqlite`
2. **连接来源**
   - 连接 URL
   - 已存在的环境变量值
   - 用户提供的环境变量名
   - SQLite 文件路径
3. **查询目标**
   - 表列表 / 表结构 / 表数据 / 自定义 SQL / schema 与代码对比

推断规则：

- 用户给了 `.db` / `.sqlite` / `.sqlite3` 文件路径时，默认按 `sqlite` 处理，除非上下文明确不是
- 用户已经明确说“看 `users` 表结构”，不要再问“你想查什么”
- 用户给了环境变量名时，先读取该变量值，再用 `--url` 传给脚本；不要把 secret 原样打印出来
- 如果直接调用脚本，优先使用 `--url-env <ENV_VAR_NAME>`
- 如果用户没给连接信息，再按需检查项目里的常见配置文件（如 `.env`、`config/database.yml`、`settings.py`、`application.properties`）

## Workflow

### 1. 收集并确认连接信息

- 优先复用用户已提供的信息
- 只在确实无法执行时，问一个最关键的补充问题
- 展示连接摘要时遮掩密码：`postgresql://user:***@host:5432/db`
- 如果是从配置文件或环境变量中找到的连接信息，明确说明来源，但不要泄露密码

### 2. 先验证连接

先跑连接测试，再做后续查询：

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<connection-url>" test
```

如果连接来自环境变量，也可以：

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url-env DATABASE_URL test
```

如果失败：

- 原样保留关键错误信息
- 给出下一步建议
- 在连接未验证成功前，不要臆测表结构或数据

常见故障处理：

- 连接被拒绝：检查 host / port / 数据库服务
- 认证失败：检查用户名 / 密码 / 权限
- SQLite 文件不存在：检查路径是否相对当前工作目录
- 依赖缺失：说明缺少的 Python 包，并在当前虚拟环境里安装后重试
- MySQL C 扩展异常：脚本已默认使用 `use_pure=True`（纯 Python 实现），无需额外处理

### 3. 选择最小操作

优先用脚本内建命令，不要先写自定义元数据查询。

**列出表**

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" tables
```

**查看表结构**

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" schema <table_name>
```

**采样表数据**

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" data <table_name> --limit 10
```

**执行只读 SQL**

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" query "<sql>"
```

格式参数是全局参数，必须放在子命令前：

```bash
python <skill-path>/scripts/db_query.py --db-type sqlite --url "./app.db" --format markdown schema users
```

选择策略：

- 看库里有什么：`tables`
- 看某张表怎么定义：`schema`
- 先看看数据长什么样：`data --limit 10`
- 用户给了明确 SQL，或需要 JOIN / 聚合 / 过滤：`query`

### 4. 写 SQL 时遵守这些约束

- 只写只读 SQL
- 如果是你代写的“探索性查询”，默认加 `LIMIT 100`，除非用户明确需要更多或语义上不适合加
- 不要改写用户已经明确给出的 SQL，除非你是在补一个显然安全且用户意图就是“采样看一下”的 `LIMIT`
- 表名必须来自用户输入或已列出的真实表名，不要猜

### 5. 结果整理与输出

不要直接把终端原始输出一股脑贴给用户；先整理，再展示。

默认输出规则：

- **表列表**
  - 表少时直接列出表名
  - 如果脚本返回了 `row_count`，一起展示
- **表结构**
  - 固定关注：字段名、类型、是否可空、默认值、主键
  - 如果脚本额外打印了索引 / 外键，把它们整理成单独小节
- **数据采样**
  - 默认展示前 10 行
  - 如果列很多，先给 1 句话总结，再给表格
- **自定义查询**
  - 先一句话说清查的是什么
  - 再展示结果
  - 结果过大时只展示前几行，并明确说明已截断

展示格式建议：

- 小结果集优先 `--format markdown`
- 宽表或长结果集优先 `table`
- 需要后处理时用 `json` / `csv`

### 6. 对比代码中的 model（可选）

如果用户是在核对 ORM / model：

1. 先用 `schema` 拿数据库定义
2. 再打开对应 model 文件
3. 只报告关键差异：
   - 字段缺失 / 多出
   - 类型不一致
   - nullable 不一致
   - 默认值或主键约束不一致

不要把整段 model 和整张 schema 大段重复粘贴。

### 7. 收尾建议

查询结束后，只给与当前上下文强相关的下一步建议，例如：

- “要不要继续看这张表的样例数据？”
- “要不要我再查一下它关联的外键表？”
- “要不要我把数据库 schema 和代码里的 model 逐项对比一下？”

## Guardrails

- 这是**只读 skill**
- 不执行 INSERT / UPDATE / DELETE / DROP / ALTER / TRUNCATE / CREATE
- 不暴露密码或完整 secret
- 不把脚本的 best-effort 只读校验当作放宽边界的理由
- `data` 默认用小 limit 做采样；避免无界查询
- 查询失败时，报告事实和建议，不编造结果
- 表不存在、字段不存在、权限不足时，明确说出失败点
