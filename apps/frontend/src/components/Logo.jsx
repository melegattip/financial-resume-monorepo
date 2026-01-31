import React from 'react';

/**
 * Componente Logo de NILOFT
 * Muestra el logo de la aplicación con diferentes tamaños y variantes
 */
const Logo = ({
  size = 'md',
  variant = 'full',
  className = '',
  showText = true,
  alt = 'NILOFT Logo'
}) => {
  // Definir tamaños disponibles
  const sizes = {
    xs: 'w-6 h-6',
    sm: 'w-8 h-8',
    md: 'w-12 h-12',
    lg: 'w-16 h-16',
    xl: 'w-24 h-24',
    '2xl': 'w-32 h-32'
  };

  // Definir clases de texto para acompañar el logo
  const textSizes = {
    xs: 'text-sm',
    sm: 'text-base',
    md: 'text-lg',
    lg: 'text-xl',
    xl: 'text-2xl',
    '2xl': 'text-3xl'
  };

  const logoClasses = [
    sizes[size] || sizes.md,
    'object-contain',
    className
  ].filter(Boolean).join(' ');

  const textClasses = [
    textSizes[size] || textSizes.md,
    'font-bold text-gray-900 dark:text-white ml-3'
  ].filter(Boolean).join(' ');

  return (
    <div className="flex items-center">
      {/* Logo Image */}
      <img
        src="/logo-niloft.png"
        alt={alt}
        className={logoClasses}
      />
      
      {/* Texto del logo */}
      {showText && (
        <span className={textClasses}>
          NILOFT
        </span>
      )}
    </div>
  );
};

export default Logo; 