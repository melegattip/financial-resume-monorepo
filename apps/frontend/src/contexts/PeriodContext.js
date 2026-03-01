import React, { createContext, useContext, useState, useCallback, useRef } from 'react';

// Crear el contexto
const PeriodContext = createContext();

// Hook personalizado para usar el contexto
export const usePeriod = () => {
  const context = useContext(PeriodContext);
  if (!context) {
    throw new Error('usePeriod debe ser usado dentro de un PeriodProvider');
  }
  return context;
};

// Provider del contexto
export const PeriodProvider = ({ children }) => {
  const [selectedYear, setSelectedYear] = useState('');
  const [selectedMonth, setSelectedMonth] = useState('');
  const [availableYears, setAvailableYears] = useState([]);
  const [availableMonths, setAvailableMonths] = useState([]);
  const [balancesHidden, setBalancesHidden] = useState(false);

  // Ref para garantizar que la auto-selección inicial ocurra solo una vez,
  // independientemente de cuántas páginas llamen a updateAvailableData.
  const hasAutoSelected = useRef(false);

  // Función para actualizar los datos disponibles.
  // Acumula los meses/años por UNION (no reemplaza) para que navegar entre
  // páginas no achique las opciones del filtro ni resetee la selección activa.
  const updateAvailableData = useCallback((expenses = [], incomes = []) => {
    const newYears = new Set();
    const newMonths = new Set();

    const addDate = (dateStr) => {
      if (!dateStr) return;
      const date = new Date(dateStr);
      if (!isNaN(date.getTime())) {
        newYears.add(date.getFullYear().toString());
        newMonths.add(date.toISOString().slice(0, 7));
      }
    };

    // Expenses: use transaction_date (business date), fallback to created_at (audit only)
    expenses.forEach(item => addDate(item.due_date || item.transaction_date || item.created_at));

    // Incomes: use received_date (business date), fallback to created_at (audit only)
    incomes.forEach(item => addDate(item.received_date || item.created_at));

    // Siempre incluir el año y mes actual
    const currentDate = new Date();
    newYears.add(currentDate.getFullYear().toString());
    newMonths.add(currentDate.toISOString().slice(0, 7));

    // Acumular con los valores previos (UNION) para no perder períodos de otras páginas
    setAvailableYears(prevYears => {
      const merged = Array.from(new Set([...prevYears, ...newYears])).sort().reverse();
      return JSON.stringify(prevYears) !== JSON.stringify(merged) ? merged : prevYears;
    });

    setAvailableMonths(prevMonths => {
      const merged = Array.from(new Set([...prevMonths, ...newMonths])).sort().reverse();
      return JSON.stringify(prevMonths) !== JSON.stringify(merged) ? merged : prevMonths;
    });

    // Auto-seleccionar el mes más reciente SOLO la primera vez (carga inicial).
    // No seleccionar meses futuros: usar el mes más reciente ≤ mes actual.
    // Usar ref en lugar de leer selectedMonth para evitar stale closure.
    if (!hasAutoSelected.current && newMonths.size > 0) {
      hasAutoSelected.current = true;
      const currentMonth = new Date().toISOString().slice(0, 7);
      const latestMonth =
        Array.from(newMonths)
          .filter(m => m <= currentMonth)
          .sort()
          .reverse()[0] || currentMonth;
      const [latestYear] = latestMonth.split('-');
      setSelectedMonth(latestMonth);
      setSelectedYear(latestYear);
    }
  }, []); // sin dependencias: usa ref + updater functions para evitar stale closures

  // Función para obtener meses disponibles para el año seleccionado
  const getMonthsForSelectedYear = useCallback(() => {
    if (!selectedYear) return availableMonths;
    
    return availableMonths.filter(month => {
      const [year] = month.split('-');
      return year === selectedYear;
    });
  }, [selectedYear, availableMonths]);

  // Función para limpiar filtros
  const clearFilters = useCallback(() => {
    setSelectedYear('');
    setSelectedMonth('');
  }, []);

  // Función para obtener parámetros de filtro para las APIs
  const getFilterParams = useCallback(() => {
    const params = {};
    
    if (selectedYear) params.year = selectedYear;
    if (selectedMonth) {
      const [year, month] = selectedMonth.split('-');
      params.year = year;
      params.month = month;
    }
    
    return params;
  }, [selectedYear, selectedMonth]);

  // Función para obtener el título del período seleccionado
  const getPeriodTitle = useCallback(() => {
    if (selectedMonth) {
      const [year, month] = selectedMonth.split('-');
      const date = new Date(parseInt(year), parseInt(month) - 1, 1);
      const formatted = date.toLocaleDateString('es-AR', { 
        year: 'numeric', 
        month: 'long' 
      });
      return formatted.charAt(0).toUpperCase() + formatted.slice(1);
    } else if (selectedYear) {
      return `Año ${selectedYear}`;
    }
    return 'Todos los períodos';
  }, [selectedMonth, selectedYear]);

  // Verificar si hay filtros activos
  const hasActiveFilters = selectedMonth || selectedYear;

  // Función para alternar visibilidad de saldos
  const toggleBalancesVisibility = useCallback(() => {
    setBalancesHidden(!balancesHidden);
  }, [balancesHidden]);

  const value = {
    // Estado
    selectedYear,
    selectedMonth,
    availableYears,
    availableMonths,
    hasActiveFilters,
    balancesHidden,
    
    // Acciones
    setSelectedYear,
    setSelectedMonth,
    updateAvailableData,
    clearFilters,
    toggleBalancesVisibility,
    
    // Utilidades
    getFilterParams,
    getPeriodTitle,
    getMonthsForSelectedYear,
  };

  return (
    <PeriodContext.Provider value={value}>
      {children}
    </PeriodContext.Provider>
  );
}; 