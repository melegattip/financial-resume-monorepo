import React, { useState, useEffect } from 'react';
import { useGamification } from '../contexts/GamificationContext';

/**
 * 🏆 DAILY CHALLENGES COMPONENT
 * 
 * Muestra los challenges diarios disponibles y el progreso del usuario
 */
const DailyChallenges = () => {
  const { 
    getDailyChallenges, 
    getWeeklyChallenges,
    processChallengeProgress,
    userProfile,
    loading 
  } = useGamification();

  const [dailyChallenges, setDailyChallenges] = useState([]);
  const [weeklyChallenges, setWeeklyChallenges] = useState([]);
  const [challengesLoading, setChallengesLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadChallenges();
  }, [userProfile]);

  const loadChallenges = async () => {
    if (!userProfile || loading) return;

    try {
      setChallengesLoading(true);
      setError(null);

      const [dailyData, weeklyData] = await Promise.all([
        getDailyChallenges(),
        getWeeklyChallenges()
      ]);

      setDailyChallenges(dailyData || []);
      setWeeklyChallenges(weeklyData || []);
    } catch (err) {
      console.error('Error loading challenges:', err);
      setError('Error cargando challenges. Inténtalo de nuevo más tarde.');
    } finally {
      setChallengesLoading(false);
    }
  };

  const handleChallengeAction = async (actionType, description = '') => {
    try {
      await processChallengeProgress(actionType, '', '', 0, description);
      // Recargar challenges para ver progreso actualizado
      await loadChallenges();
    } catch (err) {
      console.error('Error processing challenge action:', err);
    }
  };

  if (loading || challengesLoading) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
        <div className="animate-pulse">
          <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-1/3 mb-4"></div>
          <div className="space-y-3">
            {[1, 2, 3].map(i => (
              <div key={i} className="h-16 bg-gray-200 dark:bg-gray-700 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
        <div className="text-center">
          <div className="text-red-500 mb-2">⚠️</div>
          <p className="text-gray-600 dark:text-gray-400 mb-4">{error}</p>
          <button 
            onClick={loadChallenges}
            className="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 transition-colors"
          >
            Reintentar
          </button>
        </div>
      </div>
    );
  }

  const ChallengeCard = ({ challenge, type }) => (
    <div className="bg-gradient-to-r from-blue-50 to-purple-50 dark:from-gray-700 dark:to-gray-600 rounded-lg p-4 border border-gray-200 dark:border-gray-600">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="text-2xl">{challenge.icon || '🏆'}</div>
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-white">
              {challenge.name || challenge.challenge_key}
            </h4>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {challenge.description}
            </p>
          </div>
        </div>
        <div className="text-right">
          <div className="text-lg font-bold text-blue-600 dark:text-blue-400">
            +{challenge.xp_reward || 0} XP
          </div>
          <div className="text-xs text-gray-500">
            {challenge.progress || 0} / {challenge.target || 1}
          </div>
        </div>
      </div>
      
      {/* Progress Bar */}
      <div className="mt-3">
        <div className="flex justify-between text-xs text-gray-600 dark:text-gray-400 mb-1">
          <span>Progreso</span>
          <span>{Math.round(((challenge.progress || 0) / (challenge.target || 1)) * 100)}%</span>
        </div>
        <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
          <div 
            className="bg-blue-500 h-2 rounded-full transition-all duration-300"
            style={{ 
              width: `${Math.min(((challenge.progress || 0) / (challenge.target || 1)) * 100, 100)}%` 
            }}
          ></div>
        </div>
      </div>

      {challenge.completed && (
        <div className="mt-2 flex items-center text-green-600 dark:text-green-400 text-sm">
          <span className="mr-1">✅</span>
          Completado
        </div>
      )}
    </div>
  );

  return (
    <div className="space-y-6">
      {/* Daily Challenges */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            🌅 Challenges Diarios
          </h3>
          <button 
            onClick={loadChallenges}
            className="text-blue-500 hover:text-blue-600 text-sm"
          >
            🔄 Actualizar
          </button>
        </div>
        
        {dailyChallenges.length > 0 ? (
          <div className="space-y-3">
            {dailyChallenges.map((challenge, index) => (
              <ChallengeCard 
                key={challenge.id || index} 
                challenge={challenge} 
                type="daily" 
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <div className="text-gray-400 text-4xl mb-2">🎯</div>
            <p className="text-gray-600 dark:text-gray-400">
              No hay challenges diarios disponibles
            </p>
          </div>
        )}
      </div>

      {/* Weekly Challenges */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          📅 Challenges Semanales
        </h3>
        
        {weeklyChallenges.length > 0 ? (
          <div className="space-y-3">
            {weeklyChallenges.map((challenge, index) => (
              <ChallengeCard 
                key={challenge.id || index} 
                challenge={challenge} 
                type="weekly" 
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <div className="text-gray-400 text-4xl mb-2">📊</div>
            <p className="text-gray-600 dark:text-gray-400">
              No hay challenges semanales disponibles
            </p>
          </div>
        )}
      </div>

      {/* Test Actions */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          🧪 Acciones de Prueba
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
          <button 
            onClick={() => handleChallengeAction('create_expense', 'Test expense')}
            className="bg-green-500 text-white px-3 py-2 rounded text-sm hover:bg-green-600 transition-colors"
          >
            Crear Gasto
          </button>
          <button 
            onClick={() => handleChallengeAction('create_category', 'Test category')}
            className="bg-blue-500 text-white px-3 py-2 rounded text-sm hover:bg-blue-600 transition-colors"
          >
            Crear Categoría
          </button>
          <button 
            onClick={() => handleChallengeAction('view_dashboard', 'Dashboard visit')}
            className="bg-purple-500 text-white px-3 py-2 rounded text-sm hover:bg-purple-600 transition-colors"
          >
            Ver Dashboard
          </button>
          <button 
            onClick={() => handleChallengeAction('daily_login', 'Daily login')}
            className="bg-orange-500 text-white px-3 py-2 rounded text-sm hover:bg-orange-600 transition-colors"
          >
            Login Diario
          </button>
          <button 
            onClick={() => handleChallengeAction('view_analytics', 'Analytics view')}
            className="bg-indigo-500 text-white px-3 py-2 rounded text-sm hover:bg-indigo-600 transition-colors"
          >
            Ver Analytics
          </button>
          <button 
            onClick={() => handleChallengeAction('assign_category', 'Assign category')}
            className="bg-teal-500 text-white px-3 py-2 rounded text-sm hover:bg-teal-600 transition-colors"
          >
            Asignar Categoría
          </button>
        </div>
      </div>
    </div>
  );
};

export default DailyChallenges; 