import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { FaEnvelope, FaSpinner, FaCheckCircle, FaArrowLeft } from 'react-icons/fa';
import Logo from '../components/Logo';
import environments from '../config/environments';

const ForgotPassword = () => {
  const [email, setEmail] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!email.trim()) {
      setError('Ingresá tu dirección de email');
      return;
    }
    setIsSubmitting(true);
    setError('');

    try {
      const baseURL = environments.USERS_API_URL;
      await fetch(`${baseURL}/auth/request-password-reset`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email.trim().toLowerCase() }),
      });
      // Siempre mostramos éxito (anti-enumeración de cuentas)
      setSubmitted(true);
    } catch {
      setError('Error de conexión. Intentá de nuevo.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* Panel izquierdo */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-fr-primary to-fr-secondary flex-col justify-center px-12">
        <div className="max-w-md">
          <h1 className="text-4xl font-bold text-white mb-6">¿Olvidaste tu contraseña?</h1>
          <p className="text-fr-primary-light text-lg">
            No te preocupes. Te enviamos un enlace a tu correo para que puedas crear una nueva contraseña de forma segura.
          </p>
        </div>
      </div>

      {/* Panel derecho */}
      <div className="flex-1 flex flex-col justify-center px-6 py-12 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-md">
          <div className="flex justify-center mb-6">
            <Logo size="lg" showText={false} />
          </div>

          {submitted ? (
            <div className="text-center">
              <div className="flex justify-center mb-4">
                <FaCheckCircle className="w-14 h-14 text-green-500" />
              </div>
              <h2 className="text-2xl font-bold text-fr-gray-900 mb-3">Revisá tu correo</h2>
              <p className="text-fr-gray-600 mb-8">
                Si existe una cuenta con <strong>{email}</strong>, vas a recibir un enlace para restablecer tu contraseña en los próximos minutos.
              </p>
              <Link
                to="/login"
                className="btn-primary inline-flex items-center gap-2"
              >
                <FaArrowLeft className="w-4 h-4" />
                Volver al inicio de sesión
              </Link>
            </div>
          ) : (
            <>
              <h2 className="text-center text-3xl font-bold text-fr-gray-900 mb-2">
                Restablecer contraseña
              </h2>
              <p className="text-center text-fr-gray-600 mb-8">
                Ingresá tu email y te enviamos un enlace para crear una nueva contraseña.
              </p>

              {error && (
                <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-fr">
                  <p className="text-red-600 text-sm">{error}</p>
                </div>
              )}

              <form onSubmit={handleSubmit} className="space-y-6">
                <div>
                  <label htmlFor="email" className="block text-sm font-medium text-fr-gray-700 mb-2">
                    Email
                  </label>
                  <div className="relative">
                    <FaEnvelope className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-fr-gray-400" />
                    <input
                      id="email"
                      type="email"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      className="input pl-10"
                      placeholder="tu@email.com"
                      disabled={isSubmitting}
                      autoFocus
                    />
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
                      <span>Enviando...</span>
                    </>
                  ) : (
                    <span>Enviar enlace</span>
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

export default ForgotPassword;
