/**
 * 🎮 GAMIFICATION API SERVICE
 * 
 * Servicio para conectar con la API de gamificación real
 * Reemplaza el sistema local por llamadas al backend
 */

import apiClient from './apiClient';

class GamificationAPI {
  constructor() {
    this.baseURL = '/gamification'; // apiClient ya incluye /api/v1
  }

  // 📊 ENDPOINTS PÚBLICOS (no requieren autenticación)
  
  /**
   * Obtiene los tipos de acciones disponibles
   */
  async getActionTypes() {
    try {
      const response = await apiClient.get(`${this.baseURL}/action-types`);
      return response.data;
    } catch (error) {
      console.error('Error fetching action types:', error);
      throw error;
    }
  }

  /**
   * Obtiene información de todos los niveles
   */
  async getLevels() {
    try {
      const response = await apiClient.get(`${this.baseURL}/levels`);
      return response.data;
    } catch (error) {
      console.error('Error fetching levels:', error);
      throw error;
    }
  }

  // 🔐 ENDPOINTS PROTEGIDOS (requieren autenticación)

  /**
   * Obtiene el perfil de gamificación del usuario
   */
  async getUserProfile() {
    try {
      const response = await apiClient.get(`${this.baseURL}/profile`);
      return response.data;
    } catch (error) {
      console.error('Error fetching user profile:', error);
      throw error;
    }
  }

  /**
   * Obtiene las estadísticas detalladas del usuario
   */
  async getUserStats() {
    try {
      const response = await apiClient.get(`${this.baseURL}/stats`);
      return response.data;
    } catch (error) {
      console.error('Error fetching user stats:', error);
      throw error;
    }
  }

  /**
   * Obtiene los achievements del usuario
   */
  async getUserAchievements() {
    try {
      const response = await apiClient.get(`${this.baseURL}/achievements`);
      return response.data;
    } catch (error) {
      console.error('Error fetching user achievements:', error);
      throw error;
    }
  }

  // 🔒 FEATURE GATES ENDPOINTS

  /**
   * Obtiene todas las features del usuario (desbloqueadas y bloqueadas)
   */
  async getUserFeatures() {
    try {
      const response = await apiClient.get(`${this.baseURL}/features`);
      return response.data;
    } catch (error) {
      console.error('Error fetching user features:', error);
      throw error;
    }
  }

  /**
   * Verifica acceso a una feature específica
   * @param {string} featureKey - Clave de la feature (SAVINGS_GOALS, BUDGETS, AI_INSIGHTS)
   */
  async checkFeatureAccess(featureKey) {
    try {
      const response = await apiClient.get(`${this.baseURL}/features/${featureKey}/access`);
      return response.data;
    } catch (error) {
      // Solo mostrar error si no es 404 (endpoint no implementado)
      if (error.response?.status !== 404) {
        console.error(`Error checking access to feature ${featureKey}:`, error);
      }
      throw error;
    }
  }

  // 🏆 CHALLENGES ENDPOINTS

  /**
   * Obtiene los challenges diarios del usuario y su progreso
   */
  async getDailyChallenges() {
    try {
      const response = await apiClient.get(`${this.baseURL}/challenges/daily`);
      return response.data;
    } catch (error) {
      console.error('Error fetching daily challenges:', error);
      throw error;
    }
  }

  /**
   * Obtiene los challenges semanales del usuario y su progreso
   */
  async getWeeklyChallenges() {
    try {
      const response = await apiClient.get(`${this.baseURL}/challenges/weekly`);
      return response.data;
    } catch (error) {
      console.error('Error fetching weekly challenges:', error);
      throw error;
    }
  }

  /**
   * Procesa el progreso de challenges para una acción específica
   * @param {string} actionType - Tipo de acción
   * @param {string} entityType - Tipo de entidad (opcional)
   * @param {string} entityId - ID de la entidad (opcional)
   * @param {number} xpEarned - XP ganado por la acción
   * @param {string} description - Descripción de la acción
   */
  async processChallengeProgress(actionType, entityType = '', entityId = '', xpEarned = 0, description = '') {
    try {
      const response = await apiClient.post(`${this.baseURL}/challenges/progress`, {
        action_type: actionType,
        entity_type: entityType,
        entity_id: entityId,
        xp_earned: xpEarned,
        description
      });
      return response.data;
    } catch (error) {
      console.error('Error processing challenge progress:', error);
      throw error;
    }
  }

  /**
   * Registra una acción del usuario y otorga XP
   * @param {string} actionType - Tipo de acción (view_insight, understand_insight, etc.)
   * @param {string} entityType - Tipo de entidad (insight, suggestion, pattern, etc.)
   * @param {string} entityId - ID de la entidad
   * @param {string} description - Descripción de la acción
   */
  async recordAction(actionType, entityType, entityId, description = '') {
    try {
      // El userID se extrae automáticamente del JWT token en el backend
      // NO necesitamos enviarlo en el payload
      const response = await apiClient.post(`${this.baseURL}/actions`, {
        action_type: actionType,
        entity_type: entityType,
        entity_id: entityId,
        description: description
      });
      return response.data;
    } catch (error) {
      console.error('Error recording action:', error);
      throw error;
    }
  }



  // 🎯 MÉTODOS DE CONVENIENCIA

  /**
   * Registra que el usuario vio un insight
   */
  async recordViewInsight(insightId, description = 'User viewed insight') {
    return this.recordAction('view_insight', 'insight', insightId, description);
  }

  /**
   * Registra que el usuario entendió un insight
   */
  async recordUnderstandInsight(insightId, description = 'User understood insight') {
    return this.recordAction('understand_insight', 'insight', insightId, description);
  }

  /**
   * Registra que el usuario completó una acción
   */
  async recordCompleteAction(actionId, description = 'User completed action') {
    return this.recordAction('complete_action', 'action', actionId, description);
  }

  /**
   * Registra que el usuario vio un patrón de gastos
   */
  async recordViewPattern(patternId, description = 'User viewed spending pattern') {
    return this.recordAction('view_pattern', 'pattern', patternId, description);
  }

  /**
   * Registra que el usuario aplicó una sugerencia
   */
  async recordUseSuggestion(suggestionId, description = 'User applied suggestion') {
    return this.recordAction('use_suggestion', 'suggestion', suggestionId, description);
  }

  // 📊 ACCIONES DE NAVEGACIÓN Y TRANSACCIONES (Sistema Actualizado)

  /**
   * Registra navegación al dashboard
   */
  async recordViewDashboard() {
    return this.recordAction('view_dashboard', 'dashboard', 'main-dashboard', 'User viewed dashboard');
  }

  /**
   * Registra visualización de gastos
   */
  async recordViewExpenses() {
    return this.recordAction('view_expenses', 'expense', 'expense-list', 'User viewed expenses');
  }

  /**
   * Registra visualización de ingresos
   */
  async recordViewIncomes() {
    return this.recordAction('view_incomes', 'income', 'income-list', 'User viewed incomes');
  }

  /**
   * Registra visualización de categorías
   */
  async recordViewCategories() {
    return this.recordAction('view_categories', 'category', 'category-list', 'User viewed categories');
  }

  /**
   * Registra visualización de analytics
   */
  async recordViewAnalytics(analyticsType = 'general') {
    return this.recordAction('view_analytics', 'analytics', analyticsType, `User viewed ${analyticsType} analytics`);
  }

  // 💰 ACCIONES DE TRANSACCIONES (Motor Principal de XP)

  /**
   * Registra creación de gasto
   */
  async recordCreateExpense(expenseId, description = 'User created expense') {
    return this.recordAction('create_expense', 'expense', expenseId, description);
  }

  /**
   * Registra creación de ingreso
   */
  async recordCreateIncome(incomeId, description = 'User created income') {
    return this.recordAction('create_income', 'income', incomeId, description);
  }

  /**
   * Registra actualización de gasto
   */
  async recordUpdateExpense(expenseId, description = 'User updated expense') {
    return this.recordAction('update_expense', 'expense', expenseId, description);
  }

  /**
   * Registra actualización de ingreso
   */
  async recordUpdateIncome(incomeId, description = 'User updated income') {
    return this.recordAction('update_income', 'income', incomeId, description);
  }

  /**
   * Registra eliminación de gasto
   */
  async recordDeleteExpense(expenseId, description = 'User deleted expense') {
    return this.recordAction('delete_expense', 'expense', expenseId, description);
  }

  /**
   * Registra eliminación de ingreso
   */
  async recordDeleteIncome(incomeId, description = 'User deleted income') {
    return this.recordAction('delete_income', 'income', incomeId, description);
  }

  // 🏷️ ACCIONES DE ORGANIZACIÓN

  /**
   * Registra creación de categoría
   */
  async recordCreateCategory(categoryId, description = 'User created category') {
    return this.recordAction('create_category', 'category', categoryId, description);
  }

  /**
   * Registra actualización de categoría
   */
  async recordUpdateCategory(categoryId, description = 'User updated category') {
    return this.recordAction('update_category', 'category', categoryId, description);
  }

  /**
   * Registra asignación de categoría a transacción
   */
  async recordAssignCategory(transactionId, categoryId, description = 'User assigned category') {
    return this.recordAction('assign_category', 'transaction', transactionId, `${description} - Category: ${categoryId}`);
  }

  // 🎯 ACCIONES DE ENGAGEMENT

  /**
   * Registra login diario
   */
  async recordDailyLogin() {
    return this.recordAction('daily_login', 'user', 'daily-login', 'User daily login');
  }

  /**
   * Registra racha semanal
   */
  async recordWeeklyStreak() {
    return this.recordAction('weekly_streak', 'user', 'weekly-streak', 'User maintained weekly streak');
  }

  /**
   * Registra racha mensual
   */
  async recordMonthlyStreak() {
    return this.recordAction('monthly_streak', 'user', 'monthly-streak', 'User maintained monthly streak');
  }

  /**
   * Registra completar perfil
   */
  async recordCompleteProfile() {
    return this.recordAction('complete_profile', 'user', 'profile-complete', 'User completed profile');
  }

  // 🏆 ACCIONES DE CHALLENGES

  /**
   * Registra completar challenge diario
   */
  async recordDailyChallengeComplete(challengeId, description = 'User completed daily challenge') {
    return this.recordAction('daily_challenge_complete', 'challenge', challengeId, description);
  }

  /**
   * Registra completar challenge semanal
   */
  async recordWeeklyChallengeComplete(challengeId, description = 'User completed weekly challenge') {
    return this.recordAction('weekly_challenge_complete', 'challenge', challengeId, description);
  }

  // 🔓 ACCIONES DE FEATURES DESBLOQUEABLES

  /**
   * Registra creación de meta de ahorro (Nivel 3+)
   */
  async recordCreateSavingsGoal(goalId, description = 'User created savings goal') {
    return this.recordAction('create_savings_goal', 'savings_goal', goalId, description);
  }

  /**
   * Registra depósito en meta de ahorro
   */
  async recordDepositSavings(goalId, amount, description = 'User deposited in savings goal') {
    return this.recordAction('deposit_savings', 'savings_goal', goalId, `${description} - Amount: ${amount}`);
  }

  /**
   * Registra creación de presupuesto (Nivel 5+)
   */
  async recordCreateBudget(budgetId, description = 'User created budget') {
    return this.recordAction('create_budget', 'budget', budgetId, description);
  }

  /**
   * Registra mantenerse dentro del presupuesto
   */
  async recordStayWithinBudget(budgetId, description = 'User stayed within budget') {
    return this.recordAction('stay_within_budget', 'budget', budgetId, description);
  }

  /**
   * Registra uso de análisis IA (Nivel 7+)
   */
  async recordUseAIAnalysis(analysisId, description = 'User used AI analysis') {
    return this.recordAction('use_ai_analysis', 'ai_analysis', analysisId, description);
  }

  /**
   * Registra aplicar sugerencia de IA
   */
  async recordApplyAISuggestion(suggestionId, description = 'User applied AI suggestion') {
    return this.recordAction('apply_ai_suggestion', 'ai_suggestion', suggestionId, description);
  }
}

// 🌟 SINGLETON PATTERN
let gamificationAPIInstance = null;

export const getGamificationAPI = () => {
  if (!gamificationAPIInstance) {
    gamificationAPIInstance = new GamificationAPI();
  }
  return gamificationAPIInstance;
};

export default GamificationAPI; 