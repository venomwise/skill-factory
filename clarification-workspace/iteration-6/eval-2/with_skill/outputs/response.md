**My understanding**: Add role-based access control with admin/user roles, different page visibility per role, and API endpoint protection.

**Why this is larger than a small clarification**:

This is a **permission model**, not a localized tweak. It involves:

- **Security architecture**: Where roles are stored, how they're verified, session/token handling
- **Multiple subsystems**: Frontend routing guards, backend middleware/decorators, database schema
- **Authorization boundaries**: Which pages/APIs need protection, what each role can access
- **Migration strategy**: How existing users get roles, backward compatibility
- **Error handling**: Redirect flows, 403 responses, unauthorized state management
- **Testing**: Security test coverage, role switching, edge cases

Before I can implement this safely, I need to understand:
- Your current auth system (if any)
- Tech stack (framework, database, frontend library)
- Existing user model
- Which specific pages and APIs need protection
- Role assignment workflow

**Recommendation**: Use the `brainstorming` skill to design the RBAC system properly. This will help us define:
- Role model and storage
- Authentication/authorization flow
- Frontend and backend enforcement strategy
- Migration plan
- Security testing approach

Would you like me to start a brainstorming session to design this RBAC system?
