import React from 'react';
import Logo from './Logo';

/**
 * Componente Brand de NILOFT
 * Combina el logo con el nombre de la marca para uso en headers y footers
 */
const Brand = ({
  size = 'md',
  layout = 'horizontal', // 'horizontal' | 'vertical'
  showTagline = false,
  className = '',
  onClick
}) => {
  const containerClasses = [
    'flex items-center',
    layout === 'vertical' ? 'flex-col space-y-2' : 'space-x-3',
    onClick ? 'cursor-pointer hover:opacity-80 transition-opacity' : '',
    className
  ].filter(Boolean).join(' ');

  const textSizes = {
    sm: 'text-lg',
    md: 'text-xl',
    lg: 'text-2xl',
    xl: 'text-3xl'
  };

  const taglineSizes = {
    sm: 'text-xs',
    md: 'text-sm',
    lg: 'text-base',
    xl: 'text-lg'
  };

  return (
    <div className={containerClasses} onClick={onClick}>
      <Logo size={size} showText={false} />
      <div className={layout === 'vertical' ? 'text-center' : ''}>
        <h1 className={`font-bold text-gray-900 dark:text-white ${textSizes[size] || textSizes.md}`}>
          Niloft
        </h1>
        {showTagline && (
          <p className={`text-gray-600 dark:text-gray-400 ${taglineSizes[size] || taglineSizes.md}`}>
            Tu asistente financiero
          </p>
        )}
      </div>
    </div>
  );
};

export default Brand; 