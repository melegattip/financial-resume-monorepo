import React, { createContext, useContext, useState, useCallback, useRef } from 'react';

const PeriodContext = createContext();

export const usePeriod = () => {
  const context = useContext(PeriodContext);
  if (!context) {
    throw new Error('usePeriod debe ser usado dentro de un PeriodProvider');
  }
  return context;
};

export const PeriodProvider = ({ children }) => {
  const [selectedYear, setSelectedYear] = useState('');
  // selectedMonth: kept for single-month compat (last toggled month or '')
  const [selectedMonth, setSelectedMonth] = useState('');
  // selectedMonths: array of 'YYYY-MM' strings for multi-month mode
  const [selectedMonths, setSelectedMonths] = useState([]);
  const [availableYears, setAvailableYears] = useState([]);
  const [availableMonths, setAvailableMonths] = useState([]);
  const [balancesHidden, setBalancesHidden] = useState(false);

  const hasAutoSelected = useRef(false);

  const updateAvailableData = useCallback((expenses = [], incomes = []) => {
    const newYears = new Set();
    const newMonths = new Set();

    const addDate = (dateStr) => {
      if (!dateStr) return;
      const date = new Date(dateStr);
      const year = date.getFullYear();
      // Skip Go zero-time values and unrealistic years
      if (!isNaN(date.getTime()) && year >= 2000 && year <= 2100) {
        newYears.add(year.toString());
        newMonths.add(date.toISOString().slice(0, 7));
      }
    };

    expenses.forEach(item => addDate(item.due_date || item.transaction_date || item.created_at));
    incomes.forEach(item => addDate(item.received_date || item.created_at));

    const currentDate = new Date();
    newYears.add(currentDate.getFullYear().toString());
    newMonths.add(currentDate.toISOString().slice(0, 7));

    setAvailableYears(prevYears => {
      const merged = Array.from(new Set([...prevYears, ...newYears])).sort().reverse();
      return JSON.stringify(prevYears) !== JSON.stringify(merged) ? merged : prevYears;
    });

    setAvailableMonths(prevMonths => {
      const merged = Array.from(new Set([...prevMonths, ...newMonths])).sort().reverse();
      return JSON.stringify(prevMonths) !== JSON.stringify(merged) ? merged : prevMonths;
    });

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
      setSelectedMonths([latestMonth]);
      setSelectedYear(latestYear);
    }
  }, []);

  const getMonthsForSelectedYear = useCallback(() => {
    if (!selectedYear) return availableMonths;
    return availableMonths.filter(month => {
      const [year] = month.split('-');
      return year === selectedYear;
    });
  }, [selectedYear, availableMonths]);

  // Toggle a month in/out of selectedMonths (multi-select)
  const toggleMonth = useCallback((month) => {
    const [year] = month.split('-');
    setSelectedYear(year);
    setSelectedMonth(month);
    setSelectedMonths(prev => {
      if (prev.includes(month)) {
        const next = prev.filter(m => m !== month);
        // If all deselected, keep the last one to avoid empty state
        return next.length === 0 ? [month] : next;
      }
      return [...prev, month].sort();
    });
  }, []);

  const clearFilters = useCallback(() => {
    setSelectedYear('');
    setSelectedMonth('');
    setSelectedMonths([]);
  }, []);

  // Returns filter params. For multi-month, returns the first month (legacy compat).
  // Components that support multi-month use selectedMonths directly.
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

  const getPeriodTitle = useCallback(() => {
    if (selectedMonths.length > 1) {
      const names = selectedMonths.map(m => {
        const [year, month] = m.split('-');
        const date = new Date(parseInt(year), parseInt(month) - 1, 1);
        return date.toLocaleDateString('es-AR', { month: 'short' });
      });
      return `${names.join(', ')} ${selectedYear}`;
    }
    if (selectedMonths.length === 1) {
      const [year, month] = selectedMonths[0].split('-');
      const date = new Date(parseInt(year), parseInt(month) - 1, 1);
      const formatted = date.toLocaleDateString('es-AR', { year: 'numeric', month: 'long' });
      return formatted.charAt(0).toUpperCase() + formatted.slice(1);
    }
    if (selectedYear) {
      return `Año ${selectedYear}`;
    }
    return 'Todos los períodos';
  }, [selectedMonths, selectedYear]);

  const hasActiveFilters = selectedMonth || selectedYear;
  const isMultiMonth = selectedMonths.length > 1;

  const toggleBalancesVisibility = useCallback(() => {
    setBalancesHidden(prev => !prev);
  }, []);

  const value = {
    selectedYear,
    selectedMonth,
    selectedMonths,
    availableYears,
    availableMonths,
    hasActiveFilters,
    isMultiMonth,
    balancesHidden,

    setSelectedYear,
    setSelectedMonth,
    setSelectedMonths,
    updateAvailableData,
    clearFilters,
    toggleMonth,
    toggleBalancesVisibility,

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
