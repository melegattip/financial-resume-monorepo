/**
 * 🔒 LOCKED WIDGET COMPONENT
 * 
 * Componente que muestra un widget bloqueado por nivel con información
 * sobre cómo desbloquear la funcionalidad a través de gamificación
 */

import React from 'react';
import { FaLock, FaStar, FaTrophy, FaBolt, FaGem, FaRocket, FaUnlock } from 'react-icons/fa';
import { useNavigate } from 'react-router-dom';

const LockedWidget = ({ 
  featureName, 
  featureIcon, 
  description, 
  requiredLevel, 
  currentLevel, 
  currentXP, 
  requiredXP,
  benefits = [],
  tips = [],
  mode = 'detailed' // 'minimal' | 'detailed' | 'compact'
}) => {
  const navigate = useNavigate();

  const xpNeeded = requiredXP - currentXP;
  const progressPercentage = Math.min((currentXP / requiredXP) * 100, 100);

  const handleUnlockClick = () => {
    navigate('/achievements');
  };

  const getLevelColor = (level) => {
    const colors = {
      3: 'from-blue-400/20 to-amber-500/30',
      5: 'from-green-400/20 to-blue-500/30', 
      7: 'from-red-400/20 to-purple-500/30'
    };
    return colors[level] || 'from-gray-400/40 to-gray-500/50';
  };

  const getLevelIcon = (level) => {
    if (level <= 3) return <FaGem className="w-3 h-3" />;
    if (level <= 5) return <FaTrophy className="w-3 h-3" />;
    return <FaRocket className="w-3 h-3" />;
  };

  // Modo compacto para dashboard optimizado
  if (mode === 'compact') {
    return (
      <div className="relative group cursor-pointer" onClick={handleUnlockClick}>
        <div className={`relative bg-gradient-to-br ${getLevelColor(requiredLevel)} rounded-xl p-3 sm:p-4 shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-[1.02] border border-white/10`}>
          <div className="flex items-start justify-between mb-2">
            <div className="flex-1 min-w-0">
              <p className="text-xs sm:text-sm font-medium text-white/90">
                {featureName}
              </p>
              <p className="text-lg sm:text-xl font-bold text-white break-words">
                Bloqueado
              </p>
            </div>
            <div className="flex-shrink-0 p-1.5 sm:p-2 rounded-fr bg-white/15 backdrop-blur-sm ml-2">
              <FaLock className="w-3 h-3 sm:w-4 sm:h-4 text-white" />
            </div>
          </div>
          <div className="text-xs text-white/80">
            {getLevelIcon(requiredLevel)}
            <span className="ml-1">Nivel {requiredLevel} requerido</span>
          </div>
        </div>
      </div>
    );
  }

  // Modo minimalista para el dashboard
  if (mode === 'minimal') {
    return (
      <div className="relative group cursor-pointer" onClick={handleUnlockClick}>
        {/* Card minimalista - estilo elegante con color por nivel */}
        <div className={`relative bg-gradient-to-br ${getLevelColor(requiredLevel)} rounded-xl p-5 shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-[1.02] h-[200px] border border-white/10`}>
          
          {/* Contenido centrado */}
          <div className="flex flex-col items-center text-center text-white h-full justify-between">
            
            {/* Candado elegante */}
            <div className="bg-white/15 backdrop-blur-sm rounded-2xl p-3 mb-3 shadow-lg border border-white/20">
              <FaLock className="w-7 h-7 text-white" />
            </div>

            {/* Información */}
            <div className="flex-1 flex flex-col justify-center">
              <div className="flex items-center justify-center gap-2 mb-2">
                <span className="text-lg">{featureIcon}</span>
                <h3 className="font-bold text-white">{featureName}</h3>
              </div>
              
              <p className="text-sm text-white/90 mb-3 leading-relaxed line-clamp-1">
                {description}
              </p>
            </div>

            {/* Badge inferior elegante */}
            <div className="bg-white/10 backdrop-blur-sm px-3 py-2 rounded-full flex items-center gap-2 border border-white/20">
              {getLevelIcon(requiredLevel)}
              <span className="text-sm font-medium text-white">Nivel {requiredLevel} Requerido</span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Modo detallado para páginas individuales
  return (
    <div className="relative group z-30">
      {/* Contenedor principal con gradiente y efectos */}
      <div className="relative bg-gradient-to-br from-slate-800 via-slate-700 to-slate-900 rounded-2xl overflow-hidden border border-slate-600/50 shadow-2xl transform transition-all duration-300 hover:scale-[1.02] hover:shadow-3xl">
        
        {/* Efectos de fondo animados */}
        <div className="absolute inset-0 bg-gradient-to-br from-blue-500/10 via-purple-500/10 to-indigo-500/10 opacity-60"></div>
        <div className="absolute -top-1/2 -right-1/2 w-96 h-96 bg-gradient-to-br from-blue-400/20 to-purple-600/20 rounded-full blur-3xl animate-pulse"></div>
        
        {/* Header con candado prominente */}
        <div className="relative z-10 p-6 text-center">
          
          {/* Candado principal con animación */}
          <div className="relative mx-auto mb-6">
            <div className={`w-20 h-20 bg-gradient-to-br ${getLevelColor(requiredLevel)} rounded-2xl flex items-center justify-center shadow-lg transform transition-transform duration-300 group-hover:rotate-3`}>
              <FaLock className="w-8 h-8 text-white drop-shadow-lg" />
            </div>
            {/* Badge de nivel requerido */}
            <div className="absolute -top-2 -right-2 bg-gradient-to-r from-orange-400 to-red-500 text-white text-xs font-bold px-2 py-1 rounded-full shadow-md">
              {requiredLevel}
            </div>
          </div>

          {/* Título de la feature */}
          <div className="mb-4">
            <h3 className="text-xl font-bold text-white mb-2 flex items-center justify-center gap-2">
              <span className="text-2xl">{featureIcon}</span>
              {featureName}
            </h3>
            <p className="text-slate-300 text-sm leading-relaxed max-w-xs mx-auto">
              {description}
            </p>
          </div>

          {/* Badge de nivel requerido más prominente */}
          <div className="inline-flex items-center gap-2 bg-white/10 backdrop-blur-sm px-4 py-2 rounded-full border border-white/20 mb-6">
            {getLevelIcon(requiredLevel)}
            <span className="text-white font-semibold">Nivel {requiredLevel} Requerido</span>
          </div>
        </div>

        {/* Sección de progreso */}
        <div className="relative z-10 px-6 pb-6">
          <div className="bg-white/5 backdrop-blur-sm rounded-xl p-4 border border-white/10">
            
            {/* Estado actual */}
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2 text-slate-300">
                <FaTrophy className="w-4 h-4 text-amber-400" />
                <span className="text-sm font-medium">Tu Nivel: {currentLevel}</span>
              </div>
              <div className="flex items-center gap-2 text-slate-300">
                <FaBolt className="w-4 h-4 text-blue-400" />
                <span className="text-sm font-medium">{currentXP.toLocaleString()} XP</span>
              </div>
            </div>

            {/* Barra de progreso mejorada */}
            <div className="mb-4">
              <div className="flex justify-between text-xs text-slate-400 mb-2">
                <span>Progreso</span>
                <span>{Math.round(progressPercentage)}%</span>
              </div>
              <div className="relative w-full bg-slate-700 rounded-full h-3 overflow-hidden">
                <div 
                  className="absolute top-0 left-0 h-full bg-gradient-to-r from-emerald-400 via-blue-500 to-purple-600 rounded-full transition-all duration-700 ease-out shadow-lg"
                  style={{ width: `${progressPercentage}%` }}
                >
                  <div className="absolute inset-0 bg-white/30 rounded-full animate-pulse"></div>
                </div>
              </div>
              {xpNeeded > 0 && (
                <p className="text-center text-slate-300 text-sm mt-2 font-medium">
                  <FaBolt className="inline w-3 h-3 mr-1 text-yellow-400" />
                  {xpNeeded.toLocaleString()} XP restantes
                </p>
              )}
            </div>

            {/* Beneficios en grid */}
            {benefits.length > 0 && (
              <div className="mb-4">
                <h4 className="text-white font-semibold text-sm mb-3 flex items-center gap-2">
                  <FaGem className="w-4 h-4 text-emerald-400" />
                  Beneficios al Desbloquear
                </h4>
                <div className="space-y-2">
                  {benefits.slice(0, 3).map((benefit, index) => (
                    <div key={index} className="flex items-center gap-2 text-slate-300 text-xs">
                      <FaStar className="w-3 h-3 text-amber-400 flex-shrink-0" />
                      <span>{benefit}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Botón de acción mejorado */}
            <button
              onClick={handleUnlockClick}
              className="w-full bg-gradient-to-r from-blue-500 via-purple-500 to-indigo-600 hover:from-blue-600 hover:via-purple-600 hover:to-indigo-700 text-white font-semibold py-3 px-4 rounded-xl transition-all duration-300 transform hover:scale-105 hover:shadow-xl flex items-center justify-center gap-2 border border-blue-400/30"
            >
              <FaUnlock className="w-4 h-4" />
              <span>Ver Cómo Desbloquear</span>
            </button>
          </div>
        </div>

        {/* Contenido de fondo difuminado para referencia */}
        <div className="absolute inset-0 z-0">
          <div className="filter blur-[2px] opacity-20 p-6">
            <div className="flex items-start justify-between">
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-600 dark:text-gray-400">
                  {featureName}
                </p>
                <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                  ••••••
                </p>
              </div>
              <div className="flex-shrink-0 p-3 rounded-lg bg-gray-100 dark:bg-gray-700 ml-2">
                {featureIcon}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LockedWidget; 