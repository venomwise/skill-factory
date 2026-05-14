I recommend using `brainstorming` before implementation.

**Handoff context:**

- **Current goal:** Add role-based access control with admin/user differentiation for both page visibility and API protection

- **Why this exceeds a small change:**
  - Security and authorization are cross-cutting concerns affecting multiple subsystems
  - Requires coordinated frontend (page routing/visibility) and backend (API protection) behavior
  - Role model needs definition: how roles are assigned, stored, and verified
  - Enforcement points need design: middleware, route guards, component-level checks
  - API protection strategy: token-based, session-based, or other authentication mechanism
  - Missing project structure means implementation would require inventing architecture

- **Key unknowns / decisions:**
  - Role model: what roles exist beyond admin/user? How are they assigned and persisted?
  - Authentication mechanism: existing auth system or new implementation needed?
  - Frontend enforcement: route guards, conditional rendering, or redirect strategy?
  - Backend enforcement: middleware layer, decorator pattern, or per-endpoint checks?
  - Permission granularity: page-level only, or feature/action-level within pages?
  - API contract: how should protected APIs respond to unauthorized requests (401, 403, redirect)?

- **Known project context:** Empty project directory - no existing codebase to inspect

- **Suggested first brainstorming question:** What authentication system (if any) already exists, and what roles/permissions model fits your application's user types and access patterns?
