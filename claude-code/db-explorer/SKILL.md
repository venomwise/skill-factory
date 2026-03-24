---
name: db-explorer
description: >
  查询和探索数据库的表结构、字段信息和表数据。支持 PostgreSQL、MySQL、SQLite。
  当用户提到数据库查询、查看表结构、查看表数据、检查字段类型、列出数据库中的表、
  查看建表语句、采样数据、调试数据问题、了解数据库 schema 时，都应该使用这个 skill。
  即使用户只是说"看看这个表长什么样"或"帮我查一下这个数据"也要触发。
---

# DB Explorer

快速探索数据库：列出表、查看表结构、采样数据、执行自定义查询。
面向开发调试场景——帮助开发者在编码过程中快速理解数据库的结构和数据。

## When to use

- 查看数据库中有哪些表
- 查看某张表的字段名、类型、约束（主键、外键、索引等）
- 快速采样表中的数据（前 N 行）
- 执行自定义 SQL 查询并格式化输出
- 开发过程中需要确认表结构与代码中 model 定义是否一致
- 调试数据问题（检查某些记录的实际值）

## When NOT to use

- 数据迁移或 schema 变更操作（ALTER TABLE、DROP 等）——这个 skill 是只读的
- 大规模数据导出（超过几百行的数据处理）
- 生产环境的敏感数据查询——请确认安全策略

## Inputs

用户需要提供以下信息（通过对话收集）：

1. **数据库类型**: PostgreSQL / MySQL / SQLite
2. **连接信息**（以下任一方式）:
   - 连接 URL（如 `postgresql://user:pass@host:port/dbname`）
   - 环境变量名（如 `DATABASE_URL`）
   - SQLite 文件路径（如 `./data/app.db`）
   - 分别提供 host、port、user、password、database
3. **查询意图**: 想看什么——表列表？某张表的结构？某张表的数据？自定义 SQL？

## Outputs

根据查询类型，输出格式会自动适配：

- **表列表**: 简洁的表名列表，包含行数统计（如果可获取）
- **表结构**: 字段名、类型、是否可空、默认值、约束信息，以易读的表格呈现
- **数据采样**: 格式化的表格数据（默认 10 行，可指定）
- **自定义查询**: 按结果集大小选择最合适的展示方式

当数据量小时用 Markdown 表格，数据量大时用对齐的文本表格或建议导出为文件。

## Workflow

### 1. 收集连接信息

如果用户没有明确提供连接信息，通过对话确认：

- 数据库类型是什么？
- 连接方式：URL、环境变量、还是文件路径（SQLite）？
- 如果信息不完整，检查项目中常见的配置文件（`.env`、`config/database.yml`、`settings.py`、`application.properties` 等）是否包含数据库连接信息

对于 SQLite，只需要文件路径。对于 PostgreSQL/MySQL，需要完整的连接参数。

**安全提示**: 永远不要在输出中明文显示密码。如果从配置文件读取了密码，在展示连接信息时用 `***` 遮掩。

### 2. 验证连接

使用 `scripts/db_query.py` 测试连接是否成功：

```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<connection-url>" test
```

如果连接失败，分析错误信息并给出针对性建议：
- 连接被拒绝 → 检查 host/port 是否正确，数据库服务是否运行
- 认证失败 → 检查用户名密码
- 数据库不存在 → 列出可用的数据库

### 3. 执行查询

根据用户意图选择操作：

**列出所有表:**
```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" tables
```

**查看表结构:**
```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" schema <table_name>
```

**采样数据:**
```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" data <table_name> --limit <N>
```

**执行自定义 SQL:**
```bash
python <skill-path>/scripts/db_query.py --db-type <type> --url "<url>" query "<sql>"
```

所有命令都支持 `--format` 参数选择输出格式：`table`（默认）、`markdown`、`json`、`csv`。
注意 `--format` 是全局参数，必须放在子命令（tables/schema/data/query）之前：

```bash
python <skill-path>/scripts/db_query.py --db-type sqlite --url "./app.db" --format markdown schema users
```

### 4. 格式化与呈现

拿到查询脚本的原始输出后，根据场景组织展示：

- 如果用户只是想快速了解表结构，直接展示表格即可
- 如果用户在对比代码中的 model 定义，可以将表结构和代码 model 并排对比
- 如果采样数据量较大（超过 20 行），建议写入文件而非直接输出到终端
- 对于自定义查询，先展示结果行数，再展示数据

### 5. 后续建议

查询完成后，根据上下文给出有价值的后续建议：

- 如果在查看表结构，提示可以查看关联表或索引
- 如果在采样数据，提示可以用 WHERE 条件过滤
- 如果发现表结构和代码不一致，指出差异

## Safety & guardrails

- **只读操作**: 这个 skill 只执行 SELECT 查询。脚本会拒绝任何包含 INSERT、UPDATE、DELETE、DROP、ALTER、TRUNCATE 的语句
- **结果限制**: 数据查询默认限制 100 行，避免意外拉取整张大表
- **密码保护**: 不在输出中展示数据库密码
- **超时控制**: 查询超时默认 30 秒，防止慢查询阻塞
- **注入防护**: 脚本使用参数化查询处理用户输入（表名除外——表名通过白名单验证）
