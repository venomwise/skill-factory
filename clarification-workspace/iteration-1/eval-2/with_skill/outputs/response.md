## Clarification

I inspected the current directory and found only:

- `README.md`

The README says:

> “A placeholder web application. No role model has been defined yet.”

There is currently no application code, routing layer, auth/session implementation, API layer, framework config, or role model to modify.

## Assumptions/defaults I would use once implementation is possible

If you want me to proceed after the app exists or after you confirm a stack, I would default to:

- Roles: `admin` and `user`
- Default role: `user`
- Access behavior:
  - Admins can access admin pages and normal user pages unless told otherwise
  - Normal users cannot access admin pages
- API protection:
  - Require authentication for protected APIs
  - Enforce role checks server-side, not only in UI routing
  - Return `401` for unauthenticated requests
  - Return `403` for authenticated users without the required role

## Implementation summary / routing recommendation

I did **not** implement RBAC because there is no application code in the current directory to safely modify.

Recommended next step: provide or create the actual app scaffold first, or confirm the desired framework, for example:

- Next.js / React
- Express
- FastAPI
- Django
- Rails
- Laravel
- Other

Once the stack is known, I can add:

- User role model or role field
- Route/page guards
- API authorization middleware
- Admin/user page separation
- Tests or validation checks

## Validation

I validated the workspace contents with:

```bash
ls
find . -maxdepth 2 -type f
```

Result: only `README.md` exists, so implementation is currently not appropriate.
