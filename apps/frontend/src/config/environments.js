// Configuración de URLs por ambiente
const environments = {
  development: {
    name: 'Development',
    API_BASE_URL: 'http://localhost:8080/api/v1',
    GAMIFICATION_API_URL: 'http://localhost:8081/api/v1',
    AI_API_URL: 'http://localhost:8082/api/v1',
    USERS_API_URL: 'http://localhost:8083/api/v1',
    REDIS_URL: 'redis://localhost:6379',
    WEBSOCKET_URL: 'ws://localhost:8080/ws',
    CORS_ORIGIN: 'http://localhost:3000',
    // Configuración específica para desarrollo
    RATE_LIMIT_DISABLED: true,
    REQUEST_THROTTLE_MS: 100,  // Throttle mínimo entre requests
    CONFIG_CACHE_DISABLED: true  // Deshabilitar cache de configuración
  },
  
  docker: {
    name: 'Docker Compose',
    API_BASE_URL: 'http://financial-resume-engine:8080/api/v1',
    GAMIFICATION_API_URL: 'http://gamification-service:8081/api/v1',
    AI_API_URL: 'http://ai-service:8082/api/v1',
    USERS_API_URL: 'http://users-service:8083/api/v1',
    REDIS_URL: 'redis://financial_redis:6379',
    WEBSOCKET_URL: 'ws://financial-resume-engine:8080/ws',
    CORS_ORIGIN: 'http://localhost:3000',
    // Configuración específica para Docker
    RATE_LIMIT_DISABLED: true,
    REQUEST_THROTTLE_MS: 100,
    CONFIG_CACHE_DISABLED: true
  },

  render: {
    name: 'Render.com',
    API_BASE_URL: 'https://financial-resume-engine.onrender.com/api/v1',
    GAMIFICATION_API_URL: 'https://financial-gamification-service.onrender.com/api/v1',
    AI_API_URL: 'https://financial-ai-api.niloft.com/api/v1',
    USERS_API_URL: 'https://users-service-mp5p.onrender.com/api/v1',
    REDIS_URL: 'redis://red-d1qmg0juibrs73esqdfg:6379',
    WEBSOCKET_URL: 'wss://financial-resume-engine.onrender.com/ws',
    CORS_ORIGIN: 'https://financial-resume-engine-frontend.onrender.com',
    // Configuración de producción
    DISABLE_CONSOLE_LOGS: true,
    LOG_LEVEL: 'WARN' // ERROR, WARN, INFO, DEBUG
  },

  production: {
    name: 'Production (Niloft)',
    API_BASE_URL: 'https://financial.niloft.com/api/v1',
    GAMIFICATION_API_URL: 'https://financial-gamification-service.onrender.com/api/v1',
    AI_API_URL: 'https://financial-ai-api.niloft.com/api/v1',
    USERS_API_URL: 'https://users-service-mp5p.onrender.com/api/v1',
    REDIS_URL: 'redis://red-d1qmg0juibrs73esqdfg:6379',
    WEBSOCKET_URL: 'wss://financial.niloft.com/ws',
    CORS_ORIGIN: 'https://financial.niloft.com',
    // Configuración de producción
    DISABLE_CONSOLE_LOGS: true,
    LOG_LEVEL: 'WARN' // ERROR, WARN, INFO, DEBUG
  },
  
};

// Detectar ambiente automáticamente
const detectEnvironment = () => {
  const hostname = window.location.hostname;
  const port = window.location.port;
  
  // Si hay variable de entorno específica, usarla (prioridad máxima)
  const envFromVar = process.env.REACT_APP_ENVIRONMENT;
  if (envFromVar) {
    console.log(`🔧 [environments] Ambiente forzado por variable de entorno: ${envFromVar}`);
    return envFromVar;
  }
  
  // Si estamos en localhost, es desarrollo
  if (hostname === 'localhost' || hostname === '127.0.0.1') {
    console.log(`🔧 [environments] Ambiente detectado por hostname: development (${hostname}:${port})`);
    return 'development';
  }
  
  // Configuración específica para financial.niloft.com
  if (hostname === 'financial.niloft.com') {
    console.log(`🔧 [environments] Ambiente detectado por hostname: production (${hostname})`);
    return 'production';
  }
  
  // Solo usar render si realmente estamos en onrender.com
  if (hostname.includes('onrender.com')) {
    console.log(`🔧 [environments] Ambiente detectado por hostname: render (${hostname})`);
    return 'render';
  }
  
  // Por defecto, desarrollo
  console.log(`🔧 [environments] Ambiente por defecto: development (${hostname}:${port})`);
  return 'development';
};

// Configuración actual basada en el ambiente
const currentEnvironment = detectEnvironment();
const config = environments[currentEnvironment];

// Función para obtener configuración con fallback
const getConfig = (key, fallback = null) => {
  // Mapeo de claves a variables de entorno
  const envKeyMap = {
    'API_BASE_URL': 'REACT_APP_API_URL',
    'GAMIFICATION_API_URL': 'REACT_APP_GAMIFICATION_URL',
    'AI_API_URL': 'REACT_APP_AI_SERVICE_URL',
    'USERS_API_URL': 'REACT_APP_USERS_SERVICE_URL'
  };
  
  // Primero intentar variable de entorno específica
  const envKey = envKeyMap[key] || `REACT_APP_${key}`;
  const envVar = process.env[envKey];
  if (envVar) {
    console.log(`🔧 [environments] Usando variable de entorno ${envKey}: ${envVar}`);
    return envVar;
  }
  
  // Luego usar configuración del ambiente
  if (config && config[key]) {
    console.log(`🔧 [environments] Usando configuración del ambiente ${activeEnvironment}: ${config[key]}`);
    return config[key];
  }
  
  // Finalmente usar fallback
  console.log(`🔧 [environments] Usando fallback para ${key}: ${fallback}`);
  return fallback;
};

// Función para cambiar ambiente manualmente (útil para debug)
const setEnvironment = (env) => {
  if (environments[env]) {
    localStorage.setItem('FORCE_ENVIRONMENT', env);
    window.location.reload();
  }
};

// Función para resetear ambiente forzado
const resetEnvironment = () => {
  localStorage.removeItem('FORCE_ENVIRONMENT');
  window.location.reload();
};

// Verificar si hay ambiente forzado
const forcedEnv = localStorage.getItem('FORCE_ENVIRONMENT');
const activeEnvironment = forcedEnv && environments[forcedEnv] ? forcedEnv : currentEnvironment;
const activeConfig = environments[activeEnvironment];

// Exportar configuración
export default {
  // Información del ambiente
  ENVIRONMENT: activeEnvironment,
  ENVIRONMENT_NAME: activeConfig.name,
  
  // URLs de servicios
  API_BASE_URL: getConfig('API_BASE_URL', activeConfig.API_BASE_URL),
  GAMIFICATION_API_URL: getConfig('GAMIFICATION_API_URL', activeConfig.GAMIFICATION_API_URL),
  AI_API_URL: getConfig('AI_API_URL', activeConfig.AI_API_URL),
  USERS_API_URL: getConfig('USERS_API_URL', activeConfig.USERS_API_URL),
  REDIS_URL: getConfig('REDIS_URL', activeConfig.REDIS_URL),
  WEBSOCKET_URL: getConfig('WEBSOCKET_URL', activeConfig.WEBSOCKET_URL),
  CORS_ORIGIN: getConfig('CORS_ORIGIN', activeConfig.CORS_ORIGIN),
  
  // Configuración de desarrollo
  IS_DEVELOPMENT: activeEnvironment === 'development',
  IS_PRODUCTION: activeEnvironment !== 'development',
  
  // Funciones utilitarias
  getAllEnvironments: () => environments,
  setEnvironment,
  resetEnvironment,
  
  // Información de debug
  debug: {
    hostname: window.location.hostname,
    detectedEnv: currentEnvironment,
    activeEnv: activeEnvironment,
    forcedEnv,
    config: activeConfig
  }
};

// Debug en consola (solo en desarrollo o si está habilitado)
if (activeEnvironment === 'development' || !activeConfig.DISABLE_CONSOLE_LOGS) {
  console.log('🔧 Environment Config:', {
    environment: activeEnvironment,
    config: activeConfig,
    debug: {
      hostname: window.location.hostname,
      detectedEnv: currentEnvironment,
      forcedEnv
    }
  });
} 