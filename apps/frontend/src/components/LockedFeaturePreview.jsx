import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useGamification } from '../contexts/GamificationContext';

/**
 * 🔒 LOCKED FEATURE PREVIEW
 * 
 * Muestra un preview atractivo de features bloqueadas para motivar al usuario
 * a continuar usando la app y alcanzar el nivel requerido
 */
const LockedFeaturePreview = ({ 
  feature, 
  featureData, 
  userLevel,
  mode = 'preview' 
}) => {
  const navigate = useNavigate();
  const { LEVEL_SYSTEM, userProfile } = useGamification();

  const requiredLevel = featureData.requiredLevel;
  const levelInfo = LEVEL_SYSTEM[requiredLevel];
  const currentLevelInfo = LEVEL_SYSTEM[userLevel] || LEVEL_SYSTEM[1];
  
  // Calcular progreso usando XP real del usuario
  const currentXP = userProfile?.total_xp || 0;
  const targetXP = levelInfo.minXP;
  const xpNeeded = Math.max(0, targetXP - currentXP);
  const progress = Math.min(100, (currentXP / targetXP) * 100);

  const handleContinueUsingApp = () => {
    // Redirigir al dashboard para que continúe ganando XP
    navigate('/dashboard');
  };

  const handleViewAchievements = () => {
    navigate('/achievements');
  };

  return (
    <div className="min-h-[calc(100vh-4rem)] bg-gray-50 dark:bg-gray-900 flex items-center justify-center p-4">
      <div className="max-w-md w-full">
        
        {/* Card principal */}
        <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border dark:border-gray-700 overflow-hidden">
          
          {/* Header con gradiente */}
          <div 
            className="px-6 py-8 text-center text-white relative"
            style={{
              background: `linear-gradient(135deg, ${levelInfo.color}20, ${levelInfo.color}40)`
            }}
          >
            <div className="absolute inset-0 bg-gradient-to-br from-blue-500 to-purple-600 opacity-90"></div>
            <div className="relative z-10">
              <div className="text-6xl mb-4">{featureData.icon}</div>
              <h2 className="text-2xl font-bold mb-2">{featureData.name}</h2>
              <p className="text-blue-100">{featureData.description}</p>
            </div>
          </div>

          {/* Contenido */}
          <div className="p-6">
            
            {/* Nivel requerido */}
            <div className="text-center mb-6">
              <div className="inline-flex items-center bg-gray-100 dark:bg-gray-700 rounded-full px-4 py-2 mb-4">
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  Requerido: Nivel {requiredLevel}
                </span>
                <span 
                  className="ml-2 w-3 h-3 rounded-full"
                  style={{ backgroundColor: levelInfo.color }}
                ></span>
              </div>
              <h3 className="font-semibold text-gray-900 dark:text-gray-100">
                {levelInfo.name}
              </h3>
            </div>

            {/* Progreso hacia el unlock */}
            <div className="mb-6">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
                  Tu progreso
                </span>
                <span className="text-sm font-bold text-gray-900 dark:text-gray-100">
                  Nivel {userLevel}
                </span>
              </div>
              
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-4 mb-2">
                <div 
                  className="bg-gradient-to-r from-blue-500 to-purple-600 h-4 rounded-full transition-all duration-500 relative"
                  style={{ width: `${Math.min(progress, 100)}%` }}
                >
                  {progress > 10 && (
                    <div className="absolute inset-0 bg-white/20 rounded-full animate-pulse"></div>
                  )}
                </div>
              </div>
              
              <p className="text-center text-sm text-gray-600 dark:text-gray-400">
                <strong>{xpNeeded} XP</strong> restantes para desbloquear
              </p>
            </div>

            {/* Beneficios */}
            <div className="mb-6">
              <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-3">
                Lo que desbloquearás:
              </h4>
              <div className="space-y-2">
                {featureData.benefits.map((benefit, index) => (
                  <div key={index} className="flex items-center">
                    <div className="w-5 h-5 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center mr-3">
                      <svg className="w-3 h-3 text-green-600 dark:text-green-400" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    </div>
                    <span className="text-sm text-gray-700 dark:text-gray-300">{benefit}</span>
                  </div>
                ))}
              </div>
            </div>

            {/* Motivación */}
            <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4 mb-6">
              <div className="flex items-center mb-2">
                <span className="text-2xl mr-2">💡</span>
                <span className="font-semibold text-blue-900 dark:text-blue-100">
                  ¡Sigue ganando XP!
                </span>
              </div>
              <p className="text-sm text-blue-700 dark:text-blue-300">
                Cada transacción que registres, cada análisis que uses y cada meta que alcances 
                te acerca más a desbloquear esta funcionalidad.
              </p>
            </div>

            {/* Acciones */}
            <div className="space-y-3">
              <button
                onClick={handleContinueUsingApp}
                className="w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white py-3 px-4 rounded-lg font-semibold hover:from-blue-600 hover:to-purple-700 transition-all duration-200 transform hover:scale-[1.02]"
              >
                Continuar usando la app
              </button>
              
              <button
                onClick={handleViewAchievements}
                className="w-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 py-3 px-4 rounded-lg font-medium hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
              >
                Ver mis logros
              </button>
            </div>

            {/* Footer motivacional */}
            <div className="mt-6 text-center">
              <p className="text-xs text-gray-500 dark:text-gray-400">
                Únete a los <strong>usuarios nivel {requiredLevel}+</strong> que ya disfrutan de esta funcionalidad
              </p>
            </div>
          </div>
        </div>

        {/* Tips para ganar XP rápido - ACTUALIZADO SIN DEPENDENCIAS */}
        <div className="mt-4 bg-white dark:bg-gray-800 rounded-lg p-4 border dark:border-gray-700">
          <h5 className="font-semibold text-gray-900 dark:text-gray-100 mb-2 text-center">
            🚀 Formas disponibles de ganar XP:
          </h5>
          <div className="grid grid-cols-2 gap-2 text-xs">
            <div className="text-center p-2 bg-gray-50 dark:bg-gray-700 rounded">
              <div className="font-semibold text-blue-600 dark:text-blue-400">+8 XP</div>
              <div>Registrar transacción</div>
            </div>
            <div className="text-center p-2 bg-gray-50 dark:bg-gray-700 rounded">
              <div className="font-semibold text-green-600 dark:text-green-400">+10 XP</div>
              <div>Crear categoría</div>
            </div>
            <div className="text-center p-2 bg-gray-50 dark:bg-gray-700 rounded">
              <div className="font-semibold text-purple-600 dark:text-purple-400">+20 XP</div>
              <div>Challenge diario</div>
            </div>
            <div className="text-center p-2 bg-gray-50 dark:bg-gray-700 rounded">
              <div className="font-semibold text-orange-600 dark:text-orange-400">+25 XP</div>
              <div>Racha semanal</div>
            </div>
          </div>
          <div className="mt-3 p-2 bg-blue-50 dark:bg-blue-900/20 rounded text-center">
            <p className="text-xs text-blue-600 dark:text-blue-400 font-medium">
              💡 Completa challenges diarios para maximizar tu progreso
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LockedFeaturePreview; 