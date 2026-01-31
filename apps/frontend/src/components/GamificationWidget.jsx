/**
 * 🎮 GAMIFICATION WIDGET COMPONENT
 * 
 * Widget compacto para mostrar el estado de gamificación del usuario:
 * - Nivel actual y nombre
 * - XP total de forma simple
 * - Click para ir a página de achievements
 */

import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { FaTrophy, FaStar, FaMedal } from 'react-icons/fa';
import { useGamification } from '../contexts/GamificationContext';

const GamificationWidget = () => {
  const navigate = useNavigate();
  const { userProfile, loading, error, getLevelInfo, refreshTrigger } = useGamification();
  const [localXP, setLocalXP] = useState(0);
  const [localLevel, setLocalLevel] = useState(0);

  // Sincronizar estado local con userProfile (sin dependencias circulares)
  useEffect(() => {
    if (userProfile) {
      const newXP = userProfile.total_xp || 0;
      const newLevel = userProfile.current_level || 0;
      
      // Solo actualizar si los valores realmente cambiaron
      if (newXP !== localXP || newLevel !== localLevel) {
        console.log(`🎮 [GamificationWidget] Sincronizando: XP=${localXP}→${newXP}, Nivel=${localLevel}→${newLevel} (trigger: ${refreshTrigger})`);
        setLocalXP(newXP);
        setLocalLevel(newLevel);
      }
    }
  }, [userProfile, refreshTrigger, localXP, localLevel]);

  const handleClick = () => {
    navigate('/achievements');
  };

  if (loading) {
    return (
      <div className="flex items-center space-x-2 px-2 py-1.5 bg-gray-100 dark:bg-gray-700 rounded-lg animate-pulse">
        <div className="w-6 h-6 bg-gray-300 dark:bg-gray-600 rounded-full"></div>
        <div className="w-20 h-3 bg-gray-300 dark:bg-gray-600 rounded"></div>
      </div>
    );
  }

  if (error || !userProfile) {
    return (
      <div className="flex items-center space-x-2 px-2 py-1.5 bg-gray-100 dark:bg-gray-700 rounded-lg opacity-50">
        <FaTrophy className="w-3 h-3 text-gray-400" />
        <span className="text-xs text-gray-500">No disponible</span>
      </div>
    );
  }

  // Usar valores locales que se actualizan con useEffect para mejor reactividad
  const totalXP = localXP;
  const currentLevel = localLevel;
  const levelInfo = getLevelInfo(currentLevel);

  return (
    <div 
      onClick={handleClick}
      className="flex items-center space-x-2 px-2 py-1.5 bg-gradient-to-r from-purple-50 to-blue-50 dark:from-purple-900/20 dark:to-blue-900/20 rounded-lg cursor-pointer hover:shadow-sm transition-all duration-200 group border border-purple-100 dark:border-purple-800"
    >
      {/* Icono de nivel con badge */}
      <div className="relative">
        <FaStar 
          className="w-5 h-5"
          style={{ color: levelInfo.color }}
        />
        <div className="absolute -top-1 -right-1 w-3 h-3 bg-white dark:bg-gray-800 rounded-full flex items-center justify-center border border-gray-200 dark:border-gray-600">
          <span className="text-xs font-bold text-gray-700 dark:text-gray-300" style={{ fontSize: '8px' }}>
            {currentLevel}
          </span>
        </div>
      </div>

      {/* Información del nivel */}
      <div className="flex flex-col min-w-0">
        <div className="flex items-center space-x-1">
          <span className="text-xs font-medium text-gray-900 dark:text-gray-100 truncate">
            Nivel {currentLevel}
          </span>
          <span className="text-xs text-gray-500 dark:text-gray-400">•</span>
          <span className="text-xs text-gray-600 dark:text-gray-400 truncate">
            {levelInfo.name}
          </span>
        </div>
        <span className="text-xs text-gray-500 dark:text-gray-400">
          {totalXP.toLocaleString()} XP
        </span>
      </div>
    </div>
  );
};

export default GamificationWidget; 