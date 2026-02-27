import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { useAuth } from './AuthContext';
import tenantService from '../services/tenantService';

const TenantContext = createContext();

export const useTenant = () => {
  const context = useContext(TenantContext);
  if (!context) {
    throw new Error('useTenant must be used within a TenantProvider');
  }
  return context;
};

// Decode JWT payload client-side (no signature verification — backend already validated it).
function decodeJWTPayload(token) {
  try {
    const base64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(atob(base64));
  } catch {
    return null;
  }
}

export const TenantProvider = ({ children }) => {
  const { authState, isInitialized } = useAuth();

  const [currentTenant, setCurrentTenant] = useState(null);
  const [myRole, setMyRole] = useState(null);
  const [permissions, setPermissions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (authState === 'authenticated' && isInitialized) {
      loadTenantData();
    } else if (authState === 'unauthenticated') {
      setCurrentTenant(null);
      setMyRole(null);
      setPermissions([]);
      setError(null);
      setLoading(false);
    }
  }, [authState, isInitialized]);

  const loadTenantData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Decode role from JWT — no extra network request needed.
      const token = localStorage.getItem('auth_token');
      if (token) {
        const payload = decodeJWTPayload(token);
        if (payload?.role) {
          setMyRole(payload.role);
        }
      }

      const [tenant, perms] = await Promise.all([
        tenantService.getMyTenant().catch(err => {
          console.warn('[TenantContext] getMyTenant failed:', err.message);
          return null;
        }),
        tenantService.getMyPermissions().catch(err => {
          console.warn('[TenantContext] getMyPermissions failed:', err.message);
          return [];
        }),
      ]);

      setCurrentTenant(tenant);
      setPermissions(perms || []);
    } catch (err) {
      console.error('[TenantContext] loadTenantData error:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const hasPermission = useCallback(
    (perm) => permissions.includes(perm),
    [permissions],
  );

  const isOwner = myRole === 'owner';
  const isAdmin = myRole === 'admin' || isOwner;

  const value = {
    currentTenant,
    myRole,
    permissions,
    hasPermission,
    isOwner,
    isAdmin,
    loading,
    error,
    refreshTenant: loadTenantData,
  };

  return (
    <TenantContext.Provider value={value}>
      {children}
    </TenantContext.Provider>
  );
};

export default TenantContext;
