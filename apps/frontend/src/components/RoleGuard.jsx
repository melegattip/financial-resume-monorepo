import React from 'react';
import { useTenant } from '../contexts/TenantContext';

/**
 * RoleGuard renders `children` only when the current user has the required
 * permission (or one of the required permissions when `any` is true).
 *
 * Props:
 *   permission  {string}           – single permission key required
 *   permissions {string[]}         – multiple permission keys (requires all by default)
 *   any         {boolean}          – when true, access is granted if the user has ANY of the listed permissions
 *   role        {string}           – alternatively guard by exact role ('owner', 'admin', …)
 *   fallback    {React.ReactNode}  – rendered when access is denied (default: null)
 */
const RoleGuard = ({ permission, permissions, any = false, role, fallback = null, children }) => {
  const { hasPermission, myRole, loading } = useTenant();

  if (loading) return null;

  // Role-based guard
  if (role) {
    if (myRole === role) return children;
    return fallback;
  }

  // Permission-based guard (multiple)
  if (permissions && permissions.length > 0) {
    const granted = any
      ? permissions.some(p => hasPermission(p))
      : permissions.every(p => hasPermission(p));
    return granted ? children : fallback;
  }

  // Permission-based guard (single)
  if (permission) {
    return hasPermission(permission) ? children : fallback;
  }

  // No guard condition specified — render children
  return children;
};

export default RoleGuard;
