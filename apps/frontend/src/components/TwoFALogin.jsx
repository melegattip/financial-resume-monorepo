import React, { useState } from 'react';
import { FaQrcode, FaKey, FaEye, FaEyeSlash, FaTimesCircle, FaShieldAlt } from 'react-icons/fa';
import toast from 'react-hot-toast';

const TwoFALogin = ({ onVerify, onCancel, userEmail }) => {
  const [code, setCode] = useState('');
  const [isBackupCode, setIsBackupCode] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!code || code.length < 6) {
      setError('Ingresa un código válido');
      return;
    }

    try {
      setLoading(true);
      setError('');
      
      await onVerify(code, isBackupCode);
      toast.success('Código verificado correctamente');
    } catch (err) {
      setError('Código inválido. Verifica el código de tu aplicación de autenticación.');
      toast.error('Código inválido');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    onCancel?.();
  };

  const toggleCodeType = () => {
    setIsBackupCode(!isBackupCode);
    setCode('');
    setError('');
  };

  return (
    <div className="space-y-6">
      <div className="text-center">
        <FaShieldAlt className="w-12 h-12 text-fr-primary mx-auto mb-4" />
        <h3 className="text-xl font-semibold text-fr-gray-900 mb-2">
          Verificación de Dos Factores
        </h3>
        <p className="text-fr-gray-600">
          {isBackupCode 
            ? 'Ingresa uno de tus códigos de backup'
            : 'Ingresa el código de 6 dígitos de tu aplicación de autenticación'
          }
        </p>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-fr p-4">
          <div className="flex items-center">
            <FaTimesCircle className="w-5 h-5 text-red-400 mr-2" />
            <p className="text-red-800">{error}</p>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-fr-gray-700 mb-2">
            {isBackupCode ? 'Código de Backup' : 'Código de Verificación'}
          </label>
          <input
            type="text"
            value={code}
            onChange={(e) => {
              const value = e.target.value.replace(/\D/g, '');
              setCode(isBackupCode ? value.slice(0, 8) : value.slice(0, 6));
            }}
            className="w-full px-4 py-2 border border-fr-gray-300 rounded-fr focus:ring-2 focus:ring-fr-primary focus:border-transparent text-center text-lg font-mono"
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
        <div className="bg-blue-50 border border-blue-200 rounded-fr p-4">
          <div className="flex items-start">
            <FaQrcode className="w-5 h-5 text-blue-600 mr-2 mt-0.5" />
            <div>
              <p className="text-blue-800 font-medium">¿No tienes la app?</p>
              <p className="text-blue-700 text-sm">
                Descarga Google Authenticator, Authy, o cualquier aplicación compatible con TOTP.
              </p>
            </div>
          </div>
        </div>
      )}

      {isBackupCode && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-fr p-4">
          <div className="flex items-start">
            <FaKey className="w-5 h-5 text-yellow-600 mr-2 mt-0.5" />
            <div>
              <p className="text-yellow-800 font-medium">Códigos de Backup</p>
              <p className="text-yellow-700 text-sm">
                Usa estos códigos solo si perdiste acceso a tu aplicación de autenticación.
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default TwoFALogin; 