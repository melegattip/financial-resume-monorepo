import React, { useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { FaCheckCircle, FaTimesCircle, FaSpinner, FaEnvelope } from 'react-icons/fa';
import Logo from '../components/Logo';
import environments from '../config/environments';
import authService from '../services/authService';

const VerifyEmail = () => {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') || '';

  // 'idle' = waiting for user click | 'loading' | 'success' | 'error'
  const [status, setStatus] = useState(token ? 'idle' : 'error');
  const [errorMsg, setErrorMsg] = useState(token ? '' : 'El enlace de verificación no es válido.');
  const [resendEmail, setResendEmail] = useState('');
  const [isResending, setIsResending] = useState(false);

  const handleVerify = async () => {
    setStatus('loading');
    try {
      const baseURL = environments.USERS_API_URL;
      const res = await fetch(`${baseURL}/auth/verify-email/${encodeURIComponent(token)}`);
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.error || 'El enlace expiró o ya fue utilizado.');
      }
      setStatus('success');
    } catch (err) {
      setStatus('error');
      setErrorMsg(err.message);
    }
  };

  const handleResend = async () => {
    if (!resendEmail) return;
    setIsResending(true);
    try {
      await authService.resendVerification(resendEmail);
    } catch (_) {
      // toast already shown by authService
    } finally {
      setIsResending(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-fr-gray-50 dark:bg-gray-900 px-4">
      <div className="max-w-md w-full bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-10 text-center">
        <div className="flex justify-center mb-6">
          <Logo size="lg" showText={false} />
        </div>

        {status === 'idle' && (
          <>
            <FaEnvelope className="w-14 h-14 text-blue-500 mx-auto mb-4" />
            <h2 className="text-2xl font-bold text-fr-gray-900 dark:text-white mb-3">
              Verificá tu cuenta
            </h2>
            <p className="text-fr-gray-600 dark:text-gray-400 mb-8">
              Hacé clic en el botón para confirmar tu dirección de email y activar tu cuenta.
            </p>
            <button
              type="button"
              onClick={handleVerify}
              className="btn-primary inline-block w-full"
            >
              Verificar mi cuenta
            </button>
          </>
        )}

        {status === 'loading' && (
          <>
            <FaSpinner className="w-12 h-12 animate-spin text-blue-500 mx-auto mb-4" />
            <h2 className="text-xl font-bold text-fr-gray-900 dark:text-white mb-2">
              Verificando tu cuenta...
            </h2>
            <p className="text-fr-gray-500 dark:text-gray-400">Un momento por favor.</p>
          </>
        )}

        {status === 'success' && (
          <>
            <FaCheckCircle className="w-14 h-14 text-green-500 mx-auto mb-4" />
            <h2 className="text-2xl font-bold text-fr-gray-900 dark:text-white mb-3">
              ¡Cuenta verificada!
            </h2>
            <p className="text-fr-gray-600 dark:text-gray-400 mb-8">
              Tu email fue confirmado correctamente. Ya podés iniciar sesión y usar Niloft.
            </p>
            <Link to="/login" className="btn-primary inline-block">
              Iniciar sesión
            </Link>
          </>
        )}

        {status === 'error' && (
          <>
            <FaTimesCircle className="w-14 h-14 text-red-400 mx-auto mb-4" />
            <h2 className="text-2xl font-bold text-fr-gray-900 dark:text-white mb-3">
              Enlace inválido
            </h2>
            <p className="text-fr-gray-600 dark:text-gray-400 mb-6">{errorMsg}</p>

            <div className="mb-6 text-left">
              <p className="text-sm text-fr-gray-600 dark:text-gray-400 mb-2">
                ¿Querés recibir un nuevo enlace? Ingresá tu email:
              </p>
              <input
                type="email"
                value={resendEmail}
                onChange={(e) => setResendEmail(e.target.value)}
                placeholder="tu@email.com"
                className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-fr-gray-900 dark:text-white mb-3 focus:outline-none focus:ring-2 focus:ring-fr-primary"
              />
              <button
                type="button"
                onClick={handleResend}
                disabled={isResending || !resendEmail}
                className="w-full py-2 px-4 bg-fr-primary text-white rounded-lg text-sm font-medium hover:bg-fr-primary/90 disabled:opacity-50 flex items-center justify-center gap-2"
              >
                {isResending && <FaSpinner className="animate-spin" />}
                Reenviar email de verificación
              </button>
            </div>

            <Link to="/login" className="text-sm text-fr-primary hover:underline">
              Volver al inicio de sesión
            </Link>
          </>
        )}
      </div>
    </div>
  );
};

export default VerifyEmail;
