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
  const [availableTenants, setAvailableTenants] = useState([]);
  const [switching, setSwitching] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (authState === 'authenticated' && isInitialized) {
      loadTenantData();
    } else if (authState === 'unauthenticated') {
      setCurrentTenant(null);
      setMyRole(null);
      setPermissions([]);
      setAvailableTenants([]);
      setError(null);
      setLoading(false);
    }
  }, [authState, isInitialized]);

  const loadTenantData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Decode role and current tenant_id from JWT — no extra network request needed.
      const token = localStorage.getItem('auth_token');
      let currentTenantId = null;
      if (token) {
        const payload = decodeJWTPayload(token);
        if (payload?.role) setMyRole(payload.role);
        if (payload?.tenant_id) currentTenantId = payload.tenant_id;
      }

      const [tenant, perms, allTenants] = await Promise.all([
        tenantService.getMyTenant().catch(err => {
          console.warn('[TenantContext] getMyTenant failed:', err.message);
          return null;
        }),
        tenantService.getMyPermissions().catch(err => {
          console.warn('[TenantContext] getMyPermissions failed:', err.message);
          return [];
        }),
        tenantService.listMyTenants().catch(err => {
          console.warn('[TenantContext] listMyTenants failed:', err.message);
          return [];
        }),
      ]);

      setCurrentTenant(tenant);
      setPermissions(perms || []);
      setAvailableTenants(allTenants || []);
    } catch (err) {
      console.error('[TenantContext] loadTenantData error:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Switch the active tenant without logging out.
  // Stores the new tokens and reloads tenant data.
  const switchTenant = useCallback(async (tenantId) => {
    try {
      setSwitching(true);
      const tokens = await tenantService.switchTenant(tenantId);

      // Persist the new tokens (user data stays the same).
      localStorage.setItem('auth_token', tokens.access_token);
      if (tokens.refresh_token) {
        localStorage.setItem('auth_refresh_token', tokens.refresh_token);
      }
      if (tokens.expires_at) {
        let expiresTimestamp = tokens.expires_at;
        if (typeof tokens.expires_at === 'string' && isNaN(Number(tokens.expires_at))) {
          expiresTimestamp = Math.floor(new Date(tokens.expires_at).getTime() / 1000);
        }
        localStorage.setItem('auth_expires_at', expiresTimestamp.toString());
      }

      await loadTenantData();
    } catch (err) {
      console.error('[TenantContext] switchTenant error:', err);
      throw err;
    } finally {
      setSwitching(false);
    }
  }, []);

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
    availableTenants,
    hasPermission,
    isOwner,
    isAdmin,
    loading,
    switching,
    error,
    refreshTenant: loadTenantData,
    switchTenant,
  };

  return (
    <TenantContext.Provider value={value}>
      {children}
    </TenantContext.Provider>
  );
};

export default TenantContext;
