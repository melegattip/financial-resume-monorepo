/**
 * 🎮 GAMIFICATION SERVICE - Arquitectura Escalable
 * 
 * Este servicio está diseñado para escalar desde MVP hasta 10M+ usuarios
 * Incluye: Points, Achievements, Levels, Analytics, Persistence
 */

import { getGamificationAPI } from './gamificationAPI';

// 🎯 CONFIGURACIÓN ESCALABLE
const GAMIFICATION_CONFIG = {
  // Sistema de niveles exponencial
  levels: [
    { level: 1, xpRequired: 0, title: "Financial Newbie", color: "#94A3B8" },
    { level: 2, xpRequired: 100, title: "Budget Learner", color: "#60A5FA" },
    { level: 3, xpRequired: 250, title: "Expense Tracker", color: "#34D399" },
    { level: 4, xpRequired: 500, title: "Savings Explorer", color: "#FBBF24" },
    { level: 5, xpRequired: 1000, title: "Investment Seeker", color: "#F472B6" },
    { level: 6, xpRequired: 2000, title: "Financial Analyst", color: "#A78BFA" },
    { level: 7, xpRequired: 4000, title: "Money Master", color: "#FB7185" },
    { level: 8, xpRequired: 8000, title: "Wealth Builder", color: "#FBBF24" },
    { level: 9, xpRequired: 15000, title: "Financial Guru", color: "#10B981" },
    { level: 10, xpRequired: 30000, title: "Money Magnate", color: "#8B5CF6" }
  ],
  
  // Sistema de puntos por acción
  points: {
    // Insights interactions
    viewInsight: 1,
    understandInsight: 3,
    actionTaken: 10,
    goalAchieved: 25,
    
    // Financial behaviors
    addTransaction: 2,
    categorizeExpense: 3,
    setGoal: 15,
    achieveGoal: 50,
    
    // Engagement
    dailyLogin: 5,
    weeklyStreak: 20,
    monthlyStreak: 100,
    
    // Social & Educational
    shareAchievement: 10,
    completeLesson: 25,
    helpOtherUser: 15
  },
  
  // Achievements escalables por categoría
  achievements: {
    insights: [
      { id: 'first-insight', name: '🧠 First Insight', description: 'Ver tu primer insight de IA', points: 50, target: 1 },
      { id: 'insight-explorer', name: '📊 Insight Explorer', description: 'Ver 25 insights únicos', points: 200, target: 25 },
      { id: 'ai-partner', name: '🤖 AI Partner', description: 'Utilizar 100 insights de IA', points: 500, target: 100 },
      { id: 'insight-master', name: '💡 Insight Master', description: 'Actuar sobre 50 insights', points: 1000, target: 50 }
    ],
    
    financial: [
      { id: 'budget-ninja', name: '🎯 Budget Ninja', description: '3 meses sin exceder presupuesto', points: 750, target: 3 },
      { id: 'savings-hero', name: '💰 Savings Hero', description: 'Ahorrar $10,000 ARS', points: 1000, target: 10000 },
      { id: 'category-master', name: '📂 Category Master', description: '90% transacciones categorizadas', points: 500, target: 90 },
      { id: 'goal-crusher', name: '🏆 Goal Crusher', description: 'Alcanzar 5 metas financieras', points: 1250, target: 5 }
    ],
    
    engagement: [
      { id: 'daily-warrior', name: '⚡ Daily Warrior', description: '7 días consecutivos activo', points: 300, target: 7 },
      { id: 'monthly-champion', name: '🔥 Monthly Champion', description: '30 días consecutivos', points: 1500, target: 30 },
      { id: 'data-enthusiast', name: '📈 Data Enthusiast', description: '100 transacciones registradas', points: 400, target: 100 }
    ]
  }
};

// 🎮 CLASE PRINCIPAL ESCALABLE
class GamificationService {
  constructor() {
    this.storageKey = 'financial_gamification';
    this.apiEndpoint = '/api/v1/gamification'; // Para futura integración backend
    this.eventListeners = new Map();
    this.analyticsQueue = [];
    
    // Inicializar datos del usuario
    this.userData = this.loadUserData();
    
    // Setup analytics batching (escalable)
    this.setupAnalyticsBatching();

    this.api = getGamificationAPI();
  }
  
  // 📊 SISTEMA DE PUNTOS ESCALABLE
  addPoints(action, metadata = {}) {
    const points = GAMIFICATION_CONFIG.points[action] || 0;
    if (points === 0) return null;
    
    const previousLevel = this.getCurrentLevel();
    this.userData.totalXP += points;
    this.userData.lastActivity = new Date().toISOString();
    
    // Track specific action
    this.trackAction(action, points, metadata);
    
    // Check for level up
    const newLevel = this.getCurrentLevel();
    const leveledUp = newLevel.level > previousLevel.level;
    
    // Check for new achievements
    const newAchievements = this.checkForNewAchievements(action, metadata);
    
    // Save data
    this.saveUserData();
    
    // Analytics
    this.queueAnalytics('points_earned', {
      action,
      points,
      totalXP: this.userData.totalXP,
      level: newLevel.level,
      metadata
    });
    
    // Trigger events
    this.triggerEvent('pointsEarned', { action, points, leveledUp, newAchievements });
    
    return {
      pointsEarned: points,
      totalXP: this.userData.totalXP,
      leveledUp,
      newLevel: leveledUp ? newLevel : null,
      newAchievements
    };
  }
  
  // 🏆 SISTEMA DE ACHIEVEMENTS ESCALABLE
  checkForNewAchievements(action, metadata) {
    const newAchievements = [];
    
    // Check all achievement categories
    Object.values(GAMIFICATION_CONFIG.achievements).flat().forEach(achievement => {
      if (this.userData.achievements.includes(achievement.id)) return;
      
      if (this.isAchievementUnlocked(achievement, action, metadata)) {
        this.unlockAchievement(achievement);
        newAchievements.push(achievement);
      }
    });
    
    return newAchievements;
  }
  
  isAchievementUnlocked(achievement, action, metadata) {
    switch (achievement.id) {
      case 'first-insight':
        return action === 'viewInsight' && this.userData.stats.insightsViewed >= 1;
      
      case 'insight-explorer':
        return this.userData.stats.insightsViewed >= 25;
      
      case 'ai-partner':
        return this.userData.stats.insightsViewed >= 100;
      
      case 'insight-master':
        return this.userData.stats.actionsCompleted >= 50;
      
      case 'daily-warrior':
        return this.userData.stats.currentStreak >= 7;
      
      case 'data-enthusiast':
        return this.userData.stats.transactionsAdded >= 100;
      
      // Add more achievement logic as needed
      default:
        return false;
    }
  }
  
  unlockAchievement(achievement) {
    this.userData.achievements.push(achievement.id);
    this.userData.totalXP += achievement.points;
    this.userData.achievementsUnlockedAt[achievement.id] = new Date().toISOString();
    
    // Analytics
    this.queueAnalytics('achievement_unlocked', {
      achievementId: achievement.id,
      achievementName: achievement.name,
      points: achievement.points,
      totalAchievements: this.userData.achievements.length
    });
    
    // Trigger event
    this.triggerEvent('achievementUnlocked', achievement);
  }
  
  // 📈 SISTEMA DE NIVELES
  getCurrentLevel() {
    const currentXP = this.userData.totalXP;
    
    for (let i = GAMIFICATION_CONFIG.levels.length - 1; i >= 0; i--) {
      const level = GAMIFICATION_CONFIG.levels[i];
      if (currentXP >= level.xpRequired) {
        const nextLevel = GAMIFICATION_CONFIG.levels[i + 1];
        
        return {
          ...level,
          progress: nextLevel ? 
            ((currentXP - level.xpRequired) / (nextLevel.xpRequired - level.xpRequired)) * 100 : 100,
          xpToNext: nextLevel ? nextLevel.xpRequired - currentXP : 0,
          nextLevel: nextLevel || null
        };
      }
    }
    
    return GAMIFICATION_CONFIG.levels[0];
  }
  
  // 📊 ANALYTICS ESCALABLES (Batch processing)
  setupAnalyticsBatching() {
    // Send analytics every 30 seconds or when queue reaches 10 items
    setInterval(() => {
      if (this.analyticsQueue.length > 0) {
        this.flushAnalytics();
      }
    }, 30000);
  }
  
  queueAnalytics(event, data) {
    this.analyticsQueue.push({
      event,
      data,
      timestamp: new Date().toISOString(),
      userId: this.userData.userId || 'anonymous',
      sessionId: this.getSessionId()
    });
    
    // Flush immediately if queue is full
    if (this.analyticsQueue.length >= 10) {
      this.flushAnalytics();
    }
  }
  
  async flushAnalytics() {
    if (this.analyticsQueue.length === 0) return;
    
    const events = [...this.analyticsQueue];
    this.analyticsQueue = [];
    
    try {
      // Future: Send to backend analytics
      // await fetch(`${this.apiEndpoint}/analytics`, {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify({ events })
      // });
      
      // For now: Store locally for development
      const existingAnalytics = JSON.parse(localStorage.getItem('gamification_analytics') || '[]');
      existingAnalytics.push(...events);
      localStorage.setItem('gamification_analytics', JSON.stringify(existingAnalytics.slice(-1000))); // Keep last 1000 events
      
    } catch (error) {
      console.error('Failed to send gamification analytics:', error);
      // Re-queue events on failure
      this.analyticsQueue.unshift(...events);
    }
  }
  
  // 📱 GESTIÓN DE DATOS ESCALABLE
  loadUserData() {
    try {
      const stored = localStorage.getItem(this.storageKey);
      if (stored) {
        const data = JSON.parse(stored);
        // Migrate old data structure if needed
        return this.migrateUserData(data);
      }
    } catch (error) {
      console.error('Failed to load gamification data:', error);
    }
    
    // Default data structure
    return this.getDefaultUserData();
  }
  
  getDefaultUserData() {
    return {
      version: '1.0.0', // For future migrations
      userId: null,
      totalXP: 0,
      level: 1,
      achievements: [],
      achievementsUnlockedAt: {},
      stats: {
        insightsViewed: 0,
        actionsCompleted: 0,
        goalsAchieved: 0,
        transactionsAdded: 0,
        currentStreak: 0,
        longestStreak: 0,
        lastLoginDate: null
      },
      preferences: {
        notifications: true,
        celebrationAnimations: true,
        shareAchievements: false
      },
      createdAt: new Date().toISOString(),
      lastActivity: new Date().toISOString()
    };
  }
  
  migrateUserData(data) {
    // Handle future data migrations
    if (!data.version) {
      data.version = '1.0.0';
      // Add any migration logic here
    }
    
    return { ...this.getDefaultUserData(), ...data };
  }
  
  saveUserData() {
    try {
      localStorage.setItem(this.storageKey, JSON.stringify(this.userData));
    } catch (error) {
      console.error('Failed to save gamification data:', error);
    }
  }
  
  // 🎯 TRACKING DE ACCIONES
  trackAction(action, points, metadata) {
    switch (action) {
      case 'viewInsight':
        this.userData.stats.insightsViewed++;
        break;
      case 'actionTaken':
        this.userData.stats.actionsCompleted++;
        break;
      case 'goalAchieved':
        this.userData.stats.goalsAchieved++;
        break;
      case 'addTransaction':
        this.userData.stats.transactionsAdded++;
        break;
    }
    
    // Update streak
    this.updateLoginStreak();
  }
  
  updateLoginStreak() {
    const today = new Date().toDateString();
    const lastLogin = this.userData.stats.lastLoginDate;
    
    if (lastLogin === today) return; // Already logged in today
    
    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);
    
    if (lastLogin === yesterday.toDateString()) {
      // Consecutive day
      this.userData.stats.currentStreak++;
    } else if (lastLogin !== null) {
      // Streak broken
      this.userData.stats.currentStreak = 1;
    } else {
      // First login
      this.userData.stats.currentStreak = 1;
    }
    
    // Update longest streak
    if (this.userData.stats.currentStreak > this.userData.stats.longestStreak) {
      this.userData.stats.longestStreak = this.userData.stats.currentStreak;
    }
    
    this.userData.stats.lastLoginDate = today;
  }
  
  // 🎪 EVENT SYSTEM ESCALABLE
  addEventListener(event, callback) {
    if (!this.eventListeners.has(event)) {
      this.eventListeners.set(event, []);
    }
    this.eventListeners.get(event).push(callback);
    
    // Return unsubscribe function
    return () => {
      const listeners = this.eventListeners.get(event);
      const index = listeners.indexOf(callback);
      if (index > -1) listeners.splice(index, 1);
    };
  }
  
  triggerEvent(event, data) {
    const listeners = this.eventListeners.get(event) || [];
    listeners.forEach(callback => {
      try {
        callback(data);
      } catch (error) {
        console.error(`Error in gamification event listener for ${event}:`, error);
      }
    });
  }
  
  // 🔧 UTILIDADES
  getSessionId() {
    if (!this.sessionId) {
      this.sessionId = 'session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
    }
    return this.sessionId;
  }
  
  // 📊 GETTERS PÚBLICOS
  getUserStats() {
    return {
      ...this.userData.stats,
      totalXP: this.userData.totalXP,
      level: this.getCurrentLevel(),
      achievements: this.userData.achievements.map(id => 
        Object.values(GAMIFICATION_CONFIG.achievements).flat().find(a => a.id === id)
      ).filter(Boolean),
      achievementProgress: this.getAchievementProgress()
    };
  }
  
  getAchievementProgress() {
    return Object.values(GAMIFICATION_CONFIG.achievements).flat().map(achievement => ({
      ...achievement,
      unlocked: this.userData.achievements.includes(achievement.id),
      progress: this.getAchievementProgressPercent(achievement),
      unlockedAt: this.userData.achievementsUnlockedAt[achievement.id] || null
    }));
  }
  
  getAchievementProgressPercent(achievement) {
    switch (achievement.id) {
      case 'insight-explorer':
        return Math.min((this.userData.stats.insightsViewed / 25) * 100, 100);
      case 'ai-partner':
        return Math.min((this.userData.stats.insightsViewed / 100) * 100, 100);
      case 'insight-master':
        return Math.min((this.userData.stats.actionsCompleted / 50) * 100, 100);
      case 'daily-warrior':
        return Math.min((this.userData.stats.currentStreak / 7) * 100, 100);
      case 'data-enthusiast':
        return Math.min((this.userData.stats.transactionsAdded / 100) * 100, 100);
      default:
        return 0;
    }
  }
  
  // 🧹 CLEANUP
  destroy() {
    this.flushAnalytics();
    this.eventListeners.clear();
  }

  /**
   * Servicio para gestionar gamificación relacionada con insights
   */
  async recordInsightViewed(insightId, insightTitle) {
    try {
      const result = await this.api.recordAction(
        'view_insight',
        'insight', 
        insightId,
        `Viewed insight: ${insightTitle}`
      );
      
      console.log('✅ Insight viewed recorded:', result);
      return result;
    } catch (error) {
      console.warn('⚠️ Error recording insight view:', error);
      // No lanzar error para no interrumpir la experiencia del usuario
      return null;
    }
  }

  /**
   * Registra que el usuario entendió/revisó un insight
   */
  async recordInsightUnderstood(insightId, insightTitle) {
    try {
      const result = await this.api.recordAction(
        'understand_insight',
        'insight',
        insightId,
        `Understood insight: ${insightTitle}`
      );
      
      console.log('✅ Insight understood recorded:', result);
      return result;
    } catch (error) {
      console.warn('⚠️ Error recording insight understanding:', error);
      return null;
    }
  }

  /**
   * Registra que el usuario completó una acción sugerida
   */
  async recordActionCompleted(actionType, description) {
    try {
      const result = await this.api.recordAction(
        'complete_action',
        'suggestion',
        `action_${Date.now()}`,
        description
      );
      
      console.log('✅ Action completed recorded:', result);
      return result;
    } catch (error) {
      console.warn('⚠️ Error recording action completion:', error);
      return null;
    }
  }

  /**
   * Registra que el usuario usó el análisis "¿Puedo comprarlo?"
   */
  async recordPurchaseAnalysisUsed(itemName, amount) {
    try {
      const result = await this.api.recordAction(
        'use_suggestion',
        'suggestion',
        `purchase_analysis_${Date.now()}`,
        `Used purchase analysis for: ${itemName} ($${amount})`
      );
      
      console.log('✅ Purchase analysis recorded:', result);
      return result;
    } catch (error) {
      console.warn('⚠️ Error recording purchase analysis:', error);
      return null;
    }
  }

  /**
   * Obtiene el perfil de gamificación del usuario
   */
  async getUserProfile() {
    try {
      return await this.api.getUserProfile();
    } catch (error) {
      console.warn('⚠️ Error getting user profile:', error);
      return null;
    }
  }

  /**
   * Obtiene las estadísticas de gamificación
   */
  async getUserStats() {
    try {
      return await this.api.getUserStats();
    } catch (error) {
      console.warn('⚠️ Error getting user stats:', error);
      return null;
    }
  }

  /**
   * Obtiene los achievements del usuario
   */
  async getUserAchievements() {
    try {
      return await this.api.getUserAchievements();
    } catch (error) {
      console.warn('⚠️ Error getting achievements:', error);
      return [];
    }
  }
}

// 🌟 SINGLETON PATTERN PARA ESCALABILIDAD
let gamificationInstance = null;

export const getGamificationService = () => {
  if (!gamificationInstance) {
    gamificationInstance = new GamificationService();
  }
  return gamificationInstance;
};

// 🎯 EXPORTS PRINCIPALES
export { GAMIFICATION_CONFIG };
export default GamificationService; 