# Spring Cloud Agent Context Initializer Skill 设计草案

## 1. 背景

用户希望创建一个 skill，用于在 Spring Cloud / Spring Boot 多微服务项目中初始化 `AGENTS.md` 或 `CLAUDE.md`。

该 skill 的核心价值不是生成普通项目 README，而是为后续 AI Coding Agent 建立可执行的项目上下文，使 Agent 在修改代码前理解：

- 项目整体架构
- 微服务边界
- Maven 模块关系
- 服务间通信方式
- Spring Cloud 基础设施
- 每个服务的职责、入口和修改风险

用户明确提出文档应采用“总分”结构：

- 根目录：创建一个总览型 `AGENTS.md` / `CLAUDE.md`
- 微服务目录：在每个可独立启动的微服务 `pom.xml` 同级创建服务级 `AGENTS.md` / `CLAUDE.md`

## 2. 设计目标

### 2.1 核心目标

该 skill 应指导 AI Agent 完成以下任务：

1. 分析 Spring Cloud Maven 项目的 `pom.xml` 结构。
2. 识别哪些目录是可独立启动的微服务，哪些只是普通模块。
3. 分析 Spring Cloud 相关组件，例如：
   - 注册中心
   - 配置中心
   - 网关
   - OpenFeign
   - MQ
   - Redis
   - 数据库访问
   - 分布式事务
   - tracing / metrics
4. 分析服务间通信关系。
5. 选择创建 `AGENTS.md` 还是 `CLAUDE.md`。
6. 生成根目录总览文档。
7. 为每个微服务生成服务级文档。
8. 保留已有人工说明，避免盲目覆盖。
9. 输出可供后续 Agent 安全使用的上下文说明。

### 2.2 成功标准

一次成功执行后，应满足：

- 根目录存在一个总览型 `AGENTS.md` 或 `CLAUDE.md`。
- 每个确认的微服务目录存在一个服务级 `AGENTS.md` 或 `CLAUDE.md`。
- 根文档列出所有微服务，并链接到对应服务级文档。
- 普通共享模块不会被误判为微服务。
- 文档中的服务职责、启动类、端口、服务名、通信关系都有证据来源；缺失或不确定的信息统一收敛到 `Open Questions`，不在模板中机械填充 `Unknown`。
- 文档包含对后续 Agent 修改代码有价值的约束，而不是泛泛的项目介绍。

## 3. 非目标

该 skill 不应该做以下事情：

- 不修改业务代码。
- 不自动重构项目结构。
- 不创建新的 Spring Cloud 服务。
- 不执行数据库迁移。
- 不自动提交 Git commit。
- 不凭服务名或目录名脑补架构和职责。
- 不为所有 Maven 模块强制创建服务级文档。
- 不把生成目录，例如 `target/`，纳入分析。

## 4. 适用范围

### 4.1 当前版本范围

第一版建议只支持：

- Maven 项目
- `pom.xml`
- Spring Boot / Spring Cloud
- Java / Kotlin 源码
- 多模块或多目录微服务结构

### 4.2 暂不支持或弱支持

以下内容可作为后续增强，不建议第一版强行支持：

- Gradle 项目
- 非 JVM 微服务
- Kubernetes Helm Chart 深度分析
- CI/CD workflow 深度分析
- 自动生成架构图图片
- 根据运行时环境访问注册中心或配置中心

## 5. 触发语义设计

该 skill 应在用户表达以下意图时触发：

- “帮我给这个 Spring Cloud 项目生成 AGENTS.md”
- “初始化 CLAUDE.md”
- “分析这个微服务项目并生成 agent 上下文”
- “根据 pom.xml 分析微服务结构”
- “给每个微服务生成 CLAUDE.md”
- “生成总分式 AGENTS.md”
- “给这个多模块 Spring Boot 项目创建 agent instructions”
- “整理 Spring Cloud 项目架构文档给 AI Agent 用”

建议 skill 名称：

```yaml
name: springcloud-init
```

建议 description：

```yaml
description: >
  Use this skill when the user asks to initialize, generate, or update AGENTS.md
  or CLAUDE.md for a Spring Cloud / Spring Boot multi-service Maven project.
  It analyzes pom.xml files, Spring Boot entry points, application.yml/bootstrap.yml,
  service discovery/configuration components, Feign/Gateway/MQ communication, and
  creates a root guide plus per-microservice guide files. Trigger for requests about
  Spring Cloud architecture docs, AI agent project context, microservice analysis,
  or 总分式 AGENTS.md/CLAUDE.md generation.
```

设计原因：

- 包含目标文件：`AGENTS.md` / `CLAUDE.md`
- 包含项目类型：Spring Cloud / Spring Boot / Maven
- 包含关键分析对象：`pom.xml`、启动类、配置文件、服务通信
- 包含中文触发词：“总分式”
- 语义足够强，降低 under-triggering 风险

## 6. 文件命名策略

用户希望由 Agent 自己决定创建 `AGENTS.md` 还是 `CLAUDE.md`。

建议规则：

1. 如果根目录已有 `AGENTS.md`，继续使用 `AGENTS.md`。
2. 如果根目录没有 `AGENTS.md`，但已有 `CLAUDE.md`，继续使用 `CLAUDE.md`。
3. 如果两者都没有，默认创建 `AGENTS.md`。
4. 服务级文件使用与根目录一致的文件名。
5. 如果两者都存在：
   - 优先维护 `AGENTS.md`。
   - 读取 `CLAUDE.md`。
   - 将仍有效的规则引用或摘要到 `AGENTS.md`。
   - 不随意删除任何现有人工规则。

设计原因：

- `AGENTS.md` 更通用，适合多 Agent。
- `CLAUDE.md` 更偏 Claude Code。
- 复用已有文件可以避免破坏项目约定。

## 7. 工作流设计

### Step 1. 确认目标目录和文件类型

Agent 应先判断当前目录是否为项目根目录，检查：

- 是否存在根 `pom.xml`
- 是否存在多个子模块 `pom.xml`
- 是否已有 `AGENTS.md` / `CLAUDE.md`

若当前目录不明显是项目根目录，应询问用户或说明假设。

### Step 2. 读取已有上下文文件

如果已有 `AGENTS.md` 或 `CLAUDE.md`，必须先读取并识别：

- 代码规范
- Git 分支规范
- 安全/部署警告
- 构建命令
- 用户手写说明

更新时应保留这些内容，除非用户明确要求重写。

为支持安全复跑，Agent 生成的内容必须包裹在显式标记中：

```text
<!-- managed:springcloud-init -->
... Agent 生成内容 ...
<!-- /managed:springcloud-init -->
```

复跑或更新时：

- 只重写标记区内的内容。
- 标记区外的所有内容（人工说明、规范、警告）一律保留。
- 若已有文件没有标记（首次为人工文件补全），不删除原有内容，只在末尾追加一个带标记的托管区块。

### Step 3. 扫描 Maven 模块

递归查找 `pom.xml`，但排除：

```text
target/
.git/
node_modules/
build/
.idea/
```

对每个 `pom.xml` 建立模块记录：

```text
path
artifactId
groupId
packaging
parent
modules
internal dependencies
Spring Boot dependencies
Spring Cloud dependencies
```

### Step 4. 识别微服务与共享模块

在 `scan_modules.py` 输出的 Maven 清单基础上，对候选模块做轻量信号扫描，再进行分类。轻量信号包括启动类、`application` / `bootstrap` 配置、`spring.application.name`、`server.port` 和少量入口注解；深度源码分析留到 Step 5 之后，只针对确认的微服务进行。

不能简单认为“有 `pom.xml` 就是微服务”。

一个目录更可能是微服务，如果存在以下强证据：

- `@SpringBootApplication`
- `SpringApplication.run(...)`
- `spring.application.name`
- `server.port`
- `spring-boot-starter-web`
- `spring-boot-starter-webflux`
- `spring-cloud-starter-gateway`
- 注册中心依赖
- 配置中心依赖
- Controller / Listener / Scheduled Job / Gateway Route

一个目录更可能是共享模块，如果：

- `packaging` 是 `pom`
- 没有启动类
- 没有独立服务名或端口
- 主要包含 common / api / dto / model / util / mapper / starter 等复用代码
- 被其他服务依赖但自身不能独立启动

分类结果建议有三类：

```text
Microservice
Shared module
Needs verification
```

### Step 5. 分析 Spring Cloud 基础组件

采用广度优先策略以控制上下文与成本：Step 3 只做“只读 pom 建立全量清单”；Step 4 对候选模块做轻量信号扫描并分类；从本步开始的深度分析，只针对 Step 4 确认为 `Microservice` 的目录进行，不对 shared module 做昂贵的全量扫描。

从 `pom.xml`、配置文件、注解和代码中识别：

| 组件 | 证据来源 |
|---|---|
| Nacos Discovery | 依赖、`spring.cloud.nacos.discovery` |
| Nacos Config | 依赖、`spring.cloud.nacos.config` |
| Eureka | 依赖、`eureka.client` |
| Consul | 依赖、`spring.cloud.consul` |
| Gateway | 依赖、`spring.cloud.gateway.routes`、RouteLocator |
| OpenFeign | `@FeignClient`、starter-openfeign |
| MQ | listener 注解、producer API、配置 key |
| Redis | redis starter、`spring.redis` |
| DB | datasource、MyBatis/JPA 依赖、mapper/entity |
| Seata | seata 依赖、配置 |
| Sentinel | sentinel 依赖、配置 |
| Tracing | Sleuth、Zipkin、Micrometer、OpenTelemetry |

### Step 6. 分析服务通信

通信关系需要有证据。

主要证据来源：

- `@FeignClient(name = "...")`
- `@FeignClient(value = "...")`
- Gateway routes 中的 `lb://service-name`
- `RestTemplate` 调用
- `WebClient` 调用
- MQ topic / queue producer 和 listener
- Dubbo `@DubboReference` / `@DubboService`
- shared API module 依赖
- controller path 和包名

通信关系应分为：

```text
Inbound
Outbound
Messaging
Inferred / Needs verification
```

### Step 6.5 计划确认检查点

批量写入任何文件前，必须先产出一份计划清单供用户确认（plan-validate-execute）。原因：在 30+ 服务的 monorepo 中，一次性生成几十个文件属于中等 blast radius 的动作，且服务识别可能误判，应在动手前拦截错误。

清单应包含：

- Microservices 清单 + 每个的判定证据。
- Shared modules 清单 + 判定证据。
- `Needs verification` 清单。
- 本次选择使用的上下文文件名：`AGENTS.md` 或 `CLAUDE.md`，以及选择原因。
- 计划写入或更新的文件路径列表。

用户确认后再进入 Step 7/8 批量写入。若用户要求自主执行，可跳过确认，但仍需先输出该清单作为执行记录。

### Step 7. 生成根目录总文档

根文档主要说明整体架构和服务分布，不应展开到每个类。

建议包含：

- Project Overview
- Architecture Summary
- Microservices
- Service Guides
- Shared Modules
- Communication Map
- Spring Cloud Infrastructure
- Build and Run
- Code Guidelines for Agents
- Git and Branch Rules
- Evidence and Open Questions

### Step 8. 生成微服务级文档

每个微服务的文档放在该微服务 `pom.xml` 同级。

服务文档主要说明该服务如何运行、提供什么能力、调用谁、被谁调用以及修改风险。

建议包含：

- Service Overview
- Runtime Identity
- Responsibilities
- Inbound APIs
- Outbound Calls
- Messaging
- Scheduled Jobs
- Internal Dependencies
- Configuration
- Data Access
- Local Development
- Agent Notes
- Evidence Used
- Open Questions

### Step 9. 验证输出

写完后，将以下 checklist 复制到回复中逐项勾选，降低跳步概率：

```text
验证清单：
- [ ] 用户已确认计划清单，或用户明确要求自主执行
- [ ] 根文档引用的每个服务级文件都已存在
- [ ] 每个微服务都有启动证据（主类 / 服务名 / 端口 / web 依赖）
- [ ] 没有把 shared module 误判为 service
- [ ] 没有漏掉明显的微服务
- [ ] 已保留标记区外的人工规则
- [ ] 没有修改 managed 标记区外的人工内容
- [ ] 生成内容都在 managed 标记区内
- [ ] 所有不确定点都收敛到 Open Questions
```

## 8. 根目录文档结构设计

模板使用原则（适用于第 8、9 节所有模板）：

- 模板是可裁剪的默认结构，不是必填骨架。
- 有证据才写对应小节；没有证据的小节直接省略，不要逐字段填 `Unknown`。
- 所有不确定项收敛到唯一的 `Open Questions`，不散落成满屏 `Unknown`。
- 目标是“最小有效文档”，而不是把模板填满。

建议模板：

```md
# Project Agent Guide

## Project Overview

简述项目是 Spring Cloud / Spring Boot 微服务项目，以及核心业务域。

## Architecture Summary

- Service discovery:
- Config center:
- Gateway:
- Service communication:
- Messaging:
- Data stores:
- Observability:

## Microservices

| Service | Path | Application Name | Main Class | Port | Responsibility |
|---|---|---|---|---|---|

## Service Guides

- [service-name](./path/to/service/AGENTS.md)

## Shared Modules

| Module | Path | Purpose | Used By |
|---|---|---|---|

## Communication Map

| Source | Target | Mechanism | Evidence |
|---|---|---|---|

## Spring Cloud Infrastructure

## Build and Run

## Code Guidelines for Agents

- 修改 Feign 接口时，同步检查调用方和实现方。
- 修改 DTO / VO / API module 时，检查依赖该模块的所有服务。
- 修改配置 key 时，检查 `bootstrap.yml`、`application.yml`、配置中心和 profile 文件。
- 修改数据库实体时，检查 mapper、XML、migration、SQL 脚本。
- 不要跨服务随意改接口契约。
- 不要把业务逻辑写入 Controller。
- 修改 shared module 时，检查所有依赖方。

## Git and Branch Rules

仅记录项目已有规则或用户明确提供的规则。
如果没有发现规则，写 `Not found`，不要编造。

## Evidence and Open Questions
```

## 9. 微服务级文档结构设计

建议模板：

```md
# <service-name> Agent Guide

## Service Overview

## Runtime Identity

| Item | Value |
|---|---|
| Artifact ID | |
| Spring Application Name | |
| Main Class | |
| Port | |
| Profiles | |

## Responsibilities

## Inbound APIs

| Entry | Path / Topic | Purpose | Evidence |
|---|---|---|---|

## Outbound Calls

| Target | Mechanism | Purpose | Evidence |
|---|---|---|---|

## Messaging

## Scheduled Jobs

## Internal Dependencies

## Configuration

## Data Access

## Local Development

## Agent Notes

## Evidence Used

## Open Questions
```

## 10. Prompt Engineering 设计要点

### 10.1 避免 under-triggering

skill description 需要显式包含：

- `AGENTS.md`
- `CLAUDE.md`
- Spring Cloud
- Spring Boot
- Maven
- `pom.xml`
- 微服务分析
- 总分式
- Agent context

### 10.2 防止 Agent 编造

skill 中需要明确：

- 没有证据的小节直接省略，不确定项收敛到 `Open Questions`，不要逐字段填 `Unknown`。
- 不要根据 artifactId 脑补职责。
- 服务职责必须来自 Controller、配置、类名、包名、README、接口或依赖关系等证据。
- 通信关系必须能追溯到 Feign、Gateway、MQ、HTTP client、Dubbo 或配置。

### 10.3 使用渐进披露

主 `SKILL.md` 不宜过长。建议结构：

```text
springcloud-init/
├── SKILL.md
├── scripts/
│   └── scan_modules.py
└── references/
    ├── detection.md
    └── templates.md
```

其中：

- `SKILL.md`：核心流程和规则。
- `scripts/scan_modules.py`：只枚举不判定的 Maven 模块扫描脚本（见 12 节）。
- `references/detection.md`：微服务识别、Spring Cloud 组件识别、通信识别细则。
- `references/templates.md`：根文档和服务文档模板，附一个填好的正例和一个臃肿反例。

## 11. AI Agent Engineering 设计要点

### 11.1 Read-before-write

该 skill 应要求 Agent：

- 写入前读取已有文件。
- 理解现有人工约束。
- 精准更新而非粗暴覆盖。

### 11.2 Evidence-first

文档应尽量有 evidence，例如：

```md
## Evidence Used

- `pom.xml`: artifactId = user-service
- `src/main/java/.../UserApplication.java`: `@SpringBootApplication`
- `application.yml`: `spring.application.name=user-service`
- `UserController.java`: exposes `/users`
- `OrderFeignClient.java`: calls `order-service`
```

### 11.3 最小有效文档

服务级文档不应变成类清单。应聚焦：

- 服务身份
- 服务职责
- 入口
- 外部调用
- 配置
- 数据访问
- 修改风险

为把该原则从抽象规则变成可对照的质量锚点，`references/templates.md` 中应提供：

- 一个填好的 mini 服务文档正例（含真实证据引用）。
- 一个“类清单式臃肿文档”反例，明确标注为不应产出的形态。

### 11.4 处理不确定性

Agent 应显式输出不确定项，例如：

```md
## Open Questions

- `payment-service` has `pom.xml` and web dependency, but no Spring Boot main class was found.
- Gateway routes reference `inventory-service`, but no matching module was found in this repository.
```

## 12. 推荐实现文件结构

如果设计通过，建议创建：

```text
springcloud-init/
├── SKILL.md
├── scripts/
│   └── scan_modules.py
└── references/
    ├── detection.md
    └── templates.md
```

第一版即引入辅助脚本 `scan_modules.py`，但严格限定职责为“只枚举，不判定”：

- 递归查找 `pom.xml`（排除 `target/`、`.git/`、`node_modules/`、`build/`、`.idea/`）。
- 对每个模块输出结构化 JSON：`path`、`artifactId`、`groupId`、`packaging`、`parent`、`modules`、依赖列表。
- 不判断“是不是微服务”，该判定仍由 Agent 依据 Step 4 的证据规则完成。

如此划分的原因：

- 机械的递归扫描与 `pom.xml` 解析最耗 token 且最容易在大仓里出错，下沉到脚本可省 token、保证多次运行结果一致。
- 判定逻辑留给 Agent，保留判断过程的透明度，符合 evidence-first。
- 若后续评测发现判定不稳定，再考虑把部分判定信号也加入脚本输出。

## 13. 推荐测试用例

### Test 1：简单三服务项目

输入场景：

```text
请分析这个 Spring Cloud Maven 项目，并以总分结构生成 AGENTS.md。
项目包含 gateway-service、user-service、order-service 和 common-core。
```

期望：

- 根目录生成 `AGENTS.md`。
- gateway/user/order 各自生成 `AGENTS.md`。
- common-core 不被当作微服务。
- 根文档引用三个服务文档。

### Test 2：已有 CLAUDE.md 的项目

输入场景：

```text
这个项目已经有 CLAUDE.md，请根据 Spring Cloud 微服务结构补全它，并为每个微服务创建对应的 CLAUDE.md。
```

期望：

- 使用 `CLAUDE.md`，不强行改为 `AGENTS.md`。
- 保留已有人工规则。
- 服务级文件也使用 `CLAUDE.md`。
- 不覆盖用户已有约束。

### Test 3：复杂多层 Maven 项目

输入场景：

```text
请给这个多层 Maven Spring Cloud 项目初始化 agent 上下文。
services 目录下是微服务，api 和 common 目录下是共享模块。
```

期望：

- 递归识别 `pom.xml`。
- 正确区分 services / api / common。
- 不把 api module 误判为可启动服务。
- 根文档包含 shared modules 表。
- 服务文档包含 Feign / Gateway / MQ 线索。

## 14. 风险与应对

| 风险 | 应对 |
|---|---|
| Agent 根据服务名脑补职责 | 要求 evidence-first，无证据省略并收敛到 Open Questions |
| 把 shared module 当微服务 | 使用启动类、服务名、端口、入口注解等多信号判断 |
| 覆盖已有人工规则 | 强制 read-before-write，用 managed 标记区隔离生成内容（Step 2） |
| 文档过长不可用 | 模板可裁剪，聚焦架构、入口、通信和修改风险，不列全量类 |
| 服务通信难以完整识别 | 标注 evidence / inferred / unknown |
| 多层 Maven 结构漏扫 | `scan_modules.py` 递归枚举，排除生成目录 |
| 大仓一次性生成几十个文件 | Step 6.5 计划确认检查点 + 广度优先深挖 |
| 大仓上下文/成本膨胀 | 浅扫描建清单 + 仅对确认微服务深挖 + 枚举脚本省 token |

## 15. 已确认决策

以下决策已评审采纳：

1. skill 名称使用 `springcloud-init`。
2. 文件选择策略：已有 `AGENTS.md` 优先，否则已有 `CLAUDE.md`，都没有则创建 `AGENTS.md`。
3. 第一版只支持 Maven，不支持 Gradle（列为后续增强）。
4. Git 分支规则只提取项目已有规则，不编写通用模板。
5. 服务级文档放在微服务 `pom.xml` 同级。
6. 不为 shared module 单独建文档，只在根文档表格列出并标注影响范围。
7. 第一版即引入“只枚举不判定”的 `scan_modules.py` 脚本，以省 token、保证一致性；判定逻辑仍由 Agent 完成。
