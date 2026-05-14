## Clarification

I inspected the current directory and found only a placeholder app:

- `README.md`
- No source files, routing layer, API handlers, auth/session code, package manifest, or framework configuration.

Because no role model, framework, routes, or APIs exist yet, there is no safe implementation target for RBAC in this workspace.

## Assumptions / Defaults I would use once the app exists

If you want me to implement RBAC after the actual app code is present, I would default to:

- Roles: `admin` and `user`
- Default access policy: deny by default
- Admin pages: accessible only to `admin`
- User pages: accessible to both `user` and possibly `admin`, depending on desired behavior
- API protection: enforced server-side, never only in the UI
- Unauthorized behavior:
  - Unauthenticated users: redirect to login or return `401`
  - Authenticated but insufficient role: show forbidden page or return `403`

## Routing Recommendation

A common structure would be:

```text
/admin/*        admin only
/app/*          normal authenticated users
/api/admin/*    admin only
/api/user/*     authenticated users
```

RBAC checks should happen in:

1. Route middleware / loaders / guards for page access
2. API middleware or per-handler authorization checks
3. Optional UI filtering only as a convenience, not as security

## Implementation Summary

No code changes were made because implementation is not appropriate without an actual application codebase or framework to modify.

## Validation

Validated the workspace contents with:

```bash
ls
find . -maxdepth 3 -type f -print
```

Result: only `README.md` is present, and it states:

> A placeholder web application. No role model has been defined yet.
