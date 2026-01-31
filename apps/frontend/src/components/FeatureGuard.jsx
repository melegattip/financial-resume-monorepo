import React from 'react';
import { useGamification } from '../contexts/GamificationContext';
import LockedFeaturePreview from './LockedFeaturePreview';

/**
 * 🔒 FEATURE GUARD
 * 
 * Controla el acceso a features basado en el nivel de gamificación del usuario
 * 
 * @param {string} feature - Clave de la feature (SAVINGS_GOALS, BUDGETS, AI_INSIGHTS)
 * @param {string} mode - Modo de protección: 'full', 'preview', 'block'
 * @param {React.Component} children - Componente a proteger
 * @param {React.Component} fallback - Componente alternativo (opcional)
 * @param {boolean} showPreview - Si mostrar preview cuando está bloqueada
 */
const FeatureGuard = ({ 
  feature, 
  mode = 'preview',
  children, 
  fallback, 
  showPreview = true 
}) => {
  const { 
    userProfile, 
    isFeatureUnlocked, 
    FEATURE_GATES,
    getFeatureAccess 
  } = useGamification();

  // Verificar si la feature está desbloqueada
  const isUnlocked = isFeatureUnlocked(feature);
  // Consultar info de acceso para detectar trial (si el Context no marcó como unlocked por backend)
  // Nota: getFeatureAccess ya usa datos del backend cuando existen
  const access = getFeatureAccess ? getFeatureAccess(feature) : null;
  const featureData = FEATURE_GATES[feature];
  const currentLevel = userProfile?.current_level || 0;

  // Si la feature está desbloqueada, mostrar contenido completo
  if (isUnlocked) {
    return children;
  }

  // Si hay un fallback personalizado, usarlo
  if (fallback) {
    return fallback;
  }

  // Manejar diferentes modos de bloqueo
  switch (mode) {
    case 'block':
      // Bloqueo completo - no mostrar nada
      return null;

    case 'preview':
      // Mostrar preview de la feature
      if (showPreview) {
        return (
          <LockedFeaturePreview 
            feature={feature}
            featureData={featureData}
            userLevel={currentLevel}
            // Mostrar banner si está en trial aunque figure como bloqueada
            trialActive={Boolean(access?.trialActive)}
            trialEndsAt={access?.trialEndsAt || null}
            mode="preview"
          />
        );
      }
      return null;

    case 'limited':
      // Mostrar feature con limitaciones
      return (
        <div className="relative">
          {children}
          <FeatureLimitationOverlay 
            feature={feature}
            featureData={featureData}
            userLevel={currentLevel}
          />
        </div>
      );

    case 'full':
    default:
      // Modo por defecto - mostrar preview
      return showPreview ? (
        <LockedFeaturePreview 
          feature={feature}
          featureData={featureData}
          userLevel={currentLevel}
          trialActive={Boolean(access?.trialActive)}
          trialEndsAt={access?.trialEndsAt || null}
          mode="full"
        />
      ) : null;
  }
};

/**
 * Overlay que se muestra sobre features limitadas
 */
const FeatureLimitationOverlay = ({ feature, featureData, userLevel }) => {
  const requiredLevel = featureData.requiredLevel;
  const xpNeeded = calculateXPNeeded(userLevel, requiredLevel);

  return (
    <div className="absolute inset-0 bg-white/90 dark:bg-gray-900/90 backdrop-blur-sm flex items-center justify-center z-10">
      <div className="text-center p-6 max-w-md">
        <div className="text-4xl mb-3">{featureData.icon}</div>
        <h3 className="text-lg font-semibold mb-2">Feature Limitada</h3>
        <p className="text-gray-600 dark:text-gray-400 mb-4">
          Desbloquea funcionalidades completas en Nivel {requiredLevel}
        </p>
        <ProgressToUnlock currentLevel={userLevel} targetLevel={requiredLevel} />
        <button className="mt-4 bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600">
          Continúa usando la app ({xpNeeded} XP restantes)
        </button>
      </div>
    </div>
  );
};

/**
 * Calcula XP necesario para alcanzar un nivel específico
 */
const calculateXPNeeded = (currentLevel, targetLevel) => {
  const LEVEL_THRESHOLDS = [0, 100, 250, 500, 1000, 2000, 4000, 8000, 16000, 32000];
  const currentXP = LEVEL_THRESHOLDS[currentLevel] || 0;
  const targetXP = LEVEL_THRESHOLDS[targetLevel] || 0;
  return Math.max(0, targetXP - currentXP);
};

/**
 * Componente de progreso hacia desbloqueo
 */
const ProgressToUnlock = ({ currentLevel, targetLevel }) => {
  const LEVEL_THRESHOLDS = [0, 100, 250, 500, 1000, 2000, 4000, 8000, 16000, 32000];
  const currentXP = LEVEL_THRESHOLDS[currentLevel] || 0;
  const targetXP = LEVEL_THRESHOLDS[targetLevel] || 0;
  const progress = currentXP / targetXP * 100;

  return (
    <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3 mb-2">
      <div 
        className="bg-blue-500 h-3 rounded-full transition-all duration-300"
        style={{ width: `${Math.min(progress, 100)}%` }}
      />
    </div>
  );
};

export default FeatureGuard; 