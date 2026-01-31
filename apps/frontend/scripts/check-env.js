#!/usr/bin/env node

// Script para verificar las variables de entorno en el build
console.log('🔍 Verificando variables de entorno...\n');

const envVars = [
  'REACT_APP_API_URL',
  'REACT_APP_USERS_SERVICE_URL',
  'REACT_APP_GAMIFICATION_URL',
  'REACT_APP_AI_SERVICE_URL',
  'REACT_APP_ENV',
  'REACT_APP_ENVIRONMENT'
];

envVars.forEach(varName => {
  const value = process.env[varName];
  if (value) {
    console.log(`✅ ${varName}: ${value}`);
  } else {
    console.log(`❌ ${varName}: NO DEFINIDA`);
  }
});

console.log('\n📋 Variables de entorno del sistema:');
console.log('NODE_ENV:', process.env.NODE_ENV);
console.log('GENERATE_SOURCEMAP:', process.env.GENERATE_SOURCEMAP); 