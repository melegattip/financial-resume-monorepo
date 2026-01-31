import { useState, useEffect } from 'react';

/**
 * Hook para debounce de valores, útil para optimizar búsquedas y filtros
 * @param {any} value - El valor a debounce
 * @param {number} delay - El delay en milisegundos
 * @returns {any} El valor debounced
 */
export function useDebounce(value, delay) {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

export default useDebounce; 