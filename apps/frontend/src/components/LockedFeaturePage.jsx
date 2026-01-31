/**
 * 🔒 LOCKED FEATURE PAGE COMPONENT
 * 
 * Página completa para mostrar cuando una feature está bloqueada
 * Usa el LockedWidget en modo detallado con información completa
 */

import React from 'react';
import { useGamification } from '../contexts/GamificationContext';
import LockedWidget from './LockedWidget';

const LockedFeaturePage = ({ 
  feature,
  featureName,
  featureIcon,
  description,
  tips = []
}) => {
  const { userProfile, FEATURE_GATES } = useGamification();

  const featureData = FEATURE_GATES[feature];
  
  if (!featureData) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <p className="text-gray-500">Feature no encontrada</p>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto py-12">
      <LockedWidget
        mode="detailed"
        featureName={featureName || featureData.name}
        featureIcon={featureIcon}
        description={description || featureData.description}
        requiredLevel={featureData.requiredLevel}
        currentLevel={userProfile?.current_level || 0}
        currentXP={userProfile?.total_xp || 0}
        requiredXP={featureData.xpThreshold}
        benefits={featureData.benefits}
        tips={tips}
      />
    </div>
  );
};

export default LockedFeaturePage; 