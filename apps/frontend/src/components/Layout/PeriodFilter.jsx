import React, { useState, useRef, useEffect } from 'react';
import { FaEye, FaEyeSlash, FaCalendarAlt, FaChevronDown, FaCheck } from 'react-icons/fa';
import { usePeriod } from '../../contexts/PeriodContext';

const PeriodFilter = ({ compact = false }) => {
  const {
    selectedYear,
    selectedMonth,
    selectedMonths,
    availableYears,
    availableMonths,
    balancesHidden,
    setSelectedYear,
    setSelectedMonth,
    setSelectedMonths,
    getMonthsForSelectedYear,
    toggleBalancesVisibility,
    toggleMonth,
    clearFilters,
    getPeriodTitle,
  } = usePeriod();

  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);

  const formatMonthOnly = (monthString) => {
    const [year, month] = monthString.split('-');
    const date = new Date(parseInt(year), parseInt(month) - 1, 1);
    const formatted = date.toLocaleDateString('es-AR', { month: 'long' });
    return formatted.charAt(0).toUpperCase() + formatted.slice(1);
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleYearSelect = (year) => {
    setSelectedYear(year);
    setSelectedMonth('');
    setSelectedMonths([]);
  };

  const handleMonthToggle = (month) => {
    toggleMonth(month);
  };

  const handleAllPeriods = () => {
    clearFilters();
    setIsOpen(false);
  };

  const handleThisYear = () => {
    const currentYear = new Date().getFullYear().toString();
    const yearMonths = availableMonths.filter(m => m.startsWith(currentYear));
    setSelectedYear(currentYear);
    setSelectedMonths(yearMonths.length > 0 ? yearMonths : []);
    setSelectedMonth(yearMonths[0] || ''); // most recent month (array is descending)
    setIsOpen(false);
  };

  const multiCount = selectedMonths.length;

  return (
    <div className={`flex items-center ${compact ? 'space-x-1' : 'space-x-3'}`}>
      {/* Toggle balance visibility */}
      <button
        onClick={toggleBalancesVisibility}
        className="p-1.5 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        title={balancesHidden ? 'Mostrar saldos' : 'Ocultar saldos'}
      >
        {balancesHidden ? (
          <FaEyeSlash className={`${compact ? 'w-4 h-4' : 'w-5 h-5'} text-gray-600 dark:text-gray-400`} />
        ) : (
          <FaEye className={`${compact ? 'w-4 h-4' : 'w-5 h-5'} text-gray-600 dark:text-gray-400`} />
        )}
      </button>

      {/* Date picker dropdown */}
      <div className="relative" ref={dropdownRef}>
        <button
          onClick={() => setIsOpen(!isOpen)}
          className={compact
            ? 'p-1.5 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors'
            : 'flex items-center space-x-2 px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors text-sm font-medium text-gray-700 dark:text-gray-300 min-w-[160px]'
          }
          title={compact ? getPeriodTitle() : undefined}
        >
          {compact ? (
            <div className="relative">
              <FaCalendarAlt className="w-4 h-4 text-gray-600 dark:text-gray-400" />
              {multiCount > 1 && (
                <span className="absolute -top-1.5 -right-1.5 bg-blue-500 text-white text-[9px] rounded-full w-3.5 h-3.5 flex items-center justify-center font-bold">
                  {multiCount}
                </span>
              )}
            </div>
          ) : (
            <>
              <FaCalendarAlt className="w-4 h-4 text-gray-500 dark:text-gray-400 flex-shrink-0" />
              <span className="truncate flex-1">{getPeriodTitle()}</span>
              {multiCount > 1 && (
                <span className="bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 text-xs px-1.5 py-0.5 rounded-full font-medium">
                  {multiCount}
                </span>
              )}
              <FaChevronDown className={`w-3 h-3 text-gray-400 transition-transform flex-shrink-0 ${isOpen ? 'rotate-180' : ''}`} />
            </>
          )}
        </button>

        {isOpen && (
          <div className="absolute top-full left-0 mt-1 w-80 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 rounded-lg shadow-lg z-50">
            <div className="p-4">
              {/* Header */}
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                  Seleccionar Período
                </h3>
                <div className="flex items-center gap-2">
                  {multiCount > 1 && (
                    <span className="text-xs text-blue-600 dark:text-blue-400 font-medium">
                      {multiCount} meses
                    </span>
                  )}
                  {(selectedYear || selectedMonth || multiCount > 0) && (
                    <button
                      onClick={handleAllPeriods}
                      className="text-xs text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 font-medium"
                    >
                      Limpiar
                    </button>
                  )}
                </div>
              </div>

              {/* Quick options */}
              <div className="mb-3">
                <div className="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">Acceso rápido</div>
                <div className="grid grid-cols-2 gap-2">
                  <button
                    onClick={handleAllPeriods}
                    className={`p-2 text-xs rounded-md transition-colors ${
                      !selectedYear && !selectedMonth && multiCount === 0
                        ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300'
                        : 'bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600'
                    }`}
                  >
                    Todos los períodos
                  </button>
                  <button
                    onClick={handleThisYear}
                    className="p-2 text-xs rounded-md bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors"
                  >
                    Este año
                  </button>
                </div>
              </div>

              {/* Year selection */}
              <div className="mb-3">
                <div className="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">Año</div>
                <div className="grid grid-cols-4 gap-1">
                  {availableYears.map(year => (
                    <button
                      key={year}
                      onClick={() => handleYearSelect(year)}
                      className={`p-2 text-xs rounded-md transition-colors ${
                        selectedYear === year && multiCount === 0
                          ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300'
                          : selectedYear === year
                          ? 'bg-blue-50 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400 border border-blue-200 dark:border-blue-700'
                          : 'bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600'
                      }`}
                    >
                      {year}
                    </button>
                  ))}
                </div>
              </div>

              {/* Month multi-select */}
              {selectedYear && (
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <div className="text-xs font-medium text-gray-500 dark:text-gray-400">
                      Meses de {selectedYear}
                    </div>
                    <div className="text-xs text-gray-400 dark:text-gray-500">
                      Click = multi-selección
                    </div>
                  </div>
                  <div className="grid grid-cols-3 gap-1">
                    {getMonthsForSelectedYear().map(month => {
                      const isSelected = selectedMonths.includes(month);
                      return (
                        <button
                          key={month}
                          onClick={() => handleMonthToggle(month)}
                          className={`p-2 text-xs rounded-md transition-colors flex items-center justify-between gap-1 ${
                            isSelected
                              ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 font-medium'
                              : 'bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600'
                          }`}
                        >
                          <span className="truncate">{formatMonthOnly(month)}</span>
                          {isSelected && <FaCheck className="w-2.5 h-2.5 flex-shrink-0" />}
                        </button>
                      );
                    })}
                  </div>
                  {multiCount > 1 && (
                    <p className="text-xs text-blue-600 dark:text-blue-400 mt-2 text-center">
                      Vista comparativa de {multiCount} meses activa
                    </p>
                  )}
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default PeriodFilter;
