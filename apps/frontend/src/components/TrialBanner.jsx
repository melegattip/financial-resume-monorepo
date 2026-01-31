import React, { useEffect, useMemo, useState } from 'react';
import { useGamification } from '../contexts/GamificationContext';

// Muestra un banner con cuenta regresiva de trial para una feature
// featureKey: 'AI_INSIGHTS' | 'BUDGETS' | 'SAVINGS_GOALS'
const TrialBanner = ({ featureKey }) => {
  const { checkFeatureAccess, userProfile, isFeatureUnlocked, features } = useGamification();
  const [access, setAccess] = useState(null);
  const [now, setNow] = useState(Date.now());
  const [hasChecked, setHasChecked] = useState(false);

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      // Si la feature ya está desbloqueada, no necesitamos verificar trial
      if (isFeatureUnlocked && isFeatureUnlocked(featureKey)) {
        setAccess(null);
        setHasChecked(true);
        return;
      }

      try {
        const res = await checkFeatureAccess(featureKey);
        if (mounted) {
          setAccess(res);
          setHasChecked(true);
        }
      } catch (_) {
        if (mounted) {
          setAccess(null);
          setHasChecked(true);
        }
      }
    };
    
    if (checkFeatureAccess && !hasChecked) {
      load();
    }
    
    const timer = setInterval(() => setNow(Date.now()), 1000);
    return () => { mounted = false; clearInterval(timer); };
  }, [featureKey, checkFeatureAccess, isFeatureUnlocked, hasChecked]);

  const { trialActive, remaining } = useMemo(() => {
    // 1) Verificar datos del backend (locked_features con trial_active)
    let trial = false;
    let endTs = null;

    if (features && features.locked_features && Array.isArray(features.locked_features)) {
      const lockedFeature = features.locked_features.find(f => f.feature_key === featureKey);
      if (lockedFeature && lockedFeature.trial_active) {
        trial = true;
        endTs = lockedFeature.trial_ends_at ? new Date(lockedFeature.trial_ends_at).getTime() : null;
      }
    }

    // 2) Fallback: verificar con endpoint específico
    if (!trial && access && access.trial_active) {
      trial = true;
      endTs = access.trial_ends_at ? new Date(access.trial_ends_at).getTime() : null;
    }

    // 3) Fallback final: usar created_at del perfil de gamificación
    if (!trial && userProfile && userProfile.created_at) {
      const created = new Date(userProfile.created_at).getTime();
      const tenDays = 10 * 24 * 3600 * 1000;
      const altEnd = created + tenDays;
      trial = now < altEnd;
      if (trial && !endTs) endTs = altEnd;
    }

    if (!trial) return { trialActive: false, remaining: null };
    if (!endTs) return { trialActive: true, remaining: null };
    
    const end = endTs;
    const ms = Math.max(0, end - now);
    const d = Math.floor(ms / (24 * 3600 * 1000));
    const h = Math.floor((ms % (24 * 3600 * 1000)) / (3600 * 1000));
    const m = Math.floor((ms % (3600 * 1000)) / (60 * 1000));
    const s = Math.floor((ms % (60 * 1000)) / 1000);
    return { trialActive: ms > 0, remaining: { d, h, m, s } };
  }, [access, userProfile, now, featureKey, features]);

  if (!trialActive) return null;

  const label = featureKey === 'AI_INSIGHTS'
    ? 'IA Financiera'
    : featureKey === 'BUDGETS'
    ? 'Presupuestos'
    : 'Metas de Ahorro';

  return (
    <div className="mx-4 mt-4">
      <div className="relative overflow-hidden rounded-xl border border-blue-200 dark:border-blue-900 bg-blue-50 dark:bg-blue-900/20">
        <div className="px-4 py-3 sm:px-6 sm:py-4">
          <div className="flex items-start sm:items-center sm:justify-between flex-col sm:flex-row">
            <div className="flex items-center mb-2 sm:mb-0">
              <span className="text-2xl mr-3">✨</span>
              <div>
                <div className="text-blue-900 dark:text-blue-100 font-semibold">
                  Acceso de prueba a {label}
                </div>
                <div className="text-sm text-blue-800/80 dark:text-blue-200/80">
                  Disfruta todas las funcionalidades durante los primeros 10 días
                </div>
              </div>
            </div>
            <div className="text-sm font-mono bg-white/70 dark:bg-blue-900/40 px-3 py-1 rounded-md text-blue-900 dark:text-blue-100">
              {remaining ? (
                <span>{String(remaining.d).padStart(2,'0')}d:{String(remaining.h).padStart(2,'0')}h:{String(remaining.m).padStart(2,'0')}m:{String(remaining.s).padStart(2,'0')}s</span>
              ) : (
                <span>En curso</span>
              )}
            </div>
          </div>
        </div>
        <div className="h-1 bg-blue-200 dark:bg-blue-900">
          <div className="h-full bg-blue-500 dark:bg-blue-400 animate-pulse" style={{ width: '100%' }} />
        </div>
      </div>
    </div>
  );
};

export default TrialBanner;


