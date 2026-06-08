---
name: springcloud-init
description: >
  Use this skill when the user asks to initialize, generate, or update AGENTS.md
  or CLAUDE.md for a Spring Cloud / Spring Boot multi-service Maven project.
  It analyzes pom.xml files, Spring Boot entry points, application.yml/bootstrap.yml,
  service discovery/configuration components, Feign/Gateway/MQ communication, and
  creates a root guide plus per-microservice guide files. Trigger for requests about
  Spring Cloud architecture docs, AI agent project context, microservice analysis,
  or 总分式 AGENTS.md/CLAUDE.md generation.
---

# Spring Cloud Agent Context Initializer

Create layered AI-agent context files for Spring Cloud Maven projects:

- root guide: one overall `AGENTS.md` or `CLAUDE.md`
- service guides: one file beside each independently runnable microservice `pom.xml`

The output is not a README. It is operational context for future coding agents: architecture, service boundaries, communication, Spring Cloud infrastructure, and code-change risks.

## File Choice

Choose the guide filename once, then use it consistently at root and service level:

1. If root `AGENTS.md` exists, use `AGENTS.md`.
2. Else if root `CLAUDE.md` exists, use `CLAUDE.md`.
3. Else create `AGENTS.md`.
4. If both exist, maintain `AGENTS.md` and preserve or reference useful rules from `CLAUDE.md`.

## Managed Content

Generated content must be inside a managed block:

```md
<!-- managed:springcloud-init -->
... generated content ...
<!-- /managed:springcloud-init -->
```

Update rules:

- If a target file has this block, replace only the block contents.
- If a target file exists without the block, preserve all existing content and append a new managed block.
- Never rewrite or delete content outside the managed block unless the user explicitly asks.

## Workflow

### 1. Inspect Before Writing

Read existing root `AGENTS.md` / `CLAUDE.md` if present. Identify human-written:

- code guidelines
- branch / commit rules
- build commands
- security or deployment warnings
- project-specific caveats

Preserve these outside managed blocks.

### 2. Build the Maven Module Map

Run the scanner from this skill directory when available:

```bash
python3 springcloud-init/scripts/scan_modules.py <project-root>
```

If the script path differs, locate it from the installed skill directory. The script only enumerates Maven modules; it does not decide what is a microservice.

If the script is unavailable, use direct commands:

```bash
find . -path './target' -prune -o -path './.git' -prune -o -path './node_modules' -prune -o -name pom.xml -print
```

Record each module's path, `artifactId`, `packaging`, parent, declared modules, internal dependencies, and Spring Boot / Spring Cloud dependencies.

### 3. Classify Services Breadth-First

First build the full module list. Then do lightweight signal scans for candidate modules only:

```bash
rg "@SpringBootApplication|SpringApplication\.run" --glob '!target/**'
rg "spring\.application\.name|server\.port|spring\.cloud" --glob '!target/**'
rg "@RestController|@Controller|@FeignClient|@KafkaListener|@RabbitListener|@Scheduled|RouteLocator|routes:" --glob '!target/**'
```

Classify each module as:

- `Microservice`: evidence it can run independently
- `Shared module`: reusable code only, not independently runnable
- `Needs verification`: conflicting or insufficient evidence

See [references/detection.md](references/detection.md) for detailed heuristics.

### 4. Deep-Analyze Confirmed Microservices

Only deep-scan confirmed microservice directories. Avoid expensive full scans of shared modules.

For each microservice, identify:

- main class and runtime identity
- `spring.application.name`, port, profiles
- controllers and inbound APIs
- Feign / HTTP / WebClient outbound calls
- Gateway routes
- MQ producers/listeners
- scheduled jobs
- internal module dependencies
- important config files and keys
- data access evidence

Use evidence-first writing: do not infer relationships or responsibilities without file evidence.

### 5. Plan-Validate-Execute Checkpoint

Before writing any files, output a plan and wait for user confirmation unless the user explicitly asked for autonomous execution.

The plan must include:

- selected filename: `AGENTS.md` or `CLAUDE.md`, with reason
- microservices list with evidence
- shared modules list with evidence
- `Needs verification` list
- files to create or update

If the user asks you to proceed without confirmation, still print the plan as an execution record before writing.

### 6. Write the Root Guide

Root guide content should cover:

- project overview
- architecture summary
- microservice table
- links to service guides
- shared modules and impact scope
- communication map
- Spring Cloud infrastructure
- build and run commands
- code guidelines for agents
- Git / branch rules only if found or provided
- evidence and open questions

Use the template principles in [references/templates.md](references/templates.md).

### 7. Write Service Guides

Create service guide files beside confirmed microservice `pom.xml` files only.

Each service guide should focus on:

- service overview
- runtime identity
- responsibilities
- inbound APIs
- outbound calls
- messaging
- scheduled jobs
- internal dependencies
- configuration
- data access
- local development
- agent notes
- evidence used
- open questions

Do not create service-level guides for pure shared modules unless the user explicitly asks.

### 8. Validate and Report

Copy this checklist into the final response and mark each item:

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

Also summarize created/updated files and any unresolved questions.

## Quality Rules

- Templates are sensible defaults, not mandatory forms.
- Omit sections without useful evidence; collect uncertainty in `Open Questions`.
- Do not fill tables with repeated `Unknown` values.
- Do not list every class. Include only architecture- and modification-relevant facts.
- Do not modify business code while generating context files.
- Do not invent Git, branch, build, deployment, or service communication rules.

## References

- Detection heuristics: [references/detection.md](references/detection.md)
- Templates and examples: [references/templates.md](references/templates.md)
