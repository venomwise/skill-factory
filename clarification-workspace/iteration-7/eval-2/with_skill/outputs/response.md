My understanding: Add role-based access control with admin and normal user roles, differentiate page visibility, and protect APIs with authorization checks.

This looks larger than a small clarification because:

- **Permission model design**: Defining roles, where they're stored (database schema?), how they're assigned, and what each role can access requires architectural decisions
- **Multiple enforcement points**: Frontend routing, backend API middleware, and potentially database-level access all need coordinated protection
- **Security implications**: Authorization bugs can expose sensitive data or actions; the approach needs careful design (JWT claims? session-based? middleware order? default-deny vs default-allow?)
- **Cross-cutting changes**: Affects authentication flow, database schema, API layer, UI routing, and potentially existing user management
- **Missing context**: No existing codebase to inspect, so I can't determine current auth mechanisms, framework, or integration points

Key blind spots that would affect implementation:
- How are users currently authenticated? (sessions, JWT, OAuth, none?)
- What framework/stack? (React+Express? Django? Next.js? Go?)
- Role storage and assignment mechanism
- Specific pages and API endpoints to protect
- Behavior on unauthorized access (redirect to login? 403 error? different content?)
- Role hierarchy (just two roles, or extensible?)
- Migration path for existing users

**I recommend using the `brainstorming` skill to define the design/spec before implementation.** This will help you think through the authorization model, enforcement strategy, and edge cases before writing code that touches security-critical paths.

Would you like to start a brainstorming session to design the RBAC system properly?
