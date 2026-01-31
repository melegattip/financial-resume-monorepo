import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { FaSpinner } from 'react-icons/fa';
import { useAuthStatus } from '../contexts/AuthContext';

/**
 * Componente para proteger rutas que requieren autenticación
 */
const ProtectedRoute = ({ children, redirectTo = '/login' }) => {
  const { isAuthenticated, isLoading, isReady } = useAuthStatus();
  const location = useLocation();

  // Mostrar loading mientras se inicializa la autenticación
  if (!isReady) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-fr-gray-50">
        <div className="flex flex-col items-center space-y-4">
          <FaSpinner className="w-8 h-8 animate-spin text-fr-primary" />
          <p className="text-fr-gray-600">Verificando autenticación...</p>
        </div>
      </div>
    );
  }

  // Si no está autenticado, redireccionar al login
  if (!isAuthenticated) {
    // Guardar la ubicación actual para redireccionar después del login
    return (
      <Navigate 
        to={redirectTo} 
        state={{ from: location }} 
        replace 
      />
    );
  }

  // Si está autenticado, mostrar el contenido protegido
  return children;
};

/**
 * Componente para rutas que solo deben ser accesibles para usuarios NO autenticados
 * (como login y register)
 */
export const PublicOnlyRoute = ({ children, redirectTo = '/dashboard' }) => {
  const { isAuthenticated, isLoading, isReady } = useAuthStatus();

  // Mostrar loading mientras se inicializa la autenticación
  if (!isReady) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-fr-gray-50">
        <div className="flex flex-col items-center space-y-4">
          <FaSpinner className="w-8 h-8 animate-spin text-fr-primary" />
          <p className="text-fr-gray-600">Cargando...</p>
        </div>
      </div>
    );
  }

  // Si está autenticado, redireccionar al dashboard
  if (isAuthenticated) {
    return <Navigate to={redirectTo} replace />;
  }

  // Si no está autenticado, mostrar el contenido público
  return children;
};

export default ProtectedRoute; 