## Clarification

I inspected the current directory and found only:

- `README.md`

There is no application code, routing layer, API handlers, auth system, package manifest, or framework structure to modify.

Because no role model or web/API implementation exists yet, I did **not** make code changes in place. Implementing RBAC would require knowing the stack and existing auth/session model.

## Assumptions / Defaults I would use

If you want me to implement this once the app exists, I would default to:

- Roles: `admin`, `user`
- Auth source: existing session/JWT user object
- Admin-only UI routes: `/admin/*`
- Normal user routes: `/app/*` or `/dashboard`
- API protection:
  - Admin APIs require `role === "admin"`
  - User APIs require authenticated users
  - Shared APIs check ownership where applicable
- Unauthorized behavior:
  - UI: redirect to login or unauthorized page
  - API: return `401` if unauthenticated, `403` if authenticated but forbidden

## Routing Recommendation

A typical RBAC layout would be:

```txt
/login
/unauthorized

/app              -> authenticated users
/app/profile      -> authenticated users

/admin            -> admins only
/admin/users      -> admins only
/admin/settings   -> admins only
```

API policy:

```txt
GET /api/me              -> authenticated
GET /api/admin/users     -> admin only
POST /api/admin/users    -> admin only
GET /api/user/data       -> authenticated user, ownership checked
```

## Implementation Summary

No implementation was applied because the repository currently contains only a placeholder README and no runnable application.

Recommended next steps:

1. Add or provide the actual app source.
2. Identify the framework, e.g. Next.js, Express, FastAPI, Rails, etc.
3. Confirm how users are authenticated.
4. Confirm role storage, e.g. database column, JWT claim, session field.
5. Then add:
   - Central role-checking middleware/helper
   - Protected UI routes/pages
   - Protected API handlers
   - Unauthorized page or redirect behavior
   - Tests for `401`, `403`, and allowed access

## Validation

I validated the workspace with:

```bash
ls
find . -maxdepth 3 -type f -print
```

Result: only `README.md` exists, so there was no code path to safely modify.
