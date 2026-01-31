import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { FaEye, FaEyeSlash, FaEnvelope, FaLock, FaSignInAlt, FaSpinner } from 'react-icons/fa';
import { useAuth } from '../contexts/AuthContext';
import { validateEmail, sanitizeText } from '../utils/validation';
import Logo from '../components/Logo';
import TwoFAModal from '../components/TwoFAModal';
import environments from '../config/environments';


const Login = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login, isAuthenticated, isLoading, isInitialized } = useAuth();

  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

  const [errors, setErrors] = useState({});
  const [showPassword, setShowPassword] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [show2FAModal, setShow2FAModal] = useState(false);
  const [twoFACredentials, setTwoFACredentials] = useState(null);

  // Redireccionar si ya está autenticado
  useEffect(() => {
    if (isAuthenticated && !isLoading) {
      const from = location.state?.from?.pathname || '/dashboard';
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, isLoading, navigate, location]);

  // Manejar cambios en los inputs
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));

    // Limpiar error del campo cuando el usuario empiece a escribir
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  // Validar formulario
  const validateForm = () => {
    const loginErrors = {};
    
    // Validar email
    const emailValidation = validateEmail(formData.email);
    if (!emailValidation.isValid) {
      loginErrors.email = emailValidation.error;
    }
    
    // Validar contraseña
    if (!formData.password) {
      loginErrors.password = 'La contraseña es requerida';
    } else if (formData.password.length < 6) {
      loginErrors.password = 'La contraseña debe tener al menos 6 caracteres';
    }

    setErrors(loginErrors);
    return Object.keys(loginErrors).length === 0;
  };

  // Verificar si el usuario requiere 2FA antes de hacer login
  const check2FARequirement = async (credentials) => {
    try {
      // Usar la configuración de ambiente correcta
      const baseURL = environments.USERS_API_URL;
      console.log('🔧 [Login] Using USERS_API_URL:', baseURL);
      const response = await fetch(`${baseURL}/auth/check-2fa`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials)
      });
      
      if (!response.ok) {
        throw new Error('Error verificando 2FA');
      }
      
      const data = await response.json();
      return data.requires_2fa || false;
    } catch (error) {
      console.error('Error verificando 2FA:', error);
      return false;
    }
  };

  // Manejar envío del formulario
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);
    setErrors({});

    try {
      const credentials = {
        email: sanitizeText(formData.email.trim().toLowerCase()),
        password: formData.password
      };

      // Paso 1: Verificar si el usuario requiere 2FA
      const requires2FA = await check2FARequirement(credentials);
      
      if (requires2FA) {
        // Paso 2: Si requiere 2FA, guardar credenciales y mostrar modal
        setTwoFACredentials(credentials);
        setShow2FAModal(true);
        setIsSubmitting(false);
        return; // Importante: salir aquí para no continuar
      } else {
        // Paso 3: Si no requiere 2FA, hacer login directamente
        await login(credentials);
        // La navegación se maneja en el useEffect
      }
    } catch (error) {
      console.error('Error en login:', error);
      setErrors({
        general: error.message || 'Error al iniciar sesión'
      });
      // NO limpiar el formulario aquí
    } finally {
      setIsSubmitting(false);
    }
  };

  // Manejar verificación de 2FA
  const handle2FAVerification = async (code, isBackupCode = false) => {
    try {
      await login({
        ...twoFACredentials,
        twofa_code: code
      });
      
      setShow2FAModal(false);
      setTwoFACredentials(null);
    } catch (error) {
      // Re-lanzar el error para que el modal lo maneje
      console.error('❌ Error en verificación 2FA:', error);
      throw error;
    }
  };

  // Manejar cancelación de 2FA
  const handle2FACancel = () => {
    setShow2FAModal(false);
    setTwoFACredentials(null);
  };

  // Mostrar loading solo si está inicializando la app (no durante login)
  if (!isInitialized) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-fr-gray-50">
        <div className="flex flex-col items-center space-y-4">
          <FaSpinner className="w-8 h-8 animate-spin text-fr-primary" />
          <p className="text-fr-gray-600">Cargando...</p>
        </div>
      </div>
    );
  }



  return (
    <>
      <div className="min-h-screen flex">
        {/* Panel izquierdo - Información */}
        <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-fr-primary to-fr-secondary flex-col justify-center px-12">
          <div className="max-w-md">
            <h1 className="text-4xl font-bold text-white mb-6">
              Bienvenido
            </h1>
            <p className="text-fr-primary-light text-lg mb-8">
              Gestiona tus finanzas personales de manera inteligente. 
              Controla tus gastos, maximiza tus ingresos y alcanza tus metas financieras.
            </p>
            <div className="space-y-4">
              <div className="flex items-center space-x-3">
                <div className="w-2 h-2 bg-white rounded-full"></div>
                <span className="text-white">Dashboard con análisis en tiempo real</span>
              </div>
              <div className="flex items-center space-x-3">
                <div className="w-2 h-2 bg-white rounded-full"></div>
                <span className="text-white">Categorización inteligente de gastos</span>
              </div>
              <div className="flex items-center space-x-3">
                <div className="w-2 h-2 bg-white rounded-full"></div>
                <span className="text-white">Reportes detallados y predicciones</span>
              </div>
            </div>
          </div>
        </div>

        {/* Panel derecho - Formulario */}
        <div className="flex-1 flex flex-col justify-center px-6 py-12 lg:px-8">
          <div className="sm:mx-auto sm:w-full sm:max-w-md">
            <div className="flex justify-center mb-6">
              <Logo size="lg" showText={false} />
            </div>

            <h2 className="text-center text-3xl font-bold text-fr-gray-900 mb-2">
              Iniciar Sesión
            </h2>
            <p className="text-center text-fr-gray-600 mb-8">
              Accede a tu cuenta para continuar
            </p>

            {/* Error general */}
            {errors.general && (
              <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-fr">
                <p className="text-red-600 text-sm">{errors.general}</p>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Email */}
              <div>
                <label htmlFor="email" className="block text-sm font-medium text-fr-gray-700 mb-2">
                  Email
                </label>
                <div className="relative">
                  <FaEnvelope className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-fr-gray-400" />
                  <input
                    id="email"
                    name="email"
                    type="email"
                    autoComplete="email"
                    required
                    value={formData.email}
                    onChange={handleInputChange}
                    className={`input pl-10 ${errors.email ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : ''}`}
                    placeholder="tu@email.com"
                    disabled={isSubmitting}
                  />
                </div>
                {errors.email && (
                  <p className="mt-1 text-sm text-red-600">{errors.email}</p>
                )}
              </div>

              {/* Contraseña */}
              <div>
                <label htmlFor="password" className="block text-sm font-medium text-fr-gray-700 mb-2">
                  Contraseña
                </label>
                <div className="relative">
                  <FaLock className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-fr-gray-400" />
                  <input
                    id="password"
                    name="password"
                    type={showPassword ? 'text' : 'password'}
                    autoComplete="current-password"
                    required
                    value={formData.password}
                    onChange={handleInputChange}
                    className={`input pl-10 pr-10 ${errors.password ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : ''}`}
                    placeholder="••••••••"
                    disabled={isSubmitting}
                  />
                  <button
                    type="button"
                    className="absolute right-3 top-1/2 transform -translate-y-1/2 text-fr-gray-400 hover:text-fr-gray-600"
                    onClick={() => setShowPassword(!showPassword)}
                    disabled={isSubmitting}
                  >
                    {showPassword ? (
                      <FaEyeSlash className="w-5 h-5" />
                    ) : (
                      <FaEye className="w-5 h-5" />
                    )}
                  </button>
                </div>
                {errors.password && (
                  <p className="mt-1 text-sm text-red-600">{errors.password}</p>
                )}
              </div>

              {/* Botón de envío */}
              <button
                type="submit"
                disabled={isSubmitting}
                className="btn-primary w-full flex items-center justify-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? (
                  <>
                    <FaSpinner className="w-4 h-4 animate-spin" />
                    <span>Iniciando sesión...</span>
                  </>
                ) : (
                  <>
                    <FaSignInAlt className="w-4 h-4" />
                    <span>Iniciar Sesión</span>
                  </>
                )}
              </button>
            </form>

            {/* Enlaces adicionales */}
            <div className="mt-6 space-y-4">
              <div className="text-center">
                <span className="text-fr-gray-600">¿No tienes una cuenta? </span>
                <Link 
                  to="/register" 
                  className="text-fr-primary hover:text-fr-primary-dark font-medium transition-colors"
                >
                  Regístrate aquí
                </Link>
              </div>

              <div className="text-center">
                <Link 
                  to="/forgot-password" 
                  className="text-sm text-fr-gray-500 hover:text-fr-gray-700 transition-colors"
                >
                  ¿Olvidaste tu contraseña?
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Modal de verificación 2FA */}
      <TwoFAModal
        isOpen={show2FAModal}
        onClose={handle2FACancel}
        onVerify={handle2FAVerification}
        userEmail={twoFACredentials?.email}
      />
    </>
  );
};

export default Login; 