import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { FaPlus, FaSearch, FaTag, FaEdit, FaTrash, FaChartBar } from 'react-icons/fa';
import { useOptimizedAPI } from '../hooks/useOptimizedAPI';
import { useGamification } from '../contexts/GamificationContext';
import { usePeriod } from '../contexts/PeriodContext';
import ValidatedInput from '../components/ValidatedInput';
import { validateCategoryName } from '../utils/validation';
import { analyticsAPI, expensesAPI, incomesAPI, formatCurrency } from '../services/api';
import dataService from '../services/dataService';
import toast from 'react-hot-toast';
import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from 'recharts';

const CHART_COLORS = ['#009ee3', '#00a650', '#ff6900', '#e53e3e', '#6b7280', '#8b5cf6', '#f59e0b', '#06b6d4', '#ec4899'];

const Categories = () => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingCategory, setEditingCategory] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [formData, setFormData] = useState({
    name: '',
    priority: '',
  });

  const [formErrors, setFormErrors] = useState({});
  const [isFormValid, setIsFormValid] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [spendingByCategory, setSpendingByCategory] = useState([]);
  const [spendingLoading, setSpendingLoading] = useState(false);

  // Monthly chart data
  const [monthlyChartData, setMonthlyChartData] = useState([]);
  const [chartLoading, setChartLoading] = useState(false);
  const [selectedChartCategories, setSelectedChartCategories] = useState([]);

  const { categories: categoriesAPI } = useOptimizedAPI();
  const { getFilterParams, getPeriodTitle, updateAvailableData } = usePeriod();
  const { recordAction } = useGamification();

  const loadCategories = useCallback(async () => {
    try {
      setLoading(true);
      const response = await categoriesAPI.list();
      const categoriesData = response.data?.data || response.data || response || [];
      setCategories(Array.isArray(categoriesData) ? categoriesData : []);
    } catch (error) {
      console.error('❌ [Categories] Error loading categories:', error);
      setCategories([]);
    } finally {
      setLoading(false);
    }
  }, [categoriesAPI]);

  const loadSpendingAnalytics = useCallback(async () => {
    setSpendingLoading(true);
    try {
      const analyticsParams = dataService.toAnalyticsDateParams(getFilterParams());
      const response = await analyticsAPI.categories(analyticsParams);
      const items = response?.data?.data || response?.data?.Categories || [];
      setSpendingByCategory(Array.isArray(items) ? items : []);
    } catch (err) {
      console.error('Error loading category spending:', err);
      setSpendingByCategory([]);
    } finally {
      setSpendingLoading(false);
    }
  }, [getFilterParams]);

  const loadMonthlyChart = useCallback(async (cats) => {
    setChartLoading(true);
    try {
      // Fetch all expenses and incomes (no period filter) for full historical evolution
      const [expResponse, incResponse] = await Promise.all([
        expensesAPI.list(),
        incomesAPI.list(),
      ]);

      const expenses = expResponse?.data?.expenses || expResponse?.data?.data || expResponse?.data || [];
      const incomes = incResponse?.data?.incomes || incResponse?.data?.data || incResponse?.data || [];

      // Feed period context so the filter dropdown has available years/months
      updateAvailableData(
        Array.isArray(expenses) ? expenses : [],
        Array.isArray(incomes) ? incomes : [],
      );

      if (!Array.isArray(expenses) || expenses.length === 0) {
        setMonthlyChartData([]);
        return;
      }

      // Build a category id→name lookup from the loaded categories
      const catMap = {};
      (cats || []).forEach(c => { catMap[c.id] = c.name; });

      // Group by month: expenses per category + total incomes as "Ingresos"
      const byMonth = {};

      const addToMonth = (rawDate, key, amount) => {
        if (!rawDate) return;
        const d = new Date(rawDate);
        if (isNaN(d.getTime())) return;
        const month = d.toLocaleDateString('es-AR', { month: 'short', year: '2-digit' });
        if (!byMonth[month]) byMonth[month] = { month };
        byMonth[month][key] = (byMonth[month][key] || 0) + (Number(amount) || 0);
      };

      expenses.forEach((exp) => {
        const rawDate = exp.transaction_date || exp.due_date || exp.created_at;
        const catName = exp.category_name || catMap[exp.category_id] || 'Sin categoría';
        addToMonth(rawDate, catName, exp.amount);
      });

      if (Array.isArray(incomes)) {
        incomes.forEach((inc) => {
          const rawDate = inc.received_date || inc.transaction_date || inc.created_at;
          addToMonth(rawDate, 'Ingresos', inc.amount);
        });
      }

      // Sort months chronologically
      const sortedMonths = Object.keys(byMonth).sort((a, b) => {
        const parse = (s) => {
          const [mon, yr] = s.split(' ');
          const monthMap = { ene: 0, feb: 1, mar: 2, abr: 3, may: 4, jun: 5, jul: 6, ago: 7, sep: 8, oct: 9, nov: 10, dic: 11 };
          return new Date(2000 + parseInt(yr), monthMap[mon.toLowerCase()] ?? 0);
        };
        return parse(a) - parse(b);
      });

      setMonthlyChartData(sortedMonths.map((m) => byMonth[m]));
    } catch (err) {
      console.error('Error loading monthly chart:', err);
      setMonthlyChartData([]);
    } finally {
      setChartLoading(false);
    }
  }, [updateAvailableData]);

  useEffect(() => {
    const init = async () => {
      // Load categories first so we can pass them to the chart for id→name lookup
      const response = await categoriesAPI.list().catch(() => null);
      const cats = response?.data?.data || response?.data || [];
      const catsArray = Array.isArray(cats) ? cats : [];
      setCategories(catsArray);
      setLoading(false);
      loadSpendingAnalytics();
      loadMonthlyChart(catsArray);
    };
    init();
  }, [loadSpendingAnalytics, loadMonthlyChart, categoriesAPI]);

  const validateForm = useCallback(() => {
    const errors = {};
    let valid = true;
    const nameValidation = validateCategoryName(formData.name);
    if (!nameValidation.isValid) {
      errors.name = nameValidation.error;
      valid = false;
    }
    setFormErrors(errors);
    setIsFormValid(valid);
    return valid;
  }, [formData]);

  useEffect(() => {
    validateForm();
  }, [validateForm]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (isSubmitting) return;
    if (!validateForm()) {
      toast.error('Por favor corrige los errores en el formulario');
      return;
    }
    setIsSubmitting(true);
    try {
      if (editingCategory) {
        await categoriesAPI.update(editingCategory.id, {
          new_name: formData.name,
          priority: parseInt(formData.priority) || 0,
        });
      } else {
        const result = await categoriesAPI.create({
          name: formData.name,
          priority: parseInt(formData.priority) || 0,
        });
        if (result?.data) {
          const categoryId = result.data.id || result.data.category_id;
          if (categoryId) {
            await recordAction('create_category', 'category', categoryId, `Created category: ${formData.name}`);
          }
        }
      }
      setShowModal(false);
      setEditingCategory(null);
      setFormData({ name: '', priority: '' });
      await loadCategories();
      await loadSpendingAnalytics();
    } catch (error) {
      console.error('Error en handleSubmit:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEdit = (category) => {
    setEditingCategory(category);
    setFormData({
      name: category.name,
      priority: category.priority > 0 ? String(category.priority) : '',
    });
    setShowModal(true);
    setIsSubmitting(false);
  };

  const handleDelete = async (category) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar esta categoría?')) {
      try {
        await categoriesAPI.delete(category.id);
        await loadCategories();
        await loadSpendingAnalytics();
      } catch (error) {
        console.error('Error en handleDelete:', error);
      }
    }
  };

  const filteredCategories = Array.isArray(categories)
    ? categories.filter(category =>
        category.name.toLowerCase().includes(searchTerm.toLowerCase())
      )
    : [];

  // All unique series names present in monthly data (expense categories + "Ingresos")
  const chartCategoryNames = monthlyChartData.length > 0
    ? [...new Set(monthlyChartData.flatMap(row => Object.keys(row).filter(k => k !== 'month')))]
    : [];

  // When chart data first loads (or reloads), select all series
  useEffect(() => {
    if (chartCategoryNames.length > 0) {
      setSelectedChartCategories(chartCategoryNames);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [monthlyChartData]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="spinner"></div>
        <span className="ml-2 text-fr-gray-600 dark:text-gray-400">Cargando categorías...</span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="card">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
          <div className="relative">
            <FaSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-fr-gray-400" />
            <input
              type="text"
              placeholder="Buscar categorías..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10 input w-full sm:w-64"
            />
          </div>
          <button
            onClick={() => { setShowModal(true); setIsSubmitting(false); }}
            className="btn-primary flex items-center space-x-2"
          >
            <FaPlus className="w-4 h-4" />
            <span>Nueva Categoría</span>
          </button>
        </div>
      </div>

      {/* Gastos por categoría — barras + acciones inline */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-base font-semibold text-fr-gray-900 dark:text-gray-100 flex items-center space-x-2">
            <FaChartBar className="w-4 h-4 text-fr-primary dark:text-blue-400" />
            <span>Gastos por categoría — {getPeriodTitle()}</span>
          </h3>
        </div>

        {spendingLoading ? (
          <div className="space-y-2">
            {[1, 2, 3].map(i => (
              <div key={i} className="h-8 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" />
            ))}
          </div>
        ) : spendingByCategory.length === 0 ? (
          <p className="text-sm text-fr-gray-500 dark:text-gray-400 py-4 text-center">
            Sin gastos registrados con categoría en este período.
          </p>
        ) : (
          <div className="space-y-3">
            {spendingByCategory.map((cat, index) => {
              const pct = cat.percentage || 0;
              const color = CHART_COLORS[index % CHART_COLORS.length];
              // Match to category object for edit/delete
              const matchedCategory = categories.find(
                c => c.id === cat.category_id || c.name === cat.category_name
              );
              return (
                <div key={cat.category_id || index}>
                  <div className="flex items-center justify-between text-sm mb-1">
                    <span className="font-medium text-fr-gray-800 dark:text-gray-200 truncate max-w-xs">
                      {cat.category_name || 'Sin nombre'}
                    </span>
                    <div className="flex items-center space-x-2 flex-shrink-0 ml-4">
                      <span className="text-fr-gray-500 dark:text-gray-400 text-xs">{pct.toFixed(1)}%</span>
                      <span className="font-semibold text-fr-gray-900 dark:text-gray-100 min-w-[80px] text-right">
                        {formatCurrency(cat.amount || 0)}
                      </span>
                      {matchedCategory && (
                        <>
                          <button
                            onClick={() => handleEdit(matchedCategory)}
                            className="p-1.5 rounded text-fr-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                            title="Editar"
                          >
                            <FaEdit className="w-3.5 h-3.5" />
                          </button>
                          <button
                            onClick={() => handleDelete(matchedCategory)}
                            className="p-1.5 rounded text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 transition-colors"
                            title="Eliminar"
                          >
                            <FaTrash className="w-3.5 h-3.5" />
                          </button>
                        </>
                      )}
                    </div>
                  </div>
                  <div className="w-full bg-gray-100 dark:bg-gray-700 rounded-full h-2 overflow-hidden">
                    <div
                      className="h-2 rounded-full transition-all duration-500"
                      style={{ width: `${Math.min(pct, 100)}%`, backgroundColor: color }}
                    />
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* Categories without spending data — show with just edit/delete */}
        {filteredCategories.length > 0 && (() => {
          const spentIds = new Set(spendingByCategory.map(c => c.category_id));
          const unspent = filteredCategories.filter(c => !spentIds.has(c.id));
          if (unspent.length === 0) return null;
          return (
            <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
              <p className="text-xs text-fr-gray-500 dark:text-gray-400 mb-2">Sin movimientos en el período</p>
              <div className="flex flex-wrap gap-2">
                {unspent.map(category => (
                  <div key={category.id} className="flex items-center space-x-1.5 bg-gray-50 dark:bg-gray-700/50 rounded-lg px-3 py-1.5">
                    <FaTag className="w-3 h-3 text-fr-primary dark:text-blue-400" />
                    <span className="text-sm text-fr-gray-700 dark:text-gray-300">{category.name}</span>
                    {category.priority > 0 && (
                      <span className="text-xs text-fr-gray-400 dark:text-gray-500">#{category.priority}</span>
                    )}
                    <button
                      onClick={() => handleEdit(category)}
                      className="p-1 rounded text-fr-gray-400 hover:text-fr-gray-600 dark:hover:text-gray-200 transition-colors"
                    >
                      <FaEdit className="w-3 h-3" />
                    </button>
                    <button
                      onClick={() => handleDelete(category)}
                      className="p-1 rounded text-red-400 hover:text-red-600 transition-colors"
                    >
                      <FaTrash className="w-3 h-3" />
                    </button>
                  </div>
                ))}
              </div>
            </div>
          );
        })()}
      </div>

      {/* Evolución mensual por categoría */}
      <div className="card">
        <h3 className="text-base font-semibold text-fr-gray-900 dark:text-gray-100 mb-4">
          Evolución de movimientos por categoría
        </h3>
        {chartLoading ? (
          <div className="h-48 bg-gray-100 dark:bg-gray-700 rounded animate-pulse" />
        ) : monthlyChartData.length === 0 ? (
          <p className="text-sm text-fr-gray-500 dark:text-gray-400 py-4 text-center">
            Sin datos suficientes para mostrar la evolución.
          </p>
        ) : (
          <>
            {/* Category filter toggles */}
            <div className="flex flex-wrap gap-2 mb-4">
              {chartCategoryNames.map((name, i) => {
                const isSelected = selectedChartCategories.includes(name);
                const color = CHART_COLORS[i % CHART_COLORS.length];
                const isIncome = name === 'Ingresos';
                return (
                  <button
                    key={name}
                    onClick={() =>
                      setSelectedChartCategories(prev =>
                        prev.includes(name) ? prev.filter(n => n !== name) : [...prev, name]
                      )
                    }
                    className="px-3 py-1 rounded-full text-xs font-medium transition-all border"
                    style={
                      isSelected
                        ? { backgroundColor: color, borderColor: color, color: '#fff' }
                        : { backgroundColor: 'transparent', borderColor: color, color: color }
                    }
                  >
                    {isIncome ? '↑ Ingresos' : name}
                  </button>
                );
              })}
            </div>

            <ResponsiveContainer width="100%" height={260}>
              <LineChart data={monthlyChartData} margin={{ top: 4, right: 16, left: 0, bottom: 4 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(107,114,128,0.2)" />
                <XAxis dataKey="month" tick={{ fontSize: 11 }} stroke="rgba(107,114,128,0.5)" />
                <YAxis tick={{ fontSize: 11 }} stroke="rgba(107,114,128,0.5)" tickFormatter={(v) => `$${v}`} />
                <Tooltip
                  formatter={(value, name) => [formatCurrency(value), name]}
                  contentStyle={{ backgroundColor: 'var(--color-bg, #1f2937)', border: '1px solid rgba(107,114,128,0.3)', borderRadius: '8px', fontSize: '12px' }}
                />
                <Legend wrapperStyle={{ fontSize: '12px' }} />
                {chartCategoryNames.map((name, i) =>
                  selectedChartCategories.includes(name) ? (
                    <Line
                      key={name}
                      type="monotone"
                      dataKey={name}
                      stroke={CHART_COLORS[i % CHART_COLORS.length]}
                      strokeWidth={name === 'Ingresos' ? 2.5 : 2}
                      strokeDasharray={name === 'Ingresos' ? '6 3' : undefined}
                      dot={{ r: 3 }}
                      activeDot={{ r: 5 }}
                    />
                  ) : null
                )}
              </LineChart>
            </ResponsiveContainer>
          </>
        )}
      </div>

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-800 rounded-fr-lg max-w-md w-full p-6">
            <h2 className="text-xl font-bold text-fr-gray-900 dark:text-gray-100 mb-6">
              {editingCategory ? 'Editar Categoría' : 'Nueva Categoría'}
            </h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <ValidatedInput
                type="text"
                name="name"
                label="Nombre"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                validator={validateCategoryName}
                validateOnChange={true}
                required={true}
                placeholder="Ej: Alimentación, Transporte, Entretenimiento"
                helpText="Nombre único para identificar la categoría"
                maxLength={50}
              />
              <div>
                <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-1">
                  Prioridad
                </label>
                <input
                  type="number"
                  min="0"
                  value={formData.priority}
                  onChange={(e) => setFormData({ ...formData, priority: e.target.value })}
                  placeholder="0 = sin prioridad"
                  className="input w-full"
                />
                <p className="text-xs text-fr-gray-500 dark:text-gray-400 mt-1">
                  Número entero — menor valor = mayor prioridad (1, 2, 3...). 0 = sin prioridad.
                </p>
              </div>
              <div className="flex space-x-4 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowModal(false);
                    setEditingCategory(null);
                    setFormData({ name: '', priority: '' });
                    setIsSubmitting(false);
                  }}
                  className="btn-outline flex-1"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className={`btn-primary flex-1 ${(!isFormValid || isSubmitting) ? 'opacity-50 cursor-not-allowed' : ''}`}
                  disabled={!isFormValid || isSubmitting}
                >
                  {isSubmitting ? 'Enviando...' : (editingCategory ? 'Actualizar' : 'Crear')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Categories;
