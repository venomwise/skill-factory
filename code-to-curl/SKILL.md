---
name: code-to-curl
description: >
  Convert HTTP interface/controller code to `curl` command and standard HTTP message format
  (raw HTTP/1.1 request) examples. Use when given source code (Spring MVC, FastAPI, Express,
  NestJS, Gin, etc.) to generate ready-to-run API test commands. Triggered by requests like
  "生成 curl 请求", "转成 http 格式", "给我这个接口的请求示例", or "how do I call this endpoint".
---

# Code to curl / HTTP message

Read HTTP endpoint definitions from source code and output a `curl` command and a standard HTTP/1.1 request message for every endpoint in scope.

## When to use

- User provides a controller, handler, or route file and wants request examples
- User pastes HTTP interface code and asks for curl or http commands
- User wants to quickly test or document endpoints without writing commands manually

## When not to use

- No source code is provided — ask the user to share the file path or paste the code
- User wants full API documentation (use `hld-generator` or an OpenAPI tool)
- User has curl commands and wants to convert them to code (reverse direction — do it directly, no skill needed)

## Inputs

- Source file path(s) or pasted code with HTTP endpoint definitions
- Optional: base URL override (defaults to framework-typical local default)
- Optional: list of specific endpoints to focus on (defaults to all endpoints in the provided scope)

## Outputs

One section per endpoint containing a `curl` command and a standard HTTP/1.1 request message.

## Workflow

### 1. Locate the code

- File path given → read the file.
- Pasted code → analyze directly.
- Class name / feature name given → search for the matching file first.

### 2. Identify the framework

Detect from imports, annotations, or decorators:

| Framework | Signals |
|-----------|---------|
| Spring MVC | `@RestController`, `@GetMapping`, `@PostMapping`, `@PutMapping`, `@DeleteMapping`, `@PatchMapping`, `@RequestMapping` |
| FastAPI | `@app.get`, `@app.post`, `@router.get`, `@router.post` (Python, type hints) |
| Flask | `@app.route`, `@bp.route` with `methods=[...]` |
| Django REST | `@api_view`, `APIView`, `ViewSet` with `list`/`create`/`retrieve` methods |
| NestJS | `@Controller`, `@Get()`, `@Post()`, `@Put()`, `@Delete()`, `@Patch()` |
| Express / Koa | `router.get(`, `router.post(`, `app.get(`, `app.post(` |
| Gin (Go) | `r.GET(`, `r.POST(`, `r.PUT(`, group routes via `r.Group(` |
| net/http (Go) | `http.HandleFunc`, `mux.Handle`, `mux.HandleFunc` |
| Laravel | `Route::get(`, `Route::post(`, `Route::resource(` |

### 3. Extract per endpoint

For each endpoint collect:

- **HTTP method** — GET / POST / PUT / DELETE / PATCH
- **Full URL path** — combine class-level prefix + method-level path; keep `{variable}` or `:variable` placeholders
- **Path variables** — names from the path template
- **Query parameters** — from `@RequestParam`, `Query`, `query.Get(`, `request.args`, etc.
- **Request body** — from `@RequestBody`, Pydantic model param, `req.body`, `c.ShouldBindJSON`, etc.  
  Identify content-type: JSON (default), form (`application/x-www-form-urlencoded`), or multipart (`multipart/form-data`)
- **Headers** — any explicitly required headers; note auth requirements from `@PreAuthorize`, security config, `AuthGuard`, JWT middleware, etc.

If the request body type is a DTO / Pydantic model / interface defined in another file, **read that file** to get the full field list before generating the body.

### 4. Determine base URL

Use the user-provided override when given. Otherwise:

| Framework | Default |
|-----------|---------|
| Spring Boot | `http://localhost:8080` |
| FastAPI | `http://localhost:8000` |
| Flask | `http://localhost:5000` |
| Django REST | `http://localhost:8000` |
| NestJS / Express / Koa | `http://localhost:3000` |
| Gin / net/http (Go) | `http://localhost:8080` |
| Laravel | `http://localhost:8000` |

### 5. Generate realistic placeholder values

| Field type / name pattern | Example value |
|--------------------------|---------------|
| `id`, `userId`, `*Id` | `1` |
| `name`, `username` | `"alice"` |
| `email` | `"user@example.com"` |
| `password` | `"P@ssw0rd"` |
| `phone` | `"13800138000"` |
| `date` / `time` | `"2024-01-01"` / `"12:00:00"` |
| `status` | `"active"` |
| `boolean` | `true` |
| `number` / `amount` | `100` |
| `page` / `pageSize` | `1` / `10` |
| Unknown string | field name as value |

Path variable `{id}` → substitute `1`. Path variable `{slug}` → substitute `"example-slug"`.

### 6. Output format

Output one section per endpoint:

````
## POST /api/users

**功能**（可从方法/参数名推断时填写）

**curl:**
```bash
curl -X POST "http://localhost:8080/api/users" \
  -H "Content-Type: application/json" \
  -d '{
  "username": "alice",
  "email": "user@example.com",
  "age": 25
}'
```

**HTTP:**
```http
POST /api/users HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "username": "alice",
  "email": "user@example.com",
  "age": 25
}
```
````

**curl conventions:**
- Always `-X METHOD` and quote the URL
- `\` for line continuation
- JSON body: `-H "Content-Type: application/json" -d '{ ... }'`
- Form body: `-H "Content-Type: application/x-www-form-urlencoded" --data-urlencode "key=value"`
- Multipart: `-F "file=@/path/to/file" -F "field=value"`
- Query params: append to URL string (`?key=value&key2=value2`)
- Auth placeholder: `-H "Authorization: Bearer <token>"`

**HTTP message conventions (RFC 7230):**
- Request line: `METHOD /path?query HTTP/1.1`
- `Host` header is always required (hostname + port if non-standard)
- One header per line: `Header-Name: value`
- Blank line between headers and body
- Body only for POST / PUT / PATCH; omit for GET / DELETE
- JSON body: include `Content-Type: application/json`
- Form body: `Content-Type: application/x-www-form-urlencoded`, body as `key=value&key2=value2`
- Auth placeholder: `Authorization: Bearer <token>`
- Omit `Content-Length` (it's optional for readability)

### 7. After all endpoints

State any assumptions made, e.g.:
- Default port used
- Auth requirement inferred from security config
- DTO fields read from a secondary file

## Verification

- [ ] All endpoints in the specified scope are covered
- [ ] curl command is syntactically valid — correct quoting, well-formed JSON body
- [ ] HTTP message has correct request line, `Host` header, blank separator line, and body
- [ ] Path variables are substituted with realistic values
- [ ] Request body matches the DTO / model field structure (read the type definition file if needed)
- [ ] Content-Type matches the actual body format (JSON / form / multipart)

## Safety & guardrails

- Use `<token>` or `<api-key>` as auth placeholders — do not invent real credentials.
- Output only — do not execute any generated command.
- If a class-level `@RequestMapping` prefix exists, always prepend it to all method-level paths.
- When the code spans multiple files (base controller + subclass, router modules), read all relevant files before generating.
