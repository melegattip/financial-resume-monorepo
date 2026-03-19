import React, { useState, useEffect, useCallback } from 'react';
import { format, subMonths, startOfMonth, endOfMonth } from 'date-fns';
import { es } from 'date-fns/locale';
import { useNavigate } from 'react-router-dom';
import { FaSpinner, FaRedo, FaArrowRight } from 'react-icons/fa';
import { aiAPI, analyticsAPI, budgetsAPI, savingsGoalsAPI } from '../../../services/api';
import { getGamificationAPI } from '../../../services/gamificationAPI';

const sentimentConfig = {
  positivo:   { label: 'Mes Positivo',   bg: 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-700', badge: 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300', emoji: '✨' },
  neutral:    { label: 'Mes Neutral',    bg: 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-700', badge: 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300', emoji: '⚖️' },
  desafiante: { label: 'Mes Desafiante', bg: 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-700', badge: 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300', emoji: '💪' },
};

const MonthlyCoachingTab = () => {
  const navigate = useNavigate();
  const [report, setReport] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const previousMonth = format(subMonths(new Date(), 1), 'yyyy-MM');
  const previousMonthLabel = format(subMonths(new Date(), 1), 'MMMM yyyy', { locale: es }).replace(/^\w/, c => c.toUpperCase());
  const currentMonthLabel = format(new Date(), 'MMMM yyyy', { locale: es }).replace(/^\w/, c => c.toUpperCase());

  const loadReport = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const lastMonthStart = startOfMonth(subMonths(new Date(), 1));
      const lastMonthEnd = endOfMonth(subMonths(new Date(), 1));
      const periodParam = { from: lastMonthStart.toISOString(), to: lastMonthEnd.toISOString() };

      const gamAPI = getGamificationAPI();
      const [catRes, incomeRes, budgetRes, goalsRes, profileRes] = await Promise.allSettled([
        analyticsAPI.categories(periodParam),
        analyticsAPI.incomes(periodParam),
        budgetsAPI.getDashboard(),
        savingsGoalsAPI.list({ status: 'active' }),
        gamAPI.getBehaviorProfile(),
      ]);

      const financialData = {};

      if (catRes.status === 'fulfilled') {
        const cats = catRes.value.data?.data || [];
        financialData.expenses_by_category = Object.fromEntries(
          cats.map(c => [c.category_name || 'Sin categoría', c.amount || 0])
        );
        financialData.total_expenses = cats.reduce((s, c) => s + (c.amount || 0), 0);
      }

      if (incomeRes.status === 'fulfilled') {
        const d = incomeRes.value.data;
        financialData.total_income = d?.total_amount || d?.data?.total_amount || 0;
        const count = d?.count || d?.data?.count || 1;
        financialData.income_stability = count >= 3 ? 0.9 : count === 2 ? 0.65 : 0.35;
      }

      const ti = financialData.total_income || 0;
      const te = financialData.total_expenses || 0;
      financialData.savings_rate = ti > 0 ? (ti - te) / ti : 0;

      if (budgetRes.status === 'fulfilled') {
        const s = budgetRes.value.data?.summary || {};
        if (s.total_budgets > 0) {
          financialData.budgets_summary = {
            total_budgets: s.total_budgets || 0,
            total_allocated: s.total_allocated || 0,
            total_spent: s.total_spent || 0,
            on_track_count: s.on_track_count || 0,
            warning_count: s.warning_count || 0,
            exceeded_count: s.exceeded_count || 0,
            average_usage: s.average_usage || 0,
          };
        }
      }

      if (goalsRes.status === 'fulfilled') {
        const goals = goalsRes.value.data?.data?.goals || [];
        if (goals.length > 0) {
          financialData.savings_goals = goals.map(g => ({
            name: g.name,
            target_amount: g.target_amount,
            current_amount: g.current_amount,
            progress: g.progress || 0,
          }));
        }
      }

      if (profileRes.status === 'fulfilled' && profileRes.value) {
        const bp = profileRes.value;
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

      financialData.period = previousMonth;
      financialData.financial_score = 0;

      const res = await aiAPI.getMonthlyCoaching(financialData, previousMonth);
      setReport(res.report);

      // Record gamification action
      try {
        const gamAPI2 = getGamificationAPI();
        await gamAPI2.recordAction('complete_monthly_review', null, 'monthly_coaching');
      } catch (_) {}

    } catch (err) {
      console.error('Error loading monthly coaching:', err);
      setError('Error cargando el reporte mensual. Por favor intentá de nuevo.');
    } finally {
      setLoading(false);
    }
  }, [previousMonth]);

  useEffect(() => {
    loadReport();
  }, [loadReport]);

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-16 space-y-4">
        <FaSpinner className="w-10 h-10 animate-spin text-blue-500 dark:text-blue-400" />
        <p className="text-gray-600 dark:text-gray-400">Analizando tu mes de {previousMonthLabel}...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600 dark:text-gray-400 mb-4">{error}</p>
        <button
          onClick={loadReport}
          className="inline-flex items-center px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          <FaRedo className="w-4 h-4 mr-2" />
          Reintentar
        </button>
      </div>
    );
  }

  if (!report) return null;

  const sentiment = sentimentConfig[report.sentiment] || sentimentConfig.neutral;

  return (
    <div className="space-y-5">
      {/* Header con sentiment badge */}
      <div className={`rounded-2xl border p-5 ${sentiment.bg}`}>
        <div className="flex items-center justify-between mb-2">
          <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${sentiment.badge}`}>
            {sentiment.emoji} {sentiment.label}
          </span>
          <span className="text-xs text-gray-500 dark:text-gray-400">{previousMonthLabel}</span>
        </div>
        <p className="text-gray-700 dark:text-gray-300 leading-relaxed">{report.summary}</p>
      </div>

      {/* Lo que hiciste bien + Para mejorar — side by side */}
      {((report.wins && report.wins.length > 0) || (report.improvements && report.improvements.length > 0)) && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {report.wins && report.wins.length > 0 && (
            <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 p-5">
              <h3 className="text-sm font-semibold text-green-700 dark:text-green-400 mb-3 flex items-center uppercase tracking-wide">
                <span className="mr-2">🏆</span> Lo que hiciste bien
              </h3>
              <div className="space-y-2">
                {report.wins.map((win, i) => (
                  <div key={i} className="bg-green-50 dark:bg-green-900/20 rounded-xl p-3">
                    <p className="font-semibold text-gray-900 dark:text-gray-100 text-sm mb-0.5">{win.title}</p>
                    <p className="text-xs text-gray-600 dark:text-gray-400 leading-relaxed">{win.description}</p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {report.improvements && report.improvements.length > 0 && (
            <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 p-5">
              <h3 className="text-sm font-semibold text-yellow-700 dark:text-yellow-400 mb-3 flex items-center uppercase tracking-wide">
                <span className="mr-2">📈</span> Para mejorar
              </h3>
              <div className="space-y-2">
                {report.improvements.map((imp, i) => (
                  <div key={i} className="bg-yellow-50 dark:bg-yellow-900/20 rounded-xl p-3">
                    <p className="font-semibold text-gray-900 dark:text-gray-100 text-sm mb-0.5">{imp.title}</p>
                    <p className="text-xs text-gray-600 dark:text-gray-400 leading-relaxed">{imp.description}</p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Plan para el mes actual */}
      {report.actions && report.actions.length > 0 && (
        <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 p-5">
          <h3 className="text-sm font-semibold text-indigo-700 dark:text-indigo-400 mb-3 flex items-center uppercase tracking-wide">
            <span className="mr-2">🎯</span> Tu plan para {currentMonthLabel}
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            {report.actions.map((action, i) => (
              <div key={i} className="bg-indigo-50 dark:bg-indigo-900/20 rounded-xl p-4 flex flex-col">
                <p className="font-semibold text-gray-900 dark:text-gray-100 text-sm mb-1">{action.title}</p>
                <p className="text-xs text-gray-600 dark:text-gray-400 leading-relaxed flex-1">{action.detail}</p>
                {action.deep_link && (
                  <button
                    onClick={() => navigate(action.deep_link)}
                    className="mt-3 inline-flex items-center justify-center px-3 py-1.5 bg-indigo-500 dark:bg-indigo-600 text-white text-xs rounded-lg hover:bg-indigo-600 dark:hover:bg-indigo-700 transition-colors font-medium"
                  >
                    Ir <FaArrowRight className="w-3 h-3 ml-1" />
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Behavior note */}
      {report.behavior_note && (
        <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 px-5 py-4 flex items-start space-x-3">
          <span className="text-xl mt-0.5">📊</span>
          <div>
            <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Tu patrón financiero</p>
            <p className="text-sm text-gray-700 dark:text-gray-300 leading-relaxed">{report.behavior_note}</p>
          </div>
        </div>
      )}

      <p className="text-xs text-gray-400 dark:text-gray-500 text-center">
        Reporte generado con tus datos reales de {previousMonthLabel} · Actualización mensual
      </p>
    </div>
  );
};

export default MonthlyCoachingTab;
