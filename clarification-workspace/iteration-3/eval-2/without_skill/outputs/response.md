Clarification: I checked the current directory and found only a `README.md`; there is no web application source code, routing layer, API handlers, auth/session code, or role model to modify.

Assumptions/defaults I documented:
- Roles: `admin` and `user`
- Server-side enforcement required for both pages and APIs
- Deny-by-default for missing/unknown roles
- `401` for unauthenticated requests
- `403` for authenticated users without the required role
- UI hiding is not sufficient; APIs must enforce authorization independently

Implementation/routing recommendation:
- Since implementation is not appropriate without an actual app/framework, I added an RBAC recommendation to `README.md`.
- Recommended route split:
  - `/admin/**`: admin only
  - normal app/dashboard pages: admin and user
  - `/api/admin/**`: admin only
  - `/api/user/**`: admin and user
  - public pages: open

Validation:
- Inspected repository contents with `ls` and `find`.
- Confirmed only `README.md` exists.
- Updated `README.md` in place with concrete RBAC defaults and policy guidance.
