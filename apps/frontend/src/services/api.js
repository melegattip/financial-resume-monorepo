import axios from 'axios';
import toast from '../utils/notifications';
import configService from './configService';
import envConfig from '../config/environments';

// Configuración base de axios - usa configuración dinámica
const api = axios.create({
  baseURL: envConfig.API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Función para inicializar la configuración dinámica
let configInitialized = false;

const initializeConfig = async () => {
  if (configInitialized) return;
  
  try {
    // Initializing dynamic configuration
    const config = await configService.loadConfig();
    
    // Actualizar la baseURL de axios con la configuración dinámica
    api.defaults.baseURL = config.api_base_url;
    configInitialized = true;
    
    // Dynamic configuration initialized
  } catch (error) {
    console.error('Error initializing configuration:', error);
    // Mantener la configuración por defecto
  }
};

// Inicializar configuración al cargar el módulo
initializeConfig();

// Función para obtener el token desde localStorage
const getAuthToken = () => {
  return localStorage.getItem('auth_token');
};

// Función para obtener el usuario actual desde localStorage
const getCurrentUser = () => {
  try {
    const userData = localStorage.getItem('auth_user');
    return userData ? JSON.parse(userData) : null;
  } catch (error) {
    console.error('Error parsing user data:', error);
    return null;
  }
};

// Interceptor para agregar el token JWT
api.interceptors.request.use(
  async (config) => {
    // Asegurar que la configuración esté inicializada antes de cada request
    if (!configInitialized) {
      await initializeConfig();
    }
    
    const token = getAuthToken();
    const user = getCurrentUser();
    
    // Agregar Authorization header si tenemos token
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    
    // Agregar X-Caller-ID si tenemos usuario autenticado (para compatibilidad con backend)
    if (user && user.id) {
      config.headers['X-Caller-ID'] = user.id.toString();
    }
    
    console.log('🔧 API Request:', {
      url: config.url,
      method: config.method,
      baseURL: config.baseURL,
      hasAuth: !!token,
      hasCallerId: !!(user && user.id),
      userId: user?.id,
      userIdType: typeof user?.id,
      fullUser: user,
      callerIdHeader: config.headers['X-Caller-ID']
    });
    
    return config;
  },
  (error) => {
    console.error('🔧 Request interceptor error:', error);
    return Promise.reject(error);
  }
);

// Interceptor para responses
api.interceptors.response.use(
  (response) => {
    // Log successful responses
    if (response.config.url) {
      console.log(`✅ API Response: {url: '${response.config.url}', status: ${response.status}}`);
    }
    
    return response;
  },
  (error) => {
    console.error(`❌ API Error: {url: '${error.config?.url}', status: ${error.response?.status}, message: '${error.message}'}`);
    
    // Si es un error de red sin respuesta
    if (!error.response) {
      toast.error('Error de conexión - Verifica tu conexión a internet');
      return Promise.reject(error);
    }
    
    const { status, data } = error.response;
    const message = data?.error || error.message || 'Error desconocido';
    
    if (status === 401) {
      // Para errores 401, no mostrar toast aquí ya que el authService lo maneja
      console.log('🔒 Token inválido o expirado, el authService manejará la renovación');
    } else if (status === 404) {
      toast.error('Recurso no encontrado');
    } else if (status === 500) {
      // Manejar errores 500 de manera más específica
      if (error.config?.url?.includes('/gamification/')) {
        console.warn('⚠️ Error 500 en servicio de gamificación, reintentando...');
        // No mostrar toast para errores de gamificación, ya que son temporales
      } else {
        toast.error('Error interno del servidor');
      }
    } else if (status === 409) {
      // Manejo específico para conflictos (duplicados)
      const conflictMessage = data?.message || data?.error || 'Ya existe un elemento con esos datos';
      toast.error(`Conflicto: ${conflictMessage}`);
    } else if (status >= 400 && status < 500) {
      toast.error(message);
    } else if (status >= 500) {
      toast.error('Error del servidor - Inténtalo de nuevo en unos momentos');
    } else {
      toast.error(message);
    }
    
    return Promise.reject(error);
  }
);

// Servicios de Categorías
export const categoriesAPI = {
  list: () => api.get('/categories'),
  get: (id) => api.get(`/categories/${id}`),
  create: (data) => api.post('/categories', data, { timeout: 30000 }),
  update: (id, data) => api.patch(`/categories/${id}`, data, { timeout: 30000 }),
  delete: (id) => api.delete(`/categories/${id}`, { timeout: 30000 }),
};

// Servicios de Gastos
export const expensesAPI = {
  list: (params) => api.get('/expenses', { params }),
  listUnpaid: () => api.get('/expenses/unpaid'),
  get: (userId, id) => api.get(`/expenses/${id}`),
  create: (data) => api.post('/expenses', data, { timeout: 30000 }),
  update: (userId, id, data) => api.put(`/expenses/${id}`, data, { timeout: 30000 }),
  delete: (userId, id) => api.delete(`/expenses/${id}`, { timeout: 30000 }),
};

// Servicios de Ingresos
export const incomesAPI = {
  list: () => api.get('/incomes'),
  get: (userId, id) => api.get(`/incomes/${id}`),
  create: (data) => api.post('/incomes', data, { timeout: 30000 }),
  update: (userId, id, data) => api.put(`/incomes/${id}`, data, { timeout: 30000 }),
  delete: (userId, id) => api.delete(`/incomes/${id}`, { timeout: 30000 }),
};

// Servicios de Reportes
export const reportsAPI = {
  generate: (startDate, endDate) => 
    api.get('/reports', {
      params: {
        start_date: startDate,
        end_date: endDate,
      },
    }),
};

// Servicios de Dashboard y Analytics
export const dashboardAPI = {
  overview: (params) => api.get('/dashboard', { params }),
};

export const analyticsAPI = {
  expenses: (params) => api.get('/analytics/expenses', { params }),
  incomes: (params) => api.get('/analytics/incomes', { params }),
  categories: (params) => api.get('/analytics/categories', { params }),
};

// Utilidades
export const formatCurrency = (amount) => {
  // Validar que amount sea un número válido
  const numericAmount = Number(amount);
  if (isNaN(numericAmount) || amount === null || amount === undefined) {
    return '$0,00';
  }
  
  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
  }).format(numericAmount);
};

export const formatDate = (date) => {
  // Validar que date sea válida
  if (!date || date === null || date === undefined) {
    return 'Fecha no disponible';
  }
  
  const dateObj = new Date(date);
  if (isNaN(dateObj.getTime())) {
    return 'Fecha inválida';
  }
  
  return new Intl.DateTimeFormat('es-AR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(dateObj);
};

export const formatPercentage = (percentage) => {
  // Validar que percentage sea un número válido
  const numericPercentage = Number(percentage);
  if (isNaN(numericPercentage) || percentage === null || percentage === undefined) {
    return '0.0%';
  }
  
  return `${numericPercentage.toFixed(1)}%`;
};

// Servicios de IA
export const aiAPI = {
  // Obtener insights generados por IA (backend: POST /ai/insights)
  // financialData: objeto con total_income, total_expenses, expenses_by_category, savings_goals, budgets_summary, etc.
  getInsights: async (financialData = {}) => {
    // AI requests can take up to 60 seconds — override the global 10s timeout
    const response = await api.post('/ai/insights', financialData, { timeout: 60000 });
    // Backend returns { success: true, data: [...insights], generated_at: "..." }
    const outerData = response.data || {};
    const insightsData = outerData.data || [];
    return {
      insights: Array.isArray(insightsData) ? insightsData : (insightsData.insights || []),
      generated_at: outerData.generated_at || new Date().toISOString(),
    };
  },

  // Analizar si puedes permitirte una compra
  canIBuy: async (purchaseData) => {
    const response = await api.post('/ai/can-i-buy', purchaseData);
    return response.data?.data || response.data;
  },

  // Obtener plan de mejora crediticia (backend: POST /ai/credit-plan)
  getCreditImprovementPlan: async (year = null, month = null) => {
    const body = {};
    if (year) body.year = year;
    if (month) body.month = month;

    const response = await api.post('/ai/credit-plan', body);
    return response.data?.data || response.data;
  },

  // Obtener puntuación de salud financiera
  // params: optional behavioral query params { streak, days_active, budgets_created, ... }
  getHealthScore: async (params = {}) => {
    const response = await api.get('/insights/financial-health', { params });
    return response.data?.data || response.data;
  }
};

// Servicios de Presupuestos
export const budgetsAPI = {
  list: (params) => api.get('/budgets', { params }),
  get: (id) => api.get(`/budgets/${id}`),
  create: (data) => api.post('/budgets', data),
  update: (id, data) => api.put(`/budgets/${id}`, data),
  delete: (id) => api.delete(`/budgets/${id}`),
  getStatus: (params) => api.get('/budgets/status', { params }),
  getDashboard: (params) => api.get('/budgets/dashboard', { params }),
};

// Servicios de Metas de Ahorro
export const savingsGoalsAPI = {
  list: (params) => api.get('/savings-goals', { params }),
  get: (id) => api.get(`/savings-goals/${id}`),
  create: (data) => api.post('/savings-goals', data),
  update: (id, data) => api.put(`/savings-goals/${id}`, data),
  delete: (id) => api.delete(`/savings-goals/${id}`),
  deposit: (id, data) => api.post(`/savings-goals/${id}/deposit`, data),
  withdraw: (id, data) => api.post(`/savings-goals/${id}/withdraw`, data),
  pause: (id) => api.post(`/savings-goals/${id}/pause`),
  resume: (id) => api.post(`/savings-goals/${id}/resume`),
  getDashboard: () => api.get('/savings-goals/dashboard'),
  // Historial de movimientos de una meta
  getTransactions: (id, { limit = 50, offset = 0 } = {}) =>
    api.get(`/savings-goals/${id}/transactions`, { params: { limit, offset } }),
};

// Servicios de Transacciones Recurrentes
export const recurringTransactionsAPI = {
  list: (params) => api.get('/recurring-transactions', { params }),
  get: (id) => api.get(`/recurring-transactions/${id}`),
  create: (data) => api.post('/recurring-transactions', data),
  update: (id, data) => api.put(`/recurring-transactions/${id}`, data),
  delete: (id) => api.delete(`/recurring-transactions/${id}`),
  pause: (id) => api.post(`/recurring-transactions/${id}/pause`),
  resume: (id) => api.post(`/recurring-transactions/${id}/resume`),
  execute: (id, data = {}) => api.post(`/recurring-transactions/${id}/execute`, data),
  getDashboard: () => api.get('/recurring-transactions/dashboard'),
  getProjection: (months = 6) => api.get('/recurring-transactions/projection', { 
    params: { months } 
  }),
  processPending: () => api.post('/recurring-transactions/batch/process'),
  sendNotifications: () => api.post('/recurring-transactions/batch/notify'),
};

export default api; 