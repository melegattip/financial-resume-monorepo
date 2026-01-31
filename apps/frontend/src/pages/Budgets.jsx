import React, { useState, useEffect } from 'react';
import { FaCheckCircle, FaExclamationTriangle } from 'react-icons/fa';
import { budgetsAPI, categoriesAPI, formatCurrency } from '../services/api';
import TrialBanner from '../components/TrialBanner';
import { usePeriod } from '../contexts/PeriodContext';
import { formatAmount } from '../utils/formatters';
import toast from '../utils/notifications';

const Budgets = () => {
  const { getFilterParams, balancesHidden } = usePeriod();
  const [budgets, setBudgets] = useState([]);
  const [categories, setCategories] = useState([]);
  const [dashboard, setDashboard] = useState(null);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingBudget, setEditingBudget] = useState(null);
  const [filters, setFilters] = useState({
    category_id: '',
    period: '',
    status: '',
    sort_by: 'created_at',
    sort_order: 'desc'
  });

  const [formData, setFormData] = useState({
    name: '',
    amount: '',
    category_id: '',
    period: 'monthly',
    alert_threshold: 80,
    start_date: '',
    end_date: ''
  });

  useEffect(() => {
    loadData();
  }, [filters, getFilterParams]);

  const loadData = async () => {
    try {
      setLoading(true);
      
      // Combinar filtros locales con filtros de período global
      const periodParams = getFilterParams();
      const combinedFilters = { ...filters, ...periodParams };
      
      const [budgetsRes, categoriesRes, dashboardRes] = await Promise.all([
        budgetsAPI.list(combinedFilters),
        categoriesAPI.list(),
        budgetsAPI.getDashboard(periodParams)
      ]);
      
      setBudgets(budgetsRes.data.data?.budgets || []);
      setCategories(categoriesRes.data.data || []);
      setDashboard(dashboardRes.data.data);
    } catch (error) {
      console.error('Error loading budgets:', error);
      toast.error('Error cargando presupuestos');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const data = {
        category_id: formData.category_id,
        amount: parseFloat(formData.amount),
        period: formData.period,
        alert_at: parseInt(formData.alert_threshold) / 100 // Convertir de porcentaje (0-100) a decimal (0-1)
      };

      if (editingBudget) {
        await budgetsAPI.update(editingBudget.id, data);
        toast.success('Presupuesto actualizado exitosamente');
      } else {
        await budgetsAPI.create(data);
        toast.success('Presupuesto creado exitosamente');
      }

      setShowModal(false);
      setEditingBudget(null);
      resetForm();
      loadData();
    } catch (error) {
      console.error('Error saving budget:', error);
      toast.error('Error guardando presupuesto');
    }
  };

  const handleEdit = (budget) => {
    setEditingBudget(budget);
    setFormData({
      name: budget.category_name || '',
      amount: budget.amount.toString(),
      category_id: budget.category_id || '',
      period: budget.period,
      alert_threshold: Math.round((budget.alert_at || 0.8) * 100),
      start_date: budget.period_start ? budget.period_start.split('T')[0] : '',
      end_date: budget.period_end ? budget.period_end.split('T')[0] : ''
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar este presupuesto?')) {
      try {
        await budgetsAPI.delete(id);
        toast.success('Presupuesto eliminado exitosamente');
        loadData();
      } catch (error) {
        console.error('Error deleting budget:', error);
        toast.error('Error eliminando presupuesto');
      }
    }
  };

  const resetForm = () => {
    setFormData({
      name: '',
      amount: '',
      category_id: '',
      period: 'monthly',
      alert_threshold: 80,
      start_date: '',
      end_date: ''
    });
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'on_track': return 'text-green-600 bg-green-100 dark:text-green-400 dark:bg-green-900/30';
      case 'warning': return 'text-yellow-600 bg-yellow-100 dark:text-yellow-400 dark:bg-yellow-900/30';
      case 'exceeded': return 'text-red-600 bg-red-100 dark:text-red-400 dark:bg-red-900/30';
      default: return 'text-gray-600 bg-gray-100 dark:text-gray-400 dark:bg-gray-800/50';
    }
  };

  const getStatusText = (status) => {
    switch (status) {
      case 'on_track': return 'En Meta';
      case 'warning': return 'Alerta';
      case 'exceeded': return 'Excedido';
      default: return 'Sin datos';
    }
  };

  const getCategoryName = (categoryId) => {
    const category = categories.find(c => c.id === categoryId);
    return category ? category.name : 'Sin categoría';
  };

  const getPeriodText = (period) => {
    switch (period) {
      case 'weekly': return 'Semanal';
      case 'monthly': return 'Mensual';
      case 'quarterly': return 'Trimestral';
      case 'yearly': return 'Anual';
      default: return period;
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="spinner"></div>
        <span className="ml-2 text-gray-600 dark:text-gray-400">Cargando presupuestos...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <TrialBanner featureKey="BUDGETS" />
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Presupuestos</h1>
          <p className="text-gray-600 dark:text-gray-400">Gestiona tus límites de gasto por categoría</p>
        </div>
        <button
          onClick={() => setShowModal(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600 transition-colors"
        >
          Nuevo Presupuesto
        </button>
      </div>

      {/* Dashboard Summary */}
      {dashboard && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="card">
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <div className="flex items-center space-x-6">
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Total</p>
                    <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">{dashboard.summary?.total_budgets || 0}</p>
                  </div>
                  <div className="h-12 w-px bg-fr-gray-200 dark:bg-gray-600"></div>
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">En Meta</p>
                    <p className="text-2xl font-bold text-green-600 dark:text-green-400">{dashboard.summary?.on_track_count || 0}</p>
                  </div>
                </div>
              </div>
              <div className="flex-shrink-0 p-3 rounded-fr bg-green-100 dark:bg-green-900/30 ml-4">
                <FaCheckCircle className="w-6 h-6 text-green-600 dark:text-green-400" />
              </div>
            </div>
          </div>
          <div className="card">
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <div className="flex items-center space-x-6">
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Con Alerta</p>
                    <p className="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{dashboard.summary?.warning_count || 0}</p>
                  </div>
                  <div className="h-12 w-px bg-fr-gray-200 dark:bg-gray-600"></div>
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Excedidos</p>
                    <p className="text-2xl font-bold text-red-600 dark:text-red-400">{dashboard.summary?.exceeded_count || 0}</p>
                  </div>
                </div>
              </div>
              <div className="flex-shrink-0 p-3 rounded-fr bg-yellow-100 dark:bg-yellow-900/30 ml-4">
                <FaExclamationTriangle className="w-6 h-6 text-yellow-600 dark:text-yellow-400" />
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Filters */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow border dark:border-gray-700">
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
          <select
            value={filters.category_id}
            onChange={(e) => setFilters({...filters, category_id: e.target.value})}
            className="border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
          >
            <option value="">Todas las categorías</option>
            {categories.map(category => (
              <option key={category.id} value={category.id}>{category.name}</option>
            ))}
          </select>
          
          <select
            value={filters.period}
            onChange={(e) => setFilters({...filters, period: e.target.value})}
            className="border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
          >
            <option value="">Todos los períodos</option>
            <option value="weekly">Semanal</option>
            <option value="monthly">Mensual</option>
            <option value="quarterly">Trimestral</option>
            <option value="yearly">Anual</option>
          </select>

          <select
            value={filters.status}
            onChange={(e) => setFilters({...filters, status: e.target.value})}
            className="border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
          >
            <option value="">Todos los estados</option>
            <option value="on_track">En Meta</option>
            <option value="warning">Alerta</option>
            <option value="exceeded">Excedido</option>
          </select>

          <select
            value={filters.sort_by}
            onChange={(e) => setFilters({...filters, sort_by: e.target.value})}
            className="border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
          >
            <option value="created_at">Fecha de creación</option>
            <option value="name">Nombre</option>
            <option value="amount">Monto</option>
            <option value="usage_percentage">% Usado</option>
          </select>

          <select
            value={filters.sort_order}
            onChange={(e) => setFilters({...filters, sort_order: e.target.value})}
            className="border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
          >
            <option value="desc">Descendente</option>
            <option value="asc">Ascendente</option>
          </select>
        </div>
      </div>

      {/* Budgets List */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden border dark:border-gray-700">
        {budgets.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-gray-400 dark:text-gray-500 mb-4">
              <svg className="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">No hay presupuestos</h3>
            <p className="text-gray-600 dark:text-gray-400 mb-4">Crea tu primer presupuesto para controlar tus gastos</p>
            <button
              onClick={() => setShowModal(true)}
              className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600 transition-colors"
            >
              Crear Presupuesto
            </button>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead className="bg-gray-50 dark:bg-gray-700">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Presupuesto
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Categoría
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Período
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Presupuesto
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Gastado
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Estado
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Acciones
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                {budgets.map((budget) => (
                  <tr key={budget.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>
                        <div className="text-sm font-medium text-gray-900 dark:text-gray-100">{budget.name}</div>
                        <div className="text-sm text-gray-500 dark:text-gray-400">
                          {budget.period_start ? new Date(budget.period_start).toLocaleDateString() : ''} - {budget.period_end ? new Date(budget.period_end).toLocaleDateString() : 'Sin fin'}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>
                        <div className="text-sm font-medium text-gray-900 dark:text-gray-100">{budget.category_name}</div>
                        <div className="text-sm text-gray-500 dark:text-gray-400">
                          {budget.period_start ? new Date(budget.period_start).toLocaleDateString() : ''} - {budget.period_end ? new Date(budget.period_end).toLocaleDateString() : 'Sin fin'}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                      {getPeriodText(budget.period)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                      {formatAmount(budget.amount, balancesHidden)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900 dark:text-gray-100">
                        {formatAmount(budget.spent_amount, balancesHidden)}
                      </div>
                      <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2 mt-1">
                        <div
                          className={`h-2 rounded-full ${
                            budget.spent_percentage >= 1 ? 'bg-red-500' :
                            budget.spent_percentage >= budget.alert_at ? 'bg-yellow-500' :
                            'bg-green-500'
                          }`}
                          style={{ width: `${Math.min((budget.spent_percentage || 0) * 100, 100)}%` }}
                        ></div>
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        {((budget.spent_percentage || 0) * 100).toFixed(1)}% usado
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(budget.status)}`}>
                        {getStatusText(budget.status)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button
                        onClick={() => handleEdit(budget)}
                        className="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300 mr-3"
                      >
                        Editar
                      </button>
                      <button
                        onClick={() => handleDelete(budget.id)}
                        className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
                      >
                        Eliminar
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md border dark:border-gray-700">
            <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">
              {editingBudget ? 'Editar Presupuesto' : 'Nuevo Presupuesto'}
            </h2>
            
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Nombre
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Monto
                </label>
                <input
                  type="number"
                  step="0.01"
                  value={formData.amount}
                  onChange={(e) => setFormData({...formData, amount: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Categoría
                </label>
                <select
                  value={formData.category_id}
                  onChange={(e) => setFormData({...formData, category_id: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                >
                  <option value="">Seleccionar categoría</option>
                  {categories.map(category => (
                    <option key={category.id} value={category.id}>{category.name}</option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Período
                </label>
                <select
                  value={formData.period}
                  onChange={(e) => setFormData({...formData, period: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  required
                >
                  <option value="weekly">Semanal</option>
                  <option value="monthly">Mensual</option>
                  <option value="quarterly">Trimestral</option>
                  <option value="yearly">Anual</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Alerta al (%)
                </label>
                <input
                  type="number"
                  min="1"
                  max="100"
                  value={formData.alert_threshold}
                  onChange={(e) => setFormData({...formData, alert_threshold: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Fecha de inicio
                </label>
                <input
                  type="date"
                  value={formData.start_date}
                  onChange={(e) => setFormData({...formData, start_date: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Fecha de fin (opcional)
                </label>
                <input
                  type="date"
                  value={formData.end_date}
                  onChange={(e) => setFormData({...formData, end_date: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                />
              </div>

              <div className="flex justify-end space-x-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowModal(false);
                    setEditingBudget(null);
                    resetForm();
                  }}
                  className="px-4 py-2 text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600"
                >
                  {editingBudget ? 'Actualizar' : 'Crear'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Budgets; 