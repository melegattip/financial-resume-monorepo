import axios from 'axios';
import configService from './configService';

// Configuración base de axios - SOLO comunicación HTTP
// Función para determinar baseURL inicial por ambiente
const getInitialBaseURL = () => {
  // Variable de entorno tiene prioridad
  if (process.env.REACT_APP_API_URL) {
    return process.env.REACT_APP_API_URL;
  }
  
  // Detección por hostname
  const hostname = window.location.hostname;
  if (hostname.includes('onrender.com') || hostname === 'financial.niloft.com') {
    return 'https://financial-resume-monorepo-l71a.onrender.com/api/v1';  // Render monolith
  } else {
    return 'http://localhost:8080/api/v1';  // Development
  }
};

const apiClient = axios.create({
  baseURL: getInitialBaseURL(),
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Función para inicializar la configuración dinámica
let configInitialized = false;

// Cache de requests para evitar duplicados en desarrollo
const requestCache = new Map();
const CACHE_DURATION = 1000; // 1 segundo en desarrollo

const initializeConfig = async () => {
  if (configInitialized) return;
  
  try {
    console.log('🔄 [apiClient] Inicializando configuración dinámica...');
    const config = await configService.loadConfig();
    
    // Actualizar la baseURL de axios con la configuración dinámica
    apiClient.defaults.baseURL = config.api_base_url;
    configInitialized = true;
    
    console.log('✅ [apiClient] Configuración dinámica inicializada:', {
      baseURL: apiClient.defaults.baseURL,
      environment: config.environment,
      version: config.version
    });
  } catch (error) {
    console.error('❌ [apiClient] Error inicializando configuración:', error);
    // Mantener la configuración por defecto
  }
};

// Inicializar configuración al cargar el módulo
initializeConfig();

// Función para obtener headers de autenticación
const getAuthHeaders = () => {
  const token = localStorage.getItem('auth_token');
  const userData = localStorage.getItem('auth_user');
  
  const headers = {};
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  if (userData) {
    try {
      const user = JSON.parse(userData);
      if (user?.id) {
        headers['X-Caller-ID'] = user.id.toString();
      }
    } catch (error) {
      console.error('Error parsing user data:', error);
    }
  }
  
  return headers;
};

// Interceptor para agregar headers de autenticación automáticamente
apiClient.interceptors.request.use(
  async (config) => {
    // En desarrollo, evitar múltiples inicializaciones de configuración
    const isDevelopment = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    
    if (!configInitialized && !isDevelopment) {
      // Solo en producción intentar reconfigurar en cada request
      await initializeConfig();
    }
    
    // En desarrollo, agregar throttling para evitar rate limiting
    if (isDevelopment) {
      const requestKey = `${config.method?.toUpperCase()}_${config.url}`;
      const now = Date.now();
      const lastRequest = requestCache.get(requestKey);
      
      if (lastRequest && (now - lastRequest) < CACHE_DURATION) {
        // Agregar un pequeño delay mínimo para evitar spam
        const delay = 50; // Solo 50ms de delay
        await new Promise(resolve => setTimeout(resolve, delay));
      }
      
      requestCache.set(requestKey, now);
      
      // Limpiar cache viejo cada minuto
      if (requestCache.size > 100) {
        const cutoff = now - 60000; // 1 minuto
        for (const [key, timestamp] of requestCache.entries()) {
          if (timestamp < cutoff) {
            requestCache.delete(key);
          }
        }
      }
    }
    
    const authHeaders = getAuthHeaders();
    config.headers = { ...config.headers, ...authHeaders };
    return config;
  },
  (error) => Promise.reject(error)
);

// Interceptor para manejar errores HTTP (sin lógica de negocio)
apiClient.interceptors.response.use(
  (response) => {
    // En desarrollo, logear todas las requests exitosas
    const isDevelopment = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    if (isDevelopment) {
      console.log(`✅ [apiClient] ${response.config.method?.toUpperCase()} ${response.config.url} - ${response.status}`);
    }
    return response;
  },
  (error) => {
    const isDevelopment = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    
    // Logging extendido en desarrollo
    if (isDevelopment) {
      console.error(`❌ [apiClient] ${error.config?.method?.toUpperCase()} ${error.config?.url}:`, {
        status: error.response?.status,
        statusText: error.response?.statusText,
        message: error.response?.data?.error || error.message,
        headers: error.response?.headers,
        isRateLimit: error.response?.status === 429
      });
      
      // Alerta específica para rate limiting
      if (error.response?.status === 429) {
        console.warn(`🚫 [apiClient] RATE LIMIT detectado en ${error.config?.url}. Headers de desarrollo agregados.`);
      }
    } else {
      // Solo logging básico en producción
      console.error('API Error:', {
        url: error.config?.url,
        status: error.response?.status,
        message: error.response?.data?.error || error.message,
      });
    }
    
    return Promise.reject(error);
  }
);

// Cliente HTTP puro - solo comunicación
export default apiClient; 