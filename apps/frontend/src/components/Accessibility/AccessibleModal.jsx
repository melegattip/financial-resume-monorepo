import React, { useEffect, useRef } from 'react';
import { FaTimes } from 'react-icons/fa';
import { FocusManager } from './FocusManager';

const AccessibleModal = ({
  isOpen,
  onClose,
  title,
  description,
  children,
  size = 'medium',
  closeOnOverlay = true,
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
    if (closeOnOverlay && event.target === overlayRef.current) {
      onClose();
    }
  };

  const sizeClasses = {
    small: 'max-w-md',
    medium: 'max-w-2xl',
    large: 'max-w-4xl',
    fullscreen: 'max-w-full h-full'
  };

  if (!isOpen) return null;

  return (
    <div
      ref={overlayRef}
      className={`
        fixed inset-0 z-50 flex items-center justify-center
        bg-black bg-opacity-50 p-4
        ${className}
      `}
      onClick={handleOverlayClick}
      role="dialog"
      aria-modal="true"
      aria-labelledby={title ? "modal-title" : undefined}
      aria-describedby={description ? "modal-description" : undefined}
    >
      <FocusManager trapFocus autoFocus restoreFocus>
        <div
          ref={modalRef}
          className={`
            bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full
            ${sizeClasses[size]}
            max-h-[90vh] overflow-y-auto
            transform transition-all duration-300 ease-out
            animate-modal-enter
          `}
          role="document"
        >
          {/* Header */}
          <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
            <div className="flex-1">
              {title && (
                <h2 
                  id="modal-title" 
                  className="text-xl font-semibold text-gray-900 dark:text-gray-100"
                >
                  {title}
                </h2>
              )}
              {description && (
                <p 
                  id="modal-description" 
                  className="mt-1 text-sm text-gray-600 dark:text-gray-400"
                >
                  {description}
                </p>
              )}
            </div>
            
            <button
              onClick={onClose}
              className="
                ml-4 p-2 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300
                rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors
                focus:outline-none focus:ring-2 focus:ring-fr-primary
              "
              aria-label="Cerrar modal"
            >
              <FaTimes className="w-5 h-5" />
            </button>
          </div>

          {/* Content */}
          <div className="p-6">
            {children}
          </div>
        </div>
      </FocusManager>
    </div>
  );
};

export default AccessibleModal; 