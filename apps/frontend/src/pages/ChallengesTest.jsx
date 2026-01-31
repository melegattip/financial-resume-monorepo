import React from 'react';
import DailyChallenges from '../components/DailyChallenges';

/**
 * 🧪 CHALLENGES TEST PAGE
 * 
 * Página temporal para probar la integración de challenges
 */
const ChallengesTest = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            🧪 Prueba de Challenges
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            Página temporal para probar la integración del sistema de challenges con el backend
          </p>
        </div>
        
        <DailyChallenges />
      </div>
    </div>
  );
};

export default ChallengesTest; 