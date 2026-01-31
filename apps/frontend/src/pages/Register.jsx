import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { FaEye, FaEyeSlash, FaEnvelope, FaLock, FaUser, FaUserPlus, FaSpinner, FaCheckCircle } from 'react-icons/fa';
import { useAuth } from '../contexts/AuthContext';
import { validateEmail, sanitizeText } from '../utils/validation';
import Logo from '../components/Logo';

const Register = () => {
  const navigate = useNavigate();
  const { register, isAuthenticated, isLoading } = useAuth();

  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    password: '',
    confirmPassword: '',
  });

  const [errors, setErrors] = useState({});
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState({
    score: 0,
    message: '',
    color: 'gray',
    isValid: false
  });

  // Redireccionar si ya está autenticado
  useEffect(() => {
    if (isAuthenticated && !isLoading) {
      navigate('/dashboard', { replace: true });
    }
  }, [isAuthenticated, isLoading, navigate]);

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

    // Evaluar fortaleza de contraseña
    if (name === 'password') {
      evaluatePasswordStrength(value);
    }
  };

  // Evaluar fortaleza de la contraseña
  const evaluatePasswordStrength = (password) => {
    if (!password) {
      setPasswordStrength({ score: 0, message: '', color: 'gray' });
      return;
    }

    let score = 0;
    let message = '';
    let color = 'red';
    let isValid = true;
    let errors = [];

    // Criterios de evaluación del backend
    if (password.length >= 8) score += 1;
    if (/[a-z]/.test(password)) score += 1;
    if (/[A-Z]/.test(password)) score += 1;
    if (/[0-9]/.test(password)) score += 1;
    if (/[^A-Za-z0-9]/.test(password)) score += 1;

    // Validaciones específicas del backend
    if (!/[A-Z]/.test(password)) {
      errors.push('Debe contener al menos una letra mayúscula');
      isValid = false;
    }
    if (!/[^A-Za-z0-9]/.test(password)) {
      errors.push('Debe contener al menos un carácter especial');
      isValid = false;
    }
    
    // Verificar caracteres secuenciales (123, abc, etc.)
    const sequentialPatterns = ['123', 'abc', 'ABC', 'qwe', 'QWE'];
    const hasSequential = sequentialPatterns.some(pattern => 
      password.toLowerCase().includes(pattern)
    );
    if (hasSequential) {
      errors.push('No puede contener caracteres secuenciales');
      isValid = false;
    }

    // Determinar mensaje y color
    if (!isValid) {
      message = errors.join(', ');
      color = 'red';
    } else {
      switch (score) {
        case 0:
        case 1:
          message = 'Muy débil';
          color = 'red';
          break;
        case 2:
          message = 'Débil';
          color = 'orange';
          break;
        case 3:
          message = 'Moderada';
          color = 'yellow';
          break;
        case 4:
          message = 'Fuerte';
          color = 'green';
          break;
        case 5:
          message = 'Muy fuerte';
          color = 'green';
          break;
        default:
          message = 'Muy débil';
          color = 'red';
      }
    }

    setPasswordStrength({ score, message, color, isValid });
  };

  // Validar formulario
  const validateForm = () => {
    const registerErrors = {};
    
    // Validar nombre
    if (!formData.firstName.trim()) {
      registerErrors.firstName = 'El nombre es requerido';
    } else if (formData.firstName.trim().length < 2) {
      registerErrors.firstName = 'El nombre debe tener al menos 2 caracteres';
    } else if (!/^[a-zA-ZÀ-ÿ\s]+$/.test(formData.firstName)) {
      registerErrors.firstName = 'El nombre solo puede contener letras';
    }
    
    // Validar apellido
    if (!formData.lastName.trim()) {
      registerErrors.lastName = 'El apellido es requerido';
    } else if (formData.lastName.trim().length < 2) {
      registerErrors.lastName = 'El apellido debe tener al menos 2 caracteres';
    } else if (!/^[a-zA-ZÀ-ÿ\s]+$/.test(formData.lastName)) {
      registerErrors.lastName = 'El apellido solo puede contener letras';
    }
    
    // Validar email
    const emailValidation = validateEmail(formData.email);
    if (!emailValidation.isValid) {
      registerErrors.email = emailValidation.error;
    }
    
    // Validar contraseña
    if (!formData.password) {
      registerErrors.password = 'La contraseña es requerida';
    } else if (formData.password.length < 8) {
      registerErrors.password = 'La contraseña debe tener al menos 8 caracteres';
    } else if (!passwordStrength.isValid) {
      registerErrors.password = passwordStrength.message;
    }
    
    // Validar confirmación de contraseña
    if (!formData.confirmPassword) {
      registerErrors.confirmPassword = 'Confirma tu contraseña';
    } else if (formData.confirmPassword !== formData.password) {
      registerErrors.confirmPassword = 'Las contraseñas no coinciden';
    }

    setErrors(registerErrors);
    return Object.keys(registerErrors).length === 0;
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
      await register({
        firstName: sanitizeText(formData.firstName.trim()),
        lastName: sanitizeText(formData.lastName.trim()),
        email: sanitizeText(formData.email.trim().toLowerCase()),
        password: formData.password
      });
      
      // La navegación se maneja en el useEffect
    } catch (error) {
      setErrors({
        general: error.message || 'Error al registrar usuario'
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  // Mostrar loading si está inicializando
  if (isLoading) {
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
    <div className="min-h-screen flex">
      {/* Panel izquierdo - Información */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-fr-secondary to-fr-primary flex-col justify-center px-12">
        <div className="max-w-md">
          <h1 className="text-4xl font-bold text-white mb-6">
            Únete a nosotros
          </h1>
          <p className="text-fr-primary-light text-lg mb-8">
            Comienza tu viaje hacia la libertad financiera. 
            Registra tu cuenta y descubre una nueva forma de gestionar tu dinero.
          </p>
          <div className="space-y-4">
            <div className="flex items-center space-x-3">
              <FaCheckCircle className="w-5 h-5 text-white" />
              <span className="text-white">Cuenta gratuita por 30 días</span>
            </div>
            <div className="flex items-center space-x-3">
              <FaCheckCircle className="w-5 h-5 text-white" />
              <span className="text-white">Sin tarjeta de crédito requerida</span>
            </div>
            <div className="flex items-center space-x-3">
              <FaCheckCircle className="w-5 h-5 text-white" />
              <span className="text-white">Acceso completo a todas las funciones</span>
            </div>
            <div className="flex items-center space-x-3">
              <FaCheckCircle className="w-5 h-5 text-white" />
              <span className="text-white">Soporte 24/7</span>
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
            Crear Cuenta
          </h2>
          <p className="text-center text-fr-gray-600 mb-8">
            Completa el formulario para comenzar
          </p>

          {/* Error general */}
          {errors.general && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-fr">
              <p className="text-red-600 text-sm">{errors.general}</p>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Nombre y Apellido */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label htmlFor="firstName" className="block text-sm font-medium text-fr-gray-700 mb-2">
                  Nombre
                </label>
                <div className="relative">
                  <FaUser className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-fr-gray-400" />
                  <input
                    id="firstName"
                    name="firstName"
                    type="text"
                    autoComplete="given-name"
                    required
                    value={formData.firstName}
                    onChange={handleInputChange}
                    className={`input pl-10 ${errors.firstName ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : ''}`}
                    placeholder="Juan"
                    disabled={isSubmitting}
                  />
                </div>
                {errors.firstName && (
                  <p className="mt-1 text-sm text-red-600">{errors.firstName}</p>
                )}
              </div>

              <div>
                <label htmlFor="lastName" className="block text-sm font-medium text-fr-gray-700 mb-2">
                  Apellido
                </label>
                <div className="relative">
                  <FaUser className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-fr-gray-400" />
                  <input
                    id="lastName"
                    name="lastName"
                    type="text"
                    autoComplete="family-name"
                    required
                    value={formData.lastName}
                    onChange={handleInputChange}
                    className={`input pl-10 ${errors.lastName ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : ''}`}
                    placeholder="Pérez"
                    disabled={isSubmitting}
                  />
                </div>
                {errors.lastName && (
                  <p className="mt-1 text-sm text-red-600">{errors.lastName}</p>
                )}
              </div>
            </div>

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
                  placeholder="juan@ejemplo.com"
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
                  autoComplete="new-password"
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
              
              {/* Indicador de fortaleza de contraseña */}
              {formData.password && (
                <div className="mt-2">
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 h-2 bg-gray-200 rounded-full">
                      <div
                        className={`h-2 rounded-full transition-all duration-300 ${
                          passwordStrength.color === 'red' ? 'bg-red-500' :
                          passwordStrength.color === 'orange' ? 'bg-orange-500' :
                          passwordStrength.color === 'yellow' ? 'bg-yellow-500' :
                          'bg-green-500'
                        }`}
                        style={{ width: `${(passwordStrength.score / 5) * 100}%` }}
                      ></div>
                    </div>
                    <span className={`text-xs font-medium ${
                      passwordStrength.color === 'red' ? 'text-red-600' :
                      passwordStrength.color === 'orange' ? 'text-orange-600' :
                      passwordStrength.color === 'yellow' ? 'text-yellow-600' :
                      'text-green-600'
                    }`}>
                      {passwordStrength.message}
                    </span>
                  </div>
                </div>
              )}
              
              {errors.password && (
                <p className="mt-1 text-sm text-red-600">{errors.password}</p>
              )}
            </div>

            {/* Confirmar Contraseña */}
            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-fr-gray-700 mb-2">
                Confirmar Contraseña
              </label>
              <div className="relative">
                <FaLock className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-fr-gray-400" />
                <input
                  id="confirmPassword"
                  name="confirmPassword"
                  type={showConfirmPassword ? 'text' : 'password'}
                  autoComplete="new-password"
                  required
                  value={formData.confirmPassword}
                  onChange={handleInputChange}
                  className={`input pl-10 pr-10 ${errors.confirmPassword ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : ''}`}
                  placeholder="••••••••"
                  disabled={isSubmitting}
                />
                <button
                  type="button"
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-fr-gray-400 hover:text-fr-gray-600"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  disabled={isSubmitting}
                >
                  {showConfirmPassword ? (
                    <FaEyeSlash className="w-5 h-5" />
                  ) : (
                    <FaEye className="w-5 h-5" />
                  )}
                </button>
              </div>
              {errors.confirmPassword && (
                <p className="mt-1 text-sm text-red-600">{errors.confirmPassword}</p>
              )}
            </div>

            {/* Botón de envío */}
            <button
              type="submit"
              disabled={isSubmitting}
              className="btn-primary w-full flex items-center justify-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed mt-6"
            >
              {isSubmitting ? (
                <>
                  <FaSpinner className="w-4 h-4 animate-spin" />
                  <span>Creando cuenta...</span>
                </>
              ) : (
                <>
                  <FaUserPlus className="w-4 h-4" />
                  <span>Crear Cuenta</span>
                </>
              )}
            </button>
          </form>

          {/* Enlaces adicionales */}
          <div className="mt-6 text-center">
            <span className="text-fr-gray-600">¿Ya tienes una cuenta? </span>
            <Link 
              to="/login" 
              className="text-fr-primary hover:text-fr-primary-dark font-medium transition-colors"
            >
              Inicia sesión aquí
            </Link>
          </div>

          {/* Términos y condiciones */}
          <div className="mt-4 text-center">
            <p className="text-xs text-fr-gray-500">
              Al crear una cuenta, aceptas nuestros{' '}
              <Link to="/terms" className="text-fr-primary hover:underline">
                Términos de Servicio
              </Link>
              {' '}y{' '}
              <Link to="/privacy" className="text-fr-primary hover:underline">
                Política de Privacidad
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Register; 