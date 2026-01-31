import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { getGamificationAPI } from '../services/gamificationAPI';
import { useGamificationNotifications } from '../components/GamificationNotification';
import { useAuth } from './AuthContext';
/**
 * 🎮 GAMIFICATION CONTEXT
 * 
 * Contexto global para manejar:
 * - Estado de gamificación del usuario
 * - Notificaciones de XP y logros
 * - Registro de acciones
 * - Features desbloqueables por nivel
 * - Cache y sincronización
 */

// 🔒 FEATURE GATES - Definición de features desbloqueables
const FEATURE_GATES = {
  SAVINGS_GOALS: {
    name: 'Metas de Ahorro',
    description: 'Crea y gestiona objetivos de ahorro personalizados',
    requiredLevel: 3,
    icon: '🎯',
    benefits: ['Objetivos personalizados', 'Seguimiento de progreso', 'Auto-ahorro'],
    xpThreshold: 200  // ACTUALIZADO: Más fácil de alcanzar
  },
  BUDGETS: {
    name: 'Presupuestos',
    description: 'Controla tus gastos con límites inteligentes por categoría',
    requiredLevel: 5,
    icon: '📊',
    benefits: ['Límites por categoría', 'Alertas automáticas', 'Control de gastos'],
    xpThreshold: 700  // ACTUALIZADO: Más fácil de alcanzar
  },
  AI_INSIGHTS: {
    name: 'IA Financiera',
    description: 'Análisis inteligente con IA para decisiones financieras',
    requiredLevel: 7,
    icon: '🧠',
    benefits: ['Análisis de compras', 'Score crediticio', 'Insights personalizados'],
    xpThreshold: 1800  // ACTUALIZADO: Más fácil de alcanzar
  }
};

// 🏆 NIVELES DEL SISTEMA - REBALANCEADOS PARA PROGRESIÓN SIN DEPENDENCIAS
const LEVEL_SYSTEM = {
  1: { name: 'Financial Newbie', minXP: 0, color: '#9CA3AF' },
  2: { name: 'Money Tracker', minXP: 75, color: '#10B981' },      // REDUCIDO: 100 → 75
  3: { name: 'Smart Saver', minXP: 200, color: '#3B82F6' },      // REDUCIDO: 300 → 200 🔓 METAS
  4: { name: 'Budget Master', minXP: 400, color: '#8B5CF6' },    // REDUCIDO: 600 → 400
  5: { name: 'Financial Planner', minXP: 700, color: '#F59E0B' }, // REDUCIDO: 1000 → 700 🔓 PRESUPUESTOS
  6: { name: 'Investment Seeker', minXP: 1200, color: '#EF4444' }, // REDUCIDO: 1500 → 1200
  7: { name: 'Wealth Builder', minXP: 1800, color: '#EC4899' },   // REDUCIDO: 2200 → 1800 🔓 IA
  8: { name: 'Financial Strategist', minXP: 2600, color: '#06B6D4' },
  9: { name: 'Money Mentor', minXP: 3600, color: '#84CC16' },
  10: { name: 'Financial Magnate', minXP: 5000, color: '#F97316' }
};

const GamificationContext = createContext();

export const useGamification = () => {
  const context = useContext(GamificationContext);
  if (!context) {
    throw new Error('useGamification must be used within a GamificationProvider');
  }
  return context;
};

export const GamificationProvider = ({ children }) => {
  // Estado del usuario
  const [userProfile, setUserProfile] = useState(null);
  const [achievements, setAchievements] = useState([]);
  const [stats, setStats] = useState(null);
  const [features, setFeatures] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Cache para evitar llamadas innecesarias
  const [lastUpdate, setLastUpdate] = useState(null);
  const [pendingActions, setPendingActions] = useState([]);
  const [featureAccessCache, setFeatureAccessCache] = useState(new Map());
  
  // Forzar re-render de componentes dependientes
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  // Contexto de autenticación
  const auth = useAuth();

  // Notificaciones
  const {
    showXPGained,
    showLevelUp,
    showAchievementUnlocked,
    GamificationNotification
  } = useGamificationNotifications();

  const api = getGamificationAPI();

  // Cargar datos iniciales solo cuando el usuario esté autenticado
  useEffect(() => {
    // Solo cargar si el usuario está autenticado y la auth está inicializada
    if (auth.authState === 'authenticated' && auth.isInitialized) {
      loadGamificationData();
    } else if (auth.authState === 'unauthenticated') {
      // Limpiar datos si el usuario no está autenticado
      setUserProfile(null);
      setAchievements([]);
      setStats(null);
      setFeatures([]);
      setError(null);
      setLoading(false);
    }
  }, [auth.authState, auth.isInitialized]);

  const loadGamificationData = async () => {
    try {
      setLoading(true);
      
      const [profileData, achievementsData, statsData, featuresData] = await Promise.all([
        api.getUserProfile().catch(err => {
          console.warn('Error loading user profile:', err);
          return null;
        }),
        api.getUserAchievements().catch(err => {
          console.warn('Error loading achievements:', err);
          return [];
        }),
        api.getUserStats().catch(err => {
          console.warn('Error loading stats:', err);
          return null;
        }),
        api.getUserFeatures().catch(err => {
          console.warn('Error loading features:', err);
          return { unlocked_features: [], locked_features: [] };
        })
      ]);

      setUserProfile(profileData);
      setAchievements(achievementsData || []);
      setStats(statsData);
      setFeatures(featuresData || { unlocked_features: [], locked_features: [] });
      setLastUpdate(Date.now());
      setError(null);
    } catch (err) {
      console.error('Error loading gamification data:', err);
      setError(err.message);
      // Establecer valores por defecto en caso de error
      setUserProfile(null);
      setAchievements([]);
      setStats(null);
      setFeatures({ unlocked_features: [], locked_features: [] });
    } finally {
      setLoading(false);
    }
  };

  // Registrar acción y manejar resultado
  const recordAction = useCallback(async (actionType, entityType, entityId, description) => {
    // Evitar acciones duplicadas
    const actionKey = `${actionType}-${entityType}-${entityId}`;
    
    try {
      if (pendingActions.includes(actionKey)) {
        return null;
      }

      setPendingActions(prev => [...prev, actionKey]);

      const result = await api.recordAction(
        actionType,
        entityType,
        entityId,
        description || `User ${actionType} ${entityType}`
      );

      // Process action result

              // Actualizar datos locales inmediatamente solo si hay cambios reales
        setUserProfile(prev => {
          const prevXP = prev?.total_xp || 0;
          const prevLevel = prev?.current_level || 0;
          const newXP = result.total_xp || prevXP;
          // 🔧 FIX: Usar current_level del resultado, no new_level
          const newLevel = result.current_level || result.new_level || prevLevel;

          // Solo actualizar si realmente cambió algo
          if (newXP !== prevXP || newLevel !== prevLevel) {

            return {
              ...prev,
              total_xp: newXP,
              current_level: newLevel,
              last_updated: Date.now()
            };
          }
          

          return prev; // No cambios, no re-render
        });

        setStats(prev => {
          const newTotalXP = result.total_xp || prev?.total_xp || 0;
          // 🔧 FIX: Usar current_level del resultado, no new_level
          const newLevel = result.current_level || result.new_level || prev?.current_level || 0;
          
          // Solo actualizar si hay cambios
          if (newTotalXP !== (prev?.total_xp || 0) || newLevel !== (prev?.current_level || 0)) {
            return {
              ...prev,
              total_xp: newTotalXP,
              current_level: newLevel,
              last_updated: Date.now()
            };
          }
          
          return prev; // No cambios
        });

      if (result.xp_earned > 0 && actionType !== 'view_insight') {

        // Mostrar notificación de XP ganado
        showXPGained(result.xp_earned, `¡Has ganado ${result.xp_earned} XP!`);
      }

      // Mostrar notificación de subida de nivel
      if (result.level_up) {
        const actualNewLevel = result.current_level || result.new_level;

        const levelInfo = getLevelInfo(actualNewLevel);
        showLevelUp(actualNewLevel, levelInfo.name);
      }

      // Mostrar notificaciones de nuevos logros
      if (result.new_achievements && result.new_achievements.length > 0) {

        result.new_achievements.forEach(achievement => {
          showAchievementUnlocked(achievement.name, achievement.description);
        });

        // Actualizar logros locales
        setAchievements(prev => 
          prev.map(a => {
            const updated = result.new_achievements.find(na => na.id === a.id);
            return updated ? { ...a, ...updated } : a;
          })
        );
      }

      // Actualizar userProfile local con los nuevos datos
      if (userProfile && result.total_xp !== undefined) {
        setUserProfile(prev => ({
          ...prev,
          total_xp: result.total_xp,
          current_level: result.new_level || result.current_level || prev.current_level,
          achievements_count: prev.achievements_count + (result.new_achievements?.length || 0)
        }));

      }

      // Limpiar acción pendiente
      setPendingActions(prev => prev.filter(p => p !== actionKey));

      // Forzar re-render del contexto actualizando el timestamp
      setLastUpdate(Date.now());
      setRefreshTrigger(prev => prev + 1); // Forzar re-render de componentes dependientes
            
      return result;
    } catch (err) {
      console.error('Error recording gamification action:', err);
      setPendingActions(prev => prev.filter(p => p !== actionKey));
      return null;
    }
  }, [api, pendingActions, showXPGained, showLevelUp, showAchievementUnlocked]);

  // 📊 Métodos de conveniencia para acciones básicas
  const recordViewDashboard = useCallback(() => {
    return recordAction('view_dashboard', 'dashboard', 'main-dashboard', 'User viewed dashboard');
  }, [recordAction]);

  const recordViewExpenses = useCallback(() => {
    return recordAction('view_expenses', 'expense', 'expense-list', 'User viewed expenses');
  }, [recordAction]);

  const recordViewIncomes = useCallback(() => {
    return recordAction('view_incomes', 'income', 'income-list', 'User viewed incomes');
  }, [recordAction]);

  const recordViewCategories = useCallback(() => {
    return recordAction('view_categories', 'category', 'category-list', 'User viewed categories');
  }, [recordAction]);

  const recordViewAnalytics = useCallback((analyticsType = 'general') => {
    return recordAction('view_analytics', 'analytics', analyticsType, `User viewed ${analyticsType} analytics`);
  }, [recordAction]);

  // 💰 Métodos de conveniencia para transacciones
  const recordCreateExpense = useCallback((expenseId, description = 'User created expense') => {
    return recordAction('create_expense', 'expense', expenseId, description);
  }, [recordAction]);

  const recordCreateIncome = useCallback((incomeId, description = 'User created income') => {
    return recordAction('create_income', 'income', incomeId, description);
  }, [recordAction]);

  const recordUpdateExpense = useCallback((expenseId, description = 'User updated expense') => {
    return recordAction('update_expense', 'expense', expenseId, description);
  }, [recordAction]);

  const recordDeleteExpense = useCallback((expenseId, description = 'User deleted expense') => {
    return recordAction('delete_expense', 'expense', expenseId, description);
  }, [recordAction]);

  // 🏷️ Métodos de conveniencia para organización
  const recordCreateCategory = useCallback((categoryId, description = 'User created category') => {
    return recordAction('create_category', 'category', categoryId, description);
  }, [recordAction]);

  const recordAssignCategory = useCallback((transactionId, categoryId, description = 'User assigned category') => {
    return recordAction('assign_category', 'transaction', transactionId, `${description} - Category: ${categoryId}`);
  }, [recordAction]);

  // 🎯 Métodos de conveniencia para engagement
  const recordDailyLogin = useCallback(() => {
    return recordAction('daily_login', 'user', 'daily-login', 'User daily login');
  }, [recordAction]);

  // 🤖 Métodos de conveniencia para IA (Legacy - mantener compatibilidad)
  const recordInsightViewed = useCallback((insightId, insightTitle) => {
    return recordAction('view_insight', 'insight', insightId, `Viewed insight: ${insightTitle}`);
  }, [recordAction]);

  const recordInsightUnderstood = useCallback((insightId, insightTitle) => {
    return recordAction('understand_insight', 'insight', insightId, `Understood insight: ${insightTitle}`);
  }, [recordAction]);

  const recordActionCompleted = useCallback((actionId, actionDescription) => {
    return recordAction('complete_action', 'suggestion', actionId, `Completed action: ${actionDescription}`);
  }, [recordAction]);

  const recordPatternViewed = useCallback((patternId, patternType) => {
    return recordAction('view_pattern', 'pattern', patternId, `Viewed pattern: ${patternType}`);
  }, [recordAction]);

  const recordSuggestionUsed = useCallback((suggestionId, suggestionTitle) => {
    return recordAction('use_suggestion', 'suggestion', suggestionId, `Applied suggestion: ${suggestionTitle}`);
  }, [recordAction]);

  // Información de niveles
  const getLevelInfo = (level) => {
    const levels = [
      { level: 0, name: "Financial Newbie", color: "#94A3B8", emoji: "🌱" },
      { level: 1, name: "Money Aware", color: "#60A5FA", emoji: "👀" },
      { level: 2, name: "Budget Tracker", color: "#34D399", emoji: "📊" },
      { level: 3, name: "Savings Starter", color: "#FBBF24", emoji: "💰" },
      { level: 4, name: "Financial Explorer", color: "#F472B6", emoji: "🧭" },
      { level: 5, name: "Money Manager", color: "#A78BFA", emoji: "💼" },
      { level: 6, name: "Investment Learner", color: "#FB7185", emoji: "📈" },
      { level: 7, name: "Financial Guru", color: "#10B981", emoji: "🧠" },
      { level: 8, name: "Money Master", color: "#8B5CF6", emoji: "👑" },
      { level: 9, name: "Financial Magnate", color: "#EF4444", emoji: "💎" }
    ];
    
    return levels[level] || levels[0];
  };

  // Refresh de datos (útil después de acciones importantes)
  const refreshData = useCallback(async () => {
    // Solo refrescar si han pasado al menos 30 segundos desde la última actualización
    if (lastUpdate && Date.now() - lastUpdate < 30000) {
      return;
    }
    
    await loadGamificationData();
  }, [lastUpdate]);

  // 🔒 Funciones de Feature Gates
  const isFeatureUnlocked = useCallback((featureKey) => {
    if (!userProfile) return false;
    
    // Usar datos del backend si están disponibles
    if (features && features.unlocked_features && Array.isArray(features.unlocked_features)) {
      if (features.unlocked_features.includes(featureKey)) return true;
      // Si aparece como bloqueada pero el backend informa trial activo, considerarla desbloqueada
      if (features.locked_features && Array.isArray(features.locked_features)) {
        const lf = features.locked_features.find(f => f.feature_key === featureKey);
        if (lf && (lf.trial_active === true)) return true;
      }
      return false;
    }
    
    // Fallback usando nivel local
    const feature = FEATURE_GATES[featureKey];
    if (!feature) return true; // Si la feature no existe, permitir acceso
    
    const userLevel = userProfile.current_level || 0;
    return userLevel >= feature.requiredLevel;
  }, [userProfile, features]);

  const getFeatureAccess = useCallback(async (featureKey) => {
    // Si tenemos datos del backend, usar esos
    if (features && features.locked_features && Array.isArray(features.locked_features)) {
      const lockedFeature = features.locked_features.find(f => f.feature_key === featureKey);
      if (lockedFeature) {
        return {
          unlocked: false,
          requiredLevel: lockedFeature.required_level,
          userLevel: lockedFeature.current_level,
          xpNeeded: lockedFeature.xp_needed,
          featureName: lockedFeature.feature_name,
          featureIcon: lockedFeature.feature_icon,
          description: lockedFeature.description,
          trialActive: Boolean(lockedFeature.trial_active),
          trialEndsAt: lockedFeature.trial_ends_at || null
        };
      }
    }
    
    // Si la feature está desbloqueada o como fallback
    const feature = FEATURE_GATES[featureKey];
    const userLevel = userProfile?.current_level || 0;
    const unlocked = userLevel >= (feature?.requiredLevel || 0);
    
    return {
      unlocked,
      requiredLevel: feature?.requiredLevel || 0,
      userLevel,
      xpNeeded: unlocked ? 0 : LEVEL_SYSTEM[feature?.requiredLevel]?.minXP - (userProfile?.total_xp || 0),
      featureName: feature?.name || featureKey,
      featureIcon: feature?.icon || '🔒',
      description: feature?.description || 'Feature description',
      trialActive: false,
      trialEndsAt: null
    };
  }, [userProfile, features]);

  // 🔒 Verificar acceso a feature específica (con llamada al backend si es necesario)
  const checkFeatureAccess = useCallback(async (featureKey) => {
    // Verificar cache primero (cache por 5 minutos)
    const cached = featureAccessCache.get(featureKey);
    const now = Date.now();
    if (cached && (now - cached.timestamp) < 5 * 60 * 1000) {
      return cached.data;
    }

    try {
      const access = await api.checkFeatureAccess(featureKey);
      
      // Guardar en cache
      setFeatureAccessCache(prev => new Map(prev.set(featureKey, {
        data: access,
        timestamp: now
      })));
      
      return access;
    } catch (error) {
      // Solo mostrar error si no es 404 (endpoint no implementado)
      if (error.response?.status !== 404) {
        console.error(`Error checking feature access for ${featureKey}:`, error);
      }
      
      // Fallback silencioso a verificación local
      const localAccess = getFeatureAccess(featureKey);
      
      // También cachear el fallback local
      setFeatureAccessCache(prev => new Map(prev.set(featureKey, {
        data: localAccess,
        timestamp: now
      })));
      
      return localAccess;
    }
  }, [api, getFeatureAccess, featureAccessCache]);

  const value = {
    // Estado
    userProfile,
    achievements,
    stats,
    features,
    loading,
    error,
    refreshTrigger, // Para forzar re-renders

    // Acciones principales
    recordAction,
    
    // 📊 Acciones básicas
    recordViewDashboard,
    recordViewExpenses,
    recordViewIncomes,
    recordViewCategories,
    recordViewAnalytics,
    
    // 💰 Acciones de transacciones
    recordCreateExpense,
    recordCreateIncome,
    recordUpdateExpense,
    recordDeleteExpense,
    
    // 🏷️ Acciones de organización
    recordCreateCategory,
    recordAssignCategory,
    
    // 🎯 Acciones de engagement
    recordDailyLogin,
    
    // 🤖 Acciones de IA (Legacy)
    recordInsightViewed,
    recordInsightUnderstood,
    recordActionCompleted,
    recordPatternViewed,
    recordSuggestionUsed,

    // Feature Gates
    isFeatureUnlocked,
    getFeatureAccess,
    checkFeatureAccess, // 🔒 Verificación con backend
    FEATURE_GATES,
    LEVEL_SYSTEM,

    // 🏆 Daily Challenges
    getDailyChallenges: api.getDailyChallenges.bind(api),
    getWeeklyChallenges: api.getWeeklyChallenges.bind(api),
    processChallengeProgress: api.processChallengeProgress.bind(api),

    // Utilidades
    getLevelInfo,
    refreshData,
    loadGamificationData,

    // Componente de notificaciones
    GamificationNotification,

    // Estado de cache
    lastUpdate,
    pendingActions
  };

  return (
    <GamificationContext.Provider value={value}>
      {children}
      <GamificationNotification />
    </GamificationContext.Provider>
  );
};

export default GamificationContext; 