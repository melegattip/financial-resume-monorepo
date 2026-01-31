import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import authService from '../services/authService';

// Crear el contexto
const AuthContext = createContext(null);

// Hook personalizado para usar el contexto
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    console.error('useAuth called outside AuthProvider');
    throw new Error('useAuth debe ser usado dentro de un AuthProvider');
  }
  return context;
};

// Estados posibles de autenticación
export const AUTH_STATES = {
  LOADING: 'loading',
  AUTHENTICATED: 'authenticated',
  UNAUTHENTICATED: 'unauthenticated',
  ERROR: 'error',
};

/**
 * Proveedor del contexto de autenticación
 */
export const AuthProvider = ({ children }) => {
  
  const [authState, setAuthState] = useState(AUTH_STATES.LOADING);
  const [user, setUser] = useState(null);
  const [isInitialized, setIsInitialized] = useState(false);

  // Verificar autenticación inicial
  useEffect(() => {
    const initializeAuth = async () => {
      try {
        if (authService.isAuthenticated()) {
          const currentUser = authService.getCurrentUser();
          setUser(currentUser);
          setAuthState(AUTH_STATES.AUTHENTICATED);
        } else {
          setAuthState(AUTH_STATES.UNAUTHENTICATED);
        }
      } catch (error) {
        console.error('Error initializing auth:', error);
        setAuthState(AUTH_STATES.UNAUTHENTICATED);
      } finally {
        setIsInitialized(true);
      }
    };

    initializeAuth();
  }, []);

  // Funciones de autenticación
  const login = useCallback(async (credentials) => {
    try {
      setAuthState(AUTH_STATES.LOADING);
      
      const result = await authService.login(credentials);
      setUser(result.data.user);
      setAuthState(AUTH_STATES.AUTHENTICATED);
      
      return result;
    } catch (error) {
      console.error('Login error:', error);
      // Para errores de 2FA, mantener estado UNAUTHENTICATED (no LOADING)
      setAuthState(AUTH_STATES.UNAUTHENTICATED);
      throw error;
    }
  }, []);

  const register = useCallback(async (userData) => {
    try {
      setAuthState(AUTH_STATES.LOADING);
      
      const result = await authService.register(userData);
      setUser(result.data.user);
      setAuthState(AUTH_STATES.AUTHENTICATED);
      
      return result;
    } catch (error) {
      console.error('Registration error:', error);
      setAuthState(AUTH_STATES.UNAUTHENTICATED);
      throw error;
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      await authService.logout();
    } catch (error) {
      console.warn('Error during logout:', error);
    } finally {
      setUser(null);
      setAuthState(AUTH_STATES.UNAUTHENTICATED);
    }
  }, []);

  const refreshToken = useCallback(async () => {
    try {
      const result = await authService.refreshToken();
      return result;
    } catch (error) {
      console.error('Token refresh error:', error);
      await logout();
      throw error;
    }
  }, [logout]);

  const changePassword = useCallback(async (passwordData) => {
    try {
      const result = await authService.changePassword(passwordData);
      return result;
    } catch (error) {
      console.error('Password change error:', error);
      throw error;
    }
  }, []);

  const updateProfile = useCallback(async (profileData) => {
    try {
      const updatedUser = await authService.updateProfile(profileData);
      
      // Actualizar el estado global con los datos reales del backend
      setUser(updatedUser);
      
      return { success: true, user: updatedUser };
    } catch (error) {
      console.error('Profile update error:', error);
      return { success: false, error: error.message };
    }
  }, []);

  const uploadAvatar = useCallback(async (file) => {
    try {
      const result = await authService.uploadAvatar(file);
      
      // Recargar el perfil del usuario para obtener la URL del avatar actualizado
      const updatedUser = await authService.getProfile();
      
      setUser(updatedUser);
      return { success: true, result };
    } catch (error) {
      console.error('Avatar upload error:', error);
      return { success: false, error: error.message };
    }
  }, []);

  // Valor del contexto
  const contextValue = {
    // Estado
    authState,
    user,
    isInitialized,
    
    // Estados derivados
    isAuthenticated: authState === AUTH_STATES.AUTHENTICATED,
    isLoading: authState === AUTH_STATES.LOADING,
    isError: authState === AUTH_STATES.ERROR,
    
    // Acciones
    login,
    register,
    logout,
    refreshToken,
    changePassword,
    updateProfile,
    uploadAvatar,
    
    // Utilidades
    hasRole: (role) => user?.roles?.includes(role) || false,
    getAuthHeaders: () => authService.getAuthHeaders(),
    timeUntilExpiry: authService.getSessionInfo().timeUntilExpiry,
    expiresAt: authService.getSessionInfo().expiresAt,
  };

  // Debug info available in development

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
};

/**
 * Hook para verificar autenticación
 */
export const useAuthStatus = () => {
  const { isAuthenticated, isLoading, isInitialized } = useAuth();
  
  return {
    isAuthenticated,
    isLoading: isLoading || !isInitialized,
    isReady: isInitialized && !isLoading,
  };
};

/**
 * Hook para obtener datos del usuario
 */
export const useUser = () => {
  const { user, updateProfile } = useAuth();
  
  const fullName = user 
    ? `${user.first_name || ''} ${user.last_name || ''}`.trim() || 'Usuario'
    : 'Usuario';
    
  const initials = user 
    ? `${user.first_name?.[0] || ''}${user.last_name?.[0] || ''}` 
    : 'U';
  
  return {
    user,
    fullName,
    initials,
    updateProfile,
  };
};

/**
 * Hook para acciones de autenticación
 */
export const useAuthActions = () => {
  const { login, register, logout, changePassword } = useAuth();
  
  return {
    login,
    register,
    logout,
    changePassword,
  };
};

export default AuthContext; 