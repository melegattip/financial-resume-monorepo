import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { FaSpinner, FaRedo } from 'react-icons/fa';
import { aiAPI, analyticsAPI, budgetsAPI, savingsGoalsAPI } from '../../../services/api';
import { getGamificationAPI } from '../../../services/gamificationAPI';

const topicConfig = {
  emergencia:  { icon: '🛡️', color: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300' },
  presupuesto: { icon: '📊', color: 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300' },
  deuda:       { icon: '⚠️', color: 'bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-300' },
  ahorro:      { icon: '💰', color: 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300' },
  'inversión': { icon: '📈', color: 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300' },
  impuestos:   { icon: '📋', color: 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300' },
};

const difficultyConfig = {
  'básico':      'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300',
  'intermedio':  'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-300',
  'avanzado':    'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300',
};

const EducationTab = () => {
  const navigate = useNavigate();
  const [cards, setCards] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [clickedCards, setClickedCards] = useState(new Set());

  const loadCards = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const gamAPI = getGamificationAPI();
      const [catRes, incomeRes, budgetRes, goalsRes, profileRes] = await Promise.allSettled([
        analyticsAPI.categories(),
        analyticsAPI.incomes(),
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
      }

      const ti = financialData.total_income || 0;
      const te = financialData.total_expenses || 0;
      financialData.savings_rate = ti > 0 ? (ti - te) / ti : 0;

      if (budgetRes.status === 'fulfilled') {
        const s = budgetRes.value.data?.summary || {};
        if (s.total_budgets > 0) {
          financialData.budgets_summary = { total_budgets: s.total_budgets || 0 };
        }
      }

      if (goalsRes.status === 'fulfilled') {
        const goals = goalsRes.value.data?.data?.goals || [];
        if (goals.length > 0) {
          financialData.savings_goals = goals.map(g => ({ name: g.name, progress: g.progress || 0 }));
        }
      }

      if (profileRes.status === 'fulfilled' && profileRes.value) {
        const bp = profileRes.value;
        financialData.behavior_profile = {
          current_level: bp.current_level,
          level_name: bp.level_name,
          discipline_score: bp.discipline_score,
          ai_recommendations_applied: bp.ai_recommendations_applied,
        };
      }

      financialData.financial_score = 0;

      const res = await aiAPI.getEducationCards(financialData);
      setCards(res.cards || []);
    } catch (err) {
      console.error('Error loading education cards:', err);
      setError('Error cargando las tarjetas educativas. Por favor intentá de nuevo.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadCards();
  }, [loadCards]);

  const handleCTA = async (card, idx) => {
    if (!clickedCards.has(idx)) {
      setClickedCards(prev => new Set([...prev, idx]));
      try {
        const gamAPI = getGamificationAPI();
        await gamAPI.recordAction('read_education_card', null, `education_card_${idx}`);
      } catch (_) {}
    }
    if (card.deep_link) {
      navigate(card.deep_link);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-16 space-y-4">
        <FaSpinner className="w-10 h-10 animate-spin text-blue-500 dark:text-blue-400" />
        <p className="text-gray-600 dark:text-gray-400">Preparando tu contenido educativo personalizado...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600 dark:text-gray-400 mb-4">{error}</p>
        <button
          onClick={loadCards}
          className="inline-flex items-center px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          <FaRedo className="w-4 h-4 mr-2" />
          Reintentar
        </button>
      </div>
    );
  }

  if (cards.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600 dark:text-gray-400">No hay tarjetas educativas disponibles.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <p className="text-sm text-gray-500 dark:text-gray-400 text-center">
        Contenido personalizado a tu perfil financiero · Se actualiza semanalmente
      </p>
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        {cards.map((card, idx) => {
          const topic = topicConfig[card.topic] || { icon: '📚', color: 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300' };
          const diffClass = difficultyConfig[card.difficulty] || difficultyConfig['básico'];
          return (
            <div key={idx} className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 p-5 flex flex-col space-y-4">
              {/* Topic + difficulty badges */}
              <div className="flex items-center justify-between">
                <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${topic.color}`}>
                  {topic.icon} {card.topic}
                </span>
                <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${diffClass}`}>
                  {card.difficulty}
                </span>
              </div>

              {/* Title */}
              <h3 className="text-base font-semibold text-gray-900 dark:text-gray-100 leading-snug">{card.title}</h3>

              {/* Summary */}
              <p className="text-gray-600 dark:text-gray-400 text-sm leading-relaxed flex-1">{card.summary}</p>

              {/* Key concept callout */}
              {card.key_concept && (
                <div className="border-l-4 border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20 px-3 py-2 rounded-r-lg">
                  <p className="text-xs font-semibold text-indigo-700 dark:text-indigo-300 uppercase tracking-wide mb-0.5">Concepto clave</p>
                  <p className="text-sm text-indigo-800 dark:text-indigo-200">{card.key_concept}</p>
                </div>
              )}

              {/* CTA button */}
              <button
                onClick={() => handleCTA(card, idx)}
                className="w-full px-4 py-2.5 bg-indigo-500 dark:bg-indigo-600 text-white rounded-xl hover:bg-indigo-600 dark:hover:bg-indigo-700 transition-colors font-medium text-sm mt-auto"
              >
                {card.cta || 'Explorar'} {clickedCards.has(idx) ? '✓' : '→'}
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default EducationTab;
