import React from 'react';
import AIInsights from '../components/AIInsights';
import TrialBanner from '../components/TrialBanner';

const FinancialInsights = () => {
  return (
    <div className="space-y-6">
      <TrialBanner featureKey="AI_INSIGHTS" />
      {/* Componente principal de insights */}
      <AIInsights />
    </div>
  );
};

export default FinancialInsights; 