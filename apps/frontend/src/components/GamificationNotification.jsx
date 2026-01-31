/**
 * 🎮 GAMIFICATION NOTIFICATION COMPONENT
 * 
 * Componente para mostrar notificaciones de gamificación:
 * - XP ganado
 * - Level ups
 * - Nuevos achievements
 * - Progreso general
 */

import React, { useState, useEffect } from 'react';
import { FaTrophy, FaStar, FaBullseye, FaBolt, FaTimes } from 'react-icons/fa';
import { createPortal } from 'react-dom';

const GamificationNotification = ({ 
  isVisible, 
  onClose, 
  type = 'xp', 
  points = 0, 
  title = '', 
  description = '',
  achievement = null 
}) => {
  const [shouldShow, setShouldShow] = useState(false);

  useEffect(() => {
    if (isVisible) {
      setShouldShow(true);
      // Auto-cerrar después de 4 segundos
      const timer = setTimeout(() => {
        handleClose();
      }, 4000);
      return () => clearTimeout(timer);
    }
  }, [isVisible]);

  const handleClose = () => {
    setShouldShow(false);
    setTimeout(() => {
      onClose();
    }, 300);
  };

  if (!isVisible) return null;

  const getNotificationConfig = () => {
    switch (type) {
      case 'achievement':
        return {
          icon: FaTrophy,
          bgColor: 'bg-gradient-to-r from-yellow-400 to-orange-500',
          iconColor: 'text-white',
          title: achievement?.name || 'Logro Desbloqueado',
          description: achievement?.description || 'Has completado un nuevo logro'
        };
      case 'level_up':
        return {
          icon: FaStar,
          bgColor: 'bg-gradient-to-r from-purple-500 to-pink-500',
          iconColor: 'text-white',
          title: title || 'Nivel Subido',
          description: description || `Felicidades, has alcanzado un nuevo nivel`
        };
      case 'xp':
      default:
        return {
          icon: FaBolt,
          bgColor: 'bg-gradient-to-r from-blue-500 to-green-500',
          iconColor: 'text-white',
          title: title || 'XP Ganado',
          description: `Has ganado ${points} puntos de experiencia`
        };
    }
  };

  const config = getNotificationConfig();
  const Icon = config.icon;

  const notificationContent = (
    <div className={`fixed top-4 right-4 z-50 transform transition-all duration-300 ${
      shouldShow ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'
    }`}>
      <div className={`${config.bgColor} rounded-lg shadow-lg p-4 max-w-sm min-w-[300px] text-white`}>
        <div className="flex items-start space-x-3">
          <div className="flex-shrink-0">
            <div className="p-2 bg-white/20 rounded-full backdrop-blur-sm">
              <Icon className={`w-6 h-6 ${config.iconColor}`} />
            </div>
          </div>
          
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm font-semibold text-white">
                  {config.title}
                </p>
                <p className="text-xs text-white/90 mt-1">
                  {config.description}
                </p>
                {points > 0 && (
                  <div className="flex items-center space-x-1 mt-2">
                    <FaBolt className="w-4 h-4 text-yellow-200" />
                    <span className="text-sm font-medium text-yellow-200">
                      +{points} XP
                    </span>
                  </div>
                )}
              </div>
              
              <button
                onClick={handleClose}
                className="flex-shrink-0 p-1 hover:bg-white/20 rounded-full transition-colors"
              >
                <FaTimes className="w-4 h-4 text-white/80" />
              </button>
            </div>
          </div>
        </div>
        
        {/* Barra de progreso animada */}
        <div className="mt-3 h-1 bg-white/20 rounded-full overflow-hidden">
          <div 
            className="h-full bg-white/40 rounded-full transition-all duration-4000 ease-out"
            style={{ 
              width: shouldShow ? '100%' : '0%',
              transitionDuration: '4000ms'
            }}
          />
        </div>
      </div>
    </div>
  );

  return createPortal(notificationContent, document.body);
};

// Hook para manejar notificaciones de gamificación
export const useGamificationNotifications = () => {
  const [notification, setNotification] = useState(null);

  const showNotification = (notificationData) => {
    console.log(`📢 [GamificationNotification] setNotification called with:`, notificationData);
    setNotification(notificationData);
  };

  const hideNotification = () => {
    setNotification(null);
  };

  // Métodos de conveniencia
  const showXPGained = (xp, message = '') => {
    console.log(`🎮 [GamificationNotification] showXPGained llamado con:`, { xp, message });
    
    showNotification({
      type: 'xp',
      points: xp,
      title: message || `XP Ganado`,
      description: message || `Has ganado ${xp} puntos de experiencia`
    });
  };

  const showLevelUp = (newLevel, levelName, xpProgress = null) => {
    showNotification({
      type: 'level_up',
      title: levelName,
      description: `Felicidades, has alcanzado un nuevo nivel`
    });
  };

  const showAchievementUnlocked = (achievementName, description = '') => {
    showNotification({
      type: 'achievement',
      achievement: {
        name: achievementName,
        description: description || `Has desbloqueado: ${achievementName}`
      }
    });
  };

  return {
    notification,
    showNotification,
    hideNotification,
    showXPGained,
    showLevelUp,
    showAchievementUnlocked,
    GamificationNotification: (props) => (
      <GamificationNotification
        isVisible={!!notification}
        onClose={hideNotification}
        {...notification}
        {...props}
      />
    )
  };
};

export default GamificationNotification; 