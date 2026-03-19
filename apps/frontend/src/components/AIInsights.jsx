import React, { useState, useEffect, useCallback, useRef } from 'react';
import { FaBrain, FaSpinner, FaRedo, FaLightbulb, FaShoppingCart, FaCheck, FaChevronRight, FaCalculator, FaExclamationTriangle, FaCheckCircle, FaChevronDown, FaChevronUp, FaBullseye } from 'react-icons/fa';
import { aiAPI, savingsGoalsAPI, budgetsAPI, analyticsAPI, dashboardAPI, formatCurrency } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import { usePeriod } from '../contexts/PeriodContext';
import { useGamification } from '../contexts/GamificationContext';
import { getGamificationAPI } from '../services/gamificationAPI';
import MonthlyCoachingTab from '../pages/insights/tabs/MonthlyCoachingTab';
import EducationTab from '../pages/insights/tabs/EducationTab';

const AIInsights = () => {
  const { user, isAuthenticated } = useAuth();
  usePeriod(); // keep context subscription even if not directly used here
  const { recordInsightViewed, recordInsightUnderstood, recordSuggestionUsed } = useGamification();
  const [insights, setInsights] = useState([]);
  const [purchaseAnalysis, setPurchaseAnalysis] = useState(null);
  const [loading, setLoading] = useState(false);
  const [purchaseLoading, setPurchaseLoading] = useState(false);
  const [error, setError] = useState(null);
  const [purchaseError, setPurchaseError] = useState(null);
  const [healthScore, setHealthScore] = useState(0);
  const [healthDetails, setHealthDetails] = useState(null);
  const [healthScoreLoading, setHealthScoreLoading] = useState(false);
  const [lastEvaluationDate, setLastEvaluationDate] = useState(null);
  const [dashboardData, setDashboardData] = useState(null);
  const [savingsGoals, setSavingsGoals] = useState([]);
  const [behaviorProfile, setBehaviorProfile] = useState(null);
  const behaviorProfileRef = useRef(null); // ref to avoid circular dep in loadHealthScore
  const [appliedInsights, setAppliedInsights] = useState(new Set());
  
  const [activeTab, setActiveTab] = useState('monthly');
  
  // Estados para el análisis de compra - ahora se inicializan dinámicamente
  const [purchaseForm, setPurchaseForm] = useState({
    itemName: '',
    amount: '',
    description: '',
    paymentTypes: [], // Cambiado de paymentType a paymentTypes array
    isNecessary: false,
    currentBalance: 0,
    monthlyIncome: 0,
    monthlyExpenses: 0,
    savingsGoal: 0
  });

  // Estados para gamificación - UX optimizada (seguridad real en backend)
  const [viewedInsights, setViewedInsights] = useState(new Set());
  const [understoodInsights, setUnderstoodInsights] = useState(() => {
    // UX: Cargar estado desde sessionStorage para mejor experiencia
    // 🛡️ SEGURIDAD REAL: El backend verifica duplicados en base de datos
    try {
      const userId = user?.id || 'guest';
      const stored = sessionStorage.getItem(`understood_insights_${userId}`);
      return stored ? new Set(JSON.parse(stored)) : new Set();
    } catch (error) {
      console.warn('⚠️ Error loading understood insights from sessionStorage:', error);
      return new Set();
    }
  });

  const paymentTypes = [
    { value: 'contado', label: 'Pago de contado' },
    { value: 'cuotas', label: 'Plan de pagos/cuotas' },
    { value: 'ahorro', label: 'Necesito ahorrar para esto' }
  ];

  // Función para manejar cambios en tipos de pago múltiples
  const handlePaymentTypeChange = (paymentType, isChecked) => {
    setPurchaseForm(prev => {
      const currentTypes = prev.paymentTypes || [];
      if (isChecked) {
        // Agregar tipo si no existe
        if (!currentTypes.includes(paymentType)) {
          return { ...prev, paymentTypes: [...currentTypes, paymentType] };
        }
      } else {
        // Remover tipo si existe
        return { ...prev, paymentTypes: currentTypes.filter(type => type !== paymentType) };
      }
      return prev;
    });
  };

  // Cargar datos del dashboard directamente desde la API (campos canónicos del backend)
  const loadDashboardData = useCallback(async () => {
    try {
      const response = await dashboardAPI.overview();
      const d = response.data;
      setDashboardData(d);
      setPurchaseForm(prev => ({
        ...prev,
        currentBalance: d?.current_month_balance || 0,
        monthlyIncome: d?.current_month_incomes || 0,
        monthlyExpenses: d?.current_month_expenses || 0,
      }));
    } catch (error) {
      console.error('Error loading dashboard data:', error);
    }
  }, []);

  // Cargar metas de ahorro
  const loadSavingsGoals = useCallback(async () => {
    if (!isAuthenticated) return;
    
    try {
      const response = await savingsGoalsAPI.list({ status: 'active' });
      const goals = response.data?.data?.goals || [];
      setSavingsGoals(goals);
    } catch (error) {
      console.error('Error loading savings goals:', error);
      setSavingsGoals([]);
    }
  }, [isAuthenticated]);

  // Cache removido - el backend maneja su propio cache de 20 horas

  // Función para cargar el health score, pasando parámetros conductuales opcionales
  const loadHealthScore = useCallback(async (profile = null) => {
    if (!isAuthenticated) return;

    setHealthScoreLoading(true);
    try {
      // Use explicit profile arg, then ref (avoids circular dep — ref never causes re-render)
      const bp = profile || behaviorProfileRef.current;
      const params = bp ? {
        streak: bp.current_streak,
        days_active: bp.days_active,
        budgets_created: bp.budgets_created,
        budget_compliance: bp.budget_compliance_events,
        savings_goals: bp.savings_goals_created,
        savings_deposits: bp.savings_deposits,
        recurring_setups: bp.recurring_setups,
        ai_applied: bp.ai_recommendations_applied,
      } : {};

      const response = await aiAPI.getHealthScore(params);
      // Backend returns score on 0-1000 scale — round to avoid floating-point display artefacts
      setHealthScore(Math.round(response.score || response.health_score || 0));
      setHealthDetails(response);
    } catch (err) {
      console.error('Error loading health score:', err.message);
      setHealthScore(0);
    } finally {
      setHealthScoreLoading(false);
    }
  }, [isAuthenticated]); // behaviorProfile removed — using ref instead to break the dep cycle

  const loadAIInsights = useCallback(async () => {
    if (!isAuthenticated) {
      setError('Debes iniciar sesión para ver el análisis inteligente');
      return;
    }

    setLoading(true);
    setError(null);
    try {
      console.log('🔍 Cargando análisis inteligente para usuario:', user?.email);

      // Fetch all financial context data in parallel (including behavioral profile)
      const gamificationAPI = getGamificationAPI();
      const [categoriesResult, incomesResult, budgetResult, goalsResult, behaviorResult] = await Promise.allSettled([
        analyticsAPI.categories(),
        analyticsAPI.incomes(),
        budgetsAPI.getDashboard(),
        savingsGoalsAPI.list({ status: 'active' }),
        gamificationAPI.getBehaviorProfile(),
      ]);

      const financialData = {};

      // Expense categories breakdown (aggregated by category name)
      if (categoriesResult.status === 'fulfilled') {
        const cats = categoriesResult.value.data?.data || [];
        financialData.expenses_by_category = Object.fromEntries(
          cats.map(c => [c.category_name || 'Sin categoría', c.amount || 0])
        );
        financialData.total_expenses = cats.reduce((s, c) => s + (c.amount || 0), 0);
      }

      // Total income from analytics summary
      if (incomesResult.status === 'fulfilled') {
        const incomesData = incomesResult.value.data;
        financialData.total_income = incomesData?.total_amount || incomesData?.data?.total_amount || 0;
        // Income stability: based on number of income sources (more sources = more stable)
        const incomeCount = incomesData?.count || incomesData?.data?.count || 1;
        financialData.income_stability = incomeCount >= 3 ? 0.9 : incomeCount === 2 ? 0.65 : 0.35;
      }

      // Savings rate
      const ti = financialData.total_income || 0;
      const te = financialData.total_expenses || 0;
      financialData.savings_rate = ti > 0 ? (ti - te) / ti : 0;

      // Budget compliance summary
      if (budgetResult.status === 'fulfilled') {
        const summary = budgetResult.value.data?.summary || {};
        if (summary.total_budgets > 0) {
          financialData.budgets_summary = {
            total_budgets: summary.total_budgets || 0,
            total_allocated: summary.total_allocated || 0,
            total_spent: summary.total_spent || 0,
            on_track_count: summary.on_track_count || 0,
            warning_count: summary.warning_count || 0,
            exceeded_count: summary.exceeded_count || 0,
            average_usage: summary.average_usage || 0,
          };
        }
      }

      // Active savings goals
      if (goalsResult.status === 'fulfilled') {
        const goals = goalsResult.value.data?.data?.goals || [];
        if (goals.length > 0) {
          financialData.savings_goals = goals.map(g => ({
            name: g.name,
            target_amount: g.target_amount,
            current_amount: g.current_amount,
            progress: g.progress || 0,
            target_date: g.target_date,
          }));
        }
      }

      // Behavioral profile from gamification (best-effort)
      let bp = null;
      if (behaviorResult.status === 'fulfilled' && behaviorResult.value) {
        bp = behaviorResult.value;
        setBehaviorProfile(bp);
        behaviorProfileRef.current = bp; // sync ref so loadHealthScore can use it immediately
        financialData.behavior_profile = {
          current_level: bp.current_level,
          level_name: bp.level_name,
          current_streak: bp.current_streak,
          days_active: bp.days_active,
          budgets_created: bp.budgets_created,
          budget_compliance_events: bp.budget_compliance_events,
          savings_goals_created: bp.savings_goals_created,
          savings_deposits: bp.savings_deposits,
          savings_goals_achieved: bp.savings_goals_achieved,
          recurring_setups: bp.recurring_setups,
          ai_recommendations_applied: bp.ai_recommendations_applied,
          consistency_score: bp.consistency_score,
          discipline_score: bp.discipline_score,
          engagement_score: bp.engagement_score,
        };
      }

      // Dynamic period: current month name + year
      const now = new Date();
      const monthName = now.toLocaleString('es-AR', { month: 'long', year: 'numeric' });
      financialData.period = monthName.charAt(0).toUpperCase() + monthName.slice(1);
      financialData.financial_score = 0;

      const response = await aiAPI.getInsights(financialData);
      const newInsights = response.insights || [];
      setInsights(newInsights);

      // Reload health score now that we have the behavioral profile
      loadHealthScore(bp);

      // Usar el timestamp del backend (generated_at)
      const backendTimestamp = response.generated_at ? new Date(response.generated_at) : new Date();
      setLastEvaluationDate(backendTimestamp);
      console.log('💾 Análisis cargado desde backend - Timestamp:', backendTimestamp.toISOString());
    } catch (err) {
      console.error('Error loading AI insights:', err.message);
      setError('Error conectando con el análisis de IA. Usando datos de ejemplo.');
      // Usar datos de ejemplo
      const fallbackInsights = [
        {
          title: "Excelente capacidad de ahorro",
          description: "Estás ahorrando 32% de tus ingresos, superando el promedio nacional. Considera explorar opciones de inversión para hacer crecer tu dinero.",
          impact: "high",
          score: 920,
          action_type: "invest",
          category: "ahorro"
        },
        {
          title: "Mayor gasto: Alimentación",
          description: "La Alimentación representa 42.4% de tus gastos ($137,000). Revisa si hay oportunidades de optimización en esta categoría.",
          impact: "medium",
          score: 400,
          action_type: "optimize",
          category: "Alimentación"
        },
        {
          title: "Ingresos variables",
          description: "Tus ingresos muestran variabilidad. Considera diversificar fuentes de ingresos o crear un fondo de emergencia más robusto.",
          impact: "medium",
          score: 600,
          action_type: "save",
          category: "ingresos"
        },
        {
          title: "Oportunidad de inversión",
          description: "Tienes $50,000 disponibles que podrías invertir en instrumentos de bajo riesgo para generar ingresos pasivos.",
          impact: "high",
          score: 850,
          action_type: "invest",
          category: "inversión"
        },
        {
          title: "Control de gastos hormiga",
          description: "Los pequeños gastos diarios suman $15,000 mensuales. Considera usar una app de control de gastos.",
          impact: "low",
          score: 300,
          action_type: "optimize",
          category: "gastos"
        }
      ];
      setInsights(fallbackInsights);

      // Usar timestamp actual para datos de fallback
      const fallbackTimestamp = new Date();
      setLastEvaluationDate(fallbackTimestamp);
      console.log('💾 Usando datos de fallback - Timestamp:', fallbackTimestamp.toISOString());
    } finally {
      setLoading(false);
    }
  }, [isAuthenticated, user?.email, loadHealthScore]);

  // Función simplificada - siempre llama al backend (que tiene su propio cache de 20h)
  const loadAIInsightsSimple = useCallback(async () => {
    if (!isAuthenticated) {
      setError('Debes iniciar sesión para ver el análisis inteligente');
      return;
    }

    // Siempre llamar al backend - él maneja su propio cache de 20 horas
    await loadAIInsights();
  }, [isAuthenticated, loadAIInsights]);

  useEffect(() => {
    if (!isAuthenticated) {
      console.warn('⚠️ Usuario no autenticado, no se cargarán insights de IA');
      setError('Debes iniciar sesión para ver el análisis inteligente');
      return;
    }
    // Cache del frontend deshabilitado - confiamos en el cache del backend (20 horas)
    // if (process.env.NODE_ENV === 'development') {
    //   localStorage.removeItem('ai_insights_cache');
    //   localStorage.removeItem('health_score_cache');
    //   console.log('🧹 Cache limpiado para desarrollo');
    // }
    loadAIInsightsSimple(); // internally calls loadHealthScore(bp) after getting behavior profile
    loadDashboardData();
    loadSavingsGoals();
  }, [isAuthenticated, loadAIInsightsSimple, loadDashboardData, loadSavingsGoals]);

  // Función para filtrar metas de ahorro relevantes
  const getRelevantSavingsGoals = (itemName, description) => {
    if (!savingsGoals || savingsGoals.length === 0) return [];

    const itemNameLower = itemName.toLowerCase();
    const descriptionLower = description.toLowerCase();

    return savingsGoals.filter(goal => {
      if (goal.status !== 'active') return false;

      const goalNameLower = goal.name.toLowerCase();
      const goalCategoryLower = goal.category.toLowerCase();

      // 1. Coincidencia directa en nombre o categoría
      if (itemNameLower.includes(goalCategoryLower) ||
          goalNameLower.includes(itemNameLower) ||
          goalCategoryLower.includes(itemNameLower)) {
        return true;
      }

      // 2. Coincidencias específicas por categoría
      switch (goal.category) {
        case 'car':
          if (itemNameLower.includes('auto') ||
              itemNameLower.includes('carro') ||
              itemNameLower.includes('vehículo') ||
              itemNameLower.includes('vehiculo') ||
              descriptionLower.includes('auto') ||
              descriptionLower.includes('carro')) {
            return true;
          }
          break;
        case 'house':
          if (itemNameLower.includes('casa') ||
              itemNameLower.includes('vivienda') ||
              itemNameLower.includes('inmueble') ||
              itemNameLower.includes('propiedad')) {
            return true;
          }
          break;
        case 'vacation':
          if (itemNameLower.includes('viaje') ||
              itemNameLower.includes('vacaciones') ||
              itemNameLower.includes('turismo')) {
            return true;
          }
          break;
        case 'education':
          if (itemNameLower.includes('curso') ||
              itemNameLower.includes('educación') ||
              itemNameLower.includes('educacion') ||
              itemNameLower.includes('estudio')) {
            return true;
          }
          break;
      }

      // 3. Análisis de descripción para palabras clave
      if (description && (descriptionLower.includes(goalNameLower) ||
                         descriptionLower.includes(goalCategoryLower))) {
        return true;
      }

      return false;
    }).map(goal => ({
      name: goal.name,
      category: goal.category,
      current_amount: goal.current_amount,
      target_amount: goal.target_amount,
      progress: goal.progress
    }));
  };

  const analyzePurchase = async () => {
    if (!purchaseForm.itemName || !purchaseForm.amount || purchaseForm.paymentTypes.length === 0) {
      setPurchaseError('Por favor completa el nombre del artículo, el monto y selecciona al menos un tipo de pago');
      return;
    }

    setPurchaseLoading(true);
    setPurchaseError(null);
    try {
      // Obtener metas de ahorro relevantes
      const relevantSavingsGoals = getRelevantSavingsGoals(purchaseForm.itemName, purchaseForm.description);
      
      const response = await aiAPI.canIBuy({
        item_name: purchaseForm.itemName,
        amount: parseFloat(purchaseForm.amount),
        description: purchaseForm.description,
        payment_types: purchaseForm.paymentTypes, // Enviar array de tipos de pago
        is_necessary: purchaseForm.isNecessary,
        current_balance: purchaseForm.currentBalance,
        monthly_income: purchaseForm.monthlyIncome,
        monthly_expenses: purchaseForm.monthlyExpenses,
        savings_goal: purchaseForm.savingsGoal,
        savings_goals: relevantSavingsGoals // Agregar metas de ahorro relevantes
      });
      setPurchaseAnalysis(response);
      
      // 🎮 Registrar acción de gamificación
      await recordSuggestionUsed(
        `purchase-analysis-${Date.now()}`,
        `Purchase analysis: ${purchaseForm.itemName}`
      );
      
    } catch (err) {
      console.error('Error analyzing purchase:', err.message);
      
      // Si es un error de IA no configurada, mostrar mensaje específico
      if (err.message.includes('IA no configurada') || err.message.includes('no disponible')) {
        setPurchaseError('❌ Análisis de compra no disponible: IA no configurada. Esta función requiere OpenAI para funcionar correctamente.');
      } else {
        setPurchaseError('❌ Error conectando con la IA. Verifica tu conexión e intenta nuevamente.');
      }
      
      setPurchaseAnalysis(null);
    } finally {
      setPurchaseLoading(false);
    }
  };

  // 🎮 Funciones de gamificación mejoradas
  const handleViewInsight = async (insightId, insightTitle) => {
    // Registrar solo una vez por sesión para cada insightId
    const sessionKey = `gami_view_insight_${insightId}`;
    try {
      if (sessionStorage.getItem(sessionKey)) return;
      sessionStorage.setItem(sessionKey, '1');
    } catch (_) {}

    if (!viewedInsights.has(insightId)) {
      setViewedInsights(prev => new Set([...prev, insightId]));
    }
    await recordInsightViewed(String(insightId), insightTitle);
  };

  const handleUnderstandInsight = async (insightId, insightTitle) => {
    if (!understoodInsights.has(insightId)) {
      const newUnderstoodInsights = new Set([...understoodInsights, insightId]);
      setUnderstoodInsights(newUnderstoodInsights);
      
      // UX: Persistir en sessionStorage para mejor experiencia de usuario
      // 🛡️ SEGURIDAD: El backend tiene la verificación real de duplicados
      try {
        const userId = user?.id || 'guest';
        sessionStorage.setItem(
          `understood_insights_${userId}`,
          JSON.stringify([...newUnderstoodInsights])
        );
        console.log('✅ [UX] Insight marcado como entendido:', insightId);
      } catch (error) {
        console.warn('⚠️ Error saving understood insights to sessionStorage:', error);
      }
      
      // Registrar en gamificación - EL BACKEND VERIFICA DUPLICADOS
      try {
        await recordInsightUnderstood(String(insightId), insightTitle);
        console.log('✅ [Gamification] XP procesado para insight:', insightId);
      } catch (error) {
        console.warn('⚠️ [Gamification] Error al procesar XP:', error);
        // Si falla el backend, revertir el estado local
        setUnderstoodInsights(prev => {
          const reverted = new Set(prev);
          reverted.delete(insightId);
          return reverted;
        });
      }
    } else {
      console.log('🔄 [UX] Insight ya marcado en esta sesión:', insightId);
    }
  };

  const handleApplyInsight = async (index, insight) => {
    if (appliedInsights.has(index)) return;
    setAppliedInsights(prev => new Set([...prev, index]));
    try {
      const gamificationAPI = getGamificationAPI();
      await gamificationAPI.recordApplyAIRecommendation(String(index), insight.next_action || insight.title);
    } catch (err) {
      console.warn('Error recording apply_ai_recommendation:', err);
    }
  };

  const getScoreColor = (score) => {
    if (score >= 800) return 'text-green-600 dark:text-green-400';
    if (score >= 600) return 'text-blue-600 dark:text-blue-400';
    if (score >= 400) return 'text-yellow-600 dark:text-yellow-400';
    return 'text-red-600 dark:text-red-400';
  };



  const getImpactIcon = (impact) => {
    switch (impact) {
      case 'high': return '🔥';
      case 'medium': return '⚡';
      case 'low': return '💡';
      default: return '📊';
    }
  };

  // Always show exactly 3 insights (backend generates exactly 3)
  const displayedInsights = insights.slice(0, 3);

  // Componente de salud financiera optimizado
  const HealthScoreDisplay = ({ score, maxScore = 1000, details = null, loading = false }) => {
    const [expanded, setExpanded] = useState(false);
    const percentage = (score / maxScore) * 100;

    const getScoreLevel = (score) => {
      if (score >= 800) return { level: 'Excelente', message: '¡Tu salud financiera es excepcional!', color: 'text-green-600 dark:text-green-400', bgColor: 'bg-green-50 dark:bg-green-900/20', borderColor: 'border-green-200 dark:border-green-700' };
      if (score >= 600) return { level: 'Bueno', message: 'Tu situación financiera es sólida', color: 'text-blue-600 dark:text-blue-400', bgColor: 'bg-blue-50 dark:bg-blue-900/20', borderColor: 'border-blue-200 dark:border-blue-700' };
      if (score >= 400) return { level: 'Regular', message: 'Hay oportunidades de mejora', color: 'text-yellow-600 dark:text-yellow-400', bgColor: 'bg-yellow-50 dark:bg-yellow-900/20', borderColor: 'border-yellow-200 dark:border-yellow-700' };
      return { level: 'Mejorable', message: 'Enfócate en las recomendaciones', color: 'text-red-600 dark:text-red-400', bgColor: 'bg-red-50 dark:bg-red-900/20', borderColor: 'border-red-200 dark:border-red-700' };
    };

    const { level, message, color, bgColor, borderColor } = getScoreLevel(score);

    const incomes = details?.current_month_incomes || 0;
    const expenses = details?.current_month_expenses || 0;
    const balance = details?.current_month_balance || 0;
    const productiveExpenses = details?.productive_expenses || 0;
    const consumptionExpenses = details?.consumption_expenses ?? (expenses - productiveExpenses);
    const ratio = incomes > 0 ? consumptionExpenses / incomes : null;
    const savingsRate = details?.savings_rate != null ? details.savings_rate : null;

    const cashFlowScore = details?.cash_flow_score ?? null;
    const planningScore = details?.planning_score ?? null;
    const consistencyScore = details?.consistency_score ?? null;
    const engagementScore = details?.engagement_score ?? null;
    const hasDimensions = cashFlowScore !== null;

    if (loading) {
      return (
        <div className="bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-xl p-6">
          <div className="text-center mb-4">
            <div className="w-16 h-16 bg-gray-300 dark:bg-gray-600 rounded-full animate-pulse mx-auto mb-2"></div>
            <div className="w-12 h-4 bg-gray-300 dark:bg-gray-600 rounded animate-pulse mx-auto mb-2"></div>
            <div className="w-20 h-6 bg-gray-300 dark:bg-gray-600 rounded-full animate-pulse mx-auto"></div>
          </div>
          <div className="w-full bg-gray-200 dark:bg-gray-600 rounded-full h-3 mb-4">
            <div className="h-full bg-gray-300 dark:bg-gray-500 rounded-full animate-pulse w-1/2"></div>
          </div>
          <div className="flex justify-between text-xs text-gray-300 dark:text-gray-500 mb-4">
            <span>0</span><span>250</span><span>500</span><span>750</span><span>1000</span>
          </div>
          <div className="text-center">
            <div className="w-48 h-4 bg-gray-300 dark:bg-gray-600 rounded animate-pulse mx-auto"></div>
          </div>
        </div>
      );
    }

    return (
      <div className={`${bgColor} ${borderColor} border rounded-xl p-6`}>
        {/* Score principal */}
        <div className="text-center mb-4">
          <div className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-1">{score}</div>
          <div className="text-gray-500 dark:text-gray-400 text-sm">/ {maxScore}</div>
          <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium mt-2 ${color} ${bgColor}`}>
            {level}
          </div>
        </div>

        {/* Barra de progreso */}
        <div className="w-full bg-gray-200 dark:bg-gray-600 rounded-full h-3 mb-4 overflow-hidden">
          <div
            className="h-full bg-gradient-to-r from-blue-500 to-green-500 rounded-full transition-all duration-1000 ease-out"
            style={{ width: `${percentage}%` }}
          />
        </div>

        {/* Etiquetas de referencia */}
        <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mb-4">
          <span>0</span><span>250</span><span>500</span><span>750</span><span>1000</span>
        </div>

        {/* Mensaje */}
        <div className="text-center">
          <p className="text-gray-600 dark:text-gray-400 text-sm leading-relaxed">{message}</p>
        </div>

        {/* Desplegable de transparencia */}
        {details && (
          <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-600">
            <button
              onClick={() => setExpanded(!expanded)}
              className="w-full flex items-center justify-between text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors"
            >
              <span className="flex items-center space-x-1">
                <FaCalculator className="w-3 h-3" />
                <span>¿Cómo se calculó este puntaje?</span>
              </span>
              {expanded ? <FaChevronUp className="w-3 h-3" /> : <FaChevronDown className="w-3 h-3" />}
            </button>

            {expanded && (
              <div className="mt-4 space-y-3">
                {/* Datos del mes */}
                <div className="grid grid-cols-3 gap-2 text-center">
                  <div className="bg-white/60 dark:bg-gray-700/60 rounded-lg p-3">
                    <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">Ingresos</div>
                    <div className="text-sm font-semibold text-green-600 dark:text-green-400">{formatCurrency(incomes)}</div>
                  </div>
                  <div className="bg-white/60 dark:bg-gray-700/60 rounded-lg p-3">
                    <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">Egresos totales</div>
                    <div className="text-sm font-semibold text-red-600 dark:text-red-400">{formatCurrency(expenses)}</div>
                  </div>
                  <div className="bg-white/60 dark:bg-gray-700/60 rounded-lg p-3">
                    <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">Balance</div>
                    <div className={`text-sm font-semibold ${balance >= 0 ? 'text-blue-600 dark:text-blue-400' : 'text-red-600 dark:text-red-400'}`}>
                      {formatCurrency(balance)}
                    </div>
                  </div>
                </div>

                {/* Desglose de egresos: consumo vs productivos */}
                {productiveExpenses > 0 && (
                  <div className="bg-white/60 dark:bg-gray-700/60 rounded-lg p-3 space-y-2">
                    <div className="text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Desglose de egresos</div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-gray-500 dark:text-gray-400">🛒 Consumo</span>
                      <span className="font-semibold text-orange-600 dark:text-orange-400">{formatCurrency(consumptionExpenses)}</span>
                    </div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-gray-500 dark:text-gray-400">📈 Inversión / activos</span>
                      <span className="font-semibold text-green-600 dark:text-green-400">{formatCurrency(productiveExpenses)}</span>
                    </div>
                    <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">
                      El puntaje se calcula sobre el consumo neto, excluyendo inversiones y activos.
                    </p>
                  </div>
                )}

                {/* Ratio consumo/ingresos */}
                {ratio !== null && (
                  <div className="bg-white/60 dark:bg-gray-700/60 rounded-lg p-3">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-xs text-gray-500 dark:text-gray-400">
                        Ratio consumo / ingresos{productiveExpenses > 0 ? ' (neto)' : ''}
                      </span>
                      <span className="text-sm font-bold text-gray-700 dark:text-gray-200">
                        {(ratio * 100).toFixed(1)}%
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 dark:bg-gray-600 rounded-full h-2 overflow-hidden">
                      <div
                        className={`h-2 rounded-full transition-all duration-700 ${ratio < 0.5 ? 'bg-green-500' : ratio < 0.7 ? 'bg-blue-500' : ratio < 0.9 ? 'bg-yellow-500' : 'bg-red-500'}`}
                        style={{ width: `${Math.min(ratio * 100, 100)}%` }}
                      />
                    </div>
                  </div>
                )}

                {/* Desglose multi-dimensional */}
                {hasDimensions && (
                  <div className="bg-white/60 dark:bg-gray-700/60 rounded-lg p-3">
                    <div className="text-xs font-medium text-gray-600 dark:text-gray-300 mb-2">Desglose del puntaje</div>
                    <div className="space-y-2">
                      {[
                        { label: '💵 Flujo de caja', value: cashFlowScore, weight: '40%' },
                        { label: '🗓️ Planificación', value: planningScore, weight: '30%' },
                        { label: '🔁 Consistencia', value: consistencyScore, weight: '20%' },
                        { label: '🤖 Engagement IA', value: engagementScore, weight: '10%' },
                      ].map(({ label, value, weight }) => (
                        <div key={label}>
                          <div className="flex items-center justify-between text-xs mb-1">
                            <span className="text-gray-600 dark:text-gray-400">{label}</span>
                            <span className="font-semibold text-gray-700 dark:text-gray-200">{value}/100 <span className="text-gray-400 font-normal">({weight})</span></span>
                          </div>
                          <div className="w-full bg-gray-200 dark:bg-gray-600 rounded-full h-1.5 overflow-hidden">
                            <div
                              className={`h-1.5 rounded-full transition-all duration-700 ${value >= 70 ? 'bg-green-500' : value >= 40 ? 'bg-blue-500' : 'bg-yellow-500'}`}
                              style={{ width: `${Math.min(value, 100)}%` }}
                            />
                          </div>
                        </div>
                      ))}
                    </div>
                    <p className="text-xs text-gray-400 dark:text-gray-500 mt-2">
                      El puntaje combina tu flujo de caja, planificación, consistencia de uso y engagement con IA.
                    </p>
                  </div>
                )}

                {/* Tasa de ahorro real */}
                {savingsRate !== null && (
                  <div className="bg-white/60 dark:bg-gray-700/60 rounded-lg p-3">
                    <div className="flex items-center justify-between text-xs">
                      <div>
                        <span className="text-gray-500 dark:text-gray-400">Tasa de ahorro real del mes</span>
                        {productiveExpenses > 0 && (
                          <p className="text-gray-400 dark:text-gray-500 mt-0.5">
                            Balance + inversiones / ingresos
                          </p>
                        )}
                      </div>
                      <span className={`font-semibold text-sm ${savingsRate >= 20 ? 'text-green-600 dark:text-green-400' : savingsRate >= 10 ? 'text-blue-600 dark:text-blue-400' : 'text-yellow-600 dark:text-yellow-400'}`}>
                        {savingsRate.toFixed(1)}%
                      </span>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    );
  };

  // Loading Skeleton Component
  const InsightSkeleton = () => (
    <div className="bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-xl p-4 animate-pulse">
      <div className="flex items-start space-x-3">
        <div className="w-8 h-8 bg-gray-300 dark:bg-gray-600 rounded"></div>
        <div className="flex-1 space-y-2">
          <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-3/4"></div>
          <div className="h-3 bg-gray-300 dark:bg-gray-600 rounded w-full"></div>
          <div className="h-3 bg-gray-300 dark:bg-gray-600 rounded w-5/6"></div>
          <div className="flex justify-between items-center mt-3">
            <div className="h-6 bg-gray-300 dark:bg-gray-600 rounded w-20"></div>
            <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-16"></div>
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <div className="space-y-6">
      {/* Salud Financiera */}
      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm dark:shadow-gray-900/20 border border-gray-100 dark:border-gray-700 p-6 md:p-8">
        <div className="text-center mb-6">
          <h2 className="text-xl md:text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Tu salud financiera</h2>
          <p className="text-gray-600 dark:text-gray-400">Evaluación inteligente powered by IA</p>
        </div>
        <HealthScoreDisplay score={healthScore} details={healthDetails} loading={healthScoreLoading} />
      </div>

      {/* Tabs Navigation */}
      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm dark:shadow-gray-900/20 border border-gray-100 dark:border-gray-700 overflow-hidden">
        <div className="border-b border-gray-200 dark:border-gray-700">
          <nav className="flex space-x-8 px-6" aria-label="Tabs">
            <button
              onClick={() => setActiveTab('monthly')}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === 'monthly'
                  ? 'border-purple-500 text-purple-600 dark:text-purple-400'
                  : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
              }`}
            >
              <div className="flex items-center space-x-2">
                <span className="text-base">📅</span>
                <span>Reporte del Mes</span>
              </div>
            </button>
            <button
              onClick={() => setActiveTab('education')}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === 'education'
                  ? 'border-indigo-500 text-indigo-600 dark:text-indigo-400'
                  : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
              }`}
            >
              <div className="flex items-center space-x-2">
                <span className="text-base">📚</span>
                <span>Educación</span>
              </div>
            </button>
            <button
              onClick={() => setActiveTab('purchase')}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === 'purchase'
                  ? 'border-green-500 text-green-600 dark:text-green-400'
                  : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
              }`}
            >
              <div className="flex items-center space-x-2">
                <FaShoppingCart className="w-4 h-4" />
                <span>¿Puedo comprarlo?</span>
                <span className="bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-300 text-xs px-2 py-1 rounded-full font-medium">
                  BETA
                </span>
              </div>
            </button>
          </nav>
        </div>

        <div className="p-6">
          {activeTab === 'monthly' && <MonthlyCoachingTab />}
          {activeTab === 'education' && <EducationTab />}

          {activeTab === 'purchase' && (
            <div className="space-y-5">
              {/* Datos financieros resumen */}
              {dashboardData ? (
                <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                  {[
                    { label: 'Ingresos del mes', value: dashboardData.current_month_incomes || 0, color: 'text-green-700 dark:text-green-400' },
                    { label: 'Gastos del mes', value: dashboardData.current_month_expenses || 0, color: 'text-red-600 dark:text-red-400' },
                    { label: 'Balance actual', value: dashboardData.current_month_balance || 0, color: 'text-blue-700 dark:text-blue-400' },
                    { label: 'Disponible', value: (dashboardData.current_month_incomes || 0) - (dashboardData.current_month_expenses || 0), color: 'text-indigo-700 dark:text-indigo-400' },
                  ].map(({ label, value, color }) => (
                    <div key={label} className="bg-gray-50 dark:bg-gray-700 rounded-xl p-3 text-center">
                      <p className="text-xs text-gray-500 dark:text-gray-400 mb-1">{label}</p>
                      <p className={`font-bold text-base ${color}`}>${Math.abs(value).toLocaleString()}</p>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="flex items-center space-x-2 text-gray-500 dark:text-gray-400 text-sm">
                  <FaSpinner className="w-4 h-4 animate-spin" />
                  <span>Cargando datos financieros...</span>
                </div>
              )}

              {/* Form + Result — side by side on md+ */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                {/* Left: form */}
                <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-5 space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                      ¿Qué quieres comprar? *
                    </label>
                    <input
                      type="text"
                      value={purchaseForm.itemName}
                      onChange={(e) => setPurchaseForm({...purchaseForm, itemName: e.target.value})}
                      placeholder="Ej: iPhone 15, Notebook, Vacaciones..."
                      className="w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:placeholder-gray-400 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                      Monto *
                    </label>
                    <div className="relative">
                      <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500 dark:text-gray-400 text-sm">$</span>
                      <input
                        type="number"
                        value={purchaseForm.amount}
                        onChange={(e) => setPurchaseForm({...purchaseForm, amount: e.target.value})}
                        placeholder="150000"
                        className="w-full pl-7 pr-4 py-2.5 border border-gray-300 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:placeholder-gray-400 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                      />
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                      Descripción (opcional)
                    </label>
                    <textarea
                      value={purchaseForm.description}
                      onChange={(e) => setPurchaseForm({...purchaseForm, description: e.target.value})}
                      placeholder="¿Para qué lo necesitas? ¿Es urgente?"
                      rows="2"
                      className="w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:placeholder-gray-400 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm resize-none"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Tipo de pago
                    </label>
                    <div className="space-y-1.5">
                      {paymentTypes.map(type => (
                        <label key={type.value} className="flex items-center space-x-2 cursor-pointer p-1.5 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors">
                          <input
                            type="checkbox"
                            value={type.value}
                            checked={purchaseForm.paymentTypes.includes(type.value)}
                            onChange={(e) => handlePaymentTypeChange(type.value, e.target.checked)}
                            className="w-4 h-4 text-blue-500 focus:ring-blue-500 rounded"
                          />
                          <span className="text-sm text-gray-700 dark:text-gray-300">{type.label}</span>
                        </label>
                      ))}
                    </div>
                  </div>

                  <label className="flex items-center space-x-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={purchaseForm.isNecessary}
                      onChange={(e) => setPurchaseForm({...purchaseForm, isNecessary: e.target.checked})}
                      className="w-4 h-4 text-blue-500 focus:ring-blue-500 rounded"
                    />
                    <span className="text-sm text-gray-700 dark:text-gray-300">Es una necesidad urgente</span>
                  </label>

                  <button
                    onClick={analyzePurchase}
                    disabled={purchaseLoading || !purchaseForm.itemName || !purchaseForm.amount || purchaseForm.paymentTypes.length === 0}
                    className="w-full px-6 py-3 bg-green-500 dark:bg-green-600 text-white rounded-lg hover:bg-green-600 dark:hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2 text-sm font-medium"
                  >
                    {purchaseLoading ? (
                      <>
                        <FaSpinner className="w-4 h-4 animate-spin" />
                        <span>Analizando con IA...</span>
                      </>
                    ) : (
                      <>
                        <FaCalculator className="w-4 h-4" />
                        <span>Analizar compra</span>
                      </>
                    )}
                  </button>
                </div>

                {/* Right: result */}
                <div>
                  {purchaseError && (
                    <div className="bg-red-50 dark:bg-red-900/20 border-l-4 border-red-400 dark:border-red-500 rounded-lg p-4">
                      <div className="flex">
                        <FaExclamationTriangle className="w-5 h-5 text-red-400 dark:text-red-500 flex-shrink-0" />
                        <div className="ml-3">
                          <p className="text-red-700 dark:text-red-300 font-medium text-sm">Error en el análisis</p>
                          <p className="text-red-600 dark:text-red-400 text-xs mt-1">{purchaseError}</p>
                        </div>
                      </div>
                    </div>
                  )}

                  {!purchaseAnalysis && !purchaseError && (
                    <div className="flex flex-col items-center justify-center h-full py-12 text-center text-gray-400 dark:text-gray-500">
                      <FaCalculator className="w-10 h-10 mb-3 opacity-30" />
                      <p className="text-sm">Completá el formulario y<br />analizá si podés hacer la compra</p>
                    </div>
                  )}

                  {purchaseAnalysis && (
                    <div className={`rounded-xl p-5 border-2 space-y-4 ${
                      purchaseAnalysis.can_buy
                        ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-700'
                        : 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-700'
                    }`}>
                      <div className="flex items-center space-x-3">
                        {purchaseAnalysis.can_buy ? (
                          <FaCheckCircle className="w-7 h-7 text-green-600 dark:text-green-400 flex-shrink-0" />
                        ) : (
                          <FaExclamationTriangle className="w-7 h-7 text-red-600 dark:text-red-400 flex-shrink-0" />
                        )}
                        <div>
                          <h4 className={`font-bold ${purchaseAnalysis.can_buy ? 'text-green-800 dark:text-green-300' : 'text-red-800 dark:text-red-300'}`}>
                            {purchaseAnalysis.can_buy ? '✅ ¡Podés comprarlo!' : '❌ Te recomendamos esperar'}
                          </h4>
                          <p className="text-xs text-gray-500 dark:text-gray-400">
                            Confianza: {Math.round(purchaseAnalysis.confidence * 100)}% · Impacto: {purchaseAnalysis.impact_score} pts
                          </p>
                        </div>
                      </div>

                      <p className="text-sm text-gray-700 dark:text-gray-300 leading-relaxed bg-white dark:bg-gray-800 rounded-lg p-3">
                        {purchaseAnalysis.reasoning}
                      </p>

                      {purchaseAnalysis.alternatives && purchaseAnalysis.alternatives.length > 0 && (
                        <div>
                          <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2 flex items-center">
                            <FaLightbulb className="w-3 h-3 mr-1 text-yellow-500" /> Alternativas
                          </p>
                          <ul className="space-y-1.5">
                            {purchaseAnalysis.alternatives.map((alt, index) => (
                              <li key={index} className="flex items-start text-sm text-gray-600 dark:text-gray-400">
                                <FaChevronRight className="w-3 h-3 mr-2 mt-0.5 text-gray-400 flex-shrink-0" />
                                <span>{alt}</span>
                              </li>
                            ))}
                          </ul>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Footer con información de actualización */}
      <div className="text-center">
        <div className="inline-flex items-center px-4 py-2 bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 rounded-lg text-sm">
          <FaBrain className="w-4 h-4 mr-2" />
          <span>
            Analizamos tu situación financiera una vez por día
            {lastEvaluationDate && (
              <span className="ml-1">
                • Último análisis: {lastEvaluationDate.toLocaleDateString('es-ES', { 
                  day: 'numeric', 
                  month: 'short',
                  hour: '2-digit',
                  minute: '2-digit'
                })}
              </span>
            )}
          </span>
        </div>
      </div>
    </div>
  );
};

export default AIInsights; 