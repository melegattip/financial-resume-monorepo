import { dashboardAPI, analyticsAPI, expensesAPI, incomesAPI, categoriesAPI } from './api';
import toast from 'react-hot-toast';

/**
 * Servicio centralizado para manejo optimizado de datos
 */
class DataService {
  constructor() {
    this.cache = new Map();
    this.cacheTimeout = 5 * 60 * 1000; // 5 minutos
  }

  /**
   * Genera una clave de cache basada en los parámetros
   */
  getCacheKey(endpoint, params = {}) {
    const sortedParams = Object.keys(params)
      .sort()
      .reduce((result, key) => {
        result[key] = params[key];
        return result;
      }, {});
    
    return `${endpoint}_${JSON.stringify(sortedParams)}`;
  }

  /**
   * Verifica si los datos en cache son válidos
   */
  isCacheValid(cacheEntry) {
    return cacheEntry && (Date.now() - cacheEntry.timestamp) < this.cacheTimeout;
  }

  /**
   * Obtiene datos del cache o hace la llamada a la API
   */
  async getCachedData(cacheKey, apiCall) {
    const cachedData = this.cache.get(cacheKey);
    
    if (this.isCacheValid(cachedData)) {
      // Using cached data
      return cachedData.data;
    }

    // Making API call
    const data = await apiCall();
    
    // Guardar en cache
    this.cache.set(cacheKey, {
      data,
      timestamp: Date.now()
    });

    return data;
  }

  /**
   * Limpia el cache (útil después de crear/actualizar/eliminar datos)
   */
  clearCache(pattern = null) {
    if (pattern) {
      // Limpiar solo entradas que coincidan con el patrón
      let deletedCount = 0;
      for (const [key] of this.cache.entries()) {
        if (key.includes(pattern)) {
          this.cache.delete(key);
          console.log(`🗑️ [DataService] Entrada eliminada: ${key}`);
          deletedCount++;
        }
      }
      console.log(`🧹 [DataService] Cache limpiado para patrón: ${pattern} (${deletedCount} entradas)`);
    } else {
      // Limpiar todo el cache
      const totalEntries = this.cache.size;
      this.cache.clear();
      console.log(`🧹 [DataService] Cache limpiado completamente (${totalEntries} entradas)`);
    }
  }

  /**
   * Converts month/year filter params to analytics API date params (from/to ISO strings).
   * The analytics backend endpoints use parsePeriod() which expects `from` and `to`,
   * not `month` and `year`.
   */
  toAnalyticsDateParams(filterParams = {}) {
    if (filterParams.month && filterParams.year) {
      const year = parseInt(filterParams.year);
      const month = parseInt(filterParams.month);
      const from = new Date(Date.UTC(year, month - 1, 1)).toISOString();
      const to = new Date(Date.UTC(year, month, 0, 23, 59, 59)).toISOString();
      return { from, to };
    }
    if (filterParams.year) {
      const year = parseInt(filterParams.year);
      const from = new Date(Date.UTC(year, 0, 1)).toISOString();
      const to = new Date(Date.UTC(year, 11, 31, 23, 59, 59)).toISOString();
      return { from, to };
    }
    return {};
  }

  /**
   * Carga datos del dashboard de forma optimizada
   */
  async loadDashboardData(filterParams = {}, useOptimizedEndpoints = true) {
    try {
      if (useOptimizedEndpoints) {
        console.log('🚀 Cargando dashboard con endpoints optimizados...');

        // Convert month/year to from/to for analytics endpoints
        const analyticsParams = this.toAnalyticsDateParams(filterParams);

        // Llamadas paralelas con cache.
        // - /dashboard: métricas pre-calculadas del backend
        // - /expenses y /incomes: items individuales (transactions module)
        // - /analytics/categories: breakdown por categoría (usa from/to, no month/year)
        // - /categories: lista para dropdowns
        const [dashboard, expenses, incomes, categories, categoriesDropdown] = await Promise.all([
          this.getCachedData(
            this.getCacheKey('dashboard', filterParams),
            () => dashboardAPI.overview(filterParams)
          ),
          this.getCachedData(
            this.getCacheKey('expenses_legacy', {}),
            () => expensesAPI.list()
          ),
          this.getCachedData(
            this.getCacheKey('incomes_legacy', {}),
            () => incomesAPI.list()
          ),
          this.getCachedData(
            this.getCacheKey('analytics_categories', analyticsParams),
            () => analyticsAPI.categories(analyticsParams)
          ),
          this.getCachedData(
            this.getCacheKey('categories_list', {}),
            () => categoriesAPI.list()
          )
        ]);

        return this.normalizeOptimizedData(dashboard.data, expenses.data, incomes.data, categories.data, categoriesDropdown.data, filterParams);
      } else {
        return await this.loadDashboardDataLegacy(filterParams);
      }
    } catch (error) {
      console.warn('⚠️ Error con endpoints optimizados, usando fallback:', error.message);
      return await this.loadDashboardDataLegacy(filterParams);
    }
  }

  /**
   * Carga datos usando endpoints legacy (fallback)
   */
  async loadDashboardDataLegacy(filterParams = {}) {
    console.log('🔄 Cargando dashboard con endpoints legacy...');
    
    const [expenses, incomes, categories] = await Promise.all([
      this.getCachedData(
        this.getCacheKey('expenses_legacy', {}),
        () => expensesAPI.list()
      ),
      this.getCachedData(
        this.getCacheKey('incomes_legacy', {}),
        () => incomesAPI.list()
      ),
      this.getCachedData(
        this.getCacheKey('categories_legacy', {}),
        () => categoriesAPI.list()
      )
    ]);

    return this.normalizeLegacyData(expenses.data, incomes.data, categories.data, filterParams);
  }

  /**
   * Normaliza datos de endpoints optimizados.
   * - dashboard: DashboardSummary del backend (snake_case)
   * - expenses: respuesta de GET /expenses → { expenses: [...], total, limit, offset }
   * - incomes: respuesta de GET /incomes → { incomes: [...], total, limit, offset }
   * - categories: respuesta de GET /analytics/categories → { data: [...], total }
   * - categoriesDropdown: respuesta de GET /categories → array o { data: [...] }
   */
  normalizeOptimizedData(dashboard, expenses, incomes, categories, categoriesDropdown, filterParams = {}) {
    // expenses viene de GET /expenses → { expenses: [...] }
    const expensesArray = expenses?.expenses || expenses?.Expenses || [];
    const normalizedExpenses = expensesArray.map(expense => ({
      id: expense.ID || expense.id,
      user_id: expense.UserID || expense.user_id,
      amount: expense.Amount || expense.amount || 0,
      amount_paid: expense.AmountPaid ?? expense.amount_paid ?? 0,
      pending_amount: expense.PendingAmount ?? expense.pending_amount ?? 0,
      description: expense.Description || expense.description,
      category_id: expense.CategoryID || expense.category_id,
      paid: expense.Paid !== undefined ? expense.Paid : (expense.paid || false),
      due_date: expense.DueDate || expense.due_date || expense.TransactionDate || expense.transaction_date,
      percentage: expense.PercentageOfIncome || expense.percentage || 0,
      created_at: expense.CreatedAt || expense.created_at,
      updated_at: expense.UpdatedAt || expense.updated_at
    }));

    // incomes viene de GET /incomes → { incomes: [...] }
    const incomesArray = incomes?.incomes || incomes?.Incomes || [];
    const normalizedIncomes = incomesArray.map(income => ({
      id: income.ID || income.id,
      user_id: income.UserID || income.user_id,
      amount: income.Amount || income.amount || 0,
      description: income.Description || income.description,
      category_id: income.CategoryID || income.category_id,
      received_date: income.ReceivedDate || income.received_date,
      created_at: income.CreatedAt || income.created_at,
      updated_at: income.UpdatedAt || income.updated_at
    }));

    const categoriesForDropdown = categoriesDropdown?.data || categoriesDropdown || [];

    // Aplicar filtro de fecha client-side (igual que normalizeLegacyData)
    const filteredExpenses = this.filterDataByMonthAndYear(normalizedExpenses, filterParams.month, filterParams.year);
    const filteredIncomes = this.filterDataByMonthAndYear(normalizedIncomes, filterParams.month, filterParams.year);

    // Totales calculados desde los datos filtrados (respetan el filtro de fecha)
    const totalExpenses = filteredExpenses.reduce((sum, e) => sum + (e.amount || 0), 0);
    const totalIncome = filteredIncomes.reduce((sum, i) => sum + (i.amount || 0), 0);
    const balance = totalIncome - totalExpenses;

    return {
      // Métricas calculadas desde datos filtrados
      totalIncome,
      totalExpenses,
      balance,

      // Transacciones filtradas por fecha
      expenses: filteredExpenses,
      incomes: filteredIncomes,

      // Categorías para dropdown
      categories: Array.isArray(categoriesForDropdown) ? categoriesForDropdown : [],

      // Datos del /dashboard backend (para health score y otros widgets que lo consuman directamente)
      dashboardMetrics: dashboard || {},
      expensesSummary: expenses?.Summary || {},
      // GET /analytics/categories → axios response: { data: { data: [...], total: N } }
      // categories.data = response body = { data: [...], total: N }
      // categories.data.data = the actual array of CategorySummary items
      categoriesAnalytics: categories?.data?.data || categories?.data?.Categories || [],

      // Datos sin filtrar (para cálculos globales)
      allExpenses: normalizedExpenses,
      allIncomes: normalizedIncomes,

      // Metadata
      source: 'optimized',
      timestamp: Date.now()
    };
  }

  /**
   * Normaliza datos de endpoints legacy
   */
  normalizeLegacyData(expenses, incomes, categories, filterParams) {
    // Extraer datos de las respuestas
    const expensesData = expenses?.expenses || expenses || [];
    const incomesData = incomes?.incomes || incomes || [];
    const categoriesData = categories?.data || categories || [];
    
    const expensesArray = Array.isArray(expensesData) ? expensesData : [];
    const incomesArray = Array.isArray(incomesData) ? incomesData : [];
    const categoriesArray = Array.isArray(categoriesData) ? categoriesData : [];

    // Filtrar datos client-side si hay filtros activos
    const filteredExpenses = this.filterDataByMonthAndYear(expensesArray, filterParams.month, filterParams.year);
    const filteredIncomes = this.filterDataByMonthAndYear(incomesArray, filterParams.month, filterParams.year);

    const totalExpenses = filteredExpenses.reduce((sum, expense) => sum + expense.amount, 0);
    const totalIncome = filteredIncomes.reduce((sum, income) => sum + income.amount, 0);

    return {
      totalIncome,
      totalExpenses,
      balance: totalIncome - totalExpenses,
      expenses: filteredExpenses,
      incomes: filteredIncomes,
      categories: categoriesArray,
      dashboardMetrics: {},
      expensesSummary: {},
      categoriesAnalytics: [],
      
      // Datos sin filtrar
      allExpenses: expensesArray,
      allIncomes: incomesArray,
      
      // Metadata
      source: 'legacy',
      timestamp: Date.now()
    };
  }

  /**
   * Filtra datos por mes y año (client-side)
   */
  filterDataByMonthAndYear(dataArray, monthFilter, yearFilter) {
    if (!Array.isArray(dataArray)) return [];

    return dataArray.filter(item => {
      // Use the business date field; fall back to created_at only for audit purposes
      const dateStr = item.due_date || item.transaction_date || item.received_date || item.created_at;
      if (!dateStr) return true;

      const itemDate = new Date(dateStr);
      const itemYear = itemDate.getUTCFullYear();
      const itemMonth = itemDate.getUTCMonth() + 1;

      if (yearFilter && itemYear !== parseInt(yearFilter)) return false;
      if (monthFilter && itemMonth !== parseInt(monthFilter)) return false;

      return true;
    });
  }

  /**
   * Invalida cache después de operaciones CRUD
   */
  invalidateAfterMutation(type) {
    switch (type) {
      case 'expense':
        this.clearCache('expenses');
        this.clearCache('analytics_expenses');
        this.clearCache('analytics_categories');
        this.clearCache('dashboard');
        break;
      case 'income':
        this.clearCache('incomes');
        this.clearCache('analytics_incomes');
        this.clearCache('dashboard');
        break;
      case 'category':
        console.log('🗑️ [DataService] Invalidando caché de categorías...');
        this.clearCache('categories');
        this.clearCache('analytics_categories');
        this.clearCache('categories_list'); // Agregado: caché específico para lista
        console.log('✅ [DataService] Caché de categorías invalidado');
        break;
      case 'recurring_transaction':
        // Cuando se ejecuta una transacción recurrente, puede crear gastos o ingresos
        this.clearCache('expenses');
        this.clearCache('incomes');
        this.clearCache('analytics_expenses');
        this.clearCache('analytics_incomes');
        this.clearCache('dashboard');
        this.clearCache('recurring'); // Para el propio dashboard de recurrentes
        console.log('🔄 Cache invalidado después de ejecutar transacción recurrente');
        break;
      default:
        this.clearCache(); // Limpiar todo
    }
    
    // Emitir evento personalizado para notificar a componentes
    this.notifyDataChange(type);
  }

  /**
   * Emite un evento personalizado para notificar cambios de datos
   */
  notifyDataChange(type) {
    const timestamp = Date.now();
    
    // 1. Emitir evento en la ventana actual (para misma pestaña)
    const event = new CustomEvent('dataChanged', {
      detail: { type, timestamp }
    });
    window.dispatchEvent(event);
    console.log(`📡 Evento 'dataChanged' emitido para tipo: ${type}`);
    
    // 2. Guardar en localStorage para comunicación entre pestañas
    const storageEvent = {
      type,
      timestamp,
      id: Math.random().toString(36).substr(2, 9) // ID único para evitar loops
    };
    localStorage.setItem('dataChanged', JSON.stringify(storageEvent));
    console.log(`💾 Evento guardado en localStorage para comunicación entre pestañas:`, storageEvent);
    
    // 3. También emitir con un delay para casos donde el backend necesita tiempo
    setTimeout(() => {
      const delayedEvent = new CustomEvent('dataChanged', {
        detail: { type, timestamp: Date.now(), delayed: true }
      });
      window.dispatchEvent(delayedEvent);
      console.log(`📡 Evento 'dataChanged' DIFERIDO emitido para tipo: ${type}`);
      
      // También actualizar localStorage con evento diferido
      const delayedStorageEvent = {
        type,
        timestamp: Date.now(),
        delayed: true,
        id: Math.random().toString(36).substr(2, 9)
      };
      localStorage.setItem('dataChanged', JSON.stringify(delayedStorageEvent));
      console.log(`💾 Evento DIFERIDO guardado en localStorage:`, delayedStorageEvent);
    }, 1000);
  }
}

// Crear instancia singleton
const dataService = new DataService();

export default dataService; 