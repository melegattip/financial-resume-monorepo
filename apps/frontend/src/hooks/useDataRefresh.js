import { useEffect, useCallback, useRef } from 'react';

/**
 * Hook personalizado para escuchar cambios de datos y refrescar automáticamente
 */
export const useDataRefresh = (refreshFunction, dataTypes = ['expense', 'income', 'recurring_transaction']) => {
  const processedEventsRef = useRef(new Set());

  const handleDataChange = useCallback((event) => {
    const { type } = event.detail;
    
    console.log(`📡 Evento 'dataChanged' recibido. Tipo: ${type}, DataTypes esperados:`, dataTypes);
    
    // Solo refrescar si el tipo de dato nos interesa
    if (dataTypes.includes(type)) {
      console.log(`🔄 Refrescando datos debido a cambio en: ${type}`);
      if (typeof refreshFunction === 'function') {
        console.log(`🚀 Ejecutando función de refresh...`);
        refreshFunction();
      } else {
        console.warn(`⚠️ refreshFunction no es una función válida:`, refreshFunction);
      }
    } else {
      console.log(`ℹ️ Tipo '${type}' no está en la lista de tipos esperados, ignorando evento`);
    }
  }, [refreshFunction, dataTypes]);

  // Manejar cambios en localStorage (para comunicación entre pestañas)
  const handleStorageChange = useCallback((event) => {
    if (event.key === 'dataChanged' && event.newValue) {
      try {
        const data = JSON.parse(event.newValue);
        const { type, id } = data;
        
        // Evitar procesar el mismo evento múltiples veces
        if (processedEventsRef.current.has(id)) {
          return;
        }
        processedEventsRef.current.add(id);
        
        console.log(`💾 Evento localStorage recibido entre pestañas. Tipo: ${type}, ID: ${id}`);
        
        // Solo refrescar si el tipo de dato nos interesa
        if (dataTypes.includes(type)) {
          console.log(`🔄 Refrescando datos debido a cambio entre pestañas: ${type}`);
          if (typeof refreshFunction === 'function') {
            console.log(`🚀 Ejecutando función de refresh desde localStorage...`);
            refreshFunction();
          }
        }
        
        // Limpiar IDs antiguos para evitar memory leaks
        if (processedEventsRef.current.size > 50) {
          processedEventsRef.current.clear();
        }
      } catch (error) {
        console.warn('Error parsing localStorage dataChanged event:', error);
      }
    }
  }, [refreshFunction, dataTypes]);

  useEffect(() => {
    // Escuchar el evento personalizado (misma pestaña)
    window.addEventListener('dataChanged', handleDataChange);
    
    // Escuchar cambios en localStorage (entre pestañas)
    window.addEventListener('storage', handleStorageChange);
    
    // Cleanup
    return () => {
      window.removeEventListener('dataChanged', handleDataChange);
      window.removeEventListener('storage', handleStorageChange);
    };
  }, [handleDataChange, handleStorageChange]);
};

export default useDataRefresh; 