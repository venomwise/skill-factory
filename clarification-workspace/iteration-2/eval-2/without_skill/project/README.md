# RBAC fixture

A placeholder web application. No concrete framework, authentication provider, or persisted role model has been defined yet, so this repository now contains a small framework-agnostic RBAC helper that can be wired into the eventual app.

## Default role model

- `admin`: can access admin pages/APIs and normal user pages/APIs.
- `user`: can access only normal user pages/APIs.
- Unknown or missing roles are denied.
- Protected routes are deny-by-default unless they have an explicit access rule.

## Added files

- `src/rbac.js` — role constants, page/API access rules, authorization helpers, and Express/Connect-compatible middleware.

## Page routing recommendation

Use separate routes/layouts for normal and admin experiences:

- Normal user pages: `/dashboard`, `/account`
- Admin pages: `/admin`, `/admin/*`

At page-render or route-guard time, call:

```js
const { authorizePage } = require('./src/rbac');

if (!authorizePage(request.path, request.user)) {
  // redirect unauthenticated users to login; show 403 for authenticated users
}
```

## API protection recommendation

Apply RBAC after authentication has populated `req.user`, for example:

```js
const { API_ACCESS, requireRouteAccess } = require('./src/rbac');

app.use(authenticateRequest);
app.use('/api', requireRouteAccess(API_ACCESS));
```

Or protect individual endpoints:

```js
const { ROLES, requireRole } = require('./src/rbac');

app.get('/api/admin/reports', requireRole([ROLES.ADMIN]), handler);
app.get('/api/me', requireRole([ROLES.USER, ROLES.ADMIN]), handler);
```

## Current default route rules

Pages:

- `/`, `/dashboard`, `/account`: `user`, `admin`
- `/admin`, `/admin/*`: `admin`

APIs:

- `/api/me`, `/api/user/*`: `user`, `admin`
- `/api/admin`, `/api/admin/*`: `admin`
