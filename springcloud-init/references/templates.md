# Spring Cloud Agent Guide Templates

Templates are sensible defaults, not mandatory forms.

Rules:

- Omit sections with no useful evidence.
- Do not fill tables with repeated `Unknown`.
- Collect uncertainty in one `Open Questions` section.
- Keep generated content inside the managed block.
- Prefer concise evidence-backed facts over exhaustive class lists.

## Root Guide Template

```md
# Project Agent Guide

<!-- managed:springcloud-init -->
## Project Overview

Briefly describe the Spring Cloud / Spring Boot system and business domain.

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

| Module | Path | Purpose | Used By / Impact |
|---|---|---|---|

## Communication Map

| Source | Target | Mechanism | Evidence |
|---|---|---|---|

## Spring Cloud Infrastructure

## Build and Run

## Code Guidelines for Agents

- Modify Feign/API contracts only after checking both callers and providers.
- When changing DTO/VO/API modules, check all dependent services.
- When changing config keys, check `bootstrap.*`, `application.*`, profiles, and config-center usage.
- When changing entities or SQL-facing code, check mappers, XML, migrations, and SQL scripts.
- Keep business logic out of controllers unless the project already uses that pattern.
- When changing shared modules, check all services that depend on them.

## Git and Branch Rules

Only include rules found in project files or explicitly provided by the user. Do not invent branch conventions.

## Evidence Used

## Open Questions
<!-- /managed:springcloud-init -->
```

## Service Guide Template

```md
# <service-name> Agent Guide

<!-- managed:springcloud-init -->
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
<!-- /managed:springcloud-init -->
```

## Positive Mini Example

This is the desired level of detail: concise, evidence-backed, and focused on change risks.

```md
# user-service Agent Guide

<!-- managed:springcloud-init -->
## Service Overview

`user-service` manages user profiles and exposes user lookup APIs. Evidence: `UserController` exposes `/users`, and `pom.xml` artifactId is `user-service`.

## Runtime Identity

| Item | Value |
|---|---|
| Artifact ID | `user-service` |
| Spring Application Name | `user-service` |
| Main Class | `com.example.user.UserApplication` |
| Port | `8081` |

## Inbound APIs

| Entry | Path | Purpose | Evidence |
|---|---|---|---|
| `UserController` | `/users` | User lookup and profile APIs | `src/main/java/.../UserController.java` |

## Outbound Calls

| Target | Mechanism | Purpose | Evidence |
|---|---|---|---|
| `order-service` | OpenFeign | Query user order summary | `OrderClient.java` has `@FeignClient("order-service")` |

## Internal Dependencies

- `common-core`: shared response/error types (`pom.xml` dependency)
- `user-api`: user DTO/API contracts (`pom.xml` dependency)

## Configuration

- `src/main/resources/bootstrap.yml`: `spring.application.name=user-service`
- `src/main/resources/application-dev.yml`: `server.port=8081`, datasource config

## Agent Notes

- If changing user DTOs, check `user-api` and Feign consumers.
- If changing `/users/**` endpoints, check gateway routes and downstream callers.

## Evidence Used

- `user-service/pom.xml`: artifactId and dependencies
- `UserApplication.java`: `@SpringBootApplication`
- `bootstrap.yml`: `spring.application.name=user-service`
- `application-dev.yml`: `server.port=8081`
- `OrderClient.java`: outbound Feign call to `order-service`

## Open Questions

- No gateway route file was found in this repository; inbound gateway mapping may live in config center.
<!-- /managed:springcloud-init -->
```

## Anti-Example: Bloated Class List

Do not produce this style. It spends tokens on inventory instead of architecture and change guidance.

```md
# user-service Agent Guide

## Classes

- UserApplication
- UserController
- UserService
- UserServiceImpl
- UserMapper
- UserEntity
- UserDTO
- UserVO
- UserRequest
- UserResponse
- UserConfig
- UserConstant
- UserException
- UserConverter
- UserRepository
...

## Unknowns

Port: Unknown
Messaging: Unknown
Scheduling: Unknown
Database: Unknown
Gateway: Unknown
```

Why this is bad:

- It lists classes without explaining service responsibility or boundaries.
- It scatters `Unknown` instead of using `Open Questions`.
- It lacks evidence paths and code-change warnings.
- It does not help an agent decide what to inspect before editing.
