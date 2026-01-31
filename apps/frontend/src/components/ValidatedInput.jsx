import React, { useState, useEffect, useCallback } from 'react';
import { formatNumericInput } from '../utils/validation';

/**
 * Componente de Input con Validación Integrada
 * Proporciona validación en tiempo real, formateo automático y feedback visual
 */
const ValidatedInput = ({
  // Props básicas del input
  type = 'text',
  value = '',
  onChange,
  onBlur,
  placeholder,
  disabled = false,
  readOnly = false,
  autoComplete,
  id,
  name,
  
  // Props de validación
  validator,
  validateOnChange = false,
  validateOnBlur = true,
  showErrorOnTouch = true,
  
  // Props de presentación
  label,
  required = false,
  helpText,
  icon,
  iconPosition = 'left', // 'left' | 'right'
  
  // Props de estilo
  className = '',
  inputClassName = '',
  labelClassName = '',
  errorClassName = '',
  size = 'md', // 'sm' | 'md' | 'lg'
  variant = 'default', // 'default' | 'outlined' | 'filled'
  
  // Props específicas para números
  allowNegative = false,
  maxDecimals = 2,
  formatOnBlur = true,
  
  // Props adicionales
  ...restProps
}) => {
  // Estados internos
  const [error, setError] = useState(null);
  const [touched, setTouched] = useState(false);
  const [focused, setFocused] = useState(false);

  // Determinar si mostrar error
  const shouldShowError = error && (touched || !showErrorOnTouch);

  // Formatear valor para tipos numéricos
  const formatValue = useCallback((val) => {
    if (type === 'number' || type === 'currency') {
      return formatNumericInput(val);
    }
    return val;
  }, [type]);

  // Validar valor
  const validateValue = useCallback((val) => {
    if (!validator) return null;
    
    const validation = validator(val);
    return validation.isValid ? null : validation.error;
  }, [validator]);

  // Manejar cambio de valor
  const handleChange = useCallback((e) => {
    let newValue = e.target.value;
    
    // Formatear para tipos numéricos
    if (type === 'number' || type === 'currency') {
      newValue = formatValue(newValue);
      
      // Manejar números negativos
      if (!allowNegative && newValue.startsWith('-')) {
        newValue = newValue.substring(1);
      }
    }
    
    // Validar en tiempo real si está habilitado
    if (validateOnChange && touched) {
      const validationError = validateValue(newValue);
      setError(validationError);
    }
    
    // Llamar onChange del padre
    if (onChange) {
      // Crear evento sintético con el valor formateado
      const syntheticEvent = {
        ...e,
        target: {
          ...e.target,
          value: newValue,
          name: name || e.target.name,
          id: id || e.target.id
        }
      };
      onChange(syntheticEvent);
    }
  }, [type, formatValue, allowNegative, validateOnChange, touched, validateValue, onChange, name, id]);

  // Manejar blur
  const handleBlur = useCallback((e) => {
    setTouched(true);
    setFocused(false);
    
    let finalValue = e.target.value;
    
    // Formatear en blur para números si está habilitado
    if ((type === 'number' || type === 'currency') && formatOnBlur && finalValue) {
      const numValue = parseFloat(finalValue);
      if (!isNaN(numValue)) {
        finalValue = numValue.toFixed(maxDecimals);
        
        // Actualizar el valor formateado
        if (onChange) {
          const syntheticEvent = {
            ...e,
            target: {
              ...e.target,
              value: finalValue,
              name: name || e.target.name,
              id: id || e.target.id
            }
          };
          onChange(syntheticEvent);
        }
      }
    }
    
    // Validar en blur si está habilitado
    if (validateOnBlur) {
      const validationError = validateValue(finalValue);
      setError(validationError);
    }
    
    // Llamar onBlur del padre
    if (onBlur) {
      onBlur(e);
    }
  }, [type, formatOnBlur, maxDecimals, validateOnBlur, validateValue, onChange, onBlur, name, id]);

  // Manejar focus
  const handleFocus = useCallback((e) => {
    setFocused(true);
  }, []);

  // Validar cuando cambia el valor externamente
  useEffect(() => {
    if (touched && validator) {
      const validationError = validateValue(value);
      setError(validationError);
    }
  }, [value, touched, validator, validateValue]);

  // Clases CSS dinámicas
  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-3 py-2 text-base',
    lg: 'px-4 py-3 text-lg'
  };

  const variantClasses = {
    default: 'border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800',
    outlined: 'border-2 border-gray-300 dark:border-gray-600 bg-transparent',
    filled: 'border-0 bg-gray-100 dark:bg-gray-700'
  };

  const inputClasses = [
    // Clases base
    'w-full rounded-md shadow-sm transition-colors duration-200',
    'focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
    'disabled:bg-gray-100 disabled:text-gray-500 disabled:cursor-not-allowed',
    'dark:text-white dark:placeholder-gray-400',
    
    // Tamaño
    sizeClasses[size],
    
    // Variante
    variantClasses[variant],
    
    // Estados de error
    shouldShowError
      ? 'border-red-500 dark:border-red-400 bg-red-50 dark:bg-red-900/20 focus:ring-red-500 focus:border-red-500'
      : '',
    
    // Estados de focus
    focused && !shouldShowError
      ? 'ring-2 ring-blue-500 border-blue-500'
      : '',
    
    // Padding para iconos
    icon && iconPosition === 'left' ? 'pl-10' : '',
    icon && iconPosition === 'right' ? 'pr-10' : '',
    
    // Clases personalizadas
    inputClassName
  ].filter(Boolean).join(' ');

  const labelClasses = [
    'block text-sm font-medium mb-2 transition-colors duration-200',
    shouldShowError
      ? 'text-red-600 dark:text-red-400'
      : 'text-gray-700 dark:text-gray-300',
    labelClassName
  ].filter(Boolean).join(' ');

  const errorClasses = [
    'mt-1 text-sm transition-all duration-200',
    'text-red-600 dark:text-red-400',
    errorClassName
  ].filter(Boolean).join(' ');

  const containerClasses = [
    'mb-4',
    className
  ].filter(Boolean).join(' ');

  return (
    <div className={containerClasses}>
      {/* Label */}
      {label && (
        <label 
          htmlFor={id || name} 
          className={labelClasses}
        >
          {label}
          {required && (
            <span className="text-red-500 ml-1" aria-label="requerido">*</span>
          )}
        </label>
      )}

      {/* Input Container */}
      <div className="relative">
        {/* Icono izquierdo */}
        {icon && iconPosition === 'left' && (
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span className={`${shouldShowError ? 'text-red-400' : 'text-gray-400'} transition-colors duration-200`}>
              {icon}
            </span>
          </div>
        )}

        {/* Input */}
        <input
          {...restProps}
          type={type === 'currency' ? 'text' : type}
          id={id || name}
          name={name}
          value={value}
          onChange={handleChange}
          onBlur={handleBlur}
          onFocus={handleFocus}
          placeholder={placeholder}
          disabled={disabled}
          readOnly={readOnly}
          autoComplete={autoComplete}
          className={inputClasses}
          aria-invalid={shouldShowError ? 'true' : 'false'}
          aria-describedby={
            [
              shouldShowError ? `${id || name}-error` : null,
              helpText ? `${id || name}-help` : null
            ].filter(Boolean).join(' ') || undefined
          }
        />

        {/* Icono derecho */}
        {icon && iconPosition === 'right' && (
          <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span className={`${shouldShowError ? 'text-red-400' : 'text-gray-400'} transition-colors duration-200`}>
              {icon}
            </span>
          </div>
        )}

        {/* Indicador de estado */}
        {shouldShowError && (
          <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <svg 
              className="h-5 w-5 text-red-500" 
              fill="currentColor" 
              viewBox="0 0 20 20"
              aria-hidden="true"
            >
              <path 
                fillRule="evenodd" 
                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" 
                clipRule="evenodd" 
              />
            </svg>
          </div>
        )}
      </div>

      {/* Texto de ayuda */}
      {helpText && !shouldShowError && (
        <p 
          id={`${id || name}-help`}
          className="mt-1 text-sm text-gray-500 dark:text-gray-400"
        >
          {helpText}
        </p>
      )}

      {/* Mensaje de error */}
      {shouldShowError && (
        <p 
          id={`${id || name}-error`}
          className={errorClasses}
          role="alert"
          aria-live="polite"
        >
          {error}
        </p>
      )}
    </div>
  );
};

export default ValidatedInput; 