import React, { useState, useEffect } from 'react';
import { FaCalendar, FaDownload, FaChartBar, FaArrowUp, FaArrowDown } from 'react-icons/fa';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, LineChart, Line } from 'recharts';
import { reportsAPI, formatCurrency } from '../services/api';
import { usePeriod } from '../contexts/PeriodContext';
import toast from 'react-hot-toast';

const Reports = () => {
  const [loading, setLoading] = useState(false);
  const [reportData, setReportData] = useState(null);
  const [dateRange, setDateRange] = useState({
    start_date: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0],
    end_date: new Date().toISOString().split('T')[0],
  });

  // Usar el contexto global para ocultar saldos
  const { balancesHidden } = usePeriod();

  const formatAmount = (amount) => {
    if (balancesHidden) return '••••••';
    return formatCurrency(amount);
  };

  const generateReport = React.useCallback(async () => {
    try {
      setLoading(true);
      const response = await reportsAPI.generate(dateRange.start_date, dateRange.end_date);
      setReportData(response.data);
    } catch (error) {
      console.warn('⚠️ API no disponible:', error.message);
      
      // Establecer datos vacíos
      setReportData({
        total_income: 0,
        total_expenses: 0,
        transactions: [],
        category_summary: []
      });
      
      toast.error('Error al generar el reporte', {
        duration: 3000,
      });
    } finally {
      setLoading(false);
    }
  }, [dateRange.start_date, dateRange.end_date]);

  useEffect(() => {
    generateReport();
  }, [generateReport]);

  const handleDateChange = (field, value) => {
    setDateRange(prev => ({ ...prev, [field]: value }));
  };

  // Función para calcular datos reales del gráfico de tendencia mensual
  const calculateMonthlyData = () => {
    if (!reportData || !reportData.transactions) {
      return [];
    }

    // Como no tenemos datos históricos, mostramos una comparación simple del período seleccionado
    const startDate = new Date(dateRange.start_date);
    const endDate = new Date(dateRange.end_date);
    const periodName = startDate.getMonth() === endDate.getMonth() 
      ? startDate.toLocaleDateString('es-AR', { month: 'short', year: 'numeric' })
      : `${startDate.toLocaleDateString('es-AR', { month: 'short' })} - ${endDate.toLocaleDateString('es-AR', { month: 'short' })}`;
    
    return [
      { 
        month: periodName, 
        ingresos: reportData.total_income || 0, 
        gastos: reportData.total_expenses || 0 
      }
    ];
  };

  // Función para calcular datos reales del gráfico de categorías
  const calculateCategoryData = () => {
    if (!reportData || !reportData.category_summary) {
      return [];
    }

    return reportData.category_summary.map(category => ({
      category: category.category_name,
      amount: category.total_amount,
      percentage: Math.round(category.percentage),
      transactions: reportData.transactions?.filter(t => 
        t.category_id === category.category_id && t.type === 'expense'
      ).length || 0
    }));
  };

  // Datos para gráficos usando datos reales del reporte
  const monthlyData = calculateMonthlyData();
  const categoryData = calculateCategoryData();

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="spinner"></div>
        <span className="ml-2 text-fr-gray-600 dark:text-gray-400">Generando reporte...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Title */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-fr-gray-900 dark:text-gray-100">Reportes</h1>
      </div>

      {/* Controles de fecha */}
      <div className="card">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
          <div className="flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
            <div>
              <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
                Fecha inicio
              </label>
              <input
                type="date"
                value={dateRange.start_date}
                onChange={(e) => handleDateChange('start_date', e.target.value)}
                className="input"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
                Fecha fin
              </label>
              <input
                type="date"
                value={dateRange.end_date}
                onChange={(e) => handleDateChange('end_date', e.target.value)}
                className="input"
              />
            </div>
          </div>

          <div className="flex space-x-4">
            <button
              onClick={generateReport}
              className="btn-primary flex items-center space-x-2"
            >
              <FaChartBar className="w-4 h-4" />
              <span>Generar Reporte</span>
            </button>
            <button className="btn-outline flex items-center space-x-2">
              <FaDownload className="w-4 h-4" />
              <span>Exportar</span>
            </button>
          </div>
        </div>
      </div>

      {/* Métricas principales */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="card">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <div className="flex items-center space-x-6">
                <div>
                  <p className="text-sm font-medium text-fr-gray-600 dark:text-gray-400">Total Ingresos</p>
                  <p className="text-2xl font-bold text-fr-secondary dark:text-green-400">
                    {formatAmount(reportData?.total_income || 0)}
                  </p>
                </div>
                <div className="h-12 w-px bg-fr-gray-200 dark:bg-gray-600"></div>
                <div>
                  <p className="text-sm font-medium text-fr-gray-600 dark:text-gray-400">Total Gastos</p>
                  <p className="text-2xl font-bold text-fr-gray-900 dark:text-gray-100">
                    {formatAmount(reportData?.total_expenses || 0)}
                  </p>
                </div>
              </div>
            </div>
            <div className="flex-shrink-0 p-3 rounded-fr bg-blue-100 dark:bg-blue-900/30 ml-4">
              <FaChartBar className="w-6 h-6 text-fr-primary dark:text-blue-400" />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <div className="flex items-center space-x-6">
                <div>
                  <p className="text-sm font-medium text-fr-gray-600 dark:text-gray-400">Balance</p>
                  <p className="text-2xl font-bold text-fr-secondary dark:text-green-400">
                    {formatAmount((reportData?.total_income || 0) - (reportData?.total_expenses || 0))}
                  </p>
                </div>
                <div className="h-12 w-px bg-fr-gray-200 dark:bg-gray-600"></div>
                <div>
                  <p className="text-sm font-medium text-fr-gray-600 dark:text-gray-400">Transacciones</p>
                  <p className="text-2xl font-bold text-fr-gray-900 dark:text-gray-100">
                    {reportData?.transactions?.length || 0}
                  </p>
                </div>
              </div>
            </div>
            <div className="flex-shrink-0 p-3 rounded-fr bg-purple-100 dark:bg-purple-900/30 ml-4">
              <FaCalendar className="w-6 h-6 text-purple-600 dark:text-purple-400" />
            </div>
          </div>
        </div>
      </div>

      {/* Gráficos */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Tendencia mensual */}
        <div className="card">
          <h3 className="text-lg font-semibold text-fr-gray-900 dark:text-gray-100 mb-6">Tendencia Mensual</h3>
          {monthlyData.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={monthlyData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                <XAxis dataKey="month" stroke="#6b7280" />
                <YAxis stroke="#6b7280" />
                <Tooltip 
                  formatter={(value) => [formatCurrency(value)]}
                  contentStyle={{
                    backgroundColor: 'white',
                    border: '1px solid #e5e7eb',
                    borderRadius: '8px',
                    boxShadow: '0 4px 12px 0 rgba(0, 0, 0, 0.15)',
                  }}
                />
                <Line 
                  type="monotone" 
                  dataKey="ingresos" 
                  stroke="#00a650" 
                  strokeWidth={3}
                  dot={{ fill: '#00a650', strokeWidth: 2, r: 4 }}
                />
                <Line 
                  type="monotone" 
                  dataKey="gastos" 
                  stroke="#ff6900" 
                  strokeWidth={3}
                  dot={{ fill: '#ff6900', strokeWidth: 2, r: 4 }}
                />
              </LineChart>
            </ResponsiveContainer>
          ) : (
            <div className="flex flex-col items-center justify-center h-64 text-center">
              <FaChartBar className="w-12 h-12 text-fr-gray-400 mb-4" />
              <h4 className="text-lg font-medium text-fr-gray-600 dark:text-gray-400 mb-2">No hay datos para el período</h4>
              <p className="text-sm text-fr-gray-500 dark:text-gray-500">
                Selecciona un rango de fechas con transacciones para ver la tendencia
              </p>
            </div>
          )}
        </div>

        {/* Gastos por categoría */}
        <div className="card">
          <h3 className="text-lg font-semibold text-fr-gray-900 dark:text-gray-100 mb-6">Gastos por Categoría</h3>
          {categoryData.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={categoryData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                <XAxis dataKey="category" stroke="#6b7280" />
                <YAxis stroke="#6b7280" />
                <Tooltip 
                  formatter={(value) => [formatCurrency(value)]}
                  contentStyle={{
                    backgroundColor: 'white',
                    border: '1px solid #e5e7eb',
                    borderRadius: '8px',
                    boxShadow: '0 4px 12px 0 rgba(0, 0, 0, 0.15)',
                  }}
                />
                <Bar dataKey="amount" fill="#009ee3" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <div className="flex flex-col items-center justify-center h-64 text-center">
              <FaChartBar className="w-12 h-12 text-fr-gray-400 mb-4" />
              <h4 className="text-lg font-medium text-fr-gray-600 dark:text-gray-400 mb-2">No hay gastos por categorías</h4>
              <p className="text-sm text-fr-gray-500 dark:text-gray-500">
                Los gastos con categorías aparecerán aquí una vez que generes el reporte
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Tabla de transacciones por categoría */}
      <div className="card">
        <h3 className="text-lg font-semibold text-fr-gray-900 dark:text-gray-100 mb-6">Detalle por Categoría</h3>
        {categoryData.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-fr-gray-200 dark:border-gray-700">
                  <th className="text-left py-3 px-4 font-medium text-fr-gray-900 dark:text-gray-100">Categoría</th>
                  <th className="text-right py-3 px-4 font-medium text-fr-gray-900 dark:text-gray-100">Monto</th>
                  <th className="text-right py-3 px-4 font-medium text-fr-gray-900 dark:text-gray-100">Porcentaje</th>
                  <th className="text-right py-3 px-4 font-medium text-fr-gray-900 dark:text-gray-100">Transacciones</th>
                </tr>
              </thead>
              <tbody>
                {categoryData.map((item, index) => (
                  <tr key={index} className="border-b border-fr-gray-100 dark:border-gray-700 hover:bg-fr-gray-50 dark:hover:bg-gray-700">
                    <td className="py-3 px-4 text-fr-gray-900 dark:text-gray-100">{item.category}</td>
                    <td className="py-3 px-4 text-right font-medium text-fr-gray-900 dark:text-gray-100">
                      {formatAmount(item.amount)}
                    </td>
                    <td className="py-3 px-4 text-right text-fr-gray-600 dark:text-gray-400">
                      {item.percentage}%
                    </td>
                    <td className="py-3 px-4 text-right text-fr-gray-600 dark:text-gray-400">
                      {item.transactions}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-32 text-center">
            <h4 className="text-lg font-medium text-fr-gray-600 dark:text-gray-400 mb-2">No hay datos de categorías</h4>
            <p className="text-sm text-fr-gray-500 dark:text-gray-500">
              El detalle aparecerá cuando haya gastos con categorías en el período seleccionado
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Reports; 