import apiClient from './apiClient';

// Servicios de API - SOLO comunicación HTTP
// NO contiene lógica de negocio, formateo, ni validaciones

export const categoriesService = {
  getAll: () => apiClient.get('/categories'),
  getById: (id) => apiClient.get(`/categories/${id}`),
  create: (data) => apiClient.post('/categories', data),
  update: (id, data) => apiClient.patch(`/categories/${id}`, data),
  delete: (id) => apiClient.delete(`/categories/${id}`),
};

export const expensesService = {
  getAll: () => apiClient.get('/expenses'),
  getUnpaid: () => apiClient.get('/expenses/unpaid'),
  getById: (id) => apiClient.get(`/expenses/${id}`),
  create: (data) => apiClient.post('/expenses', data),
  update: (id, data) => apiClient.patch(`/expenses/${id}`, data),
  delete: (id) => apiClient.delete(`/expenses/${id}`),
};

export const incomesService = {
  getAll: () => apiClient.get('/incomes'),
  getById: (id) => apiClient.get(`/incomes/${id}`),
  create: (data) => apiClient.post('/incomes', data),
  update: (id, data) => apiClient.patch(`/incomes/${id}`, data),
  delete: (id) => apiClient.delete(`/incomes/${id}`),
};

export const dashboardService = {
  getOverview: (params) => apiClient.get('/dashboard', { params }),
};

export const analyticsService = {
  getExpenses: (params) => apiClient.get('/analytics/expenses', { params }),
  getIncomes: (params) => apiClient.get('/analytics/incomes', { params }),
  getCategories: (params) => apiClient.get('/analytics/categories', { params }),
};

export const reportsService = {
  generate: (startDate, endDate) => 
    apiClient.get('/reports', {
      params: {
        start_date: startDate,
        end_date: endDate,
      },
    }),
};

export const authService = {
  login: (credentials) => apiClient.post('/auth/login', credentials),
  register: (userData) => apiClient.post('/auth/register', userData),
  logout: () => apiClient.post('/auth/logout'),
  getProfile: () => apiClient.get('/auth/profile'),
  refreshToken: () => apiClient.post('/auth/refresh'),
  changePassword: (passwordData) => apiClient.put('/users/security/change-password', passwordData),
}; 