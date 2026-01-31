/**
 * 🏆 ACHIEVEMENTS PAGE
 * 
 * Página principal de gamificación que muestra:
 * - Progreso del usuario y nivel actual
 * - Logros completados y por completar
 * - Estadísticas detalladas
 * - Historial de acciones recientes
 */

import React, { useState, useEffect } from 'react';
import { FaTrophy, FaStar, FaBolt, FaChartLine, FaCalendarAlt, FaEye, FaCheckCircle, FaTimesCircle, FaCrown, FaRocket, FaFire } from 'react-icons/fa';
import { getGamificationAPI } from '../services/gamificationAPI';
import { useGamificationNotifications } from '../components/GamificationNotification';

const Achievements = () => {
  const [loading, setLoading] = useState(true);
  const [gamificationData, setGamificationData] = useState(null);
  const [achievements, setAchievements] = useState([]);
  const [stats, setStats] = useState(null);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState('overview');

  const { GamificationNotification } = useGamificationNotifications();

  useEffect(() => {
    loadGamificationData();
  }, []);

  const loadGamificationData = async () => {
    try {
      setLoading(true);
      const api = getGamificationAPI();
      
      const [profile, achievementsData, statsData] = await Promise.all([
        api.getUserProfile(),
        api.getUserAchievements(),
        api.getUserStats()
      ]);

      setGamificationData(profile);
      setAchievements(achievementsData);
      setStats(statsData);
      setError(null);
    } catch (err) {
      console.error('Error loading gamification data:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const getLevelInfo = (level) => {
    const levels = [
      { level: 0, name: "Financial Newbie", color: "#94A3B8", emoji: "🌱", description: "Comenzando tu viaje financiero" },
      { level: 1, name: "Money Aware", color: "#60A5FA", emoji: "👀", description: "Consciente de tus finanzas" },
      { level: 2, name: "Budget Tracker", color: "#34D399", emoji: "📊", description: "Seguimiento activo del presupuesto" },
      { level: 3, name: "Savings Starter", color: "#FBBF24", emoji: "💰", description: "Iniciando el camino del ahorro" },
      { level: 4, name: "Financial Explorer", color: "#F472B6", emoji: "🧭", description: "Explorando nuevas estrategias" },
      { level: 5, name: "Money Manager", color: "#A78BFA", emoji: "💼", description: "Gestión avanzada del dinero" },
      { level: 6, name: "Investment Learner", color: "#FB7185", emoji: "📈", description: "Aprendiendo sobre inversiones" },
      { level: 7, name: "Financial Guru", color: "#10B981", emoji: "🧠", description: "Maestro de las finanzas" },
      { level: 8, name: "Money Master", color: "#8B5CF6", emoji: "👑", description: "Dominio total del dinero" },
      { level: 9, name: "Financial Magnate", color: "#EF4444", emoji: "💎", description: "Magnate financiero supremo" }
    ];
    
    return levels[level] || levels[0];
  };

  const getAchievementIcon = (type) => {
    switch (type) {
      case 'ai_partner': return '🤖';
      case 'action_taker': return '🎯';
      case 'data_explorer': return '📊';
      case 'quick_learner': return '⚡';
      case 'insight_master': return '💡';
      case 'streak_keeper': return '🔥';
      default: return '🏆';
    }
  };

  const getProgressBarColor = (progress) => {
    if (progress >= 100) return 'from-green-500 to-emerald-500';
    if (progress >= 75) return 'from-blue-500 to-cyan-500';
    if (progress >= 50) return 'from-yellow-500 to-orange-500';
    return 'from-gray-400 to-gray-500';
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="bg-white dark:bg-gray-800 rounded-lg p-6">
          <div className="animate-pulse space-y-4">
            <div className="h-8 bg-gray-300 dark:bg-gray-600 rounded w-1/3"></div>
            <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-2/3"></div>
            <div className="h-20 bg-gray-300 dark:bg-gray-600 rounded"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <div className="flex items-center space-x-3">
            <FaTimesCircle className="w-6 h-6 text-red-600 dark:text-red-400" />
            <div>
              <h3 className="text-lg font-semibold text-red-900 dark:text-red-200">
                Error al cargar gamificación
              </h3>
              <p className="text-red-700 dark:text-red-300">
                {error}
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  const currentLevel = getLevelInfo(gamificationData?.current_level || 0);
  const nextLevel = getLevelInfo((gamificationData?.current_level || 0) + 1);
  const completedAchievements = achievements?.filter(a => a.completed) || [];
  const pendingAchievements = achievements?.filter(a => !a.completed) || [];

  return (
    <div className="space-y-6">
      <GamificationNotification />
      
      {/* Header con título */}
      <div className="bg-gradient-to-r from-purple-600 to-blue-600 rounded-lg p-6 text-white">
        <div className="flex items-center space-x-4">
          <div className="p-3 bg-white/20 rounded-full">
            <FaTrophy className="w-8 h-8" />
          </div>
          <div>
            <h1 className="text-2xl font-bold">Logros y Progreso</h1>
            <p className="text-purple-100">Tu evolución en el manejo financiero inteligente</p>
          </div>
        </div>
      </div>

      {/* Tabs de navegación */}
      <div className="bg-white dark:bg-gray-800 rounded-lg overflow-hidden">
        <div className="border-b border-gray-200 dark:border-gray-700">
          <nav className="flex space-x-8 px-6">
            {[
              { id: 'overview', label: 'Resumen', icon: FaChartLine },
              { id: 'achievements', label: 'Logros', icon: FaTrophy },
              { id: 'stats', label: 'Estadísticas', icon: FaStar }
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center space-x-2 py-4 border-b-2 transition-colors ${
                  activeTab === tab.id
                    ? 'border-purple-500 text-purple-600 dark:text-purple-400'
                    : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
                }`}
              >
                <tab.icon className="w-4 h-4" />
                <span className="font-medium">{tab.label}</span>
              </button>
            ))}
          </nav>
        </div>

        <div className="p-6">
          {/* Tab: Resumen */}
          {activeTab === 'overview' && (
            <div className="space-y-6">
              {/* Progreso del nivel */}
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="bg-gradient-to-br from-purple-50 to-blue-50 dark:from-purple-900/20 dark:to-blue-900/20 rounded-lg p-6 border border-purple-100 dark:border-purple-800">
                  <div className="flex items-center space-x-4 mb-4">
                    <div 
                      className="w-16 h-16 rounded-full flex items-center justify-center text-white text-2xl font-bold shadow-lg"
                      style={{ backgroundColor: currentLevel.color }}
                    >
                      {currentLevel.emoji}
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                        {currentLevel.name}
                      </h3>
                      <p className="text-sm text-gray-600 dark:text-gray-300">
                        {currentLevel.description}
                      </p>
                    </div>
                  </div>
                  
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                        Progreso al siguiente nivel
                      </span>
                      <span className="text-sm text-gray-600 dark:text-gray-400">
                        {stats?.progress_percent || 0}%
                      </span>
                    </div>
                    <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                      <div 
                        className="h-full bg-gradient-to-r from-purple-500 to-blue-500 rounded-full transition-all duration-500"
                        style={{ width: `${stats?.progress_percent || 0}%` }}
                      />
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <div className="flex items-center space-x-1">
                        <FaBolt className="w-3 h-3 text-yellow-500" />
                        <span className="text-gray-600 dark:text-gray-300">
                          {(stats?.total_xp || 0).toLocaleString()} XP
                        </span>
                      </div>
                      {stats?.xp_to_next_level > 0 && (
                        <span className="text-gray-500 dark:text-gray-400">
                          +{stats.xp_to_next_level} para {nextLevel.name}
                        </span>
                      )}
                    </div>
                  </div>
                </div>

                {/* Resumen de logros */}
                <div className="bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 rounded-lg p-6 border border-green-100 dark:border-green-800">
                  <div className="flex items-center space-x-3 mb-4">
                    <FaTrophy className="w-6 h-6 text-green-600 dark:text-green-400" />
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                      Logros Desbloqueados
                    </h3>
                  </div>
                  
                  <div className="grid grid-cols-2 gap-4">
                    <div className="text-center">
                      <div className="text-3xl font-bold text-green-600 dark:text-green-400">
                        {completedAchievements.length}
                      </div>
                      <div className="text-sm text-gray-600 dark:text-gray-300">
                        Completados
                      </div>
                    </div>
                    <div className="text-center">
                      <div className="text-3xl font-bold text-gray-400 dark:text-gray-500">
                        {pendingAchievements.length}
                      </div>
                      <div className="text-sm text-gray-600 dark:text-gray-300">
                        Pendientes
                      </div>
                    </div>
                  </div>

                  {completedAchievements.length > 0 && (
                    <div className="mt-4 pt-4 border-t border-green-200 dark:border-green-700">
                      <div className="flex -space-x-2">
                        {completedAchievements.slice(0, 5).map((achievement, index) => (
                          <div
                            key={achievement.id}
                            className="w-8 h-8 bg-green-100 dark:bg-green-800 rounded-full flex items-center justify-center border-2 border-white dark:border-gray-800 text-sm"
                            title={achievement.name}
                          >
                            {getAchievementIcon(achievement.type)}
                          </div>
                        ))}
                        {completedAchievements.length > 5 && (
                          <div className="w-8 h-8 bg-gray-100 dark:bg-gray-700 rounded-full flex items-center justify-center border-2 border-white dark:border-gray-800 text-xs font-medium text-gray-600 dark:text-gray-300">
                            +{completedAchievements.length - 5}
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              </div>

              {/* Próximos logros */}
              <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center space-x-2">
                  <FaRocket className="w-5 h-5 text-blue-600 dark:text-blue-400" />
                  <span>Próximos Logros</span>
                </h3>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {pendingAchievements.slice(0, 4).map((achievement) => {
                    const progress = Math.min((achievement.progress / achievement.target) * 100, 100);
                    return (
                      <div 
                        key={achievement.id}
                        className="border border-gray-200 dark:border-gray-600 rounded-lg p-4 hover:shadow-md transition-shadow"
                      >
                        <div className="flex items-center space-x-3 mb-3">
                          <div className="text-2xl">
                            {getAchievementIcon(achievement.type)}
                          </div>
                          <div className="flex-1 min-w-0">
                            <h4 className="font-medium text-gray-900 dark:text-white truncate">
                              {achievement.name}
                            </h4>
                            <p className="text-sm text-gray-500 dark:text-gray-400 truncate">
                              {achievement.description}
                            </p>
                          </div>
                        </div>
                        
                        <div className="space-y-2">
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-gray-600 dark:text-gray-300">
                              Progreso
                            </span>
                            <span className="font-medium">
                              {achievement.progress}/{achievement.target}
                            </span>
                          </div>
                          <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                            <div 
                              className={`h-full bg-gradient-to-r ${getProgressBarColor(progress)} transition-all duration-300`}
                              style={{ width: `${progress}%` }}
                            />
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>
          )}

          {/* Tab: Logros */}
          {activeTab === 'achievements' && (
            <div className="space-y-6">
              {/* Logros completados */}
              {completedAchievements.length > 0 && (
                <div>
                  <div className="flex items-center space-x-2 mb-4">
                    <FaCheckCircle className="w-5 h-5 text-green-600 dark:text-green-400" />
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                      Logros Completados ({completedAchievements.length})
                    </h3>
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {completedAchievements.map((achievement) => (
                      <div 
                        key={achievement.id}
                        className="bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 border border-green-200 dark:border-green-700 rounded-lg p-4"
                      >
                        <div className="flex items-center space-x-3 mb-3">
                          <div className="text-3xl">
                            {getAchievementIcon(achievement.type)}
                          </div>
                          <div className="flex-1 min-w-0">
                            <h4 className="font-semibold text-green-900 dark:text-green-200">
                              {achievement.name}
                            </h4>
                            <p className="text-sm text-green-700 dark:text-green-300">
                              {achievement.description}
                            </p>
                          </div>
                        </div>
                        
                        <div className="flex items-center justify-between">
                          <div className="flex items-center space-x-1">
                            <FaBolt className="w-3 h-3 text-yellow-500" />
                            <span className="text-sm font-medium text-green-800 dark:text-green-200">
                              +{achievement.points} XP
                            </span>
                          </div>
                          {achievement.unlocked_at && (
                            <div className="flex items-center space-x-1 text-xs text-green-600 dark:text-green-400">
                              <FaCalendarAlt className="w-3 h-3" />
                              <span>
                                {new Date(achievement.unlocked_at).toLocaleDateString()}
                              </span>
                            </div>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Logros pendientes */}
              {pendingAchievements.length > 0 && (
                <div>
                  <div className="flex items-center space-x-2 mb-4">
                    <FaCrown className="w-5 h-5 text-blue-600 dark:text-blue-400" />
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                      Logros Pendientes ({pendingAchievements.length})
                    </h3>
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {pendingAchievements.map((achievement) => {
                      const progress = Math.min((achievement.progress / achievement.target) * 100, 100);
                      return (
                        <div 
                          key={achievement.id}
                          className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 rounded-lg p-4 hover:shadow-md transition-shadow"
                        >
                          <div className="flex items-center space-x-3 mb-3">
                            <div className="text-3xl opacity-60">
                              {getAchievementIcon(achievement.type)}
                            </div>
                            <div className="flex-1 min-w-0">
                              <h4 className="font-semibold text-gray-900 dark:text-white">
                                {achievement.name}
                              </h4>
                              <p className="text-sm text-gray-500 dark:text-gray-400">
                                {achievement.description}
                              </p>
                            </div>
                          </div>
                          
                          <div className="space-y-2">
                            <div className="flex items-center justify-between text-sm">
                              <span className="text-gray-600 dark:text-gray-300">
                                Progreso
                              </span>
                              <span className="font-medium">
                                {achievement.progress}/{achievement.target}
                              </span>
                            </div>
                            <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                              <div 
                                className={`h-full bg-gradient-to-r ${getProgressBarColor(progress)} transition-all duration-300`}
                                style={{ width: `${progress}%` }}
                              />
                            </div>
                            <div className="flex items-center justify-between">
                              <div className="flex items-center space-x-1">
                                <FaBolt className="w-3 h-3 text-yellow-500" />
                                <span className="text-sm text-gray-600 dark:text-gray-300">
                                  +{achievement.points} XP
                                </span>
                              </div>
                              <span className="text-xs text-gray-500 dark:text-gray-400">
                                {Math.round(progress)}% completado
                              </span>
                            </div>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Tab: Estadísticas */}
          {activeTab === 'stats' && stats && (
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {[
                  {
                    label: 'XP Total',
                    value: stats.total_xp.toLocaleString(),
                    icon: FaBolt,
                    color: 'text-yellow-600 dark:text-yellow-400',
                    bg: 'bg-yellow-50 dark:bg-yellow-900/20'
                  },
                  {
                    label: 'Nivel Actual',
                    value: stats.current_level,
                    icon: FaStar,
                    color: 'text-purple-600 dark:text-purple-400',
                    bg: 'bg-purple-50 dark:bg-purple-900/20'
                  },
                  {
                    label: 'Logros Completados',
                    value: `${stats.completed_achievements}/${stats.total_achievements}`,
                    icon: FaTrophy,
                    color: 'text-green-600 dark:text-green-400',
                    bg: 'bg-green-50 dark:bg-green-900/20'
                  },
                  {
                    label: 'Racha Actual',
                    value: `${stats.current_streak} días`,
                    icon: FaFire,
                    color: 'text-red-600 dark:text-red-400',
                    bg: 'bg-red-50 dark:bg-red-900/20'
                  }
                ].map((stat, index) => (
                  <div key={index} className={`${stat.bg} border border-gray-200 dark:border-gray-600 rounded-lg p-4`}>
                    <div className="flex items-center space-x-3">
                      <div className={`p-2 ${stat.bg} rounded-lg`}>
                        <stat.icon className={`w-5 h-5 ${stat.color}`} />
                      </div>
                      <div>
                        <div className="text-2xl font-bold text-gray-900 dark:text-white">
                          {stat.value}
                        </div>
                        <div className="text-sm text-gray-600 dark:text-gray-300">
                          {stat.label}
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* Información adicional */}
              <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-6">
                <h4 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                  Detalles del Progreso
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="text-gray-600 dark:text-gray-300">Última actividad:</span>
                    <div className="font-medium text-gray-900 dark:text-white">
                      {stats.last_activity ? new Date(stats.last_activity).toLocaleString() : 'Sin actividad'}
                    </div>
                  </div>
                  <div>
                    <span className="text-gray-600 dark:text-gray-300">XP para siguiente nivel:</span>
                    <div className="font-medium text-gray-900 dark:text-white">
                      {stats.xp_to_next_level > 0 ? stats.xp_to_next_level : 'Nivel máximo alcanzado'}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Achievements; 