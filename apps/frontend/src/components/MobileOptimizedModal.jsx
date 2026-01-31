import React, { useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import { FaTimes } from 'react-icons/fa';

/**
 * Modal optimizado para dispositivos móviles
 * Proporciona mejor UX en pantallas pequeñas con:
 * - Altura completa en móvil
 * - Mejor padding y spacing
 * - Scroll interno optimizado
 * - Botones con touch targets adecuados
 */
const MobileOptimizedModal = ({
  isOpen,
  onClose,
  title,
  children,
  maxWidth = 'md',
  showCloseButton = true,
  preventBackgroundClose = false,
  className = ''
}) => {
  const modalRef = useRef(null);
  const overlayRef = useRef(null);

  // Prevenir scroll del body cuando el modal está abierto
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
      return () => {
        document.body.style.overflow = 'auto';
      };
    }
  }, [isOpen]);

  // Manejo de teclas globales
  useEffect(() => {
    const handleKeyDown = (event) => {
      if (event.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener('keydown', handleKeyDown);
      return () => document.removeEventListener('keydown', handleKeyDown);
    }
  }, [isOpen, onClose]);

  const handleOverlayClick = (event) => {
    if (!preventBackgroundClose && event.target === overlayRef.current) {
      onClose();
    }
  };

  const maxWidthClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
    '2xl': 'max-w-2xl',
    '3xl': 'max-w-3xl',
    full: 'max-w-full'
  };

  if (!isOpen) return null;

  return createPortal(
    <div
      ref={overlayRef}
      className="fixed inset-0 bg-black bg-opacity-50 z-[9999] flex items-end sm:items-center justify-center"
      onClick={handleOverlayClick}
      role="dialog"
      aria-modal="true"
      aria-labelledby={title ? "mobile-modal-title" : undefined}
    >
      <div
        ref={modalRef}
        className={`
          bg-white dark:bg-gray-800 w-full h-full sm:h-auto sm:max-h-[90vh] 
          sm:rounded-t-2xl sm:rounded-b-none md:rounded-xl 
          transform transition-all duration-300 ease-out
          ${maxWidthClasses[maxWidth]} sm:m-4
          flex flex-col
          ${className}
        `}
        role="document"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        {(title || showCloseButton) && (
          <div className="flex items-center justify-between p-4 sm:p-6 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
            <div className="flex-1 min-w-0">
              {title && (
                <h2 
                  id="mobile-modal-title" 
                  className="text-lg sm:text-xl font-semibold text-gray-900 dark:text-gray-100 truncate"
                >
                  {title}
                </h2>
              )}
            </div>
            
            {showCloseButton && (
              <button
                onClick={onClose}
                className="
                  ml-4 p-2 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300
                  rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors
                  focus:outline-none focus:ring-2 focus:ring-fr-primary
                  min-w-10 min-h-10 flex items-center justify-center
                "
                aria-label="Cerrar modal"
              >
                <FaTimes className="w-5 h-5" />
              </button>
            )}
          </div>
        )}

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-4 sm:p-6">
          <div className="space-y-4 sm:space-y-6">
            {children}
          </div>
        </div>
      </div>
    </div>,
    document.body
  );
};

/**
 * Componente de Form Group optimizado para móvil
 * Proporciona spacing y layout consistente para formularios
 */
export const MobileFormGroup = ({ 
  label, 
  children, 
  error, 
  helpText, 
  required = false,
  className = '' 
}) => {
  return (
    <div className={`space-y-2 ${className}`}>
      {label && (
        <label className="block text-sm sm:text-base font-medium text-gray-700 dark:text-gray-300">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      
      <div className="space-y-1">
        {children}
        
        {error && (
          <p className="text-sm text-red-600 dark:text-red-400 flex items-start space-x-1">
            <span className="text-red-500 mt-0.5">⚠️</span>
            <span>{error}</span>
          </p>
        )}
        
        {helpText && !error && (
          <p className="text-xs sm:text-sm text-gray-500 dark:text-gray-400">
            {helpText}
          </p>
        )}
      </div>
    </div>
  );
};

/**
 * Botones de acción optimizados para móvil
 * Proporciona touch targets adecuados y mejor UX
 */
export const MobileActionButtons = ({ 
  primaryAction, 
  secondaryAction, 
  layout = 'horizontal', // 'horizontal' | 'vertical' | 'stacked'
  className = '' 
}) => {
  const layoutClasses = {
    horizontal: 'flex flex-row space-x-3',
    vertical: 'flex flex-col space-y-3',
    stacked: 'flex flex-col-reverse sm:flex-row sm:space-y-0 sm:space-x-3 space-y-3'
  };

  return (
    <div className={`${layoutClasses[layout]} pt-4 sm:pt-6 border-t border-gray-200 dark:border-gray-600 ${className}`}>
      {secondaryAction && (
        <button
          type={secondaryAction.type || 'button'}
          onClick={secondaryAction.onClick}
          disabled={secondaryAction.disabled}
          className="
            flex-1 px-4 py-3 sm:py-2.5 text-sm sm:text-base font-medium 
            text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 
            border border-gray-300 dark:border-gray-600 rounded-lg 
            hover:bg-gray-50 dark:hover:bg-gray-600 
            focus:outline-none focus:ring-2 focus:ring-fr-primary 
            disabled:opacity-50 disabled:cursor-not-allowed
            transition-colors duration-200
            min-h-12 sm:min-h-10
          "
        >
          {secondaryAction.label}
        </button>
      )}
      
      {primaryAction && (
        <button
          type={primaryAction.type || 'button'}
          onClick={primaryAction.onClick}
          disabled={primaryAction.disabled}
          className={`
            flex-1 px-4 py-3 sm:py-2.5 text-sm sm:text-base font-medium 
            text-white rounded-lg 
            focus:outline-none focus:ring-2 focus:ring-opacity-50 
            disabled:opacity-50 disabled:cursor-not-allowed
            transition-colors duration-200
            min-h-12 sm:min-h-10
            ${primaryAction.variant === 'secondary' 
              ? 'bg-fr-secondary hover:bg-green-600 focus:ring-fr-secondary' 
              : 'bg-fr-primary hover:bg-blue-600 focus:ring-fr-primary'
            }
          `}
        >
          {primaryAction.loading ? (
            <div className="flex items-center justify-center space-x-2">
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              <span>{primaryAction.loadingLabel || 'Procesando...'}</span>
            </div>
          ) : (
            primaryAction.label
          )}
        </button>
      )}
    </div>
  );
};

export default MobileOptimizedModal; 