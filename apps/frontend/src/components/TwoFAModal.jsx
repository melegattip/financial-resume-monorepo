import React, { useState } from 'react';
import { FaQrcode, FaKey, FaTimes, FaShieldAlt } from 'react-icons/fa';
import toast from 'react-hot-toast';

const TwoFAModal = ({ isOpen, onClose, onVerify, userEmail }) => {
  const [code, setCode] = useState('');
  const [isBackupCode, setIsBackupCode] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!code || code.length < (isBackupCode ? 8 : 6)) {
      setError('Ingresa un código válido');
      return;
    }

    try {
      setLoading(true);
      setError('');
      
      await onVerify(code, isBackupCode);
      toast.success('Código verificado correctamente');
      onClose();
    } catch (err) {
      console.log('❌ Error capturado en TwoFAModal:', err);
      const errorMessage = '❌ Código 2FA inválido. Por favor, verifica el código de tu aplicación de autenticación e intenta nuevamente.';
      setError(errorMessage);
      toast.error('Código 2FA inválido', {
        duration: 4000,
        position: 'top-center',
        style: {
          background: '#ef4444',
          color: '#ffffff',
          fontSize: '14px',
          fontWeight: '600'
        }
      });
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    setCode('');
    setError('');
    setIsBackupCode(false);
    onClose();
  };

  const toggleCodeType = () => {
    setIsBackupCode(!isBackupCode);
    setCode('');
    setError('');
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md mx-4">
        <div className="flex justify-between items-center mb-4">
          <div className="flex items-center space-x-2">
            <FaShieldAlt className="w-6 h-6 text-fr-primary" />
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              Verificación de Dos Factores
            </h3>
          </div>
          <button
            onClick={handleCancel}
            className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
          >
            <FaTimes className="w-5 h-5" />
          </button>
        </div>

        <div className="space-y-4">
          <p className="text-gray-600 dark:text-gray-300">
            {isBackupCode 
              ? 'Ingresa uno de tus códigos de backup'
              : 'Ingresa el código de 6 dígitos de tu aplicación de autenticación'
            }
          </p>

          {error && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 animate-pulse">
              <div className="flex items-start space-x-2">
                <div className="flex-shrink-0">
                  <svg className="w-5 h-5 text-red-600 dark:text-red-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                  </svg>
                </div>
                <div>
                  <p className="text-red-800 dark:text-red-200 text-sm font-medium">Error de Verificación</p>
                  <p className="text-red-700 dark:text-red-300 text-xs mt-1">{error}</p>
                </div>
              </div>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {isBackupCode ? 'Código de Backup' : 'Código de Verificación'}
              </label>
              <input
                type="text"
                value={code}
                onChange={(e) => {
                  const value = e.target.value.replace(/\D/g, '');
                  setCode(isBackupCode ? value.slice(0, 8) : value.slice(0, 6));
                }}
                className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded focus:ring-2 focus:ring-fr-primary focus:border-transparent text-center text-lg font-mono bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                placeholder={isBackupCode ? "XXXX-XXXX" : "123456"}
                maxLength={isBackupCode ? 8 : 6}
                autoFocus
              />
            </div>

            <div className="flex justify-center">
              <button
                type="button"
                onClick={toggleCodeType}
                className="text-fr-primary hover:text-fr-primary-dark text-sm flex items-center space-x-1"
              >
                {isBackupCode ? <FaQrcode className="w-4 h-4" /> : <FaKey className="w-4 h-4" />}
                <span>
                  {isBackupCode ? 'Usar código de aplicación' : 'Usar código de backup'}
                </span>
              </button>
            </div>

            <div className="flex space-x-3">
              <button
                type="submit"
                disabled={loading || code.length < (isBackupCode ? 8 : 6)}
                className="btn-primary flex-1"
              >
                {loading ? 'Verificando...' : 'Verificar'}
              </button>
              <button
                type="button"
                onClick={handleCancel}
                className="btn-outline flex-1"
              >
                Cancelar
              </button>
            </div>
          </form>

          {!isBackupCode && (
            <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded p-3">
              <div className="flex items-start">
                <FaQrcode className="w-4 h-4 text-blue-600 dark:text-blue-400 mr-2 mt-0.5" />
                <div>
                  <p className="text-blue-800 dark:text-blue-200 font-medium text-sm">¿No tienes la app?</p>
                  <p className="text-blue-700 dark:text-blue-300 text-xs">
                    Descarga Google Authenticator, Authy, o cualquier aplicación compatible con TOTP.
                  </p>
                </div>
              </div>
            </div>
          )}

          {isBackupCode && (
            <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded p-3">
              <div className="flex items-start">
                <FaKey className="w-4 h-4 text-yellow-600 dark:text-yellow-400 mr-2 mt-0.5" />
                <div>
                  <p className="text-yellow-800 dark:text-yellow-200 font-medium text-sm">Códigos de Backup</p>
                  <p className="text-yellow-700 dark:text-yellow-300 text-xs">
                    Usa estos códigos solo si perdiste acceso a tu aplicación de autenticación.
                  </p>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default TwoFAModal; 