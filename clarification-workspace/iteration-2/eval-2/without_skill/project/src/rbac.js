'use strict';

/**
 * Small, framework-agnostic RBAC helpers.
 *
 * Default model:
 * - admin: can access admin and user areas
 * - user: can access only normal user areas
 *
 * These helpers can be used from server-rendered page routes, SPA route guards,
 * and Express/Connect-compatible API middleware.
 */

const ROLES = Object.freeze({
  ADMIN: 'admin',
  USER: 'user',
});

const ROLE_HIERARCHY = Object.freeze({
  [ROLES.USER]: 1,
  [ROLES.ADMIN]: 2,
});

const PAGE_ACCESS = Object.freeze([
  { path: '/', roles: [ROLES.USER, ROLES.ADMIN] },
  { path: '/dashboard', roles: [ROLES.USER, ROLES.ADMIN] },
  { path: '/account', roles: [ROLES.USER, ROLES.ADMIN] },
  { path: '/admin', roles: [ROLES.ADMIN] },
  { path: '/admin/*', roles: [ROLES.ADMIN] },
]);

const API_ACCESS = Object.freeze([
  { method: '*', path: '/api/me', roles: [ROLES.USER, ROLES.ADMIN] },
  { method: '*', path: '/api/user/*', roles: [ROLES.USER, ROLES.ADMIN] },
  { method: '*', path: '/api/admin', roles: [ROLES.ADMIN] },
  { method: '*', path: '/api/admin/*', roles: [ROLES.ADMIN] },
]);

function normalizePath(path) {
  if (!path) return '/';
  const [pathname] = String(path).split('?');
  if (pathname.length > 1 && pathname.endsWith('/')) return pathname.slice(0, -1);
  return pathname || '/';
}

function pathMatches(pattern, path) {
  const normalizedPattern = normalizePath(pattern);
  const normalizedPath = normalizePath(path);

  if (normalizedPattern.endsWith('/*')) {
    const prefix = normalizedPattern.slice(0, -2);
    return normalizedPath === prefix || normalizedPath.startsWith(`${prefix}/`);
  }

  return normalizedPattern === normalizedPath;
}

function getUserRole(user) {
  return user && typeof user.role === 'string' ? user.role : undefined;
}

function canAccess(userOrRole, allowedRoles) {
  const role = typeof userOrRole === 'string' ? userOrRole : getUserRole(userOrRole);
  return Boolean(role && Array.isArray(allowedRoles) && allowedRoles.includes(role));
}

function hasMinimumRole(userOrRole, minimumRole) {
  const role = typeof userOrRole === 'string' ? userOrRole : getUserRole(userOrRole);
  return Boolean(role && ROLE_HIERARCHY[role] >= ROLE_HIERARCHY[minimumRole]);
}

function findAccessRule(rules, path, method = '*') {
  const upperMethod = String(method || '*').toUpperCase();
  return rules.find((rule) => {
    const ruleMethod = String(rule.method || '*').toUpperCase();
    return (ruleMethod === '*' || ruleMethod === upperMethod) && pathMatches(rule.path, path);
  });
}

function requiredRolesForPath(rules, path, method = '*') {
  const rule = findAccessRule(rules, path, method);
  return rule ? rule.roles : undefined;
}

function authorizePath(rules, path, user, method = '*') {
  const roles = requiredRolesForPath(rules, path, method);

  // Deny by default when a protected route has no explicit rule.
  if (!roles) return false;

  return canAccess(user, roles);
}

function authorizePage(path, user) {
  return authorizePath(PAGE_ACCESS, path, user);
}

function authorizeApi(path, user, method = '*') {
  return authorizePath(API_ACCESS, path, user, method);
}

function requireRole(allowedRoles) {
  return function rbacMiddleware(req, res, next) {
    if (canAccess(req.user, allowedRoles)) return next();

    const status = req.user ? 403 : 401;
    return res.status(status).json({ error: status === 401 ? 'unauthorized' : 'forbidden' });
  };
}

function requireRouteAccess(rules) {
  return function rbacRouteMiddleware(req, res, next) {
    const path = req.path || req.url;
    const method = req.method || '*';

    if (authorizePath(rules, path, req.user, method)) return next();

    const status = req.user ? 403 : 401;
    return res.status(status).json({ error: status === 401 ? 'unauthorized' : 'forbidden' });
  };
}

module.exports = {
  ROLES,
  ROLE_HIERARCHY,
  PAGE_ACCESS,
  API_ACCESS,
  normalizePath,
  pathMatches,
  getUserRole,
  canAccess,
  hasMinimumRole,
  findAccessRule,
  requiredRolesForPath,
  authorizePath,
  authorizePage,
  authorizeApi,
  requireRole,
  requireRouteAccess,
};
