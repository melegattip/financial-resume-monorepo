import React from 'react';

/**
 * Componente ResponsiveTable que renderiza tabla en desktop y cards en móvil
 * @param {Array} data - Array de objetos con los datos a mostrar
 * @param {Array} columns - Array de objetos con configuración de columnas
 * @param {Function} renderMobileCard - Función para renderizar card móvil personalizado
 * @param {string} emptyMessage - Mensaje cuando no hay datos
 * @param {string} className - Clases CSS adicionales
 */
const ResponsiveTable = ({
  data = [],
  columns = [],
  renderMobileCard = null,
  emptyMessage = "No hay datos disponibles",
  className = "",
  loading = false,
  onRowClick = null,
  showMobileCards = true,
  currentSortBy = null,
  currentSortOrder = 'asc',
  onSortChange = null,
}) => {
  
  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-fr-primary"></div>
        <span className="ml-3 text-fr-gray-600 dark:text-gray-400">Cargando...</span>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-fr-gray-400 dark:text-gray-500 text-4xl mb-4">📋</div>
        <p className="text-fr-gray-500 dark:text-gray-400 mb-6">{emptyMessage}</p>
      </div>
    );
  }

  const defaultMobileCard = (item, index) => (
    <div 
      key={index}
      className={`
        card mb-4 border-l-4 border-fr-primary transition-all duration-200
        ${onRowClick ? 'cursor-pointer hover:shadow-lg transform hover:-translate-y-1' : ''}
      `}
      onClick={() => onRowClick && onRowClick(item)}
    >
      <div className="space-y-3">
        {columns.map((column, colIndex) => {
          if (column.hideOnMobile) return null;
          
          const value = column.accessor ? item[column.accessor] : '';
          const displayValue = column.render ? column.render(value, item) : value;
          
          return (
            <div key={colIndex} className="flex justify-between items-start">
              <span className="text-sm font-medium text-fr-gray-600 dark:text-gray-400 min-w-0 flex-shrink-0 mr-3">
                {column.header}:
              </span>
              <div className="text-sm text-fr-gray-900 dark:text-gray-100 text-right min-w-0 flex-1">
                {displayValue}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );

  return (
    <div className={`responsive-table ${className}`}>
      {/* Vista Desktop - Tabla */}
      <div className="hidden lg:block">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-800">
              <tr>
                {columns.map((column, index) => {
                  const sortKey = column.sortKey || column.accessor;
                  const isActive = column.sortable && sortKey && currentSortBy === sortKey;
                  const sortIcon = !column.sortable
                    ? null
                    : isActive
                      ? (currentSortOrder === 'asc' ? '▲' : '▼')
                      : '↕';

                  const handleClick = () => {
                    if (!column.sortable || !onSortChange || !sortKey) return;
                    // Toggle order if same column, else asc
                    const nextOrder = isActive ? (currentSortOrder === 'asc' ? 'desc' : 'asc') : 'asc';
                    onSortChange(sortKey, nextOrder);
                  };

                  return (
                    <th
                      key={index}
                      onClick={handleClick}
                      className={`
                        px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 
                        uppercase tracking-wider select-none
                        ${column.align === 'right' ? 'text-right' : ''}
                        ${column.align === 'center' ? 'text-center' : ''}
                        ${column.sortable ? 'cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700' : ''}
                      `}
                    >
                      <span className="inline-flex items-center gap-1">
                        {column.header}
                        {column.sortable && (
                          <span className="text-[10px] opacity-80">{sortIcon}</span>
                        )}
                      </span>
                    </th>
                  );
                })}
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              {data.map((item, rowIndex) => (
                <tr 
                  key={rowIndex}
                  className={`
                    hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors
                    ${onRowClick ? 'cursor-pointer' : ''}
                  `}
                  onClick={() => onRowClick && onRowClick(item)}
                >
                  {columns.map((column, colIndex) => {
                    const value = column.accessor ? item[column.accessor] : '';
                    const displayValue = column.render ? column.render(value, item) : value;
                    
                    return (
                      <td
                        key={colIndex}
                        className={`
                          px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100
                          ${column.align === 'right' ? 'text-right' : ''}
                          ${column.align === 'center' ? 'text-center' : ''}
                        `}
                      >
                        {displayValue}
                      </td>
                    );
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Vista Tablet - Cards condensadas */}
      <div className="hidden md:block lg:hidden">
        <div className="space-y-4">
          {data.map((item, index) => (
            <div 
              key={index}
              className={`
                card border-l-4 border-fr-primary transition-all duration-200
                ${onRowClick ? 'cursor-pointer hover:shadow-lg' : ''}
              `}
              onClick={() => onRowClick && onRowClick(item)}
            >
              <div className="grid grid-cols-2 gap-4">
                {columns.filter(col => !col.hideOnTablet).map((column, colIndex) => {
                  const value = column.accessor ? item[column.accessor] : '';
                  const displayValue = column.render ? column.render(value, item) : value;
                  
                  return (
                    <div key={colIndex} className="min-w-0">
                      <span className="text-xs font-medium text-fr-gray-600 dark:text-gray-400 block">
                        {column.header}
                      </span>
                      <div className="text-sm text-fr-gray-900 dark:text-gray-100 mt-1">
                        {displayValue}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Vista Móvil - Cards completas */}
      {showMobileCards && (
        <div className="block md:hidden">
          <div className="space-y-4">
            {data.map((item, index) => 
              renderMobileCard ? renderMobileCard(item, index) : defaultMobileCard(item, index)
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default ResponsiveTable; 