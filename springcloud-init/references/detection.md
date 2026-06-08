# Spring Cloud Detection Heuristics

Use these rules to classify Maven modules and identify Spring Cloud architecture evidence. Prefer multiple weak signals over one weak signal. Mark ambiguous cases as `Needs verification`.

## Module Classification

### Microservice Evidence

Strong evidence that a module can run independently:

- Java/Kotlin source contains `@SpringBootApplication`.
- Source contains `SpringApplication.run(...)`.
- Config contains `spring.application.name`.
- Config contains `server.port`.
- `pom.xml` depends on web/runtime starters:
  - `spring-boot-starter-web`
  - `spring-boot-starter-webflux`
  - `spring-cloud-starter-gateway`
- Module contains runtime entry points:
  - `@RestController`
  - `@Controller`
  - `@KafkaListener`
  - `@RabbitListener`
  - `@RocketMQMessageListener`
  - `@Scheduled`
- Module depends on discovery/config clients:
  - Nacos discovery/config
  - Eureka client
  - Consul discovery/config
  - Spring Cloud Config client

A module with a main class plus service identity (`spring.application.name`) is usually a microservice.

### Shared Module Evidence

A module is usually shared code when:

- `packaging` is `pom`.
- No Spring Boot main class is found.
- No service name or port is found.
- It mainly contains common, api, dto, model, entity, mapper, util, sdk, starter, or autoconfigure code.
- It is depended on by services but cannot run independently.

Do not create service-level guides for shared modules unless the user asks.

### Needs Verification

Use `Needs verification` when evidence conflicts, for example:

- Web dependency exists but no main class or application name is found.
- `spring.application.name` exists in a shared config module.
- A gateway route references a service not present in the repository.
- A module looks runnable but is not listed in root Maven modules.

## Spring Cloud Infrastructure Signals

| Component | Evidence |
|---|---|
| Nacos Discovery | `spring-cloud-starter-alibaba-nacos-discovery`, `spring.cloud.nacos.discovery` |
| Nacos Config | `spring-cloud-starter-alibaba-nacos-config`, `spring.cloud.nacos.config` |
| Eureka | `spring-cloud-starter-netflix-eureka-client`, `eureka.client` |
| Consul | `spring-cloud-starter-consul-discovery`, `spring.cloud.consul` |
| Spring Cloud Config | `spring-cloud-starter-config`, `spring.cloud.config` |
| Gateway | `spring-cloud-starter-gateway`, `spring.cloud.gateway.routes`, `RouteLocator` |
| OpenFeign | `spring-cloud-starter-openfeign`, `@FeignClient` |
| Sentinel | sentinel dependency, `spring.cloud.sentinel` |
| Seata | seata dependency, `seata.*` config |
| Tracing | Sleuth, Zipkin, Micrometer Tracing, OpenTelemetry dependencies/config |
| Redis | Redis starter, `spring.redis`, `spring.data.redis` |
| Database | datasource config, MyBatis/JPA dependencies, mapper/entity packages |

## Communication Detection

### Inbound

Evidence:

- `@RestController` / `@Controller`
- class-level and method-level request mappings
- Gateway routes targeting the service
- MQ listeners
- scheduled jobs if they trigger service-side work

### Outbound

Evidence:

- `@FeignClient(name = "...")`, `@FeignClient(value = "...")`
- `RestTemplate` calls
- `WebClient` calls
- MQ producer APIs or templates
- Dubbo references
- direct service names in config
- internal API module dependencies

### Gateway Routes

Look for:

- `spring.cloud.gateway.routes`
- `lb://service-name`
- `RouteLocator`
- route predicates and filters

Record route evidence as:

```text
Gateway route: path/predicate -> lb://target-service (file:line if available)
```

## Evidence Rules

- Include file paths for important claims.
- Prefer exact evidence: class name, annotation, config key, dependency artifact.
- Use `Inferred` only when the relationship is plausible but not directly proven.
- Put uncertainty in `Open Questions`; do not scatter `Unknown` through tables.
