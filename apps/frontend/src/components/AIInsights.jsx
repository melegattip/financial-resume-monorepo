import React, { useState, useEffect, useCallback } from 'react';
import { FaBrain, FaSpinner, FaRedo, FaLightbulb, FaShoppingCart, FaCheck, FaChevronRight, FaCalculator, FaExclamationTriangle, FaCheckCircle, FaChevronDown, FaChevronUp, FaBullseye } from 'react-icons/fa';
import { aiAPI, savingsGoalsAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import { usePeriod } from '../contexts/PeriodContext';
import { useGamification } from '../contexts/GamificationContext';

const AIInsights = () => {
  const { user, isAuthenticated } = useAuth();
  const { updateAvailableData } = usePeriod();
  const { recordInsightViewed, recordInsightUnderstood, recordSuggestionUsed } = useGamification();
  const [insights, setInsights] = useState([]);
  const [purchaseAnalysis, setPurchaseAnalysis] = useState(null);
  const [loading, setLoading] = useState(false);
  const [purchaseLoading, setPurchaseLoading] = useState(false);
  const [error, setError] = useState(null);
  const [purchaseError, setPurchaseError] = useState(null);
  const [healthScore, setHealthScore] = useState(0);
  const [healthScoreLoading, setHealthScoreLoading] = useState(false);
  const [lastEvaluationDate, setLastEvaluationDate] = useState(null);
  const [dashboardData, setDashboardData] = useState(null);
  const [savingsGoals, setSavingsGoals] = useState([]);
  
  // Estados para Progressive Disclosure
  const [showAllInsights, setShowAllInsights] = useState(false);
  const [activeTab, setActiveTab] = useState('insights'); // 'insights' | 'purchase'
  
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

  // Cargar datos del dashboard para obtener información financiera real
  const loadDashboardData = useCallback(async () => {
    try {
      // Usar el servicio de datos para obtener información completa
      const dataService = (await import('../services/dataService')).default;
      const dashboardData = await dataService.loadDashboardData({}, isAuthenticated && user);
      
      setDashboardData(dashboardData);
      
      // Actualizar datos disponibles en el contexto de períodos
      updateAvailableData(
        dashboardData.allExpenses || dashboardData.expenses || [], 
        dashboardData.allIncomes || dashboardData.incomes || []
      );
      
      // Actualizar el formulario de compra con datos reales
      const newFormData = {
        currentBalance: dashboardData.balance || 0,
        monthlyIncome: dashboardData.totalIncome || 0,
        monthlyExpenses: dashboardData.totalExpenses || 0,
        savingsGoal: dashboardData.savings_goal || 50000
      };
      
      setPurchaseForm(prev => ({
        ...prev,
        ...newFormData
      }));
      
      console.log('✅ AIInsights: Datos del dashboard y períodos actualizados');
    } catch (error) {
      console.error('Error loading dashboard data:', error);
      // Mantener valores por defecto en caso de error
    }
  }, [updateAvailableData, isAuthenticated, user]);

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

  // Función para cargar el health score
  const loadHealthScore = useCallback(async () => {
    if (!isAuthenticated) return;

    setHealthScoreLoading(true);
    try {
      const response = await aiAPI.getHealthScore();
      setHealthScore(response.health_score || 0);
    } catch (err) {
      console.error('Error loading health score:', err.message);
      // Mantener el valor por defecto de 0 en caso de error
      setHealthScore(0);
    } finally {
      setHealthScoreLoading(false);
    }
  }, [isAuthenticated, user?.email]);

  const loadAIInsights = useCallback(async () => {
    if (!isAuthenticated) {
      setError('Debes iniciar sesión para ver el análisis inteligente');
      return;
    }

    setLoading(true);
    setError(null);
    try {
      console.log('🔍 Cargando análisis inteligente para usuario:', user?.email);
      const response = await aiAPI.getInsights();
      const newInsights = response.insights || [];
      setInsights(newInsights);
      
      // Usar el timestamp del backend (generated_at)
      const backendTimestamp = response.generated_at ? new Date(response.generated_at) : new Date();
      setLastEvaluationDate(backendTimestamp);
      console.log('💾 Análisis cargado desde backend - Timestamp:', backendTimestamp.toISOString());
    } catch (err) {
      console.error('Error loading AI insights:', err.message);
      setError('Error conectando con GPT-4. Usando datos de ejemplo.');
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
  }, [isAuthenticated, user?.email]);

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
    loadAIInsightsSimple();
    loadHealthScore();
    loadDashboardData();
    loadSavingsGoals();
  }, [isAuthenticated, loadAIInsightsSimple, loadHealthScore, loadDashboardData, loadSavingsGoals]);

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

  // Progressive Disclosure: mostrar solo los primeros 3 insights
  const displayedInsights = showAllInsights ? insights : insights.slice(0, 3);

  // Componente de salud financiera optimizado
  const HealthScoreDisplay = ({ score, maxScore = 1000, loading = false }) => {
    const percentage = (score / maxScore) * 100;
    
    const getScoreLevel = (score) => {
      if (score >= 800) return { level: 'Excelente', message: '¡Tu salud financiera es excepcional!', color: 'text-green-600 dark:text-green-400', bgColor: 'bg-green-50 dark:bg-green-900/20', borderColor: 'border-green-200 dark:border-green-700' };
      if (score >= 600) return { level: 'Bueno', message: 'Tu situación financiera es sólida', color: 'text-blue-600 dark:text-blue-400', bgColor: 'bg-blue-50 dark:bg-blue-900/20', borderColor: 'border-blue-200 dark:border-blue-700' };
      if (score >= 400) return { level: 'Regular', message: 'Hay oportunidades de mejora', color: 'text-yellow-600 dark:text-yellow-400', bgColor: 'bg-yellow-50 dark:bg-yellow-900/20', borderColor: 'border-yellow-200 dark:border-yellow-700' };
      return { level: 'Mejorable', message: 'Enfócate en las recomendaciones', color: 'text-red-600 dark:text-red-400', bgColor: 'bg-red-50 dark:bg-red-900/20', borderColor: 'border-red-200 dark:border-red-700' };
    };

    const { level, message, color, bgColor, borderColor } = getScoreLevel(score);

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
          <span>0</span>
          <span>250</span>
          <span>500</span>
          <span>750</span>
          <span>1000</span>
        </div>

        {/* Mensaje */}
        <div className="text-center">
          <p className="text-gray-600 dark:text-gray-400 text-sm leading-relaxed">{message}</p>
        </div>
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
        <HealthScoreDisplay score={healthScore} loading={healthScoreLoading} />
      </div>

      {/* Tabs Navigation */}
      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm dark:shadow-gray-900/20 border border-gray-100 dark:border-gray-700 overflow-hidden">
        <div className="border-b border-gray-200 dark:border-gray-700">
          <nav className="flex space-x-8 px-6" aria-label="Tabs">
            <button
              onClick={() => setActiveTab('insights')}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === 'insights'
                  ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                  : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
              }`}
            >
              <div className="flex items-center space-x-2">
                <FaLightbulb className="w-4 h-4" />
                <span>Recomendaciones</span>
                {insights.length > 0 && (
                  <span className="bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-300 text-xs px-2 py-1 rounded-full">
                    {insights.length}
                  </span>
                )}
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
          {activeTab === 'insights' && (
            <div className="space-y-6">
              {loading ? (
                <div className="space-y-4">
                  {[...Array(3)].map((_, i) => (
                    <InsightSkeleton key={i} />
                  ))}
                </div>
              ) : error && insights.length === 0 ? (
                <div className="text-center py-12">
                  <FaSpinner className="w-16 h-16 mx-auto mb-4 text-gray-400 dark:text-gray-500" />
                  <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">Sin conexión</h3>
                  <p className="text-gray-600 dark:text-gray-400 mb-6 max-w-md mx-auto">{error}</p>
                  <button
                    onClick={loadAIInsights}
                    className="inline-flex items-center px-6 py-3 bg-blue-500 dark:bg-blue-600 text-white rounded-lg hover:bg-blue-600 dark:hover:bg-blue-700 transition-colors font-medium"
                  >
                    <FaRedo className="w-4 h-4 mr-2" />
                    Reintentar análisis
                  </button>
                </div>
              ) : insights.length === 0 ? (
                <div className="text-center py-12">
                  <FaBullseye className="w-16 h-16 mx-auto mb-4 text-gray-400 dark:text-gray-500" />
                  <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">¡Perfecto!</h3>
                  <p className="text-gray-600 dark:text-gray-400 max-w-md mx-auto">
                    Tu situación financiera está tan bien que no tenemos recomendaciones urgentes. 
                    Sigue así y revisa periódicamente.
                  </p>
                </div>
              ) : (
                <>
                  <div className="grid gap-4">
                    {displayedInsights.map((insight, index) => (
                      <div
                        key={index}
                        className="bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-xl p-4 hover:shadow-md hover:border-gray-300 dark:hover:border-gray-500 transition-all cursor-pointer group"
                        onMouseEnter={() => {/* no-op: evitar múltiples registros por hover */}}
                        onFocus={() => {/* no-op */}}
                        onClick={() => handleViewInsight(index, insight.title)}
                      >
                        <div className="flex items-start space-x-3">
                          <div className="text-2xl group-hover:scale-110 transition-transform">
                            {getImpactIcon(insight.impact)}
                          </div>
                          <div className="flex-1 min-w-0">
                            <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-2 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                              {insight.title}
                            </h3>
                            <p className="text-gray-600 dark:text-gray-400 text-sm leading-relaxed mb-3">
                              {insight.description}
                            </p>
                            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-2 sm:space-y-0">
                              <div className="flex items-center space-x-3">
                                <span className="inline-flex items-center px-2 py-1 bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-lg text-xs font-medium">
                                  📊 {insight.category}
                                </span>
                                <span className={`text-sm font-bold ${getScoreColor(insight.score)}`}>
                                  {insight.score} pts
                                </span>
                              </div>
                              <div className="flex items-center space-x-2">
                                {!understoodInsights.has(index) && (
                                  <button
                                    onClick={() => handleUnderstandInsight(index, insight.title)}
                                    className="inline-flex items-center px-3 py-1.5 bg-blue-500 dark:bg-blue-600 text-white text-xs rounded-lg hover:bg-blue-600 dark:hover:bg-blue-700 transition-colors font-medium"
                                  >
                                    <FaCheck className="w-3 h-3 mr-1" />
                                    Marcar como revisado
                                  </button>
                                )}
                                {understoodInsights.has(index) && (
                                  <div className="inline-flex items-center px-3 py-1.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300 text-xs rounded-lg font-medium">
                                    <FaCheckCircle className="w-3 h-3 mr-1" />
                                    ¡Revisado!
                                  </div>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>

                  {/* Progressive Disclosure */}
                  {insights.length > 3 && (
                    <div className="text-center pt-4">
                      <button
                        onClick={() => setShowAllInsights(!showAllInsights)}
                        className="inline-flex items-center px-6 py-3 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors font-medium"
                      >
                        {showAllInsights ? (
                          <>
                            <FaChevronUp className="w-4 h-4 mr-2" />
                            Mostrar menos recomendaciones
                          </>
                        ) : (
                          <>
                            <FaChevronDown className="w-4 h-4 mr-2" />
                            Ver todas las recomendaciones ({insights.length - 3} más)
                          </>
                        )}
                      </button>
                    </div>
                  )}
                </>
              )}
            </div>
          )}

          {activeTab === 'purchase' && (
            <div className="max-w-2xl mx-auto space-y-6">
              <div className="text-center mb-6">
                <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Análisis de Compra Inteligente</h3>
                <p className="text-gray-600 dark:text-gray-400">Te ayudamos a tomar decisiones financieras informadas</p>
              </div>

              {/* Información financiera automática */}
              {dashboardData ? (
                <div className="bg-blue-50 dark:bg-blue-900/20 rounded-xl p-4 mb-4">
                  <h4 className="font-medium text-blue-900 dark:text-blue-300 mb-3 flex items-center">
                    <FaCalculator className="w-4 h-4 mr-2" />
                    Datos financieros actuales (automáticos)
                  </h4>
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="text-blue-700 dark:text-blue-400">Balance actual:</span>
                      <p className="font-semibold text-blue-900 dark:text-blue-200">${(dashboardData?.balance || dashboardData?.Metrics?.Balance || dashboardData?.metrics?.balance || 0).toLocaleString()}</p>
                    </div>
                    <div>
                      <span className="text-blue-700 dark:text-blue-400">Ingresos mensuales:</span>
                      <p className="font-semibold text-blue-900 dark:text-blue-200">${(dashboardData?.totalIncome || dashboardData?.Metrics?.TotalIncome || dashboardData?.metrics?.total_income || 0).toLocaleString()}</p>
                    </div>
                    <div>
                      <span className="text-blue-700 dark:text-blue-400">Gastos mensuales:</span>
                      <p className="font-semibold text-blue-900 dark:text-blue-200">${(dashboardData?.totalExpenses || dashboardData?.Metrics?.TotalExpenses || dashboardData?.metrics?.total_expenses || 0).toLocaleString()}</p>
                    </div>
                    <div>
                      <span className="text-blue-700 dark:text-blue-400">Disponible/mes:</span>
                      <p className="font-semibold text-green-700 dark:text-green-400">${((dashboardData?.totalIncome || dashboardData?.Metrics?.TotalIncome || dashboardData?.metrics?.total_income || 0) - (dashboardData?.totalExpenses || dashboardData?.Metrics?.TotalExpenses || dashboardData?.metrics?.total_expenses || 0)).toLocaleString()}</p>
                    </div>
                  </div>
                  <div className="flex items-center justify-between mt-2">
                    <p className="text-xs text-blue-600 dark:text-blue-400">
                      💡 Estos datos se calculan automáticamente basándose en tus transacciones
                    </p>
                    <button
                      onClick={loadDashboardData}
                      className="text-xs bg-blue-500 dark:bg-blue-600 text-white px-2 py-1 rounded hover:bg-blue-600 dark:hover:bg-blue-700 transition-colors"
                    >
                      🔄 Actualizar
                    </button>
                  </div>
                </div>
              ) : (
                <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-4 mb-4">
                  <div className="flex items-center justify-center space-x-2 text-gray-600 dark:text-gray-400">
                    <FaSpinner className="w-4 h-4 animate-spin" />
                    <span className="text-sm">Cargando datos financieros...</span>
                  </div>
                </div>
              )}

              {/* Formulario optimizado para móvil */}
              <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-4 md:p-6 space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    ¿Qué quieres comprar? *
                  </label>
                  <input
                    type="text"
                    value={purchaseForm.itemName}
                    onChange={(e) => setPurchaseForm({...purchaseForm, itemName: e.target.value})}
                    placeholder="Ej: iPhone 15, Notebook, Vacaciones..."
                    className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:placeholder-gray-400 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-base"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Monto *
                  </label>
                  <div className="relative">
                    <span className="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-500 dark:text-gray-400">$</span>
                    <input
                      type="number"
                      value={purchaseForm.amount}
                      onChange={(e) => setPurchaseForm({...purchaseForm, amount: e.target.value})}
                      placeholder="150000"
                      className="w-full pl-8 pr-4 py-3 border border-gray-300 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:placeholder-gray-400 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-base"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Descripción (opcional)
                  </label>
                  <textarea
                    value={purchaseForm.description}
                    onChange={(e) => setPurchaseForm({...purchaseForm, description: e.target.value})}
                    placeholder="¿Para qué lo necesitas? ¿Es urgente?"
                    rows="3"
                    className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:placeholder-gray-400 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-base resize-none"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                    Tipo de pago (puedes seleccionar varios)
                  </label>
                  <div className="space-y-2">
                    {paymentTypes.map(type => (
                      <label key={type.value} className="flex items-center space-x-3 cursor-pointer p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors">
                        <input
                          type="checkbox"
                          value={type.value}
                          checked={purchaseForm.paymentTypes.includes(type.value)}
                          onChange={(e) => handlePaymentTypeChange(type.value, e.target.checked)}
                          className="w-4 h-4 text-blue-500 focus:ring-blue-500 rounded"
                        />
                        <span className="text-sm text-gray-700 dark:text-gray-300 font-medium">{type.label}</span>
                      </label>
                    ))}
                  </div>
                  {purchaseForm.paymentTypes.length === 0 && (
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      Selecciona al menos un tipo de pago
                    </p>
                  )}
                </div>

                <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-3">
                  <label className="flex items-start space-x-3 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={purchaseForm.isNecessary}
                      onChange={(e) => setPurchaseForm({...purchaseForm, isNecessary: e.target.checked})}
                      className="w-4 h-4 text-blue-500 focus:ring-blue-500 rounded mt-0.5"
                    />
                    <div>
                      <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                        Es una necesidad urgente
                      </span>
                      <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                        Marca esto solo si es esencial para tu trabajo, salud o seguridad
                      </p>
                    </div>
                  </label>
                </div>
                
                <button
                  onClick={analyzePurchase}
                  disabled={purchaseLoading || !purchaseForm.itemName || !purchaseForm.amount || purchaseForm.paymentTypes.length === 0}
                  className="w-full px-6 py-4 bg-green-500 dark:bg-green-600 text-white rounded-lg hover:bg-green-600 dark:hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2 text-base font-medium shadow-lg"
                >
                  {purchaseLoading ? (
                    <>
                      <FaSpinner className="w-5 h-5 animate-spin" />
                      <span>Analizando con IA...</span>
                    </>
                  ) : (
                    <>
                      <FaCalculator className="w-5 h-5" />
                      <span>Analizar compra</span>
                    </>
                  )}
                </button>
              </div>

              {/* Resultado del análisis mejorado */}
              {purchaseError && (
                <div className="bg-red-50 dark:bg-red-900/20 border-l-4 border-red-400 dark:border-red-500 rounded-lg p-4">
                  <div className="flex">
                    <FaExclamationTriangle className="w-5 h-5 text-red-400 dark:text-red-500" />
                    <div className="ml-3">
                      <p className="text-red-700 dark:text-red-300 font-medium">Error en el análisis</p>
                      <p className="text-red-600 dark:text-red-400 text-sm mt-1">{purchaseError}</p>
                    </div>
                  </div>
                </div>
              )}

              {purchaseAnalysis && (
                <div className={`rounded-xl p-6 border-2 shadow-lg ${
                  purchaseAnalysis.can_buy 
                    ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-700' 
                    : 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-700'
                }`}>
                  <div className="flex items-start space-x-4 mb-4">
                    {purchaseAnalysis.can_buy ? (
                      <div className="p-2 bg-green-100 dark:bg-green-900/30 rounded-full">
                        <FaCheckCircle className="w-6 h-6 text-green-600 dark:text-green-400" />
                      </div>
                    ) : (
                      <div className="p-2 bg-red-100 dark:bg-red-900/30 rounded-full">
                        <FaExclamationTriangle className="w-6 h-6 text-red-600 dark:text-red-400" />
                      </div>
                    )}
                    <div className="flex-1">
                      <h4 className={`text-lg font-bold mb-1 ${
                        purchaseAnalysis.can_buy ? 'text-green-800 dark:text-green-300' : 'text-red-800 dark:text-red-300'
                      }`}>
                        {purchaseAnalysis.can_buy ? '✅ ¡Puedes comprarlo!' : '❌ Te recomendamos esperar'}
                      </h4>
                      <p className={`text-sm ${
                        purchaseAnalysis.can_buy ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
                      }`}>
                        Confianza del análisis: {Math.round(purchaseAnalysis.confidence * 100)}%
                      </p>
                    </div>
                  </div>
                  
                  <div className={`bg-white dark:bg-gray-800 rounded-lg p-4 mb-4 ${
                    purchaseAnalysis.can_buy ? 'border border-green-200 dark:border-green-700' : 'border border-red-200 dark:border-red-700'
                  }`}>
                    <p className="text-gray-700 dark:text-gray-300 leading-relaxed">
                      {purchaseAnalysis.reasoning}
                    </p>
                  </div>
                  
                  {purchaseAnalysis.alternatives && purchaseAnalysis.alternatives.length > 0 && (
                    <div className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
                      <h5 className="text-sm font-semibold text-gray-900 dark:text-gray-100 mb-3 flex items-center">
                        <FaLightbulb className="w-4 h-4 mr-2 text-yellow-500 dark:text-yellow-400" />
                        Alternativas sugeridas
                      </h5>
                      <ul className="space-y-2">
                        {purchaseAnalysis.alternatives.map((alt, index) => (
                          <li key={index} className="flex items-center text-sm text-gray-600 dark:text-gray-400">
                            <FaChevronRight className="w-4 h-4 mr-2 text-gray-400 dark:text-gray-500" />
                            <span>{alt}</span>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                  
                  <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <div className="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
                      <span>Impacto en tu presupuesto</span>
                      <span className="font-medium">{purchaseAnalysis.impact_score} pts</span>
                    </div>
                  </div>
                </div>
              )}
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