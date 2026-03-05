import React, { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { FaCheckCircle, FaTimesCircle, FaSpinner } from 'react-icons/fa';
import Logo from '../components/Logo';
import environments from '../config/environments';

const VerifyEmail = () => {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') || '';

  const [status, setStatus] = useState('loading'); // 'loading' | 'success' | 'error'
  const [errorMsg, setErrorMsg] = useState('');

  useEffect(() => {
    if (!token) {
      setStatus('error');
      setErrorMsg('El enlace de verificación no es válido.');
      return;
    }

    const verify = async () => {
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

    verify();
  }, [token]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-fr-gray-50 dark:bg-gray-900 px-4">
      <div className="max-w-md w-full bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-10 text-center">
        <div className="flex justify-center mb-6">
          <Logo size="lg" showText={false} />
        </div>

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
            <p className="text-fr-gray-600 dark:text-gray-400 mb-8">{errorMsg}</p>
            <Link to="/login" className="btn-primary inline-block">
              Volver al inicio de sesión
            </Link>
          </>
        )}
      </div>
    </div>
  );
};

export default VerifyEmail;
