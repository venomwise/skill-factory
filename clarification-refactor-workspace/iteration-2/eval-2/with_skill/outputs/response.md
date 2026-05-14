This is larger than a small clarification because it involves:

- **Authorization and security**: Access control is a security boundary that affects the entire application
- **Multiple subsystems**: Frontend (different pages per role) and backend (API protection) need coordinated enforcement
- **Missing design decisions**: Role model, permission matrix, enforcement points, session/token handling, and fallback behavior are undefined
- **Cross-cutting concern**: Every protected page and API endpoint needs consistent authorization logic

I recommend using `brainstorming` before implementation.

**Handoff context:**

- **Current goal**: Add role-based access control with admin/user roles, different page visibility, and API protection

- **Why this exceeds a small change**:
  - Security boundaries require consistent enforcement across frontend and backend
  - Role model undefined (how are roles assigned? stored? verified?)
  - Authorization strategy unclear (middleware? decorators? guards? per-route checks?)
  - Session/authentication integration unknown
  - Error handling for unauthorized access needs definition (redirect? 403? 401?)
  - Frontend and backend must stay synchronized

- **Key unknowns / decisions**:
  - What roles exist beyond admin/user? Are roles hierarchical?
  - How are roles assigned and persisted (database? JWT claims? session?)
  - Where is authentication currently handled?
  - What happens when unauthorized users try to access protected resources?
  - Should some pages/APIs be partially accessible with degraded functionality?
  - Do we need audit logging for authorization failures?

- **Known project context**: None inspected yet. Need to understand current tech stack, authentication system, and project structure.

- **Suggested first brainstorming question**: "What authentication system is currently in place, and how should roles be stored and verified?"
