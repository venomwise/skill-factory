Clarification: this repo only contained a `README.md` placeholder and no actual web framework, router, auth/session code, or API handlers. Because there was nothing concrete to wire into, I added a framework-agnostic RBAC module plus usage/routing guidance instead of inventing an app structure.

Assumptions/defaults:
- Roles are `admin` and `user`.
- `admin` can access admin and normal user areas.
- `user` can access only normal user areas.
- Missing/unknown roles are denied.
- Protected routes are deny-by-default unless explicitly listed.
- Auth is expected to populate `req.user.role` before RBAC runs.

Implementation summary:
- Added `src/rbac.js`
  - Role constants: `ROLES.ADMIN`, `ROLES.USER`
  - Page access rules:
    - `/`, `/dashboard`, `/account`: `user`, `admin`
    - `/admin`, `/admin/*`: `admin`
  - API access rules:
    - `/api/me`, `/api/user/*`: `user`, `admin`
    - `/api/admin`, `/api/admin/*`: `admin`
  - Helpers:
    - `authorizePage(path, user)`
    - `authorizeApi(path, user, method)`
    - `canAccess(...)`
    - `hasMinimumRole(...)`
    - `requireRole(...)`
    - `requireRouteAccess(...)`
- Updated `README.md` with integration examples for pages and APIs.

Routing recommendation:
- Keep normal user pages under routes like `/dashboard` and `/account`.
- Keep admin pages under `/admin` and `/admin/*`.
- Run authentication first, then RBAC.
- For APIs, apply `requireRouteAccess(API_ACCESS)` after authentication or protect endpoints individually with `requireRole([...])`.

Validation:
- Ran a Node smoke test validating:
  - Admin can access `/admin`.
  - Normal user cannot access `/admin`.
  - Normal user can access `/dashboard`.
  - Admin can access `/api/admin/reports`.
  - Normal user cannot access `/api/admin/reports`.
  - Unknown API routes are denied by default.

Validation passed.
