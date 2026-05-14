Clarification: You asked for RBAC so admins and normal users see different pages and APIs are protected.

Assumptions/defaults: I inspected the current directory and found only:

- `README.md`
- No application source files, framework config, routes, API handlers, auth/session code, package files, or existing role model.

Implementation summary / routing recommendation: I did **not** implement code changes because there is no actual web application to modify, and the README explicitly says no role model has been defined yet. Adding RBAC safely requires confirming at least:

- Framework/runtime, e.g. Next.js, Express, Django, Rails, etc.
- Existing authentication/session mechanism
- Role source, e.g. `user.role`, JWT claim, database relation
- Role names and permissions
- Which pages/routes are admin-only vs user-only/shared
- API protection behavior, e.g. 401 unauthenticated, 403 unauthorized

Recommended next step: define a short RBAC design/spec first, then implement once the app structure or target framework exists.

Validation: Verified current project contents with filesystem inspection and read `README.md`; there are no code files available to change.
