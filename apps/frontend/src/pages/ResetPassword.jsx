import React, { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { FaLock, FaEye, FaEyeSlash, FaSpinner, FaCheckCircle, FaArrowLeft } from 'react-icons/fa';
import Logo from '../components/Logo';
import environments from '../config/environments';

const ResetPassword = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get('token') || '';

  const [formData, setFormData] = useState({ password: '', confirm: '' });
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [done, setDone] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (!token) {
      setError('El enlace de restablecimiento no es válido. Solicitá uno nuevo.');
      return;
    }
    if (formData.password.length < 8) {
      setError('La contraseña debe tener al menos 8 caracteres.');
      return;
    }
    if (formData.password !== formData.confirm) {
      setError('Las contraseñas no coinciden.');
      return;
    }

    setIsSubmitting(true);
    try {
      const baseURL = environments.USERS_API_URL;
      const res = await fetch(`${baseURL}/auth/reset-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token, new_password: formData.password }),
      });

      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.error || 'El enlace expiró o no es válido. Solicitá uno nuevo.');
      }

      setDone(true);
      setTimeout(() => navigate('/login'), 3000);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!token) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-fr-gray-50">
        <div className="text-center max-w-md px-6">
          <h2 className="text-2xl font-bold text-fr-gray-900 mb-3">Enlace inválido</h2>
          <p className="text-fr-gray-600 mb-6">Este enlace de restablecimiento no es válido o ya fue utilizado.</p>
          <Link to="/forgot-password" className="btn-primary">
            Solicitar nuevo enlace
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex">
      {/* Panel izquierdo */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-fr-primary to-fr-secondary flex-col justify-center px-12">
        <div className="max-w-md">
          <h1 className="text-4xl font-bold text-white mb-6">Nueva contraseña</h1>
          <p className="text-fr-primary-light text-lg">
            Elegí una contraseña segura de al menos 8 caracteres para proteger tu cuenta.
          </p>
        </div>
      </div>

      {/* Panel derecho */}
      <div className="flex-1 flex flex-col justify-center px-6 py-12 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-md">
          <div className="flex justify-center mb-6">
            <Logo size="lg" showText={false} />
          </div>

          {done ? (
            <div className="text-center">
              <div className="flex justify-center mb-4">
                <FaCheckCircle className="w-14 h-14 text-green-500" />
              </div>
              <h2 className="text-2xl font-bold text-fr-gray-900 mb-3">¡Contraseña actualizada!</h2>
              <p className="text-fr-gray-600 mb-8">
                Tu contraseña fue restablecida correctamente. Redirigiendo al inicio de sesión...
              </p>
              <Link to="/login" className="btn-primary inline-flex items-center gap-2">
                <FaArrowLeft className="w-4 h-4" />
                Ir al inicio de sesión
              </Link>
            </div>
          ) : (
            <>
              <h2 className="text-center text-3xl font-bold text-fr-gray-900 mb-2">
                Nueva contraseña
              </h2>
              <p className="text-center text-fr-gray-600 mb-8">
                Ingresá y confirmá tu nueva contraseña.
              </p>

              {error && (
                <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-fr">
                  <p className="text-red-600 text-sm">{error}</p>
                </div>
              )}

              <form onSubmit={handleSubmit} className="space-y-5">
                <div>
                  <label className="block text-sm font-medium text-fr-gray-700 mb-2">
                    Nueva contraseña
                  </label>
                  <div className="relative">
                    <FaLock className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-fr-gray-400" />
                    <input
                      type={showPassword ? 'text' : 'password'}
                      value={formData.password}
                      onChange={(e) => setFormData(p => ({ ...p, password: e.target.value }))}
                      className="input pl-10 pr-10"
                      placeholder="Mínimo 8 caracteres"
                      disabled={isSubmitting}
                      autoFocus
                    />
                    <button
                      type="button"
                      className="absolute right-3 top-1/2 transform -translate-y-1/2 text-fr-gray-400 hover:text-fr-gray-600"
                      onClick={() => setShowPassword(v => !v)}
                    >
                      {showPassword ? <FaEyeSlash className="w-5 h-5" /> : <FaEye className="w-5 h-5" />}
                    </button>
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-fr-gray-700 mb-2">
                    Confirmar contraseña
                  </label>
                  <div className="relative">
                    <FaLock className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-fr-gray-400" />
                    <input
                      type={showConfirm ? 'text' : 'password'}
                      value={formData.confirm}
                      onChange={(e) => setFormData(p => ({ ...p, confirm: e.target.value }))}
                      className="input pl-10 pr-10"
                      placeholder="Repetí la contraseña"
                      disabled={isSubmitting}
                    />
                    <button
                      type="button"
                      className="absolute right-3 top-1/2 transform -translate-y-1/2 text-fr-gray-400 hover:text-fr-gray-600"
                      onClick={() => setShowConfirm(v => !v)}
                    >
                      {showConfirm ? <FaEyeSlash className="w-5 h-5" /> : <FaEye className="w-5 h-5" />}
                    </button>
                  </div>
                </div>

                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="btn-primary w-full flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isSubmitting ? (
                    <>
                      <FaSpinner className="w-4 h-4 animate-spin" />
                      <span>Guardando...</span>
                    </>
                  ) : (
                    <span>Guardar nueva contraseña</span>
                  )}
                </button>
              </form>

              <div className="mt-6 text-center">
                <Link
                  to="/login"
                  className="inline-flex items-center gap-1.5 text-sm text-fr-gray-500 hover:text-fr-gray-700 transition-colors"
                >
                  <FaArrowLeft className="w-3.5 h-3.5" />
                  Volver al inicio de sesión
                </Link>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default ResetPassword;
