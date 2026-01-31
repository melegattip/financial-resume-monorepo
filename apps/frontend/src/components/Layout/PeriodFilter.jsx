import React, { useState, useRef, useEffect } from 'react';
import { FaEye, FaEyeSlash, FaCalendarAlt, FaChevronDown } from 'react-icons/fa';
import { usePeriod } from '../../contexts/PeriodContext';

const PeriodFilter = () => {
  const {
    selectedYear,
    selectedMonth,
    availableYears,
    balancesHidden,
    setSelectedYear,
    setSelectedMonth,
    getMonthsForSelectedYear,
    toggleBalancesVisibility,
    getPeriodTitle,
  } = usePeriod();

  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);

  const formatMonthOnly = (monthString) => {
    const [year, month] = monthString.split('-');
    const date = new Date(parseInt(year), parseInt(month) - 1, 1);
    const formatted = date.toLocaleDateString('es-AR', { 
      month: 'long' 
    });
    return formatted.charAt(0).toUpperCase() + formatted.slice(1);
  };

  // Cerrar dropdown al hacer click fuera
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handlePeriodSelect = (year, month) => {
    setSelectedYear(year);
    setSelectedMonth(month);
    setIsOpen(false);
  };

  const clearFilters = () => {
    setSelectedYear('');
    setSelectedMonth('');
    setIsOpen(false);
  };

  return (
    <div className="flex items-center space-x-3">
      {/* Botón para ocultar/mostrar saldos */}
      <button
        onClick={toggleBalancesVisibility}
        className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        title={balancesHidden ? "Mostrar saldos" : "Ocultar saldos"}
      >
        {balancesHidden ? (
          <FaEyeSlash className="w-5 h-5 text-gray-600 dark:text-gray-400" />
        ) : (
          <FaEye className="w-5 h-5 text-gray-600 dark:text-gray-400" />
        )}
      </button>

      {/* Date Picker Moderno */}
      <div className="relative" ref={dropdownRef}>
        <button
          onClick={() => setIsOpen(!isOpen)}
          className="flex items-center space-x-2 px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors text-sm font-medium text-gray-700 dark:text-gray-300 min-w-[160px]"
        >
          <FaCalendarAlt className="w-4 h-4 text-gray-500 dark:text-gray-400" />
          <span className="truncate">{getPeriodTitle()}</span>
          <FaChevronDown className={`w-3 h-3 text-gray-400 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
        </button>

        {/* Dropdown Panel */}
        {isOpen && (
          <div className="absolute top-full left-0 mt-1 w-80 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 rounded-lg shadow-lg z-50">
            <div className="p-4">
              {/* Header */}
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                  Seleccionar Período
                </h3>
                {(selectedYear || selectedMonth) && (
                  <button
                    onClick={clearFilters}
                    className="text-xs text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 font-medium"
                  >
                    Limpiar
                  </button>
                )}
              </div>

              {/* Quick Options */}
              <div className="mb-4">
                <div className="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">Acceso rápido</div>
                <div className="grid grid-cols-2 gap-2">
                  <button
                    onClick={() => handlePeriodSelect('', '')}
                    className={`p-2 text-xs rounded-md transition-colors ${
                      !selectedYear && !selectedMonth
                        ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300'
                        : 'bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600'
                    }`}
                  >
                    Todos los períodos
                  </button>
                  <button
                    onClick={() => {
                      const currentYear = new Date().getFullYear().toString();
                      handlePeriodSelect(currentYear, '');
                    }}
                    className="p-2 text-xs rounded-md bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors"
                  >
                    Este año
                  </button>
                </div>
              </div>

              {/* Year Selection */}
              <div className="mb-4">
                <div className="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">Año</div>
                <div className="grid grid-cols-4 gap-1 max-h-24 overflow-y-auto">
                  {availableYears.map(year => (
                    <button
                      key={year}
                      onClick={() => handlePeriodSelect(year, '')}
                      className={`p-2 text-xs rounded-md transition-colors ${
                        selectedYear === year && !selectedMonth
                          ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300'
                          : 'bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600'
                      }`}
                    >
                      {year}
                    </button>
                  ))}
                </div>
              </div>

              {/* Month Selection */}
              {selectedYear && (
                <div>
                  <div className="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">
                    Mes de {selectedYear}
                  </div>
                  <div className="grid grid-cols-3 gap-1 max-h-32 overflow-y-auto">
                    {getMonthsForSelectedYear().map(month => (
                      <button
                        key={month}
                        onClick={() => handlePeriodSelect(selectedYear, month)}
                        className={`p-2 text-xs rounded-md transition-colors ${
                          selectedMonth === month
                            ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300'
                            : 'bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600'
                        }`}
                      >
                        {formatMonthOnly(month)}
                      </button>
                    ))}
                  </div>
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