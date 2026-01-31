import React, { useState, useEffect } from 'react';
import { FaQrcode, FaKey, FaEye, FaEyeSlash, FaCheckCircle, FaTimesCircle, FaDownload, FaPrint } from 'react-icons/fa';
import { authAPI_instance as authAPI } from '../services/authService';
import toast from 'react-hot-toast';

const TwoFASetup = ({ onComplete, onCancel }) => {
  const [step, setStep] = useState('setup'); // 'setup' | 'verify' | 'backup' | 'complete'
  const [setupData, setSetupData] = useState(null);
  const [verificationCode, setVerificationCode] = useState('');
  const [backupCodes, setBackupCodes] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [showBackupCodes, setShowBackupCodes] = useState(false);

  useEffect(() => {
    initializeSetup();
  }, []);

  const initializeSetup = async () => {
    try {
      setLoading(true);
      setError('');
      
      const response = await authAPI.post('/users/security/2fa/setup');
      setSetupData(response.data);
      setBackupCodes(response.data.backup_codes || []);
      setStep('setup');
    } catch (err) {
      setError('Error configurando 2FA. Inténtalo de nuevo.');
      toast.error('Error configurando 2FA');
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyCode = async () => {
    if (!verificationCode || verificationCode.length !== 6) {
      setError('Ingresa un código de 6 dígitos');
      return;
    }

    try {
      setLoading(true);
      setError('');
      
      await authAPI.post('/users/security/2fa/enable', {
        code: verificationCode
      });
      
      setStep('backup');
      toast.success('¡2FA activado correctamente!');
    } catch (err) {
      setError('Código inválido. Verifica el código de tu aplicación de autenticación.');
      toast.error('Código inválido');
    } finally {
      setLoading(false);
    }
  };

  const handleComplete = () => {
    setStep('complete');
    onComplete?.();
  };

  const handleCancel = () => {
    onCancel?.();
  };

  const downloadBackupCodes = () => {
    const codes = backupCodes.join('\n');
    const blob = new Blob([codes], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'backup-codes.txt';
    a.click();
    URL.revokeObjectURL(url);
    toast.success('Códigos de backup descargados');
  };

  const printBackupCodes = () => {
    const printWindow = window.open('', '_blank');
    const codes = backupCodes.join('\n');
    printWindow.document.write(`
      <html>
        <head><title>Códigos de Backup 2FA</title></head>
        <body>
          <h2>Códigos de Backup - Autenticación de Dos Factores</h2>
          <p><strong>Guarda estos códigos en un lugar seguro:</strong></p>
          <pre style="font-family: monospace; font-size: 14px; background: #f5f5f5; padding: 20px; border: 1px solid #ddd;">${codes}</pre>
          <p><em>Usa estos códigos si pierdes acceso a tu aplicación de autenticación.</em></p>
        </body>
      </html>
    `);
    printWindow.document.close();
    printWindow.print();
  };

  if (loading && !setupData) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-fr-primary mx-auto mb-4"></div>
          <p className="text-fr-gray-600 dark:text-gray-300">Configurando 2FA...</p>
        </div>
      </div>
    );
  }

  if (step === 'setup') {
    return (
      <div className="space-y-6">
        <div className="text-center">
          <FaQrcode className="w-12 h-12 text-fr-primary mx-auto mb-4" />
          <h3 className="text-xl font-semibold text-fr-gray-900 dark:text-gray-100 mb-2">
            Configurar Autenticación de Dos Factores
          </h3>
          <p className="text-fr-gray-600 dark:text-gray-300">
            Escanea el código QR con tu aplicación de autenticación
          </p>
        </div>

        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-fr p-4">
            <div className="flex items-center">
              <FaTimesCircle className="w-5 h-5 text-red-400 dark:text-red-300 mr-2" />
              <p className="text-red-800 dark:text-red-200">{error}</p>
            </div>
          </div>
        )}

        <div className="bg-white dark:bg-gray-800 border border-fr-gray-200 dark:border-gray-600 rounded-fr p-6 text-center">
          <div className="mb-4">
            <p className="text-sm text-fr-gray-600 dark:text-gray-300 mb-2">Código QR para escanear:</p>
            <div className="bg-gray-100 dark:bg-gray-700 p-4 rounded-fr inline-block">
              <img 
                src={`data:image/png;base64,${setupData?.qr_code_image}`} 
                alt="QR Code 2FA"
                className="w-48 h-48 mx-auto"
              />
            </div>
          </div>
          
          <div className="mb-4">
            <p className="text-sm text-fr-gray-600 dark:text-gray-300 mb-2">O ingresa manualmente:</p>
            <code className="bg-gray-100 dark:bg-gray-700 px-3 py-2 rounded text-sm font-mono break-all text-gray-800 dark:text-gray-200">
              {setupData?.secret}
            </code>
          </div>
        </div>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
              Código de verificación (6 dígitos)
            </label>
            <input
              type="text"
              value={verificationCode}
              onChange={(e) => setVerificationCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
              className="w-full px-4 py-2 border border-fr-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 rounded-fr focus:ring-2 focus:ring-fr-primary focus:border-transparent"
              placeholder="123456"
              maxLength={6}
            />
          </div>

          <div className="flex space-x-3">
            <button
              onClick={handleVerifyCode}
              disabled={loading || verificationCode.length !== 6}
              className="btn-primary flex-1"
            >
              {loading ? 'Verificando...' : 'Verificar y Activar'}
            </button>
            <button
              onClick={handleCancel}
              className="btn-outline flex-1"
            >
              Cancelar
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (step === 'backup') {
    return (
      <div className="space-y-6">
        <div className="text-center">
          <FaKey className="w-12 h-12 text-fr-primary mx-auto mb-4" />
          <h3 className="text-xl font-semibold text-fr-gray-900 dark:text-gray-100 mb-2">
            Códigos de Backup
          </h3>
          <p className="text-fr-gray-600 dark:text-gray-300">
            Guarda estos códigos en un lugar seguro. Los necesitarás si pierdes acceso a tu aplicación de autenticación.
          </p>
        </div>

        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-fr p-4">
          <div className="flex items-start">
            <FaCheckCircle className="w-5 h-5 text-yellow-600 dark:text-yellow-400 mr-2 mt-0.5" />
            <div>
              <p className="text-yellow-800 dark:text-yellow-200 font-medium">¡Importante!</p>
              <p className="text-yellow-700 dark:text-yellow-300 text-sm">
                Estos códigos solo se muestran una vez. Guárdalos de forma segura.
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 border border-fr-gray-200 dark:border-gray-600 rounded-fr p-6">
          <div className="flex justify-between items-center mb-4">
            <h4 className="font-medium text-fr-gray-900 dark:text-gray-100">Códigos de Backup</h4>
            <button
              onClick={() => setShowBackupCodes(!showBackupCodes)}
              className="text-fr-primary hover:text-fr-primary-dark dark:text-blue-400 dark:hover:text-blue-300 text-sm flex items-center space-x-1"
            >
              {showBackupCodes ? <FaEyeSlash /> : <FaEye />}
              <span>{showBackupCodes ? 'Ocultar' : 'Mostrar'}</span>
            </button>
          </div>
          
          <div className="grid grid-cols-2 gap-3">
            {backupCodes.map((code, index) => (
              <div key={index} className="bg-gray-100 dark:bg-gray-700 p-3 rounded text-center">
                <code className="font-mono text-sm text-gray-800 dark:text-gray-200">
                  {showBackupCodes ? code : '••••-••••'}
                </code>
              </div>
            ))}
          </div>
        </div>

        <div className="flex space-x-3">
          <button
            onClick={downloadBackupCodes}
            className="btn-outline flex-1 flex items-center justify-center space-x-2"
          >
            <FaDownload className="w-4 h-4" />
            <span>Descargar</span>
          </button>
          <button
            onClick={printBackupCodes}
            className="btn-outline flex-1 flex items-center justify-center space-x-2"
          >
            <FaPrint className="w-4 h-4" />
            <span>Imprimir</span>
          </button>
        </div>

        <button
          onClick={handleComplete}
          className="btn-primary w-full"
        >
          Completar Configuración
        </button>
      </div>
    );
  }

  if (step === 'complete') {
    return (
      <div className="text-center space-y-6">
        <FaCheckCircle className="w-16 h-16 text-green-500 dark:text-green-400 mx-auto" />
        <h3 className="text-xl font-semibold text-fr-gray-900 dark:text-gray-100">
          ¡2FA Configurado Exitosamente!
        </h3>
        <p className="text-fr-gray-600 dark:text-gray-300">
          Tu cuenta ahora está protegida con autenticación de dos factores.
        </p>
        <button
          onClick={handleComplete}
          className="btn-primary"
        >
          Continuar
        </button>
      </div>
    );
  }

  return null;
};

export default TwoFASetup; 