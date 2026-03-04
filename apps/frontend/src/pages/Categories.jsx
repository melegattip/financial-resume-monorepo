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

// Pure helper: group raw expenses + incomes by an arbitrary date key function
const groupTransactions = (expenses, incomes, catMap, getKey) => {
  const byKey = {};
  expenses.forEach((exp) => {
    const rawDate = exp.transaction_date || exp.due_date || exp.created_at;
    if (!rawDate) return;
    // Use noon local time to avoid UTC→local day-shift (e.g. UTC-3: "2026-03-01T00:00:00Z" → Feb 28)
    const d = new Date(rawDate.split('T')[0] + 'T12:00:00');
    if (isNaN(d.getTime())) return;
    const key = getKey(d);
    const catName = exp.category_name || catMap[exp.category_id] || 'Sin categoría';
    if (!byKey[key]) byKey[key] = { month: key };
    byKey[key][catName] = (byKey[key][catName] || 0) + (Number(exp.amount) || 0);
  });
  incomes.forEach((inc) => {
    const rawDate = inc.received_date || inc.transaction_date || inc.created_at;
    if (!rawDate) return;
    const d = new Date(rawDate.split('T')[0] + 'T12:00:00');
    if (isNaN(d.getTime())) return;
    const key = getKey(d);
    if (!byKey[key]) byKey[key] = { month: key };
    byKey[key]['Ingresos'] = (byKey[key]['Ingresos'] || 0) + (Number(inc.amount) || 0);
  });
  return Object.keys(byKey).sort().map((k) => byKey[k]);
};

const monthKey = (d) =>
  `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
const dayKey = (d) =>
  `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;

const Categories = () => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingCategory, setEditingCategory] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState('priority');
  const [sortOrder, setSortOrder] = useState('asc');
  const [formData, setFormData] = useState({
    name: '',
    priority: '',
  });

  const [formErrors, setFormErrors] = useState({});
  const [isFormValid, setIsFormValid] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [spendingByCategory, setSpendingByCategory] = useState([]);
  const [spendingLoading, setSpendingLoading] = useState(false);

  // Raw chart data (no grouping yet — grouping happens in displayedChartData)
  const [allExpenses, setAllExpenses] = useState([]);
  const [allIncomes, setAllIncomes] = useState([]);
  const [chartLoading, setChartLoading] = useState(false);
  const [selectedChartCategories, setSelectedChartCategories] = useState([]);

  const { categories: categoriesAPI } = useOptimizedAPI();
  const { getFilterParams, getPeriodTitle, updateAvailableData, selectedYear, selectedMonth } = usePeriod();
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

  const loadMonthlyChart = useCallback(async () => {
    setChartLoading(true);
    try {
      const [expResponse, incResponse] = await Promise.all([
        expensesAPI.list(),
        incomesAPI.list(),
      ]);

      const expenses = expResponse?.data?.expenses || expResponse?.data?.data || expResponse?.data || [];
      const incomes = incResponse?.data?.incomes || incResponse?.data?.data || incResponse?.data || [];

      const safeExp = Array.isArray(expenses) ? expenses : [];
      const safeInc = Array.isArray(incomes) ? incomes : [];

      // Feed period context so the filter dropdown has available years/months
      updateAvailableData(safeExp, safeInc);

      // Store raw data — grouping happens in displayedChartData useMemo
      setAllExpenses(safeExp);
      setAllIncomes(safeInc);
    } catch (err) {
      console.error('Error loading chart data:', err);
      setAllExpenses([]);
      setAllIncomes([]);
    } finally {
      setChartLoading(false);
    }
  }, [updateAvailableData]);

  useEffect(() => {
    const init = async () => {
      const response = await categoriesAPI.list().catch(() => null);
      const cats = response?.data?.data || response?.data || [];
      setCategories(Array.isArray(cats) ? cats : []);
      setLoading(false);
      loadSpendingAnalytics();
      loadMonthlyChart();
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
    ? categories
        .filter(category => category.name.toLowerCase().includes(searchTerm.toLowerCase()))
        .sort((a, b) => {
          if (sortBy === 'priority') {
            const pA = a.priority || 0;
            const pB = b.priority || 0;
            if (pA === 0 && pB === 0) return a.name.localeCompare(b.name);
            if (pA === 0) return 1;
            if (pB === 0) return -1;
            return sortOrder === 'asc' ? pA - pB : pB - pA;
          } else if (sortBy === 'name') {
            return sortOrder === 'asc' ? a.name.localeCompare(b.name) : b.name.localeCompare(a.name);
          }
          return 0;
        })
    : [];

  // id→name map built from loaded categories
  const catMap = useMemo(() => {
    const map = {};
    categories.forEach((c) => { map[c.id] = c.name; });
    return map;
  }, [categories]);

  const sortedSpendingByCategory = useMemo(() => {
    let list = [...spendingByCategory];
    if (searchTerm.trim()) {
      const q = searchTerm.toLowerCase();
      list = list.filter(cat => (cat.category_name || '').toLowerCase().includes(q));
    }
    if (sortBy === 'priority') {
      return list.sort((a, b) => {
        const catA = categories.find(c => c.id === a.category_id);
        const catB = categories.find(c => c.id === b.category_id);
        const pA = catA?.priority || 0;
        const pB = catB?.priority || 0;
        if (pA === 0 && pB === 0) return 0;
        if (pA === 0) return 1;
        if (pB === 0) return -1;
        return sortOrder === 'asc' ? pA - pB : pB - pA;
      });
    } else if (sortBy === 'name') {
      return list.sort((a, b) =>
        sortOrder === 'asc'
          ? (a.category_name || '').localeCompare(b.category_name || '')
          : (b.category_name || '').localeCompare(a.category_name || '')
      );
    } else if (sortBy === 'amount') {
      return list.sort((a, b) =>
        sortOrder === 'asc' ? (a.amount || 0) - (b.amount || 0) : (b.amount || 0) - (a.amount || 0)
      );
    }
    return list;
  }, [spendingByCategory, sortBy, sortOrder, categories, searchTerm]);

  // Adaptive grouping: daily when a month is selected, monthly otherwise
  const displayedChartData = useMemo(() => {
    if (allExpenses.length === 0 && allIncomes.length === 0) return [];
    if (selectedMonth) {
      // Daily granularity, filtered to the selected month
      return groupTransactions(allExpenses, allIncomes, catMap, dayKey)
        .filter((r) => r.month.startsWith(selectedMonth + '-'));
    }
    const monthly = groupTransactions(allExpenses, allIncomes, catMap, monthKey);
    if (selectedYear) return monthly.filter((r) => r.month.startsWith(selectedYear + '-'));
    return monthly;
  }, [allExpenses, allIncomes, catMap, selectedMonth, selectedYear]);

  // All unique series names present in the filtered data (expense categories + "Ingresos")
  const chartCategoryNames = useMemo(
    () =>
      displayedChartData.length > 0
        ? [...new Set(displayedChartData.flatMap((row) => Object.keys(row).filter((k) => k !== 'month')))]
        : [],
    [displayedChartData],
  );

  // When visible data changes, reset toggles to all-selected
  useEffect(() => {
    if (chartCategoryNames.length > 0) {
      setSelectedChartCategories(chartCategoryNames);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [displayedChartData]);

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
      {/* Categorías — toolbar compacta + gastos */}
      <div className="card p-0 overflow-hidden">
        <div className="flex flex-wrap items-center gap-2 px-3 py-2 border-b border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
          <div className="relative flex-1 min-w-[160px]">
            <FaSearch className="absolute left-2.5 top-1/2 -translate-y-1/2 w-3 h-3 text-gray-400 dark:text-gray-500" />
            <input
              type="text"
              placeholder="Buscar categorías..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-8 pr-3 py-1.5 text-sm border border-gray-200 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 w-full"
            />
          </div>
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="text-sm border border-gray-200 dark:border-gray-600 rounded-lg px-2 py-1.5 bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 focus:outline-none focus:ring-1 focus:ring-blue-500"
          >
            <option value="priority">Prioridad</option>
            <option value="name">Nombre</option>
            <option value="amount">Monto</option>
          </select>
          <select
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value)}
            className="text-sm border border-gray-200 dark:border-gray-600 rounded-lg px-2 py-1.5 bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 focus:outline-none focus:ring-1 focus:ring-blue-500"
          >
            <option value="asc">↑ Asc</option>
            <option value="desc">↓ Desc</option>
          </select>
          <button
            onClick={() => { setShowModal(true); setIsSubmitting(false); }}
            className="ml-auto btn-primary flex items-center gap-1.5 py-1.5 px-3 text-sm"
          >
            <FaPlus className="w-3 h-3" />
            <span>Nueva</span>
          </button>
        </div>

        <div className="p-4">
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
            {sortedSpendingByCategory.map((cat, index) => {
              const pct = cat.percentage || 0;
              const color = CHART_COLORS[index % CHART_COLORS.length];
              // Match to category object for edit/delete
              const matchedCategory = categories.find(
                c => c.id === cat.category_id || c.name === cat.category_name
              );
              return (
                <div key={cat.category_id || index}>
                  <div className="flex items-center justify-between text-sm mb-1">
                    <div className="flex items-center gap-1.5 min-w-0 flex-1">
                      {matchedCategory?.priority > 0 && (
                        <span className="flex-shrink-0 text-xs font-mono font-semibold text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/30 px-1.5 py-0.5 rounded">
                          #{matchedCategory.priority}
                        </span>
                      )}
                      <span className="font-medium text-fr-gray-800 dark:text-gray-200 truncate">
                        {cat.category_name || 'Sin nombre'}
                      </span>
                    </div>
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
                      <span className="text-xs font-mono font-semibold text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/30 px-1.5 py-0.5 rounded">#{category.priority}</span>
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
        </div>{/* end p-4 */}
      </div>{/* end card */}

      {/* Evolución mensual por categoría */}
      <div className="card">
        <h3 className="text-base font-semibold text-fr-gray-900 dark:text-gray-100 mb-4">
          Evolución de movimientos por categoría
        </h3>
        {chartLoading ? (
          <div className="h-48 bg-gray-100 dark:bg-gray-700 rounded animate-pulse" />
        ) : displayedChartData.length === 0 ? (
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
              <LineChart data={displayedChartData} margin={{ top: 4, right: 16, left: 0, bottom: 4 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(107,114,128,0.2)" />
                <XAxis
                  dataKey="month"
                  tick={{ fontSize: 11 }}
                  stroke="rgba(107,114,128,0.5)"
                  tickFormatter={(value) => {
                    const parts = value.split('-');
                    if (parts.length === 3) {
                      // YYYY-MM-DD → "15 mar."
                      const d = new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]));
                      return d.toLocaleDateString('es-AR', { day: 'numeric', month: 'short' });
                    }
                    // YYYY-MM → "Mar. 26"
                    const d = new Date(parseInt(parts[0]), parseInt(parts[1]) - 1);
                    const label = d.toLocaleDateString('es-AR', { month: 'short', year: '2-digit' });
                    return label.charAt(0).toUpperCase() + label.slice(1);
                  }}
                />
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
                      connectNulls={true}
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
